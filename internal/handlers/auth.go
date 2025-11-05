package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/database"
	"github.com/mygallery/mygallery/internal/middleware"
	"github.com/mygallery/mygallery/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "参数错误",
			Message: err.Error(),
		})
		return
	}
	
	// 查找用户
	var user models.User
	if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "用户名或密码错误",
		})
		return
	}
	
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "用户名或密码错误",
		})
		return
	}
	
	// 生成 token
	token, err := middleware.GenerateToken(h.cfg, user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "生成令牌失败",
			Message: err.Error(),
		})
		return
	}
	
	// 返回响应
	var resp models.LoginResponse
	resp.Token = token
	resp.User.ID = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email
	resp.User.Role = user.Role
	
	c.JSON(http.StatusOK, resp)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "参数错误",
			Message: err.Error(),
		})
		return
	}
	
	// 获取用户
	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "用户不存在"})
		return
	}
	
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "原密码错误"})
		return
	}
	
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "密码加密失败",
			Message: err.Error(),
		})
		return
	}
	
	// 更新密码
	user.Password = string(hashedPassword)
	if err := database.GetDB().Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "更新密码失败",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码修改成功",
	})
}

