package redisctl

import (
	redisctldto "app/app/modules/redisctl/dto"
	"app/app/response"
	"app/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	Service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		Service: svc,
	}
}

func (ctl *Controller) Create(ctx *gin.Context) {
	session := map[string]interface{}{}
	if err := ctx.Bind(&session); err != nil {
		logger.Err(err.Error())

		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	key := uuid.New().String()
	err := ctl.Service.SetJSON(ctx, key, session, time.Hour)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, key)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	key := redisctldto.GetByIDRedis{}
	if err := ctx.BindUri(&key); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	session := map[string]interface{}{}
	err := ctl.Service.GetJSON(ctx, key.ID, &session)
	if err != nil {
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, session)
}
