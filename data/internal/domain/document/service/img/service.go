package img

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Service struct {
	ocr ocr
}

func New(ocr ocr) *Service {
	return &Service{
		ocr: ocr,
	}
}

func (s *Service) Parse(reader io.ReadSeeker) (string, error) {
	f, err := os.CreateTemp("", "ocr-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(f.Name())
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read content: %w", err)
	}
	if _, err := f.Write(content); err != nil {
		log.Fatal(err)
	}
	text, err := s.ocr.Process(f.Name())
	if err != nil {
		return "", fmt.Errorf("failed to get text: %w", err)
	}
	return text, nil
}
