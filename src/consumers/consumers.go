package consumers

import (
	"context"

	"github.com/tam-code/lrn/src/kafka"
	"github.com/tam-code/lrn/src/repositories"
)

type Consumers struct {
	imageUploadedConsumer ImageUploadedConsumer
}

func NewConsumers(c kafka.Consumer, repositories *repositories.Repositories) *Consumers {
	return &Consumers{
		imageUploadedConsumer: NewImageUploadedConsumer(c, repositories),
	}
}

func (c *Consumers) Run() {
	go c.imageUploadedConsumer.Consume(context.Background())
}
