package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
)

const getDefaultScenario = `
	select id, title, user_id, use_multiquery, n_queries, query_model_name, use_rerank, reranker_model_name, reranker_max_length, reranker_top_k, llm_model_name, temperature, top_k, top_p, system_prompt, top_n, threshold, search_by_query, created_at, updated_at
	from domain.scenario
	where title = $1
		and user_id = $2;
`

// GetDefaultScenario returns default scenario for domain.
func (r *Repo) GetDefaultScenario(ctx context.Context, title string, userID int64) (model.ScenarioDao, error) {
	var scenario model.ScenarioDao
	if err := r.pg.Query(ctx, &scenario, getDefaultScenario, title, userID); err != nil {
		return scenario, errs.WrapErr(err, "get default scenario")
	}
	return scenario, nil
}
