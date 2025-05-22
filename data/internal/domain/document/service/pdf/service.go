package pdf

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const pageSeparator = "\n\n\n\n\n"
const minImageWidth = 10
const minImageHeight = 10

type Service struct {
	ocr ocr
}

func New(ocr ocr) *Service {
	return &Service{
		ocr: ocr,
	}
}

func (s Service) Parse(reader io.ReadSeeker) (string, error) {
	pdf, err := os.CreateTemp("", "ocr-*.pdf")
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
	pdf.Close()

	tmpDir, err := os.MkdirTemp(os.TempDir(), "pdf-images-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPrefix := filepath.Join(tmpDir, "page")
	cmd := exec.Command(
		"pdftoppm",
		"-jpeg",
		"-r", "300",
		"-cropbox",
		pdf.Name(),
		outputPrefix,
	)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to convert PDF to images: %w", err)
	}

	generatedImages, err := filepath.Glob(filepath.Join(tmpDir, "page-*.jpg"))
	if err != nil {
		return "", fmt.Errorf("failed to list image files: %w", err)
	}

	var validTexts []string
	for _, file := range generatedImages {
		imgFile, err := os.Open(file)
		if err != nil {
			return "", fmt.Errorf("failed to open image file: %w", err)
		}

		img, _, err := image.DecodeConfig(imgFile)
		imgFile.Close()
		if err != nil {
			return "", fmt.Errorf("failed to decode image: %w", err)
		}

		if !valid(img) {
			log.Printf("Skipping small image %s (%dx%d)", file, img.Width, img.Height)
			continue
		}

		pageText, err := s.ocr.Process(file)
		if err != nil {
			return "", fmt.Errorf("failed to get text: %w", err)
		}
		validTexts = append(validTexts, pageText)
	}

	return strings.Join(validTexts, pageSeparator), nil
}

func valid(img image.Config) bool {
	return !(img.Width < minImageWidth || img.Height < minImageHeight)
}
