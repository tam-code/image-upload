package producers

import (
	"context"
	"encoding/json"

	"github.com/tam-code/lrn/src/kafka"
)

type (
	ImageUploadedProducer interface {
		Publish(images []string) error
	}

	imageUploadedProducer struct {
		producer kafka.Producer
	}
)

func NewImageUploadedProducer(p kafka.Producer) ImageUploadedProducer {
	return &imageUploadedProducer{p}
}

func (p *imageUploadedProducer) Publish(images []string) error {
	jsonImage, err := json.Marshal(images)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Value: jsonImage,
	}

	return p.producer.WriteMessages(context.Background(), message)
}
