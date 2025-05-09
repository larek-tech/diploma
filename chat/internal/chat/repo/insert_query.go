package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertQuery = `
	insert into chat.query(user_id, chat_id, content, domain_id, source_ids, scenario_id, metadata)
	values ($1, $2, $3, $4, $5, $6, $7)
	returning id;
`

// InsertQuery creates new query in chat.
func (r *Repo) InsertQuery(ctx context.Context, q model.QueryDao) (int64, error) {
	var queryID int64
	if err := r.pg.Query(
		ctx,
		&queryID,
		insertQuery,
		q.UserID,
		q.ChatID,
		q.Content,
		q.DomainID,
		q.SourceIDs,
		q.ScenarioID,
		q.Metadata,
	); err != nil {
		return 0, errs.WrapErr(err, "insert query")
	}
	return queryID, nil
}
