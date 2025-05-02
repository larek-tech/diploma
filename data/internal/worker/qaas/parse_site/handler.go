package parse_site

import (
	"context"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
)

type Handler struct {
	siteStore     siteStore
	pagePublisher pagePublisher
}

func New(
	siteStore siteStore,
	pagePublisher pagePublisher,
) *Handler {
	return &Handler{
		siteStore:     siteStore,
		pagePublisher: pagePublisher,
	}
}

func (h Handler) Handle(ctx context.Context, job messages.SiteJob) error {
	slog.Info("handled site job", "job", job)
	currentSite := job.Payload
	if err := h.siteStore.Save(ctx, currentSite); err != nil {
		slog.Error("failed to save site", "site", currentSite, "error", err)
		return err
	}

	indexPage, err := site.NewPage(currentSite.ID, currentSite.URL)
	if err != nil {
		slog.Error("failed to create index page", "site", currentSite, "error", err)
		return err
	}

	err = h.pagePublisher.Publish(ctx, messages.PageJob{
		Type:    messages.ParsePage,
		Payload: indexPage,
		Delay:   0,
	})
	if err != nil {
		slog.Error("failed to publish page job", "page", indexPage, "error", err)
		return err
	}

	return nil
}
