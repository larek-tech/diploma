package controller

import (
	"context"

	"github.com/larek-tech/diploma/auth/internal/auth/model"
	"github.com/larek-tech/diploma/auth/pkg/jwt"
	"go.opentelemetry.io/otel/trace"
)

type authRepo interface {
	FindOneByEmail(ctx context.Context, email string) (model.UserDao, error)
	FindUserRoles(ctx context.Context, userID int64) ([]int64, error)
}

// Controller implements logic for authorization.
type Controller struct {
	tracer trace.Tracer
	ar     authRepo
	jwt    *jwt.Provider
}

// New creates new Controller.
func New(tracer trace.Tracer, ar authRepo, jwt *jwt.Provider) *Controller {
	return &Controller{
		tracer: tracer,
		ar:     ar,
		jwt:    jwt,
	}
}
