package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/role/model"
	"github.com/yogenyslav/pkg/errs"
)

const getRole = `
	select id, name, created_at
	from auth.role
	where id = $1
		and is_deleted = false;
`

// GetRole returns role.
func (r *Repo) GetRole(ctx context.Context, id int64) (model.RoleDao, error) {
	var role model.RoleDao
	if err := r.pg.Query(ctx, &role, getRole, id); err != nil {
		return role, errs.WrapErr(err, "get role")
	}
	return role, nil
}
