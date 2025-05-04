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

// UpdateDomain updates domain data.
func (h *Handler) UpdateDomain(ctx context.Context, req *pb.UpdateDomainRequest) (*pb.Domain, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.dc.UpdateDomain(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("update domain")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "domain not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to update domain")
		}
	}

	return resp, status.Error(codes.OK, "updated domain successfully")
}
