package controller

import (
	"context"
	"slices"

	"github.com/larek-tech/diploma/domain/internal/auth"
	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteRole soft delete user.
func (ctrl *Controller) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteRole",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("roleID", req.GetRoleId()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return errs.WrapErr(auth.ErrRequireAdmin, "delete user")
	}

	if err := ctrl.rr.DeleteRole(ctx, req.GetRoleId()); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
