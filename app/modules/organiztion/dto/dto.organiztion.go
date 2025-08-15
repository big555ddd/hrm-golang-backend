package organizationdto

type CreateOrganization struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateOrganization struct {
	CreateOrganization
}

type ListOrganizationRequest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"searchBy"`
	SortBy   string `form:"sortBy"`
	OrderBy  string `form:"orderBy"`
}

type GetByIDOrganization struct {
	ID string `uri:"id" binding:"required"`
}

type ListOrganizationResponse struct {
	ID          string `bun:"id" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	CreatedAt   int64  `bun:"created_at" json:"created_at"`
	UpdatedAt   int64  `bun:"updated_at" json:"updated_at"`
}
