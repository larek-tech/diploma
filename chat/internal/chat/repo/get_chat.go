package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/errs"
)

const getChat = `
	select 
		(q.id, q.user_id, q.chat_id, q.content, q.domain_id, q.source_ids, q.scenario_id, q.metadata, q.created_at) as query,
		(r.id, r.query_id, r.chat_id, r.content, r.status, r.metadata, r.created_at, r.updated_at) as response,
		c.id, c.user_id, c.title, c.created_at, c.updated_at
	from chat.query q
	join
	    chat.response r
		on q.id = r.query_id
	join 
	    chat.chat c
		on q.chat_id = c.id
	where
	    q.chat_id = $1 
	  	and c.is_deleted = false
	order by q.id;
`

// GetChat returns chat by id.
func (r *Repo) GetChat(ctx context.Context, chatID uuid.UUID) (model.ChatDao, error) {
	var chat model.ChatDao
	if err := r.pg.Query(ctx, &chat, getChat, chatID); err != nil {
		return chat, errs.WrapErr(err, "get chat")
	}
	return chat, nil
}
