package handler

import (
	"context"
	"errors"

	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/source/controller"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPermittedRoles returns list of roles permitted to source.
func (h *Handler) GetPermittedRoles(ctx context.Context, req *pb.GetResourcePermissionsRequest) (*pb.PermittedRoles, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.GetPermittedRoles(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get permitted roles")
		if errors.Is(err, controller.ErrNoAccessToSource) {
			return nil, status.Error(codes.PermissionDenied, "user doesn't have enough rights")
		}
		return nil, status.Error(codes.Internal, "failed to get permitted roles")
	}

	return resp, status.Error(codes.OK, "got permitted roles successfully")
}
