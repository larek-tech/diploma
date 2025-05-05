package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
)

const deleteOldUserPermissions = `
	delete from domain.domain_permitted_users
	where domain_id = $1;
`

const insertNewUserPermissions = `
	insert into domain.domain_permitted_users (domain_id, user_id)
	select $1, u.id
	from unnest($2::bigint[]) as uid
	join auth.user u on u.id = uid
	returning uid;
`

// UpdatePermittedUsers replaces old user domain permissions with new.
func (r *Repo) UpdatePermittedUsers(ctx context.Context, domainID int64, userIDs []int64) ([]int64, error) {
	var resUserIDs []int64

	ctx, err := r.pg.BeginSerializable(ctx)
	if err != nil {
		return nil, errs.WrapErr(err, "start tx")
	}
	defer func() {
		if e := r.pg.RollbackTx(ctx); e != nil {
			log.Warn().Err(errs.WrapErr(err)).Msg("rollback tx")
		}
	}()

	if _, err = r.pg.ExecTx(ctx, deleteOldUserPermissions, domainID); err != nil {
		return nil, errs.WrapErr(err, "delete old user permissions")
	}

	if err = r.pg.QuerySliceTx(ctx, &resUserIDs, insertNewUserPermissions, domainID, userIDs); err != nil {
		return nil, errs.WrapErr(err, "insert new user permissions")
	}

	if err = r.pg.CommitTx(ctx); err != nil {
		return nil, errs.WrapErr(err, "commit tx")
	}

	return resUserIDs, nil
}
