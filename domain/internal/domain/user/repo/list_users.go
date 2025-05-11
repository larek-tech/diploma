package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/user/model"
	"github.com/yogenyslav/pkg/errs"
)

const listUsers = `
	select id, email, created_at, updated_at
	from auth.user
	where is_deleted = false
	order by created_at desc, updated_at desc
	offset $1
	limit $2;
`

// ListUsers returns paginated list of users.
func (r *Repo) ListUsers(ctx context.Context, offset, limit uint64) ([]model.UserDao, error) {
	var users []model.UserDao
	if err := r.pg.QuerySlice(ctx, &users, listUsers, offset, limit); err != nil {
		return nil, errs.WrapErr(err, "list users")
	}
	return users, nil
}
