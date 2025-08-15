package model

import (
	"github.com/uptrace/bun"
)

type DocumentLeave struct {
	bun.BaseModel `bun:"table:document_leaves"`

	ID          string  `bun:",default:gen_random_uuid(),pk"`
	DocumentID  string  `bun:"document_id,notnull"`
	LeaveID     string  `bun:"leave_id,notnull"`
	Description string  `bun:"description"`
	StartDate   int64   `bun:"start_date,notnull"`
	EndDate     int64   `bun:"end_date,notnull"`
	LeaveHours  float64 `bun:"leave_hours,notnull"`
	UsedQuota   float64 `bun:"used_quota,notnull"`

	CreateUpdateUnixTimestamp
}
