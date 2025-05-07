package vector_search

import (
	"context"

	"github.com/larek-tech/diploma/data/internal/data/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedDataServiceServer
	chunkStore chunkStorage
	embedder   embedder
}

func New(chunkStore chunkStorage, embedder embedder) *Handler {
	return &Handler{chunkStore: chunkStore, embedder: embedder}
}

func (h Handler) VectorSearch(ctx context.Context, in *pb.VectorSearchRequest) (*pb.VectorSearchResponse, error) {

	query, err := h.embedder.CreateEmbedding(ctx, []string{in.Query})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "embedding error: %v", err)
	}

	res, err := h.chunkStore.Search(ctx, query[0], in.SourceIds, in.Threshold, int(in.TopK))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search error: %v", err)
	}

	var results []*pb.DocumentChunk
	for _, r := range res {
		results = append(results, &pb.DocumentChunk{
			Id:         r.ID,
			Index:      int64(r.Index),
			Content:    r.Content,
			Metadata:   r.Metadata,
			Similarity: r.CosineSimilarity,
		})
	}

	return &pb.VectorSearchResponse{Chunks: results}, status.New(codes.OK, "ok").Err()
}
