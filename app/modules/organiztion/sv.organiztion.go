package organization

import (
	message "app/app/messsage"
	"app/app/model"
	organizationdto "app/app/modules/organiztion/dto"
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

func (s *Service) Create(ctx context.Context, req *organizationdto.CreateOrganization) (*model.Organization, error) {
	m := &model.Organization{
		Name:        req.Name,
		Description: req.Description,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, errors.New("email already exists")
		}
	}

	return m, err

}

func (s *Service) Update(ctx context.Context, req *organizationdto.UpdateOrganization, id string) (*model.Organization, error) {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, err
	}

	m := &model.Organization{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?", req.Name).
		Set("description = ?", req.Description).
		Set("updated_at = ?", m.UpdatedAt).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return m, err
}

func (s *Service) List(ctx context.Context, req *organizationdto.ListOrganizationRequest) ([]organizationdto.ListOrganizationResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []organizationdto.ListOrganizationResponse{}
	query := s.db.NewSelect().
		TableExpr("organizations AS o").
		Column("o.id", "o.name", "o.description", "o.created_at", "o.updated_at").
		Where("o.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(o.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(o.name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("o.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*model.Organization, error) {
	m := model.Organization{}
	err := s.db.NewSelect().
		Model(&m).
		Where("id = ?", id).Scan(ctx, &m)
	return &m, err
}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Organization{}).Where("id = ?", id).Exists(ctx)
	return ex, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New(message.OrganizationNotFound)
	}

	//check if organization has use by branch
	ex, err = s.db.NewSelect().Model(&model.Branch{}).Where("organization_id = ?", id).Exists(ctx)
	if err != nil {
		return err
	}

	if ex {
		return errors.New(message.OrganizationInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Organization{}).Where("id = ?", id).Exec(ctx)
	return err

}
