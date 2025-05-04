package model

import (
	"time"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// ToProto converts dao model into protobuf format.
func (s *SourceDao) ToProto() *pb.Source {
	var updateParams *pb.UpdateParams = nil
	updateParamsDto := s.AssembleUpdateParams()
	if updateParamsDto != nil {
		updateParams = updateParamsDto.ToProto()
	}

	return &pb.Source{
		Id:           s.ID,
		UserId:       s.UserID,
		Title:        s.Title,
		Content:      s.Content,
		Typ:          pb.SourceType(s.Type),
		UpdateParams: updateParams,
		Credentials:  s.Credentials,
		Status:       pb.SourceStatus(s.Status),
		CreatedAt:    timestamppb.New(s.CreatedAt),
		UpdatedAt:    timestamppb.New(s.UpdatedAt),
	}
}

func (s *SourceDao) setUpdateParamsDefaults() {
	s.UpdateEveryPeriod = -1
	s.CronWeekDay = -1
	s.CronMonth = -1
	s.CronDay = -1
	s.CronHour = -1
	s.CronMinute = -1
}

// FillUpdateParams sets update params value from protobuf format or uses default value -1.
func (s *SourceDao) FillUpdateParams(updateParams *pb.UpdateParams) {
	s.setUpdateParamsDefaults()
	if updateParams != nil {
		if updateParams.EveryPeriod != nil {
			s.UpdateEveryPeriod = updateParams.GetEveryPeriod()
		} else {
			cron := updateParams.GetCron()
			s.CronWeekDay = cron.GetDayOfWeek()
			s.CronMonth = cron.GetMonth()
			s.CronDay = cron.GetDayOfMonth()
			s.CronHour = cron.GetHour()
			s.CronMinute = cron.GetMinute()
		}
	}
}

// AssembleUpdateParams returns separate update params assembled into a single struct.
func (s *SourceDao) AssembleUpdateParams() *UpdateParams {
	var updateParams UpdateParams

	switch {
	case s.UpdateEveryPeriod != -1:
		updateParams.EveryPeriod = &s.UpdateEveryPeriod
	case s.UpdateEveryPeriod != -1 && s.CronWeekDay != -1 && s.CronMonth != -1 && s.CronDay != -1 && s.CronHour != -1 && s.CronMinute != -1:
		updateParams.Cron = &Cron{
			WeekDay: s.CronWeekDay,
			Month:   s.CronMonth,
			Day:     s.CronDay,
			Hour:    s.CronHour,
			Minute:  s.CronMinute,
		}
	default:
		return nil
	}

	return &updateParams
}
