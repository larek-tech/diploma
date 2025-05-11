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

// DeleteUser soft delete user.
func (ctrl *Controller) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteUser",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("request userID", req.GetUserId()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return errs.WrapErr(auth.ErrRequireAdmin, "delete user")
	}

	if err := ctrl.ur.DeleteUser(ctx, req.GetUserId()); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
