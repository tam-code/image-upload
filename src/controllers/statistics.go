package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/tam-code/image-upload/src/models"
	"github.com/tam-code/image-upload/src/repositories"
)

type (
	StatisticsController interface {
		GetStatistics(w http.ResponseWriter, r *http.Request)
	}

	statisticsController struct {
		statisticsRepository repositories.StatisticsRepository
	}

	statistics struct {
		MostPopularImageFormat  []models.Statistics `json:"mostPopularImageFormat"`
		MostPopularCameraModels []models.Statistics `json:"mostPopularCameraModels"`
		UploadFrequencyPerDay   []models.Statistics `json:"uploadFrequencyPerDay"`
	}
)

func NewStatisticsController(repositories *repositories.Repositories) StatisticsController {
	return &statisticsController{
		statisticsRepository: repositories.Statistics,
	}
}

func (c *statisticsController) GetStatistics(w http.ResponseWriter, r *http.Request) {
	mostPopularImageFormat, err := c.statisticsRepository.GetStatisticsSortedByCount(models.ImageFormatType, 1)
	if err != nil {
		http.Error(w, "Error getting most popular image format", http.StatusInternalServerError)
		return
	}

	mostPopularCameraModels, err := c.statisticsRepository.GetStatisticsSortedByCount(models.CameraModelType, 10)
	if err != nil {
		http.Error(w, "Error getting most popular camera models", http.StatusInternalServerError)
		return
	}

	uploadFrequencyPerDay, err := c.statisticsRepository.GetStatisticsFrequency(models.DateFrequencyType, 30)
	if err != nil {
		http.Error(w, "Error getting upload frequency per day", http.StatusInternalServerError)
		return
	}

	statistics := statistics{
		MostPopularCameraModels: mostPopularCameraModels,
		UploadFrequencyPerDay:   uploadFrequencyPerDay,
		MostPopularImageFormat:  mostPopularImageFormat,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statistics)
}
