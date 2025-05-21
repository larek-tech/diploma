package service

import (
	"context"
	"fmt"
	"io"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/file"
	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/larek-tech/diploma/data/pkg/metric"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO: remove fileExt from Process func
func (s Service) Process(ctx context.Context, obj io.ReadSeeker, fileExt document.FileExtension, sourceObj any, sourceID string, metadata map[string]any) error {
	ctx, span := s.tracer.Start(ctx, "embeddingService.Process", trace.WithAttributes(
		attribute.String("sourceID", sourceID),
		attribute.String("fileExt", string(fileExt)),
		attribute.String("metadata", fmt.Sprintf("%v", metadata)),
	))
	defer span.End()
	objectID, docType := getObjectData(sourceObj)
	doc, err := s.parse(ctx, obj, fileExt)
	metric.IncrementDocumentsParsed(objectID, string(fileExt), sourceID, err)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to parse document: %w", err)
	}
	doc.ObjectID = objectID
	doc.ObjectType = docType
	doc.SourceID = sourceID
	doc.Metadata = metadata

	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		txErr := s.documentStorage.Save(ctx, doc)
		if err != nil {
			return fmt.Errorf("failed to save document: %w", txErr)
		}
		chunks, txErr := s.embed(ctx, doc)
		if txErr != nil {
			return fmt.Errorf("failed to embed document: %w", txErr)
		}
		txErr = s.chunkStorage.Update(ctx, doc.ID, chunks)
		metric.IncrementChunksCreated(doc.ID, doc.SourceID, string(doc.ObjectType), txErr, len(chunks))
		if txErr != nil {
			return fmt.Errorf("failed to update chunks: %w", txErr)
		}

		// Add check for nil questionService
		questions, txErr := s.questionService.GenerateQuestions(ctx, chunks)
		if txErr != nil {
			return fmt.Errorf("failed to generate questions: %w", txErr)
		}
		if questions != nil && len(questions) > 0 {
			txErr = s.questionStorage.Save(ctx, questions)
			if txErr != nil {
				return fmt.Errorf("failed to save questions: %w", txErr)
			}
		}

		doc.Chunks = lo.Map(chunks, func(chunk *document.Chunk, _ int) string {
			return chunk.ID
		})
		txErr = s.documentStorage.Save(ctx, doc)
		if txErr != nil {
			return fmt.Errorf("failed to save document: %w", txErr)
		}
		return nil
	})
	metric.IncrementDocumentsProcessed(string(doc.ObjectType), doc.SourceID, err)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to process document: %w", err)
	}
	return nil
}

func getObjectData(source any) (string, document.Type) {
	switch v := source.(type) {
	case *site.Page:
		return v.ID, document.TypePage
	case *file.File:
		return v.ID, document.TypeFile
	default:
		return "", document.TypeFile
	}
}
