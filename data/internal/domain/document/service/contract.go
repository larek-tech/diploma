package service

import (
	"context"
	"io"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/question"
)

type (
	documentStorage interface {
		Save(ctx context.Context, doc *document.Document) error
	}
	chunkStorage interface {
		Update(ctx context.Context, documentID string, chunks []*document.Chunk) error
		Delete(ctx context.Context, documentID string) error
	}
	embedder interface {
		CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
	}
	questionStorage interface {
		Save(ctx context.Context, questions []*question.Questions) error
	}
	trManager interface {
		Do(context.Context, func(ctx context.Context) error) error
	}
	parser interface {
		Parse(io.ReadSeeker) (string, error)
	}
	ocr interface {
		Process(string) (string, error)
	}
	questionService interface {
		GenerateQuestions(ctx context.Context, chunks []*document.Chunk) ([]*question.Questions, error)
	}
)
