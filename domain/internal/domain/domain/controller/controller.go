package controller

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"go.opentelemetry.io/otel/trace"
)

type domainRepo interface {
	InsertDomain(ctx context.Context, d model.DomainDao) (int64, error)
	GetDomainByID(ctx context.Context, id, userID int64, roleIDs []int64) (model.DomainDao, error)
	UpdateDomain(ctx context.Context, d model.DomainDao, userID int64, roleIDs []int64) error
	DeleteDomain(ctx context.Context, id, userID int64, roleIDs []int64) error
	ListDomains(ctx context.Context, userID int64, roleIDs []int64, offset, limit uint64) ([]model.DomainDao, error)
}

// Controller implements domain methods on logic layer.
type Controller struct {
	dr     domainRepo
	tracer trace.Tracer
}

// New creates new Controller.
func New(dr domainRepo, tracer trace.Tracer) *Controller {
	return &Controller{
		dr:     dr,
		tracer: tracer,
	}
}
