package handler

import (
	"fmt"
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/presentation/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService service.ReportService
	exportService service.ExportService
}

func NewReportHandler(rService service.ReportService, eService service.ExportService) *ReportHandler {
	return &ReportHandler{
		reportService: rService,
		exportService: eService,
	}
}

func (h *ReportHandler) GetPMIStats(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, isAdmin := middleware.GetUserInfo(c)
	stats, err := h.reportService.GetProjectPMIStats(c.Request.Context(), id, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *ReportHandler) ExportWBS(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	c.Header("Content-Disposition", "attachment; filename=wbs_export.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := h.exportService.ExportWBSExcel(c.Request.Context(), id, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xuất file Excel: " + err.Error()})
		return
	}
}

func (h *ReportHandler) ExportSummary(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, isAdmin := middleware.GetUserInfo(c)

	stats, err := h.reportService.GetProjectPMIStats(c.Request.Context(), id, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy thông số báo cáo: " + err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=project_summary_%d.pdf", id))
	c.Header("Content-Type", "application/pdf")

	if err := h.exportService.ExportProjectSummaryPDF(c.Request.Context(), id, stats, userID, isAdmin, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xuất file PDF: " + err.Error()})
		return
	}
}

func (h *ReportHandler) ExportProjectList(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	c.Header("Content-Disposition", "attachment; filename=project_list_export.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	userID, isAdmin := middleware.GetUserInfo(c)
	if err := h.exportService.ExportProjectListExcel(c.Request.Context(), search, status, userID, isAdmin, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xuất danh sách dự án: " + err.Error()})
		return
	}
}
