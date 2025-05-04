package repo

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/source/model"
	"github.com/yogenyslav/pkg/errs"
)

const insertSource = `
	insert into domain.source(user_id, title, content, type, update_every_period, cron_week_day, cron_month, cron_day, cron_hour, cron_minute, credentials, status)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	returning internal_id;
`

// InsertSource create new source record.
func (r *Repo) InsertSource(ctx context.Context, s model.SourceDao) (int64, error) {
	var sourceID int64
	err := r.pg.Query(
		ctx,
		&sourceID,
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
		return 0, errs.WrapErr(err, "insert source")
	}
	return sourceID, nil
}
