package auth

import (
	"app/app/helper"
	authdto "app/app/modules/auth/dto"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		Service: svc,
	}
}

func (ctl *Controller) Login(ctx *gin.Context) {
	req := authdto.LoginRequest{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	// Call the service to login the user
	data, err := ctl.Service.Login(ctx, &req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, data)
}

func (ctl *Controller) ForgotPassword(ctx *gin.Context) {
	req := authdto.ForgotPasswordRequest{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	data, err := ctl.Service.ForgotPassword(ctx, &req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, data)
}

func (ctl *Controller) ResetPassword(ctx *gin.Context) {
	req := authdto.ResetPasswordRequest{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	err := ctl.Service.ResetPassword(ctx, &req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, nil)
}

func (ctl *Controller) Info(ctx *gin.Context) {
	user, _ := helper.GetUserByToken(ctx)
	data, err := ctl.Service.Info(ctx, user.Data.ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, data)
}
