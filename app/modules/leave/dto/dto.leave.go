package leavedto

type CreateLeave struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Year            int      `json:"year"`
	Amount          int64    `json:"amount"`
	OrganizationIDs []string `json:"organizationIds"`
}

type UpdateLeave struct {
	CreateLeave
}

type ListLeaveRequest struct {
	Page           int    `form:"page"`
	Size           int    `form:"size"`
	Search         string `form:"search"`
	SearchBy       string `form:"searchBy"`
	SortBy         string `form:"sortBy"`
	OrderBy        string `form:"orderBy"`
	OrganizationID string `form:"organizationId"`
	Year           int    `form:"year"`
}

type GetByIDLeave struct {
	ID string `uri:"id" binding:"required"`
}

type ListLeaveResponse struct {
	ID            string                      `bun:"id" json:"id"`
	Name          string                      `bun:"name" json:"name"`
	Description   string                      `bun:"description" json:"description"`
	Year          int                         `json:"year"`
	Amount        int64                       `json:"amount"`
	Organizations []OrganizationLeaveResponse `bun:"organizations" json:"organizations,omitempty"`
	CreatedAt     int64                       `bun:"created_at" json:"created_at"`
	UpdatedAt     int64                       `bun:"updated_at" json:"updated_at"`
}

type LeaveUserResponse struct {
	ID         string  `bun:"id" json:"id"`
	Name       string  `bun:"name" json:"name"`
	Amount     int64   `bun:"amount" json:"amount"`
	UsedAmount float64 `bun:"used_amount" json:"used_amount"`
}

type OrganizationLeaveResponse struct {
	ID   string `bun:"id" json:"id"`
	Name string `bun:"name" json:"name"`
}

type SetLeavePermissions struct {
	LeaveID       string   `json:"LeaveId" binding:"required"`
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}

type LeavePermissionResponse struct {
	PermissionID   string `bun:"permission_id" json:"permission_id"`
	PermissionName string `bun:"permission_name" json:"permission_name"`
}
