package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const listSources = `
	select source.internal_id, external_id, user_id, title, content, type, update_every_period, cron_week_day, cron_month, cron_day, cron_hour, cron_minute, credentials, status, created_at, updated_at
	from domain.source
		where internal_id = (
			select internal_source_id
			from domain.get_permitted_sources($1, $2)
		);
`

// ListSources returns list of sources available for user.
func (r *SourceRepo) ListSources(ctx context.Context, userID int64, roleIDs []int64) ([]model.SourceDao, error) {
	var sources []model.SourceDao
	if err := r.pg.QuerySlice(ctx, &sources, listSources, userID, roleIDs); err != nil {
		return sources, errs.WrapErr(err, "list sources")
	}
	return sources, nil
}
