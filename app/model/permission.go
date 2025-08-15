package model

import (
	"github.com/uptrace/bun"
)

type Permission struct {
	bun.BaseModel `bun:"table:permissions"`

	ID          string `bun:",default:gen_random_uuid(),pk"`
	Name        string `bun:"name,notnull,unique"`
	Description string `bun:"description"`

	CreateUpdateUnixTimestamp
}
