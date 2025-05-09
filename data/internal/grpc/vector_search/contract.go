package vector_search

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/document"
)

type (
	chunkStorage interface {
		Search(ctx context.Context, query []float32, sourceIDs []string, threshold float32, limit int) ([]*document.SearchResult, error)
	}
	embedder interface {
		CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
	}
)
