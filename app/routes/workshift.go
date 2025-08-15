package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func workshift(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	workshift := router.Group("", md, log)
	{
		workshift.POST("", module.WorkShift.Ctl.Create)
		workshift.POST("/change", module.WorkShift.Ctl.ChangeWorkShift)
		workshift.PATCH("/:id", module.WorkShift.Ctl.Update)
		workshift.GET("", module.WorkShift.Ctl.List)
		workshift.GET("/:id", module.WorkShift.Ctl.Get)
		workshift.DELETE("/:id", module.WorkShift.Ctl.Delete)

	}
}
