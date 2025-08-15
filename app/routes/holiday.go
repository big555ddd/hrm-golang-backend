package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func holiday(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	holiday := router.Group("", md, log)
	{
		holiday.POST("", module.Holiday.Ctl.Create)
		holiday.PATCH("/:id", module.Holiday.Ctl.Update)
		holiday.GET("", module.Holiday.Ctl.List)
		holiday.GET("/:id", module.Holiday.Ctl.Get)
		holiday.DELETE("/:id", module.Holiday.Ctl.Delete)

	}
}
