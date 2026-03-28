package controllers

import (
	"gblog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetArticleDetail godoc
// @Summary 获取文章详情
// @Description 根据文章 ID 获取文章详细信息
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Success 200 {object} object{code=int,msg=string,data=object} "查询成功"
// @Failure 400 {object} object{code=int,msg=string} "ID 格式错误"
// @Failure 500 {object} object{code=int,msg=string} "查询失败或文章不存在"
// @Router /api/v1/article/{id} [get]
func GetArticleDetail(c *gin.Context) {
	// 1. 获取路径参数 id (字符串)
	idStr := c.Param("id")

	// 2. 将字符串转换为数字 (uint32)
	// Atoi 是转为 int，如果是 uint32 建议用 ParseUint
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// 如果转换失败（比如输入了非数字），返回 400
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID 格式错误，请输入数字",
		})
		return
	}

	// 3. 转换为 uint32 并调用模型层方法
	id := uint32(idUint64)
	article, err := models.GetArticleDetail(id)
	if err != nil {
		// 处理数据库查询错误或文章不存在
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败或文章不存在",
		})
		return
	}

	// 4. 返回成功数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": article,
		"msg":  "查询成功",
	})
}
