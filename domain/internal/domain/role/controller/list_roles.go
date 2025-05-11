package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListRoles returns paginated list of roles.
func (ctrl *Controller) ListRoles(
	ctx context.Context,
	req *pb.ListRolesRequest,
	meta *authpb.UserAuthMetadata,
) (*pb.ListRolesResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListRoles",
		trace.WithAttributes(
			attribute.Int64("roleID", meta.GetUserId()),
			attribute.Int("offset", int(req.GetOffset())),
			attribute.Int("limit", int(req.GetLimit())),
		),
	)
	defer span.End()

	rolesDB, err := ctrl.rr.ListRoles(ctx, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	roles := make([]*pb.Role, len(rolesDB))
	for idx := range rolesDB {
		roles[idx] = rolesDB[idx].ToProto()
	}

	return &pb.ListRolesResponse{Roles: roles}, nil
}
