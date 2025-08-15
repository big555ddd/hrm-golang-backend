package model

import (
	"app/app/enum"

	"github.com/uptrace/bun"
)

type Document struct {
	bun.BaseModel `bun:"table:documents"`

	ID          string              `bun:",default:gen_random_uuid(),pk"`
	UserID      string              `bun:"user_id,notnull"`
	Status      enum.StatusDocument `bun:"status,notnull"`
	Type        enum.DocumentType   `bun:"type,notnull"`
	Description string              `bun:"description"`
	Approved    []string            `bun:"approved,type:jsonb"`
	Rejected    []string            `bun:"rejected,type:jsonb"`

	CreateUpdateUnixTimestamp
}
