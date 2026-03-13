package entity

import "time"

type ProjectSnapshot struct {
	ID         int       `json:"id"`
	ProjectID  int       `json:"project_id"`
	SPI        float64   `json:"spi"`
	CPI        float64   `json:"cpi"`
	EV         float64   `json:"ev"`
	AC         float64   `json:"ac"`
	PV         float64   `json:"pv"`
	Progress   int       `json:"progress"`
	CapturedAt time.Time `json:"captured_at"`
}

type MilestoneSnapshot struct {
	ID            int        `json:"id"`
	ProjectID     int        `json:"project_id"`
	NodeID        int        `json:"node_id"`
	MilestoneName string     `json:"milestone_name"`
	PlannedDate   time.Time  `json:"planned_date"`
	ActualDate    *time.Time `json:"actual_date"`
	CapturedAt    time.Time  `json:"captured_at"`
}
