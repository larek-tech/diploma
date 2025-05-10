package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const getResponseByID = `
	select id, query_id, chat_id, content, status, created_at, updated_at
	from chat.response
	where id = $1;
`

// GetResponseByID returns response by id.
func (r *Repo) GetResponseByID(ctx context.Context, respID int64) (model.ResponseDao, error) {
	var resp model.ResponseDao
	if err := r.pg.Query(ctx, &resp, getResponseByID, respID); err != nil {
		return resp, errs.WrapErr(err, "get response by id")
	}
	return resp, nil
}
