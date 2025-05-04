package service

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/document"
)

type (
	documentStorage interface {
		Save(ctx context.Context, doc *document.Document) error
		Get(ctx context.Context, id string) (*document.Document, error)
	}
	chunkStorage interface {
		Update(ctx context.Context, documentID string, chunks []*document.Chunk) error
		Delete(ctx context.Context, documentID string) error
	}
	embedder interface {
		CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
	}
	trManager interface {
		Do(context.Context, func(ctx context.Context) error) error
	}
)
