package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/document"
)

const systemPrompt = `
You are a question generation model. Your task is to generate questions based on the provided content.
`

func (s Service) generateQuestions(ctx context.Context, chunks []*document.Chunk) ([]*document.Questions, error) {
	wg := sync.WaitGroup{}
	questions := make([]*document.Questions, len(chunks))
	for i, chunk := range chunks {
		wg.Add(1)
		go func(chunk *document.Chunk) {
			defer wg.Done()
			question, err := s.llm.Call(ctx, systemPrompt+chunk.Content)
			if err != nil {
				slog.Error("failed to generate question", "error", err)
				return
			}
			embeds, err := s.embedder.CreateEmbedding(ctx, []string{question})
			if err != nil {
				slog.Error("failed to create embedding for question", "error", err)
				return
			}

			questions[i] = &document.Questions{
				ID:         uuid.NewString(),
				ChunkID:    chunk.ID,
				Question:   question,
				Embeddings: embeds[0],
			}
		}(chunk)
	}
	wg.Wait()
	return nil, nil
}
