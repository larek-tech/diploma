package parse_page

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"go.dataddo.com/pgq"
)

type Handler struct {
	pageStore   pageStore
	pageService pageService
	publisher   publisher
}

func New(
	pageStore pageStore,
	pageService pageService,
	publisher publisher,
) *Handler {
	return &Handler{
		pageService: pageService,
		pageStore:   pageStore,
		publisher:   publisher,
	}
}

func (h Handler) Handle(ctx context.Context, msg *pgq.MessageIncoming) (bool, error) {
	// job qaas.PageJob
	var job qaas.PageJob
	err := json.Unmarshal(msg.Payload, &job)
	if err != nil {
		return true, fmt.Errorf("failed to unmarshal parsepage payload: %w", err)
	}

	if job.Metadata == nil {
		return true, fmt.Errorf("failed to get siteJobID from job")
	}
	siteJobID, ok := job.Metadata["siteJobID"]
	if !ok || siteJobID == "" {
		return true, fmt.Errorf("failed to get siteJobID from job")
	}

	// FIXME: find place with empty uuid and how to prevent that?
	slog.Info("handled page job", "jon", job)
	page := job.Payload
	if page == nil {
		return true, fmt.Errorf("failed to get page from job")
	}
	_, _, err = h.pageService.ParsePage(ctx, page, siteJobID.(string))
	if err != nil {
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
		slog.Error("failed to publish result message for page", "err", err, "pageID", page.ID)
	}

	return true, nil
}
