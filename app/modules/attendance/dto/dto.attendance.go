package attendancedto

type CreateAttendance struct {
	UserID      string `json:"user_id"`
	WorkShiftID string `json:"work_shift_id"`
	CheckIn     int64  `json:"check_in"`
	CheckOut    int64  `json:"check_out"`
	Date        int64  `json:"date"`
	IsOnTime    bool   `json:"is_on_time"`
	IsLate      bool   `json:"is_late"`
	IsLeave     bool   `json:"is_leave"`
}

type UpdateAttendance struct {
	CreateAttendance
}

type CheckInAttendanceRequest struct {
	CheckInTime int64   `json:"checkInTime" binding:"required"`
	Date        int64   `json:"date" binding:"required"`
	LocationX   float64 `json:"locationX"`
	LocationY   float64 `json:"locationY"`
}

type ListAttendanceRequest struct {
	Page           int    `form:"page"`
	Size           int    `form:"size"`
	Search         string `form:"search"`
	SearchBy       string `form:"searchBy"`
	SortBy         string `form:"sortBy"`
	OrderBy        string `form:"orderBy"`
	UserID         string `form:"userId"`
	WorkShiftID    string `form:"workShiftId"`
	OrganizationID string `form:"organizationId"`
	BranchID       string `form:"branchId"`
	DepartmentID   string `form:"departmentId"`
	Date           int64  `form:"date"`
	IsOnTime       bool   `form:"isOnTime"`
	IsLate         bool   `form:"isLate"`
	IsLeave        bool   `form:"isLeave"`
}

type GetByIDAttendance struct {
	ID string `uri:"id" binding:"required"`
}

type ListAttendanceResponse struct {
	ID            string `bun:"id" json:"id"`
	UserID        string `json:"user_id"`
	EmpCode       string `bun:"emp_code" json:"emp_code"`
	FirstName     string `bun:"first_name" json:"first_name"`
	LastName      string `bun:"last_name" json:"last_name"`
	WorkShiftID   string `json:"work_shift_id"`
	WorkShiftName string `bun:"work_shift_name" json:"work_shift_name"`
	CheckIn       int64  `json:"check_in"`
	CheckOut      int64  `json:"check_out"`
	Date          int64  `json:"date"`
	IsOnTime      bool   `json:"is_on_time"`
	IsLate        bool   `json:"is_late"`
	IsLeave       bool   `json:"is_leave"`
	CreatedAt     int64  `bun:"created_at" json:"created_at"`
	UpdatedAt     int64  `bun:"updated_at" json:"updated_at"`
}

type SetAttendancePermissions struct {
	AttendanceID  string   `json:"AttendanceId" binding:"required"`
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}

type AttendancePermissionResponse struct {
	PermissionID   string `bun:"permission_id" json:"permission_id"`
	PermissionName string `bun:"permission_name" json:"permission_name"`
}

type AttendanceCountResponse struct {
	CheckInCount int `json:"count"`
	LateCount    int `json:"late_count"`
	AbsentCount  int `json:"absent_count"`
}

type AttendanceCountRequest struct {
	OrganizationID string `form:"organizationId"`
	BranchID       string `form:"branchId"`
	DepartmentID   string `form:"departmentId"`
	UserID         string `form:"userId"`
	Date           int64  `form:"date"`
}
