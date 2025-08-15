package holiday

import (
	message "app/app/messsage"
	"app/app/model"
	holidaydto "app/app/modules/holiday/dto"
	organization "app/app/modules/organiztion"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	db           *bun.DB
	organization *organization.Module
}

func NewService(db *bun.DB, organization *organization.Module) *Service {
	return &Service{
		db:           db,
		organization: organization,
	}
}

func (s *Service) Create(ctx context.Context, req *holidaydto.CreateHoliday) (*model.Holiday, error) {
	m := &model.Holiday{
		Name:        req.Name,
		IsActive:    req.IsActive,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().Model(m).Exec(ctx)
		if err != nil {
			return err
		}
		holidays := []model.HolidayOrganization{}
		for _, orgID := range req.OrganizationIDs {
			holidayOrg := &model.HolidayOrganization{
				HolidayID:      m.ID,
				OrganizationID: orgID,
			}
			holidays = append(holidays, *holidayOrg)
		}
		_, err = tx.NewInsert().Model(&holidays).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return m, err

}

func (s *Service) Update(ctx context.Context, req *holidaydto.UpdateHoliday, id string) (*model.Holiday, bool, error) {
	holiday, err := s.Get(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, true, errors.New(message.HolidayNotFound)
		}
		return nil, false, err
	}
	if req.StartDate != holiday.StartDate || req.EndDate != holiday.EndDate {
		if holiday.StartDate < time.Now().Unix() {
			return nil, true, errors.New(message.HolidayInUse)
		}
	}

	m := &model.Holiday{
		ID:          id,
		Name:        req.Name,
		IsActive:    req.IsActive,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
	m.SetUpdateNow()
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewUpdate().Model(m).
			Set("name = ?", req.Name).
			Set("is_active = ?", req.IsActive).
			Set("description = ?", req.Description).
			Set("start_date = ?", req.StartDate).
			Set("end_date = ?", req.EndDate).
			Set("updated_at = ?", m.UpdatedAt).
			WherePK().
			OmitZero().
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}
		holidays := []model.HolidayOrganization{}
		for _, orgID := range req.OrganizationIDs {
			holidayOrg := &model.HolidayOrganization{
				HolidayID:      m.ID,
				OrganizationID: orgID,
			}
			holidays = append(holidays, *holidayOrg)
		}
		_, err = tx.NewInsert().Model(&holidays).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	return m, false, nil
}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Holiday{}).Where("id= ?", id).Exists(ctx)
	return ex, err
}

func (s *Service) List(ctx context.Context, req *holidaydto.ListHolidayRequest) ([]model.Holiday, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []model.Holiday{}
	query := s.db.NewSelect().
		Model(&m)

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(d.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(d.name) LIKE ?", search)
		}
	}

	if req.OrganizationID != "" {
		query.Join("LEFT JOIN holiday_organization AS ho").
			JoinOn("ho.holiday_id = id").
			Where("ho.organization_id = ?", req.OrganizationID)
	}
	if req.Year > 0 {
		start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.Local).Unix()
		end := time.Date(req.Year+1, 1, 1, 0, 0, 0, 0, time.Local).Unix()
		query.Where("start_date >= ? AND end_date < ?", start, end)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*model.Holiday, error) {
	m := model.Holiday{}
	err := s.db.NewSelect().
		Model(&m).
		Where("id = ?", id).Scan(ctx)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id string) (bool, error) {
	holiday, err := s.Get(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, errors.New(message.HolidayNotFound)
		}
		return false, err
	}
	// Check if the holiday has passed (using Unix timestamps)
	currentTime := time.Now().Unix()
	if holiday.StartDate < currentTime {
		return true, errors.New(message.HolidayInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Holiday{}).Where("id = ?", id).Exec(ctx)
	return false, err

}

func (s *Service) GetHolidaysByOrganization(ctx context.Context, organizationID string) ([]model.Holiday, error) {
	holidays := []model.Holiday{}
	err := s.db.NewSelect().
		Model(&holidays).
		Join("LEFT JOIN holiday_organization AS ho ON ho.holiday_id = id").
		Where("ho.organization_id = ?", organizationID).
		Scan(ctx)
	return holidays, err
}
