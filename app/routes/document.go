package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func document(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	document := router.Group("", md, log)
	{
		document.POST("", module.Document.Ctl.Create)
		document.GET("", module.Document.Ctl.List)
		document.GET("/:id", module.Document.Ctl.Get)
		document.DELETE("/:id", module.Document.Ctl.Delete)
		document.PATCH("/:id/approve", module.Document.Ctl.Approved)
		document.PATCH("/:id/reject", module.Document.Ctl.Rejected)

	}
}
