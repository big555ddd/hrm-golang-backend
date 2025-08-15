package activitylog

import "github.com/uptrace/bun"

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB) *Module {
	svc := NewService(db)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
