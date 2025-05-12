package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const deleteDomain = `
	delete from domain.domain
	where id in (
	    select id
	    from domain.get_permitted_domains($2, $3)
	    where id = $1
	);
`

// DeleteDomain deletes domain by ID.
func (r *Repo) DeleteDomain(ctx context.Context, id, userID int64, roleIDs []int64) error {
	rows, err := r.pg.Exec(ctx, deleteDomain, id, userID, roleIDs)
	if err != nil {
		return errs.WrapErr(err, "delete domain")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "domain not found")
	}

	return nil
}
