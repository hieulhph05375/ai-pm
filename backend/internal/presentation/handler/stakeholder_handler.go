package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StakeholderHandler struct {
	stakeholderService service.StakeholderService
}

func NewStakeholderHandler(ss service.StakeholderService) *StakeholderHandler {
	return &StakeholderHandler{stakeholderService: ss}
}

func (h *StakeholderHandler) Create(c *gin.Context) {
	var s entity.Stakeholder
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	if err := h.stakeholderService.CreateStakeholder(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo stakeholder: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (h *StakeholderHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	s, err := h.stakeholderService.GetStakeholder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy stakeholder"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StakeholderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var s entity.Stakeholder
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.ID = id

	if err := h.stakeholderService.UpdateStakeholder(c.Request.Context(), &s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StakeholderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.stakeholderService.DeleteStakeholder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

func (h *StakeholderHandler) List(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	results, total, err := h.stakeholderService.ListStakeholdersPaginated(c.Request.Context(), search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
	})
}

// Project Mapping Handlers

func (h *StakeholderHandler) ListByProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	results, err := h.stakeholderService.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *StakeholderHandler) AssignToProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		StakeholderID int    `json:"stakeholder_id" binding:"required"`
		Role          string `json:"role"`
		RoleID        uint   `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.stakeholderService.AssignToProject(c.Request.Context(), projectID, req.StakeholderID, req.Role, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Gán thành công"})
}

func (h *StakeholderHandler) UnassignFromProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	stakeholderID, _ := strconv.Atoi(c.Param("sid"))

	if err := h.stakeholderService.UnassignFromProject(c.Request.Context(), projectID, stakeholderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hủy gán thành công"})
}
