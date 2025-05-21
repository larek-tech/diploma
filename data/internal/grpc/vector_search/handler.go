package vector_search

import (
	"context"
	"log/slog"
	"strings"

	"github.com/larek-tech/diploma/data/internal/data/pb"
	grpcSpan "github.com/larek-tech/diploma/data/internal/infrastructure/grpc/span"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	chunkStore chunkStorage
	embedder   embedder
	tracer     trace.Tracer
}

func New(chunkStore chunkStorage, embedder embedder, tracer trace.Tracer) *Handler {
	return &Handler{chunkStore: chunkStore, embedder: embedder, tracer: tracer}
}

func (h Handler) VectorSearch(ctx context.Context, in *pb.VectorSearchRequest) (*pb.VectorSearchResponse, error) {
	ctx, err := grpcSpan.GetTraceCtx(ctx)
	if err != nil {
		slog.Error("failed to get trace context", "error", err)
	}
	ctx, span := h.tracer.Start(ctx, "VectorSearch", trace.WithAttributes(
		attribute.String("query", in.Query),
		attribute.String("sourceIds", strings.Join(in.SourceIds, ",")),
		attribute.Float64("threshold", float64(in.Threshold)),
		attribute.Int64("topK", int64(in.TopK)),
		attribute.Bool("useQuestions", in.UseQuestions),
	))
	defer span.End()

	query, err := h.embedder.CreateEmbedding(ctx, []string{in.Query})
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "embedding error: %v", err)
	}

	res, err := h.chunkStore.Search(ctx, query[0], in.SourceIds, in.Threshold, int(in.TopK), in.UseQuestions)
	if err != nil {
		span.RecordError(err)
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
