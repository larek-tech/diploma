package controller

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	mlpb "github.com/larek-tech/diploma/chat/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// ProcessQuery start query processing with streaming response.
func (ctrl *Controller) ProcessQuery(
	ctx context.Context,
	req *pb.ProcessQueryRequest,
	out chan *pb.ChunkedResponse,
	errCh chan error,
) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ProcessQuery",
		trace.WithAttributes(attribute.Int64("userID", req.GetUserId())),
	)
	defer span.End()

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		errCh <- errs.WrapErr(err, "invalid chat id")
		return
	}

	creatorID, err := ctrl.cr.GetChatUserID(ctx, chatID)
	if err != nil {
		errCh <- errs.WrapErr(err)
		return
	}

	if creatorID != req.GetUserId() {
		errCh <- errs.WrapErr(ErrNoAccessToChat, "process query")
		return
	}

	processCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var queryMeta []byte
	if m := req.GetMetadata(); m == nil {
		queryMeta = []byte("{}")
	} else {
		queryMeta = m
	}
	q := model.QueryDao{
		UserID:     req.GetUserId(),
		ChatID:     chatID,
		Content:    req.GetContent(),
		DomainID:   req.GetDomainId(),
		SourceIDs:  req.GetSourceIds(),
		ScenarioID: req.GetScenarioId(),
		Metadata:   queryMeta,
	}
	queryID, err := ctrl.cr.InsertQuery(ctx, q)
	if err != nil {
		errCh <- errs.WrapErr(err, "insert query")
		return
	}

	ctrl.mu.Lock()
	ctrl.processing[queryID] = cancel
	ctrl.mu.Unlock()

	var scenario *mlpb.Scenario = nil
	if q.Metadata != nil {
		scenario = &mlpb.Scenario{}
		if err = json.Unmarshal(q.Metadata, &scenario); err != nil {
			errCh <- errs.WrapErr(err, "get query scenario from metadata")
			return
		}
	}

	respCreate := model.ResponseDao{
		QueryID: queryID,
		ChatID:  chatID,
		Status:  model.StatusCreated,
	}
	respID, err := ctrl.cr.InsertResponse(ctx, respCreate)
	if err != nil {
		errCh <- errs.WrapErr(err)
		return
	}

	resp, err := ctrl.cr.GetResponseByID(ctx, respID)
	if err != nil {
		errCh <- errs.WrapErr(err)
		return
	}

	defer func() {
		if err == nil {
			return
		}
		ctx, span := ctrl.tracer.Start(
			context.Background(),
			"Controller.SetResponseStatusError",
			trace.WithAttributes(
				attribute.Int64("queryID", resp.QueryID),
				attribute.String("chatID", resp.ChatID.String()),
			),
		)
		defer span.End()

		log.Err(errs.WrapErr(err)).Msg("processing query")
		resp.Status = model.StatusError

		span.SetAttributes(attribute.Int("status", int(resp.Status)))

		if e := ctrl.cr.UpdateResponse(ctx, resp); e != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("set response status error")
		}
	}()

	mlReq := &mlpb.ProcessQueryRequest{
		Query: &mlpb.Query{
			Id:      queryID,
			UserId:  q.UserID,
			Content: q.Content,
		},
		Scenario:  scenario,
		SourceIds: q.SourceIDs,
	}
	stream, err := ctrl.mlService.ProcessQuery(ctx, mlReq)
	if err != nil {
		errCh <- errs.WrapErr(err, "start stream")
		return
	}

	var (
		contentBuff        strings.Builder
		streamCtx, success = context.WithCancel(ctx)
	)
	defer success()

	for {
		select {
		case <-processCtx.Done():
			return
		case <-streamCtx.Done():
			log.Debug().
				Int64("queryID", resp.QueryID).
				Any("sourceIDs", q.SourceIDs).
				Msg("got chunk")
			log.Info().Int64("queryID", queryID).Msg("stream successfully finished")
			return
		default:
			if err = ctrl.receiveChunk(streamCtx, success, stream, out, &resp, &contentBuff); err != nil {
				errCh <- errs.WrapErr(err)
				return
			}
		}
	}
}

func (ctrl *Controller) receiveChunk(
	ctx context.Context,
	success context.CancelFunc,
	stream grpc.ServerStreamingClient[mlpb.ProcessQueryResponse],
	out chan *pb.ChunkedResponse,
	resp *model.ResponseDao,
	buff *strings.Builder,
) error {
	r, err := stream.Recv()
	if err == io.EOF {
		resp.Status = model.StatusSuccess
		if err = ctrl.cr.UpdateResponse(ctx, *resp); err != nil {
			return errs.WrapErr(err, "process last chunk")
		}
		success()
		return nil
	}

	if err != nil {
		return errs.WrapErr(err, "streaming error")
	}

	content := r.GetChunk().Content

	sourceIDs := r.GetSourceIds()
	if sourceIDs != nil {
		content += "\nИсточники: [" + strings.Join(sourceIDs, ", ") + "]"
	}

	_, err = buff.WriteString(content)
	if err != nil {
		return errs.WrapErr(err, "write chunk")
	}

	log.Debug().Int64("queryID", resp.QueryID).Str("content", content).Msg("got chunk")

	out <- &pb.ChunkedResponse{
		QueryId:   resp.QueryID,
		Content:   content,
		SourceIds: sourceIDs,
	}

	resp.Content = buff.String()
	resp.Status = model.StatusProcessing
	if err = ctrl.cr.UpdateResponse(ctx, *resp); err != nil {
		return errs.WrapErr(err, "append chunk to response")
	}

	return nil
}
