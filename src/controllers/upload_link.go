package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tam-code/lrn/src/models"
	"github.com/tam-code/lrn/src/repositories"
)

type (
	UploadLinkController interface {
		CreateUploadLink(w http.ResponseWriter, r *http.Request)
	}

	uploadLinkController struct {
		uploadLinkRepo repositories.UploadLinkRepository
		uploadLinkPath string
	}
)

func NewUploadLinkController(repositories *repositories.Repositories, UploadLinkPath string) UploadLinkController {
	return &uploadLinkController{
		uploadLinkRepo: repositories.UploadLink,
		uploadLinkPath: UploadLinkPath,
	}
}

func (c *uploadLinkController) CreateUploadLink(w http.ResponseWriter, r *http.Request) {
	// validate request body
	expiration := r.FormValue("expiration")
	if expiration == "" {
		http.Error(w, "Expiration is required", http.StatusBadRequest)
		return
	}

	expirationTime, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		http.Error(w, "Invalid expiration, it must be ISO8601 format e.g. 2007-10-09T22:50:01.23Z", http.StatusBadRequest)
		return
	}

	if expirationTime.Before(time.Now()) {
		http.Error(w, "Expiration must be in the future", http.StatusBadRequest)
		return
	}

	// create upload link
	uploadLink, err := c.uploadLinkRepo.CreateUploadLink(models.UploadLink{
		ExpirationTime: expirationTime,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return upload link
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Host + c.uploadLinkPath + "/" + uploadLink.ID)
}
