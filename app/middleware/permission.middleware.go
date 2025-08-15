package middleware

import (
	"app/app/helper"
	"app/app/modules"
	"app/app/response"
	"slices"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := helper.GetUserByToken(ctx)
		if err != nil {
			response.Unauthorized(ctx, "unauthorized", nil)
			ctx.Abort()
			return
		}
		permissionsdata, err := modules.New().User.Svc.GetUserPermissionsName(ctx, user.Data.ID)
		if err != nil {
			response.Unauthorized(ctx, "unauthorized", err)
			ctx.Abort()
			return
		}
		//check permission include permissionsdata
		if !slices.Contains(permissionsdata, permission) {
			response.Forbidden(ctx, "forbidden", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
