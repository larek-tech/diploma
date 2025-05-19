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

// ListSourcesByDomain returns paginated list of sources by specified domain.
func (h *Handler) ListSourcesByDomain(ctx context.Context, req *pb.ListSourcesByDomainRequest) (*pb.ListSourcesResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.ListSourcesByDomain(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("list sources by domain")
		return nil, status.Error(codes.Internal, "failed to list sources by domain")
	}

	return resp, status.Error(codes.OK, "got sources list by domain successfully")
}
