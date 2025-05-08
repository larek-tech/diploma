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

func (ctrl *Controller) ProcessQuery(
	ctx context.Context,
	req *pb.ProcessQueryRequest,
	out chan *pb.ChunkedResponse,
	errCh chan error,
	cancel chan struct{},
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

	q := model.QueryDao{
		UserID:     req.GetUserId(),
		ChatID:     chatID,
		Content:    req.GetContent(),
		DomainID:   req.GetDomainId(),
		SourceIDs:  req.GetSourceIds(),
		ScenarioID: req.GetScenarioId(),
		Metadata:   req.GetMetadata(),
	}
	queryID, err := ctrl.cr.InsertQuery(ctx, q)
	if err != nil {
		errCh <- errs.WrapErr(err, "insert query")
		return
	}

	scenario := mlpb.Scenario{}
	if err = json.Unmarshal(q.Metadata, &scenario); err != nil {
		errCh <- errs.WrapErr(err, "get query scenario from metadata")
		return
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

		errCh <- errs.WrapErr(err, "processing query")
		resp.Status = model.StatusError

		if e := ctrl.cr.UpdateResponse(ctx, resp); e != nil {
			e = errs.WrapErr(e)
			log.Warn().Err(e).Msg("set response status error")
			errCh <- e
		}
	}()

	mlReq := &mlpb.ProcessQueryRequest{
		Query: &mlpb.Query{
			Id:      queryID,
			UserId:  q.UserID,
			Content: q.Content,
		},
		Scenario:  &scenario,
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
		case <-cancel:
			if err = ctrl.cancel(ctx, &resp); err != nil {
				errCh <- errs.WrapErr(err)
			} else {
				log.Info().Int64("queryID", queryID).Msg("processing canceled")
			}
			return
		case <-streamCtx.Done():
			log.Info().Int64("queryID", queryID).Msg("processing successfully finished")
			out <- &pb.ChunkedResponse{
				QueryId:   queryID,
				SourceIds: q.SourceIDs,
			}
			return
		default:
			if err = ctrl.receiveChunk(streamCtx, success, stream, out, &resp, &contentBuff); err != nil {
				errCh <- errs.WrapErr(err)
				return
			}
		}
	}
}

func (ctrl *Controller) cancel(ctx context.Context, resp *model.ResponseDao) error {
	resp.Status = model.StatusCanceled
	if err := ctrl.cr.UpdateResponse(ctx, *resp); err != nil {
		return errs.WrapErr(err, "set response status cancel")
	}
	return nil
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
	_, err = buff.WriteString(content)
	if err != nil {
		return errs.WrapErr(err, "write chunk")
	}

	curContent := buff.String()
	out <- &pb.ChunkedResponse{
		QueryId:   resp.QueryID,
		Content:   curContent,
		SourceIds: nil,
	}

	resp.Content = curContent
	resp.Status = model.StatusProcessing
	if err = ctrl.cr.UpdateResponse(ctx, *resp); err != nil {
		return errs.WrapErr(err, "append chunk to response")
	}

	return nil
}
