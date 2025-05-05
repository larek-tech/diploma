package handler

import (
	"context"
	"errors"

	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/domain/controller"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePermittedRoles updates list of domain permitted roles.
func (h *Handler) UpdatePermittedRoles(ctx context.Context, req *pb.PermittedRoles) (*pb.PermittedRoles, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.dc.UpdatePermittedRoles(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("update permitted roles")
		if errors.Is(err, controller.ErrNoAccessToDomain) {
			return nil, status.Error(codes.PermissionDenied, "user doesn't have enough rights")
		}
		return nil, status.Error(codes.Internal, "failed to update permitted roles")
	}

	return resp, status.Error(codes.OK, "updated permitted roles successfully")
}
