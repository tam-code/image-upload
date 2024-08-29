package main

import (
	"fmt"
	"net/http"

	"github.com/tam-code/image-upload/config"
	"github.com/tam-code/image-upload/src/consumers"
	"github.com/tam-code/image-upload/src/databases"
	"github.com/tam-code/image-upload/src/kafka"
	"github.com/tam-code/image-upload/src/producers"
	"github.com/tam-code/image-upload/src/repositories"
	"github.com/tam-code/image-upload/src/routes"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	mongodb, err := databases.NewMongoDB(config.MongoDB)
	if err != nil {
		panic(err)
	}

	repositories := repositories.NewRepositories(mongodb)

	consumers := consumers.NewConsumers(kafka.NewConsumer(kafka.NewKafkaReader(config.Kafka)), repositories)
	consumers.Run()

	producers := producers.NewProducers(kafka.NewProducer(kafka.NewKafkaWriter(config.Kafka)))

	http.ListenAndServe(fmt.Sprintf(":%v", config.APIPort), routes.SetupRoutes(repositories, producers))
}
