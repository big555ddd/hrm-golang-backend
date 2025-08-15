package routes

import (
	"app/app/middleware"
	"app/app/modules"

	"github.com/gin-gonic/gin"
)

func attendance(router *gin.RouterGroup) {
	module := modules.New()
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	attendance := router.Group("", md, log)
	{
		attendance.POST("/check-in", module.Attendance.Ctl.CheckIn)
		attendance.GET("", module.Attendance.Ctl.List)
		attendance.GET("/check-in/count", module.Attendance.Ctl.AttendanceCount)

	}
}
