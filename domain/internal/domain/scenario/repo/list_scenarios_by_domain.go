package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
)

const listScenariosByDomainQuery = `
	select id, user_id, domain_id, context_size use_multiquery, n_queries, query_model_name, use_rerank, 
	       reranker_model_name, reranker_max_length, reranker_top_k, llm_model_name, temperature, top_k, top_p, 
	       system_prompt, top_n, threshold, search_by_query, created_at, updated_at
	from domain.scenario
		where domain_id = $3
	order by created_at desc, updated_at desc
	offset $1
	limit $2;
`

// ListScenariosByDomain returns list of scenarios available for user by domain.
func (r *Repo) ListScenariosByDomain(
	ctx context.Context,
	domainID int64,
	offset, limit uint64,
) ([]model.ScenarioDao, error) {
	var scenarios []model.ScenarioDao
	if err := r.pg.QuerySlice(ctx, &scenarios, listScenariosByDomainQuery, offset, limit, domainID); err != nil {
		return scenarios, errs.WrapErr(err, "list scenarios by domain")
	}
	return scenarios, nil
}
