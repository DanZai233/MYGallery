package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/models"
)

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler { return &CommentHandler{} }

func (h *CommentHandler) GetComments(c *gin.Context) {
	photoID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var comments []models.Comment
	database.GetDB().Where("photo_id = ?", photoID).Order("created_at DESC").Find(&comments)
	c.JSON(http.StatusOK, gin.H{"comments": comments, "total": len(comments)})
}

func (h *CommentHandler) AddComment(c *gin.Context) {
	photoID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var body struct {
		Nickname string `json:"nickname" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称和内容不能为空"})
		return
	}
	body.Nickname = strings.TrimSpace(body.Nickname)
	body.Content = strings.TrimSpace(body.Content)
	if len(body.Nickname) > 50 || len(body.Content) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称最多50字，内容最多1000字"})
		return
	}

	comment := models.Comment{
		PhotoID:   uint(photoID),
		Nickname:  body.Nickname,
		Content:   body.Content,
		IPAddress: c.ClientIP(),
	}
	database.GetDB().Create(&comment)
	c.JSON(http.StatusOK, gin.H{"success": true, "comment": comment})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("commentId")
	database.GetDB().Delete(&models.Comment{}, id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
