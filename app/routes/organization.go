package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func organization(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	organization := router.Group("", md, log)
	{
		organization.POST("", module.Organization.Ctl.Create)
		organization.PATCH("/:id", module.Organization.Ctl.Update)
		organization.GET("", module.Organization.Ctl.List)
		organization.GET("/:id", module.Organization.Ctl.Get)
		organization.DELETE("/:id", module.Organization.Ctl.Delete)

	}
}
