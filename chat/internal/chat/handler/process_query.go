package handler

import (
	"context"
	"time"

	"github.com/larek-tech/diploma/chat/internal/auth"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProcessQuery starts query processing pipeline.
func (h *Handler) ProcessQuery(req *pb.ProcessQueryRequest, stream grpc.ServerStreamingServer[pb.ChunkedResponse]) error {
	ctx, span := h.tracer.Start(
		context.Background(),
		"Handler.ProcessQuery",
		trace.WithAttributes(
			attribute.Int64("userID", req.GetUserId()),
			attribute.String("chatID", req.GetChatId()),
			attribute.StringSlice("sourceIDs", req.GetSourceIds()),
		),
	)
	defer span.End()

	if len(req.GetSourceIds()) == 0 {
		log.Err(errs.WrapErr(ErrEmptySourceIDs)).Msg("process query")
		return status.Error(codes.InvalidArgument, "at least 1 source id is required")
	}

	_, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return status.Error(codes.Unauthenticated, "unauthorized")
	}

	out := make(chan *pb.ChunkedResponse)
	errCh := make(chan error)
	defer close(out)
	defer close(errCh)

	go h.cc.ProcessQuery(ctx, req, out, errCh)

	for {
		select {
		case chunk := <-out:
			if err = stream.Send(chunk); err != nil {
				log.Err(errs.WrapErr(err)).Msg("process query")
				return status.Error(codes.Internal, "failed sending response chunk in stream")
			}
			if chunk.GetSourceIds() != nil {
				return status.Error(codes.OK, "processed query successfully")
			}
		case e := <-errCh:
			log.Err(errs.WrapErr(e)).Msg("process query")
			return status.Errorf(codes.Internal, "failed processing query")
		case <-time.After(time.Minute):
			log.Err(ErrProcessQueryTimeout).Msg("process query timeout")
			return status.Error(codes.DeadlineExceeded, "processing query timeout")
		}
	}
}
