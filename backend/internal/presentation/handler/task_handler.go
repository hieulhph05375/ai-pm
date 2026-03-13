package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(ts service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: ts}
}

// extractActorID extracts the current user's ID from the Gin context (set by auth middleware)
func extractActorID(c *gin.Context) uint {
	v, exists := c.Get("userID")
	if !exists {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return uint(val)
	case uint:
		return val
	case int:
		return uint(val)
	}
	return 0
}

// List lấy danh sách task của user hiện tại (Kanban/Gantt cá nhân) với phân trang
func (h *TaskHandler) List(c *gin.Context) {
	userID := extractActorID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	tasks, total, err := h.taskService.ListTasksByUser(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách công việc"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": tasks,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Create tạo công việc mới cho user hiện tại
func (h *TaskHandler) Create(c *gin.Context) {
	actorID := extractActorID(c)
	if actorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập"})
		return
	}

	var t entity.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	if t.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tiêu đề là bắt buộc"})
		return
	}

	if err := h.taskService.CreateTask(c.Request.Context(), &t, actorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo công việc"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// GetByID lấy chi tiết một task
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy công việc"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy thông tin công việc"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Update cập nhật thông tin task
func (h *TaskHandler) Update(c *gin.Context) {
	actorID := extractActorID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var t entity.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	t.ID = uint(id)

	if err := h.taskService.UpdateTask(c.Request.Context(), &t, actorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật công việc"})
		return
	}
	c.JSON(http.StatusOK, t)
}

// Delete xóa task
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := h.taskService.DeleteTask(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa công việc"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListActivities lấy lịch sử hoạt động của task
func (h *TaskHandler) ListActivities(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	activities, err := h.taskService.ListActivities(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy lịch sử hoạt động"})
		return
	}
	c.JSON(http.StatusOK, activities)
}

// AddComment thêm comment cho task
func (h *TaskHandler) AddComment(c *gin.Context) {
	actorID := extractActorID(c)
	if actorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nội dung comment là bắt buộc"})
		return
	}

	if err := h.taskService.AddComment(c.Request.Context(), uint(id), actorID, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể thêm bình luận"})
		return
	}
	c.Status(http.StatusCreated)
}
