package model

import (
	"github.com/uptrace/bun"
)

type RolePermission struct {
	bun.BaseModel `bun:"table:role_permissions"`

	RoleID       string `bun:"role_id,notnull"`
	PermissionID string `bun:"permission_id,notnull"`
}
