package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/domain/document"
)

const (
	// ChunkSize is the size of each chunk in characters.
	ChunkSize = 1024
	// ChunkOverlap is the number of overlapping characters between chunks.
	ChunkOverlap  = 100
	EmbeddingSize = 1024
)

// embed embeds the document content into chunks and returns them.
func (s Service) embed(ctx context.Context, doc *document.Document) ([]*document.Chunk, error) {
	err := validateDocument(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to validate document: %w", err)
	}
	content := document.CleanUTF8(doc.Content)
	rawChunks := characterTextSplitter(content, ChunkSize, ChunkOverlap)
	if len(rawChunks) == 0 {
		return nil, nil
	}
	embeddings, err := s.embedder.CreateEmbedding(ctx, rawChunks)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}
	// TODO: add generation of questions using llm model
	metadata, err := json.Marshal(doc)
	if err != nil {
		metadata = []byte{}
	}
	chunks := make([]*document.Chunk, 0, len(rawChunks))
	for i, rawChunk := range rawChunks {
		if len(embeddings[i]) > EmbeddingSize {
			embeddings[i] = embeddings[i][:EmbeddingSize]
		} else if len(embeddings[i]) < EmbeddingSize {
			return nil, fmt.Errorf("embedding at index %d is smaller than %d", i, EmbeddingSize)
		}

		chunk := &document.Chunk{
			ID:         uuid.NewString(),
			DocumentID: doc.ID,
			SourceID:   doc.SourceID,
			Content:    rawChunk,
			Index:      i,
			Embeddings: embeddings[i],
			Metadata:   metadata,
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func validateDocument(doc *document.Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	if doc.Content == "" {
		return fmt.Errorf("document content is empty")
	}
	if doc.ID == "" {
		return fmt.Errorf("document ID is empty")
	}
	return nil
}

func characterTextSplitter(text string, chunkSize, overlap int) []string {
	if chunkSize <= 0 || len(text) == 0 {
		return []string{}
	}

	var chunks []string
	start := 0

	for start < len(text) {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}

		chunk := text[start:end]
		chunks = append(chunks, chunk)

		start += chunkSize - overlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}
