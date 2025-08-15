package documentdto

import "app/app/enum"

type CreateDocument struct {
	UserID      string            `json:"userId" binding:"required"`
	Type        enum.DocumentType `json:"type" binding:"required"`
	Description string            `json:"description"`
	Leave       *LeaveReq         `json:"leave,omitempty" binding:"omitempty"`
	OverTime    *OverTimeReq      `json:"overtime,omitempty" binding:"omitempty"`
}

type LeaveReq struct {
	LeaveID   string `json:"leaveId" binding:"required"`
	StartDate int64  `json:"startDate" binding:"required"`
	EndDate   int64  `json:"endDate" binding:"required"`
}

type OverTimeReq struct {
	OverTimeType enum.OverTimeType `bun:"overtime_type" json:"overtimeType" binding:"required"`
	StartDate    int64             `json:"startDate" binding:"required"`
	EndDate      int64             `json:"endDate" binding:"required"`
}

type UpdateDocument struct {
	CreateDocument
}

type ListDocumentRequest struct {
	Page           int               `form:"page"`
	Size           int               `form:"size"`
	Search         string            `form:"search"`
	SearchBy       string            `form:"searchBy"`
	SortBy         string            `form:"sortBy"`
	OrderBy        string            `form:"orderBy"`
	UserID         string            `form:"userId"`
	OrganizationID string            `form:"organizationId"`
	BranchID       string            `form:"branchId"`
	DepartmentID   string            `form:"departmentId"`
	Type           enum.DocumentType `form:"type"`
	LeaveID        string            `form:"leaveId"`
	OvertimeType   enum.OverTimeType `form:"overtimeType"`
	DocumentID     string            `form:"documentId"`
}

type GetByIDDocument struct {
	ID string `uri:"id" binding:"required"`
}

type ListDocumentResponse struct {
	ID              string                    `json:"id"`
	UserID          string                    `json:"user_id"`
	FirstName       string                    `json:"first_name"`
	LastName        string                    `json:"last_name"`
	EmpCode         string                    `json:"emp_code"`
	Status          enum.StatusDocument       `json:"status"`
	Type            enum.DocumentType         `json:"type"`
	Approved        []string                  `json:"approved,omitempty"`
	Description     string                    `json:"description"`
	ApprovedCount   int                       `json:"approved_count"`
	LeaveDetails    *LeaveDocumentResponse    `json:"leave_details,omitempty"`
	OverTimeDetails *DocumentOvertimeResponse `json:"overtime_details,omitempty"`
	CreatedAt       int64                     `json:"created_at"`
	UpdatedAt       int64                     `json:"updated_at"`
}

type GetDocumentResponse struct {
	ListDocumentResponse
	Rejected   []string                  `json:"rejected,omitempty"`
	ApprovedBy []GetDocumentUserResponse `json:"approved_by"`
	RejectedBy []GetDocumentUserResponse `json:"rejected_by"`
}

type GetDocumentUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LeaveDocumentResponse struct {
	LeaveReq
	Description string  `json:"description"`
	LeaveHours  float64 `json:"leave_hours"`
	UsedQuota   float64 `json:"used_quota"`
	LeaveName   string  `json:"leave_name"`
	LeaveAmount int64   `json:"leave_amount"`
	Used        float64 `json:"used"`
}

type DocumentOvertimeResponse struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	DurationMinutes int64  `json:"durationMinutes"`
	OverTimeReq
}
