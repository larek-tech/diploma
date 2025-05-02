package document

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/data/internal/domain/document"
)

type Storage struct {
	db db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) Save(ctx context.Context, doc *document.Document) error {
	if doc == nil {
		return nil
	}
	// check if document with given ID already exists
	var res document.Document
	sql := `
SELECT 
id
FROM documents
WHERE id = $1;
`
	err := s.db.QueryStruct(ctx, &res, sql, doc.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sql = `
INSERT INTO documents (id, object_id, object_type, source_id, name, content, metadata, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
			err = s.db.Exec(ctx, sql, doc.ID, doc.ObjectID, doc.ObjectType, doc.SourceID, doc.Name, doc.Content, doc.Metadata, doc.CreatedAt, doc.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existance of document: %w", err)
		}
	}
	// if document with given ID already exists, update it
	sql = `
UPDATE documents
SET 
    source_id = $1,
    object_id = $2,
    object_type = $3,
    name = $4,
    content = $5,
    metadata = $6,
    updated_at = $7
WHERE id = $8`
	err = s.db.Exec(ctx, sql, doc.SourceID, doc.ObjectID, doc.ObjectType, doc.Name, doc.Content, doc.Metadata, doc.UpdatedAt, doc.ID)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	return nil
}
