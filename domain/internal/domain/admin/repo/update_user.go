package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/domain/internal/domain/admin/model"
	"github.com/yogenyslav/pkg/errs"
)

const updateUser = `
	update auth.user
	set email = $2,
		hash_password = $3
	where id = $1
		and is_deleted = false;
`

// UpdateUser updates user data.
func (r *Repo) UpdateUser(ctx context.Context, u model.UserDao) error {
	rows, err := r.pg.Exec(
		ctx,
		updateUser,
		u.ID,
		u.Email,
		u.HashPassword,
	)

	if err != nil {
		return errs.WrapErr(err, "update user")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "update user not found")
	}

	return nil
}
