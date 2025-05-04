package parse_page

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/samber/lo"
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

	// FIXME: find place with empty uuid and how to prevent that?
	slog.Info("handled page job", "jon", job)
	page := job.Payload
	if page == nil {
		return true, fmt.Errorf("failed to get page from job")
	}
	outgoingPages, _, err := h.pageService.ParsePage(ctx, page)
	if err != nil {
		return true, fmt.Errorf("failed to handle page job: %w", err)
	}
	pagesToParse := lo.Map(outgoingPages, func(page *site.Page, index int) any {
		return qaas.PageJob{
			Payload: page,
		}
	})
	publishOptions := []qaas.PublishOption{
		qaas.WithQueue(qaas.ParsePageQueue),
	}

	_, err = h.publisher.Publish(ctx, pagesToParse, publishOptions...)
	if err != nil {
		slog.Error("failed to publish outgoing pages to qaas: %w", err)
	}

	publishOptions = []qaas.PublishOption{
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
