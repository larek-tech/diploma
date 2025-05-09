package server

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/data/pb"
)

type (
	VectorSearchHandler interface {
		VectorSearch(ctx context.Context, in *pb.VectorSearchRequest) (*pb.VectorSearchResponse, error)
	}
	GetDocumentsHandler interface {
		GetDocuments(context.Context, *pb.GetDocumentsIn) (*pb.GetDocumentsOut, error)
	}
)
