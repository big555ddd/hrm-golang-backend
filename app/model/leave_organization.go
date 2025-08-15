package model

import (
	"github.com/uptrace/bun"
)

type LeaveOrganization struct {
	bun.BaseModel `bun:"table:leave_organizations"`

	LeaveID        string `bun:"leave_id,notnull"`
	OrganizationID string `bun:"organization_id,notnull"`
}
