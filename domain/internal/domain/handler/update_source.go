package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) UpdateSource(ctx context.Context, req *pb.UpdateSourceRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, status.Error(codes.Unimplemented, "")
}
