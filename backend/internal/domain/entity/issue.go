package entity

import "time"

type Issue struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	TypeID      uint      `json:"type_id"`
	PriorityID  uint      `json:"priority_id"`
	StatusID    uint      `json:"status_id"`
	Type        *Category `json:"type,omitempty"`
	Priority    *Category `json:"priority,omitempty"`
	Status      *Category `json:"status,omitempty"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	AssigneeID  *int      `json:"assignee_id"`
	ReporterID  *int      `json:"reporter_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
