package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func auth(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	auth := router.Group("")
	{
		auth.POST("/login", module.Auth.Ctl.Login)
		auth.GET("/info", md, module.Auth.Ctl.Info)
		auth.POST("/forgot-password", module.Auth.Ctl.ForgotPassword)
		auth.POST("/reset-password", module.Auth.Ctl.ResetPassword)
	}
}
