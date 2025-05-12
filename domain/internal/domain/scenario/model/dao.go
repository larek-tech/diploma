package model

import (
	"time"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ScenarioDao is a model for scenario on data layer.
type ScenarioDao struct {
	ID                int64     `db:"id"`
	Title             string    `db:"title"`
	UserID            int64     `db:"user_id"`
	UseMultiquery     bool      `db:"use_multiquery"`
	NQueries          int64     `db:"n_queries"`
	QueryModelName    string    `db:"query_model_name"`
	UseRerank         bool      `db:"use_rerank"`
	RerankerModelName string    `db:"reranker_model_name"`
	RerankerMaxLength int64     `db:"reranker_max_length"`
	RerankerTopK      int64     `db:"reranker_top_k"`
	LlmModelName      string    `db:"llm_model_name"`
	Temperature       float32   `db:"temperature"`
	TopK              int64     `db:"top_k"`
	TopP              float32   `db:"top_p"`
	SystemPrompt      string    `db:"system_prompt"`
	TopN              int64     `db:"top_n"`
	Threshold         float32   `db:"threshold"`
	SearchByQuery     bool      `db:"search_by_query"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// ToProto converts dao model into protobuf format.
func (s *ScenarioDao) ToProto() *pb.Scenario {
	return &pb.Scenario{
		Id:    s.ID,
		Title: s.Title,
		MultiQuery: &pb.MultiQuery{
			UseMultiquery:  s.UseMultiquery,
			NQueries:       s.NQueries,
			QueryModelName: &s.QueryModelName,
		},
		Reranker: &pb.Reranker{
			UseRerank:         s.UseRerank,
			RerankerModel:     s.RerankerModelName,
			RerankerMaxLength: s.RerankerMaxLength,
			TopK:              s.RerankerTopK,
		},
		VectorSearch: &pb.VectorSearch{
			TopN:          s.TopN,
			Threshold:     s.Threshold,
			SearchByQuery: s.SearchByQuery,
		},
		Model: &pb.LlmModel{
			ModelName:    s.LlmModelName,
			Temperature:  s.Temperature,
			TopK:         s.TopK,
			TopP:         s.TopP,
			SystemPrompt: s.SystemPrompt,
		},
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
}
