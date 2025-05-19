package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
)

const updateScenario = `
	update domain.scenario
	set use_multiquery=$3,
	    n_queries=$4,
	    query_model_name=$5,
	    use_rerank=$6,
	    reranker_model_name=$7,
	    reranker_max_length=$8,
	    reranker_top_k=$9,
	    llm_model_name=$10,
	    temperature=$11,
	    top_k=$12,
	    top_p=$13,
	    system_prompt=$14,
	    top_n=$15,
	    threshold=$16,
	    search_by_query=$17,
		title=$18,
		context_size=$19
	where id = $1
		and user_id = $2;
`

// UpdateScenario updates data for scenario.
func (r *Repo) UpdateScenario(ctx context.Context, s model.ScenarioDao, userID int64) error {
	rows, err := r.pg.Exec(
		ctx,
		updateScenario,
		s.ID,
		userID,
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
		s.Title,
		s.ContextSize,
	)
	if err != nil {
		return errs.WrapErr(err, "update scenario")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "scenario not found")
	}

	return nil
}
