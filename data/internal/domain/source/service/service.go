package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/domain/sitemap"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/samber/lo"
)

type Service struct {
	sitemapParser sitemapParser
	sourceStorage sourceStorage
	pub           publisher
	trManager     transactionalManager
}

func New(sourceStorage sourceStorage, sitemapParser sitemapParser, pub publisher, trManager transactionalManager) *Service {
	return &Service{
		sitemapParser: sitemapParser,
		sourceStorage: sourceStorage,
		pub:           pub,
		trManager:     trManager,
	}
}

func (s Service) CreateSource(ctx context.Context, message source.DataMessage) (*source.Source, error) {
	src := &source.Source{
		ID:    uuid.NewString(),
		Title: message.Title,
		Type:  message.Type,
	}
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		err := s.sourceStorage.Save(ctx, src)
		if err != nil {
			return err
		}
		switch src.Type {
		case source.Web:
			var webSource *site.Site
			webSource, err = s.createSite(src, message)
			if err != nil {
				return fmt.Errorf("failed to create site: %w", err)
			}
			publishOptions := []qaas.PublishOption{
				qaas.WithQueue(qaas.ParseSiteQueue),
			}

			_, err = s.pub.Publish(ctx, []any{qaas.SiteJob{
				Payload: webSource,
				Delay:   0,
				Metadata: map[string]any{
					"siteJobID": uuid.NewString(),
				},
			}}, publishOptions...)
			if err != nil {
				return fmt.Errorf("failed to publish site job: %w", err)
			}
		default:
			return fmt.Errorf("unsupported source type: %v", src.Type)
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return src, nil
}

func (s Service) createSite(src *source.Source, message source.DataMessage) (*site.Site, error) {
	var decodedURL []byte
	_, err := base64.StdEncoding.Decode(decodedURL, message.Content)
	if err != nil {
		return nil, err
	}

	siteURL, err := url.Parse(string(decodedURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url for web source: %w", err)
	}
	site, err := site.NewSite(src.ID, siteURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	availableURLs, err := s.sitemapParser.GetAndParseSitemap(siteURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap: %w", err)
	}
	site.AvailablePages = lo.Map(availableURLs, func(v sitemap.URLResult, _ int) string { return v.URL })

	return site, nil
}
