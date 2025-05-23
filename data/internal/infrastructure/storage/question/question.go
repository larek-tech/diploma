package question

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/larek-tech/diploma/data/internal/domain/question"
)

type Storage struct {
	db        db
	trManager trManager
}

func New(db db, trManager trManager) *Storage {
	return &Storage{
		db:        db,
		trManager: trManager,
	}
}

func prepareVector(embeddings []float32) string {
	if len(embeddings) == 0 {
		return "[]"
	}
	embeddingsBytes, _ := json.Marshal(embeddings)

	return string(embeddingsBytes)
}

func (s Storage) Save(ctx context.Context, questions []*question.Questions) error {
	return s.trManager.Do(ctx, func(txCtx context.Context) error {
		for _, q := range questions {
			if err := s.db.Exec(
				txCtx,
				`INSERT INTO chunk_questions (id, chunk_id, question, embeddings)
				 VALUES ($1, $2, $3, $4)`,
				q.ID, q.ChunkID, q.Question, prepareVector(q.Embeddings),
			); err != nil {
				return fmt.Errorf("failed to insert question: %w", err)
			}
		}
		return nil
	})
}
