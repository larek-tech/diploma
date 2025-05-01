package crawler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/samber/lo"
)

// ParsePage parses the page and returns a list of outgoing pages.
func (s Service) ParsePage(ctx context.Context, page *site.Page) ([]*site.Page, error) {
	err := validate(page)
	if err != nil {
		return nil, fmt.Errorf("failed to validate page: %w", err)
	}
	siteInfo, err := s.siteStore.GetByID(ctx, page.SiteID)
	if err != nil {
		return nil, err
	}
	sameDomain, err := isSameDomain(page.URL, siteInfo.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to check domain: %w", err)
	}
	if !sameDomain {
		return nil, fmt.Errorf("url %s is not in the same domain as site %s", page.URL, siteInfo.URL)
	}

	outgoingLinks, err := s.fetchContent(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}
	outgoingPages := make([]*site.Page, 0, len(outgoingLinks))
	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		for _, link := range outgoingLinks {
			outgoingPage, outGoingErr := s.pageStore.GetByURL(ctx, link)
			if outGoingErr != nil {
				slog.Error("failed to get outgoing page", "outGoingErr", outGoingErr)
				continue
			}
			if outgoingPage == nil {
				outgoingPage = &site.Page{
					ID:            uuid.NewString(),
					URL:           link,
					SiteID:        page.SiteID,
					OutgoingPages: make([]string, 0),
				}
				saveErr := s.pageStore.Save(ctx, outgoingPage)
				if saveErr != nil {
					return fmt.Errorf("failed to save outgoing page: %w", saveErr)

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
			return fmt.Errorf("failed to save page: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return outgoingPages, nil
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
