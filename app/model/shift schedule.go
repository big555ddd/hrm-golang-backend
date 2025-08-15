package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type ShiftSchedule struct {
	bun.BaseModel `bun:"table:shift_schedules"`

	WorkShiftID string   `bun:"work_shift_id,notnull"`
	Day         enum.Day `bun:"day,notnull"`
	StartTime   int64    `bun:"start_time,notnull"`
	EndTime     int64    `bun:"end_time,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
