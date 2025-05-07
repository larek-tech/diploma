package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/pkg/errs"
)

const deleteScenario = `
	delete from domain.scenario
	where id = $1
		and user_id = $2;
`

// DeleteScenario deletes scenario by ID.
func (r *Repo) DeleteScenario(ctx context.Context, id, userID int64) error {
	rows, err := r.pg.Exec(ctx, deleteScenario, id, userID)
	if err != nil {
		return errs.WrapErr(err, "delete scenario")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "scenario not found")
	}

	return nil
}
