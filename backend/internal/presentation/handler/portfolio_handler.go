package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/presentation/middleware"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	portfolioService service.PortfolioService
}

func NewPortfolioHandler(ps service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: ps}
}

func (h *PortfolioHandler) GetOverview(c *gin.Context) {
	userID, isAdmin := middleware.GetUserInfo(c)
	overview, err := h.portfolioService.GetPortfolioOverview(c.Request.Context(), userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải tổng quan portfolio"})
		return
	}
	c.JSON(http.StatusOK, overview)
}
