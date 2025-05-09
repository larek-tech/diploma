package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ListSources(ctx context.Context, req *pb.ListSourcesRequest) (*pb.ListSourcesResponse, error) {
	return &pb.ListSourcesResponse{}, status.Error(codes.Unimplemented, "")
}
