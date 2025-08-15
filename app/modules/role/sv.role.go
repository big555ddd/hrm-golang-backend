package role

import (
	message "app/app/messsage"
	"app/app/model"
	roledto "app/app/modules/role/dto"
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

func (s *Service) Create(ctx context.Context, req *roledto.CreateRole) (*model.Role, error) {
	m := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return m, err

}

func (s *Service) Update(ctx context.Context, req *roledto.UpdateRole, id string) (*model.Role, bool, error) {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return nil, false, err
	}
	if !ex {
		return nil, true, errors.New(message.RoleNotFound)
	}

	m := &model.Role{
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

func (s *Service) List(ctx context.Context, req *roledto.ListRoleRequest) ([]roledto.ListRoleResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []roledto.ListRoleResponse{}
	query := s.db.NewSelect().
		TableExpr("roles AS r").
		Column("r.id", "r.name", "r.description", "r.created_at", "r.updated_at").
		Where("r.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(r.name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*roledto.ListRoleResponse, error) {
	m := roledto.ListRoleResponse{}
	err := s.db.NewSelect().
		TableExpr("roles AS r").
		Column("r.id", "r.name", "r.description", "r.created_at", "r.updated_at").
		Where("r.id = ?", id).Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return err
	}
	if !ex {
		return errors.New(message.RoleNotFound)
	}

	//check if role is used by any user
	ex, err = s.db.NewSelect().Model(&model.UserRole{}).Where("role_id = ?", id).Exists(ctx)
	if err != nil {
		return err
	}

	if ex {
		return errors.New(message.RoleInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Role{}).Where("id = ?", id).Exec(ctx)
	return err

}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Role{}).Where("id = ?", id).Exists(ctx)
	return ex, err
}

func (s *Service) SetRolePermissions(ctx context.Context, req *roledto.SetRolePermissions) error {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// First, delete existing permissions for the role
		_, err := tx.NewDelete().Model(&model.RolePermission{}).Where("role_id = ?", req.RoleID).Exec(ctx)
		if err != nil {
			return err
		}
		rolePermission := []*model.RolePermission{}
		// Then, insert new permissions
		for _, permissionID := range req.PermissionIDs {
			rolePermission = append(rolePermission, &model.RolePermission{
				RoleID:       req.RoleID,
				PermissionID: permissionID,
			})
		}

		_, err = tx.NewInsert().Model(&rolePermission).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *Service) GetRolePermissions(ctx context.Context, roleID string) ([]roledto.RolePermissionResponse, error) {
	permissions := []roledto.RolePermissionResponse{}
	err := s.db.NewSelect().
		TableExpr("role_permissions AS rp").
		Column("rp.permission_id").
		ColumnExpr("p.name AS permission_name").
		Join("JOIN permissions AS p ON rp.permission_id = p.id").
		Where("rp.role_id = ?", roleID).
		Scan(ctx, &permissions)
	return permissions, err
}

func (s *Service) GetRolePermissionsName(ctx context.Context, roleID string) ([]string, error) {
	permissions := []string{}
	err := s.db.NewSelect().
		TableExpr("role_permissions AS rp").
		ColumnExpr("p.name AS permission_name").
		Join("JOIN permissions AS p ON rp.permission_id = p.id").
		Where("rp.role_id = ?", roleID).
		Scan(ctx, &permissions)
	return permissions, err
}

func (s *Service) Permission(ctx context.Context) ([]model.Permission, error) {
	permissions := []model.Permission{}
	err := s.db.NewSelect().Model(&permissions).Scan(ctx)
	return permissions, err
}
