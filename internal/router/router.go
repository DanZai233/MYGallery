package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/handlers"
	"github.com/mygallery/mygallery/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)
	
	r := gin.Default()
	
	// CORS 中间件
	r.Use(middleware.CORSMiddleware(cfg))
	
	// 静态文件
	r.Static("/uploads", cfg.Storage.Local.UploadDir)
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/admin", "./public/admin.html")
	r.StaticFile("/settings.html", "./public/settings.html")
	r.StaticFile("/categories.html", "./public/categories.html")
	r.StaticFile("/map.html", "./public/map.html")
	r.StaticFile("/test", "./test-frontend.html")
	r.Static("/assets", "./public/assets")
	
	// 创建 handlers
	authHandler := handlers.NewAuthHandler(cfg)
	photoHandler := handlers.NewPhotoHandler(cfg)
	settingsHandler := handlers.NewSettingsHandler(cfg)
	reactionHandler := handlers.NewReactionHandler()
	
	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关（无需 token）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}
		
		// 照片相关（公开接口）
		photos := api.Group("/photos")
		{
			photos.GET("", photoHandler.GetPhotos)       // 获取照片列表
			photos.GET("/:id", photoHandler.GetPhoto)    // 获取单张照片
		}
		
		// 设置相关（公开接口）
		api.GET("/settings", settingsHandler.GetSettings)
		api.GET("/categories", settingsHandler.GetCategories)

		// 表态相关（公开接口，通过指纹识别用户）
		api.GET("/photos/:id/reactions", reactionHandler.GetReactions)
		api.POST("/photos/:id/reactions", reactionHandler.AddReaction)
		api.DELETE("/photos/:id/reactions", reactionHandler.DeleteReaction)
		
		// 需要认证的接口
		authRequired := api.Group("")
		authRequired.Use(middleware.AuthMiddleware(cfg))
		{
			// 用户相关
			authRequired.POST("/auth/change-password", authHandler.ChangePassword)
			
			// 照片管理
			authRequired.POST("/photos", photoHandler.UploadPhoto)
			authRequired.PUT("/photos/:id", photoHandler.UpdatePhoto)
			authRequired.DELETE("/photos/:id", photoHandler.DeletePhoto)
			
			// 设置管理
			authRequired.PUT("/settings", settingsHandler.UpdateSettings)
			
			// 分类管理
			authRequired.POST("/categories", settingsHandler.CreateCategory)
			authRequired.PUT("/categories/:id", settingsHandler.UpdateCategory)
			authRequired.DELETE("/categories/:id", settingsHandler.DeleteCategory)
		}
	}
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})
	
	return r
}

