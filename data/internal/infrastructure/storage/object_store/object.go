package object_store

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/domain/object_store"
	"github.com/larek-tech/diploma/data/internal/infrastructure/s3"
	postgres "github.com/larek-tech/diploma/data/internal/infrastructure/storage"
)

const (
	ObjectBucketName = "objects"
	ObjectKeyPrefix  = "object/"
)

type Storage struct {
	db          db
	objectStore objectStore
}

func NewStorage(db db, objectStore objectStore) *Storage {
	return &Storage{
		db:          db,
		objectStore: objectStore,
	}
}

func getObjectStoreKey(object *object_store.Object) string {
	return ObjectKeyPrefix + object.ID + object.ContentType
}

func assembleObject(object *object_store.Object) *s3.Object {
	metadata := map[string]string{
		"object_id": object.ID,
	}
	obj := s3.NewObject(ObjectBucketName, getObjectStoreKey(object), object.Data, s3.ContentTypeUndefined, metadata)
	return obj
}

func (s Storage) Save(ctx context.Context, object *object_store.Object) error {
	currentObject, err := s.GetByID(ctx, object.ID)
	if err != nil && !postgres.IsNoRowsError(err) {
		return err
	}

	if currentObject != nil {
		object.RawContentID = getObjectStoreKey(object)
		object.ContentType = currentObject.ContentType
		object.Size = currentObject.Size
		object.CreatedAt = currentObject.CreatedAt
		object.UpdatedAt = currentObject.UpdatedAt
		err = s.db.Exec(ctx, `
UPDATE object_storage_files
SET
	object_storage_id = $1,
	raw_object_id = $2,
	content_type = $3,
	content_size = $4,
		`)
		if err != nil {
			return err
		}
	}

	if object.Data != nil {
		err = s.objectStore.Upload(ctx, assembleObject(object))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) GetByID(ctx context.Context, id string) (*object_store.Object, error) {
	var obj object_store.Object
	sql := `
SELECT
	id,
	object_storage_id,
	raw_object_id,
	content_type,
	content_size,
	created_at,
	updated_at
FROM
	object_storage_files
WHERE id = $1
`
	err := s.db.QueryStruct(ctx, &obj, sql, id)
	if err != nil {
		return nil, err
	}
	if obj.RawContentID != "" {
		// Get the object from the object store
		object, err := s.objectStore.Download(ctx, ObjectBucketName, obj.RawContentID)
		if err != nil {
			return nil, err
		}
		obj.Data = object.GetData()
	}

	return &obj, nil
}
