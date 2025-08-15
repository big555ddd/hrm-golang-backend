package payroll

import (
	payrolldto "app/app/modules/payroll/dto"
	"app/app/response"

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

func (c *Controller) GetPayroll(ctx *gin.Context) {
	var req payrolldto.CalculatePayrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error(), nil)
		return
	}
	resp, err := c.Service.GetPayroll(ctx, &req)
	if err != nil {
		response.InternalServerError(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, resp)
}
