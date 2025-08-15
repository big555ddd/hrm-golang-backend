package attendance

import (
	"app/app/modules/user"
	"app/app/modules/workshift"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, user *user.Module, workshift *workshift.Module) *Module {
	svc := NewService(db, user, workshift)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
