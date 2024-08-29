package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/gorilla/mux"

	"github.com/tam-code/image-upload/src/models"
	"github.com/tam-code/image-upload/src/producers"
	"github.com/tam-code/image-upload/src/repositories"
)

const (
	uploadPath = "./upload/"
)

type (
	ImageController interface {
		UploadImage(w http.ResponseWriter, r *http.Request)
		GetImage(w http.ResponseWriter, r *http.Request)
	}

	imageController struct {
		uploadLinkRepo        repositories.UploadLinkRepository
		imageRepo             repositories.ImageRepository
		imageUploadedProducer producers.ImageUploadedProducer
	}
)

func NewImageController(repositories *repositories.Repositories, producers *producers.Producers) ImageController {
	return &imageController{
		uploadLinkRepo:        repositories.UploadLink,
		imageRepo:             repositories.Image,
		imageUploadedProducer: producers.ImageUploaded,
	}
}

func (c *imageController) UploadImage(w http.ResponseWriter, r *http.Request) {

	uploadLinkId := mux.Vars(r)["upload_link_id"]
	uploadLink, err := c.uploadLinkRepo.GetUploadLinkByID(uploadLinkId)
	if err != nil {
		http.Error(w, "Invalid upload link or not found", http.StatusNotFound)
		return
	}

	if uploadLink.ExpirationTime.Before(time.Now()) {
		http.Error(w, "Upload link expired", http.StatusForbidden)
		return
	}

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing form, "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "No images found upload", http.StatusBadRequest)
		return
	}

	var images []interface{}
	imagesMap := make(map[string]int)
	for _, file := range files {
		// check if file already uploaded, incase of multiple files with same name
		_, ok := imagesMap[file.Filename]
		if ok {
			continue
		}

		imagesMap[file.Filename] = 1

		image, err := c.handleFileUpload(file, uploadLinkId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if image != nil {
			images = append(images, image)
		}
	}

	insertedImages, err := c.imageRepo.InsertImages(images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(insertedImages) == 0 {
		http.Error(w, "No images uploaded", http.StatusBadRequest)
		return
	}

	// publish images uploaded event
	err = c.imageUploadedProducer.Publish(insertedImages)
	if err != nil {
		log.Printf("error publishing images uploaded event: %v", err)
	}

	// return inserted images ids
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insertedImages)
}

func (c *imageController) handleFileUpload(file *multipart.FileHeader, uploadLinkId string) (*models.Image, error) {
	// validate file
	err := validateImage(file)
	if err != nil {
		return nil, err
	}

	// check duplicate image
	imageExist, err := c.imageRepo.GetImageByNameAndUploadLinkID(file.Filename, uploadLinkId)
	if err != nil {
		return nil, fmt.Errorf("error getting image by name: %w", err)
	}

	if imageExist != nil {
		log.Printf("image already uploaded: %s", file.Filename)
		return nil, nil
	}

	// upload file
	loc, err := uploadImageSource(file, uploadLinkId)
	if err != nil {
		return nil, err
	}

	// create image model
	image := models.Image{
		Name:         file.Filename,
		Path:         loc,
		UploadLinkID: uploadLinkId,
		UploadedAt:   time.Now(),
	}

	adaptImageMetadata(&image)

	return &image, nil
}

func (c *imageController) GetImage(w http.ResponseWriter, r *http.Request) {
	imageID := mux.Vars(r)["image_id"]
	image, err := c.imageRepo.GetImageByID(imageID)
	if err != nil {
		http.Error(w, "Invalid image id or not found", http.StatusNotFound)
		return
	}

	// return inserted images ids
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}

func uploadImageSource(file *multipart.FileHeader, uploadLinkID string) (string, error) {
	loc := filepath.Join(uploadPath, uploadLinkID)
	err := os.MkdirAll(loc, os.ModePerm)
	if err != nil {
		return loc, err
	}

	loc = filepath.Join(loc, "/", file.Filename)
	dst, err := os.Create(loc)
	if err != nil {
		return loc, err
	}
	defer dst.Close()

	f, err := file.Open()
	if err != nil {
		return loc, err
	}
	defer f.Close()

	if _, err := io.Copy(dst, f); err != nil {
		return loc, err
	}

	return loc, nil
}

func validateImage(file *multipart.FileHeader) error {
	// validate file size
	if file.Size > 10<<20 {
		return fmt.Errorf("file size exceeds 10MB")
	}

	// validate file type
	ext := filepath.Ext(file.Filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" || !strings.HasPrefix(mimeType, "image") {
		return fmt.Errorf("invalid file type: %s", mimeType)
	}

	return nil
}

func adaptImageMetadata(image *models.Image) {
	// open file
	f, err := os.Open("./" + image.Path)
	if err != nil {
		log.Printf("error opening image: %v", err)
		return
	}
	defer f.Close()

	// decode image
	e, err := imagemeta.Decode(f)
	if err != nil {
		log.Printf("error decoding image: %v", err)
	}

	// update image metadata
	image.ImageWidth = int(e.ImageWidth)
	image.ImageHeight = int(e.ImageHeight)

	if e.GPS.Latitude() != 0 && e.GPS.Longitude() != 0 {
		image.Latitude = e.GPS.Latitude()
		image.Longitude = e.GPS.Longitude()
	}

	if e.Model != "" {
		image.CameraModel = e.Model
	}

	if e.ImageType.String() != "" {
		image.ImageFormat = e.ImageType.String()
	}
}
