package notification

import (
	"app/app/model"
	notificationdto "app/app/modules/notification/dto"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	db  *bun.DB
	Hub *Hub
}

func NewService(db *bun.DB) *Service {
	hub := GetHub() // Use singleton hub

	service := &Service{
		db:  db,
		Hub: hub,
	}

	return service
}

// BroadcastToAll sends message to all connected clients
func (s *Service) BroadcastToAll(message interface{}) {
	s.Hub.BroadcastToAll(message)
}

func (s *Service) List(ctx context.Context, req notificationdto.ListNotificationRequest) ([]model.Notification, int, error) {
	resp := []model.Notification{}
	query := s.db.NewSelect().
		Model(&model.Notification{})

	if req.UserID != "" {
		query.Where("user_id = ?", req.UserID)
	}
	if req.IsRead != nil {
		query.Where("is_read = ?", *req.IsRead)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return resp, 0, nil
	}

	err = query.Scan(ctx, &resp)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

// SendToUser sends message to specific user
func (s *Service) SendToUser(userID string, message interface{}) {
	s.Hub.SendToUser(userID, message)
}

func (s *Service) SendToUserMulti(notifications []model.Notification) {
	for _, notification := range notifications {
		s.SendToUser(notification.UserID, notification)
	}
}

func (s *Service) NotiToUser(ctx context.Context, tx bun.Tx, req []notificationdto.NotificationRequest) error {
	m := []model.Notification{}
	for _, notification := range req {
		m = append(m, model.Notification{
			UserID:     notification.UserID,
			Message:    notification.Message,
			Type:       notification.Action,
			IsRead:     false,
			DocumentID: notification.DocumentID,
		})
	}
	_, err := tx.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		return err
	}
	s.SendToUserMulti(m)

	return nil
}

func (s *Service) NotiToUserSingle(ctx context.Context, req *notificationdto.NotificationRequest) error {
	m := model.Notification{
		UserID:     req.UserID,
		Message:    req.Message,
		Type:       req.Action,
		IsRead:     false,
		DocumentID: req.DocumentID,
	}
	_, err := s.db.NewInsert().Model(&m).Exec(ctx)
	if err != nil {
		return err
	}
	s.SendToUser(req.UserID, m)

	return nil
}

func (s *Service) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	_, err := s.db.NewUpdate().
		Model(&model.Notification{}).
		Set("is_read = ?", true).
		Set("updated_at = ?", time.Now().Unix()).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Exec(ctx)

	if err != nil {
		return err
	}
	count, err := s.GetUnreadCount(ctx, userID)
	if err != nil {
		return err
	}

	// Emit read event to user
	readEvent := notificationdto.NotificationRead{
		Type: "notification_read",
		Data: struct {
			NotificationID string `json:"notificationId"`
			IsRead         bool   `json:"isRead"`
			UserID         string `json:"userId"`
			Count          int    `json:"count"`
		}{
			NotificationID: notificationID,
			IsRead:         true,
			UserID:         userID,
			Count:          count,
		},
	}

	s.SendToUser(userID, readEvent)
	return nil
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := s.db.NewUpdate().
		Model(&model.Notification{}).
		Set("is_read = ?", true).
		Set("updated_at = ?", time.Now().Unix()).
		Where("user_id = ? AND is_read = ?", userID, false).
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	count, err := s.db.NewSelect().
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(ctx)

	return count, err
}

// IsUserOnline checks if user has active WebSocket connections
func (s *Service) IsUserOnline(userID string) bool {
	return s.Hub.IsUserOnline(userID)
}

// GetActiveUsers returns list of users with active connections
func (s *Service) GetActiveUsers() []string {
	return s.Hub.GetActiveUsers()
}

// SendToUserWithFallback sends notification via WebSocket if online, stores in DB regardless
func (s *Service) SendToUserWithFallback(ctx context.Context, userID string, message interface{}, notificationData *model.Notification) error {
	// Always store notification in database
	if notificationData != nil {
		_, err := s.db.NewInsert().Model(notificationData).Exec(ctx)
		if err != nil {
			return err
		}
	}

	// Send via WebSocket if user is online
	if s.IsUserOnline(userID) {
		s.SendToUser(userID, message)
	}

	return nil
}
