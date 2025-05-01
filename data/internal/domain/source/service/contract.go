package service

import (
	"context"
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/source"
)

type (
	sourceStorage interface {
		GetByName(ctx context.Context, name string) (*source.Source, error)
		GetByID(ctx context.Context, id string) (*source.Source, error)
		Save(ctx context.Context, source *source.Source) error
	}
	publisher interface {
		Publish(ctx context.Context, msg any, scheduledAt ...*time.Time) error
	}
	transactionalManager interface {
		Do(context.Context, func(context.Context) error) error
	}
)
