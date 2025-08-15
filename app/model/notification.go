package model

import (
	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID         string `bun:",default:gen_random_uuid(),pk"`
	UserID     string `bun:"user_id,notnull"`
	Type       string `bun:"type,notnull"`
	Message    string `bun:"message,notnull"`
	IsRead     bool   `bun:"is_read,default:false"`
	DocumentID string `bun:"document_id"`

	CreateUpdateUnixTimestamp
}
