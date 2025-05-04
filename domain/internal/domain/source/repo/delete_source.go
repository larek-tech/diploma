package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const deleteSource = `
	delete from domain.source
	where internal_id = (
	    select internal_source_id
	    from domain.get_permitted_sources($2, $3)
	    where internal_source_id = $1
	);
`

// DeleteSource deletes source by ID.
func (r *Repo) DeleteSource(ctx context.Context, id, userID int64, roleIDs []int64) error {
	rows, err := r.pg.Exec(ctx, deleteSource, id, userID, roleIDs)
	if err != nil {
		return errs.WrapErr(err, "delete source")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "source not found")
	}

	return nil
}
