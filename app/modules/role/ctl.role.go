package role

import (
	message "app/app/messsage"
	roledto "app/app/modules/role/dto"
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
	body := roledto.CreateRole{}

	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, err := ctl.Service.Create(ctx, &body)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, message.InternalServerError, nil)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := roledto.GetByIDRole{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	body := roledto.UpdateRole{}
	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	data, mserr, err := ctl.Service.Update(ctx, &body, ID.ID)
	if err != nil {
		ms := message.InternalServerError
		if mserr {
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalServerError(ctx, ms, nil)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := roledto.ListRoleRequest{
		Page:    1,
		Size:    10,
		OrderBy: "asc",
		SortBy:  "created_at",
	}
	if err := ctx.Bind(&req); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, total, err := ctl.Service.List(ctx, &req)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := roledto.GetByIDRole{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, err := ctl.Service.Get(ctx, ID.ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := roledto.GetByIDRole{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	err := ctl.Service.Delete(ctx, ID.ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, nil)
}

func (ctl *Controller) SetPermissions(ctx *gin.Context) {
	body := roledto.SetRolePermissions{}
	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	err := ctl.Service.SetRolePermissions(ctx, &body)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) GetRolePermissions(ctx *gin.Context) {
	ID := roledto.GetByIDRole{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	data, err := ctl.Service.GetRolePermissions(ctx, ID.ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Permission(ctx *gin.Context) {
	data, err := ctl.Service.Permission(ctx)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, data)
}
