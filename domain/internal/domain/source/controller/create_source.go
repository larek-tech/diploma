package controller

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/source/model"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CreateSource creates new source record.
func (ctrl *Controller) CreateSource(ctx context.Context, req *pb.CreateSourceRequest, meta *authpb.UserAuthMetadata) (*pb.Source, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CreateSource",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("title", req.GetTitle()),
			attribute.Int("source type", int(req.GetTyp())),
		),
	)
	defer span.End()

	source := model.SourceDao{
		UserID:      meta.GetUserId(),
		Title:       req.GetTitle(),
		Content:     req.GetContent(),
		Type:        model.SourceType(req.GetTyp()),
		Credentials: req.GetCredentials(),
		Status:      model.StatusParsing,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	source.FillUpdateParams(req.GetUpdateParams())

	sourceID, err := ctrl.sr.InsertSource(ctx, source)
	if err != nil {
		return nil, errs.WrapErr(err)
	}
	source.ID = sourceID

	if err = ctrl.saveSourceData(ctx, source, meta); err != nil {
		return nil, errs.WrapErr(err, "save source data")
	}

	return source.ToProto(), nil
}

func (ctrl *Controller) saveSourceData(ctx context.Context, source model.SourceDao, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.saveSourceData")
	defer span.End()

	dataMsg := model.DataMessage{
		Title:        source.Title,
		Content:      source.Content,
		Type:         source.Type,
		Credentials:  source.Credentials,
		UpdateParams: source.AssembleUpdateParams(),
	}

	data, err := json.Marshal(dataMsg)
	if err != nil {
		return errs.WrapErr(err, "marshal data message for kafka")
	}

	ctrl.producer.SendAsyncMessage(&sarama.ProducerMessage{
		Topic:     sourceTopic,
		Key:       sarama.StringEncoder(strconv.FormatInt(source.ID, 10)),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	})

	go ctrl.getParsingResult(source, meta)

	return nil
}

func (ctrl *Controller) getParsingResult(source model.SourceDao, meta *authpb.UserAuthMetadata) {
	ctx, span := ctrl.tracer.Start(
		context.Background(),
		"Controller.getParsingResult",
		trace.WithAttributes(
			attribute.Int64("source internal id", source.ID),
			attribute.String("title", source.Title),
			attribute.Int("source type", int(source.Type)),
		),
	)
	defer span.End()

	var (
		startedParsing bool
		failed         bool
	)

	for {
		select {
		case err := <-ctrl.errCh:
			log.Err(errs.WrapErr(err)).Msg("consuming source status topic failed")
		case msg := <-ctrl.statusCh:
			if string(msg.Key) != strconv.FormatInt(source.ID, 10) {
				continue
			}

			var (
				resp   model.ParsingStatus
				status model.SourceStatus
			)

			if err := json.Unmarshal(msg.Value, &resp); err != nil {
				log.Err(errs.WrapErr(err)).Msg("unmarshal parsing status message")
				status = model.StatusFailed
			} else {
				status = resp.Status
			}

			if status == model.StatusFailed {
				failed = true
			}

			if (!startedParsing && status == model.StatusParsing) || failed {
				startedParsing = true

				source.Status = status
				if err := ctrl.sr.UpdateSource(ctx, source, meta.GetUserId(), meta.GetRoles()); err != nil {
					log.Warn().Err(errs.WrapErr(err)).Msg("can't properly update source status")
				}
			}

			if failed {
				log.Err(errs.WrapErr(ErrUpdateSourceStatus)).Msg("get parsing result")
				return
			}
		}
	}
}
