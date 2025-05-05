package repo

import (
	"context"

	"github.com/yogenyslav/pkg/errs"
)

const getPermittedRoles = `
	select role_id
	from domain.source_permitted_roles
	where internal_source_id = $1;
`

// GetPermittedRoles returns list of roles that have direct access to the source.
func (r *Repo) GetPermittedRoles(ctx context.Context, sourceID int64) ([]int64, error) {
	var userIDs []int64
	if err := r.pg.QuerySlice(ctx, &userIDs, getPermittedRoles, sourceID); err != nil {
		return nil, errs.WrapErr(err, "get permitted roles")
	}
	return userIDs, nil
}
