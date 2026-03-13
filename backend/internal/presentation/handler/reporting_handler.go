package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportingHandler struct {
	snapshotService service.SnapshotService
}

func NewReportingHandler(ss service.SnapshotService) *ReportingHandler {
	return &ReportingHandler{snapshotService: ss}
}

func (h *ReportingHandler) CaptureSnapshots(c *gin.Context) {
	if err := h.snapshotService.CaptureAllProjectsSnapshot(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Snapshots captured successfully"})
}

func (h *ReportingHandler) GetProjectTrends(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	trends, err := h.snapshotService.GetProjectTrends(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

func (h *ReportingHandler) GetMilestoneTrends(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	trends, err := h.snapshotService.GetMilestoneTrends(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}
