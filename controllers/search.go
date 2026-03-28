package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Search godoc
// @Summary 搜索接口
// @Description 搜索接口，支持搜索功能
// @Tags 搜索
// @Accept json
// @Produce json
// @Param request body object{query=string} true "搜索关键词"
// @Success 200 {object} object{message=string} "搜索成功"
// @Failure 400 {object} object{error=string} "参数错误"
// @Router /api/v1/search [post]
func Search(c *gin.Context) {
	var input struct {
		Q string `json:"query"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Q == "暗号" {
		if input.Q == "暗号" {
			// 设置一个名为 "暗号" 的 Cookie
			c.SetCookie(
				"暗号", // name: Cookie 名称
				"暗号", // value: Cookie 值
				36000000, // maxAge: 有效期（秒），3600 秒即 1 小时；设为 -1 表示即刻删除
				"/",      // path: 路径，"/" 表示全站可用
				"",       // domain: 域名，本地测试用 "localhost" 或空字符串
				false,    // secure: 是否仅通过 HTTPS 传输
				true,     // httpOnly: 是否禁止 JS 读取（建议设为 true 以增强安全性）
			)

			c.JSON(http.StatusOK, gin.H{"message": "成功啦"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "功能尚未开发"})
}
