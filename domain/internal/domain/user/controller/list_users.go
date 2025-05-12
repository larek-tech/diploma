package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListUsers returns paginated list of users.
func (ctrl *Controller) ListUsers(
	ctx context.Context,
	req *pb.ListUsersRequest,
	meta *authpb.UserAuthMetadata,
) (*pb.ListUsersResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListUsers",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int("offset", int(req.GetOffset())),
			attribute.Int("limit", int(req.GetLimit())),
		),
	)
	defer span.End()

	usersDB, err := ctrl.ur.ListUsers(ctx, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	users := make([]*pb.User, len(usersDB))
	for idx := range usersDB {
		users[idx] = usersDB[idx].ToProto()
	}

	return &pb.ListUsersResponse{Users: users}, nil
}
