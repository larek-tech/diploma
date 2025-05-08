package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const softDeleteChat = `
	update chat.chat
	set is_deleted = true
	where id = $1
		and is_deleted = false;
`

// SoftDeleteChat soft delete chat.
func (r *Repo) SoftDeleteChat(ctx context.Context, chatID uuid.UUID) error {
	rows, err := r.pg.Exec(ctx, softDeleteChat, chatID)
	if err != nil {
		return errs.WrapErr(err, "soft delete chat")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "soft delete chat")
	}

	return nil
}

const deleteChat = `
	delete from chat.chat
	where id = $1
		and is_deleted = false
		and (
		    select count(id)
		    from chat.query
		    where chat_id = $1
		) = 0;
`

// DeleteChat delete chat.
func (r *Repo) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	_, err := r.pg.Exec(ctx, deleteChat, chatID)
	if err != nil {
		return errs.WrapErr(err, "delete chat")
	}
	return nil
}
