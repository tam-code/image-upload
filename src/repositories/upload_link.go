package repositories

import (
	"context"
	"fmt"

	"github.com/tam-code/lrn/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UploadLinkRepository interface {
		CreateUploadLink(models.UploadLink) (*models.UploadLink, error)
		GetUploadLinkByID(string) (*models.UploadLink, error)
	}

	uploadLinkRepository struct {
		mongoCollection *mongo.Collection
	}
)

func newUploadLinkRepository(mongodb mongo.Database) UploadLinkRepository {
	return &uploadLinkRepository{
		mongoCollection: mongodb.Collection("upload_links"),
	}
}

func (r *uploadLinkRepository) CreateUploadLink(uploadLink models.UploadLink) (*models.UploadLink, error) {
	insertedData, err := r.mongoCollection.InsertOne(context.Background(), uploadLink)
	if err != nil {
		return nil, fmt.Errorf("error inserting status update log: %w", err)
	}

	uploadLink.ID = insertedData.InsertedID.(primitive.ObjectID).Hex()

	return &uploadLink, nil
}

func (r *uploadLinkRepository) GetUploadLinkByID(id string) (*models.UploadLink, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error converting id to object id: %w", err)
	}

	var uploadLink models.UploadLink
	err = r.mongoCollection.FindOne(context.Background(), primitive.M{"_id": objectID}).Decode(&uploadLink)
	if err != nil {
		return nil, fmt.Errorf("error getting upload link by id: %w", err)
	}

	return &uploadLink, nil
}
