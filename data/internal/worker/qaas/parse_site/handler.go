package parse_site

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"go.dataddo.com/pgq"
)

type Handler struct {
	siteStore     siteStore
	pagePublisher publisher
}

func New(
	siteStore siteStore,
	pagePublisher publisher,
) *Handler {
	return &Handler{
		siteStore:     siteStore,
		pagePublisher: pagePublisher,
	}
}

// TODO: add count of currently running jobs
func (h Handler) Handle(ctx context.Context, msg *pgq.MessageIncoming) (bool, error) {
	//qaas.SiteJob
	var job qaas.SiteJob
	err := json.Unmarshal(msg.Payload, &job)
	if err != nil {
		return true, fmt.Errorf("failed to unmarshal parsesite payload: %w", err)
	}

	slog.Info("handled site job", "job", job)
	currentSite := job.Payload
	if err := h.siteStore.Save(ctx, currentSite); err != nil {
		slog.Error("failed to save site", "site", currentSite, "error", err)
		return true, err
	}

	indexPage, err := site.NewPage(currentSite.ID, currentSite.URL)
	if err != nil {
		slog.Error("failed to create index page", "site", currentSite, "error", err)
		return true, err
	}

	publishOptions := []qaas.PublishOption{
		qaas.WithQueue(qaas.ParsePageQueue),
	}

	_, err = h.pagePublisher.Publish(ctx, []any{qaas.PageJob{
		Payload: indexPage,
		Delay:   0,
	}}, publishOptions...)
	if err != nil {
		slog.Error("failed to publish page job", "page", indexPage, "error", err)
		return true, err
	}

	return true, nil
}
