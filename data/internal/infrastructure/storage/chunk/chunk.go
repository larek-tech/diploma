package chunk

import (
	"context"
	"encoding/json"
	"fmt"
	"unicode/utf8"

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

func prepareVector(embeddings []float32) string {
	if len(embeddings) == 0 {
		return "[]"
	}
	embeddingsBytes, _ := json.Marshal(embeddings)

	return string(embeddingsBytes)
}
func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				v = append(v, '�')
				continue
			}
		}
		v = append(v, r)
	}
	return string(v)
}

func (s Storage) Update(ctx context.Context, documentID string, chunks []*document.Chunk) error {
	return s.trManager.Do(ctx, func(txCtx context.Context) error {
		if err := s.db.Exec(txCtx, "DELETE FROM chunks WHERE document_id = $1", documentID); err != nil {
			return fmt.Errorf("failed to delete old chunks: %w", err)
		}
		for _, chunk := range chunks {
			chunk.Content = document.CleanUTF8(chunk.Content)

			if err := s.db.Exec(
				txCtx,
				`INSERT INTO chunks (id, index, source_id, document_id, content, embeddings)
     VALUES ($1, $2, $3, $4, $5, $6)`,
				chunk.ID, chunk.Index, chunk.SourceID, documentID, sanitizeUTF8(chunk.Content), prepareVector(chunk.Embeddings),
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

func (s Storage) Search(ctx context.Context, query []float32, sourceIDs []string, threshold float32, limit int) ([]*document.SearchResult, error) {
	if len(query) == 0 {
		return nil, fmt.Errorf("query is empty")
	}
	if len(sourceIDs) == 0 {
		return nil, fmt.Errorf("sourceIDs is empty")
	}

	sql := `
SELECT
    id,
    index,
    source_id,
    document_id,
    content,
    1 - (embeddings <=> $1) AS cosine_similarity
FROM chunks
WHERE source_id = ANY($2) AND 1 - (embeddings <=> $1) > $3
ORDER BY 1 - (embeddings <=> $1) desc 
LIMIT $4;`

	var res []*document.SearchResult
	err := s.db.QueryStructs(ctx, &res, sql, prepareVector(query), sourceIDs, threshold, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query chunks: %w", err)
	}
	//for _, r := range res {
	//	if len(r.Metadata) == 0 {
	//		continue
	//	}
	//	decoded, err := base64.StdEncoding.DecodeString(string(r.Metadata))
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to base64 decode metadata: %w", err)
	//	}
	//	if err := json.Unmarshal(decoded, &r.Metadata); err != nil {
	//		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	//	}
	//}
	return res, nil
}
