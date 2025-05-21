package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/document"
)

var ErrFileTypeNotSupported = errors.New("file extension not supported")

// parse parses the document based on its file extension and returns a Document object.
func (s Service) parse(_ context.Context, obj io.ReadSeeker, fileExt document.FileExtension) (*document.Document, error) {
	parser, found := s.parsers[fileExt]
	if !found {
		return nil, fmt.Errorf("failed to parse file unsupported filetype %v: %w", fileExt, ErrFileTypeNotSupported)
	}
	content, err := parser.Parse(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %w", err)
	}
	// TODO: add source metadata (link)
	return &document.Document{
		ID:        uuid.NewString(),
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
