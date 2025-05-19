package pdf

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
)

const pageSeparator = "\n\n\n\n\n"

type Service struct {
	ocr ocr
}

func New(ocr ocr) *Service {
	return &Service{
		ocr: ocr,
	}
}

func (s Service) Parse(reader io.ReadSeeker) (string, error) {
	pdf, err := os.CreateTemp("", "ocr-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(pdf.Name())
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read content: %w", err)
	}
	if _, err := pdf.Write(content); err != nil {
		log.Fatal(err)
	}

	doc, err := fitz.New(pdf.Name())
	if err != nil {
		return "", fmt.Errorf("failed to split pdf: %w", err)
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "fitz")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	text := make([]string, doc.NumPage())

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			return "", fmt.Errorf("failed to get pdf page: %w", err)

		}

		f, err := os.Create(filepath.Join(tmpDir, fmt.Sprintf("test%03d.jpg", n)))
		if err != nil {
			return "", fmt.Errorf("failed to save temp file for pdf page: %w", err)

		}
		defer f.Close()

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			return "", fmt.Errorf("failed to save pdf page jpeg: %w", err)

		}
		pageText, err := s.ocr.Process(f.Name())
		if err != nil {
			return "", fmt.Errorf("failed to get text: %w", err)

		}
		text[n] = pageText

	}

	return strings.Join(text, pageSeparator), nil
}
