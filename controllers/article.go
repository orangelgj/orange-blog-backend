package controllers

import (
	"fmt"
	"gblog/models"
	"gblog/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Author     string `json:"author" binding:"required"`
	CategoryID uint32 `json:"category_id" binding:"required"`
	Summary    string `json:"summary"`
	Content    string `json:"content" binding:"required"`
}

// CreateArticle godoc
// @Summary 创建文章
// @Description 创建新文章接口
// @Tags 文章
// @Accept json
// @Produce json
// @Param request body CreateArticleRequest true "文章信息"
// @Success 200 {object} object{code=int,msg=string} "创建成功"
// @Failure 400 {object} object{code=int,msg=string} "参数错误"
// @Failure 401 {object} object{code=int,msg=string} "未认证"
// @Failure 403 {object} object{code=int,msg=string} "权限不足"
// @Failure 500 {object} object{code=int,msg=string} "创建失败"
// @Router /api/v1/articles [post]
func CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	article := models.Article{
		Title:      req.Title,
		Author:     req.Author,
		CategoryID: req.CategoryID,
		Summary:    req.Summary,
		Content:    req.Content,
	}

	if err := models.CreateArticle(&article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "文章创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "文章创建成功",
	})
	users, err := models.GetUsers()
	if err != nil {
		//报错说获取邮箱失败无法发送邮件，结构化
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("获取用户邮箱失败无法发送邮件")
		return
	}
	for _, user := range users {
		if user.Role == 1 {
			continue
		}
		err := utils.SendArticleEmail(user.Email, article.Title, strconv.Itoa(int(article.ID)), user.Username)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"error": err,
				"email": user.Email,
			}).Error("发送邮件失败")
			fmt.Printf("发送邮件失败: %s\n", user.Email)
			continue
		}
	}
}
