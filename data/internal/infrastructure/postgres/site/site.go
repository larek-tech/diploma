package site

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/postgres"
)

type Storage struct {
	db db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) Save(ctx context.Context, site *site.Site) error {
	currentSite, err := s.GetByURL(ctx, site.URL)
	if err != nil {
		// if record not found, create a new one
		if postgres.IsNoRowsError(err) {
			err = s.db.Exec(ctx, `
INSERT INTO sites (id, source_id, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);
`, site.ID, site.SourceID, site.URL, site.CreatedAt, site.UpdatedAt)
			return err
		}
		return err
	}
	if currentSite != nil {
		err = s.db.Exec(ctx, `
UPDATE sites
SET source_id = $1, url = $2, updated_at = $3
WHERE id = $4;
`, site.SourceID, site.URL, site.UpdatedAt, site.ID)
		site.ID = currentSite.ID
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) GetByID(ctx context.Context, id string) (*site.Site, error) {
	var res site.Site
	err := s.db.QueryStruct(ctx, &res, `
SELECT
    id,
    source_id,
    url,
	available_pages,
    created_at,
    updated_at
FROM sites WHERE id = $1;
`, id)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s Storage) GetByURL(ctx context.Context, url string) (*site.Site, error) {
	var res site.Site
	err := s.db.QueryStruct(ctx, &res, `
SELECT
    id,
    source_id,
    url,
	available_pages,
    created_at,
    updated_at
FROM sites WHERE url = $1;
`, url)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
