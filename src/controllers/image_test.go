package controllers

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mocksProducer "github.com/tam-code/lrn/mocks/producers"
	mocks "github.com/tam-code/lrn/mocks/repositories"
	"github.com/tam-code/lrn/src/models"
)

func TestUploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer cleanUploadTestFolder()

	mockUploadLinkRepo := mocks.NewMockUploadLinkRepository(ctrl)
	mockImageRepo := mocks.NewMockImageRepository(ctrl)
	mockImageUploadedProducer := mocksProducer.NewMockImageUploadedProducer(ctrl)

	controller := &imageController{
		uploadLinkRepo:        mockUploadLinkRepo,
		imageRepo:             mockImageRepo,
		imageUploadedProducer: mockImageUploadedProducer,
	}

	tests := []struct {
		name           string
		uploadLinkID   string
		uploadLink     *models.UploadLink
		mockRepoFunc   func()
		formData       map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "invalid upload link",
			uploadLinkID: "invalid",
			mockRepoFunc: func() {
				mockUploadLinkRepo.EXPECT().GetUploadLinkByID("invalid").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Invalid upload link or not found\n",
		},
		{
			name:         "expired upload link",
			uploadLinkID: "expired",
			uploadLink: &models.UploadLink{
				ExpirationTime: time.Now().Add(-time.Hour),
			},
			mockRepoFunc: func() {
				mockUploadLinkRepo.EXPECT().GetUploadLinkByID("expired").Return(&models.UploadLink{
					ExpirationTime: time.Now().Add(-time.Hour),
				}, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Upload link expired\n",
		},
		{
			name:         "error parsing form",
			uploadLinkID: "valid",
			uploadLink: &models.UploadLink{
				ExpirationTime: time.Now().Add(time.Hour),
			},
			mockRepoFunc: func() {
				mockUploadLinkRepo.EXPECT().GetUploadLinkByID("valid").Return(&models.UploadLink{
					ExpirationTime: time.Now().Add(time.Hour),
				}, nil)
			},
			formData:       map[string]string{"invalid": "data"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No images found upload\n",
		},
		{
			name:         "successful upload",
			uploadLinkID: "valid",
			uploadLink: &models.UploadLink{
				ExpirationTime: time.Now().Add(time.Hour),
			},
			mockRepoFunc: func() {
				mockUploadLinkRepo.EXPECT().GetUploadLinkByID("valid").Return(&models.UploadLink{
					ExpirationTime: time.Now().Add(time.Hour),
				}, nil)
				mockImageRepo.EXPECT().InsertImages(gomock.Any()).Return([]string{"image1.jpg"}, nil)
				mockImageRepo.EXPECT().GetImageByNameAndUploadLinkID(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockImageUploadedProducer.EXPECT().Publish(gomock.Any()).Return(nil)
			},
			formData:       map[string]string{"images": "image1.jpg"},
			expectedStatus: http.StatusOK,
			expectedBody:   `["image1.jpg"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFunc()

			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			for key, val := range tt.formData {
				filePart, _ := writer.CreateFormFile(key, val)
				filePart.Write([]byte("Hello, World!"))
			}
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/upload-image", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req = mux.SetURLVars(req, map[string]string{"upload_link_id": tt.uploadLinkID})
			w := httptest.NewRecorder()

			controller.UploadImage(w, req)

			resp := w.Result()
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := strings.TrimSpace(string(bodyBytes))

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, tt.expectedBody, bodyString)
		})
	}
}

func cleanUploadTestFolder() {
	os.RemoveAll(uploadPath)
}

func TestGetImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := mocks.NewMockImageRepository(ctrl)

	controller := &imageController{
		imageRepo: mockImageRepo,
	}

	tests := []struct {
		name           string
		imageID        string
		mockRepoFunc   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "image not found",
			imageID: "invalid",
			mockRepoFunc: func() {
				mockImageRepo.EXPECT().GetImageByID("invalid").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Invalid image id or not found",
		},
		{
			name:    "successful retrieval",
			imageID: "valid",
			mockRepoFunc: func() {
				mockImageRepo.EXPECT().GetImageByID("valid").Return(&models.Image{
					ID:   "valid",
					Name: "test_image",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":"valid","name":"test_image"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFunc()

			req := httptest.NewRequest(http.MethodGet, "/image/{image_id}", nil)
			req = mux.SetURLVars(req, map[string]string{"image_id": tt.imageID})
			w := httptest.NewRecorder()

			controller.GetImage(w, req)

			resp := w.Result()
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := strings.TrimSpace(string(bodyBytes))

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, bodyString, tt.expectedBody)
		})
	}
}
