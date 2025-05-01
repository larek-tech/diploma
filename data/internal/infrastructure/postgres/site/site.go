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
	currentSite, err := s.GetByID(ctx, site.ID)
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
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) GetByID(ctx context.Context, id string) (*site.Site, error) {
	var site site.Site
	err := s.db.QueryStruct(ctx, &site, `
SELECT
    id,
    source_id,
    url,
    created_at,
    updated_at
FROM sites WHERE id = $1;
`, id)
	if err != nil {
		return nil, err
	}
	pagesIDs, err := s.fetchPagesIDs(ctx, site)
	if err != nil {
		return nil, err
	}
	site.AvailablePages = pagesIDs

	return &site, nil
}

func (s Storage) GetByURL(ctx context.Context, url string) (*site.Site, error) {
	var site site.Site
	err := s.db.QueryStruct(ctx, &site, `
SELECT
    id,
    source_id,
    url,
    created_at,
    updated_at
FROM sites WHERE url = $1;
`, url)
	if err != nil {
		return nil, err
	}
	pagesIDs, err := s.fetchPagesIDs(ctx, site)
	if err != nil {
		return nil, err
	}
	site.AvailablePages = pagesIDs

	return &site, nil
}

func (s Storage) fetchPagesIDs(ctx context.Context, site site.Site) ([]string, error) {
	var pagesIDS []string
	err := s.db.QueryStructs(ctx, &pagesIDS, `
SELECT 
	id
FROM pages
WHERE site_id = $1;
`, site.ID)
	if err != nil {
		if !postgres.IsNoRowsError(err) {
			return nil, err
		}
		// No pages found, return empty slice
		return []string{}, nil
	}
	return pagesIDS, nil
}
