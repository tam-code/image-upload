package controllers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/tam-code/lrn/mocks/repositories"
	"github.com/tam-code/lrn/src/models"
)

func TestGetStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatisticsRepo := mocks.NewMockStatisticsRepository(ctrl)

	controller := &statisticsController{
		statisticsRepository: mockStatisticsRepo,
	}

	tests := []struct {
		name           string
		mockRepoFunc   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "error getting most popular image format",
			mockRepoFunc: func() {
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.ImageFormatType, 1).Return(nil, errors.New("error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Error getting most popular image format",
		},
		{
			name: "error getting most popular camera models",
			mockRepoFunc: func() {
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.ImageFormatType, 1).Return([]models.Statistics{{Name: "JPEG", Count: 100}}, nil)
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.CameraModelType, 10).Return(nil, errors.New("error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Error getting most popular camera models",
		},
		{
			name: "error getting upload frequency per day",
			mockRepoFunc: func() {
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.ImageFormatType, 1).Return([]models.Statistics{{Name: "JPEG", Count: 100}}, nil)
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.CameraModelType, 10).Return([]models.Statistics{{Name: "Canon", Count: 50}}, nil)
				mockStatisticsRepo.EXPECT().GetStatisticsFrequency(models.DateFrequencyType, 30).Return(nil, errors.New("error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Error getting upload frequency per day",
		},
		{
			name: "successful retrieval",
			mockRepoFunc: func() {
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.ImageFormatType, 1).Return([]models.Statistics{{Name: "JPEG", Count: 100}}, nil)
				mockStatisticsRepo.EXPECT().GetStatisticsSortedByCount(models.CameraModelType, 10).Return([]models.Statistics{{Name: "Canon", Count: 50}}, nil)
				mockStatisticsRepo.EXPECT().GetStatisticsFrequency(models.DateFrequencyType, 30).Return([]models.Statistics{{Name: "2023-10-01", Count: 10}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"mostPopularImageFormat":[{"name":"JPEG","count":100}],"mostPopularCameraModels":[{"name":"Canon","count":50}],"uploadFrequencyPerDay":[{"name":"2023-10-01","count":10}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFunc()

			req := httptest.NewRequest(http.MethodGet, "/statistics", nil)
			w := httptest.NewRecorder()

			controller.GetStatistics(w, req)

			resp := w.Result()
			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyString := strings.TrimSpace(string(bodyBytes))

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, bodyString, tt.expectedBody)
		})
	}
}
