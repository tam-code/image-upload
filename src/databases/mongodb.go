package databases

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tam-code/image-upload/config"
)

func NewMongoDB(cfg config.MongoDBConfig) (*mongo.Database, error) {
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI()))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	return mongoClient.Database(cfg.Database), nil
}
