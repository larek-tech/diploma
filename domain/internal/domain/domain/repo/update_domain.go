package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const updateDomains = `
	update domain.domain
	set title = $4,
	    source_ids = (
			select coalesce(array_agg(s.internal_id), '{}')
			from domain.source s
			where s.internal_id = any($5)
		),
		scenario_ids = $6
	where id in (
	    select id
	    from domain.get_permitted_domains($2, $3)
	    where id = $1
	);
`

// UpdateDomain updates data for domain.
func (r *Repo) UpdateDomain(ctx context.Context, d model.DomainDao, userID int64, roleIDs []int64) error {
	rows, err := r.pg.Exec(
		ctx,
		updateDomains,
		d.ID,
		userID,
		roleIDs,
		d.Title,
		d.SourceIDs,
		d.ScenarioIds,
	)
	if err != nil {
		return errs.WrapErr(err, "update domain")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "domain not found")
	}

	return nil
}
