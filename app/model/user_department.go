package model

import (
	"github.com/uptrace/bun"
)

type UserDepartment struct {
	bun.BaseModel `bun:"table:user_departments"`

	UserID       string `bun:"user_id,notnull"`
	DepartmentID string `bun:"department_id,notnull"`
}
