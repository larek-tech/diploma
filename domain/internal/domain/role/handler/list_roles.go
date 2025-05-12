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

// ListRoles returns paginated list of roles.
func (h *Handler) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.rc.ListRoles(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("list roles")
		return nil, status.Error(codes.Internal, "failed to list roles")
	}

	return resp, status.Error(codes.OK, "got roles list successfully")
}
