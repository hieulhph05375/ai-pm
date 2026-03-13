package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) getCtxUserID(c *gin.Context) (uint, bool) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	switch v := userIDStr.(type) {
	case string:
		sid, _ := strconv.ParseUint(v, 10, 64)
		return uint(sid), true
	case float64:
		return uint(v), true
	}
	return 0, false
}

func (h *UserHandler) GetMe(c *gin.Context) {
	id, ok := h.getCtxUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) List(c *gin.Context) {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1000"))

	users, total, err := h.userService.ListUsersPaginated(c.Request.Context(), q, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách người dùng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	var input struct {
		entity.User
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.CreateUser(c.Request.Context(), &input.User, input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input.User)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	adminID, ok := h.getCtxUserID(c)
	if !ok {
		return
	}

	var u entity.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.ID = uint(id)

	if err := h.userService.UpdateUser(c.Request.Context(), adminID, &u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var input struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, ok := h.getCtxUserID(c)
	if !ok {
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), adminID, uint(id), input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	adminID, ok := h.getCtxUserID(c)
	if !ok {
		return
	}

	if err := h.userService.ToggleStatus(c.Request.Context(), adminID, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status toggled"})
}
