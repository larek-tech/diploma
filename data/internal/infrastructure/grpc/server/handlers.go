package server

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/data/pb"
)

type Handlers struct {
	pb.UnimplementedDataServiceServer
	vh  VectorSearchHandler
	gdh GetDocumentsHandler
}

func NewHandlers(vectorSearchHandler VectorSearchHandler, getDocumentsHandler GetDocumentsHandler) *Handlers {
	return &Handlers{
		UnimplementedDataServiceServer: pb.UnimplementedDataServiceServer{},
		vh:                             vectorSearchHandler,
		gdh:                            getDocumentsHandler,
	}
}

func (h Handlers) VectorSearch(ctx context.Context, in *pb.VectorSearchRequest) (*pb.VectorSearchResponse, error) {
	return h.vh.VectorSearch(ctx, in)
}

func (h Handlers) GetDocuments(ctx context.Context, in *pb.GetDocumentsIn) (*pb.GetDocumentsOut, error) {
	return h.gdh.GetDocuments(ctx, in)
}
