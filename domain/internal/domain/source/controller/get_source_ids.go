package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetSourceIDs returns external source ids by list of internal ids.
func (ctrl *Controller) GetSourceIDs(ctx context.Context, req *pb.GetSourceIDsRequest, meta *authpb.UserAuthMetadata) (*pb.GetSourceIDsResponse, error) {
	sourceIDs := req.GetSourceIds()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetSourceIDs",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64Slice("sourceIDs", sourceIDs),
		),
	)
	defer span.End()

	sourceUUIDs := make([]string, len(sourceIDs))
	for idx := range sourceIDs {
		externalID, err := ctrl.sr.GetSourceIDs(ctx, sourceIDs[idx], meta.GetUserId(), meta.GetRoles())
		if err != nil {
			return nil, errs.WrapErr(err)
		}
		sourceUUIDs[idx] = externalID.String()
	}

	return &pb.GetSourceIDsResponse{SourceIds: sourceUUIDs}, nil
}
