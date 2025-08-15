package model

import (
	"github.com/uptrace/bun"
)

type Branch struct {
	bun.BaseModel `bun:"table:branches"`

	ID             string `bun:",default:gen_random_uuid(),pk"`
	Name           string `bun:"name,notnull"`
	Description    string `bun:"description"`
	OrganizationID string `bun:"organization_id,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
