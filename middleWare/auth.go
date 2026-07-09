package middleware

import (
	"fmt"
	"memo/common"
	"memo/controller"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		// 检查tokenString格式是否正确
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			fmt.Println("token格式不正确")
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		// 提取tokenString有效部分（去除前7位，因为Bearer 占7位）
		tokenString = tokenString[7:]
		// 解析tokenString
		token, claim, err := common.ParseToken(tokenString)
		fmt.Println("当前token对应的用户id是：", claim.UserId)
		fmt.Println("当前token是：", tokenString)
		if err != nil || !token.Valid {
			fmt.Println("错误是:", err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		// 根据token拿到用户信息
		userId := claim.UserId
		userDB := controller.GetUserDB()
		var tmp_user common.UserInfo
		userDB.First(&tmp_user, userId)

		// 用户不存在
		if tmp_user.Name == "" {
			fmt.Println("用户不存在")
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		// 用户存在，将用户信息写入上下文，方便使用
		ctx.Set("user", tmp_user)

		// 继续中间件
		ctx.Next()
	}
}
