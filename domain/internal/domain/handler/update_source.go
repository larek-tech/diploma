package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateSource(ctx context.Context, req *pb.UpdateSourceRequest) (*pb.Source, error) {
	return &pb.Source{
		Title:        req.GetTitle(),
		Content:      req.GetContent(),
		Typ:          0,
		UpdateParams: req.GetUpdateParams(),
		Credentials:  req.GetCredentials(),
	}, status.Error(codes.Unimplemented, "")
}
