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

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// --- CategoryType Handlers ---

func (h *CategoryHandler) ListTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	offset := (page - 1) * limit

	types, total, err := h.service.ListTypes(c.Request.Context(), search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch category types"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  types,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *CategoryHandler) CreateType(c *gin.Context) {
	var ct entity.CategoryType
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	if ct.Name == "" || ct.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and code are required"})
		return
	}

	if err := h.service.CreateType(c.Request.Context(), &ct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create category type"})
		return
	}
	c.JSON(http.StatusCreated, ct)
}

func (h *CategoryHandler) UpdateType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ct entity.CategoryType
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	ct.ID = uint(id)

	if err := h.service.UpdateType(c.Request.Context(), &ct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update category type"})
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (h *CategoryHandler) DeleteType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteType(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete category type"})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Category Handlers ---

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	offset := (page - 1) * limit

	var typeIDPtr *uint
	if tIDStr := c.Query("type_id"); tIDStr != "" {
		tID, err := strconv.ParseUint(tIDStr, 10, 32)
		if err == nil {
			uTID := uint(tID)
			typeIDPtr = &uTID
		}
	}

	categories, total, err := h.service.ListCategories(c.Request.Context(), typeIDPtr, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch categories: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var cat entity.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	if cat.Name == "" || cat.TypeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and type are required"})
		return
	}

	if err := h.service.CreateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create category"})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var cat entity.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	cat.ID = uint(id)

	if err := h.service.UpdateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update category"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete category"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	cat, err := h.service.GetCategory(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch category details"})
		return
	}
	c.JSON(http.StatusOK, cat)
}
