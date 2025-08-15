package model

import (
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
