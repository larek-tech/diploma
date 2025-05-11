package controller

import (
	"context"
	"slices"
	"time"

	"github.com/larek-tech/diploma/domain/internal/auth"
	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/user/model"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/secure"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CreateUser create new user.
func (ctrl *Controller) CreateUser(ctx context.Context, req *pb.CreateUserRequest, meta *authpb.UserAuthMetadata) (*pb.User, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CreateUser",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("email", req.GetEmail()),
		),
	)
	defer span.End()

	if !slices.Contains(meta.GetRoles(), auth.AdminRoleID) {
		return nil, errs.WrapErr(auth.ErrRequireAdmin, "create user")
	}

	hashPassword, err := secure.Encrypt(req.GetPassword(), ctrl.encryption)
	if err != nil {
		return nil, errs.WrapErr(err, "encrypt password")
	}

	user := model.UserDao{
		Email:        req.GetEmail(),
		HashPassword: hashPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	userID, err := ctrl.ur.InsertUser(ctx, user)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	user.ID = userID

	return user.ToProto(), nil
}
