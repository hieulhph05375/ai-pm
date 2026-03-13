package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TimesheetHandler struct {
	tsService service.TimesheetService
}

func NewTimesheetHandler(s service.TimesheetService) *TimesheetHandler {
	return &TimesheetHandler{tsService: s}
}

func extractUserID(c *gin.Context) int {
	v, exists := c.Get("userID")
	if !exists {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case uint:
		return int(val)
	}
	return 0
}

func (h *TimesheetHandler) Create(c *gin.Context) {
	userID := extractUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var t entity.Timesheet
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	t.UserID = userID

	if err := h.tsService.CreateTimesheet(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TimesheetHandler) List(c *gin.Context) {
	userID := extractUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	projectID, _ := strconv.Atoi(c.Query("project_id"))

	offset := (page - 1) * limit

	var timesheets []entity.Timesheet
	var err error
	var total int

	if projectID > 0 {
		timesheets, total, err = h.tsService.ListByProject(c.Request.Context(), projectID, limit, offset)
	} else {
		timesheets, total, err = h.tsService.ListByUser(c.Request.Context(), userID, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timesheets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": timesheets,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *TimesheetHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	t, err := h.tsService.GetTimesheet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Timesheet not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TimesheetHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var t entity.Timesheet
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}
	t.ID = id

	if err := h.tsService.UpdateTimesheet(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TimesheetHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.tsService.DeleteTimesheet(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete timesheet"})
		return
	}
	c.Status(http.StatusNoContent)
}
