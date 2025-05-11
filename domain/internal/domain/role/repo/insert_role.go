package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/role/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertRole = `
	insert into auth.role(name)
	values ($1)
	returning id;
`

// InsertRole creates new role.
func (r *Repo) InsertRole(ctx context.Context, u model.RoleDao) (int64, error) {
	var roleID int64
	if err := r.pg.Query(ctx, &roleID, insertRole, u.Name); err != nil {
		return 0, errs.WrapErr(err, "insert role")
	}
	return roleID, nil
}
