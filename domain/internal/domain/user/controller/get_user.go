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

// GetUser returns user by id.
func (ctrl *Controller) GetUser(ctx context.Context, req *pb.GetUserRequest, meta *authpb.UserAuthMetadata) (*pb.User, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetUser",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("requested userID", req.GetUserId()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return nil, errs.WrapErr(auth.ErrRequireAdmin, "get user")
	}

	user, err := ctrl.ur.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return user.ToProto(), nil
}
