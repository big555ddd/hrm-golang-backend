package department

import (
	message "app/app/messsage"
	"app/app/model"
	"app/app/modules/branch"
	departmentdto "app/app/modules/department/dto"
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
	branch       *branch.Module
}

func NewService(db *bun.DB, organization *organization.Module, branch *branch.Module) *Service {
	return &Service{
		db:           db,
		organization: organization,
		branch:       branch,
	}
}

func (s *Service) Create(ctx context.Context, req *departmentdto.CreateDepartment) (*model.Department, bool, error) {

	ex, err := s.branch.Svc.Exist(ctx, req.BranchID)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New(message.BranchNotFound)
	}

	m := &model.Department{
		Name:        req.Name,
		Description: req.Description,
		BranchID:    req.BranchID,
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, false, err
	}

	return m, false, nil

}

func (s *Service) Update(ctx context.Context, req *departmentdto.UpdateDepartment, id string) (*model.Department, bool, error) {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New(message.DepartmentNotFound)
	}

	ex, err = s.branch.Svc.Exist(ctx, req.BranchID)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New(message.BranchNotFound)
	}

	m := &model.Department{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		BranchID:    req.BranchID,
	}
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?", req.Name).
		Set("description = ?", req.Description).
		Set("branch_id = ?", req.BranchID).
		Set("updated_at = ?", m.UpdatedAt).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, false, err
	}
	return m, false, nil
}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Department{}).Where("id= ?", id).Exists(ctx)
	return ex, err
}

func (s *Service) List(ctx context.Context, req *departmentdto.ListDepartmentRequest) ([]departmentdto.ListDepartmentResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []departmentdto.ListDepartmentResponse{}
	query := s.db.NewSelect().
		TableExpr("departments AS d").
		Column("d.id", "d.name", "d.description", "d.created_at", "d.updated_at").
		ColumnExpr("d.branch_id AS branch_id").
		ColumnExpr("b.name AS branch_name").
		ColumnExpr("o.id AS organization_id").
		ColumnExpr("o.name AS organization_name").
		Join("LEFT JOIN branches AS b").
		JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id").
		Where("d.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(d.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(d.name) LIKE ?", search)
		}
	}

	if req.BranchID != "" {
		query.Where("d.branch_id = ?", req.BranchID)
	}

	if req.OrganizationID != "" {
		query.Where("b.organization_id = ?", req.OrganizationID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("d.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*departmentdto.ListDepartmentResponse, error) {
	m := departmentdto.ListDepartmentResponse{}
	err := s.db.NewSelect().
		TableExpr("departments AS d").
		Column("d.id", "d.name", "d.description", "d.created_at", "d.updated_at").
		ColumnExpr("d.branch_id AS branch_id").
		ColumnExpr("b.name AS branch_name").
		ColumnExpr("o.id AS organization_id").
		ColumnExpr("o.name AS organization_name").
		Join("LEFT JOIN branches AS b").
		JoinOn("d.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").
		JoinOn("b.organization_id = o.id").
		Where("d.deleted_at IS NULL").
		Where("d.id = ?", id).Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Department{}).Where("id = ?", id).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return false, err
	}

	if !ex {
		return true, errors.New(message.DepartmentNotFound)
	}

	// Check if department is in use by userDepartment
	ex, err = s.db.NewSelect().Model(&model.UserDepartment{}).Where("department_id = ?", id).Exists(ctx)
	if err != nil {
		return false, err
	}

	if ex {
		return true, errors.New(message.DepartmentInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Department{}).Where("id = ?", id).Exec(ctx)
	return false, err

}
