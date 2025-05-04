package embed_document

import (
	"context"
	"io"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type (
	embeddingService interface {
		Process(ctx context.Context, obj io.ReadSeeker, fileExt document.FileExtension, sourceID, objType, objectID string) error
	}
	pageStore interface {
		GetByID(ctx context.Context, id string) (*site.Page, error)
	}
	siteStore interface {
		GetByID(ctx context.Context, id string) (*site.Site, error)
	}
	documentStore interface {
		GetByID(ctx context.Context, id string) (*document.Document, error)
		Save(ctx context.Context, doc *document.Document) error
	}
)
