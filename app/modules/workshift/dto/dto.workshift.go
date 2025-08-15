package workshiftdto

import "app/app/enum"

type CreateWorkShift struct {
	Name          string          `json:"name"`
	WorkLocationX float64         `json:"workLocationX"`
	WorkLocationY float64         `json:"workLocationY"`
	Description   string          `json:"description"`
	LateMinutes   int64           `json:"lateMinutes"`
	Schedules     []ShiftSchedule `json:"shiftSchedules"`
}

type ShiftSchedule struct {
	Day       enum.Day `json:"day"`
	StartTime int64    `json:"startTime"`
	EndTime   int64    `json:"endTime"`
}

type UpdateWorkShift struct {
	CreateWorkShift
}

type ListWorkShiftRequest struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"searchBy"`
	SortBy   string `form:"sortBy"`
	OrderBy  string `form:"orderBy"`
}

type GetByIDWorkShift struct {
	ID string `uri:"id" binding:"required"`
}

type ListWorkShiftResponse struct {
	ID            string          `bun:"id" json:"id"`
	Name          string          `bun:"name" json:"name"`
	WorkLocationX float64         `bun:"work_location_x" json:"workLocationX"`
	WorkLocationY float64         `bun:"work_location_y" json:"workLocationY"`
	LateMinutes   int64           `bun:"late_minutes" json:"late_minutes"`
	Description   string          `bun:"description" json:"description"`
	Schedules     []ShiftSchedule `bun:"shift_schedules" json:"shift_schedules,omitempty"`
	CreatedAt     int64           `bun:"created_at" json:"created_at"`
	UpdatedAt     int64           `bun:"updated_at" json:"updated_at"`
}

type ChangeWorkShiftRequest struct {
	UserID      string `json:"userId" binding:"required"`
	WorkShiftID string `json:"workShiftId" binding:"required"`
}
