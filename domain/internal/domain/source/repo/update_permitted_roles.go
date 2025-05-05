package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
)

const deleteOldRolePermissions = `
	delete from domain.source_permitted_roles
	where internal_source_id = $1;
`

const insertNewRolePermissions = `
	insert into domain.source_permitted_roles (internal_source_id, role_id)
	select $1, r.id
	from unnest($2::bigint[]) as rid
	join auth.role r on r.id = rid
	returning rid;
`

// UpdatePermittedRoles replaces old role source permissions with new.
func (r *Repo) UpdatePermittedRoles(ctx context.Context, sourceID int64, roleIDs []int64) ([]int64, error) {
	var resRoleIDs []int64

	ctx, err := r.pg.BeginSerializable(ctx)
	if err != nil {
		return nil, errs.WrapErr(err, "start tx")
	}
	defer func() {
		if e := r.pg.RollbackTx(ctx); e != nil {
			log.Warn().Err(errs.WrapErr(err)).Msg("rollback tx")
		}
	}()

	if _, err = r.pg.ExecTx(ctx, deleteOldRolePermissions, sourceID); err != nil {
		return nil, errs.WrapErr(err, "delete old role permissions")
	}

	if err = r.pg.QuerySliceTx(ctx, &resRoleIDs, insertNewRolePermissions, sourceID, roleIDs); err != nil {
		return nil, errs.WrapErr(err, "insert new role permissions")
	}

	if err = r.pg.CommitTx(ctx); err != nil {
		return nil, errs.WrapErr(err, "commit tx")
	}

	return resRoleIDs, nil
}
