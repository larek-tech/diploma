package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdatePermittedUsers updates domain user permissions.
func (ctrl *Controller) UpdatePermittedUsers(ctx context.Context, req *pb.PermittedUsers, meta *authpb.UserAuthMetadata) (*pb.PermittedUsers, error) {
	userID := meta.GetUserId()
	sourceID := req.GetResourceId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdatePermittedUsers",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("sourceID", sourceID),
		),
	)
	defer span.End()

	if err := ctrl.checkSourceCreator(ctx, sourceID, meta); err != nil {
		return nil, errs.WrapErr(err)
	}

	updatedUsers, err := ctrl.sr.UpdatePermittedUsers(ctx, sourceID, req.GetUserIds())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return &pb.PermittedUsers{
		ResourceId: sourceID,
		UserIds:    updatedUsers,
	}, nil
}
