package repo

import (
	"context"
	"github.com/larek-tech/diploma/auth/internal/auth/model"
	"github.com/yogenyslav/pkg/errs"
)

const findOneByEmail = `
	select id, email, hash_password, created_at, updated_at, is_deleted
	from auth.user
	where email = $1;
`

// FindOneByEmail returns a user filtered by email.
func (r *AuthRepo) FindOneByEmail(ctx context.Context, email string) (model.UserDao, error) {
	var user model.UserDao
	if err := r.pg.Query(ctx, &user, findOneByEmail, email); err != nil {
		return user, errs.WrapErr(err, "find one by email")
	}
	return user, nil
}
