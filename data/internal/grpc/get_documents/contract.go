package get_documents

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/document"
)

type (
	documentsStore interface {
		GetMany(ctx context.Context, sourceID string, page, size int) (int, []*document.Document, error)
	}
)
