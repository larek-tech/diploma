package repo

import (
	"context"

	"github.com/yogenyslav/pkg/errs"
)

const removeRole = `
	delete from auth.user_roles
	where user_id = $1
		and role_id = $2;
`

// RemoveRole removes role to user.
func (r *Repo) RemoveRole(ctx context.Context, userID, roleID int64) error {
	if _, err := r.pg.Exec(ctx, removeRole, userID, roleID); err != nil {
		return errs.WrapErr(err, "remove role")
	}
	return nil
}
