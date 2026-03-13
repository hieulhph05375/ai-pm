package entity

type PortfolioOverview struct {
	TotalProjects     int       `json:"total_projects"`
	ActiveProjects    int       `json:"active_projects"`
	CompletedProjects int       `json:"completed_projects"`
	OnHoldProjects    int       `json:"on_hold_projects"`
	TotalBudget       float64   `json:"total_budget"`
	GreenProjects     int       `json:"green_projects"`
	YellowProjects    int       `json:"yellow_projects"`
	RedProjects       int       `json:"red_projects"`
	HighRiskProjects  []Project `json:"high_risk_projects"`
}
