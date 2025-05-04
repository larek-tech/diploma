package model

import (
	"time"
)

// SourceDao is a model for source on data layer.
type SourceDao struct {
	ID                int64        `db:"internal_id"`
	ExtID             string       `db:"external_id"`
	UserID            int64        `db:"user_id"`
	Title             string       `db:"title"`
	Content           []byte       `db:"content"`
	Type              SourceType   `db:"type"`
	UpdateEveryPeriod int64        `db:"update_every_period"`
	CronWeekDay       int32        `db:"cron_week_day"`
	CronMonth         int32        `db:"cron_month"`
	CronDay           int32        `db:"cron_day"`
	CronHour          int32        `db:"cron_hour"`
	CronMinute        int32        `db:"cron_minute"`
	Credentials       []byte       `db:"credentials"`
	Status            SourceStatus `db:"status"`
	CreatedAt         time.Time    `db:"created_at"`
	UpdatedAt         time.Time    `db:"updated_at"`
}

// SetDefaultUpdateParams sets update params to default values.
func (s *SourceDao) SetDefaultUpdateParams() {
	s.UpdateEveryPeriod = -1
	s.CronWeekDay = -1
	s.CronMonth = -1
	s.CronDay = -1
	s.CronHour = -1
	s.CronMinute = -1
}

// HasUpdateParams returns true if there is at least one non-default value in update parameters.
func (s *SourceDao) HasUpdateParams() bool {
	return s.UpdateEveryPeriod != -1 &&
		s.CronWeekDay != -1 &&
		s.CronMonth != -1 &&
		s.CronDay != -1 &&
		s.CronHour != -1 &&
		s.CronMinute != -1
}

// DomainDao is a model for domain on data layer.
type DomainDao struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	UserID    int64     `db:"user_id"`
	SourceIDs []string  `db:"source_ids"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
