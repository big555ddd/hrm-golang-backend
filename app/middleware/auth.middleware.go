package middleware

import (
	"app/app/response"
	"app/app/util/jwt"
	"app/internal/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ""
		logger.Info(ctx.Request.URL.Path)
		if ctx.Request.URL.Path == "/api/v1/notifications/ws" {
			authHeader = ctx.Query("token")
		} else {
			authHeader = ctx.GetHeader("Authorization")
		}
		if authHeader == "" {
			response.Unauthorized(ctx, "Authorization header is required", nil)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(ctx, "Authorization header format must be Bearer {token}", nil)
			ctx.Abort()
			return
		}

		token := parts[1]
		claims, _, err := jwt.Verify(token)
		if err != nil {
			response.Unauthorized(ctx, err.Error(), nil)
			ctx.Abort()
			return
		}

		// Optionally, you can set the claims in the context for further use in your handlers
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
