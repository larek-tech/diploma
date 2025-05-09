package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertChat = `
	insert into chat.chat(user_id, title)
	values ($1, $2)
	returning id;
`

// InsertChat create new chat.
func (r *Repo) InsertChat(ctx context.Context, chat model.ChatDao) (uuid.UUID, error) {
	var chatID uuid.UUID
	if err := r.pg.Query(ctx, &chatID, insertChat, chat.UserID, chat.Title); err != nil {
		return chatID, errs.WrapErr(err, "insert chat")
	}
	return chatID, nil
}
