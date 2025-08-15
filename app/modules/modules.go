package modules

import (
	"app/app/modules/attendance"
	"app/app/modules/auth"
	"app/app/modules/branch"
	"app/app/modules/department"
	"app/app/modules/document"
	"app/app/modules/holiday"
	"app/app/modules/leave"
	"app/app/modules/notification"
	organization "app/app/modules/organiztion"
	"app/app/modules/redisctl"
	"app/app/modules/role"
	"app/app/modules/user"
	"app/app/modules/workshift"
	"app/config"
)

type Module struct {
	Redis        *redisctl.Module
	Auth         *auth.Module
	User         *user.Module
	Role         *role.Module
	Department   *department.Module
	Branch       *branch.Module
	Organization *organization.Module
	Holiday      *holiday.Module
	WorkShift    *workshift.Module
	Leave        *leave.Module
	Document     *document.Module
	Attendance   *attendance.Module
	Notification *notification.Module
}

func New() *Module {
	db := config.GetDB()
	redis := config.GetRedis()
	redisctl := redisctl.NewModule(redis)

	// Create notification module first
	notification := notification.NewModule(db)

	organization := organization.NewModule(db)
	branch := branch.NewModule(db, organization)
	department := department.NewModule(db, organization, branch)
	workshift := workshift.NewModule(db)
	role := role.NewModule(db)
	user := user.NewModule(db, department, role, workshift)
	attendance := attendance.NewModule(db, user, workshift)
	leave := leave.NewModule(db, user)
	holiday := holiday.NewModule(db, organization)

	// Pass the same notification instance to all modules that need it
	auth := auth.NewModule(db, user, department, role, workshift, leave, attendance, holiday, notification)
	document := document.NewModule(db, user, leave, workshift, attendance, notification)

	return &Module{
		Redis:        redisctl,
		User:         user,
		Auth:         auth,
		Role:         role,
		Department:   department,
		Branch:       branch,
		Organization: organization,
		Holiday:      holiday,
		WorkShift:    workshift,
		Leave:        leave,
		Document:     document,
		Attendance:   attendance,
		Notification: notification,
	}
}
