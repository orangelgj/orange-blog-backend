package controllers

import (
	"gblog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategories godoc
// @Summary 获取分类列表
// @Description 获取所有分类列表
// @Tags 分类
// @Accept json
// @Produce json
// @Success 200 {object} object{code=int,msg=string,data=[]object} "获取成功"
// @Failure 500 {object} object{code=int,msg=string,data=nil} "获取失败"
// @Router /api/v1/categories [get]
func GetCategories(c *gin.Context) {
	categories, err := models.GetCategories()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"data": nil,
			"msg":  "获取失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
		"msg":  "获取成功",
	})
}
