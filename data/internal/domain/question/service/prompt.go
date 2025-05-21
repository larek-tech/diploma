package service

import (
	"embed"
	"io"
	"log"
)

//go:embed prompt
var prompt embed.FS

// GetSystemPrompt returns the system prompt from the embedded prompt.txt file
func GetSystemPrompt() string {
	// Open the embedded file
	file, err := prompt.Open("prompt/prompt.txt")
	if err != nil {
		log.Fatalf("Failed to open embedded prompt file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read embedded prompt file: %v", err)
	}

	return string(content)
}
