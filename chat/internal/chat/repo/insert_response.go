package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertResponse = `
	insert into chat.response(query_id, chat_id, content, status)
	values ($1, $2, $3, $4)
	returning id;
`

// InsertResponse creates new response in chat.
func (r *Repo) InsertResponse(ctx context.Context, resp model.ResponseDao) (int64, error) {
	var responseID int64
	if err := r.pg.Query(
		ctx,
		&responseID,
		insertResponse,
		resp.QueryID,
		resp.ChatID,
		resp.Content,
		resp.Status,
	); err != nil {
		return 0, errs.WrapErr(err, "insert response")
	}
	return responseID, nil
}
