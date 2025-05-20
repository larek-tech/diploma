package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/question"
)

type Service struct {
	embedder        embedder
	llm             llm
	questionsPrompt string
}

func New(
	llm llm,
	embedder embedder,
	questionsPrompt ...string,
) *Service {
	if len(questionsPrompt) == 0 {
		questionsPrompt = []string{GetSystemPrompt()}
	}
	return &Service{
		llm:             llm,
		embedder:        embedder,
		questionsPrompt: questionsPrompt[0],
	}
}

func (s Service) GenerateQuestions(ctx context.Context, chunks []*document.Chunk) ([]*question.Questions, error) {
	if len(chunks) == 0 {
		slog.Error("no chunks provided")
		return nil, nil
	}
	questions := make([]*question.Questions, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk == nil {
			continue
		}
		llmQuestions, err := s.llm.Call(ctx, s.questionsPrompt+chunk.Content)
		if err != nil {
			slog.Error("failed to generate question", "error", err)
			return nil, err
		}
		embeds, err := s.embedder.CreateEmbedding(ctx, []string{llmQuestions})
		if err != nil {
			slog.Error("failed to create embedding for question", "error", err)
			return nil, err
		}

		questions = append(questions, &question.Questions{
			ID:         uuid.NewString(),
			ChunkID:    chunk.ID,
			Question:   llmQuestions,
			Embeddings: embeds[0],
		})

	}
	return questions, nil
}
