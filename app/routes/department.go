package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func department(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	department := router.Group("", md, log)
	{
		department.POST("", module.Department.Ctl.Create)
		department.PATCH("/:id", module.Department.Ctl.Update)
		department.GET("", module.Department.Ctl.List)
		department.GET("/:id", module.Department.Ctl.Get)
		department.DELETE("/:id", module.Department.Ctl.Delete)

	}
}
