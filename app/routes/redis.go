package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func redis(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	redis := router.Group("", md, log)
	{
		redis.POST("", module.Redis.Ctl.Create)
		redis.GET("/:id", module.Redis.Ctl.Get)

	}
}
