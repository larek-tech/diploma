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

// CreateUser creates new user.
func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.uc.CreateUser(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create user")
		if errors.Is(err, auth.ErrRequireAdmin) {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return resp, status.Error(codes.OK, "created user successfully")
}
