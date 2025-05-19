package service

import (
	"context"
	"net/url"

	"github.com/larek-tech/diploma/data/internal/domain/file"
	"github.com/larek-tech/diploma/data/internal/domain/sitemap"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
)

type (
	sourceStorage interface {
		GetByName(ctx context.Context, name string) (*source.Source, error)
		GetByID(ctx context.Context, id string) (*source.Source, error)
		Save(ctx context.Context, source *source.Source) error
	}
	fileStorage interface {
		GetByID(ctx context.Context, id string) (*file.File, error)
		Save(ctx context.Context, file *file.File) error
	}

	publisher interface {
		Publish(ctx context.Context, rawMsg []any, opts ...qaas.PublishOption) ([]string, error)
	}
	transactionalManager interface {
		Do(context.Context, func(context.Context) error) error
	}
	sitemapParser interface {
		GetAndParseSitemap(url.URL) ([]sitemap.URLResult, error)
	}
)
