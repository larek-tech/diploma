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

// CreateSource creates new source.
func (h *Handler) CreateSource(ctx context.Context, req *pb.CreateSourceRequest) (*pb.Source, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if err = validateUpdateParams(req.GetUpdateParams()); err != nil {
		log.Err(errs.WrapErr(err)).Msg("validate update params")
		return nil, status.Error(codes.InvalidArgument, "invalid update params value")
	}

	resp, err := h.sc.CreateSource(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create source")
		return nil, status.Error(codes.Internal, "failed to create source")
	}

	return resp, status.Error(codes.OK, "source created successfully")
}
