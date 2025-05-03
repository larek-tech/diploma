package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateSource(ctx context.Context, req *pb.Source) (*pb.Source, error) {
	return req, status.Error(codes.Unimplemented, "")
}
