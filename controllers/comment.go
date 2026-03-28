package controllers

import (
	"encoding/json"
	"gblog/dto"
	"gblog/models"
	"gblog/utils"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetRootCommentList(c *gin.Context) {
	articleIDStr := c.Query("articleId")
	if articleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "缺少 articleId 参数",
		})
		return
	}

	articleID, err := strconv.ParseUint(articleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "articleId 格式错误",
		})
		return
	}

	pageNum := 1
	pageSize := 10
	if pn := c.Query("pageNum"); pn != "" {
		if v, err := strconv.Atoi(pn); err == nil && v > 0 {
			pageNum = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	comments, total, err := models.GetRootCommentsByArticleID(uint(articleID), pageSize, pageNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询评论失败",
		})
		return
	}

	userCache := make(map[uint]string)
	rootIDs := make([]uint, 0, len(comments))
	for _, comment := range comments {
		rootIDs = append(rootIDs, comment.ID)
	}

	firstReplies, _ := models.GetFirstRepliesByRootIDs(rootIDs)

	rootCommentDTOs := make([]dto.RootCommentDTO, 0, len(comments))

	for _, comment := range comments {
		userName, ok := userCache[comment.UserID]
		if !ok {
			user, err := models.GetUserByID(comment.UserID)
			if err == nil {
				userName = user.Username
				userCache[comment.UserID] = userName
			}
		}

		replyCount, _ := models.GetReplyCountByRootID(comment.ID)

		var previewReply *dto.ReplyDTO
		if reply, exists := firstReplies[comment.ID]; exists {
			replyUserName, ok := userCache[reply.UserID]
			if !ok {
				user, err := models.GetUserByID(reply.UserID)
				if err == nil {
					replyUserName = user.Username
					userCache[reply.UserID] = replyUserName
				}
			}
			replyContent := reply.Content
			if reply.Status == -1 {
				replyContent = "已删除"
			} else if reply.Status == 0 {
				replyContent = "审核中"
			}
			previewReply = &dto.ReplyDTO{
				ID:         reply.ID,
				Content:    replyContent,
				UserID:     reply.UserID,
				UserName:   replyUserName,
				ToUserID:   reply.ToUserId,
				ToUserName: reply.ToUserName,
				CreateTime: reply.CreateTime,
				LikeCount:  reply.LikeCount,
				IpLocation: reply.IpLocation,
			}
		}

		content := comment.Content
		if comment.Status == -1 {
			content = "已删除"
		} else if comment.Status == 0 {
			content = "审核中"
		}

		rootCommentDTOs = append(rootCommentDTOs, dto.RootCommentDTO{
			ID:           comment.ID,
			Content:      content,
			UserID:       comment.UserID,
			UserName:     userName,
			CreateTime:   comment.CreateTime,
			LikeCount:    comment.LikeCount,
			ReplyCount:   replyCount,
			IpLocation:   comment.IpLocation,
			PreviewReply: previewReply,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "获取成功",
		"data":  rootCommentDTOs,
		"total": total,
	})
}

func GetChildCommentList(c *gin.Context) {
	rootIDStr := c.Query("rootId")
	if rootIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "缺少 rootId 参数",
		})
		return
	}

	rootID, err := strconv.ParseUint(rootIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "rootId 格式错误",
		})
		return
	}

	pageNum := 1
	pageSize := 10
	if pn := c.Query("pageNum"); pn != "" {
		if v, err := strconv.Atoi(pn); err == nil && v > 0 {
			pageNum = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	comments, total, err := models.GetChildCommentsByRootID(uint(rootID), pageSize, pageNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询回复失败",
		})
		return
	}

	userCache := make(map[uint]string)
	replyDTOs := make([]dto.ReplyDTO, 0, len(comments))

	for _, comment := range comments {
		userName, ok := userCache[comment.UserID]
		if !ok {
			user, err := models.GetUserByID(comment.UserID)
			if err == nil {
				userName = user.Username
				userCache[comment.UserID] = userName
			}
		}

		content := comment.Content
		if comment.Status == -1 {
			content = "已删除"
		} else if comment.Status == 0 {
			content = "审核中"
		}

		replyDTOs = append(replyDTOs, dto.ReplyDTO{
			ID:         comment.ID,
			Content:    content,
			UserID:     comment.UserID,
			UserName:   userName,
			ToUserID:   comment.ToUserId,
			ToUserName: comment.ToUserName,
			CreateTime: comment.CreateTime,
			LikeCount:  comment.LikeCount,
			IpLocation: comment.IpLocation,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "获取成功",
		"data":  replyDTOs,
		"total": total,
	})
}

type CreateCommentRequest struct {
	ArticleID uint   `json:"articleId" binding:"required"`
	Content   string `json:"content" binding:"required,max=200"`
	RootId    uint   `json:"rootId"`
	ParentId  uint   `json:"parentId"`
	ToUserId  uint   `json:"toUserId"`
}

func CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误：" + err.Error(),
		})
		return
	}

	content := strings.TrimSpace(req.Content)
	if len(content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "评论内容不能为空",
		})
		return
	}

	if len(content) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "评论内容不能超过200字",
		})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户未登录",
		})
		return
	}
	user := userInterface.(models.User)

	ip := c.ClientIP()
	ipLocation := getIpLocation(ip)

	var toUserName string
	if req.ToUserId > 0 {
		toUser, err := models.GetUserByID(req.ToUserId)
		if err == nil {
			toUserName = toUser.Username
		}
	}

	comment := &models.Comment{
		ArticleID:  req.ArticleID,
		UserID:     uint(user.ID),
		RootId:     req.RootId,
		ParentId:   req.ParentId,
		ToUserId:   req.ToUserId,
		ToUserName: toUserName,
		Content:    content,
		Status:     1,
		IpLocation: ipLocation,
	}

	if err := models.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "评论发表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "评论发表成功",
		"data": dto.ReplyDTO{
			ID:         comment.ID,
			Content:    comment.Content,
			UserID:     comment.UserID,
			UserName:   user.Username,
			ToUserID:   comment.ToUserId,
			ToUserName: comment.ToUserName,
			CreateTime: comment.CreateTime,
			LikeCount:  comment.LikeCount,
			IpLocation: comment.IpLocation,
		},
	})

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    user.ID,
		"username":   user.Username,
		"comment_id": comment.ID,
		"article_id": req.ArticleID,
		"ip":         c.ClientIP(),
		"action":     "create_comment",
	}).Info("Comment created")
}

func getIpLocation(ip string) string {
	ipStr := ip
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return "未知" // 或者处理空字符串
	}

	// IsLoopback: 检查是否是 127.0.0.1 或 ::1
	// IsPrivate: 检查是否是 RFC 1918 (10/8, 172.16/12, 192.168/16) 或 RFC 4193 (IPv6 ULA)
	// IsLinkLocalUnicast: 检查是否是 169.254.x.x
	if parsedIP.IsLoopback() || parsedIP.IsPrivate() || parsedIP.IsLinkLocalUnicast() {
		return "本地"
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + url.QueryEscape(ip) + "?lang=zh-CN")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Country string `json:"country"`
		City    string `json:"city"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	if result.Status != "success" {
		return ""
	}

	location := result.Country
	if result.City != "" && result.City != result.Country {
		location += " " + result.City
	}

	return location
}

func DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	if commentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "缺少评论ID",
		})
		return
	}

	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "评论ID格式错误",
		})
		return
	}

	comment, err := models.GetCommentByID(uint(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评论不存在",
		})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户未登录",
		})
		return
	}
	user := userInterface.(models.User)

	roleInterface, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未认证",
		})
		return
	}

	role, ok := roleInterface.(int8)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "权限类型错误",
		})
		return
	}

	if role != 1 && uint(user.ID) != comment.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "权限不足",
		})
		return
	}

	if err := models.DeleteComment(uint(commentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})

	utils.Logger.WithFields(logrus.Fields{
		"user_id":    user.ID,
		"username":   user.Username,
		"comment_id": commentID,
		"ip":         c.ClientIP(),
		"action":     "delete_comment",
	}).Info("Comment deleted")
}
