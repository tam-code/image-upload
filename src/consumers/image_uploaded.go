package consumers

import (
	"context"
	"log"

	"github.com/tam-code/image-upload/src/handlers"
	"github.com/tam-code/image-upload/src/kafka"
	"github.com/tam-code/image-upload/src/repositories"
)

type (
	ImageUploadedConsumer interface {
		Consume(context.Context)
	}

	imageUploadedConsumer struct {
		consumer             kafka.Consumer
		imageUploadedHandler handlers.ImageUploadedHandler
	}
)

func NewImageUploadedConsumer(c kafka.Consumer, repositories *repositories.Repositories) ImageUploadedConsumer {
	return &imageUploadedConsumer{
		consumer:             c,
		imageUploadedHandler: handlers.NewImageUploadedHandler(repositories),
	}
}

func (c *imageUploadedConsumer) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Fatalf("consumer canceled", ctx.Err())
			return
		default:
			// Fetch message from kafka
			msg, err := c.consumer.FetchMessage(context.Background())
			if err != nil {
				log.Fatalf("error fetching message: %v", err)
				return
			}

			c.imageUploadedHandler.Handle(msg.Value)

			// Commit message
			if err := c.consumer.CommitMessages(context.Background(), msg); err != nil {
				log.Fatalf("error committing message: %v", err)
			}
		}
	}
}
