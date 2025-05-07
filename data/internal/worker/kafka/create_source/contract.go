package create_source

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/larek-tech/diploma/data/internal/domain/source"
)

type (
	service interface {
		CreateSource(ctx context.Context, message source.DataMessage) (*source.Source, error)
	}
	kafkaProducer interface {
		Produce(context.Context, *kafka.Message) error
	}
)
