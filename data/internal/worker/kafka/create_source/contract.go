package create_source

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/source"
)

type (
	service interface {
		CreateSource(ctx context.Context, message source.DataMessage) (*source.Source, error)
	}
	kafkaProducer interface {
		Produce(ctx context.Context, topic string, key []byte, value []byte) error
	}
)
