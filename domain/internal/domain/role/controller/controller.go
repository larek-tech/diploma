package controller

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/role/model"
	"go.opentelemetry.io/otel/trace"
)

type roleRepo interface {
	InsertRole(ctx context.Context, u model.RoleDao) (int64, error)
	GetRole(ctx context.Context, id int64) (model.RoleDao, error)
	DeleteRole(ctx context.Context, id int64) error
	ListRoles(ctx context.Context, offset, limit uint64) ([]model.RoleDao, error)
	SetRole(ctx context.Context, userID, roleID int64) error
	RemoveRole(ctx context.Context, userID, roleID int64) error
}

// Controller implements role methods on logic layer.
type Controller struct {
	rr     roleRepo
	tracer trace.Tracer
}

// New creates new controller.
func New(rr roleRepo, tracer trace.Tracer) *Controller {
	return &Controller{
		rr:     rr,
		tracer: tracer,
	}
}
