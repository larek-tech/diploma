package question

import (
	"context"
)

type (
	db interface {
		Exec(ctx context.Context, sql string, args ...interface{}) error
		QueryStruct(ctx context.Context, dst interface{}, sql string, args ...interface{}) error
		QueryStructs(ctx context.Context, dst interface{}, sql string, args ...interface{}) error
	}
	trManager interface {
		Do(context.Context, func(context.Context) error) error
	}
)
