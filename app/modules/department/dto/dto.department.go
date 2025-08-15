package departmentdto

type CreateDepartment struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BranchID    string `json:"branchId"`
}

type UpdateDepartment struct {
	CreateDepartment
}

type ListDepartmentRequest struct {
	Page           int    `form:"page"`
	Size           int    `form:"size"`
	Search         string `form:"search"`
	SearchBy       string `form:"searchBy"`
	SortBy         string `form:"sortBy"`
	OrderBy        string `form:"orderBy"`
	BranchID       string `form:"branchId"`
	OrganizationID string `form:"organizationId"`
}

type GetByIDDepartment struct {
	ID string `uri:"id" binding:"required"`
}

type ListDepartmentResponse struct {
	ID               string `bun:"id" json:"id"`
	Name             string `bun:"name" json:"name"`
	Description      string `bun:"description" json:"description"`
	BranchID         string `bun:"branch_id" json:"branch_id"`
	BranchName       string `bun:"branch_name" json:"branch_name"`
	OrganizationID   string `bun:"organization_id" json:"organization_id"`
	OrganizationName string `bun:"organization_name" json:"organization_name"`
	CreatedAt        int64  `bun:"created_at" json:"created_at"`
	UpdatedAt        int64  `bun:"updated_at" json:"updated_at"`
}
