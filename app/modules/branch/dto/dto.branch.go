package branchdto

type CreateBranch struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	OrganizationID string `json:"organizationId" binding:"required"`
}

type UpdateBranch struct {
	CreateBranch
}

type ListBranchRequest struct {
	Page           int    `form:"page"`
	Size           int    `form:"size"`
	Search         string `form:"search"`
	SearchBy       string `form:"searchBy"`
	SortBy         string `form:"sortBy"`
	OrderBy        string `form:"orderBy"`
	OrganizationID string `form:"organizationId"`
}

type GetByIDBranch struct {
	ID string `uri:"id" binding:"required"`
}

type ListBranchResponse struct {
	ID               string `bun:"id" json:"id"`
	Name             string `bun:"name" json:"name"`
	Description      string `bun:"description" json:"description"`
	OrganizationID   string `bun:"organization_id" json:"organization_id"`
	OrganizationName string `bun:"organization_name" json:"organization_name"`
	CreatedAt        int64  `bun:"created_at" json:"created_at"`
	UpdatedAt        int64  `bun:"updated_at" json:"updated_at"`
}
