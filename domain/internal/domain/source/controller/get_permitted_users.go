package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetPermittedUsers returns list of users permitted to source.
func (ctrl *Controller) GetPermittedUsers(ctx context.Context, req *pb.GetResourcePermissionsRequest, meta *authpb.UserAuthMetadata) (*pb.PermittedUsers, error) {
	userID := meta.GetUserId()
	sourceID := req.GetResourceId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetPermittedUsers",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("sourceID", sourceID),
		),
	)
	defer span.End()

	if err := ctrl.checkSourceCreator(ctx, sourceID, meta); err != nil {
		return nil, errs.WrapErr(err)
	}

	permittedUsers, err := ctrl.sr.GetPermittedUsers(ctx, req.GetResourceId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return &pb.PermittedUsers{
		ResourceId: sourceID,
		UserIds:    permittedUsers,
	}, nil
}
