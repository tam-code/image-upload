package kafka

import (
	"context"

	segmentio "github.com/segmentio/kafka-go"
)

type (
	Producer interface {
		WriteMessages(ctx context.Context, msgs ...Message) error
		Close() error
	}
	producerWrapper struct {
		writer *segmentio.Writer
	}
)

func NewProducer(w *segmentio.Writer) Producer {
	return &producerWrapper{w}
}

func (p *producerWrapper) WriteMessages(ctx context.Context, msgs ...Message) error {
	return p.writer.WriteMessages(ctx, msgs...)
}

func (p *producerWrapper) Close() error {
	return p.writer.Close()
}
