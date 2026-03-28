package controllers

import (
	"fmt"
	"gblog/models"
	"gblog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码，返回 JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{username=string,password=string} true "登录信息"
// @Success 200 {object} object{code=int,msg=string,data=object{token=string,user=string}} "登录成功"
// @Failure 400 {object} object{code=int,msg=string} "参数格式错误"
// @Failure 401 {object} object{code=int,msg=string} "用户不存在或密码错误"
// @Failure 500 {object} object{code=int,msg=string} "Token 生成失败"
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	var input struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	// 1. 绑定前端传来的 JSON 参数
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数格式错误"})
		return
	}
	user, err := models.GetUserByUsername(input.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}
	// 4. 生成 Token
	// 传入数据库中的用户 ID、角色和密码更新时间戳
	token, err := utils.GenerateToken(user.ID, user.Role, user.UpdateTime.Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Token 生成失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
			"username": user.Username,
			"userId":   user.ID,
			"role":     user.Role,
		},
	})

	utils.Logger.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"ip":       c.ClientIP(),
		"action":   "login",
	}).Info("User logged in")
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口，创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{username=string,password=string} true "注册信息"
// @Success 200 {object} object{code=int,msg=string} "注册成功"
// @Failure 400 {object} object{code=int,msg=string} "参数格式错误"
// @Router /api/v1/register [post]
func Register(c *gin.Context) {
	// 1. 定义带有校验标签的结构体
	var input struct {
		UserName    string `json:"username" binding:"required,min=2,max=20"`
		Password    string `json:"password" binding:"required,min=6"`
		Email       string `json:"email"    binding:"required,email"`
		Description string `json:"description" binding:"max=200"`
	}

	// 2. 参数绑定与校验
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败：" + err.Error()})
		return
	}

	// 3. 权限检查（记得 return！）
	cookieValue, err := c.Cookie("暗号")
	if err != nil || cookieValue != "暗号" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "权限不足，无法注册"})
		return // 必须 return，否则会继续往下执行
	}

	// 4. 【关键】密码哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密失败"})
		return
	}

	user := models.User{
		Username:    input.UserName,
		Password:    string(hashedPassword), // 存储哈希后的密码
		Email:       input.Email,
		Description: input.Description,
		Role:        2,
		// CreateTime 和 UpdateTime 依赖数据库 DEFAULT，这里不需要传
	}

	// 5. 写入数据库
	if err := models.CreateUser(&user); err != nil {
		// 这里建议在 models 层判断是否是唯一键冲突
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "用户名或邮箱已被注册"})
		return
	}

	// 6. 成功返回
	msg := fmt.Sprintf("欢迎你 %s 呀！！！嘻嘻！", input.UserName)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
	})

	utils.Logger.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"ip":       c.ClientIP(),
		"action":   "register",
	}).Info("User registered")
}

// UpdateUsername godoc
// @Summary 更改用户名
// @Description 用户更改用户名接口，需要验证旧密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{old_password=string,new_username=string} true "更改用户名信息"
// @Success 200 {object} object{code=int,msg=string} "更改成功"
// @Failure 400 {object} object{code=int,msg=string} "参数格式错误"
// @Failure 401 {object} object{code=int,msg=string} "旧密码错误"
// @Failure 409 {object} object{code=int,msg=string} "用户名已被占用"
// @Router /api/v1/user/username [post]
func UpdateUsername(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	currentUser := user.(models.User)

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewUsername string `json:"new_username" binding:"required,min=2,max=20"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败：" + err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "旧密码错误"})
		return
	}

	existingUser, err := models.GetUserByUsername(input.NewUsername)
	if err == nil && existingUser.ID != currentUser.ID {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "用户名已被占用"})
		return
	}

	if err := models.UpdateUsername(currentUser.ID, input.NewUsername); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "用户名更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户名修改成功，请重新登录",
	})
}

// UpdatePassword godoc
// @Summary 更改密码
// @Description 用户更改密码接口，需要验证旧密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{old_password=string,new_password=string} true "更改密码信息"
// @Success 200 {object} object{code=int,msg=string} "更改成功"
// @Failure 400 {object} object{code=int,msg=string} "参数格式错误"
// @Failure 401 {object} object{code=int,msg=string} "旧密码错误"
// @Router /api/v1/user/password [post]
func UpdatePassword(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	currentUser := user.(models.User)

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败：" + err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "旧密码错误"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密失败"})
		return
	}

	if err := models.UpdatePassword(currentUser.ID, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码修改成功，请重新登录",
	})
}
