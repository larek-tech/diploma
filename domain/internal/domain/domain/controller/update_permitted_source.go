package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdatePermittedRoles updates domain roles permissions.
func (ctrl *Controller) UpdatePermittedRoles(ctx context.Context, req *pb.PermittedRoles, meta *authpb.UserAuthMetadata) (*pb.PermittedRoles, error) {
	userID := meta.GetUserId()
	domainID := req.GetResourceId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdatePermittedRoles",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("domainID", domainID),
		),
	)
	defer span.End()

	if err := ctrl.checkDomainCreator(ctx, domainID, meta); err != nil {
		return nil, errs.WrapErr(err)
	}

	updatedRoles, err := ctrl.dr.UpdatePermittedRoles(ctx, domainID, req.GetRoleIds())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return &pb.PermittedRoles{
		ResourceId: domainID,
		RoleIds:    updatedRoles,
	}, nil
}
