package authdto

import (
	"app/app/model"
	leavedto "app/app/modules/leave/dto"
	roledto "app/app/modules/role/dto"
	userdto "app/app/modules/user/dto"
	workshiftdto "app/app/modules/workshift/dto"
)

type LoginRequest struct {
	EmpCode  string `json:"empCode"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	Ref   string `json:"ref"`
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Otp         string `json:"otp"`
	Ref         string `json:"ref"`
	NewPassword string `json:"newPassword"`
}

type UserInfoResponse struct {
	userdto.ListUserResponse
	Permissions       []roledto.RolePermissionResponse `json:"permissions"`
	ShiftSchedules    []workshiftdto.ShiftSchedule     `json:"shiftSchedules"`
	Leaves            []leavedto.LeaveUserResponse     `json:"leaves"`
	AttendanceToday   model.Attendance                 `json:"attendanceToday"`
	Holidays          []model.Holiday                  `json:"holidays"`
	NotificationCount int                              `json:"notificationCount"`
}
