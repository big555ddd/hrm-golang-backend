package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func leave(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	leave := router.Group("", md, log)
	{
		leave.POST("", module.Leave.Ctl.Create)
		leave.PATCH("/:id", module.Leave.Ctl.Update)
		leave.GET("", module.Leave.Ctl.List)
		leave.GET("/:id", module.Leave.Ctl.Get)
		leave.DELETE("/:id", module.Leave.Ctl.Delete)

	}
}
