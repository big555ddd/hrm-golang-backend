package userdto

import (
	"fmt"
	"strings"
)

type CreateUser struct {
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	RoleID       string  `json:"roleId"`
	Email        string  `json:"email" binding:"required,email"`
	Password     string  `json:"password" binding:"min=8"`
	Phone        string  `json:"phone"`
	Salary       float64 `json:"salary"`
	DepartmentID string  `json:"departmentId"`
	WorkShiftID  string  `json:"workShiftId"`
}

type UpdateUser struct {
	CreateUser
}

type ListUserRequest struct {
	Page           int    `form:"page"`
	Size           int    `form:"size"`
	Search         string `form:"search"`
	SearchBy       string `form:"searchBy"`
	SortBy         string `form:"sortBy"`
	OrderBy        string `form:"orderBy"`
	OrganizationID string `form:"organizationId"`
	BranchID       string `form:"branchId"`
	DepartmentID   string `form:"departmentId"`
	WorkShiftID    string `form:"workShiftId"`
	RoleID         string `form:"roleId"`
	ForNoti        string
	WithPermission string
}

type GetByIDUser struct {
	ID string `uri:"id" binding:"required"`
}
type ListUserResponse struct {
	ID               string  `bun:"id" json:"id"`
	FirstName        string  `bun:"first_name" json:"first_name"`
	LastName         string  `bun:"last_name" json:"last_name"`
	Email            string  `bun:"email" json:"email"`
	EmpCode          string  `bun:"emp_code" json:"emp_code"`
	Phone            string  `bun:"phone" json:"phone"`
	IsActive         bool    `bun:"is_active" json:"is_active"`
	RoleID           string  `bun:"role_id" json:"role_id"`
	RoleName         string  `bun:"role_name" json:"role_name"`
	DepartmentID     string  `bun:"department_id" json:"department_id"`
	DepartmentName   string  `bun:"department_name" json:"department_name"`
	BranchID         string  `bun:"branch_id" json:"branch_id"`
	BranchName       string  `bun:"branch_name" json:"branch_name"`
	OrganizationID   string  `bun:"organization_id" json:"organization_id"`
	OrganizationName string  `bun:"organization_name" json:"organization_name"`
	WorkShiftID      string  `bun:"work_shift_id" json:"work_shift_id"`
	WorkShiftName    string  `bun:"work_shift_name" json:"work_shift_name"`
	Salary           float64 `bun:"salary" json:"salary"`
	CreatedAt        int64   `bun:"created_at" json:"created_at"`
	UpdatedAt        int64   `bun:"updated_at" json:"updated_at"`
}

type ChangePasswordRequest struct {
	UserID   string `json:"userId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDayoffResponse struct {
	DayOffWeek []int64 `json:"day_off_week"`
	Holidays   []int64 `json:"holidays"`
}

func (r *ListUserRequest) Validator() error {
	if r.Page == 0 {
		r.Page = 1
	}

	if r.Size == 0 {
		r.Size = 10
	}

	if r.OrderBy == "" {
		r.OrderBy = "asc"
	}

	if r.SortBy == "" {
		r.SortBy = "created_at"
	} else {
		sortForm := []string{"created_at", "first_name"}
		if !strings.Contains(strings.Join(sortForm, ","), r.SortBy) {
			return fmt.Errorf("invalid sort by field: %s", r.SortBy)
		}
	}

	if r.SearchBy != "" {
		// Define valid fields and their corresponding prefixes
		searchByMap := map[string]string{
			"id":                "u.",
			"first_name":        "u.",
			"last_name":         "u.",
			"email":             "u.",
			"emp_code":          "u.",
			"phone":             "u.",
			"salary":            "u.",
			"created_at":        "u.",
			"updated_at":        "u.",
			"role_name":         "r.",
			"department_name":   "d.",
			"branch_name":       "b.",
			"organization_name": "o.",
			"work_shift_name":   "w.",
		}
		prefix, ok := searchByMap[r.SearchBy]
		if !ok {
			return fmt.Errorf("invalid search by field: %s", r.SearchBy)
		}
		r.SearchBy = prefix + r.SearchBy
		return nil
	}
	return nil
}
