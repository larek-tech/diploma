package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
)

const listScenarios = `
	select id, user_id, use_multiquery, n_queries, query_model_name, use_rerank, reranker_model_name, reranker_max_length, reranker_top_k, llm_model_name, temperature, top_k, top_p, system_prompt, top_n, threshold, search_by_query, created_at, updated_at
	from domain.scenario
		where user_id = $1
	order by created_at desc, updated_at desc
	offset $2
	limit $3;
`

// ListScenarios returns list of scenarios available for user.
func (r *Repo) ListScenarios(ctx context.Context, userID int64, offset, limit uint64) ([]model.ScenarioDao, error) {
	var scenarios []model.ScenarioDao
	if err := r.pg.QuerySlice(ctx, &scenarios, listScenarios, userID, offset, limit); err != nil {
		return scenarios, errs.WrapErr(err, "list scenarios")
	}
	return scenarios, nil
}
