package page

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/postgres"
)

type Store struct {
	db db
}

func New(db db) *Store {
	return &Store{
		db: db,
	}
}
func (s Store) Save(ctx context.Context, page *site.Page) error {
	currentPage, err := s.GetByID(ctx, page.ID)
	if err != nil {
		return err
	}
	if currentPage != nil {
		err = s.db.Exec(ctx, `
UPDATE pages
SET
	url = $1,
	metadata = $2,
	raw = $3,
	content = $4,
	updated_at = now()
WHERE id = $5;
`, page.URL, page.Metadata, page.Raw, page.Content, page.ID)
		if err != nil {
			return err
		}
		return nil
	}
	err = s.db.Exec(ctx, `
INSERT INTO pages (id, url, metadata, raw, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now());
`, page.ID, page.URL, page.Metadata, page.Raw, page.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) GetByURL(ctx context.Context, url string) (*site.Page, error) {
	var page site.Page
	err := s.db.QueryStruct(ctx, &page, `
SELECT
	id,
	url,
	metadata,
	raw,
	content,
	created_at,
	updated_at
FROM pages
WHERE url = $1;
`, url)
	if err != nil {
		if postgres.IsNoRowsError(err) {
			return nil, nil
		}
	}
	return &page, err
}

func (s Store) GetByID(ctx context.Context, id string) (*site.Page, error) {
	var page site.Page
	err := s.db.QueryStruct(ctx, &page, `
SELECT
	id,
	url,
metadata,
raw,
content,
created_at,
updated_at
FROM pages
WHERE id = $1;
`, id)
	if err != nil {
		if postgres.IsNoRowsError(err) {
			return nil, nil
		}
	}
	return &page, err
}
