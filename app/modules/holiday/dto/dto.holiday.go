package holidaydto

type CreateHoliday struct {
	Name            string   `json:"name"`
	IsActive        bool     `json:"isActive"`
	Description     string   `json:"description"`
	StartDate       int64    `json:"startDate"`
	EndDate         int64    `json:"endDate"`
	OrganizationIDs []string `json:"organizationIds"`
}

type UpdateHoliday struct {
	CreateHoliday
}

type ListHolidayRequest struct {
	Page           int    `form:"page"`
	Size           int    `form:"size"`
	Search         string `form:"search"`
	SearchBy       string `form:"searchBy"`
	SortBy         string `form:"sortBy"`
	OrderBy        string `form:"orderBy"`
	OrganizationID string `form:"organizationId"`
	Year           int    `form:"year"`
}

type GetByIDHoliday struct {
	ID string `uri:"id" binding:"required"`
}
