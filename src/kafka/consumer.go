package kafka

import (
	"context"

	segmentio "github.com/segmentio/kafka-go"
)

type (
	Consumer interface {
		FetchMessage(ctx context.Context) (Message, error)
		CommitMessages(ctx context.Context, msgs ...Message) error
		Close() error
	}
	consumerWrapper struct {
		reader *segmentio.Reader
	}
)

func NewConsumer(r *segmentio.Reader) Consumer {
	return &consumerWrapper{r}
}

func (c *consumerWrapper) FetchMessage(ctx context.Context) (Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *consumerWrapper) CommitMessages(ctx context.Context, msgs ...Message) error {
	return c.reader.CommitMessages(ctx, msgs...)
}

func (c *consumerWrapper) Close() error {
	return c.reader.Close()
}
