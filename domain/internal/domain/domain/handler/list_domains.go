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

// ListDomains returns the paginated list of available domains.
func (h *Handler) ListDomains(ctx context.Context, req *pb.ListDomainsRequest) (*pb.ListDomainsResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.dc.ListDomains(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("list domains")
		return nil, status.Error(codes.Internal, "failed to list domains")
	}

	return resp, status.Error(codes.OK, "got sources domain successfully")
}
