package repo

import (
	"context"

	"github.com/yogenyslav/pkg/errs"
)

const getPermittedUsers = `
	select user_id
	from domain.domain_permitted_users
	where domain_id = $1;
`

// GetPermittedUsers returns list of users that have access to the domain.
func (r *Repo) GetPermittedUsers(ctx context.Context, domainID int64) ([]int64, error) {
	var userIDs []int64
	if err := r.pg.QuerySlice(ctx, &userIDs, getPermittedUsers, domainID); err != nil {
		return nil, errs.WrapErr(err, "get permitted users")
	}
	return userIDs, nil
}
