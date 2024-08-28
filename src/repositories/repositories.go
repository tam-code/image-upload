package repositories

import "go.mongodb.org/mongo-driver/mongo"

type Repositories struct {
	UploadLink UploadLinkRepository
	Image      ImageRepository
	Statistics StatisticsRepository
}

func NewRepositories(mongodb *mongo.Database) *Repositories {
	return &Repositories{
		UploadLink: newUploadLinkRepository(*mongodb),
		Image:      newImageRepository(*mongodb),
		Statistics: newStatisticsRepository(*mongodb),
	}
}
