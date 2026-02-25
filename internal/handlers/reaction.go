package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/models"
)

var validReactions = map[string]bool{
	"like": true, "love": true, "amazing": true, "funny": true,
	"wow": true, "sad": true, "fire": true, "sparkle": true,
}

type ReactionHandler struct{}

func NewReactionHandler() *ReactionHandler {
	return &ReactionHandler{}
}

func getFingerprint(c *gin.Context) string {
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	lang := c.GetHeader("Accept-Language")
	raw := fmt.Sprintf("%s|%s|%s", ip, ua, lang)
	h := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", h[:12])
}

// GetReactions returns reaction counts + current user's reaction for a photo
func (h *ReactionHandler) GetReactions(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
		return
	}

	db := database.GetDB()
	fp := getFingerprint(c)

	var counts []struct {
		ReactionType string
		Count        int64
	}
	db.Model(&models.PhotoReaction{}).
		Select("reaction_type, count(*) as count").
		Where("photo_id = ?", photoID).
		Group("reaction_type").
		Find(&counts)

	reactions := make(map[string]int64)
	for _, t := range []string{"like", "love", "amazing", "funny", "wow", "sad", "fire", "sparkle"} {
		reactions[t] = 0
	}
	for _, c := range counts {
		reactions[c.ReactionType] = c.Count
	}

	var userReaction models.PhotoReaction
	userReactionType := ""
	if err := db.Where("photo_id = ? AND fingerprint = ?", photoID, fp).First(&userReaction).Error; err == nil {
		userReactionType = userReaction.ReactionType
	}

	c.JSON(http.StatusOK, models.ReactionResponse{
		PhotoID:      uint(photoID),
		Reactions:    reactions,
		UserReaction: userReactionType,
	})
}

// AddReaction creates or updates a reaction
func (h *ReactionHandler) AddReaction(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
		return
	}

	var body struct {
		ReactionType string `json:"reaction_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reaction_type is required"})
		return
	}

	if !validReactions[body.ReactionType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reaction type"})
		return
	}

	db := database.GetDB()
	fp := getFingerprint(c)

	var existing models.PhotoReaction
	if err := db.Where("photo_id = ? AND fingerprint = ?", photoID, fp).First(&existing).Error; err == nil {
		db.Model(&existing).Update("reaction_type", body.ReactionType)
		c.JSON(http.StatusOK, gin.H{"success": true, "action": "updated", "reaction_type": body.ReactionType})
		return
	}

	reaction := models.PhotoReaction{
		PhotoID:      uint(photoID),
		ReactionType: body.ReactionType,
		Fingerprint:  fp,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
	}
	db.Create(&reaction)
	c.JSON(http.StatusOK, gin.H{"success": true, "action": "created", "reaction_type": body.ReactionType})
}

// DeleteReaction removes the user's reaction
func (h *ReactionHandler) DeleteReaction(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
		return
	}

	db := database.GetDB()
	fp := getFingerprint(c)

	result := db.Where("photo_id = ? AND fingerprint = ?", photoID, fp).Delete(&models.PhotoReaction{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no reaction found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "action": "deleted"})
}
