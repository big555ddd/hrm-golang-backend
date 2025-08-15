package redisctl

import (
	"github.com/redis/go-redis/v9"
)

type Module struct {
	Ctl *Controller
	Svc *Service
}

func NewModule(rd *redis.Client) *Module {
	svc := NewService(rd)
	return &Module{
		Ctl: NewController(svc),
		Svc: svc,
	}
}
