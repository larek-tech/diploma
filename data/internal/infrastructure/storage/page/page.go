package page

import (
	"context"
	"fmt"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/s3"
	postgres "github.com/larek-tech/diploma/data/internal/infrastructure/storage"
)

const (
	PageBucketName = "pages"
	PageKeyPrefix  = "page/"
)

type Store struct {
	db          db
	objectStore objectStore
}

func New(db db, objectStore objectStore) *Store {
	return &Store{
		db:          db,
		objectStore: objectStore,
	}
}

func getObjectStoreKey(page *site.Page) string {
	return PageKeyPrefix + page.ID + ".html"
}

func assembleObject(page *site.Page) *s3.Object {
	metadata := map[string]string{
		"site_id": page.SiteID,
		"page_id": page.ID,
		"url":     page.URL,
	}
	obj := s3.NewObject(PageBucketName, getObjectStoreKey(page), []byte(page.Raw), s3.ContentTypeHTML, metadata)
	return obj
}

func (s Store) Save(ctx context.Context, page *site.Page) error {
	page.RawObjectID = getObjectStoreKey(page)

	currentPage, err := s.GetByURL(ctx, page.URL)
	if err != nil && !postgres.IsNoRowsError(err) {
		return err
	}

	// Determine whether to update or insert
	if currentPage != nil {
		err = s.db.Exec(ctx, `
UPDATE pages
SET
	site_id = $1,
	url = $2,
	raw_object_id = $3,
	metadata = $4,
	content = $5,
	updated_at = now()
WHERE id = $6;
`, page.SiteID, page.URL, page.RawObjectID, page.Metadata, page.Content, currentPage.ID)

		if err == nil {
			*page = *currentPage // update all fields in the passed-in page pointer
		}
	} else {
		err = s.db.Exec(ctx, `
INSERT INTO pages (id, site_id, url, raw_object_id, metadata, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now());
`, page.ID, page.SiteID, page.URL, page.RawObjectID, page.Metadata, page.Content)
	}

	if err != nil {
		return err
	}

	// Upload the object after DB operation if we have a raw object ID
	if page.RawObjectID != "" {
		if storeErr := s.objectStore.Upload(ctx, assembleObject(page)); storeErr != nil {
			return fmt.Errorf("failed to upload page raw content: %w", storeErr)
		}
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
	raw_object_id,
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
	if page.RawObjectID != "" {
		obj, err := s.objectStore.Download(ctx, PageBucketName, page.RawObjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to download page raw content: %w", err)
		}
		page.Raw = string(obj.GetData())
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
	raw_object_id,
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
	if page.RawObjectID != "" {
		obj, err := s.objectStore.Download(ctx, PageBucketName, page.RawObjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to download page raw content: %w", err)
		}
		page.Raw = string(obj.GetData())
	}
	return &page, err
}
