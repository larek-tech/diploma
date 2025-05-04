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

// UpdateSource updates source data.
func (ctrl *Controller) UpdateSource(ctx context.Context, req *pb.UpdateSourceRequest, meta *authpb.UserAuthMetadata) (*pb.Source, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdateSource",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("sourceID", req.GetSourceId()),
			attribute.String("title", req.GetTitle()),
		),
	)
	defer span.End()

	source, err := ctrl.sr.GetSourceByID(ctx, req.GetSourceId(), meta.GetUserId(), meta.GetRoles())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	source.Title = req.GetTitle()
	source.Content = req.GetContent()
	source.Credentials = req.GetCredentials()
	source.FillUpdateParams(req.UpdateParams)
	source.UpdatedAt = time.Now()

	if err = ctrl.sr.UpdateSource(ctx, source, meta.GetUserId(), meta.GetRoles()); err != nil {
		return nil, errs.WrapErr(err)
	}

	return source.ToProto(), nil
}
