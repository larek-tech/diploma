package controller

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/user/model"
	"go.opentelemetry.io/otel/trace"
)

type userRepo interface {
	InsertUser(ctx context.Context, u model.UserDao) (int64, error)
	GetUser(ctx context.Context, id int64) (model.UserDao, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, offset, limit uint64) ([]model.UserDao, error)
}

// Controller implements user methods on logic layer.
type Controller struct {
	ur         userRepo
	tracer     trace.Tracer
	encryption string
}

// New creates new Controller.
func New(ur userRepo, tracer trace.Tracer, encryption string) *Controller {
	return &Controller{
		ur:         ur,
		tracer:     tracer,
		encryption: encryption,
	}
}
