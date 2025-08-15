package model

import (
	"github.com/uptrace/bun"
)

type WorkShift struct {
	bun.BaseModel `bun:"table:work_shifts"`

	ID            string  `bun:",default:gen_random_uuid(),pk"`
	Name          string  `bun:"name,notnull"`
	WorkLocationX float64 `bun:"work_location_x,notnull"`
	WorkLocationY float64 `bun:"work_location_y,notnull"`
	Description   string  `bun:"description"`
	LateMinutes   int64   `bun:"late_minutes,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
