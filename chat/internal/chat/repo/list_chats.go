package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const listChats = `
	select id, user_id, title, created_at, updated_at
	from chat.chat
	where user_id = $1
		and is_deleted = false
	offset $2
	limit $3;
`

// ListChats returns list of user active chats.
func (r *Repo) ListChats(ctx context.Context, offset, limit uint64, userID int64) ([]model.ChatDao, error) {
	var chats []model.ChatDao
	if err := r.pg.QuerySlice(ctx, &chats, listChats, userID, offset, limit); err != nil {
		return nil, errs.WrapErr(err, "list chats")
	}
	return chats, nil
}
