package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListSourcesByDomain returns paginated list of sources by specified domain.
func (ctrl *Controller) ListSourcesByDomain(
	ctx context.Context,
	req *pb.ListSourcesByDomainRequest,
	meta *authpb.UserAuthMetadata,
) (*pb.ListSourcesResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"ListSourcesByDomain",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("domainID", req.GetDomainId()),
			attribute.Int64("offset", int64(req.GetOffset())),
			attribute.Int64("limit", int64(req.GetLimit())),
		),
	)
	defer span.End()

	sourcesDB, err := ctrl.sr.ListSourcesByDomain(
		ctx,
		meta.GetUserId(),
		req.GetDomainId(),
		meta.GetRoles(),
		req.GetOffset(),
		req.GetLimit(),
	)
	if err != nil {
		return nil, errs.WrapErr(err, "list sources by domain")
	}

	sources := make([]*pb.Source, len(sourcesDB))
	for idx := range sourcesDB {
		sources[idx] = sourcesDB[idx].ToProto()
	}

	return &pb.ListSourcesResponse{
		Sources: sources,
	}, nil
}
