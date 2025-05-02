package grpcclient

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// PushOutMeta pushes value with a key into outgoing grpc context metadata.
//
// Returns result context.
func PushOutMeta(ctx context.Context, key, val string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, val)
}
