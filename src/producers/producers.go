package producers

import "github.com/tam-code/image-upload/src/kafka"

type (
	Producers struct {
		ImageUploaded ImageUploadedProducer
	}
)

func NewProducers(writer kafka.Producer) *Producers {
	return &Producers{
		ImageUploaded: NewImageUploadedProducer(writer),
	}
}
