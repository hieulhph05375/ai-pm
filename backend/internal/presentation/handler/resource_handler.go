package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"time"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	resourceService service.ResourceService
}

func NewResourceHandler(rs service.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService: rs}
}

func (h *ResourceHandler) GetWorkload(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Default to current month if not provided
	if startDate == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	if endDate == "" {
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}

	overview, err := h.resourceService.GetWorkload(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải dữ liệu phân bổ nguồn lực"})
		return
	}

	c.JSON(http.StatusOK, overview)
}
