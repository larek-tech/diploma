package service

import (
	"context"
	"fmt"
	"io"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/site"
	"github.com/samber/lo"
)

// TODO: remove fileExt from Process func
func (s Service) Process(ctx context.Context, obj io.ReadSeeker, fileExt document.FileExtension, sourceObj any, sourceID string) error {
	doc, err := s.parse(ctx, obj, fileExt)
	if err != nil {
		return fmt.Errorf("failed to parse document: %w", err)
	}
	doc.SourceID = sourceID
	setDocData(doc, sourceObj)

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
		if txErr != nil {
			return fmt.Errorf("failed to update chunks: %w", txErr)
		}
		//questions, txErr := s.generateQuestions(ctx, chunks)
		//if txErr != nil {
		//	return fmt.Errorf("failed to generate questions: %w", txErr)
		//}
		//txErr = s.questionStorage.Save(ctx, questions)
		//if txErr != nil {
		//	return fmt.Errorf("failed to save questions: %w", txErr)
		//}
		doc.Chunks = lo.Map(chunks, func(chunk *document.Chunk, _ int) string {
			return chunk.ID
		})
		txErr = s.documentStorage.Save(ctx, doc)
		if txErr != nil {
			return fmt.Errorf("failed to save document: %w", txErr)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to process document: %w", err)
	}
	return nil
}

func setDocData(doc *document.Document, source any) {
	switch v := source.(type) {
	case *site.Page:
		doc.ObjectID = v.ID
		doc.ObjectType = document.TypePage
	default:
		doc.ObjectID = ""
		doc.ObjectType = document.TypeFile
	}
}
