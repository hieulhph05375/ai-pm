package entity

import "time"

type User struct {
	ID             uint       `json:"id"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"-"`
	FullName       string     `json:"full_name"`
	IsActive       bool       `json:"is_active"`
	RoleID         uint       `json:"role_id"`
	IsAdmin        bool       `json:"is_admin"`
	Role           Role       `json:"role,omitempty"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Role struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Permission struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
