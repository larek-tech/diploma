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

// SetRole adds role to user.
func (h *Handler) SetRole(ctx context.Context, req *pb.UpdateRoleRequest) (*emptypb.Empty, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get role meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.rc.SetRole(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("set role")
		if errors.Is(err, auth.ErrRequireAdmin) {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}
		return nil, status.Error(codes.Internal, "failed to set role")
	}

	return &emptypb.Empty{}, status.Error(codes.OK, "set role successfully")
}
