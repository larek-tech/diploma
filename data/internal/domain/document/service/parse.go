package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/russross/blackfriday/v2"
	"github.com/unidoc/unidoc/pdf/contentstream"
	"github.com/unidoc/unidoc/pdf/model"
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

	return &document.Document{
		ID:        uuid.NewString(),
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func ParseMarkdown(content io.ReadSeeker) (string, error) {
	// Convert Markdown to HTML
	rawBytes, err := io.ReadAll(content)
	if err != nil {
		return "", fmt.Errorf("error reading markdown: %w", err)
	}

	html := blackfriday.Run(rawBytes)
	return string(html), nil
}

func ParsePDF(content io.ReadSeeker) (string, error) {
	reader, err := model.NewPdfReader(content)
	if err != nil {
		return "", fmt.Errorf("error creating PDF reader: %w", err)
	}

	var extractedText strings.Builder
	numPages, err := reader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("error getting page count: %w", err)
	}

	// Extract text from each page
	for i := 1; i <= numPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			return "", fmt.Errorf("error getting page %d: %w", i, err)
		}

		// Extract text from the page
		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return "", fmt.Errorf("error getting content streams: %w", err)
		}

		// Process content streams
		for _, cstream := range contentStreams {
			cstreamParser := contentstream.NewContentStreamParser(cstream)
			operations, err := cstreamParser.Parse()
			if err != nil {
				return "", fmt.Errorf("error parsing content stream: %w", err)
			}

			extractedText.Write(operations.Bytes())
		}
	}

	return extractedText.String(), nil
}

func ParseTXT(content io.ReadSeeker) (string, error) {
	rawBytes, err := io.ReadAll(content)
	if err != nil {
		return "", fmt.Errorf("error reading txt content: %w", err)
	}
	return string(rawBytes), nil
}
