package handler

import (
	"net/http"
	"project-mgmt/backend/internal/application/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingService service.SettingService
}

func NewSettingHandler(ss service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: ss}
}

// GetAll returns all system settings
func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.settingService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Convert to map for easy frontend consumption
	result := make(map[string]any)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, result)
}

// Update updates a setting
func (h *SettingHandler) Update(c *gin.Context) {
	key := c.Param("key")
	var body struct {
		Value any `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingService.Set(c.Request.Context(), key, body.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": body.Value})
}
