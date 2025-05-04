package handler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateSource updates source data.
func (h *Handler) UpdateSource(ctx context.Context, req *pb.UpdateSourceRequest) (*pb.Source, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if err = validateUpdateParams(req.GetUpdateParams()); err != nil {
		log.Err(errs.WrapErr(err)).Msg("validate update params")
		return nil, status.Error(codes.InvalidArgument, "invalid update params value")
	}

	resp, err := h.sc.UpdateSource(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("update source")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "source not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to update source")
		}
	}

	return resp, status.Error(codes.OK, "updated source successfully")
}
