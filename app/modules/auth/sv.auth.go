package auth

import (
	"app/app/helper"
	message "app/app/messsage"
	"app/app/model"
	"app/app/modules/attendance"
	authdto "app/app/modules/auth/dto"
	"app/app/modules/department"
	"app/app/modules/holiday"
	"app/app/modules/leave"
	"app/app/modules/notification"
	"app/app/modules/role"
	"app/app/modules/user"
	userdto "app/app/modules/user/dto"
	"app/app/modules/workshift"
	workshiftdto "app/app/modules/workshift/dto"
	"app/app/util/jwt"
	"app/config"
	"app/internal/logger"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	db           *bun.DB
	user         *user.Module
	department   *department.Module
	role         *role.Module
	workshift    *workshift.Module
	leave        *leave.Module
	attendance   *attendance.Module
	holiday      *holiday.Module
	notification *notification.Module
}

func NewService(db *bun.DB, user *user.Module,
	department *department.Module, role *role.Module,
	workshift *workshift.Module, leave *leave.Module,
	attendance *attendance.Module, holiday *holiday.Module,
	notification *notification.Module) *Service {
	return &Service{
		db:           db,
		user:         user,
		department:   department,
		role:         role,
		workshift:    workshift,
		leave:        leave,
		attendance:   attendance,
		holiday:      holiday,
		notification: notification,
	}
}

func (s *Service) Login(ctx context.Context, req *authdto.LoginRequest) (string, error) {
	user, err := s.user.Svc.GetByEmpCode(ctx, req.EmpCode)
	if err != nil {
		return "", err
	}

	if err := s.user.Svc.CheckPassword(ctx, user, req.Password); err != nil {
		return "", err
	}

	data := jwt.ClaimData{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	token, err := jwt.CreateToken(data)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ForgotPassword(ctx context.Context, req *authdto.ForgotPasswordRequest) (*authdto.ForgotPasswordResponse, error) {
	resp := new(authdto.ForgotPasswordResponse)
	user, err := s.user.Svc.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	ref, err := helper.GenerateREFCode(6)
	if err != nil {
		return nil, err
	}

	// Check if reference code already exists, regenerate if it does
	for {
		exists, err := s.db.NewSelect().
			Table("user_forgots").
			Where("ref = ?", ref).
			Where("used = false").
			Where("expires > ?", time.Now().Unix()).
			Exists(ctx)
		if err != nil {
			return nil, err
		}

		if !exists {
			break // Reference code is unique, exit loop
		}

		// Regenerate reference code if it already exists
		ref, err = helper.GenerateREFCode(6)
		if err != nil {
			return nil, err
		}
	}

	otp, err := helper.GenerateOTPCode(6)
	if err != nil {
		return nil, err
	}

	//save to user_forgot table
	m := &model.UserForgot{
		UserID:  user.ID,
		Ref:     ref,
		Otp:     otp,
		Expires: time.Now().Unix() + 3600, // 1 hour expiration
		Used:    false,
	}
	_, err = s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, err
	}
	resp.Ref = ref

	resp.Email = helper.MaskEmail(user.Email)

	//send email with otp using HTML template
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Forgot Password - OTP Code</title>
		<style>
			body {
				margin: 0;
				padding: 0;
				font-family: Arial, sans-serif;
				background-color: #f5f5f5;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				background-color: white;
				border-radius: 10px;
				overflow: hidden;
				box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
			}
			.header {
				background: linear-gradient(135deg, #dc3545, #e74c3c);
				color: white;
				padding: 30px;
				text-align: center;
			}
			.header h1 {
				margin: 0;
				font-size: 28px;
				font-weight: bold;
			}
			.content {
				padding: 40px 30px;
				background-color: white;
			}
			.otp-container {
				background-color: #f8f9fa;
				border: 2px dashed #dc3545;
				border-radius: 8px;
				padding: 20px;
				text-align: center;
				margin: 20px 0;
			}
			.otp-code {
				font-size: 32px;
				font-weight: bold;
				color: #dc3545;
				letter-spacing: 4px;
				margin: 10px 0;
			}
			.ref-code {
				font-size: 14px;
				color: #6c757d;
				margin-top: 10px;
			}
			.info-text {
				color: #495057;
				line-height: 1.6;
				margin: 20px 0;
			}
			.warning {
				background-color: #fff5f5;
				border-left: 4px solid #dc3545;
				padding: 15px;
				margin: 20px 0;
				border-radius: 4px;
			}
			.warning p {
				margin: 0;
				color: #721c24;
			}
			.footer {
				background-color: #dc3545;
				color: white;
				padding: 20px;
				text-align: center;
				font-size: 14px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>🔐 Password Reset Request</h1>
			</div>
			<div class="content">
				<p class="info-text">Hello,</p>
				<p class="info-text">We received a request to reset your password. Please use the OTP code below to complete the password reset process:</p>
				
				<div class="otp-container">
					<div class="otp-code">` + otp + `</div>
					<div class="ref-code">Reference: ` + ref + `</div>
				</div>
				
				<div class="warning">
					<p><strong>⚠️ Important:</strong></p>
					<p>• This OTP code will expire in 1 hour</p>
					<p>• Do not share this code with anyone</p>
					<p>• If you didn't request this, please ignore this email</p>
				</div>
				
				<p class="info-text">If you're having trouble with the password reset process, please contact our support team.</p>
			</div>
			<div class="footer">
				<p>© 2025 HRM System. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>`

	err = config.SendEmail(req.Email, "HRM System", "Password Reset - OTP Code", htmlBody)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) ResetPassword(ctx context.Context, req *authdto.ResetPasswordRequest) error {
	m := &model.UserForgot{
		Otp: req.Otp,
		Ref: req.Ref,
	}

	err := s.db.NewSelect().Model(m).
		Where("otp = ?", req.Otp).
		Where("ref = ?", req.Ref).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(message.OTPCodeInvalid)
		}
		return err
	}
	if m.Used {
		return errors.New(message.OTPCodeUsed)
	}

	if m.Expires < time.Now().Unix() {
		return errors.New(message.OTPCodeExpired)
	}
	// Delete the OTP
	_, err = s.db.NewDelete().Model(m).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	changePassword := &userdto.ChangePasswordRequest{
		UserID:   m.UserID,
		Password: req.NewPassword,
	}
	// Update user password
	err = s.user.Svc.ChangePassword(ctx, changePassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Info(ctx context.Context, userID string) (*authdto.UserInfoResponse, error) {
	user, err := s.user.Svc.Get(ctx, userID)
	if err != nil {
		logger.Err(err.Error())
		return nil, err
	}

	//get permiossions for the user
	permissions, err := s.role.Svc.GetRolePermissions(ctx, user.RoleID)
	if err != nil {
		logger.Err(err.Error())
		return nil, err
	}

	//get shift schedules for the user
	shiftSchedules, err := s.workshift.Svc.Get(ctx, user.WorkShiftID)
	if err != nil && err != sql.ErrNoRows {
		logger.Err(err.Error())
		return nil, err
	}

	year := time.Now().Year()

	//get leaves for the user
	leaves, err := s.leave.Svc.LeaveByUser(ctx, user.ID, int64(year))
	if err != nil {
		logger.Err(err.Error())
		return nil, err
	}

	attendanceToday, err := s.attendance.Svc.AttendanceToday(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		logger.Err(err.Error())
		return nil, err
	}

	resp := &authdto.UserInfoResponse{
		ListUserResponse: *user,
		Permissions:      permissions,
		Leaves:           leaves,
	}
	if shiftSchedules != nil {
		resp.ShiftSchedules = shiftSchedules.Schedules
	} else {
		resp.ShiftSchedules = []workshiftdto.ShiftSchedule{}
	}

	if attendanceToday != nil {
		resp.AttendanceToday = *attendanceToday
	}

	//get holiday for user
	holidays, err := s.holiday.Svc.GetHolidaysByOrganization(ctx, user.OrganizationID)
	if err != nil {
		logger.Err(err.Error())
		return nil, err
	}
	resp.Holidays = holidays

	notificationCount, err := s.notification.Svc.GetUnreadCount(ctx, user.ID)
	if err != nil {
		logger.Err(err.Error())
		return nil, err
	}
	resp.NotificationCount = notificationCount

	return resp, nil
}
