package ollama

import (
	"context"

	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/samber/lo"
)

type Service struct {
	client *api.Client
}

const (
	EmbeddingSize   = 8192
	EmbeddingsModel = "bge-m3"
	LLMModel        = "llama3.1:latest"
	LLMContextSize  = 16000
)

var keepAlive = time.Hour * 24

func New(host string) (*Service, error) {
	ollamaURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(ollamaURL, http.DefaultClient)

	return &Service{client: client}, nil
}

func (s Service) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(inputTexts))
	for i, t := range inputTexts {
		req := &api.EmbeddingRequest{
			Model:     EmbeddingsModel,
			Prompt:    t,
			KeepAlive: &api.Duration{Duration: keepAlive},
			Options:   map[string]any{"num_ctx": EmbeddingSize},
		}
		res, err := s.client.Embeddings(ctx, req)
		if err != nil {
			return nil, err
		}
		embeddings[i] = lo.Map(res.Embedding, func(e float64, _ int) float32 { return float32(e) })
	}

	return embeddings, nil
}

func (s Service) Call(ctx context.Context, prompt string) (string, error) {
	stream := false
	req := &api.ChatRequest{
		Model:     LLMModel,
		Messages:  []api.Message{{Role: "user", Content: prompt}},
		Stream:    &stream,
		KeepAlive: &api.Duration{Duration: keepAlive},
		Options:   map[string]any{"num_ctx": LLMContextSize},
	}
	response := ""
	err := s.client.Chat(ctx, req, func(res api.ChatResponse) error {
		response += res.Message.Content
		return nil
	})
	if err != nil {
		return "", err
	}
	return "", nil
}
