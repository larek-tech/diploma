package parse_site_status

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
)

type (
	publisher interface {
		Publish(ctx context.Context, rawMsg []any, opts ...qaas.PublishOption) ([]string, error)
	}
	pageJobStore interface {
		GetProcessedPageCount(ctx context.Context, parseSiteJobID string) (int, error)
		GetUnprocessedPageCount(ctx context.Context, parseSiteJobID string) (int, error)
	}
	// TODO: add kafka producer for status messages
	kafkaProducer interface {
		Produce(ctx context.Context, msg *kafka.Message) error
	}
)
