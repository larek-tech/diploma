package crawler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/samber/lo"
)

const ParsingDelta = time.Hour * 12

// ParsePage parses the page and returns a list of outgoing pages.
func (s Service) ParsePage(ctx context.Context, page *site.Page) ([]*site.Page, bool, error) {
	err := validate(page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to validate page: %w", err)
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

	outgoingLinks, err := s.fetchContent(ctx, page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch content: %w", err)
	}
	outgoingPages := make([]*site.Page, 0, len(outgoingLinks))
	for _, link := range outgoingLinks {
		outgoingPage, outGoingErr := s.pageStore.GetByURL(ctx, link)
		if outGoingErr != nil {
			slog.Error("failed to get outgoing page", "outGoingErr", outGoingErr)
			continue
		}
		if outgoingPage == nil {
			outgoingPage, err = site.NewPage(siteInfo.ID, link)
			if err != nil {
				slog.Warn("failed to create outgoing page", "err", err)
				continue
			}
		}
		outgoingPages = append(outgoingPages, outgoingPage)
	}
	outgoingIDs := lo.Map(outgoingPages, func(p *site.Page, _ int) string {
		return p.ID
	})
	page.OutgoingPages = lo.Uniq(append(page.OutgoingPages, outgoingIDs...))

	err = s.pageStore.Save(ctx, page)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save page: %w", err)
	}

	return outgoingPages, true, nil
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
