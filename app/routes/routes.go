// app/routes/routes.go
package routes

import (
	"net/http"

	"app/internal/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router sets up all the routes for the application
func Router(app *gin.Engine) {

	// Health check endpoint
	app.GET("/healthz", func(ctx *gin.Context) {
		logger.Infof("Health check passed")
		ctx.JSON(http.StatusOK, gin.H{"status": "Health check passed.", "message": "Welcome to Project-k API."})
	})

	// Middleware
	app.Use(otelgin.Middleware(viper.GetString("APP_NAME")))
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"*"},
		AllowHeaders:           []string{"*"},
		AllowCredentials:       true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             false,
	}))

	// Create a new group for /api/v1
	apiV1 := app.Group("/api/v1")

	// Define groups of routes under /api/v1
	auth(apiV1.Group("/auth"))
	redis(apiV1.Group("/redis"))
	user(apiV1.Group("/users"))
	role(apiV1.Group("/roles"))
	organization(apiV1.Group("/organizations"))
	branch(apiV1.Group("/branches"))
	department(apiV1.Group("/departments"))
	holiday(apiV1.Group("/holidays"))
	workshift(apiV1.Group("/workshifts"))
	leave(apiV1.Group("/leaves"))
	document(apiV1.Group("/documents"))
	attendance(apiV1.Group("/attendances"))
	notification(apiV1.Group("/notifications"))

}
