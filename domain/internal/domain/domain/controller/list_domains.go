package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListDomains returns the paginated list of available domains.
func (ctrl *Controller) ListDomains(ctx context.Context, req *pb.ListDomainsRequest, meta *authpb.UserAuthMetadata) (*pb.ListDomainsResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListDomains",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("offest", int64(req.GetOffset())),
			attribute.Int64("limit", int64(req.GetLimit())),
		),
	)
	defer span.End()

	domainsDao, err := ctrl.dr.ListDomains(ctx, meta.GetUserId(), meta.GetRoles(), req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	domains := make([]*pb.Domain, len(domainsDao))
	for idx := range domainsDao {
		domains[idx] = domainsDao[idx].ToProto()
	}

	return &pb.ListDomainsResponse{Domains: domains}, nil
}
