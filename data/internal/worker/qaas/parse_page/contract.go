package parse_page

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
)

type (
	pageStore interface {
		Save(ctx context.Context, page *site.Page) error
		GetByURL(ctx context.Context, url string) (*site.Page, error)
		GetByID(ctx context.Context, id string) (*site.Page, error)
	}
	publisher interface {
		Publish(ctx context.Context, rawMsg []any, opts ...qaas.PublishOption) ([]string, error)
	}
	pageService interface {
		ParsePage(ctx context.Context, page *site.Page) ([]*site.Page, bool, error)
	}
)
