package entity

import "time"

// Holiday represents a global holiday (state or company level)
type Holiday struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Type        string    `json:"type" binding:"required"` // 'state' or 'company'
	TypeID      uint      `json:"type_id"`
	TypeCat     *Category `json:"type_cat,omitempty"`
	IsRecurring bool      `json:"is_recurring"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
