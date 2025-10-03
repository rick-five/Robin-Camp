package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"robin-camp/handlers"
	"robin-camp/models"
	"robin-camp/utils"
)

func main() {
	// 连接数据库
	db, err := utils.ConnectDatabase()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// 自动迁移数据库模式
	err = db.AutoMigrate(&models.Movie{}, &models.Rating{})
	if err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由器
	router := gin.New()
	router.Use(gin.Recovery())

	// 注册路由
	registerRoutes(router, db)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	fmt.Printf("Server starting on port %s\n", port)
	err = router.Run("0.0.0.0:" + port)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func registerRoutes(router *gin.Engine, db *gorm.DB) {
	// 健康检查端点
	router.GET("/healthz", func(c *gin.Context) {
		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Database connection error"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Database ping failed"})
			return
		}

		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建处理器
	movieHandler := handlers.NewMovieHandler(db)
	ratingHandler := handlers.NewRatingHandler(db)

	// 电影路由
	movies := router.Group("/movies")
	{
		// 需要认证的路由组
		authMovies := movies.Group("/")
		authMovies.Use(authMiddleware())
		authMovies.POST("/", movieHandler.CreateMovie)

		// 公开路由
		movies.GET("/", movieHandler.ListMovies)
		
		// 通过ID获取单个电影
		movies.GET("/:id", movieHandler.GetMovie)
		
		// 评分路由 - 使用不同的路径前缀避免与ID路由冲突
		movies.POST("/title/:title/ratings", ratingHandler.SubmitRating)
		movies.GET("/title/:title/rating", ratingHandler.GetRatingAggregate)
	}
}

// 认证中间件 - 用于电影创建
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := os.Getenv("AUTH_TOKEN")
		if authToken == "" {
			c.JSON(500, gin.H{"error": "AUTH_TOKEN not configured"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		expectedHeader := "Bearer " + authToken

		if authHeader != expectedHeader {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 评分认证中间件 - 检查X-Rater-Id头
func raterAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raterID := c.GetHeader("X-Rater-Id")
		if raterID == "" {
			c.JSON(401, gin.H{"error": "Missing X-Rater-Id header"})
			c.Abort()
			return
		}

		c.Next()
	}
}
