package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/role/model"
	"github.com/yogenyslav/pkg/errs"
)

const listRoles = `
	select id, name, created_at
	from auth.role
	where is_deleted = false
	order by created_at desc
	offset $1
	limit $2;
`

// ListRoles returns paginated list of roles.
func (r *Repo) ListRoles(ctx context.Context, offset, limit uint64) ([]model.RoleDao, error) {
	var roles []model.RoleDao
	if err := r.pg.QuerySlice(ctx, &roles, listRoles, offset, limit); err != nil {
		return nil, errs.WrapErr(err, "list roles")
	}
	return roles, nil
}
