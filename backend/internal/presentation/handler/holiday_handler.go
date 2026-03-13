package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type HolidayHandler struct {
	holidayService service.HolidayService
}

func NewHolidayHandler(hs service.HolidayService) *HolidayHandler {
	return &HolidayHandler{holidayService: hs}
}

// holidayRequest is a DTO to accept date as a plain string (YYYY-MM-DD or ISO)
type holidayRequest struct {
	Name        string `json:"name" binding:"required"`
	Date        string `json:"date" binding:"required"`
	Type        string `json:"type" binding:"required"`
	TypeID      uint   `json:"type_id"`
	IsRecurring bool   `json:"is_recurring"`
}

func parseDateFlexible(s string) (time.Time, error) {
	// Try plain date format first, then fallback to RFC3339
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

func (h *HolidayHandler) Create(c *gin.Context) {
	var req holidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := parseDateFlexible(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ngày không hợp lệ, vui lòng dùng định dạng YYYY-MM-DD"})
		return
	}

	hol := &entity.Holiday{
		Name:        req.Name,
		Date:        date,
		Type:        req.Type,
		TypeID:      req.TypeID,
		IsRecurring: req.IsRecurring,
	}

	if err := h.holidayService.CreateHoliday(c.Request.Context(), hol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hol)
}

func (h *HolidayHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	hol, err := h.holidayService.GetHoliday(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy ngày nghỉ"})
		return
	}
	c.JSON(http.StatusOK, hol)
}

func (h *HolidayHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req holidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := parseDateFlexible(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ngày không hợp lệ, vui lòng dùng định dạng YYYY-MM-DD"})
		return
	}

	hol := &entity.Holiday{
		ID:          id,
		Name:        req.Name,
		Date:        date,
		Type:        req.Type,
		TypeID:      req.TypeID,
		IsRecurring: req.IsRecurring,
	}

	if err := h.holidayService.UpdateHoliday(c.Request.Context(), hol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hol)
}

func (h *HolidayHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.holidayService.DeleteHoliday(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

func (h *HolidayHandler) List(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1000"))

	var start, end time.Time
	if startStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr != "" {
		end, _ = time.Parse("2006-01-02", endStr)
	}

	results, total, err := h.holidayService.ListHolidaysPaginated(c.Request.Context(), start, end, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
	})
}
