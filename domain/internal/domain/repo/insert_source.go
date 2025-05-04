package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertSource = `
	insert into domain.source(user_id, title, content, type, update_every_period, cron_week_day, cron_month, cron_day, cron_hour, cron_minute, credentials, status)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
`

// InsertSource create new source record.
func (r *SourceRepo) InsertSource(ctx context.Context, s model.SourceDao) error {
	_, err := r.pg.Exec(
		ctx,
		insertSource,
		s.UserID,
		s.Title,
		s.Content,
		s.Type,
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
		return errs.WrapErr(err, "insert source")
	}
	return nil
}
