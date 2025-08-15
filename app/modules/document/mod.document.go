package document

import (
	"app/app/modules/attendance"
	"app/app/modules/leave"
	"app/app/modules/notification"
	"app/app/modules/user"
	"app/app/modules/workshift"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, user *user.Module,
	leave *leave.Module, workshift *workshift.Module,
	attendance *attendance.Module, notification *notification.Module) *Module {
	svc := NewService(db, user, leave, workshift, attendance, notification)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
