package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/domain/internal/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const updateSource = `
	update domain.source
	set external_id = $4,
	    title = $5,
	    content = $6,
	    update_every_period = $7,
	    cron_week_day = $8,
	    cron_month = $9,
	    cron_day = $10,
	    cron_hour = $11,
	    cron_minute = $12,
		credentials = $13,
		status = $14
	where internal_id = (
	    select internal_source_id
	    from domain.get_permitted_sources($2, $3)
	    where internal_source_id = $1
	);
`

// UpdateSource updates data for source.
func (r *SourceRepo) UpdateSource(ctx context.Context, s model.SourceDao, userID int64, roleIDs []int64) error {
	rows, err := r.pg.Exec(
		ctx,
		updateSource,
		s.ID,
		userID,
		roleIDs,
		s.ExtID,
		s.Title,
		s.Content,
		s.UpdateEveryPeriod,
		s.CronWeekDay,
		s.CronMonth,
		s.CronDay,
		s.CronHour,
		s.CronMinute,
		s.Credentials,
		s.Status,
	)
	if err != nil {
		return errs.WrapErr(err, "update source")
	}

	if rows == 0 {
		return errs.WrapErr(pgx.ErrNoRows, "source not found")
	}

	return nil
}
