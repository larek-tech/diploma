package embed_document

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"go.dataddo.com/pgq"
)

type Handler struct {
	embeddingService embeddingService
	pageStore        pageStore
	siteStore        siteStore
	fileStore        fileStore
}

func New(embeddingService embeddingService, pageStore pageStore, siteStore siteStore, fileStore fileStore) *Handler {
	return &Handler{
		embeddingService: embeddingService,
		pageStore:        pageStore,
		siteStore:        siteStore,
		fileStore:        fileStore,
	}
}

func (h Handler) Handle(ctx context.Context, msg *pgq.MessageIncoming) (bool, error) {
	objType, ok := msg.Metadata["sourceQueue"]
	if !ok {
		return false, nil
	}
	switch qaas.Queue(objType) {
	case qaas.ParsePageQueue:
		var job qaas.PageJob
		if err := json.Unmarshal(msg.Payload, &job); err != nil {
			return true, fmt.Errorf("failed to unmarshal PageJob: %w", err)
		}
		page, err := h.pageStore.GetByID(ctx, job.Payload.ID)
		if err != nil {
			return true, fmt.Errorf("failed to get page in embed_document: %w", err)
		}
		if page == nil {
			return true, fmt.Errorf("page not found in embed_document")
		}
		site, err := h.siteStore.GetByID(ctx, page.SiteID)
		if err != nil {
			return true, fmt.Errorf("failed to get site in embed_document: %w", err)
		}
		err = h.embeddingService.Process(
			ctx,
			strings.NewReader(page.Raw),
			document.HTML,
			page,
			site.SourceID,
		)
		if err != nil {
			return true, fmt.Errorf("failed to process page in embed_document: %w", err)
		}
		return true, nil
	case qaas.ParseFileQueue:
		var job qaas.FileJob
		if err := json.Unmarshal(msg.Payload, &job); err != nil {
			return true, fmt.Errorf("failed to unmarshal FileJob: %w", err)
		}
		file, err := h.fileStore.GetByID(ctx, job.Payload.ID)
		if err != nil {
			return true, fmt.Errorf("failed o get file from db: %w", err)
		}
		if file == nil {
			return true, fmt.Errorf("got empty file from storage: %v", job)
		}
		ext, err := getFileExt(file.Filename)
		if err != nil {
			return true, fmt.Errorf("failed to determine file extension")
		}
		err = h.embeddingService.Process(
			ctx,
			bytes.NewReader(file.Raw),
			ext,
			file,
			file.SourceID,
		)
		if err != nil {
			return true, fmt.Errorf("failed to process file in embed_document: %w", err)
		}
	default:
		return true, fmt.Errorf("unknown job type: %T", msg)
	}
	return true, nil
}

func getFileExt(title string) (document.FileExtension, error) {
	parts := strings.Split(title, ".")
	if len(parts) < 2 {
		return document.PNG, errors.New("file extension not found in title")
	}
	ext := "." + strings.ToLower(parts[len(parts)-1])
	if val, ok := document.FileExtensionMap[ext]; ok {
		return val, nil
	}
	return document.PNG, errors.New("unsupported file extension: " + ext)
}

func UnmarshalJob(objType string, payload []byte) (any, error) {
	switch objType {
	case "SiteJob":
		var job qaas.SiteJob
		if err := json.Unmarshal(payload, &job); err != nil {
			return nil, fmt.Errorf("failed to unmarshal SiteJob: %w", err)
		}
		return job, nil
	case "PageJob":
		var job qaas.PageJob
		if err := json.Unmarshal(payload, &job); err != nil {
			return nil, fmt.Errorf("failed to unmarshal PageJob: %w", err)
		}
		return job, nil
	case "EmbedJob":
		var job qaas.EmbedJob
		if err := json.Unmarshal(payload, &job); err != nil {
			return nil, fmt.Errorf("failed to unmarshal EmbedJob: %w", err)
		}
		return job, nil
	default:
		return nil, fmt.Errorf("unsupported objType: %s", objType)
	}
}
