package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type DocumentOvertime struct {
	bun.BaseModel `bun:"table:document_overtimes"`

	ID              string            `bun:",default:gen_random_uuid(),pk"`
	DocumentID      string            `bun:"document_id,notnull"`
	OverTimeType    enum.OverTimeType `bun:"overtime_type,notnull"`
	Description     string            `bun:"description"`
	StartDate       int64             `bun:"start_date,notnull"`
	EndDate         int64             `bun:"end_date,notnull"`
	DurationMinutes int64             `bun:"duration_minutes,notnull"`

	CreateUpdateUnixTimestamp
}
