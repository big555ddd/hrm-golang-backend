package document

import (
	"app/app/helper"
	message "app/app/messsage"
	documentdto "app/app/modules/document/dto"
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
	body := documentdto.CreateDocument{}

	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	_, err := ctl.Service.Create(ctx, &body)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, message.InternalServerError, nil)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := documentdto.ListDocumentRequest{
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
	ID := documentdto.GetByIDDocument{}
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
	ID := documentdto.GetByIDDocument{}
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

func (ctl *Controller) Approved(ctx *gin.Context) {
	ID := documentdto.GetByIDDocument{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	user, _ := helper.GetUserByToken(ctx)

	err := ctl.Service.Approved(ctx, ID.ID, user.Data.ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, nil)
}

func (ctl *Controller) Rejected(ctx *gin.Context) {
	ID := documentdto.GetByIDDocument{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error(), nil)
		return
	}

	user, _ := helper.GetUserByToken(ctx)

	err := ctl.Service.Rejected(ctx, ID.ID, user.Data.ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, nil)
}
