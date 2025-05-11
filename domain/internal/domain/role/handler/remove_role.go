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
	"google.golang.org/protobuf/types/known/emptypb"
)

// RemoveRole removes role from user.
func (h *Handler) RemoveRole(ctx context.Context, req *pb.UpdateRoleRequest) (*emptypb.Empty, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get role meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.rc.RemoveRole(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("remove role")
		if errors.Is(err, auth.ErrRequireAdmin) {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}
		return nil, status.Error(codes.PermissionDenied, "failed to remove role")
	}

	return &emptypb.Empty{}, status.Error(codes.OK, "removed role successfully")
}
