package repositories

import (
	"context"
	"fmt"

	"github.com/tam-code/lrn/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ImageRepository interface {
		InsertImages([]interface{}) ([]string, error)
		GetImageByID(string) (*models.Image, error)
		GetImageByName(string) (*models.Image, error)
		GetImagesByIDs([]string) ([]models.Image, error)
		UpdateImage(*models.Image) error
		GetImageByNameAndUploadLinkID(string, string) (*models.Image, error)
	}

	imageRepository struct {
		mogoCollection *mongo.Collection
	}
)

func newImageRepository(mongoDB mongo.Database) ImageRepository {
	return &imageRepository{
		mogoCollection: mongoDB.Collection("images"),
	}
}

func (r *imageRepository) InsertImages(images []interface{}) ([]string, error) {

	insertedData, err := r.mogoCollection.InsertMany(context.Background(), images)
	if err != nil {
		return nil, err
	}

	var insertedImages []string
	for _, id := range insertedData.InsertedIDs {
		insertedImages = append(insertedImages, id.(primitive.ObjectID).Hex())
	}

	return insertedImages, nil
}

func (r *imageRepository) GetImageByID(id string) (*models.Image, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error converting id to object id: %w", err)
	}

	var image models.Image
	err = r.mogoCollection.FindOne(context.Background(), primitive.M{"_id": objectID}).Decode(&image)
	if err != nil {
		return nil, fmt.Errorf("error getting image by id: %w", err)
	}

	image.ID = objectID.Hex()

	return &image, nil
}

func (r *imageRepository) GetImageByName(name string) (*models.Image, error) {
	var image models.Image
	err := r.mogoCollection.FindOne(context.Background(), primitive.M{"name": name}).Decode(&image)
	if err != nil {
		return nil, fmt.Errorf("error getting image by name: %w", err)
	}

	return &image, nil
}

func (r *imageRepository) GetImagesByIDs(ids []string) ([]models.Image, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("error converting id to object id: %w", err)
		}

		objectIDs = append(objectIDs, objectID)
	}

	cursor, err := r.mogoCollection.Find(context.Background(), primitive.M{"_id": primitive.M{"$in": objectIDs}})
	if err != nil {
		return nil, fmt.Errorf("error getting images by ids: %w", err)
	}

	var images []models.Image
	err = cursor.All(context.Background(), &images)
	if err != nil {
		return nil, fmt.Errorf("error getting images by ids: %w", err)
	}

	return images, nil
}

func (r *imageRepository) UpdateImage(image *models.Image) error {
	objectID, err := primitive.ObjectIDFromHex(image.ID)
	if err != nil {
		return fmt.Errorf("error converting id to object id: %w", err)
	}

	_, err = r.mogoCollection.UpdateOne(context.Background(), primitive.M{"_id": objectID}, primitive.M{"$set": image})
	if err != nil {
		return fmt.Errorf("error updating image: %w", err)
	}

	return nil
}

func (r *imageRepository) GetImageByNameAndUploadLinkID(name, uploadLinkID string) (*models.Image, error) {
	var image models.Image
	err := r.mogoCollection.FindOne(context.Background(), primitive.M{"name": name, "upload_link_id": uploadLinkID}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting image by name and upload link id: %w", err)
	}

	return &image, nil
}
