package model

import (
	"github.com/uptrace/bun"
)

type UserWorkShift struct {
	bun.BaseModel `bun:"table:user_work_shifts"`

	UserID      string `bun:"user_id,notnull"`
	WorkShiftID string `bun:"work_shift_id,notnull"`
}
