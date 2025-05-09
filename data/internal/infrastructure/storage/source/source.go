package source

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/source"
	storage "github.com/larek-tech/diploma/data/internal/infrastructure/storage"
)

type Storage struct {
	db db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) Save(ctx context.Context, src *source.Source) error {
	var currentSource *source.Source
	currentSource, err := s.GetByName(ctx, src.Title)
	if err != nil {
		return err
	}
	if currentSource == nil {
		err = s.db.Exec(ctx, `
INSERT INTO sources (id, title, type, credentials)
VALUES ($1, $2, $3, $4);
`, src.ID, src.Title, src.Type, src.Credentials)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.db.Exec(ctx, `
UPDATE sources
SET title = $1, type = $2, credentials = $3
WHERE id = $4;
`, src.Title, src.Type, src.Credentials, src.ID)
	src.ID = currentSource.ID
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) GetByName(ctx context.Context, name string) (*source.Source, error) {
	var res source.Source
	err := s.db.QueryStruct(ctx, &res, `
SELECT
	id,
	title,
	type,
	credentials
FROM sources 
WHERE title = $1;
`, name)
	if err != nil {
		if storage.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (s Storage) GetByID(ctx context.Context, id string) (*source.Source, error) {
	var res source.Source
	err := s.db.QueryStruct(ctx, &res, `
SELECT
	id,
	title,
	type,
	credentials
FROM sources
WHERE id = $1;
`, id)
	if err != nil {
		if storage.IsNoRowsError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}
