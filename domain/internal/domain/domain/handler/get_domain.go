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

// GetDomain returns domain by id.
func (h *Handler) GetDomain(ctx context.Context, req *pb.GetDomainRequest) (*pb.GetDomainResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.dc.GetDomain(ctx, req.GetDomainId(), meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get domain")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "domain not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to get domain")
		}
	}

	return resp, status.Error(codes.OK, "got domain successfully")
}
