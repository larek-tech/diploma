package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertDomain = `
	insert into domain.domain(title, user_id, source_ids)
	values ($1, $2, (
        select coalesce(array_agg(s.internal_id), '{}')
        from domain.source s
        where s.internal_id = any($3)
    ))
	returning id;
`

// InsertDomain create new domain record.
func (r *Repo) InsertDomain(ctx context.Context, d model.DomainDao) (int64, error) {
	var domainID int64
	err := r.pg.Query(
		ctx,
		&domainID,
		insertDomain,
		d.Title,
		d.UserID,
		d.SourceIDs,
	)
	if err != nil {
		return 0, errs.WrapErr(err, "insert domain")
	}
	return domainID, nil
}
