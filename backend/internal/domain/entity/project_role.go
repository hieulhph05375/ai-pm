package entity

import "time"

type ProjectPermission struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Module      string    `json:"module"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectRole struct {
	ID          uint                `json:"id"`
	ProjectID   uint                `json:"project_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Color       string              `json:"color"`
	IsDefault   bool                `json:"is_default"`
	Permissions []ProjectPermission `json:"permissions,omitempty" gorm:"many2many:project_role_permissions"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
