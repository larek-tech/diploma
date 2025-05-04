package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListSources returns the paginated list of available sources.
func (h *Handler) ListSources(ctx context.Context, req *pb.ListSourcesRequest) (*pb.ListSourcesResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.ListSources(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("list sources")
		return nil, status.Error(codes.Internal, "failed to list sources")
	}

	return resp, status.Error(codes.OK, "got sources list successfully")
}
