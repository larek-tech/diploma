package repo

import (
	"context"
	"github.com/yogenyslav/pkg/errs"
)

const findUserRoles = `
	select role_id
	from auth.user_role
	where user_id = $1;
`

// FindUserRoles returns a slices of role ids of a certain user.
func (r *AuthRepo) FindUserRoles(ctx context.Context, userID int64) ([]int64, error) {
	var roles []int64
	if err := r.pg.QuerySlice(ctx, &roles, findUserRoles, userID); err != nil {
		return nil, errs.WrapErr(err, "find user roles")
	}
	return roles, nil
}
