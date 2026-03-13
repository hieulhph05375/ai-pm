package entity

import "time"

type Category struct {
	ID          uint          `json:"id"`
	TypeID      uint          `json:"type_id"`
	ParentID    *uint         `json:"parent_id"`
	Type        *CategoryType `json:"type,omitempty"`
	Parent      *Category     `json:"parent,omitempty"`
	Name        string        `json:"name"`
	Color       *string       `json:"color"`
	Icon        *string       `json:"icon"`
	Description *string       `json:"description"`
	IsActive    *bool         `json:"is_active"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
