package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/chat/internal/auth/pb"
	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CancelProcessing cancels processing context and dependant jobs.
func (ctrl *Controller) CancelProcessing(ctx context.Context, req *pb.CancelProcessingRequest, meta *authpb.UserAuthMetadata) error {
	queryID := req.GetQueryId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CancelProcessing",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("queryID", queryID),
		),
	)
	defer span.End()

	resp, err := ctrl.cr.GetResponseByQueryID(ctx, queryID)
	if err != nil {
		return errs.WrapErr(err)
	}

	creatorID, err := ctrl.cr.GetChatUserID(ctx, resp.ChatID)
	if err != nil {
		return errs.WrapErr(err)
	}

	if meta.GetUserId() != creatorID {
		return errs.WrapErr(ErrNoAccessToChat, "cancel processing query")
	}

	resp.Status = model.StatusCanceled
	if err = ctrl.cr.UpdateResponse(ctx, resp); err != nil {
		return errs.WrapErr(err, "set response status cancel")
	}

	ctrl.mu.Lock()
	cancel := ctrl.processing[queryID]
	cancel()
	delete(ctrl.processing, queryID)
	ctrl.mu.Unlock()

	return nil
}
