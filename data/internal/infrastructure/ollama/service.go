package ollama

import (
	"context"
	"encoding/json"

	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/samber/lo"
)

type Service struct {
	client *api.Client
	cfg    *Config
}

type Config struct {
	EmbeddingSize   int
	EmbeddingsModel string
	LLMModel        string
	LLMContextSize  int
}

const (
	EmbeddingSize   = 8192
	EmbeddingsModel = "bge-m3"
	LLMModel        = "llama3.1:latest"
	LLMContextSize  = 32000
)

func NewDefaultConfig() *Config {
	return &Config{
		EmbeddingSize:   EmbeddingSize,
		EmbeddingsModel: EmbeddingsModel,
		LLMModel:        LLMModel,
		LLMContextSize:  LLMContextSize,
	}
}

var keepAlive = time.Hour * 24

func New(host string, cfg ...*Config) (*Service, error) {
	if len(cfg) == 0 {
		cfg = append(cfg, NewDefaultConfig())
	}

	ollamaURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(ollamaURL, http.DefaultClient)

	return &Service{client: client, cfg: cfg[0]}, nil
}

func (s Service) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(inputTexts))
	for i, t := range inputTexts {
		req := &api.EmbeddingRequest{
			Model:     s.cfg.EmbeddingsModel,
			Prompt:    t,
			KeepAlive: &api.Duration{Duration: keepAlive},
			Options:   map[string]any{"num_ctx": s.cfg.EmbeddingSize},
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
	req := &api.GenerateRequest{
		Model:     s.cfg.LLMModel,
		Stream:    &stream,
		KeepAlive: &api.Duration{Duration: keepAlive},
		Options:   map[string]any{"num_ctx": s.cfg.LLMContextSize},
		Prompt:    prompt,
		Context:   []int{},
		Raw:       false,
		Format:    json.RawMessage{},
		Images:    []api.ImageData{},
	}
	response := ""
	err := s.client.Generate(ctx, req, func(res api.GenerateResponse) error {
		response += res.Response
		return nil
	})
	if err != nil {
		return "", err
	}
	return response, nil
}
