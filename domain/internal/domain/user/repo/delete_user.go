package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const deleteUser = `
	update auth.user
	set is_deleted = true
	where id = $1
		and is_deleted = false;
`

// DeleteUser soft delete user.
func (r *Repo) DeleteUser(ctx context.Context, id int64) error {
	rows, err := r.pg.Exec(ctx, deleteUser, id)

	if err != nil {
		return errs.WrapErr(err, "delete user")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "delete user not found")
	}

	return nil
}
