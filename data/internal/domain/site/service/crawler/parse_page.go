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
	err := validate(page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to validate page: %w", err)
	}
	if parsed, err := s.pageJobStore.IsAlreadyParsed(ctx, page.URL); err != nil {
		return nil, false, fmt.Errorf("failed to check if page is already parsed: %w", err)
	} else if parsed {
		return nil, false, errors.New("page already parsed")
	}

	//oldPage, err := s.pageStore.GetByURL(ctx, page.URL)
	//if err != nil {
	//	slog.Debug("page not found", "err", err)
	//}
	// TODO: create mechanism to check if page was parsed in the last ParsingDelta
	//if oldPage != nil {
	//	if time.Since(oldPage.UpdatedAt) < ParsingDelta {
	//		slog.Debug("page already parsed", "page", page.URL)
	//		return nil, true, nil
	//	}
	//	slog.Debug("page already parsed, but outdated", "page", page.URL)
	//	page.ID = oldPage.ID
	//	page.OutgoingPages = lo.Uniq(append(page.OutgoingPages, oldPage.OutgoingPages...))
	//	err = s.pageStore.Save(ctx, page)
	//	if err != nil {
	//		return nil, false, fmt.Errorf("failed to save page: %w", err)
	//	}
	//	return nil, true, nil
	//}

	siteInfo, err := s.siteStore.GetByID(ctx, page.SiteID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get site: %w", err)
	}
	sameDomain, err := isSameDomain(page.URL, siteInfo.URL)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check domain: %w", err)
	}
	if !sameDomain {
		return nil, false, fmt.Errorf("url %s is not in the same domain as site %s", page.URL, siteInfo.URL)
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
