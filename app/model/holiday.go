package model

import (
	"github.com/uptrace/bun"
)

type Holiday struct {
	bun.BaseModel `bun:"table:holidays"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	Name        string `bun:"name,notnull"`
	IsActive    bool   `bun:"is_active,notnull"`
	Description string `bun:"description"`
	StartDate   int64  `bun:"start_date,notnull"`
	EndDate     int64  `bun:"end_date,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
