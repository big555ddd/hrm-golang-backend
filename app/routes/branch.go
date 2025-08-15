package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func branch(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	branch := router.Group("", md, log)
	{
		branch.POST("", module.Branch.Ctl.Create)
		branch.PATCH("/:id", module.Branch.Ctl.Update)
		branch.GET("", module.Branch.Ctl.List)
		branch.GET("/:id", module.Branch.Ctl.Get)
		branch.DELETE("/:id", module.Branch.Ctl.Delete)

	}
}
