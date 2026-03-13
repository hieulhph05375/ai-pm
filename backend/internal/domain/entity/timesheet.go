package entity

import "time"

const (
	TimesheetStatusDraft     = "DRAFT"
	TimesheetStatusSubmitted = "SUBMITTED"
	TimesheetStatusApproved  = "APPROVED"
	TimesheetStatusRejected  = "REJECTED"
)

type Timesheet struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	ProjectID   *int      `json:"project_id,omitempty" db:"project_id"`
	NodeID      *int      `json:"node_id,omitempty" db:"node_id"`
	TaskID      *int      `json:"task_id,omitempty" db:"task_id"`
	WorkDate    time.Time `json:"work_date" db:"work_date"`
	Hours       float64   `json:"hours" db:"hours"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Joins
	UserName    string `json:"user_name,omitempty" db:"user_name"`
	ProjectName string `json:"project_name,omitempty" db:"project_name"`
	NodeTitle   string `json:"node_title,omitempty" db:"node_title"`
	TaskTitle   string `json:"task_title,omitempty" db:"task_title"`
}
