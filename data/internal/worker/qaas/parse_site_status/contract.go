package parse_site_status

import (
	"context"

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

	kafkaProducer interface {
		Produce(ctx context.Context, topic string, key []byte, value []byte) error
	}
)
