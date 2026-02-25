package handlers

import (
	"bytes"
	"fmt"
	"log"
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
	search := c.Query("search")

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

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"title LIKE ? OR description LIKE ? OR tags LIKE ? OR location LIKE ? OR original_name LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	query.Count(&total)
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&photos)

	stor := storage.GetStorage()
	for i := range photos {
		photos[i].URL = stor.GetURL(photos[i].StoragePath)
		if photos[i].ThumbnailPath != "" {
			photos[i].ThumbnailURL = stor.GetURL("thumbnails/" + photos[i].ThumbnailPath)
		}
		if photos[i].IsLivePhoto && photos[i].LivePhotoPath != "" {
			photos[i].LivePhotoURL = stor.GetURL(photos[i].LivePhotoPath)
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

	stor := storage.GetStorage()
	photo.URL = stor.GetURL(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		photo.ThumbnailURL = stor.GetURL("thumbnails/" + photo.ThumbnailPath)
	}
	if photo.IsLivePhoto && photo.LivePhotoPath != "" {
		photo.LivePhotoURL = stor.GetURL(photo.LivePhotoPath)
	}

	database.GetDB().Model(&photo).UpdateColumn("views", photo.Views+1)

	c.JSON(http.StatusOK, photo)
}

// UploadPhoto 上传照片（支持云存储缩略图 + Live Photo）
func (h *PhotoHandler) UploadPhoto(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "获取上传文件失败",
			Message: err.Error(),
		})
		return
	}
	defer file.Close()

	if !utils.IsAllowedImageType(header.Header.Get("Content-Type"), h.cfg.Image.AllowedTypes) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "不支持的文件类型",
		})
		return
	}

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

	filename := utils.GenerateFilename(header.Filename)

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

	exifData, _ := utils.ExtractEXIF(tempPath)

	thumbnailFilename := "thumb_" + filename
	stor := storage.GetStorage()

	if localStor, ok := stor.(*storage.LocalStorage); ok {
		thumbnailPath := filepath.Join(localStor.GetThumbnailDir(), thumbnailFilename)
		utils.GenerateThumbnail(
			tempPath, thumbnailPath,
			h.cfg.Image.Thumbnail.Width, h.cfg.Image.Thumbnail.Height, h.cfg.Image.Thumbnail.Quality,
		)
	} else {
		thumbBytes, err := utils.GenerateThumbnailBytes(
			tempPath,
			h.cfg.Image.Thumbnail.Width, h.cfg.Image.Thumbnail.Height, h.cfg.Image.Thumbnail.Quality,
		)
		if err == nil {
			stor.UploadThumbnail(thumbnailFilename, bytes.NewReader(thumbBytes))
		} else {
			log.Printf("生成缩略图失败: %v", err)
		}
	}

	file.Seek(0, 0)
	storagePath, err := stor.Upload(filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "上传文件失败",
			Message: err.Error(),
		})
		return
	}

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
		photo.Software = exifData.Software
		photo.Orientation = exifData.Orientation
		photo.WhiteBalance = exifData.WhiteBalance
		photo.Flash = exifData.Flash
		photo.ExposureMode = exifData.ExposureMode
		photo.MeteringMode = exifData.MeteringMode
		photo.ExposureBias = exifData.ExposureBias
		photo.ColorSpace = exifData.ColorSpace
		photo.SceneType = exifData.SceneType
	}

	// Live Photo 配套视频上传
	if liveVideo, liveHeader, err := c.Request.FormFile("live_photo"); err == nil {
		defer liveVideo.Close()
		if utils.IsLivePhotoVideo(liveHeader.Filename) {
			liveFilename := "live_" + utils.GenerateFilename(liveHeader.Filename)
			if livePath, err := stor.Upload(liveFilename, liveVideo); err == nil {
				photo.IsLivePhoto = true
				photo.LivePhotoPath = livePath
			}
		}
	}

	if err := database.GetDB().Create(&photo).Error; err != nil {
		stor.Delete(storagePath)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "保存照片信息失败",
			Message: err.Error(),
		})
		return
	}

	photo.URL = stor.GetURL(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		photo.ThumbnailURL = stor.GetURL("thumbnails/" + photo.ThumbnailPath)
	}
	if photo.IsLivePhoto && photo.LivePhotoPath != "" {
		photo.LivePhotoURL = stor.GetURL(photo.LivePhotoPath)
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

	updates := make(map[string]interface{})

	formFields := map[string]string{
		"title": "title", "description": "description", "tags": "tags",
		"location": "location", "copyright": "copyright", "category": "category",
		"camera_make": "camera_make", "camera_model": "camera_model",
		"lens_model": "lens_model", "focal_length": "focal_length",
		"aperture": "aperture", "shutter_speed": "shutter_speed",
		"iso": "iso", "date_taken": "date_taken",
		"software": "software", "white_balance": "white_balance",
		"flash": "flash", "exposure_mode": "exposure_mode",
		"metering_mode": "metering_mode", "exposure_bias": "exposure_bias",
		"color_space": "color_space", "scene_type": "scene_type",
	}

	for formKey, dbKey := range formFields {
		if val := c.PostForm(formKey); val != "" {
			updates[dbKey] = val
		}
	}

	if err := database.GetDB().Model(&photo).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "更新照片信息失败",
			Message: err.Error(),
		})
		return
	}

	database.GetDB().First(&photo, id)

	stor := storage.GetStorage()
	photo.URL = stor.GetURL(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		photo.ThumbnailURL = stor.GetURL("thumbnails/" + photo.ThumbnailPath)
	}
	if photo.IsLivePhoto && photo.LivePhotoPath != "" {
		photo.LivePhotoURL = stor.GetURL(photo.LivePhotoPath)
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

	stor := storage.GetStorage()
	stor.Delete(photo.StoragePath)
	if photo.ThumbnailPath != "" {
		stor.Delete("thumbnails/" + photo.ThumbnailPath)
	}
	if photo.LivePhotoPath != "" {
		stor.Delete(photo.LivePhotoPath)
	}

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
