package crawler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/site"
)

const ParsingDelta = time.Hour * 12

// ParsePage parses the page and returns a list of outgoing pages.
func (s Service) ParsePage(ctx context.Context, page *site.Page, parseSiteJobID string) ([]*site.Page, bool, error) {
	ctx, span := s.tracer.Start(ctx, "ParsePage")
	defer span.End()

	err := validate(page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to validate page: %w", err)
	}
	if parsed, err := s.pageJobStore.IsAlreadyParsed(ctx, page.URL); err != nil {
		return nil, false, fmt.Errorf("failed to check if page is already parsed: %w", err)
	} else if parsed {
		return nil, false, errors.New("page already parsed")
	}

	_, err = s.fetchContent(ctx, page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch content: %w", err)
	}
	if page.Raw == "" {
		return nil, false, fmt.Errorf("page raw content is empty")
	}

	err = s.pageStore.Save(ctx, page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save page: %w", err)
	}

	return nil, true, nil
}

func validate(page *site.Page) error {
	if page == nil {
		return fmt.Errorf("page is nil")
	}
	if page.URL == "" {
		return fmt.Errorf("page URL is empty")
	}
	if page.SiteID == "" {
		return fmt.Errorf("page SiteID is empty")
	}
	return nil
}
