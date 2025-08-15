package department

import (
	"app/app/modules/branch"
	organization "app/app/modules/organiztion"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, organization *organization.Module, branch *branch.Module) *Module {
	svc := NewService(db, organization, branch)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
