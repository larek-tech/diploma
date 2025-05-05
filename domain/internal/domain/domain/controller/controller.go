package controller

import (
	"context"
	"errors"
	"slices"

	"github.com/larek-tech/diploma/domain/internal/auth"
	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrNoAccessToDomain is an error when user can't edit domain.
	ErrNoAccessToDomain = errors.New("user has no access to edit domain")
)

type domainRepo interface {
	InsertDomain(ctx context.Context, d model.DomainDao) (int64, error)
	GetDomainByID(ctx context.Context, id, userID int64, roleIDs []int64) (model.DomainDao, error)
	UpdateDomain(ctx context.Context, d model.DomainDao, userID int64, roleIDs []int64) error
	DeleteDomain(ctx context.Context, id, userID int64, roleIDs []int64) error
	ListDomains(ctx context.Context, userID int64, roleIDs []int64, offset, limit uint64) ([]model.DomainDao, error)
	GetPermittedUsers(ctx context.Context, domainID int64) ([]int64, error)
	GetPermittedRoles(ctx context.Context, domainID int64) ([]int64, error)
	UpdatePermittedUsers(ctx context.Context, domainID int64, userIDs []int64) ([]int64, error)
	UpdatePermittedRoles(ctx context.Context, domainID int64, roleIDs []int64) ([]int64, error)
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

func (ctrl *Controller) checkDomainCreator(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) error {
	userID := meta.GetUserId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.checkDomainCreator",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("domainID", domainID),
		),
	)
	defer span.End()

	roles := meta.GetRoles()
	if !slices.Contains(roles, auth.AdminUserID) {
		domain, err := ctrl.dr.GetDomainByID(ctx, domainID, userID, roles)
		if err != nil {
			return errs.WrapErr(err)
		}

		if domain.UserID != userID {
			return errs.WrapErr(ErrNoAccessToDomain, "check domain creator")
		}
	}
	return nil
}
