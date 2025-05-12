package controller

import (
	"context"
	"time"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateDomain updates domain data.
func (ctrl *Controller) UpdateDomain(ctx context.Context, req *pb.UpdateDomainRequest, meta *authpb.UserAuthMetadata) (*pb.Domain, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdateDomain",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("domainID", req.GetDomainId()),
			attribute.String("title", req.GetTitle()),
		),
	)
	defer span.End()

	domain, err := ctrl.dr.GetDomainByID(ctx, req.GetDomainId(), meta.GetUserId(), meta.GetRoles())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	domain.Title = req.GetTitle()
	domain.SourceIDs = req.GetSourceIds()
	domain.ScenarioIds = req.GetScenarioIds()
	domain.UpdatedAt = time.Now()

	if err = ctrl.dr.UpdateDomain(ctx, domain, meta.GetUserId(), meta.GetRoles()); err != nil {
		return nil, errs.WrapErr(err)
	}

	return domain.ToProto(), nil
}
