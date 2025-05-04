package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const getSourceByID = `
	select internal_id, external_id, user_id, title, content, type, update_every_period, cron_week_day, cron_month, cron_day, cron_hour, cron_minute, credentials, status, created_at, updated_at
	from domain.source
		where internal_id = (
			select internal_source_id
			from domain.get_permitted_sources($2, $3)
			where internal_source_id = $1
		);
`

// GetSourceByID returns source by ID.
func (r *SourceRepo) GetSourceByID(ctx context.Context, id, userID int64, roleIDs []int64) (model.SourceDao, error) {
	var source model.SourceDao
	if err := r.pg.Query(ctx, &source, getSourceByID, id, userID, roleIDs); err != nil {
		return source, errs.WrapErr(err, "get source by id")
	}
	return source, nil
}
