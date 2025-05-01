package parse_site

import (
	"context"
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type (
	siteStore interface {
		Save(ctx context.Context, site *site.Site) error
	}
	pagePublisher interface {
		Publish(context.Context, any, ...*time.Time) error
	}
)
