package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
)

type Service struct {
	sourceStorage sourceStorage
	pub           publisher
	trManager     transactionalManager
}

func New(sourceStorage sourceStorage, pub publisher, trManager transactionalManager) *Service {
	return &Service{
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
			webSource, err = createSite(src, message)
			if err != nil {
				return fmt.Errorf("failed to create site: %w", err)
			}
			publishOptions := []qaas.PublishOption{
				qaas.WithQueue(qaas.ParseSiteQueue),
			}

			_, err = s.pub.Publish(ctx, []any{qaas.SiteJob{
				Payload: webSource,
				Delay:   0,
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

func createSite(src *source.Source, message source.DataMessage) (*site.Site, error) {
	siteURL, err := url.Parse(string(message.Content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url for web source: %w", err)
	}

	return &site.Site{
		ID:             uuid.NewString(),
		SourceID:       src.ID,
		URL:            siteURL.String(),
		AvailablePages: nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
