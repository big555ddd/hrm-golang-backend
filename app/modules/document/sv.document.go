package document

import (
	"app/app/enum"
	message "app/app/messsage"
	"app/app/model"
	"app/app/modules/attendance"
	documentdto "app/app/modules/document/dto"
	"app/app/modules/leave"
	"app/app/modules/notification"
	notificationdto "app/app/modules/notification/dto"
	"app/app/modules/user"
	userdto "app/app/modules/user/dto"
	"app/app/modules/workshift"
	"context"
	"errors"
	"fmt"
	"strings"

	"slices"

	"github.com/uptrace/bun"
)

type Service struct {
	db           *bun.DB
	user         *user.Module
	leave        *leave.Module
	workshift    *workshift.Module
	attendance   *attendance.Module
	notification *notification.Module
}

func NewService(db *bun.DB, user *user.Module,
	leave *leave.Module, workshift *workshift.Module,
	attendance *attendance.Module, notification *notification.Module) *Service {
	return &Service{
		db:           db,
		user:         user,
		leave:        leave,
		workshift:    workshift,
		attendance:   attendance,
		notification: notification,
	}
}

func (s *Service) Create(ctx context.Context, req *documentdto.CreateDocument) (*model.Document, error) {
	user, err := s.user.Svc.Get(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	m := &model.Document{
		UserID:      req.UserID,
		Status:      enum.STATUS_DOCUMENT_PENDING,
		Type:        req.Type,
		Description: req.Description,
	}
	var documentDetail any
	if req.Type == enum.DOCUMENT_TYPE_LEAVE {
		quota, hours, err := s.calculateLeaveQuotaAndHours(ctx, req.Leave, req.UserID)
		if err != nil {
			return nil, err
		}
		documentDetail = &model.DocumentLeave{
			LeaveID:     req.Leave.LeaveID,
			Description: req.Description,
			StartDate:   req.Leave.StartDate,
			EndDate:     req.Leave.EndDate,
			LeaveHours:  hours,
			UsedQuota:   quota,
		}
		//calculate leave Quota and hours
	} else if req.Type == enum.DOCUMENT_TYPE_OVERTIME {
		minutes, err := s.calculateOvertimeDuration(req.OverTime)
		if err != nil {
			return nil, err
		}
		documentDetail = &model.DocumentOvertime{
			OverTimeType:    req.OverTime.OverTimeType,
			Description:     req.Description,
			StartDate:       req.OverTime.StartDate,
			EndDate:         req.OverTime.EndDate,
			DurationMinutes: minutes,
		}

	} else {
		documentDetail = nil
	}
	doctype := ""
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(m).Exec(ctx)
		if err != nil {
			return err
		}
		if documentDetail != nil {
			switch v := documentDetail.(type) {
			case *model.DocumentLeave:
				v.DocumentID = m.ID
				doctype = "leave"
			case *model.DocumentOvertime:
				v.DocumentID = m.ID
				doctype = "overtime"
			}
			_, err = tx.NewInsert().Model(documentDetail).Exec(ctx)
			if err != nil {
				return err
			}
		}
		//notification to user who have permission document:update
		listReq := userdto.ListUserRequest{
			Page:           1,
			Size:           999,
			OrganizationID: user.OrganizationID,
			ForNoti:        user.ID,
			WithPermission: "document:update",
			OrderBy:        "asc",
			SortBy:         "created_at",
		}
		alluser, count, err := s.user.Svc.List(ctx, listReq)
		if err != nil {
			return err
		}

		if count == 0 {
			return nil
		}
		notimessage := fmt.Sprintf("%s %s %s has created a new %s document", user.EmpCode, user.FirstName, user.LastName, doctype)

		notiReq := []notificationdto.NotificationRequest{}
		for _, u := range alluser {
			notiReq = append(notiReq, notificationdto.NotificationRequest{
				UserID:     u.ID,
				Action:     "document:create",
				Message:    notimessage,
				DocumentID: m.ID,
			})
		}
		err = s.notification.Svc.NotiToUser(ctx, tx, notiReq)
		if err != nil {
			return err
		}
		return nil
	})
	return m, err

}

func (s *Service) calculateOvertimeDuration(overtime *documentdto.OverTimeReq) (int64, error) {
	// Convert milliseconds to seconds
	startDate := overtime.StartDate
	endDate := overtime.EndDate

	// Validate start and end times
	if startDate >= endDate {
		return 0, errors.New(message.InvalidRequest)
	}

	// Calculate duration in seconds
	durationSeconds := endDate - startDate

	// Convert to minutes and round down (truncate any remaining seconds)
	durationMinutes := durationSeconds / 60

	if durationMinutes <= 0 {
		return 0, errors.New(message.InvalidRequest)
	}

	return durationMinutes, nil
}

func (s *Service) List(ctx context.Context, req *documentdto.ListDocumentRequest) ([]documentdto.ListDocumentResponse, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []documentdto.ListDocumentResponse{}
	query := s.db.NewSelect().
		TableExpr("documents AS d").
		Column("d.id", "d.user_id", "d.status", "d.type", "d.approved",
			"d.description", "d.created_at", "d.updated_at").
		ColumnExpr("CASE WHEN jsonb_typeof(d.approved) = 'array' THEN jsonb_array_length(d.approved) ELSE 0 END AS approved_count").
		ColumnExpr("u.first_name").
		ColumnExpr("u.last_name").
		ColumnExpr("u.emp_code").
		Join("LEFT JOIN users AS u").JoinOn("u.id = d.user_id").
		Join("LEFT JOIN user_departments AS ud").JoinOn("u.id = ud.user_id").
		Join("LEFT JOIN departments AS dpm").JoinOn("ud.department_id = dpm.id").
		Join("LEFT JOIN branches AS b").JoinOn("dpm.branch_id = b.id").
		Join("LEFT JOIN organizations AS o").JoinOn("b.organization_id = o.id")

	if req.Search != "" {
		search := fmt.Sprint("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(d.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(d.description) LIKE ?", search)
		}
	}

	if req.UserID != "" {
		query.Where("d.user_id = ?", req.UserID)
	}

	if req.Type != "" {
		query.Where("d.type = ?", req.Type)
	}

	if req.LeaveID != "" {
		query.Join("LEFT JOIN document_leaves AS dl").
			JoinOn("dl.document_id = d.id").
			Where("dl.leave_id = ?", req.LeaveID)
	}

	if req.OvertimeType != "" {
		query.Join("LEFT JOIN document_overtimes AS dot").
			JoinOn("dot.document_id = d.id").
			Where("dot.overtime_type = ?", req.OvertimeType)
	}

	if req.OrganizationID != "" {
		query.Where("o.id = ?", req.OrganizationID)
	}

	if req.BranchID != "" {
		query.Where("b.id = ?", req.BranchID)
	}

	if req.DepartmentID != "" {
		query.Where("dpm.id = ?", req.DepartmentID)
	}

	if req.DocumentID != "" {
		query.Where("d.id = ?", req.DocumentID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Custom order by status priority and pending approval count
	// Priority: 1/2 pending (has some approvals), 0/2 pending (no approvals), approved, rejected
	err = query.OrderExpr("CASE WHEN d.status = 'pending' AND jsonb_typeof(d.approved) = 'array' AND jsonb_array_length(d.approved) > 0 THEN 1 WHEN d.status = 'pending' AND (d.approved IS NULL OR jsonb_typeof(d.approved) != 'array' OR jsonb_array_length(d.approved) = 0) THEN 2 WHEN d.status = 'approved' THEN 3 WHEN d.status = 'rejected' THEN 4 ELSE 5 END ASC").
		OrderExpr("d.created_at DESC").
		Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err

}

func (s *Service) Approved(ctx context.Context, id string, userID string) error {
	doc := &model.Document{}
	err := s.db.NewSelect().Model(doc).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	if doc.Status != enum.STATUS_DOCUMENT_PENDING {
		return errors.New(message.DocumentInUse)
	}

	if len(doc.Approved) != 0 {
		//check if user already approved
		if slices.Contains(doc.Approved, userID) {
			return errors.New(message.DocumentAlreadyApproved)
		}
	}

	doc.Approved = append(doc.Approved, userID)
	if len(doc.Approved) >= 2 {
		doc.Status = enum.STATUS_DOCUMENT_APPROVED
		if doc.Type == enum.DOCUMENT_TYPE_LEAVE {
			err := s.createLeaveRecord(ctx, doc.ID)
			if err != nil {
				return err
			}
		}
	} else {
		doc.Status = enum.STATUS_DOCUMENT_PENDING
	}
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		_, err = tx.NewUpdate().Model(doc).Set("approved = ?", doc.Approved).
			Set("status = ?", doc.Status).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return err
		}
		if doc.Status == enum.STATUS_DOCUMENT_APPROVED {
			user, err := s.user.Svc.Get(ctx, doc.UserID)
			if err != nil {
				return err
			}
			//notification to user
			notimessage := fmt.Sprintf("%s %s document has been approved", user.FirstName, user.LastName)
			notiReq := notificationdto.NotificationRequest{
				UserID:     user.ID,
				Action:     "document:approved",
				Message:    notimessage,
				DocumentID: doc.ID,
			}

			err = s.notification.Svc.NotiToUser(ctx, tx, []notificationdto.NotificationRequest{notiReq})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *Service) Rejected(ctx context.Context, id string, userID string) error {
	doc := &model.Document{}
	err := s.db.NewSelect().Model(doc).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	if doc.Status != enum.STATUS_DOCUMENT_PENDING {
		return errors.New(message.DocumentInUse)
	}

	doc.Rejected = append(doc.Rejected, userID)
	doc.Status = enum.STATUS_DOCUMENT_REJECTED

	_, err = s.db.NewUpdate().Model(doc).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if doc.Status == enum.STATUS_DOCUMENT_REJECTED {
		user, err := s.user.Svc.Get(ctx, doc.UserID)
		if err != nil {
			return err
		}
		//notification to user
		notimessage := fmt.Sprintf("%s %s document has been rejected", user.FirstName, user.LastName)
		notiReq := notificationdto.NotificationRequest{
			UserID:     user.ID,
			Action:     "document:rejected",
			Message:    notimessage,
			DocumentID: doc.ID,
		}

		err = s.notification.Svc.NotiToUserSingle(ctx, &notiReq)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Service) Get(ctx context.Context, id string) (*documentdto.GetDocumentResponse, error) {
	m := documentdto.GetDocumentResponse{}
	err := s.db.NewSelect().
		TableExpr("documents AS d").
		Column("d.id", "d.user_id", "d.status", "d.type",
			"d.description", "d.approved", "d.rejected",
			"d.created_at", "d.updated_at").
		ColumnExpr("CASE WHEN jsonb_typeof(d.approved) = 'array' THEN jsonb_array_length(d.approved) ELSE 0 END AS approved_count").
		ColumnExpr("u.first_name").
		ColumnExpr("u.last_name").
		ColumnExpr("u.emp_code").
		Join("LEFT JOIN users AS u").JoinOn("u.id = d.user_id").
		Where("d.id = ?", id).Scan(ctx, &m)
	if err != nil {
		return nil, err
	}

	// Fetch approved by users details
	if len(m.Approved) > 0 {
		approvedBy := []documentdto.GetDocumentUserResponse{}
		err := s.db.NewSelect().
			TableExpr("users AS u").
			Column("u.id", "u.first_name", "u.last_name").
			Where("u.id IN (?)", bun.In(m.Approved)).
			Scan(ctx, &approvedBy)
		if err != nil {
			return nil, err
		}
		m.ApprovedBy = approvedBy
		m.Approved = nil
	} else {
		m.ApprovedBy = []documentdto.GetDocumentUserResponse{}
	}

	// Fetch rejected by users details
	if len(m.Rejected) > 0 {
		rejectedBy := []documentdto.GetDocumentUserResponse{}
		err := s.db.NewSelect().
			TableExpr("users AS u").
			Column("u.id", "u.first_name", "u.last_name").
			Where("u.id IN (?)", bun.In(m.Rejected)).
			Scan(ctx, &rejectedBy)
		if err != nil {
			return nil, err
		}
		m.RejectedBy = rejectedBy
		m.Rejected = nil
	} else {
		m.RejectedBy = []documentdto.GetDocumentUserResponse{}
	}

	// If the document type is leave, fetch leave details
	if m.Type == enum.DOCUMENT_TYPE_LEAVE {
		leaveDetails := &documentdto.LeaveDocumentResponse{}
		err := s.db.NewSelect().
			TableExpr("document_leaves AS dl").
			Column("dl.leave_id", "dl.start_date", "dl.end_date", "dl.description",
				"dl.leave_hours", "dl.used_quota").
			ColumnExpr("l.name AS leave_name").
			ColumnExpr("l.amount AS leave_amount").
			Where("dl.document_id = ?", m.ID).Join("LEFT JOIN leaves AS l").JoinOn("l.id = dl.leave_id").Scan(ctx, leaveDetails)
		if err != nil {
			return nil, err
		}
		m.LeaveDetails = leaveDetails
		//get used quota for leave
		usedQuota, err := s.leave.Svc.GetLeaveCountByUser(ctx, m.UserID, leaveDetails.LeaveID)
		if err != nil {
			return nil, err
		}
		leaveDetails.Used = usedQuota
	} else if m.Type == enum.DOCUMENT_TYPE_OVERTIME {
		overTimeDetails := &documentdto.DocumentOvertimeResponse{}
		err := s.db.NewSelect().
			TableExpr("document_overtimes AS dot").
			Column("dot.overtime_type", "dot.start_date", "dot.end_date", "dot.id", "dot.description",
				"dot.duration_minutes").
			Where("dot.document_id = ?", m.ID).Scan(ctx, overTimeDetails)
		if err != nil {
			return nil, err
		}
		m.OverTimeDetails = overTimeDetails
	}
	return &m, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	doc, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if doc.Status != enum.STATUS_DOCUMENT_PENDING {
		return errors.New(message.DocumentInUse)
	}

	_, err = s.db.NewDelete().Model(&model.Document{}).Where("id = ?", id).Exec(ctx)
	return err

}

func (s *Service) Exist(ctx context.Context, id string) (bool, error) {
	ex, err := s.db.NewSelect().Model(&model.Document{}).Where("id = ?", id).Exists(ctx)
	return ex, err
}
