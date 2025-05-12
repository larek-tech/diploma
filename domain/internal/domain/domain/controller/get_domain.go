package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetDomain returns domain by id.
func (ctrl *Controller) GetDomain(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) (*pb.Domain, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetDomain",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("domainID", domainID),
		),
	)
	defer span.End()

	domain, err := ctrl.dr.GetDomainByID(ctx, domainID, meta.GetUserId(), meta.GetRoles())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return domain.ToProto(), nil
}
