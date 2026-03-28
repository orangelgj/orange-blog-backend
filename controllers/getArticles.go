package controllers

import (
	"gblog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetArticles godoc
// @Summary 获取文章列表
// @Description 获取文章列表接口，返回文章列表，支持按分类筛选
// @Tags 文章
// @Accept json
// @Produce json
// @Param categoryId query int false "分类ID，0表示全部"
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认10"
// @Success 200 {object} object{code=int,msg=string,data=[]object} "获取成功"
// @Failure 500 {object} object{code=int,msg=string,data=nil} "获取失败"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	// 获取查询参数
	categoryIDStr := c.DefaultQuery("categoryId", "0")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// 转换参数
	categoryID, _ := strconv.Atoi(categoryIDStr)
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 模拟从数据库查出来的文章列表
	articles, err := models.GetArticleList(pageSize, page, categoryID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"data": nil,
			"msg":  "获取失败,数据库查询问题",
		})
		return
	}
	// 2. 返回 JSON 格式
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": articles,
		"msg":  "获取成功",
	})
}
