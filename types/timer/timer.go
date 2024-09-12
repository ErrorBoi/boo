package timer

import (
	"database/sql"
	"time"
)

type TimerType string

func (t TimerType) String() string {
	return string(t)
}

const (
	Daily      TimerType = "daily"
	Periodical TimerType = "periodical"
)

type RepeatType string

func (t *RepeatType) String() string {
	if t == nil {
		return ""
	}

	return string(*t)
}

const (
	OnceMissed RepeatType = "once missed"
	Always     RepeatType = "always"
)

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type Timer struct {
	ID            int64          `db:"id"`
	TgID          int64          `db:"tg_id"`
	Name          string         `db:"name"`
	Type          *TimerType     `db:"timer_type"`
	RepeatType    RepeatType     `db:"repeat_type"`
	CreatedAt     time.Time      `db:"created_at"`
	Description   sql.NullString `db:"description"`
	PeriodSeconds sql.NullInt64  `db:"period_seconds"`
	TriggerTime   sql.NullString `db:"trigger_time"`
	LastTrigger   time.Time      `db:"last_trigger"`
	Status        Status         `db:"status"`
	LastAck       time.Time      `db:"last_ack"`
	NextTrigger   time.Time      `db:"next_trigger"`
	Link          sql.NullString `db:"link"`
	ImgUrl        sql.NullString `db:"img_url"`
	NotifyID      sql.NullInt64  `db:"notify_id"`
}

func NewTimer(tgID int64, name string) *Timer {
	return &Timer{
		TgID: tgID,
		Name: name,
	}
}

func (t *Timer) IsOwnedBy(userID int64) bool {
	return t.TgID == userID
}
