package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/models"
)

type SettingsHandler struct {
	cfg *config.Config
}

func NewSettingsHandler(cfg *config.Config) *SettingsHandler {
	return &SettingsHandler{cfg: cfg}
}

// GetSettings 获取网站设置
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	var settings models.Settings
	
	// 获取第一条设置记录，如果不存在则创建
	if err := database.GetDB().First(&settings).Error; err != nil {
		settings = models.Settings{
			SiteTitle:       "MYGallery",
			SiteDescription: "个人照片展示与管理系统",
		}
		database.GetDB().Create(&settings)
	}
	
	c.JSON(http.StatusOK, settings)
}

// UpdateSettings 更新网站设置
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var settings models.Settings
	
	// 获取现有设置
	database.GetDB().First(&settings)
	
	// 绑定更新数据
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "参数错误",
			Message: err.Error(),
		})
		return
	}
	
	// 更新
	if err := database.GetDB().Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "更新失败",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "设置更新成功",
		"settings": settings,
	})
}

// GetCategories 获取分类列表
func (h *SettingsHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	database.GetDB().Order("sort_order ASC, created_at DESC").Find(&categories)
	
	c.JSON(http.StatusOK, categories)
}

// CreateCategory 创建分类
func (h *SettingsHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "参数错误",
			Message: err.Error(),
		})
		return
	}
	
	if err := database.GetDB().Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "创建失败",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "分类创建成功",
		"category": category,
	})
}

// UpdateCategory 更新分类
func (h *SettingsHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	
	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "分类不存在"})
		return
	}
	
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "参数错误",
			Message: err.Error(),
		})
		return
	}
	
	if err := database.GetDB().Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "更新失败",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "分类更新成功",
		"category": category,
	})
}

// DeleteCategory 删除分类
func (h *SettingsHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	
	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "分类不存在"})
		return
	}
	
	if err := database.GetDB().Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "删除失败",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分类删除成功",
	})
}

