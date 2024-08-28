package models

import "time"

type Image struct {
	ID           string    `json:"id" bson:"-"`
	Name         string    `json:"name" bson:"name"`
	UploadLinkID string    `json:"uploadLinkID" bson:"upload_link_id"`
	CameraModel  string    `json:"cameraModel" bson:"camera_model"`
	ImageFormat  string    `json:"imageFormat" bson:"image_format"`
	ImageWidth   int       `json:"imageWidth" bson:"image_width"`
	ImageHeight  int       `json:"imageHeight" bson:"image_height"`
	Path         string    `json:"path" bson:"path"`
	Latitude     float64   `json:"latitude" bson:"latitude"`
	Longitude    float64   `json:"longitude" bson:"longitude"`
	UploadedAt   time.Time `json:"uploadTime" bson:"upload_time"`
}
