package roledto

type CreateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRole struct {
	CreateRole
}

type ListRoleRequest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"searchBy"`
	SortBy   string `form:"sortBy"`
	OrderBy  string `form:"orderBy"`
}

type GetByIDRole struct {
	ID string `uri:"id" binding:"required"`
}

type ListRoleResponse struct {
	ID          string `bun:"id" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	CreatedAt   int64  `bun:"created_at" json:"created_at"`
	UpdatedAt   int64  `bun:"updated_at" json:"updated_at"`
}

type SetRolePermissions struct {
	RoleID        string   `json:"roleId" binding:"required"`
	PermissionIDs []string `json:"permissionIds" binding:"required"`
}

type RolePermissionResponse struct {
	PermissionID   string `bun:"permission_id" json:"permission_id"`
	PermissionName string `bun:"permission_name" json:"permission_name"`
}
