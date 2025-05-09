package storage

import (
	"context"

	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/larek-tech/diploma/data/internal/domain/document"
)

type JSONStorage struct {
	storagePath string
	mu          sync.RWMutex
}

func NewJSONStorage(storagePath string) (*JSONStorage, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &JSONStorage{
		storagePath: storagePath,
	}, nil
}

func (s *JSONStorage) GetByID(ctx context.Context, id string) (*document.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.storagePath, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Document not found
		}
		return nil, fmt.Errorf("failed to read document file: %w", err)
	}

	var doc document.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return &doc, nil
}

func (s *JSONStorage) Save(ctx context.Context, doc *document.Document) error {
	slog.Info("json: saving document")
	if doc == nil {
		return fmt.Errorf("cannot save nil document")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	filePath := filepath.Join(s.storagePath, doc.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write document file: %w", err)
	}

	return nil
}

type MockStorage struct {
	docs map[string]*document.Document
	mu   sync.RWMutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		docs: make(map[string]*document.Document),
	}
}

func (m *MockStorage) GetByID(ctx context.Context, id string) (*document.Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	doc, exists := m.docs[id]
	if !exists {
		return nil, nil
	}
	return doc, nil
}

func (m *MockStorage) Save(ctx context.Context, doc *document.Document) error {
	if doc == nil {
		return fmt.Errorf("cannot save nil document")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.docs[doc.ID] = doc
	return nil
}
