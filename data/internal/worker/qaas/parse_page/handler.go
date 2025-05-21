package parse_page

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"go.dataddo.com/pgq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	pageStore   pageStore
	pageService pageService
	publisher   publisher
	tracer      trace.Tracer
}

func New(
	pageStore pageStore,
	pageService pageService,
	publisher publisher,
	tracer trace.Tracer,
) *Handler {
	return &Handler{
		pageService: pageService,
		pageStore:   pageStore,
		publisher:   publisher,
		tracer:      tracer,
	}
}

func (h Handler) Handle(ctx context.Context, msg *pgq.MessageIncoming) (bool, error) {
	ctx, span := h.tracer.Start(ctx, "parse_page.Handle")
	defer span.End()
	// job qaas.PageJob
	var job qaas.PageJob
	err := json.Unmarshal(msg.Payload, &job)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal parsepage payload: %w", err)
		span.RecordError(err)
		return true, err
	}

	if job.Metadata == nil {
		err = fmt.Errorf("failed to get siteJobID from job")
		span.RecordError(err)
		return true, err
	}
	siteJobID, ok := job.Metadata["siteJobID"]
	if !ok || siteJobID == "" {
		err = fmt.Errorf("failed to get siteJobID from job")
		span.RecordError(err)
		return true, err
	}

	slog.Debug("handled page job", "jon", job)
	page := job.Payload
	if page == nil {
		err = fmt.Errorf("failed to get page from job")
		span.RecordError(err)
		return true, err
	}
	_, _, err = h.pageService.ParsePage(ctx, page, siteJobID.(string))
	if err != nil {
		err = fmt.Errorf("failed to handle page job: %w", err)
		span.RecordError(err,
			trace.WithAttributes(
				attribute.String("siteJobID", siteJobID.(string)),
				attribute.String("pageID", page.ID),
			),
		)
		return true, fmt.Errorf("failed to handle page job: %w", err)
	}

	publishOptions := []qaas.PublishOption{
		qaas.WithQueue(qaas.ParsePageResultQueue),
		qaas.WithSourceQueue(qaas.ParsePageQueue),
	}
	_, err = h.publisher.Publish(ctx, []any{qaas.PageResultJob{
		Payload: page,
		Delay:   0,
	}}, publishOptions...)
	if err != nil {
		err = fmt.Errorf("failed to publish result message for page: %w", err)
		span.RecordError(err)
		slog.Error("failed to publish result message for page", "err", err, "pageID", page.ID)
	}

	return true, nil
}
