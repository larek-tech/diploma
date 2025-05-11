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

// GetRole returns role by id.
func (ctrl *Controller) GetRole(ctx context.Context, req *pb.GetRoleRequest, meta *authpb.UserAuthMetadata) (*pb.Role, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetRole",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("roleID", req.GetRoleId()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return nil, errs.WrapErr(auth.ErrRequireAdmin, "get role")
	}

	role, err := ctrl.rr.GetRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return role.ToProto(), nil
}
