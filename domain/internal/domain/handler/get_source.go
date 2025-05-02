package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetSource(ctx context.Context, req *pb.GetSourceRequest) (*pb.GetSourceResponse, error) {
	return &pb.GetSourceResponse{}, status.Error(codes.Unimplemented, "")
}
