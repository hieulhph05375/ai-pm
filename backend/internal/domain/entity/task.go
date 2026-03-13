package entity

import "time"

// TaskStatus constants
const (
	TaskStatusTodo       = "TODO"
	TaskStatusInProgress = "IN_PROGRESS"
	TaskStatusDone       = "DONE"
)

// TaskPriority constants
const (
	TaskPriorityLow    = "LOW"
	TaskPriorityMedium = "MEDIUM"
	TaskPriorityHigh   = "HIGH"
	TaskPriorityUrgent = "URGENT"
)

// Task represents a personal task (independent from WBS nodes - ADR-011)
type Task struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	StatusID    uint       `json:"status_id"`
	PriorityID  uint       `json:"priority_id"`
	StatusCat   *Category  `json:"status_cat,omitempty"`
	PriorityCat *Category  `json:"priority_cat,omitempty"`
	AssigneeID  *uint      `json:"assignee_id"`
	Assignee    *User      `json:"assignee,omitempty"`
	CreatedBy   *uint      `json:"created_by"`
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`
	Progress    int        `json:"progress"` // 0-100
	Labels      []string   `json:"labels"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskActivity is an audit log entry for changes on a Task
type TaskActivity struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	ActorID   *uint     `json:"actor_id"`
	Action    string    `json:"action"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}
