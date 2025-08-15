package redisctl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	rd *redis.Client
}

func NewService(rd *redis.Client) *Service {
	return &Service{
		rd: rd,
	}
}

func (s *Service) GetJSON(ctx context.Context, key string, value any) error {
	b, err := s.rd.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, value)
}

func (s *Service) SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.rd.Set(ctx, key, b, expiration).Err()
}
