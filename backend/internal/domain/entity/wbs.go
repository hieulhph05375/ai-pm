package entity

import "time"

type WBSNodeType string

const (
	TypePhase     WBSNodeType = "Phase"
	TypeMilestone WBSNodeType = "Milestone"
	TypeTask      WBSNodeType = "Task"
	TypeSubTask   WBSNodeType = "Sub-task"
)

type DependencyType string

const (
	DepFS DependencyType = "FS"
	DepSF DependencyType = "SF"
	DepSS DependencyType = "SS"
	DepFF DependencyType = "FF"
)

// WBSNode represents a node in the Work Breakdown Structure
type WBSNode struct {
	ID               int         `json:"id" db:"id"`
	ProjectID        int         `json:"project_id" db:"project_id"`
	Title            string      `json:"title" db:"title"`
	Type             WBSNodeType `json:"type" db:"type"`
	TypeID           *uint       `json:"type_id" db:"type_id"`
	TypeCat          *Category   `json:"type_cat,omitempty"`
	Path             string      `json:"path" db:"path"`
	OrderIndex       int         `json:"order_index" db:"order_index"`
	PlannedStartDate *time.Time  `json:"planned_start_date" db:"planned_start_date"`
	PlannedEndDate   *time.Time  `json:"planned_end_date" db:"planned_end_date"`
	ActualStartDate  *time.Time  `json:"actual_start_date" db:"actual_start_date"`
	ActualEndDate    *time.Time  `json:"actual_end_date" db:"actual_end_date"`
	Progress         float64     `json:"progress" db:"progress"`
	PlannedValue     float64     `json:"planned_value" db:"planned_value"`
	ActualCost       float64     `json:"actual_cost" db:"actual_cost"`
	EstimatedEffort  float64     `json:"estimated_effort" db:"estimated_effort"`
	ActualEffort     float64     `json:"actual_effort" db:"actual_effort"`
	AssignedTo       *int        `json:"assigned_to" db:"assigned_to"`
	Description      *string     `json:"description" db:"description"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`

	// Optional fields for frontend nested tree representation
	Children    []*WBSNode `json:"children,omitempty"`
	HasChildren bool       `json:"has_children"`
}

// WBSDependency represents the relationship between two WBS nodes
type WBSDependency struct {
	ID            int            `json:"id" db:"id"`
	ProjectID     int            `json:"project_id" db:"project_id"`
	PredecessorID int            `json:"predecessor_id" db:"predecessor_id"`
	SuccessorID   int            `json:"successor_id" db:"successor_id"`
	Type          DependencyType `json:"type" db:"type"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
}

type WBSComment struct {
	ID        int       `json:"id" db:"id"`
	ProjectID int       `json:"project_id" db:"project_id"`
	NodeID    int       `json:"node_id" db:"node_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name" db:"user_name"` // Joined from users table
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WBSFilter struct {
	Search     string
	AssignedTo *int
	Status     string   // 'todo', 'doing', 'done'
	ParentPath string   // For lazy loading: children of this path
	Fields     []string // For field selection: specific columns to fetch
	Page       int      // Pagination: current page (1-based)
	Limit      int      // Pagination: items per page
	FetchAll   bool     // Bypass depth and pagination limits to fetch the entire tree
}

type WBSBaseline struct {
	ID          int       `json:"id" db:"id"`
	ProjectID   int       `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type WBSBaselineNode struct {
	Options          interface{} `json:"-"` // Added for generic compat
	BaselineID       int         `json:"baseline_id" db:"baseline_id"`
	NodeID           int         `json:"node_id" db:"node_id"`
	Path             string      `json:"path" db:"path"`
	PlannedStartDate *time.Time  `json:"planned_start_date" db:"planned_start_date"`
	PlannedEndDate   *time.Time  `json:"planned_end_date" db:"planned_end_date"`
	Progress         float64     `json:"progress" db:"progress"`
	PlannedValue     float64     `json:"planned_value" db:"planned_value"`
	ActualCost       float64     `json:"actual_cost" db:"actual_cost"`
}
