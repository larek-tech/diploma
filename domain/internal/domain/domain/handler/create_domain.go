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

// CreateDomain creates new domain.
func (h *Handler) CreateDomain(ctx context.Context, req *pb.CreateDomainRequest) (*pb.Domain, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.dc.CreateDomain(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create domain")
		return nil, status.Error(codes.Internal, "failed to create domain")
	}

	return resp, status.Error(codes.OK, "domain created successfully")
}
