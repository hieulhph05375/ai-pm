package entity

// ResourceWorkloadEntry aggregates workload data for a user on a single day
type ResourceWorkloadEntry struct {
	UserID         int     `json:"user_id"`
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	Date           string  `json:"date"` // YYYY-MM-DD
	TaskCount      int     `json:"task_count"`
	TotalHours     float64 `json:"total_hours"`     // estimated hours per day
	LoadPercentage float64 `json:"load_percentage"` // e.g., (TotalHours / 8.0) * 100
}

// ResourceWorkload groups workload entries by user
type ResourceWorkload struct {
	UserID   int                     `json:"user_id"`
	FullName string                  `json:"full_name"`
	Email    string                  `json:"email"`
	Role     string                  `json:"role"`
	Entries  []ResourceWorkloadEntry `json:"entries"`
}

// WorkloadOverview is the top-level response
type WorkloadOverview struct {
	StartDate string             `json:"start_date"`
	EndDate   string             `json:"end_date"`
	Users     []ResourceWorkload `json:"users"`
}
