package auth

import (
	"app/app/modules/attendance"
	"app/app/modules/department"
	"app/app/modules/holiday"
	"app/app/modules/leave"
	"app/app/modules/notification"
	"app/app/modules/role"
	"app/app/modules/user"
	"app/app/modules/workshift"

	"github.com/uptrace/bun"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(db *bun.DB, user *user.Module,
	department *department.Module, role *role.Module,
	workshift *workshift.Module, leave *leave.Module,
	attendance *attendance.Module, holiday *holiday.Module,
	notification *notification.Module) *Module {
	svc := NewService(db, user, department, role,
		workshift, leave, attendance,
		holiday, notification)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
