package user

import (
	message "app/app/messsage"
	"app/app/model"
	"app/app/modules/department"
	"app/app/modules/role"
	userdto "app/app/modules/user/dto"
	"app/app/modules/workshift"
	"app/internal/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db         *bun.DB
	department *department.Module
	role       *role.Module
	workshift  *workshift.Module
}

func NewService(db *bun.DB, department *department.Module, role *role.Module, workshift *workshift.Module) *Service {
	return &Service{
		db:         db,
		department: department,
		role:       role,
		workshift:  workshift,
	}
}

func (s *Service) Create(ctx context.Context, req userdto.CreateUser) (*model.User, bool, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, false, err
	}
	empCode, err := s.GetNextEmpCode(ctx)
	if err != nil {
		return nil, false, err
	}

	//cheeck department exists
	department, err := s.department.Svc.Exist(ctx, req.DepartmentID)
	if err != nil {
		return nil, false, err
	}
	if !department {
		return nil, false, errors.New(message.DepartmentNotFound)
	}

	//check role exists
	Existrole, err := s.role.Svc.Exist(ctx, req.RoleID)
	if err != nil {
		return nil, false, err
	}
	if !Existrole {
		return nil, true, errors.New(message.RoleNotFound)
	}

	m := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(bytes),
		Phone:     req.Phone,
		EmpCode:   empCode,
		Salary:    req.Salary,
		IsActive:  true,
	}
	mserr := false
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewInsert().Model(m).Exec(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				mserr = true
				return errors.New(message.EmailAlreadyExists)
			}
		}
		//create userDepartment
		userDepartment := &model.UserDepartment{
			UserID:       m.ID,
			DepartmentID: req.DepartmentID,
		}
		_, err = tx.NewInsert().Model(userDepartment).Exec(ctx)
		if err != nil {
			return err
		}

		//create userRole
		userRole := &model.UserRole{
			UserID: m.ID,
			RoleID: req.RoleID,
		}
		_, err = tx.NewInsert().Model(userRole).Exec(ctx)
		if err != nil {
			return err
		}

		//create userWorkShift
		if req.WorkShiftID != "" {
			err = s.workshift.Svc.SetWorkShift(ctx, &tx, req.WorkShiftID, m.ID, true)
			if err != nil {
				logger.Err(err)
				return err
			}
		}
		return nil
	})

	return m, mserr, err

}

func (s *Service) GetNextEmpCode(ctx context.Context) (string, error) {
	var lastUser model.User
	err := s.db.NewSelect().Table("users").Order("created_at DESC").Limit(1).Scan(ctx, &lastUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return "EMP-0001", nil
		}
		return "", err
	}

	if lastUser.EmpCode == "" {
		return "EMP-0001", nil
	}

	lastEmpCode := lastUser.EmpCode
	nextEmpCode := fmt.Sprintf("EMP-%04d", (atoi(lastEmpCode[4:]) + 1))
	return nextEmpCode, nil
}

func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}

func (s *Service) Update(ctx context.Context, req userdto.UpdateUser, id userdto.GetByIDUser) (*model.User, bool, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, false, err
	}
	// ex, err := s.Exist(ctx, id.ID)
	// if err != nil {
	// 	return nil, false, err
	// }

	// if !ex {
	// 	return nil, false, err
	// }
	m := &model.User{
		ID:        id.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(bytes),
		Phone:     req.Phone,
		Salary:    req.Salary,
	}

	ex, err := m.Exist(ctx, s.db, m.ID)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New(message.UserNotFound)
	}

	//cheeck department exists
	department, err := s.department.Svc.Exist(ctx, req.DepartmentID)
	if err != nil {
		return nil, false, err
	}
	if !department {
		return nil, true, errors.New(message.DepartmentNotFound)
	}

	//check role exists
	Existrole, err := s.role.Svc.Exist(ctx, req.RoleID)
	if err != nil {
		return nil, false, err
	}
	if !Existrole {
		return nil, false, errors.New(message.RoleNotFound)
	}

	m.SetUpdateNow()
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		update := tx.NewUpdate().Model(m).
			Set("first_name = ?", m.FirstName).
			Set("last_name = ?", m.LastName).
			Set("email = ?", m.Email)
		if req.Password != "" {
			update.Set("password = ?", m.Password)
		}
		_, err := update.Set("phone = ?", m.Phone).
			Set("salary = ?", m.Salary).
			Set("updated_at = ?", m.UpdatedAt).
			WherePK().
			OmitZero().
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}
		//delete old userDepartment
		_, err = tx.NewDelete().Model((*model.UserDepartment)(nil)).Where("user_id = ?", m.ID).Exec(ctx)
		if err != nil {
			return err
		}
		//create userDepartment
		userDepartment := &model.UserDepartment{
			UserID:       m.ID,
			DepartmentID: req.DepartmentID,
		}
		_, err = tx.NewInsert().Model(userDepartment).Exec(ctx)
		if err != nil {
			return err
		}
		//delete old userRole
		_, err = tx.NewDelete().Model((*model.UserRole)(nil)).Where("user_id = ?", m.ID).Exec(ctx)
		if err != nil {
			return err
		}
		//create userRole
		userRole := &model.UserRole{
			UserID: m.ID,
			RoleID: req.RoleID,
		}
		_, err = tx.NewInsert().Model(userRole).Exec(ctx)
		if err != nil {
			return err
		}
		//delete old userWorkShift
		_, err = tx.NewDelete().Model((*model.UserWorkShift)(nil)).Where("user_id = ?", m.ID).Exec(ctx)
		if err != nil {
			return err
		}
		//create userWorkShift
		if req.WorkShiftID != "" {
			err = s.workshift.Svc.SetWorkShift(ctx, &tx, req.WorkShiftID, m.ID, false)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return m, false, err
}

func (s *Service) List(ctx context.Context, req userdto.ListUserRequest) ([]userdto.ListUserResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []userdto.ListUserResponse{}
	query := s.db.NewSelect().
		TableExpr("users AS u").
		Column("u.id", "u.first_name", "u.last_name",
			"u.email", "u.emp_code", "u.phone", "u.is_active",
			"u.salary", "u.created_at", "u.updated_at").
		ColumnExpr("r.id AS role_id").
		ColumnExpr("r.name AS role_name").
		ColumnExpr("d.id AS department_id").
		ColumnExpr("d.name AS department_name").
		ColumnExpr("b.id AS branch_id").
		ColumnExpr("b.name AS branch_name").
		ColumnExpr("o.id AS organization_id").
		ColumnExpr("o.name AS organization_name").
		ColumnExpr("uw.work_shift_id AS work_shift_id").
		ColumnExpr("ws.name AS work_shift_name").
		Join("LEFT JOIN user_roles AS ur").
		JoinOn("u.id = ur.user_id").
		Join("LEFT JOIN roles AS r").
		JoinOn("ur.role_id = r.id").
		Join("LEFT JOIN user_departments AS ud").
		JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS d").
		JoinOn("ud.department_id = d.id").
		Join("LEFT JOIN branches AS b").
		JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id").
		Join("LEFT JOIN user_work_shifts AS uw").
		JoinOn("u.id = uw.user_id").
		Join("LEFT JOIN work_shifts AS ws").
		JoinOn("uw.work_shift_id = ws.id").
		Where("u.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(u.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(u.first_name) LIKE ?", search)
		}
	}

	if req.OrganizationID != "" {
		query.Where("o.id = ?", req.OrganizationID)
	}

	if req.BranchID != "" {
		query.Where("b.id = ?", req.BranchID)
	}

	if req.DepartmentID != "" {
		query.Where("d.id = ?", req.DepartmentID)
	}

	if req.WorkShiftID != "" {
		query.Where("uw.work_shift_id = ?", req.WorkShiftID)
	}

	if req.RoleID != "" {
		query.Where("ur.role_id = ?", req.RoleID)
	}

	if req.ForNoti != "" {
		query.Where("u.id != ?", req.ForNoti)
	}

	if req.WithPermission != "" {
		query.Join("LEFT JOIN role_permissions AS rp").
			JoinOn("rp.role_id = r.id").
			Join("LEFT JOIN permissions AS p ON p.id = rp.permission_id").
			Where("p.name = ?", req.WithPermission)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*userdto.ListUserResponse, error) {
	m := userdto.ListUserResponse{}
	err := s.db.NewSelect().
		TableExpr("users AS u").
		Column("u.id", "u.first_name", "u.last_name",
			"u.email", "u.emp_code", "u.phone", "u.is_active",
			"u.salary", "u.created_at", "u.updated_at").
		ColumnExpr("r.id AS role_id").
		ColumnExpr("r.name AS role_name").
		ColumnExpr("d.id AS department_id").
		ColumnExpr("d.name AS department_name").
		ColumnExpr("b.id AS branch_id").
		ColumnExpr("b.name AS branch_name").
		ColumnExpr("o.id AS organization_id").
		ColumnExpr("o.name AS organization_name").
		ColumnExpr("uw.work_shift_id AS work_shift_id").
		ColumnExpr("ws.name AS work_shift_name").
		Join("LEFT JOIN user_roles AS ur").
		JoinOn("u.id = ur.user_id").
		Join("LEFT JOIN roles AS r").
		JoinOn("ur.role_id = r.id").
		Join("LEFT JOIN user_departments AS ud").
		JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS d").
		JoinOn("ud.department_id = d.id").
		Join("LEFT JOIN branches AS b").
		JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id").
		Join("LEFT JOIN user_work_shifts AS uw").
		JoinOn("u.id = uw.user_id").
		Join("LEFT JOIN work_shifts AS ws").
		JoinOn("uw.work_shift_id = ws.id").
		Where("u.deleted_at IS NULL").
		Where("u.id = ?", id).Scan(ctx, &m)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(message.UserNotFound)
		}
		return nil, err
	}
	return &m, nil
}

func (s *Service) Delete(ctx context.Context, id userdto.GetByIDUser) error {
	// ex, err := s.Exist(ctx, id.ID)
	// if err != nil {
	// 	return err
	// }

	// if !ex {
	// 	return errors.New(message.UserNotFound)
	// }

	_, err := s.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err

}

// func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
// 	ex, err := s.db.NewSelect().Model(&model.User{}).Where("id = ?", id).Exists(ctx)
// 	return ex, err
// }

func (s *Service) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	m := model.User{}
	err := s.db.NewSelect().Model(&m).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if m.ID == "" {
		return nil, errors.New("user not found")
	}
	return &m, nil
}

func (s *Service) GetByEmpCode(ctx context.Context, empCode string) (*model.User, error) {
	m := model.User{}
	err := s.db.NewSelect().Model(&m).Where("emp_code = ?", empCode).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if m.ID == "" {
		return nil, errors.New("user not found")
	}
	return &m, nil
}

func (s *Service) CheckPassword(ctx context.Context, req *model.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangePassword(ctx context.Context, req *userdto.ChangePasswordRequest) error {
	m := &model.User{
		ID:       req.UserID,
		Password: req.Password,
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), 14)
	if err != nil {
		return err
	}
	m.Password = string(bytes)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("password = ?", m.Password).
		Set("updated_at = ?", m.UpdatedAt).
		WherePK().
		OmitZero().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserPermissionsName(ctx context.Context, userID string) ([]string, error) {
	user, err := s.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	permission, err := s.role.Svc.GetRolePermissionsName(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

// func for go not working day and holiday for this user
// func (s *Service) GetUserHolidays(ctx context.Context, userID string) (*userdto.UserDayoffResponse, error) {
// 	resp := &userdto.UserDayoffResponse{}
// 	user, err := s.Get(ctx, userID)
// 	if err != nil {
// 		return resp, err
// 	}
// 	resp.DayOffWeek = user.DayOffWeek
// 	resp.Holidays = user.Holidays
// 	return resp, nil
// }
