package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertScenario = `
	insert into domain.scenario(title, user_id, domain_id, context_size, use_multiquery, n_queries, query_model_name, use_rerank, 
	                            reranker_model_name, reranker_max_length, reranker_top_k, llm_model_name, temperature, 
	                            top_k, top_p, system_prompt, top_n, threshold, search_by_query)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	returning id;
`

// InsertScenario create new scenario record.
func (r *Repo) InsertScenario(ctx context.Context, s model.ScenarioDao) (int64, error) {
	var scenarioID int64
	err := r.pg.Query(
		ctx,
		&scenarioID,
		insertScenario,
		s.Title,
		s.UserID,
		s.DomainID,
		s.ContextSize,
		s.UseMultiquery,
		s.NQueries,
		s.QueryModelName,
		s.UseRerank,
		s.RerankerModelName,
		s.RerankerMaxLength,
		s.RerankerTopK,
		s.LlmModelName,
		s.Temperature,
		s.TopK,
		s.TopP,
		s.SystemPrompt,
		s.TopN,
		s.Threshold,
		s.SearchByQuery,
	)
	if err != nil {
		return 0, errs.WrapErr(err, "insert scenario")
	}
	return scenarioID, nil
}
