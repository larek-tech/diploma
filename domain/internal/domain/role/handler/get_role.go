package handler

import (
	"context"
	"errors"

	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRole returns role by id.
func (h *Handler) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.Role, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get role meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.rc.GetRole(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get role")
		if errors.Is(err, auth.ErrRequireAdmin) {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}
		return nil, status.Error(codes.Internal, "failed to get role")
	}

	return resp, status.Error(codes.OK, "got role successfully")
}
