package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func user(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()

	user := router.Group("", md, log)
	{
		user.POST("", module.User.Ctl.Create)
		user.PATCH("/:id", module.User.Ctl.Update)
		user.GET("", module.User.Ctl.List)
		user.GET("/:id", module.User.Ctl.Get)
		user.DELETE("/:id", module.User.Ctl.Delete)

	}
}
