package repo

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const listChats = `
	select c.id, c.user_id, c.title, c.created_at, c.updated_at
	from chat.chat c
	join chat.query q
		on q.chat_id = c.id
	where c.user_id = $1
		and c.is_deleted = false
	order by c.updated_at desc
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
