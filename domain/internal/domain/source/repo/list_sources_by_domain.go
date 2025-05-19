package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/source/model"
	"github.com/yogenyslav/pkg/errs"
)

const listSourcesByDomainQuery = `
	with permitted_sources as (
		select internal_source_id
		from domain.get_permitted_sources($1, $2)
		intersect 
		select unnest(source_ids)
		from domain.domain
		where id = $5
	)
	select internal_id, external_id, user_id, title, content, type, update_every_period, cron_week_day, cron_month, cron_day, cron_hour, cron_minute, credentials, status, created_at, updated_at
	from domain.source
		where internal_id in (
			select internal_source_id
			from permitted_sources
		)
	order by created_at desc, updated_at desc
	offset $3
	limit $4;
`

// ListSourcesByDomain returns paginated list of sources by specified domain.
func (r *Repo) ListSourcesByDomain(
	ctx context.Context,
	userID, domainID int64,
	roles []int64,
	offset, limit uint64,
) ([]model.SourceDao, error) {
	var sources []model.SourceDao
	if err := r.pg.QuerySlice(ctx, &sources, listSourcesByDomainQuery, userID, roles, offset, limit, domainID); err != nil {
		return sources, errs.WrapErr(err, "list sources by domain")
	}
	return sources, nil
}
