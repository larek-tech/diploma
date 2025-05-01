package service

import (
	"context"
	"fmt"
	"io"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/samber/lo"
)

func (s Service) Process(ctx context.Context, obj io.ReadSeeker, fileExt document.FileExtension) error {
	doc, err := s.parse(ctx, obj, fileExt)
	if err != nil {
		return fmt.Errorf("failed to parse document: %w", err)
	}
	err = s.trManager.Do(ctx, func(ctx context.Context) error {
		txErr := s.documentStorage.Save(ctx, doc)
		if err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}
		chunks, txErr := s.embed(ctx, doc)
		if txErr != nil {
			return fmt.Errorf("failed to embed document: %w", err)
		}
		txErr = s.chunkStorage.Update(ctx, doc.ID, chunks)
		if txErr != nil {
			return fmt.Errorf("failed to update chunks: %w", txErr)
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
	if err != nil {
		return fmt.Errorf("failed to process document: %w", err)
	}
	return nil
}
