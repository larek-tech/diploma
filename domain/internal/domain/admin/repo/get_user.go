package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/admin/model"
	"github.com/yogenyslav/pkg/errs"
)

const getUser = `
	select id, email, hash_password, created_at, updated_at
	from auth.user
	where id = $1
		and is_deleted = false;
`

// GetUser returns user.
func (r *Repo) GetUser(ctx context.Context, id int64) (model.UserDao, error) {
	var user model.UserDao
	if err := r.pg.Query(ctx, &user, getUser, id); err != nil {
		return user, errs.WrapErr(err, "get user")
	}
	return user, nil
}
