package main

import (
	"fmt"
	"net/http"

	"github.com/tam-code/lrn/config"
	"github.com/tam-code/lrn/src/consumers"
	"github.com/tam-code/lrn/src/databases"
	"github.com/tam-code/lrn/src/kafka"
	"github.com/tam-code/lrn/src/producers"
	"github.com/tam-code/lrn/src/repositories"
	"github.com/tam-code/lrn/src/routes"
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
