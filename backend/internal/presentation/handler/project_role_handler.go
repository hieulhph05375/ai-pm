package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectRoleHandler struct {
	service service.ProjectRoleService
}

func NewProjectRoleHandler(service service.ProjectRoleService) *ProjectRoleHandler {
	return &ProjectRoleHandler{service: service}
}

func (h *ProjectRoleHandler) GetPermissions(c *gin.Context) {
	perms, err := h.service.GetAllPermissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}

func (h *ProjectRoleHandler) GetRolesByProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	roles, err := h.service.GetRolesByProject(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *ProjectRoleHandler) CreateRole(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	var role entity.ProjectRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role.ProjectID = uint(projectID)

	if err := h.service.CreateRole(c.Request.Context(), &role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, role)
}

func (h *ProjectRoleHandler) UpdateRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleId"))
	var role entity.ProjectRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role.ID = uint(roleID)

	if err := h.service.UpdateRole(c.Request.Context(), &role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (h *ProjectRoleHandler) DeleteRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleId"))
	if err := h.service.DeleteRole(c.Request.Context(), uint(roleID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProjectRoleHandler) SetPermissions(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleId"))
	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetRolePermissions(c.Request.Context(), uint(roleID), req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *ProjectRoleHandler) GetRolePermissions(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleId"))
	perms, err := h.service.GetRolePermissions(c.Request.Context(), uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}
