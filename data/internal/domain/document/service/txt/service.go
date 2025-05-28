package txt

import (
	"fmt"
	"io"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s Service) Parse(content io.ReadSeeker) (string, error) {
	// convert text to string
	rawBytes, err := io.ReadAll(content)
	if err != nil {
		return "", fmt.Errorf("error reading text: %w", err)
	}
	return string(rawBytes), nil
}
