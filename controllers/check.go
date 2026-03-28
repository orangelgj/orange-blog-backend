package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Check godoc
// @Summary 检查认证状态
// @Description 检查用户认证状态，验证 cookie 是否有效
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} object{auth=bool,message=string} "认证成功"
// @Failure 200 {object} object{auth=bool,message=string} "认证失败"
// @Router /api/v1/check [get]
func Check(c *gin.Context) {
	// 1. 获取名为 "暗号" 的 cookie
	cookieValue, err := c.Cookie("暗号")

	// 2. 判断是否存在且值是否等于 "暗号"
	if err != nil || cookieValue != "暗号" {
		// 校验失败：返回 401 或自定义逻辑
		c.JSON(200, gin.H{
			"auth":    false,
			"message": "嘻嘻嘻啦啦啦",
		})
		return
	}

	// 3. 校验成功
	c.JSON(http.StatusOK, gin.H{
		"auth":    true,
		"message": "成功啦",
	})
}
