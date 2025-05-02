package embed_document

import (
	"context"
	"fmt"
	"strings"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
)

type Handler struct {
	embeddingService embeddingService
	pageStore        pageStore
	siteStore        siteStore
}

func New(embeddingService embeddingService, pageStore pageStore, siteStore siteStore) *Handler {
	return &Handler{
		embeddingService: embeddingService,
		pageStore:        pageStore,
		siteStore:        siteStore,
	}
}

func (h Handler) Handle(ctx context.Context, job messages.ResultMessage) error {
	switch job.Type {
	case messages.WebResult:
		page, err := h.pageStore.GetByID(ctx, job.ObjID)
		if err != nil {
			return fmt.Errorf("failed to get page in embed_document: %w", err)
		}
		if page == nil {
			return fmt.Errorf("page not found in embed_document")
		}
		site, err := h.siteStore.GetByID(ctx, page.SiteID)
		if err != nil {
			return fmt.Errorf("failed to get site in embed_document: %w", err)
		}
		err = h.embeddingService.Process(
			ctx,
			strings.NewReader(page.Raw),
			document.HTML,
			site.SourceID,
			string(job.Type),
			page.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to process page in embed_document: %w", err)
		}
		return nil
	case messages.FileResult:
		return nil
	default:
		return fmt.Errorf("unknown message type in embed_document: %s", job.Type)
	}
}
