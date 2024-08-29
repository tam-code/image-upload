package kafka

import (
	"strings"
	"time"

	segmentio "github.com/segmentio/kafka-go"
	"github.com/tam-code/image-upload/config"
)

type Message = segmentio.Message
type Header = segmentio.Header

func NewKafkaWriter(cfg config.KafkaConfig) *segmentio.Writer {
	writerConfig := segmentio.WriterConfig{
		Brokers:  strings.Split(cfg.Brokers, ","),
		Topic:    cfg.Topic,
		Balancer: &segmentio.LeastBytes{},
	}

	return segmentio.NewWriter(writerConfig)
}

func NewKafkaReader(cfg config.KafkaConfig) *segmentio.Reader {
	readerConfig := segmentio.ReaderConfig{
		Brokers:  strings.Split(cfg.Brokers, ","),
		GroupID:  cfg.Group,
		Topic:    cfg.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  time.Duration(cfg.MaxWaitTimeoutMilliseconds) * time.Millisecond,
	}

	return segmentio.NewReader(readerConfig)
}
