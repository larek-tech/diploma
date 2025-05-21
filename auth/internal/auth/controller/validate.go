package controller

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
)

type tokenMeta struct {
	UserID int64   `json:"sub"`
	Roles  []int64 `json:"roles"`
}

// Validate validates the provided access token and returns user meta if it is correct.
func (ctrl *Controller) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Validate")
	defer span.End()

	token, err := ctrl.jwt.ParseAccessToken(req.GetToken())
	if err != nil {
		return nil, errs.WrapErr(err, "parse access token")
	}

	claims := token.Claims.(jwt.MapClaims)
	rawClaims, err := json.Marshal(claims)
	if err != nil {
		return nil, errs.WrapErr(err, "marshal token claims")
	}

	var meta tokenMeta
	if err = json.Unmarshal(rawClaims, &meta); err != nil {
		return nil, errs.WrapErr(err, "unmarshal token claims")
	}

	resp := &pb.ValidateResponse{
		Meta: &pb.UserAuthMetadata{
			UserId: meta.UserID,
			Roles:  meta.Roles,
		},
	}

	return resp, nil
}
