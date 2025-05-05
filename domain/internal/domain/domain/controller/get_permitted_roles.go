package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetPermittedRoles returns list of roles permitted to domain.
func (ctrl *Controller) GetPermittedRoles(ctx context.Context, req *pb.GetResourcePermissionsRequest, meta *authpb.UserAuthMetadata) (*pb.PermittedRoles, error) {
	userID := meta.GetUserId()
	domainID := req.GetResourceId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetPermittedRoles",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("domainID", domainID),
		),
	)
	defer span.End()

	if err := ctrl.checkDomainCreator(ctx, domainID, meta); err != nil {
		return nil, errs.WrapErr(err)
	}

	permittedRoles, err := ctrl.dr.GetPermittedRoles(ctx, req.GetResourceId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return &pb.PermittedRoles{
		ResourceId: domainID,
		RoleIds:    permittedRoles,
	}, nil
}
