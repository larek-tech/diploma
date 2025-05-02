package parse_page

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
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

func (h Handler) Handle(ctx context.Context, job messages.PageJob) error {
	// FIXME: find place with empty uuid and how to prevent that?
	slog.Info("handled page job", "jon", job)
	page := job.Payload
	if page == nil {
		return fmt.Errorf("failed to get page from job")
	}
	outgoingPages, _, err := h.pageService.ParsePage(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to handle page job: %w", err)
	}
	for _, nextPage := range outgoingPages {
		err = h.publisher.Publish(ctx, messages.PageJob{
			Type:    messages.ParsePage,
			Payload: nextPage,
		})
		if err != nil {
			slog.Error("failed to publish outgoing page", "err", err)
		}
	}
	err = h.publisher.Publish(ctx, messages.ResultMessage{
		Type:  messages.WebResult,
		ObjID: page.ID,
	})
	if err != nil {
		slog.Error("failed to publish result message", "err", err)
	}

	return nil
}
