package parse_site

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/samber/lo"
	"go.dataddo.com/pgq"
)

const (
	StatusDelay = time.Second * 10
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

func (h Handler) Handle(ctx context.Context, msg *pgq.MessageIncoming) (bool, error) {
	//qaas.SiteJob
	var job qaas.SiteJob
	err := json.Unmarshal(msg.Payload, &job)
	if err != nil {
		return true, fmt.Errorf("failed to unmarshal parsesite payload: %w", err)
	}
	if job.Metadata == nil {
		return true, fmt.Errorf("failed to get siteJobID from job")
	}
	siteJobID, ok := job.Metadata["siteJobID"]
	if !ok || siteJobID == "" {
		return true, fmt.Errorf("failed to get siteJobID from job")
	}
	externalKey, ok := job.Metadata["externalKey"]
	if !ok {
		return true, fmt.Errorf("missing external key in Metadata")
	}

	slog.Info("handled site job", "job", job)
	currentSite := job.Payload

	if err := h.siteStore.Save(ctx, currentSite); err != nil {
		slog.Error("failed to save site", "site", currentSite, "error", err)
		return true, err
	}
	metadata := map[string]any{
		"siteJobID":   siteJobID,
		"externalKey": externalKey.(string),
	}
	parseJobs := lo.Map(currentSite.AvailablePages, func(url string, _ int) any {
		page, mapErr := site.NewPage(currentSite.ID, url)
		if mapErr != nil {
			slog.Error("failed to create page", "site", currentSite, "error", mapErr)
		}

		return qaas.PageJob{
			Payload:  page,
			Delay:    0,
			Metadata: metadata,
		}
	})

	publishOptions := []qaas.PublishOption{
		qaas.WithQueue(qaas.ParsePageQueue),
	}

	parsePageJobIDs, err := h.pagePublisher.Publish(ctx, parseJobs, publishOptions...)
	if err != nil {
		slog.Error("failed to publish page job", "page", currentSite.AvailablePages, "error", err)
		return true, err
	}

	publishOptions = []qaas.PublishOption{
		qaas.WithQueue(qaas.ParseSiteStatusQueue),
	}

	statusJob := qaas.ParseStatusJob{
		ExternalKey:      externalKey.(string),
		SiteID:           currentSite.ID,
		SourceID:         currentSite.SourceID,
		ParsePageJobsIDs: parsePageJobIDs,
		Delay:            StatusDelay,
		SiteJobID:        siteJobID.(string),
	}

	_, err = h.pagePublisher.Publish(ctx, []any{statusJob}, publishOptions...)
	if err != nil {
		slog.Error("failed to publish status job", "site", currentSite, "error", err)
		return true, err
	}

	return true, nil
}
