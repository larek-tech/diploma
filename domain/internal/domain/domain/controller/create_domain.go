package controller

import (
	"context"
	"time"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CreateDomain creates new domain record.
func (ctrl *Controller) CreateDomain(ctx context.Context, req *pb.CreateDomainRequest, meta *authpb.UserAuthMetadata) (*pb.Domain, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CreateDomain",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("title", req.GetTitle()),
			attribute.Int64Slice("source IDs", req.GetSourceIds()),
		),
	)
	defer span.End()

	domain := model.DomainDao{
		Title:     req.GetTitle(),
		UserID:    meta.GetUserId(),
		SourceIDs: req.GetSourceIds(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	domainID, err := ctrl.dr.InsertDomain(ctx, domain)
	if err != nil {
		return nil, errs.WrapErr(err)
	}
	domain.ID = domainID

	return domain.ToProto(), nil
}
