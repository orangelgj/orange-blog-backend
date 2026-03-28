package middlewares

import (
	"gblog/models"
	"gblog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "错误",
			})
			c.Abort() // 重点：停止执行后续的 Handler
			return
		}
		// 2. 截取 Token 字符串
		tokenString := authHeader[7:]

		// 3. 调用我们之前写的 utils.ParseToken 解析 Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的 Token 或登录已过期",
			})
			c.Abort()
			return
		}
		user, err := models.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户不存在",
			})
			c.Abort()
			return
		}
		if user.UpdateTime.Unix() != claims.PasswordUpdateTime {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "密码已修改，请重新登录",
			})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func AdminMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未认证",
			})
			c.Abort()
			return
		}

		roleInt, ok := role.(int8)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限类型错误",
			})
			c.Abort()
			return
		}

		if roleInt != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
