package middleware

import (
	"goworks/common"
	"goworks/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401, "msg": "权限不足1",
			})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401, "msg": "权限不足2",
			})
			ctx.Abort()
			return
		}
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)
		if userId == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401, "msg": "权限不足3",
			})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}

}
