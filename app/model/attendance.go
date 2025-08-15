package model

import (
	"github.com/uptrace/bun"
)

type Attendance struct {
	bun.BaseModel `bun:"table:attendances"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	UserID      string `bun:"user_id,notnull"`
	WorkShiftID string `bun:"work_shift_id,notnull"`
	CheckIn     int64  `bun:"check_in,notnull"`
	CheckOut    int64  `bun:"check_out,notnull"`
	Date        int64  `bun:"date,notnull"`
	IsOnTime    bool   `bun:"is_on_time,default:false"`
	IsLate      bool   `bun:"is_late,default:false"`
	IsLeave     bool   `bun:"is_leave,default:false"`

	CreateUpdateUnixTimestamp
}
