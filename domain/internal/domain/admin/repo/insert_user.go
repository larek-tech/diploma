package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/admin/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertUser = `
	insert into auth.user(email, hash_password)
	values ($1, $2)
	returning id;
`

// InsertUser creates new user.
func (r *Repo) InsertUser(ctx context.Context, u model.UserDao) (int64, error) {
	var userID int64
	if err := r.pg.Query(ctx, &userID, insertUser, u.Email, u.HashPassword); err != nil {
		return 0, errs.WrapErr(err, "insert user")
	}
	return userID, nil
}
