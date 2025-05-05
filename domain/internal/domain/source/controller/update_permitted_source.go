package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdatePermittedRoles updates source roles permissions.
func (ctrl *Controller) UpdatePermittedRoles(ctx context.Context, req *pb.PermittedRoles, meta *authpb.UserAuthMetadata) (*pb.PermittedRoles, error) {
	userID := meta.GetUserId()
	sourceID := req.GetResourceId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdatePermittedRoles",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("sourceID", sourceID),
		),
	)
	defer span.End()

	if err := ctrl.checkSourceCreator(ctx, sourceID, meta); err != nil {
		return nil, errs.WrapErr(err)
	}

	updatedRoles, err := ctrl.sr.UpdatePermittedRoles(ctx, sourceID, req.GetRoleIds())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return &pb.PermittedRoles{
		ResourceId: sourceID,
		RoleIds:    updatedRoles,
	}, nil
}
