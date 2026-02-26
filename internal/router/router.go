package router

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/handlers"
	"github.com/mygallery/mygallery/internal/middleware"
	"github.com/mygallery/mygallery/internal/models"
	"github.com/mygallery/mygallery/internal/storage"
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
	r.GET("/", func(c *gin.Context) {
		photoID := c.Query("photo")
		if photoID == "" {
			c.File("./public/index.html")
			return
		}
		id, err := strconv.Atoi(photoID)
		if err != nil {
			c.File("./public/index.html")
			return
		}
		var photo models.Photo
		if err := database.GetDB().First(&photo, id).Error; err != nil {
			c.File("./public/index.html")
			return
		}
		imgURL := storage.GetStorage().GetURL(photo.StoragePath)
		if photo.ThumbnailPath != "" {
			imgURL = storage.GetStorage().GetURL("thumbnails/" + photo.ThumbnailPath)
		}
		title := photo.Title
		if title == "" {
			title = photo.OriginalName
		}
		desc := photo.Description
		if desc == "" {
			desc = "MYGallery - 个人照片墙"
		}
		htmlBytes, _ := os.ReadFile("./public/index.html")
		html := string(htmlBytes)
		ogTags := fmt.Sprintf(
			`<meta property="og:title" content="%s">`+
				`<meta property="og:description" content="%s">`+
				`<meta property="og:image" content="%s">`+
				`<meta property="og:type" content="article">`+
				`<meta name="twitter:card" content="summary_large_image">`+
				`<meta name="twitter:title" content="%s">`+
				`<meta name="twitter:description" content="%s">`+
				`<meta name="twitter:image" content="%s">`,
			title, desc, imgURL, title, desc, imgURL,
		)
		html = strings.Replace(html, "</head>", ogTags+"</head>", 1)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
	r.StaticFile("/admin", "./public/admin.html")
	r.StaticFile("/settings.html", "./public/settings.html")
	r.StaticFile("/categories.html", "./public/categories.html")
	r.StaticFile("/map.html", "./public/map.html")
	r.StaticFile("/albums.html", "./public/albums.html")
	r.StaticFile("/test", "./test-frontend.html")
	r.Static("/assets", "./public/assets")
	
	// 创建 handlers
	authHandler := handlers.NewAuthHandler(cfg)
	photoHandler := handlers.NewPhotoHandler(cfg)
	settingsHandler := handlers.NewSettingsHandler(cfg)
	reactionHandler := handlers.NewReactionHandler()
	albumHandler := handlers.NewAlbumHandler()
	
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

		// 表态相关（公开接口）
		api.GET("/photos/:id/reactions", reactionHandler.GetReactions)
		api.POST("/photos/:id/reactions", reactionHandler.AddReaction)
		api.DELETE("/photos/:id/reactions", reactionHandler.DeleteReaction)

		// 相册（公开读取）
		api.GET("/albums", albumHandler.GetAlbums)
		api.GET("/albums/:id", albumHandler.GetAlbum)
		
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

			// 相册管理
			authRequired.POST("/albums", albumHandler.CreateAlbum)
			authRequired.PUT("/albums/:id", albumHandler.UpdateAlbum)
			authRequired.DELETE("/albums/:id", albumHandler.DeleteAlbum)
			authRequired.POST("/albums/:id/photos", albumHandler.AddPhotos)
			authRequired.DELETE("/albums/:id/photos/:photoId", albumHandler.RemovePhoto)
			
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

