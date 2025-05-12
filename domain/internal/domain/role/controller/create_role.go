package controller

import (
	"context"
	"slices"
	"time"

	"github.com/larek-tech/diploma/domain/internal/auth"
	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/role/model"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CreateRole create new role.
func (ctrl *Controller) CreateRole(ctx context.Context, req *pb.CreateRoleRequest, meta *authpb.UserAuthMetadata) (*pb.Role, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CreateRole",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("role", req.GetName()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return nil, errs.WrapErr(auth.ErrRequireAdmin, "create role")
	}

	role := model.RoleDao{
		Name:      req.GetName(),
		CreatedAt: time.Now(),
	}
	roleID, err := ctrl.rr.InsertRole(ctx, role)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	role.ID = roleID

	return role.ToProto(), nil
}
