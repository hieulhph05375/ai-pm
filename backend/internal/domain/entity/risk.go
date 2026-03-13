package entity

import "time"

type Risk struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Probability int       `json:"probability"` // 1-5
	Impact      int       `json:"impact"`      // 1-5
	RiskScore   int       `json:"risk_score"`  // computed: probability * impact
	Status      string    `json:"status"`      // Open, Mitigated, Closed
	StatusID    uint      `json:"status_id"`
	StatusCat   *Category `json:"status_cat,omitempty"`
	OwnerID     *int      `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
