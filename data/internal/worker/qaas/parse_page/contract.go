package parse_page

import (
	"context"
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type (
	pageStore interface {
		Save(ctx context.Context, page *site.Page) error
		GetByURL(ctx context.Context, url string) (*site.Page, error)
		GetByID(ctx context.Context, id string) (*site.Page, error)
	}
	publisher interface {
		Publish(ctx context.Context, msg any, time ...*time.Time) error
	}
	pageService interface {
		ParsePage(ctx context.Context, page *site.Page) ([]*site.Page, error)
	}
)
