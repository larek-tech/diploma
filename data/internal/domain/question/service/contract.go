package service

import (
	"context"
)

type (
	llm interface {
		Call(ctx context.Context, prompt string) (string, error)
	}
	embedder interface {
		CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error)
	}
)
