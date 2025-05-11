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

// GetSourceIDs returns external source ids by list of internal ids.
func (h *Handler) GetSourceIDs(ctx context.Context, req *pb.GetSourceIDsRequest) (*pb.GetSourceIDsResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.GetSourceIDs(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get source ids")
		return nil, status.Error(codes.Internal, "failed to get source ids")
	}

	return resp, status.Error(codes.OK, "successfully got source ids")
}
