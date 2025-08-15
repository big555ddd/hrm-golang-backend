package leave

import (
	"app/app/enum"
	message "app/app/messsage"
	"app/app/model"
	leavedto "app/app/modules/leave/dto"
	"app/app/modules/user"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type Service struct {
	db   *bun.DB
	user *user.Module
}

func NewService(db *bun.DB, user *user.Module) *Service {
	return &Service{
		db:   db,
		user: user,
	}
}

func (s *Service) Create(ctx context.Context, req *leavedto.CreateLeave) (*model.Leave, error) {
	m := &model.Leave{
		Name:        req.Name,
		Description: req.Description,
		Year:        req.Year,
		Amount:      req.Amount,
	}
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(m).Exec(ctx)
		if err != nil {
			return err
		}

		leaveOrganizations := []model.LeaveOrganization{}
		for _, orgID := range req.OrganizationIDs {
			leaveOrg := &model.LeaveOrganization{
				LeaveID:        m.ID,
				OrganizationID: orgID,
			}
			leaveOrganizations = append(leaveOrganizations, *leaveOrg)
		}
		if len(leaveOrganizations) >= 0 {
			_, err = tx.NewInsert().Model(&leaveOrganizations).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return m, err

}

func (s *Service) Update(ctx context.Context, req *leavedto.UpdateLeave, id string) (*model.Leave, bool, error) {
	ex, err := s.Exist(ctx, id)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, true, errors.New(message.LeaveNotFound)
	}

	m := &model.Leave{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Year:        req.Year,
		Amount:      req.Amount,
	}
	m.SetUpdateNow()
	var mserr bool
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		_, err = tx.NewUpdate().Model(m).
			Set("name = ?", req.Name).
			Set("description = ?", req.Description).
			Set("year = ?", req.Year).
			Set("amount = ?", req.Amount).
			Set("updated_at = ?", m.UpdatedAt).
			WherePK().
			OmitZero().
			Returning("*").
			Exec(ctx)
		if err != nil {
			return err
		}
		// Clear existing organizations
		_, err = tx.NewDelete().Model(&model.LeaveOrganization{}).Where("leave_id = ?", m.ID).Exec(ctx)
		if err != nil {
			return err
		}
		leaveOrganizations := []model.LeaveOrganization{}
		for _, orgID := range req.OrganizationIDs {
			leaveOrg := &model.LeaveOrganization{
				LeaveID:        m.ID,
				OrganizationID: orgID,
			}
			leaveOrganizations = append(leaveOrganizations, *leaveOrg)
		}
		if len(leaveOrganizations) >= 0 {
			_, err = tx.NewInsert().Model(&leaveOrganizations).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return m, mserr, err
}

func (s *Service) List(ctx context.Context, req *leavedto.ListLeaveRequest) ([]leavedto.ListLeaveResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []leavedto.ListLeaveResponse{}
	query := s.db.NewSelect().
		TableExpr("leaves AS l").
		Column("l.id", "l.name", "l.description", "l.year", "l.amount", "l.created_at", "l.updated_at").
		Where("l.deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(l.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(l.name) LIKE ?", search)
		}
	}

	if req.Year != 0 {
		query.Where("l.year = ?", req.Year)
	}

	if req.OrganizationID != "" {
		query.Join("LEFT JOIN leave_organizations AS lo").
			JoinOn("lo.leave_id = l.id").
			Where("lo.organization_id = ?", req.OrganizationID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("l.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Get(ctx context.Context, id string) (*leavedto.ListLeaveResponse, error) {
	m := leavedto.ListLeaveResponse{}
	err := s.db.NewSelect().
		TableExpr("leaves AS l").
		Column("l.id", "l.name", "l.description", "l.year", "l.amount", "l.created_at", "l.updated_at").
		Where("l.id = ?", id).
		Where("l.deleted_at IS NULL").
		Scan(ctx, &m)
	if err != nil {
		return nil, err
	}
	// Fetch organizations
	leaveOrganizations := []leavedto.OrganizationLeaveResponse{}
	err = s.db.NewSelect().
		TableExpr("leave_organizations AS lo").
		ColumnExpr("lo.organization_id AS id").
		ColumnExpr("o.name AS name").
		Join("LEFT JOIN organizations AS o").
		JoinOn("o.id = lo.organization_id").
		Where("lo.leave_id = ?", id).
		Scan(ctx, &leaveOrganizations)
	if err != nil {
		return nil, err
	}
	if len(leaveOrganizations) == 0 {
		m.Organizations = []leavedto.OrganizationLeaveResponse{}
	} else {
		m.Organizations = leaveOrganizations
	}
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	ex, err := s.db.NewSelect().Model(&model.Leave{}).Where("id = ?", id).Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New(message.RoleNotFound)
	}

	//check if leave is used by document_leave
	ex, err = s.db.NewSelect().Model(&model.DocumentLeave{}).Where("leave_id = ?", id).Exists(ctx)
	if err != nil {
		return err
	}
	if ex {
		return errors.New(message.LeaveInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Leave{}).Where("id = ?", id).Exec(ctx)
	return err

}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Leave{}).Where("id = ?", id).Exists(ctx)
	return ex, err
}

func (s *Service) ExistOnOrg(ctx context.Context, organizationID, leaveID string) (bool, error) {
	ex, err := s.db.NewSelect().
		Model(&model.LeaveOrganization{}).
		Where("organization_id = ?", organizationID).
		Where("leave_id = ?", leaveID).
		Exists(ctx)

	return ex, err
}

func (s *Service) LeaveByUser(ctx context.Context, userID string, year int64) ([]leavedto.LeaveUserResponse, error) {
	user, err := s.user.Svc.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	m := []leavedto.LeaveUserResponse{}
	err = s.db.NewSelect().
		TableExpr("leaves AS l").
		Column("l.id", "l.name", "l.amount").
		Join("LEFT JOIN leave_organizations AS lo").
		JoinOn("lo.leave_id = l.id").
		Where("lo.organization_id = ?", user.OrganizationID).
		Where("l.year = ?", year).
		Where("l.deleted_at IS NULL").
		Scan(ctx, &m)
	if err != nil {
		return m, err
	}

	//get used amount for each leave
	//count form documemt where status is approved and type is leave then join
	// document_leave on DocumentID and where leave_id is l.id then count userQuota
	for i, leave := range m {
		count, err := s.GetLeaveCountByUser(ctx, userID, leave.ID)
		if err != nil {
			return m, err
		}
		m[i].UsedAmount = count
	}
	return m, nil
}

func (s *Service) GetLeaveCountByUser(ctx context.Context, userID string, leaveID string) (float64, error) {
	var count float64
	err := s.db.NewSelect().
		TableExpr("document_leaves AS dl").
		ColumnExpr("SUM(dl.used_quota) AS used_amount").
		Join("LEFT JOIN documents AS d").
		JoinOn("d.id = dl.document_id").
		Where("d.user_id = ?", userID).
		Where("dl.leave_id = ?", leaveID).
		Where("d.status = ?", enum.STATUS_DOCUMENT_APPROVED).
		Where("d.type = ?", enum.DOCUMENT_TYPE_LEAVE).
		Scan(ctx, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
