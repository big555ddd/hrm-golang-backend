package workshift

import (
	message "app/app/messsage"
	"app/app/model"
	workshiftdto "app/app/modules/workshift/dto"
	"context"
	"errors"
	"fmt"
	"strings"

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

func (s *Service) Create(ctx context.Context, req *workshiftdto.CreateWorkShift) (*model.WorkShift, error) {
	m := &model.WorkShift{
		Name:          req.Name,
		WorkLocationX: req.WorkLocationX,
		WorkLocationY: req.WorkLocationY,
		LateMinutes:   req.LateMinutes,
		Description:   req.Description,
	}
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(m).Exec(ctx)
		if err != nil {
			return err
		}
		shiftSchedules := []model.ShiftSchedule{}
		for _, schedule := range req.Schedules {
			shiftSchedule := &model.ShiftSchedule{
				WorkShiftID: m.ID,
				StartTime:   schedule.StartTime,
				EndTime:     schedule.EndTime,
				Day:         schedule.Day,
			}
			shiftSchedules = append(shiftSchedules, *shiftSchedule)
		}
		if len(shiftSchedules) >= 0 {
			_, err = tx.NewInsert().Model(&shiftSchedules).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return m, err

}

func (s *Service) Update(ctx context.Context, req *workshiftdto.UpdateWorkShift, id workshiftdto.GetByIDWorkShift) (*model.WorkShift, bool, error) {
	ex, err := s.db.NewSelect().Model(&model.WorkShift{}).Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New(message.WorkShiftNotFound)
	}

	m := &model.WorkShift{
		ID:            id.ID,
		Name:          req.Name,
		WorkLocationX: req.WorkLocationX,
		WorkLocationY: req.WorkLocationY,
		LateMinutes:   req.LateMinutes,
		Description:   req.Description,
	}
	m.SetUpdateNow()

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewUpdate().Model(m).
			Set("name = ?", req.Name).
			Set("work_location_x = ?", req.WorkLocationX).
			Set("work_location_y = ?", req.WorkLocationY).
			Set("description = ?", req.Description).
			Set("updated_at = ?", m.UpdatedAt).
			WherePK().
			OmitZero().
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}

		shiftSchedules := []model.ShiftSchedule{}
		for _, schedule := range req.Schedules {
			shiftSchedule := &model.ShiftSchedule{
				WorkShiftID: m.ID,
				StartTime:   schedule.StartTime,
				EndTime:     schedule.EndTime,
				Day:         schedule.Day,
			}
			shiftSchedules = append(shiftSchedules, *shiftSchedule)
		}
		// Delete existing schedules
		_, err = tx.NewDelete().Model(&model.ShiftSchedule{}).
			Where("work_shift_id = ?", m.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Insert new schedules
		if len(shiftSchedules) >= 0 {
			_, err = tx.NewInsert().Model(&shiftSchedules).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return m, false, err
}

func (s *Service) ChangeWorkShift(ctx context.Context, req *workshiftdto.ChangeWorkShiftRequest) error {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return s.SetWorkShift(ctx, &tx, req.WorkShiftID, req.UserID, false)
	})

	return err
}

func (s *Service) SetWorkShift(ctx context.Context, tx *bun.Tx, workShiftID string, userID string, isCreate bool) error {
	// Check if the work shift exists
	ex, err := s.Exist(ctx, workShiftID)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New(message.WorkShiftNotFound)
	}
	if !isCreate {
		ex, err = s.db.NewSelect().Model(&model.User{}).Where("id = ?", userID).Exists(ctx)
		if err != nil {
			return err
		}
		if !ex {
			return errors.New(message.UserNotFound)
		}
	}

	// Delete existing permissions for the work shift
	_, err = tx.NewDelete().Model(&model.UserWorkShift{}).
		Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return err
	}

	// Insert new permissions
	workShift := &model.UserWorkShift{
		WorkShiftID: workShiftID,
		UserID:      userID,
	}
	_, err = tx.NewInsert().Model(workShift).Exec(ctx)

	return err
}

func (s *Service) List(ctx context.Context, req workshiftdto.ListWorkShiftRequest) ([]workshiftdto.ListWorkShiftResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []workshiftdto.ListWorkShiftResponse{}
	query := s.db.NewSelect().
		TableExpr("work_shifts AS ws").
		Column("ws.id", "ws.name", "ws.late_minutes", "ws.description", "ws.work_location_x", "ws.work_location_y", "ws.created_at", "ws.updated_at").
		Where("ws.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(ws.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(ws.name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("ws.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*workshiftdto.ListWorkShiftResponse, error) {
	m := workshiftdto.ListWorkShiftResponse{}
	err := s.db.NewSelect().
		TableExpr("work_shifts AS ws").
		Column("ws.id", "ws.name", "ws.late_minutes", "ws.description", "ws.work_location_x", "ws.work_location_y", "ws.created_at", "ws.updated_at").
		Where("ws.id = ?", id).Scan(ctx, &m)
	if err != nil {
		return nil, err
	}
	// Fetch schedules
	schedules := []model.ShiftSchedule{}
	err = s.db.NewSelect().
		Model(&schedules).
		Where("work_shift_id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	if len(schedules) == 0 {
		m.Schedules = []workshiftdto.ShiftSchedule{}
		return &m, nil
	}
	m.Schedules = make([]workshiftdto.ShiftSchedule, len(schedules))
	for i, schedule := range schedules {
		m.Schedules[i] = workshiftdto.ShiftSchedule{
			Day:       schedule.Day,
			StartTime: schedule.StartTime,
			EndTime:   schedule.EndTime,
		}
	}
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id workshiftdto.GetByIDWorkShift) error {
	ex, err := s.Exist(ctx, id.ID)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New(message.WorkShiftInUse)
	}

	//check if the work shift is assigned to any user
	ex, err = s.db.NewSelect().Model(&model.UserWorkShift{}).Where("work_shift_id = ?", id.ID).Exists(ctx)
	if err != nil {
		return err
	}

	if ex {
		return errors.New(message.WorkShiftInUse)
	}

	_, err = s.db.NewDelete().Model(&model.WorkShift{}).Where("id = ?", id.ID).Exec(ctx)
	return err

}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.WorkShift{}).Where("id = ?", id).Exists(ctx)
	return ex, err
}
