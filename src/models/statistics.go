package models

type StatisticsType string

const (
	ImageFormatType   StatisticsType = "ImageFormatType"
	CameraModelType   StatisticsType = "CameraModelType"
	DateFrequencyType StatisticsType = "DateFrequencyType"
)

type Statistics struct {
	ID    string         `json:"-" bson:"-"`
	Type  StatisticsType `json:"-" bson:"type"`
	Name  string         `json:"name" bson:"name"`
	Count int            `json:"count" bson:"count"`
}
