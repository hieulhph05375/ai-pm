package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

func parseDate(s string) *time.Time {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return &t
		}
	}
	log.Printf("WARNING: could not parse date string: %q", s)
	return nil
}

type WBSHandler struct {
	wbsService service.WBSService
}

func NewWBSHandler(wbsService service.WBSService) *WBSHandler {
	return &WBSHandler{wbsService: wbsService}
}

func (h *WBSHandler) ListTree(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	filter := entity.WBSFilter{
		Search:     c.Query("search"),
		Status:     c.Query("status"),
		ParentPath: c.Query("parent_path"),
	}

	if fields := c.Query("fields"); fields != "" {
		filter.Fields = strings.Split(fields, ",")
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	picIDStr := c.Query("assigned_to")
	log.Printf("[WBS Handler] ListTree: project=%s search=%q status=%q assigned_to=%q parent=%q fields=%v",
		idStr, filter.Search, filter.Status, picIDStr, filter.ParentPath, filter.Fields)

	if picIDStr != "" {
		if picID, err := strconv.Atoi(picIDStr); err == nil {
			filter.AssignedTo = &picID
		}
	}

	nodes, total, err := h.wbsService.GetProjectTree(c.Request.Context(), projectID, filter)
	if err != nil {
		log.Printf("Error fetching WBS tree: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch WBS tree"})
		return
	}

	// Trả về flat array kèm pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":  nodes,
		"total": total,
	})
}

func (h *WBSHandler) CreateNode(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	var req struct {
		ParentPath       string  `json:"parent_path"`
		Title            string  `json:"title"`
		Type             string  `json:"type"`
		AssignedTo       *int    `json:"assigned_to"`
		Description      *string `json:"description"`
		PlannedStartDate *string `json:"planned_start_date"`
		PlannedEndDate   *string `json:"planned_end_date"`
		PlannedValue     float64 `json:"planned_value"`
		ActualCost       float64 `json:"actual_cost"`
		Progress         float64 `json:"progress"`
		EstimatedEffort  float64 `json:"estimated_effort"`
		ActualEffort     float64 `json:"actual_effort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	node := &entity.WBSNode{
		ProjectID:       projectID,
		Title:           req.Title,
		Type:            entity.WBSNodeType(req.Type),
		AssignedTo:      req.AssignedTo,
		Description:     req.Description,
		PlannedValue:    req.PlannedValue,
		ActualCost:      req.ActualCost,
		Progress:        req.Progress,
		EstimatedEffort: req.EstimatedEffort,
		ActualEffort:    req.ActualEffort,
	}

	if req.PlannedStartDate != nil && *req.PlannedStartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.PlannedStartDate)
		if err == nil {
			node.PlannedStartDate = &t
		}
	}
	if req.PlannedEndDate != nil && *req.PlannedEndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.PlannedEndDate)
		if err == nil {
			node.PlannedEndDate = &t
		}
	}

	err = h.wbsService.CreateNode(c.Request.Context(), node, req.ParentPath)
	if err != nil {
		log.Printf("Error creating WBS node: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": node})
}

func (h *WBSHandler) UpdateNode(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID format"})
		return
	}

	var req struct {
		Title            string  `json:"title"`
		Type             string  `json:"type"`
		OrderIndex       int     `json:"order_index"`
		Progress         float64 `json:"progress"`
		PlannedValue     float64 `json:"planned_value"`
		ActualCost       float64 `json:"actual_cost"`
		EstimatedEffort  float64 `json:"estimated_effort"`
		ActualEffort     float64 `json:"actual_effort"`
		AssignedTo       *int    `json:"assigned_to"`
		Description      *string `json:"description"`
		PlannedStartDate *string `json:"planned_start_date"`
		PlannedEndDate   *string `json:"planned_end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	node := &entity.WBSNode{
		ID:              nodeID,
		Title:           req.Title,
		Type:            entity.WBSNodeType(req.Type),
		OrderIndex:      req.OrderIndex,
		Progress:        req.Progress,
		PlannedValue:    req.PlannedValue,
		ActualCost:      req.ActualCost,
		EstimatedEffort: req.EstimatedEffort,
		ActualEffort:    req.ActualEffort,
		AssignedTo:      req.AssignedTo,
		Description:     req.Description,
	}

	if req.PlannedStartDate != nil && *req.PlannedStartDate != "" {
		node.PlannedStartDate = parseDate(*req.PlannedStartDate)
	}
	if req.PlannedEndDate != nil && *req.PlannedEndDate != "" {
		node.PlannedEndDate = parseDate(*req.PlannedEndDate)
	}

	log.Printf("[UpdateNode] nodeID=%d title=%q type=%q progress=%.0f PV=%.0f AC=%.0f assignedTo=%v startDate=%v endDate=%v",
		nodeID, req.Title, req.Type, req.Progress, req.PlannedValue, req.ActualCost, req.AssignedTo,
		req.PlannedStartDate, req.PlannedEndDate)

	err = h.wbsService.UpdateNode(c.Request.Context(), node)
	if err != nil {
		log.Printf("Error updating WBS node: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node updated successfully"})
}

func (h *WBSHandler) DeleteNode(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID format"})
		return
	}

	err = h.wbsService.DeleteNode(c.Request.Context(), nodeID)
	if err != nil {
		log.Printf("Error deleting WBS node: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node deleted"})
}

func (h *WBSHandler) GetNode(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID format"})
		return
	}

	node, err := h.wbsService.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		log.Printf("Error fetching WBS node: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": node})
}

// Dependency endpoints
func (h *WBSHandler) ListDependencies(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	deps, err := h.wbsService.ListDependencies(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dependencies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": deps})
}

func (h *WBSHandler) CreateDependency(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	var req struct {
		PredecessorID int    `json:"predecessor_id"`
		SuccessorID   int    `json:"successor_id"`
		Type          string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	depType := entity.DependencyType(req.Type)
	if depType == "" {
		depType = entity.DepFS
	}
	dep := &entity.WBSDependency{
		ProjectID:     projectID,
		PredecessorID: req.PredecessorID,
		SuccessorID:   req.SuccessorID,
		Type:          depType,
	}
	if err := h.wbsService.CreateDependency(c.Request.Context(), dep); err != nil {
		log.Printf("Error creating dependency: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dependency"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": dep})
}

func (h *WBSHandler) DeleteDependency(c *gin.Context) {
	depID, err := strconv.Atoi(c.Param("depId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dependency ID"})
		return
	}
	if err := h.wbsService.DeleteDependency(c.Request.Context(), depID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete dependency"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dependency deleted"})
}

// Comment endpoints
func (h *WBSHandler) ListComments(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}
	comments, err := h.wbsService.ListComments(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func (h *WBSHandler) AddComment(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	nodeID, _ := strconv.Atoi(c.Param("nodeId"))

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	userID_raw, exists := c.Get("userID") // Provided by auth middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// JWT numbers are floats or strings depending on claims
	var userID int
	switch v := userID_raw.(type) {
	case float64:
		userID = int(v)
	case string:
		// Convert string ID if necessary (e.g., strconv.Atoi)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	comment := &entity.WBSComment{
		ProjectID: projectID,
		NodeID:    nodeID,
		UserID:    userID,
		Content:   req.Content,
	}

	if err := h.wbsService.AddComment(c.Request.Context(), comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}

func (h *WBSHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	userID_raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isAdmin, _ := c.Get("isAdmin")

	// Get comment to check ownership
	comment, err := h.wbsService.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Permission check: only admin or the user who created the comment can delete
	if (!isAdmin.(bool)) && (int(userID_raw.(float64)) != comment.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this comment"})
		return
	}

	if err := h.wbsService.DeleteComment(c.Request.Context(), commentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

func (h *WBSHandler) UpdateComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	userID_raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isAdmin, _ := c.Get("isAdmin")

	// Get comment to check ownership
	comment, err := h.wbsService.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Permission check: only admin or the user who created the comment can update
	if (!isAdmin.(bool)) && (int(userID_raw.(float64)) != comment.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this comment"})
		return
	}

	if err := h.wbsService.UpdateComment(c.Request.Context(), commentID, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

// Baseline endpoints
func (h *WBSHandler) CreateBaseline(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	userID_raw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var userID int
	switch v := userID_raw.(type) {
	case float64:
		userID = int(v)
	case string:
		userID, _ = strconv.Atoi(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	baseline, err := h.wbsService.CreateBaseline(c.Request.Context(), projectID, req.Name, req.Description, userID)
	if err != nil {
		log.Printf("Error creating baseline: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create baseline"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": baseline})
}

func (h *WBSHandler) ListBaselines(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	baselines, err := h.wbsService.GetBaselines(c.Request.Context(), projectID)
	if err != nil {
		log.Printf("Error listing baselines: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch baselines"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": baselines})
}

func (h *WBSHandler) GetBaselineNodes(c *gin.Context) {
	baselineID, err := strconv.Atoi(c.Param("baselineId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid baseline ID"})
		return
	}

	nodes, err := h.wbsService.GetBaselineNodes(c.Request.Context(), baselineID)
	if err != nil {
		log.Printf("Error fetching baseline nodes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch baseline nodes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": nodes})
}
