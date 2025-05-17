package html

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockFile   = "mocks/html/index.html"
	resultFile = "mocks/html/index.txt"
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
