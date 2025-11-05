package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/models"
	"github.com/mygallery/mygallery/internal/storage"
	"github.com/mygallery/mygallery/internal/utils"
)

type PhotoHandler struct {
	cfg *config.Config
}

func NewPhotoHandler(cfg *config.Config) *PhotoHandler {
	return &PhotoHandler{cfg: cfg}
}

// GetPhotos 获取照片列表
func (h *PhotoHandler) GetPhotos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	category := c.Query("category")
	
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	
	offset := (page - 1) * size
	
	var photos []models.Photo
	var total int64
	
	db := database.GetDB()
	query := db.Model(&models.Photo{})
	
	// 按分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	query.Count(&total)
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&photos)
	
	// 为每个照片生成 URL
	stor := storage.GetStorage()
	for i := range photos {
		photos[i].URL = stor.GetURL(photos[i].StoragePath)
		if photos[i].ThumbnailPath != "" {
			// 缩略图路径需要加上 thumbnails/ 前缀
			photos[i].ThumbnailURL = stor.GetURL("thumbnails/" + photos[i].ThumbnailPath)
		}
	}
	
	c.JSON(http.StatusOK, models.PhotoListResponse{
		Total:  total,
		Page:   page,
		Size:   size,
		Photos: photos,
	})
}

// GetPhoto 获取单张照片
func (h *PhotoHandler) GetPhoto(c *gin.Context) {
	id := c.Param("id")
	
	var photo models.Photo
	if err := database.GetDB().First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "照片不存在"})
		return
	}
	
	// 生成 URL
	stor := storage.GetStorage()
	photo.URL = stor.GetURL(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		photo.ThumbnailURL = stor.GetURL("thumbnails/" + photo.ThumbnailPath)
	}
	
	// 增加浏览次数
	database.GetDB().Model(&photo).UpdateColumn("views", photo.Views+1)
	
	c.JSON(http.StatusOK, photo)
}

// UploadPhoto 上传照片
func (h *PhotoHandler) UploadPhoto(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// 获取上传的文件
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "获取上传文件失败",
			Message: err.Error(),
		})
		return
	}
	defer file.Close()
	
	// 检查文件类型
	if !utils.IsAllowedImageType(header.Header.Get("Content-Type"), h.cfg.Image.AllowedTypes) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "不支持的文件类型",
		})
		return
	}
	
	// 检查文件大小
	fileSize, err := utils.GetFileSize(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "获取文件大小失败"})
		return
	}
	
	if fileSize > h.cfg.Image.MaxUploadSize {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: fmt.Sprintf("文件大小超过限制 (%dMB)", h.cfg.Image.MaxUploadSize/1024/1024),
		})
		return
	}
	
	// 生成文件名
	filename := utils.GenerateFilename(header.Filename)
	
	// 保存到临时目录（用于EXIF提取）
	tempDir := filepath.Join(os.TempDir(), "mygallery")
	os.MkdirAll(tempDir, 0755)
	tempPath := filepath.Join(tempDir, filename)
	
	tempFile, err := os.Create(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "创建临时文件失败"})
		return
	}
	defer os.Remove(tempPath)
	defer tempFile.Close()
	
	file.Seek(0, 0)
	if _, err := tempFile.ReadFrom(file); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "保存临时文件失败"})
		return
	}
	tempFile.Close()
	
	// 提取 EXIF 信息
	exifData, _ := utils.ExtractEXIF(tempPath)
	
	// 生成缩略图
	thumbnailFilename := "thumb_" + filename
	var thumbnailPath string
	
	if localStor, ok := storage.GetStorage().(*storage.LocalStorage); ok {
		thumbnailPath = filepath.Join(localStor.GetThumbnailDir(), thumbnailFilename)
		utils.GenerateThumbnail(
			tempPath,
			thumbnailPath,
			h.cfg.Image.Thumbnail.Width,
			h.cfg.Image.Thumbnail.Height,
			h.cfg.Image.Thumbnail.Quality,
		)
	}
	
	// 上传到存储
	file.Seek(0, 0)
	stor := storage.GetStorage()
	storagePath, err := stor.Upload(filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "上传文件失败",
			Message: err.Error(),
		})
		return
	}
	
	// 创建照片记录
	photo := models.Photo{
		Filename:      filename,
		OriginalName:  header.Filename,
		Title:         c.PostForm("title"),
		Description:   c.PostForm("description"),
		Tags:          c.PostForm("tags"),
		Location:      c.PostForm("location"),
		StorageType:   h.cfg.Storage.Type,
		StoragePath:   storagePath,
		ThumbnailPath: thumbnailFilename,
		FileSize:      fileSize,
		MimeType:      header.Header.Get("Content-Type"),
		Copyright:     c.PostForm("copyright"),
		UserID:        userID,
	}
	
	// 设置 EXIF 数据
	if exifData != nil {
		photo.CameraMake = exifData.CameraMake
		photo.CameraModel = exifData.CameraModel
		photo.LensModel = exifData.LensModel
		photo.FocalLength = exifData.FocalLength
		photo.Aperture = exifData.Aperture
		photo.ShutterSpeed = exifData.ShutterSpeed
		photo.ISO = exifData.ISO
		photo.DateTaken = exifData.DateTaken
		photo.GPSLatitude = exifData.GPSLatitude
		photo.GPSLongitude = exifData.GPSLongitude
		photo.Width = exifData.Width
		photo.Height = exifData.Height
	}
	
	// 保存到数据库
	if err := database.GetDB().Create(&photo).Error; err != nil {
		// 删除已上传的文件
		stor.Delete(storagePath)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "保存照片信息失败",
			Message: err.Error(),
		})
		return
	}
	
	// 生成 URL
	photo.URL = stor.GetURL(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		photo.ThumbnailURL = stor.GetURL(photo.ThumbnailPath)
	}
	
	c.JSON(http.StatusOK, models.UploadResponse{
		Success: true,
		Message: "照片上传成功",
		Photo:   photo,
	})
}

// UpdatePhoto 更新照片信息
func (h *PhotoHandler) UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	
	var photo models.Photo
	if err := database.GetDB().First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "照片不存在"})
		return
	}
	
	// 更新字段
	updates := make(map[string]interface{})
	
	if title := c.PostForm("title"); title != "" {
		updates["title"] = title
	}
	if description := c.PostForm("description"); description != "" {
		updates["description"] = description
	}
	if tags := c.PostForm("tags"); tags != "" {
		updates["tags"] = tags
	}
	if location := c.PostForm("location"); location != "" {
		updates["location"] = location
	}
	if copyright := c.PostForm("copyright"); copyright != "" {
		updates["copyright"] = copyright
	}
	if category := c.PostForm("category"); category != "" {
		updates["category"] = category
	}
	
	if err := database.GetDB().Model(&photo).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "更新照片信息失败",
			Message: err.Error(),
		})
		return
	}
	
	// 重新获取更新后的数据
	database.GetDB().First(&photo, id)
	
	// 生成 URL
	stor := storage.GetStorage()
	photo.URL = stor.GetURL(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		photo.ThumbnailURL = stor.GetURL(photo.ThumbnailPath)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "照片信息更新成功",
		"photo":   photo,
	})
}

// DeletePhoto 删除照片
func (h *PhotoHandler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	
	var photo models.Photo
	if err := database.GetDB().First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "照片不存在"})
		return
	}
	
	// 删除存储的文件
	stor := storage.GetStorage()
	stor.Delete(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		stor.Delete(photo.ThumbnailPath)
	}
	
	// 从数据库删除
	if err := database.GetDB().Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "删除照片失败",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "照片删除成功",
	})
}

