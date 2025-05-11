package repo

import (
	"context"

	"github.com/yogenyslav/pkg/errs"
)

const setRole = `
	insert auth.user_role (user_id, role_id)
	value ($1, $2);
`

// SetRole adds role to user.
func (r *Repo) SetRole(ctx context.Context, userID, roleID int64) error {
	if _, err := r.pg.Exec(ctx, setRole, userID, roleID); err != nil {
		return errs.WrapErr(err, "set role")
	}
	return nil
}
