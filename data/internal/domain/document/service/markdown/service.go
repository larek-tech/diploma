package markdown

import (
	"fmt"
	"io"

	"github.com/russross/blackfriday/v2"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s Service) Parse(content io.ReadSeeker) (string, error) {
	// Convert Markdown to HTML
	rawBytes, err := io.ReadAll(content)
	if err != nil {
		return "", fmt.Errorf("error reading markdown: %w", err)
	}

	html := blackfriday.Run(rawBytes)
	return string(html), nil
}
