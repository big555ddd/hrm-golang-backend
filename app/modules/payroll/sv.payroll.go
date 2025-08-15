package payroll

import (
	"app/app/enum"
	"app/app/helper"
	payrolldto "app/app/modules/payroll/dto"
	"context"

	"github.com/uptrace/bun"
)

type Service struct {
	db *bun.DB
}

func NewService(db *bun.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetPayroll(ctx context.Context, req *payrolldto.CalculatePayrollRequest) (*payrolldto.PayRollResponse, error) {
	//prase month,year to unix timestamp
	start, end := helper.GetMonthRange(req.Month, req.Year)
	durationMinutes := payrolldto.PayrollMinutes{}
	err := s.db.NewSelect().TableExpr("document_overtimes AS do").
		Join("LEFT JOIN documents as d").JoinOn("do.document_id = d.id").
		ColumnExpr("COALESCE(SUM(CASE WHEN do.overtime_type = 'dayWork' THEN do.duration_minutes ELSE 0 END), 0) AS day_work_minutes").
		ColumnExpr("COALESCE(SUM(CASE WHEN do.overtime_type = 'dayOfWork' THEN do.duration_minutes ELSE 0 END), 0) AS day_of_work_minutes").
		ColumnExpr("COALESCE(SUM(CASE WHEN do.overtime_type = 'holiday' THEN do.duration_minutes ELSE 0 END), 0) AS holiday_minutes").
		Where("d.user_id = ?", req.UserID).
		Where("d.status = ?", enum.STATUS_DOCUMENT_APPROVED).
		Where("do.start_time >= ? AND do.end_time <= ?", start, end).
		Scan(ctx, &durationMinutes)
	return &payrolldto.PayRollResponse{
		DayWork:   payrolldto.CalculateResponse{Minutes: durationMinutes.DayWorkMinutes, Multiplier: 1.5, Total: (float64(durationMinutes.DayWorkMinutes) / 60.0) * req.Salary * 1.5},
		DayOfWork: payrolldto.CalculateResponse{Minutes: durationMinutes.DayOfWorkMinutes, Multiplier: 2.0, Total: (float64(durationMinutes.DayOfWorkMinutes) / 60.0) * req.Salary * 2.0},
		Holiday:   payrolldto.CalculateResponse{Minutes: durationMinutes.HolidayMinutes, Multiplier: 2.5, Total: (float64(durationMinutes.HolidayMinutes) / 60.0) * req.Salary * 2.5},
	}, err
}
