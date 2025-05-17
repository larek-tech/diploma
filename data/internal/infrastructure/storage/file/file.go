package file

import (
	"context"
	"fmt"

	"github.com/larek-tech/diploma/data/internal/domain/file"
	"github.com/larek-tech/diploma/data/internal/infrastructure/s3"
	postgres "github.com/larek-tech/diploma/data/internal/infrastructure/storage"
)

const (
	FileBucketName = "files"
	FileKeyPrefix  = "files/"
)

type Store struct {
	o  objectStore
	db db
}

func New(db db, objectStore objectStore) *Store {
	return &Store{
		o:  objectStore,
		db: db,
	}
}

func getObjectStoreKey(f *file.File) string {
	return FileKeyPrefix + f.ID + "." + f.Extension
}

func (s Store) Save(ctx context.Context, f *file.File) error {
	existingFile, err := s.GetByID(ctx, f.ID)
	if err != nil && !postgres.IsNoRowsError(err) {
		return err
	}
	f.ObjectKey = getObjectStoreKey(f)
	if existingFile != nil {
		err := s.db.Exec(ctx, `
UPDATE files
SET
	source_id = $1,
	filename = $2,
	extension = $3,
	object_key = $4,
	updated_at = NOW()
WHERE id = $5
`, f.SourceID, f.Filename, f.Extension, f.ObjectKey, f.ID)
		if err != nil {
			return err
		}
	} else {
		err := s.db.Exec(ctx, `
INSERT INTO files (id, source_id, filename, extension, object_key, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
`, f.ID, f.SourceID, f.Filename, f.Extension, f.ObjectKey)
		if err != nil {
			return err
		}
	}
	if f.Raw != nil {
		obj := s3.NewObject(FileBucketName, f.ObjectKey, f.Raw, s3.ContentTypeUndefined, nil)
		err := s.o.Upload(ctx, obj)
		if err != nil {
			return fmt.Errorf("failed to upload file raw content: %w", err)
		}
	}
	return nil
}

func (s Store) GetByID(ctx context.Context, id string) (*file.File, error) {
	var f file.File
	err := s.db.QueryStruct(ctx, &f, `--sql
SELECT
	id,
	source_id,
	filename,
	extension,
	object_key,
	created_at,
	updated_at
FROM files
WHERE id = $1
`, id)
	if err != nil {
		return nil, err
	}
	f.ObjectKey = getObjectStoreKey(&f)
	if f.ObjectKey != "" {
		obj, err := s.o.Download(ctx, FileBucketName, f.ObjectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to download file raw content: %w", err)
		}
		f.Raw = obj.GetData()
	}
	return &f, nil
}
