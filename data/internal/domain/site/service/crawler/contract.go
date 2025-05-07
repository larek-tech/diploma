package crawler

import (
	"context"
	"net/http"

	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type (
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
	transactionalManager interface {
		Do(context.Context, func(context.Context) error) error
	}
	siteStore interface {
		GetByID(ctx context.Context, id string) (*site.Site, error)
		GetByURL(ctx context.Context, url string) (*site.Site, error)
		Save(ctx context.Context, site *site.Site) error
	}
	pageStore interface {
		GetByID(ctx context.Context, id string) (*site.Page, error)
		GetByURL(ctx context.Context, url string) (*site.Page, error)
		Save(ctx context.Context, page *site.Page) error
	}
	pageJobStore interface {
		IsAlreadyParsed(ctx context.Context, parseSiteJobID string) (bool, error)
	}
)
