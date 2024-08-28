package models

import "time"

type UploadLink struct {
	ID             string    `json:"id" bson:"-"`
	ExpirationTime time.Time `json:"expirationTime" bson:"expiration_time"`
}
