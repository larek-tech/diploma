package service

import (
	"fmt"
	"net/url"

	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/domain/sitemap"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/samber/lo"
)

func (s Service) createSite(src *source.Source, msg source.DataMessage) (*site.Site, error) {
	rawUrl := string(msg.Content)
	if rawUrl == "" {
		return nil, fmt.Errorf("url is empty")
	}

	siteURL, err := url.Parse(string(msg.Content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url for web source: %w", err)
	}
	site, err := site.NewSite(src.ID, siteURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}
	availableURLs, err := s.sitemapParser.GetAndParseSitemap(*siteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap: %w", err)
	}
	site.AvailablePages = lo.Map(availableURLs, func(v sitemap.URLResult, _ int) string { return v.URL })

	return site, nil
}
