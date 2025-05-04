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
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteSource deletes source by id.
func (h *Handler) DeleteSource(ctx context.Context, req *pb.DeleteSourceRequest) (*emptypb.Empty, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.sc.DeleteSource(ctx, req.GetSourceId(), meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("delete source")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "source not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to delete source")
		}
	}

	return &emptypb.Empty{}, status.Error(codes.OK, "deleted source successfully")
}
