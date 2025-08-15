package model

import (
	"github.com/uptrace/bun"
)

type UserRole struct {
	bun.BaseModel `bun:"table:user_roles"`

	UserID string `bun:"user_id,notnull"`
	RoleID string `bun:"role_id,notnull"`
}
