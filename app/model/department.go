package model

import (
	"github.com/uptrace/bun"
)

type Department struct {
	bun.BaseModel `bun:"table:departments"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`
	BranchID    string `bun:"branch_id,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
