package embed_document

import (
	"context"
	"io"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/file"
	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type (
	embeddingService interface {
		Process(ctx context.Context, obj io.ReadSeeker, fileExt document.FileExtension, sourceObj any, sourceID string, metadata map[string]any) error
	}
	pageStore interface {
		GetByID(ctx context.Context, id string) (*site.Page, error)
	}
	siteStore interface {
		GetByID(ctx context.Context, id string) (*site.Site, error)
	}
	fileStore interface {
		GetByID(ctx context.Context, id string) (*file.File, error)
	}
)
