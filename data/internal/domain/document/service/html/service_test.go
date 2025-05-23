package html

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockFile   = "mocks/html/index.html"
	resultFile = "mocks/html/index.txt"
)

//go:embed test/*.html
var htmlFiles embed.FS

const (
	smallFile  = "500kb.html"
	mediumFile = "1mb.html"
	largeFile  = "2mb.html"
)

func loadFile() string {
	file, err := os.ReadFile(mockFile)
	if err != nil {
		panic(err)
	}
	return string(file)
}

func saveResult(content string) {
	err := os.WriteFile(resultFile, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func TestParsing(t *testing.T) {
	content := loadFile()
	reader := strings.NewReader(content)
	s := New()
	processed, err := s.Parse(reader)
	assert.Greater(t, len(content), len(processed))
	assert.Greater(t, len(processed), 0)
	assert.NoError(t, err)
	saveResult(processed)
}

func BenchmarkService_Parse(b *testing.B) {
	files := map[string]string{
		"small":  smallFile,
		"medium": mediumFile,
		"large":  largeFile,
	}
	for name, file := range files {
		b.Run(name, func(b *testing.B) {
			content, err := htmlFiles.ReadFile("test/" + file)
			if err != nil {
				b.Fatalf("failed to read file: %v", err)
			}
			s := New()
			for i := 0; i < b.N; i++ {
				reader := strings.NewReader(string(content))
				_, err := s.Parse(reader)
				if err != nil {
					b.Fatalf("parse error: %v", err)
				}
			}
		})
	}
}

func BenchmarkService_STDParse(b *testing.B) {
	files := map[string]string{
		"small":  smallFile,
		"medium": mediumFile,
		"large":  largeFile,
	}
	for name, file := range files {
		b.Run(name, func(b *testing.B) {
			content, err := htmlFiles.ReadFile("test/" + file)
			if err != nil {
				b.Fatalf("failed to read file: %v", err)
			}
			s := New()
			for i := 0; i < b.N; i++ {
				reader := strings.NewReader(string(content))
				_, err := s.STDParse(reader)
				if err != nil {
					b.Fatalf("parse error: %v", err)
				}
			}
		})
	}
}
