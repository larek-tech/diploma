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

// CreateRole creates new role.
func (h *Handler) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.Role, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.rc.CreateRole(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create role")
		if errors.Is(err, auth.ErrRequireAdmin) {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}
		return nil, status.Error(codes.PermissionDenied, "failed to create role")
	}

	return resp, status.Error(codes.OK, "created role successfully")
}
