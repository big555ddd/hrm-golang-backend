package holiday

import (
	organization "app/app/modules/organiztion"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, organization *organization.Module) *Module {
	svc := NewService(db, organization)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
