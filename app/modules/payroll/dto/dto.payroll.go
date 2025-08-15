package payrolldto

type GetByIDPayroll struct {
	ID string `uri:"id" binding:"required"`
}

type CalculatePayrollRequest struct {
	UserID string  `json:"userId" binding:"required"`
	Salary float64 `json:"salary" binding:"required"`
	Month  int     `json:"month" binding:"required"`
	Year   int     `json:"year" binding:"required"`
}

type PayrollMinutes struct {
	DayWorkMinutes   int64 `json:"dayWorkMinutes"`
	DayOfWorkMinutes int64 `json:"dayOfWork"`
	HolidayMinutes   int64 `json:"holidayMinutes"`
}

type PayRollResponse struct {
	DayWork   CalculateResponse `json:"dayWork"`
	DayOfWork CalculateResponse `json:"dayOfWork"`
	Holiday   CalculateResponse `json:"holiday"`
}

type CalculateResponse struct {
	Minutes    int64   `json:"minutes"`
	Multiplier float64 `json:"multiplier"`
	Total      float64 `json:"total"`
}
