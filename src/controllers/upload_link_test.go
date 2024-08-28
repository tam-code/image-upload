package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/tam-code/lrn/mocks/repositories"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/tam-code/lrn/src/models"
)

func TestCreateUploadLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUploadLinkRepository(ctrl)
	controller := &uploadLinkController{
		uploadLinkRepo: mockRepo,
		uploadLinkPath: "/api/image",
	}

	tests := []struct {
		name           string
		expiration     string
		mockRepoFunc   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing expiration",
			expiration:     "",
			mockRepoFunc:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Expiration is required",
		},
		{
			name:           "invalid expiration format",
			expiration:     "invalid-date",
			mockRepoFunc:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid expiration, it must be ISO8601 format e.g. 2007-10-09T22:50:01.23Z",
		},
		{
			name:           "expiration in the past",
			expiration:     "2006-01-02T15:04:05.999Z",
			mockRepoFunc:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Expiration must be in the future",
		},
		{
			name:       "successful creation",
			expiration: "2106-01-02T15:04:05.999Z",
			mockRepoFunc: func() {
				mockRepo.EXPECT().CreateUploadLink(gomock.Any()).Return(&models.UploadLink{ID: "tested-link"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "tested-link",
		},
		{
			name:       "unsuccessful creation",
			expiration: "2106-01-02T15:04:05.999Z",
			mockRepoFunc: func() {
				mockRepo.EXPECT().CreateUploadLink(gomock.Any()).Return(nil, errors.New("error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFunc()

			reqBody := bytes.NewBufferString("expiration=" + tt.expiration)
			req := httptest.NewRequest(http.MethodPost, "/upload-link", reqBody)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			controller.CreateUploadLink(w, req)

			resp := w.Result()
			body, _ := json.Marshal(w.Body.String())

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}
