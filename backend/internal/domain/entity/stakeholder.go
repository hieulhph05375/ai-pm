package entity

import "time"

// Stakeholder represents an organization-wide stakeholder
type Stakeholder struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Role         string    `json:"role"`
	RoleID       uint      `json:"role_id"`
	RoleCat      *Category `json:"role_cat,omitempty"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Organization string    `json:"organization"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProjectStakeholder represents a stakeholder mapped to a specific project
type ProjectStakeholder struct {
	ProjectID     int          `json:"project_id"`
	StakeholderID int          `json:"stakeholder_id"`
	ProjectRole   string       `json:"project_role"` // Role specific to this project
	RoleID        *uint        `json:"role_id"`
	RoleCat       *Category    `json:"role_cat,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	Stakeholder   *Stakeholder `json:"stakeholder,omitempty"` // populated if joined
}
