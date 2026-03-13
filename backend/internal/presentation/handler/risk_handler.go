package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RiskHandler struct {
	riskService service.RiskService
}

func NewRiskHandler(rs service.RiskService) *RiskHandler {
	return &RiskHandler{riskService: rs}
}

func (h *RiskHandler) List(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã dự án không hợp lệ"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	risks, total, err := h.riskService.ListRisks(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải danh sách rủi ro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": risks,
		"total": total,
	})
}

func (h *RiskHandler) Create(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã dự án không hợp lệ"})
		return
	}
	var r entity.Risk
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	r.ProjectID = projectID
	if err := h.riskService.CreateRisk(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo rủi ro"})
		return
	}
	c.JSON(http.StatusCreated, r)
}

func (h *RiskHandler) Update(c *gin.Context) {
	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã rủi ro không hợp lệ"})
		return
	}
	var r entity.Risk
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	r.ID = riskID
	if err := h.riskService.UpdateRisk(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật rủi ro"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *RiskHandler) Delete(c *gin.Context) {
	riskID, err := strconv.Atoi(c.Param("riskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã rủi ro không hợp lệ"})
		return
	}
	if err := h.riskService.DeleteRisk(c.Request.Context(), riskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa rủi ro"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa rủi ro thành công"})
}
