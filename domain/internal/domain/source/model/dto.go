package model

import (
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
)

// Cron contains cron-format parameters for source updates.
type Cron struct {
	WeekDay int32 `db:"cron_week_day"`
	Month   int32 `db:"cron_month"`
	Day     int32 `db:"cron_day"`
	Hour    int32 `db:"cron_hour"`
	Minute  int32 `db:"cron_minute"`
}

// ToProto converts dto model to protobuf format.
func (c *Cron) ToProto() *pb.CronFormat {
	return &pb.CronFormat{
		Minute:     c.Minute,
		Hour:       c.Hour,
		DayOfMonth: c.Day,
		Month:      c.Month,
		DayOfWeek:  c.WeekDay,
	}
}

// UpdateParams sets time conditions to parse dynamic source (not static files).
type UpdateParams struct {
	EveryPeriod *int64 `json:"every_period,omitempty"` // update every N seconds
	Cron        *Cron  `json:"cron,omitempty"`         // update on date/time (cron-format)
}

// ToProto converts dto model into protobuf format.
func (up *UpdateParams) ToProto() *pb.UpdateParams {
	var cron *pb.CronFormat = nil

	if up.Cron != nil {
		cron = up.Cron.ToProto()
	}

	return &pb.UpdateParams{
		EveryPeriod: up.EveryPeriod,
		Cron:        cron,
	}
}

// DataMessage contains information about new Source and is sent to Data service to be processed.
type DataMessage struct {
	Title        string        `json:"title"`
	Content      []byte        `json:"content"` // byte encoded url or file content
	Type         SourceType    `json:"type"`
	Credentials  []byte        `json:"credentials,omitempty"`
	UpdateParams *UpdateParams `json:"update_params,omitempty"`
}

// ParsingStatus status of processing source.
type ParsingStatus struct {
	SourceID  string       `json:"source"`
	JobID     string       `json:"job_id"`
	Processed int          `json:"processed"`
	Total     int          `json:"total"`
	Status    SourceStatus `json:"status"`
}
