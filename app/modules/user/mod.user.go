package user

import (
	"app/app/modules/department"
	"app/app/modules/role"
	"app/app/modules/workshift"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, department *department.Module, role *role.Module, workshift *workshift.Module) *Module {
	svc := NewService(db, department, role, workshift)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
