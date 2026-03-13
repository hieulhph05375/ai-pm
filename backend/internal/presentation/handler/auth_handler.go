package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hồ sơ đăng nhập không hợp lệ"})
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"fullName": user.FullName,
			"isAdmin":  user.IsAdmin,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đăng ký không hợp lệ"})
		return
	}

	user := &entity.User{
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   5, // Default to Member
	}

	if err := h.authService.Register(c.Request.Context(), user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đăng ký người dùng"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Đăng ký thành công"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token làm mới không hợp lệ"})
		return
	}

	accessToken, refreshToken, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Phiên làm việc đã hết hạn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
