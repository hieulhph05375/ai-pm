package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectMemberHandler struct {
	service service.ProjectMemberService
}

func NewProjectMemberHandler(service service.ProjectMemberService) *ProjectMemberHandler {
	return &ProjectMemberHandler{service: service}
}

func (h *ProjectMemberHandler) GetMembers(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	members, total, err := h.service.GetMembersByProject(c.Request.Context(), projectID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  members,
		"total": total,
	})
}

func (h *ProjectMemberHandler) AddMember(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		UserID int `json:"user_id" binding:"required"`
		RoleID int `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member := &entity.ProjectMember{
		ProjectID:     projectID,
		UserID:        req.UserID,
		ProjectRoleID: req.RoleID,
	}

	if err := h.service.AddMember(c.Request.Context(), member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *ProjectMemberHandler) RemoveMember(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := strconv.Atoi(c.Param("userId"))

	if err := h.service.RemoveMember(c.Request.Context(), projectID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProjectMemberHandler) UpdateMemberRole(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	// userId is the URL param for the user being updated
	userID, _ := strconv.Atoi(c.Param("userId"))

	var req struct {
		RoleID int `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMemberRole(c.Request.Context(), projectID, userID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
