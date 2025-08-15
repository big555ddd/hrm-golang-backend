package model

import (
	"github.com/uptrace/bun"
)

type Leave struct {
	bun.BaseModel `bun:"table:leaves"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`
	Year        int    `bun:"year,notnull"`
	Amount      int64  `bun:"amount,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
