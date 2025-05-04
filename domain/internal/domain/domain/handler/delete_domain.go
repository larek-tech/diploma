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

// DeleteDomain deletes domain by id.
func (h *Handler) DeleteDomain(ctx context.Context, req *pb.DeleteDomainRequest) (*emptypb.Empty, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.dc.DeleteDomain(ctx, req.GetDomainId(), meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("delete domain")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "domain not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to delete domain")
		}
	}

	return &emptypb.Empty{}, status.Error(codes.OK, "deleted domain successfully")
}
