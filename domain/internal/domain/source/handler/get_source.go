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

// GetSource returns source by id.
func (h *Handler) GetSource(ctx context.Context, req *pb.GetSourceRequest) (*pb.GetSourceResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.GetSource(ctx, req.GetSourceId(), meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get source")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "source not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to get source")
		}
	}

	return resp, status.Error(codes.OK, "got source successfully")
}
