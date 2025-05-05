package repo

import (
	"context"

	"github.com/yogenyslav/pkg/errs"
)

const getPermittedUsers = `
	select user_id
	from domain.source_permitted_users
	where internal_source_id = $1;
`

// GetPermittedUsers returns list of users that have direct access to the source.
func (r *Repo) GetPermittedUsers(ctx context.Context, sourceID int64) ([]int64, error) {
	var userIDs []int64
	if err := r.pg.QuerySlice(ctx, &userIDs, getPermittedUsers, sourceID); err != nil {
		return nil, errs.WrapErr(err, "get permitted users")
	}
	return userIDs, nil
}
