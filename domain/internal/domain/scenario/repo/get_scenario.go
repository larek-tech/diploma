package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
)

const getScenarioByID = `
	select id, title, user_id, use_multiquery, n_queries, query_model_name, use_rerank, reranker_model_name, reranker_max_length, reranker_top_k, llm_model_name, temperature, top_k, top_p, system_prompt, top_n, threshold, search_by_query, created_at, updated_at
	from domain.scenario
	where id = $1
		and user_id = $2;
`

// GetScenarioByID returns scenario by ID.
func (r *Repo) GetScenarioByID(ctx context.Context, id, userID int64) (model.ScenarioDao, error) {
	var scenario model.ScenarioDao
	if err := r.pg.Query(ctx, &scenario, getScenarioByID, id, userID); err != nil {
		return scenario, errs.WrapErr(err, "get scenario by id")
	}
	return scenario, nil
}
