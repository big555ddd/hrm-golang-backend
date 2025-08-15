package migrations

func Models() []any {
	return []any{
		// (*model.ActivityLog)(nil),
		// (*model.Attendance)(nil),
		// (*model.Branch)(nil),
		// (*model.Department)(nil),
		// (*model.DocumentLeave)(nil),
		// (*model.DocumentOvertime)(nil),
		// (*model.Document)(nil),
		// (*model.HolidayOrganization)(nil),
		// (*model.Holiday)(nil),
		// (*model.LeaveOrganization)(nil),
		// (*model.Leave)(nil),
		// (*model.Notification)(nil),
		// (*model.Organization)(nil),
		// (*model.Permission)(nil),
		// (*model.RolePermission)(nil),
		// (*model.Role)(nil),
		// (*model.UserDepartment)(nil),
		// (*model.UserRole)(nil),
		// (*model.User)(nil),
		// (*model.UserForgot)(nil),
		// (*model.UserWorkShift)(nil),
		// (*model.ShiftSchedule)(nil),
		// (*model.WorkShift)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{}
}
