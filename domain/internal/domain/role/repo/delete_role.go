package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const deleteRole = `
	update auth.role
	set is_deleted = true
	where id = $1
		and is_deleted = false;
`

// DeleteRole soft delete role.
func (r *Repo) DeleteRole(ctx context.Context, id int64) error {
	rows, err := r.pg.Exec(ctx, deleteRole, id)

	if err != nil {
		return errs.WrapErr(err, "delete role")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "delete role not found")
	}

	return nil
}
