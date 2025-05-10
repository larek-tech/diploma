package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const getResponseByQueryID = `
	select id, query_id, chat_id, content, status, created_at, updated_at
	from chat.response
	where query_id = $1;
`

// GetResponseByQueryID returns response by query id.
func (r *Repo) GetResponseByQueryID(ctx context.Context, queryID int64) (model.ResponseDao, error) {
	var resp model.ResponseDao
	if err := r.pg.Query(ctx, &resp, getResponseByQueryID, queryID); err != nil {
		return resp, errs.WrapErr(err, "get response by query id")
	}
	return resp, nil
}
