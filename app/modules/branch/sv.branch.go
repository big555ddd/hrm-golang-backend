package branch

import (
	message "app/app/messsage"
	"app/app/model"
	branchdto "app/app/modules/branch/dto"
	organization "app/app/modules/organiztion"
	"context"
	"errors"
	"fmt"
	"strings"

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

func (s *Service) Create(ctx context.Context, req *branchdto.CreateBranch) (*model.Branch, bool, error) {
	// Check if organization exists
	ex, err := s.organization.Svc.Exist(ctx, req.OrganizationID)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New(message.OrganizationNotFound)
	}
	m := &model.Branch{
		Name:           req.Name,
		Description:    req.Description,
		OrganizationID: req.OrganizationID,
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	return m, false, err

}

func (s *Service) Update(ctx context.Context, req branchdto.UpdateBranch, id string) (*model.Branch, bool, error) {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New(message.BranchNotFound)
	}

	organizationExist, err := s.organization.Svc.Exist(ctx, req.OrganizationID)
	if err != nil {
		return nil, false, err
	}
	if !organizationExist {
		return nil, true, errors.New(message.OrganizationNotFound)
	}

	m := &model.Branch{
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
		return nil, false, err

	}
	return m, false, err
}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	exists, err := s.db.NewSelect().Model(&model.Branch{}).Where("id = ?", id).Exists(ctx)
	return exists, err
}

func (s *Service) List(ctx context.Context, req branchdto.ListBranchRequest) ([]branchdto.ListBranchResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []branchdto.ListBranchResponse{}
	query := s.db.NewSelect().
		TableExpr("branches AS b").
		Column("b.id", "b.name", "b.description", "b.created_at", "b.updated_at").
		ColumnExpr("b.organization_id AS organization_id").
		ColumnExpr("o.name AS organization_name").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id").
		Where("b.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(b.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(b.name) LIKE ?", search)
		}
	}

	if req.OrganizationID != "" {
		query.Where("b.organization_id = ?", req.OrganizationID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("b.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*branchdto.ListBranchResponse, error) {
	m := branchdto.ListBranchResponse{}
	err := s.db.NewSelect().
		TableExpr("branches AS b").
		Column("b.id", "b.name", "b.description", "b.created_at", "b.updated_at").
		ColumnExpr("b.organization_id AS organization_id").
		ColumnExpr("o.name AS organization_name").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id").
		Where("b.deleted_at IS NULL").
		Where("b.id = ?", id).Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Branch{}).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return false, err
	}

	if !ex {
		return true, errors.New(message.BranchNotFound)
	}

	// Check if branch is in use by.Branch
	ex, err = s.db.NewSelect().Model(&model.Department{}).Where("branch_id = ?", id).Exists(ctx)
	if err != nil {
		return false, err
	}

	if ex {
		return true, errors.New(message.BranchInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Branch{}).Where("id = ?", id).Exec(ctx)
	return false, err

}
