package activitylog

import (
	"context"

	"app/app/model" // Ensure to import the model package
	"app/internal/logger"

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

func (s *Service) Create(ctx context.Context, req model.ActivityLog) (*model.ActivityLog, error) {
	if _, err := s.db.NewInsert().Model(&req).Exec(ctx); err != nil {
		logger.Infof("[error]: %v", err)
		return nil, err
	}
	return &req, nil
}
