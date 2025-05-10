package object_store

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/infrastructure/s3"
)

type (
	db interface {
		Exec(ctx context.Context, sql string, args ...interface{}) error
		QueryStruct(ctx context.Context, dst interface{}, sql string, args ...interface{}) error
		QueryStructs(ctx context.Context, dst interface{}, sql string, args ...interface{}) error
	}
	objectStore interface {
		Upload(ctx context.Context, object *s3.Object) error
		Download(ctx context.Context, bucketName, key string) (*s3.Object, error)
	}
)
