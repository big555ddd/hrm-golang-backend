package user

import (
	message "app/app/messsage"
	userdto "app/app/modules/user/dto"
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

func (ctl *Controller) Create(ctx *gin.Context) {
	body := userdto.CreateUser{}
	if err := ctx.Bind(&body); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	_, mserr, err := ctl.Service.Create(ctx, body)
	if err != nil {
		ms := message.InternalServerError
		if mserr {
			ms = err.Error()
		}
		logger.Err(err.Error())
		response.InternalServerError(ctx, ms, nil)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := userdto.GetByIDUser{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	body := userdto.UpdateUser{}
	if err := ctx.Bind(&body); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	_, mserr, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		ms := message.InternalServerError
		if mserr {
			ms = err.Error()
		}
		logger.Err(err.Error())
		response.InternalServerError(ctx, ms, nil)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := userdto.ListUserRequest{}
	if err := ctx.Bind(&req); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	err := req.Validator()
	if err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := userdto.GetByIDUser{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, err := ctl.Service.Get(ctx, ID.ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := userdto.GetByIDUser{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Err(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	err := ctl.Service.Delete(ctx, ID)
	if err != nil {
		logger.Err(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, nil)
}
