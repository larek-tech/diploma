package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/file"
	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/pkg/metric"
)

type Service struct {
	sitemapParser sitemapParser
	sourceStorage sourceStorage
	fileStorage   fileStorage
	pub           publisher
	trManager     transactionalManager
}

func New(sourceStorage sourceStorage, fileStorage fileStorage, sitemapParser sitemapParser, pub publisher, trManager transactionalManager) *Service {
	return &Service{
		sitemapParser: sitemapParser,
		sourceStorage: sourceStorage,
		fileStorage:   fileStorage,
		pub:           pub,
		trManager:     trManager,
	}
}

func (s Service) CreateSource(ctx context.Context, msg source.DataMessage) (*source.Source, error) {
	src := &source.Source{
		ID:    uuid.NewString(),
		Title: msg.Title,
		Type:  msg.Type,
	}
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		err := s.sourceStorage.Save(ctx, src)
		if err != nil {
			return err
		}
		switch src.Type {
		case source.Web:
			var webSource *site.Site
			webSource, err = s.createSite(src, msg)
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
					"siteJobID":   uuid.NewString(),
					"externalKey": string(msg.ExternalKey),
				},
			}}, publishOptions...)
			if err != nil {
				return fmt.Errorf("failed to publish site job: %w", err)
			}
		case source.S3WithCredentials:
			slog.Info("creating s3 source", "source", src)
		case source.SingleFile:
			var file *file.File
			file, err := s.createFile(src, msg)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			err = s.fileStorage.Save(ctx, file)
			if err != nil {
				return fmt.Errorf("failed to save file: %w", err)
			}
			publishOptions := []qaas.PublishOption{
				qaas.WithQueue(qaas.ParseFileQueue),
				qaas.WithSourceQueue(qaas.ParseFileQueue),
			}
			_, err = s.pub.Publish(ctx, []any{qaas.FileJob{
				Payload: file,
				Delay:   0,
				Metadata: map[string]any{
					"externalKey": string(msg.ExternalKey),
				},
			}}, publishOptions...)
			if err != nil {
				return fmt.Errorf("failed to publish file job: %w", err)
			}
		case source.ArchivedFiles:
			files, err := s.createArchive(src, msg)
			if err != nil {
				return fmt.Errorf("failed to create archive: %w", err)
			}
			var jobs []any
			for _, f := range files {
				err = s.fileStorage.Save(ctx, f)
				if err != nil {
					return fmt.Errorf("failed to save file: %w", err)
				}
				jobs = append(jobs, qaas.FileJob{
					Payload: f,
					Delay:   0,
					Metadata: map[string]any{
						"externalKey": string(msg.ExternalKey),
					},
				})

			}
			publishOptions := []qaas.PublishOption{
				qaas.WithQueue(qaas.ParseFileQueue),
				qaas.WithSourceQueue(qaas.ParseFileQueue),
			}
			_, err = s.pub.Publish(ctx, jobs, publishOptions...)
			if err != nil {
				return fmt.Errorf("failed to publish file job: %w", err)
			}
		default:
			return fmt.Errorf("unsupported source type: %v", src.Type)
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	metric.IncrementSourcesCreated((string(src.Type)), src.ID, err)
	return src, nil
}
