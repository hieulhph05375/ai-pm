package entity

import "time"

type ProjectMember struct {
	ID            int          `json:"id"`
	ProjectID     int          `json:"project_id"`
	UserID        int          `json:"user_id"`
	ProjectRoleID int          `json:"project_role_id"`
	JoinedAt      time.Time    `json:"joined_at"`
	CreatedAt     time.Time    `json:"created_at"`
	User          *User        `json:"user,omitempty"`
	Project       *Project     `json:"project,omitempty"`
	Role          *ProjectRole `json:"role,omitempty"`
}
