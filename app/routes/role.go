package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func role(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	role := router.Group("", md, log)
	{
		role.POST("", module.Role.Ctl.Create)
		role.PATCH("/:id", module.Role.Ctl.Update)
		role.GET("", module.Role.Ctl.List)
		role.GET("/:id", module.Role.Ctl.Get)
		role.DELETE("/:id", module.Role.Ctl.Delete)
		role.POST("/set-permissions", module.Role.Ctl.SetPermissions)
		role.GET("/:id/permissions", module.Role.Ctl.GetRolePermissions)
		role.GET("/permissions", module.Role.Ctl.Permission)

	}
}
