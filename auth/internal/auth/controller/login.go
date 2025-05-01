package controller

import (
	"context"
	"errors"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/larek-tech/diploma/auth/pkg/jwt"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/secure"
)

var (
	// ErrUserNotFound is an error when no user record was found.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is an error when provided invalid credentials.
	ErrInvalidCredentials = errors.New("invalid password or username")
)

// Login authorizes user with credentials and responds with access token.
func (ctrl *Controller) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Login")
	defer span.End()

	user, err := ctrl.ar.FindOneByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.WrapErr(errors.Join(err, ErrUserNotFound))
	}

	if !secure.VerifyPassword(user.HashPassword, req.GetPassword()) {
		return nil, errs.WrapErr(ErrInvalidCredentials, "verify password")
	}

	roles, err := ctrl.ar.FindUserRoles(ctx, user.ID)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	meta := &pb.UserAuthMetadata{
		UserId: user.ID,
		Roles:  roles,
	}
	token, err := ctrl.jwt.CreateAccessToken(meta)
	if err != nil {
		return nil, errs.WrapErr(err, "create access token")
	}

	resp := &pb.LoginResponse{
		Token: token,
		Type:  jwt.TypeBearerToken,
		Meta:  meta,
	}
	return resp, nil
}
