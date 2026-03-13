package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(s service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: s}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID := extractActorID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	items, total, err := h.svc.ListByUser(c.Request.Context(), int(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	unread, _ := h.svc.GetUnreadCount(c.Request.Context(), int(userID))

	c.JSON(http.StatusOK, gin.H{
		"items":        items,
		"total":        total,
		"unread_count": unread,
		"page":         page,
		"limit":        limit,
	})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.svc.MarkRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := extractActorID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if err := h.svc.MarkAllRead(c.Request.Context(), int(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := extractActorID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	count, err := h.svc.GetUnreadCount(c.Request.Context(), int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
