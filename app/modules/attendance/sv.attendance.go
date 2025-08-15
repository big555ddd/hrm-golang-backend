package attendance

import (
	"app/app/helper"
	message "app/app/messsage"
	"app/app/model"
	attendancedto "app/app/modules/attendance/dto"
	"app/app/modules/user"
	"app/app/modules/workshift"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	db        *bun.DB
	user      *user.Module
	workshift *workshift.Module
}

func NewService(db *bun.DB, user *user.Module, workshift *workshift.Module) *Service {
	return &Service{
		db:        db,
		user:      user,
		workshift: workshift,
	}

}

func (s *Service) Create(ctx context.Context, req *attendancedto.CreateAttendance) (*model.Attendance, error) {
	m := &model.Attendance{
		UserID:      req.UserID,
		WorkShiftID: req.WorkShiftID,
		CheckIn:     req.CheckIn,
		CheckOut:    req.CheckOut,
		Date:        req.Date,
		IsOnTime:    req.IsOnTime,
		IsLate:      req.IsLate,
		IsLeave:     req.IsLeave,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return m, err

}

func (s *Service) CheckIn(ctx context.Context, userID string, req *attendancedto.CheckInAttendanceRequest) (*model.Attendance, error) {
	user, err := s.user.Svc.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	//check if user has already checked in today
	existingAttendance, err := s.db.NewSelect().Model(&model.Attendance{}).
		Where("user_id = ? AND date = ?", userID, req.Date).
		Limit(1).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if existingAttendance {
		m := &model.Attendance{}
		_, err = s.db.NewUpdate().Model(m).
			Set("check_out = ?", req.CheckInTime).
			Where("user_id = ? AND date = ?", userID, req.Date).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return m, nil
	}
	m := &model.Attendance{
		UserID:      userID,
		WorkShiftID: user.WorkShiftID,
		CheckIn:     req.CheckInTime,
		Date:        req.Date,
		IsOnTime:    true,
		IsLate:      false,
		IsLeave:     false,
	}
	//check late for user
	workshift, err := s.workshift.Svc.Get(ctx, user.WorkShiftID)
	if err != nil {
		return nil, err
	}

	if workshift.WorkLocationX != 0 && workshift.WorkLocationY != 0 {
		allowed := helper.CalculateDistance(req.LocationX, req.LocationY, workshift.WorkLocationX, workshift.WorkLocationY)
		if !allowed {
			return nil, errors.New(message.CheckInLocationNotAllowed)
		}
	}

	day := helper.UnixToDay(req.Date)
	schedule := helper.GetScheduleForDay(workshift.Schedules, day)

	// If no schedule for this day (holidays/OT), set as on-time by default
	if schedule == nil {
		m.IsOnTime = true
		m.IsLate = false
	} else {
		// Normal working day - check if late
		checkInTime := time.Unix(req.CheckInTime, 0).In(time.FixedZone("Asia/Bangkok", 7*3600))
		LateTime := float64(schedule.StartTime) + float64(workshift.LateMinutes)/60.0 // Convert minutes to hours
		CheckIn := float64(checkInTime.Hour()) + float64(checkInTime.Minute())/60.0
		if CheckIn > LateTime {
			m.IsOnTime = false
			m.IsLate = true
		}
	}

	m.CheckIn = req.CheckInTime
	m.CheckOut = 0
	_, err = s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Service) AttendanceToday(ctx context.Context, userID string) (*model.Attendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	date := today.Unix() - 25200
	m := &model.Attendance{}
	err := s.db.NewSelect().Model(m).
		Where("user_id = ? AND date = ?", userID, date).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No attendance record for today
		}
		return nil, err
	}

	return m, nil
}

func (s *Service) AttendanceCount(ctx context.Context, req *attendancedto.AttendanceCountRequest) (*attendancedto.AttendanceCountResponse, error) {

	// Convert date to day of week for schedule checking
	day := helper.UnixToDay(req.Date)

	// Build query for check-in count
	checkInQuery := s.db.NewSelect().TableExpr("attendances AS a").
		Join("LEFT JOIN users AS u").JoinOn("u.id = a.user_id").
		Join("LEFT JOIN work_shifts AS ws").JoinOn("ws.id = a.work_shift_id").
		Join("LEFT JOIN user_departments AS ud").JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS d").JoinOn("ud.department_id = d.id").
		Join("LEFT JOIN branches AS b").JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").JoinOn("b.organization_id = o.id").
		Where("a.date = ?", req.Date)

	// Build query for late count
	lateQuery := s.db.NewSelect().TableExpr("attendances AS a").
		Join("LEFT JOIN users AS u").JoinOn("u.id = a.user_id").
		Join("LEFT JOIN work_shifts AS ws").JoinOn("ws.id = a.work_shift_id").
		Join("LEFT JOIN user_departments AS ud").JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS d").JoinOn("ud.department_id = d.id").
		Join("LEFT JOIN branches AS b").JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").JoinOn("b.organization_id = o.id").
		Where("a.date = ? AND a.is_late = true", req.Date)

	// Build query for user count - only count users who should work on this day
	userQuery := s.db.NewSelect().TableExpr("users AS u").
		Join("LEFT JOIN user_work_shifts AS uws").JoinOn("u.id = uws.user_id").
		Join("LEFT JOIN work_shifts AS ws").JoinOn("uws.work_shift_id = ws.id").
		Join("LEFT JOIN shift_schedules AS ss").JoinOn("ws.id = ss.work_shift_id").
		Join("LEFT JOIN user_departments AS ud").JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS d").JoinOn("ud.department_id = d.id").
		Join("LEFT JOIN branches AS b").JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").JoinOn("b.organization_id = o.id").
		Where("u.is_active = ? AND ss.day = ?", true, day)

	// Apply filters if provided
	if req.OrganizationID != "" {
		checkInQuery.Where("o.id = ?", req.OrganizationID)
		lateQuery.Where("o.id = ?", req.OrganizationID)
		userQuery.Where("o.id = ?", req.OrganizationID)
	}

	if req.BranchID != "" {
		checkInQuery.Where("b.id = ?", req.BranchID)
		lateQuery.Where("b.id = ?", req.BranchID)
		userQuery.Where("b.id = ?", req.BranchID)
	}

	if req.DepartmentID != "" {
		checkInQuery.Where("d.id = ?", req.DepartmentID)
		lateQuery.Where("d.id = ?", req.DepartmentID)
		userQuery.Where("d.id = ?", req.DepartmentID)
	}

	if req.UserID != "" {
		checkInQuery.Where("u.id = ?", req.UserID)
		lateQuery.Where("u.id = ?", req.UserID)
		userQuery.Where("u.id = ?", req.UserID)
	}

	// Execute queries
	checkInCount, err := checkInQuery.Count(ctx)
	if err != nil {
		return nil, err
	}

	lateCount, err := lateQuery.Count(ctx)
	if err != nil {
		return nil, err
	}

	userCount, err := userQuery.Count(ctx)
	if err != nil {
		return nil, err
	}

	resp := &attendancedto.AttendanceCountResponse{
		CheckInCount: checkInCount,
		LateCount:    lateCount,
		AbsentCount:  userCount - checkInCount,
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, req *attendancedto.UpdateAttendance, id string) (*model.Attendance, bool, error) {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New(message.RoleNotFound)
	}

	m := &model.Attendance{
		ID:          id,
		UserID:      req.UserID,
		WorkShiftID: req.WorkShiftID,
		CheckIn:     req.CheckIn,
		CheckOut:    req.CheckOut,
		Date:        req.Date,
		IsOnTime:    req.IsOnTime,
		IsLate:      req.IsLate,
		IsLeave:     req.IsLeave,
	}
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("user_id = ?", m.UserID).
		Set("work_shift_id = ?", m.WorkShiftID).
		Set("check_in = ?", m.CheckIn).
		Set("check_out = ?", m.CheckOut).
		Set("date = ?", m.Date).
		Set("is_on_time = ?", m.IsOnTime).
		Set("is_late = ?", m.IsLate).
		Set("is_leave = ?", m.IsLeave).
		Set("updated_at = ?", m.UpdatedAt).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, false, err
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req *attendancedto.ListAttendanceRequest) ([]attendancedto.ListAttendanceResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []attendancedto.ListAttendanceResponse{}
	query := s.db.NewSelect().
		TableExpr("attendances AS a").
		Column("a.id", "a.user_id", "a.work_shift_id",
			"a.check_in", "a.check_out", "a.date",
			"a.is_on_time", "a.is_late", "a.is_leave",
			"a.created_at", "a.updated_at").
		ColumnExpr("u.first_name").
		ColumnExpr("u.last_name").
		ColumnExpr("u.emp_code").
		ColumnExpr("ws.name AS work_shift_name").
		Join("LEFT JOIN users AS u").
		JoinOn("u.id = a.user_id").
		Join("LEFT JOIN work_shifts AS ws").
		JoinOn("ws.id = a.work_shift_id").
		Join("LEFT JOIN user_departments AS ud").
		JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS d").
		JoinOn("ud.department_id = d.id").
		Join("LEFT JOIN branches AS b").
		JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(a.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(a.user_id) LIKE ?", search)
		}
	}

	if req.UserID != "" {
		query.Where("a.user_id = ?", req.UserID)
	}

	if req.WorkShiftID != "" {
		query.Where("a.work_shift_id = ?", req.WorkShiftID)
	}

	if req.Date != 0 {
		query.Where("a.date = ?", req.Date)
	}

	if req.IsOnTime {
		query.Where("a.is_on_time = ?", req.IsOnTime)
	}

	if req.IsLate {
		query.Where("a.is_late = ?", req.IsLate)
	}

	if req.IsLeave {
		query.Where("a.is_leave = ?", req.IsLeave)
	}

	if req.OrganizationID != "" {
		query.Where("b.organization_id = ?", req.OrganizationID)
	}

	if req.BranchID != "" {
		query.Where("d.branch_id = ?", req.BranchID)
	}
	if req.DepartmentID != "" {
		query.Where("ud.department_id = ?", req.DepartmentID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("a.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*attendancedto.ListAttendanceResponse, error) {
	m := attendancedto.ListAttendanceResponse{}
	err := s.db.NewSelect().
		TableExpr("attendances AS a").
		Column("a.id", "a.user_id", "a.work_shift_id",
			"a.check_in", "a.check_out", "a.date",
			"a.is_on_time", "a.is_late", "a.is_leave",
			"a.created_at", "a.updated_at").
		Where("a.id = ?", id).Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New(message.RoleNotFound)
	}
	_, err = s.db.NewDelete().Model(&model.Attendance{}).Where("id = ?", id).Exec(ctx)
	return err

}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Attendance{}).Where("id = ?", id).Exists(ctx)
	return ex, err
}
