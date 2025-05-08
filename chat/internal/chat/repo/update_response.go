package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const updateResponse = `
	update chat.response
	set content = $2,
	    status = $3
	where id = $1;
`

// UpdateResponse updates response in chat.
func (r *Repo) UpdateResponse(ctx context.Context, resp model.ResponseDao) error {
	if _, err := r.pg.Exec(ctx, updateResponse, resp.ID, resp.Content, resp.Status); err != nil {
		return errs.WrapErr(err, "update response")
	}
	return nil
}
