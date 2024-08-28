package handlers

import (
	"encoding/json"
	"log"

	"github.com/tam-code/lrn/src/models"
	"github.com/tam-code/lrn/src/repositories"
)

type (
	ImageUploadedHandler interface {
		Handle([]byte)
	}

	imageUploadedHandler struct {
		imageRepository      repositories.ImageRepository
		statisticsRepository repositories.StatisticsRepository
	}
)

func NewImageUploadedHandler(repositories *repositories.Repositories) ImageUploadedHandler {
	return &imageUploadedHandler{
		imageRepository:      repositories.Image,
		statisticsRepository: repositories.Statistics,
	}
}

func (h *imageUploadedHandler) Handle(message []byte) {
	var images []string
	if err := json.Unmarshal(message, &images); err != nil {
		log.Fatalf("error unmarshalling message: %v", err)
	}

	// Do something with images
	imagesObjects, err := h.imageRepository.GetImagesByIDs(images)
	if err != nil {
		log.Fatalf("error getting images by ids: %v", err)
	}

	imageFormats := make(map[string]int)
	cameraModels := make(map[string]int)
	uploadedImagesPerDay := make(map[string]int)
	for _, image := range imagesObjects {
		if image.ImageFormat != "" {
			imageFormats[image.ImageFormat]++
		}

		if image.CameraModel != "" {
			cameraModels[image.CameraModel]++
		}

		uploadedImagesPerDay[image.UploadedAt.Format("2006-01-02")]++
	}

	h.updateCameraModelsCount(cameraModels)
	h.updateImageFormatsCount(imageFormats)
	h.updateImagesDailyUploadedCount(uploadedImagesPerDay)

}

func (h *imageUploadedHandler) updateCameraModelsCount(cameraModels map[string]int) {
	// Do something with camera models
	for model, count := range cameraModels {
		cameraModel, err := h.statisticsRepository.GetStatistics(models.CameraModelType, model)
		if err != nil {
			log.Printf("error getting camera model by name: %v", err)
		}

		if cameraModel == nil {
			cameraModel = &models.Statistics{
				Type:  models.CameraModelType,
				Name:  model,
				Count: count,
			}

			if err = h.statisticsRepository.InsertStatistics(cameraModel); err != nil {
				log.Printf("error creating camera model: %v", err)
			}
			continue
		}

		cameraModel.Count += count
		if err := h.statisticsRepository.UpdateStatistics(cameraModel); err != nil {
			log.Printf("error updating camera model: %v", err)
		}
	}
}

func (h *imageUploadedHandler) updateImageFormatsCount(imageFormats map[string]int) {
	// Do something with image formats
	for format, count := range imageFormats {
		imageFormat, err := h.statisticsRepository.GetStatistics(models.ImageFormatType, format)
		if err != nil {
			log.Printf("error getting image format by id: %v", err)
		}

		if imageFormat == nil {
			imageFormat = &models.Statistics{
				Type:  models.ImageFormatType,
				Name:  format,
				Count: count,
			}

			if err = h.statisticsRepository.InsertStatistics(imageFormat); err != nil {
				log.Printf("error creating image format: %v", err)
			}
			continue
		}

		imageFormat.Count += count
		if err := h.statisticsRepository.UpdateStatistics(imageFormat); err != nil {
			log.Printf("error updating image format: %v", err)
		}
	}
}

func (h *imageUploadedHandler) updateImagesDailyUploadedCount(uploadedImagesPerDay map[string]int) {
	// Do something with images daily uploaded
	for day, count := range uploadedImagesPerDay {
		imagesDailyUploaded, err := h.statisticsRepository.GetStatistics(models.DateFrequencyType, day)
		if err != nil {
			log.Printf("error getting images daily uploaded by day: %v", err)
		}

		if imagesDailyUploaded == nil {
			imagesDailyUploaded = &models.Statistics{
				Type:  models.DateFrequencyType,
				Name:  day,
				Count: count,
			}

			if err = h.statisticsRepository.InsertStatistics(imagesDailyUploaded); err != nil {
				log.Printf("error creating images daily uploaded: %v", err)
			}
			continue
		}

		imagesDailyUploaded.Count += count
		if err := h.statisticsRepository.UpdateStatistics(imagesDailyUploaded); err != nil {
			log.Printf("error updating images daily uploaded: %v", err)
		}
	}
}
