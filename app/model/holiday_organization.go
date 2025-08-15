package model

import (
	"github.com/uptrace/bun"
)

type HolidayOrganization struct {
	bun.BaseModel `bun:"table:holiday_organization"`

	HolidayID      string `bun:"holiday_id,notnull"`
	OrganizationID string `bun:"organization_id,notnull"`
}
