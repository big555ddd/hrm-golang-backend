package leave

import (
	"app/app/modules/user"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, user *user.Module) *Module {
	svc := NewService(db, user)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
