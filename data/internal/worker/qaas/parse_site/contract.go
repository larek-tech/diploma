package parse_site

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
)

type (
	siteStore interface {
		Save(ctx context.Context, site *site.Site) error
	}
	publisher interface {
		Publish(ctx context.Context, rawMsg []any, opts ...qaas.PublishOption) ([]string, error)
	}
)
