package repo

import (
	"context"

	"github.com/yogenyslav/pkg/errs"
)

const getPermittedRoles = `
	select role_id
	from domain.domain_permitted_roles
	where domain_id = $1;
`

// GetPermittedRoles returns list of roles that have access to the domain.
func (r *Repo) GetPermittedRoles(ctx context.Context, domainID int64) ([]int64, error) {
	var userIDs []int64
	if err := r.pg.QuerySlice(ctx, &userIDs, getPermittedRoles, domainID); err != nil {
		return nil, errs.WrapErr(err, "get permitted roles")
	}
	return userIDs, nil
}
