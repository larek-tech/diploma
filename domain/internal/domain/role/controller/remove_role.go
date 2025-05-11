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

// RemoveRole removes role from user.
func (ctrl *Controller) RemoveRole(ctx context.Context, req *pb.UpdateRoleRequest, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.RemoveRole",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("requested userID", req.GetUserId()),
			attribute.Int64("roleID", req.GetRoleId()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return errs.WrapErr(auth.ErrRequireAdmin, "remove role")
	}

	if err := ctrl.rr.RemoveRole(ctx, req.GetUserId(), req.GetRoleId()); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
