package chunk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/larek-tech/diploma/data/internal/domain/document"
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

func (s Storage) Update(ctx context.Context, documentID string, chunks []*document.Chunk) error {
	return s.trManager.Do(ctx, func(txCtx context.Context) error {
		if err := s.db.Exec(txCtx, "DELETE FROM chunks WHERE document_id = $1", documentID); err != nil {
			return fmt.Errorf("failed to delete old chunks: %w", err)
		}
		for _, chunk := range chunks {
			metadata, err := json.Marshal(chunk.Metadata)
			if err != nil {
				return fmt.Errorf("failed to marshal metadata: %w", err)
			}
			embeddingsBytes, err := json.Marshal(chunk.Embeddings)
			if err != nil {
				return fmt.Errorf("failed to marshal embeddings: %w", err)
			}
			embeddings := string(embeddingsBytes) // Convert []byte to string

			if err := s.db.Exec(
				txCtx,
				`INSERT INTO chunks (id, index, document_id, content, metadata, embeddings)
     VALUES ($1, $2, $3, $4, $5, $6)`,
				chunk.ID, chunk.Index, documentID, chunk.Content, metadata, embeddings,
			); err != nil {
				return fmt.Errorf("failed to insert chunk: %w", err)
			}
		}
		return nil
	})
}

func (s Storage) Delete(ctx context.Context, documentID string) error {
	return s.db.Exec(ctx, "DELETE FROM chunks WHERE document_id = $1", documentID)
}
