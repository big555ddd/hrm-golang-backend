package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func notification(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	notification := router.Group("", md, log)
	{

		notification.GET("/ws", module.Notification.Ctl.Connect)
		notification.PATCH("/read/:id", module.Notification.Ctl.MarkAsRead)
		notification.PATCH("/read-all", module.Notification.Ctl.MarkAllAsRead)
		notification.GET("", module.Notification.Ctl.List)

	}
}
