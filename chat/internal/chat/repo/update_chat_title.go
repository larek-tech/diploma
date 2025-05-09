package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const updateChatTitle = `
	update chat.chat
	set title = $1
	where id = $2
		and is_deleted = false;
`

// UpdateChatTitle updates chat title.
func (r *Repo) UpdateChatTitle(ctx context.Context, title string, chatID uuid.UUID) error {
	rows, err := r.pg.Exec(ctx, updateChatTitle, title, chatID)
	if err != nil {
		return errs.WrapErr(err, "update chat title")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "update chat title")
	}

	return nil
}
