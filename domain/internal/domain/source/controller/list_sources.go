package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListSources returns the paginated list of available sources.
func (ctrl *Controller) ListSources(ctx context.Context, req *pb.ListSourcesRequest, meta *authpb.UserAuthMetadata) (*pb.ListSourcesResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListSources",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("offest", int64(req.GetOffset())),
			attribute.Int64("limit", int64(req.GetLimit())),
		),
	)
	defer span.End()

	sourcesDao, err := ctrl.sr.ListSources(ctx, meta.GetUserId(), meta.GetRoles(), req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	sources := make([]*pb.Source, len(sourcesDao))
	for idx := range sourcesDao {
		sources[idx] = sourcesDao[idx].ToProto()
	}

	return &pb.ListSourcesResponse{Sources: sources}, nil
}
