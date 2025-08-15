package model

import (
	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:organizations"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
