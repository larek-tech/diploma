package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/pkg/errs"
)

const getChatUserID = `
	select user_id
	from chat.chat
	where id = $1
		and is_deleted = false;
`

// GetChatUserID returns the userID of chat.
func (r *Repo) GetChatUserID(ctx context.Context, chatID uuid.UUID) (int64, error) {
	var userID int64
	if err := r.pg.Query(ctx, &userID, getChatUserID, chatID); err != nil {
		return 0, errs.WrapErr(err, "get chat user id")
	}
	return userID, nil
}
