package repositories

import (
	"context"

	"github.com/tam-code/lrn/src/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	StatisticsRepository interface {
		InsertStatistics(statistics *models.Statistics) error
		GetStatistics(statisticsType models.StatisticsType, name string) (*models.Statistics, error)
		UpdateStatistics(statistics *models.Statistics) error
		GetStatisticsFrequency(statisticsType models.StatisticsType, limit int) ([]models.Statistics, error)
		GetStatisticsSortedByCount(statisticsType models.StatisticsType, limit int) ([]models.Statistics, error)
	}

	statisticsRepository struct {
		mongoCollection *mongo.Collection
	}
)

func newStatisticsRepository(mongoDB mongo.Database) StatisticsRepository {
	return &statisticsRepository{
		mongoCollection: mongoDB.Collection("statistics"),
	}
}

func (r *statisticsRepository) InsertStatistics(statistics *models.Statistics) error {
	_, err := r.mongoCollection.InsertOne(context.Background(), statistics)
	if err != nil {
		return err
	}

	return nil
}

func (r *statisticsRepository) GetStatistics(statisticsType models.StatisticsType, name string) (*models.Statistics, error) {
	var statistics models.Statistics
	err := r.mongoCollection.FindOne(context.Background(), map[string]interface{}{"type": statisticsType, "name": name}).Decode(&statistics)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &statistics, nil
}

func (r *statisticsRepository) UpdateStatistics(statistics *models.Statistics) error {
	_, err := r.mongoCollection.UpdateOne(context.Background(), map[string]interface{}{"type": statistics.Type, "name": statistics.Name}, map[string]interface{}{"$set": statistics})
	if err != nil {
		return err
	}

	return nil
}

func (r *statisticsRepository) GetStatisticsFrequency(statisticsType models.StatisticsType, limit int) ([]models.Statistics, error) {
	var statistics []models.Statistics
	limit64 := int64(limit)
	pipeline := mongo.Pipeline{
		{{"$match", map[string]interface{}{"type": statisticsType}}},
		{{"$sort", map[string]interface{}{"name": -1}}},
		{{"$limit", limit64}},
	}
	cursor, err := r.mongoCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &statistics); err != nil {
		return nil, err
	}

	return statistics, nil
}

func (r *statisticsRepository) GetStatisticsSortedByCount(statisticsType models.StatisticsType, limit int) ([]models.Statistics, error) {
	var statistics []models.Statistics
	limit64 := int64(limit)
	pipeline := mongo.Pipeline{
		{{"$match", map[string]interface{}{"type": statisticsType}}},
		{{"$sort", map[string]interface{}{"count": -1}}},
		{{"$limit", limit64}},
	}
	cursor, err := r.mongoCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &statistics); err != nil {
		return nil, err
	}

	return statistics, nil
}
