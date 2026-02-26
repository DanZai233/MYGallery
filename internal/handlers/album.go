package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/models"
	"github.com/mygallery/mygallery/internal/storage"
)

type AlbumHandler struct{}

func NewAlbumHandler() *AlbumHandler { return &AlbumHandler{} }

func (h *AlbumHandler) GetAlbums(c *gin.Context) {
	db := database.GetDB()
	var albums []models.Album
	db.Order("created_at DESC").Find(&albums)

	stor := storage.GetStorage()
	for i := range albums {
		db.Model(&models.PhotoAlbum{}).Where("album_id = ?", albums[i].ID).Count(&albums[i].PhotoCount)
		var pa models.PhotoAlbum
		if err := db.Where("album_id = ?", albums[i].ID).Order("created_at DESC").First(&pa).Error; err == nil {
			var photo models.Photo
			if err := db.First(&photo, pa.PhotoID).Error; err == nil {
				if photo.ThumbnailPath != "" {
					albums[i].CoverURL = stor.GetURL("thumbnails/" + photo.ThumbnailPath)
				} else {
					albums[i].CoverURL = stor.GetURL(photo.StoragePath)
				}
			}
		}
	}
	c.JSON(http.StatusOK, albums)
}

func (h *AlbumHandler) GetAlbum(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var album models.Album
	if err := db.First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}
	db.Model(&models.PhotoAlbum{}).Where("album_id = ?", album.ID).Count(&album.PhotoCount)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 { page = 1 }
	if size < 1 || size > 100 { size = 20 }

	var photoAlbums []models.PhotoAlbum
	db.Where("album_id = ?", album.ID).Order("created_at DESC").Offset((page-1)*size).Limit(size).Find(&photoAlbums)

	photoIDs := make([]uint, len(photoAlbums))
	for i, pa := range photoAlbums { photoIDs[i] = pa.PhotoID }

	var photos []models.Photo
	if len(photoIDs) > 0 {
		db.Where("id IN ?", photoIDs).Find(&photos)
		stor := storage.GetStorage()
		for i := range photos {
			photos[i].URL = stor.GetURL(photos[i].StoragePath)
			if photos[i].ThumbnailPath != "" {
				photos[i].ThumbnailURL = stor.GetURL("thumbnails/" + photos[i].ThumbnailPath)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"album":  album,
		"photos": photos,
		"total":  album.PhotoCount,
		"page":   page,
		"size":   size,
	})
}

func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}
	album := models.Album{Title: body.Title, Description: body.Description}
	database.GetDB().Create(&album)
	c.JSON(http.StatusOK, gin.H{"success": true, "album": album})
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	var album models.Album
	if err := database.GetDB().First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "相册不存在"})
		return
	}
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	c.ShouldBindJSON(&body)
	updates := map[string]interface{}{}
	if body.Title != "" { updates["title"] = body.Title }
	if body.Description != "" { updates["description"] = body.Description }
	database.GetDB().Model(&album).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"success": true, "album": album})
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	database.GetDB().Where("album_id = ?", id).Delete(&models.PhotoAlbum{})
	database.GetDB().Delete(&models.Album{}, id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AlbumHandler) AddPhotos(c *gin.Context) {
	albumID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var body struct {
		PhotoIDs []uint `json:"photo_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo_ids 不能为空"})
		return
	}
	db := database.GetDB()
	added := 0
	for _, pid := range body.PhotoIDs {
		var exists int64
		db.Model(&models.PhotoAlbum{}).Where("album_id = ? AND photo_id = ?", albumID, pid).Count(&exists)
		if exists == 0 {
			db.Create(&models.PhotoAlbum{AlbumID: uint(albumID), PhotoID: pid})
			added++
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "added": added})
}

func (h *AlbumHandler) RemovePhoto(c *gin.Context) {
	albumID := c.Param("id")
	photoID := c.Param("photoId")
	database.GetDB().Where("album_id = ? AND photo_id = ?", albumID, photoID).Delete(&models.PhotoAlbum{})
	c.JSON(http.StatusOK, gin.H{"success": true})
}
