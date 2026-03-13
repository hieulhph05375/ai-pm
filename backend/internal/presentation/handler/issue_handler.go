package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IssueHandler struct {
	issueService service.IssueService
}

func NewIssueHandler(is service.IssueService) *IssueHandler {
	return &IssueHandler{issueService: is}
}

func (h *IssueHandler) List(c *gin.Context) {
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

	issues, total, err := h.issueService.ListIssues(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tải danh sách issues"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": issues,
		"total": total,
	})
}

func (h *IssueHandler) Create(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã dự án không hợp lệ"})
		return
	}
	var i entity.Issue
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	i.ProjectID = projectID
	if err := h.issueService.CreateIssue(c.Request.Context(), &i); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo issue"})
		return
	}
	c.JSON(http.StatusCreated, i)
}

func (h *IssueHandler) Update(c *gin.Context) {
	issueID, err := strconv.Atoi(c.Param("issueId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã issue không hợp lệ"})
		return
	}
	var i entity.Issue
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	i.ID = issueID
	if err := h.issueService.UpdateIssue(c.Request.Context(), &i); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật issue"})
		return
	}
	c.JSON(http.StatusOK, i)
}

func (h *IssueHandler) Delete(c *gin.Context) {
	issueID, err := strconv.Atoi(c.Param("issueId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã issue không hợp lệ"})
		return
	}
	if err := h.issueService.DeleteIssue(c.Request.Context(), issueID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa issue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa issue thành công"})
}
