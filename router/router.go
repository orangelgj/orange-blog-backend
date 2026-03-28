package router

import (
	"gblog/controllers"
	"gblog/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.GlobalRateLimiter())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/search", controllers.Search)
		v1.GET("/articles", controllers.GetArticles)
		v1.GET("/article/:id", controllers.GetArticleDetail)
		v1.GET("/categories", controllers.GetCategories)

		v1.GET("/check", controllers.Check)
		v1.POST("/login", middlewares.StrictRateLimiter(), controllers.Login)
		v1.POST("/register", middlewares.StrictRateLimiter(), controllers.Register)
	}

	auth := v1.Group("")
	auth.Use(middlewares.AuthMiddleWare())
	{
		auth.GET("/comments", controllers.GetRootCommentList)
		auth.GET("/comments/replies", controllers.GetChildCommentList)
		auth.POST("/comments", controllers.CreateComment)
		auth.DELETE("/comments/:id", controllers.DeleteComment)
		auth.POST("/user/username", controllers.UpdateUsername)
		auth.POST("/user/password", controllers.UpdatePassword)
	}

	admin := v1.Group("")
	admin.Use(middlewares.AuthMiddleWare(), middlewares.AdminMiddleWare())
	{
		admin.POST("/articles", controllers.CreateArticle)
	}

	return r
}
