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
	currentPage, err := s.GetByURL(ctx, page.URL)
	if err != nil {
		return err
	}
	if currentPage != nil {
		err = s.db.Exec(ctx, `
UPDATE pages
SET
    site_id = $1,
    url = $2,
    metadata = $3,
    raw = $4,
    content = $5,
    updated_at = now()
WHERE id = $6;
`, page.SiteID, page.URL, page.Metadata, page.Raw, page.Content, currentPage.ID)
		if err != nil {
			return err
		}
		*page = *currentPage // update all fields in the passed-in page pointer
		return nil
	}
	err = s.db.Exec(ctx, `
INSERT INTO pages (id, site_id, url, metadata, raw, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now());
`, page.ID, page.SiteID, page.URL, page.Metadata, page.Raw, page.Content)
	if err != nil {
		return err
	}
	return nil
}

// FIXME: multiple parse_pages can return same page_url as outgoing and it goes to parse_page forever
func (s Store) GetByURL(ctx context.Context, url string) (*site.Page, error) {
	var page site.Page
	err := s.db.QueryStruct(ctx, &page, `
SELECT
	id,
	site_id,
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
	site_id,
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
