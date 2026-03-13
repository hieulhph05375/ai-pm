package entity

import (
	"time"
)

type Project struct {
	ID                   int     `json:"id"`
	ProjectID            string  `json:"project_id" binding:"required"`
	ProjectName          string  `json:"project_name" binding:"required"`
	Description          *string `json:"description"`
	ProjectManager       *string `json:"project_manager"`
	Sponsor              *string `json:"sponsor"`
	RequestingDepartment *string `json:"requesting_department"`
	CurrentPhase         *string `json:"current_phase"`
	ProjectStatus        *string `json:"project_status"`
	StrategicGoal        *string `json:"strategic_goal"`
	PortfolioCategory    *string `json:"portfolio_category"`
	StrategicScore       int     `json:"strategic_score"`
	PriorityLevel        *string `json:"priority_level"`

	// Category IDs (Phase 23)
	ProjectStatusID     *uint `json:"project_status_id"`
	CurrentPhaseID      *uint `json:"current_phase_id"`
	PortfolioCategoryID *uint `json:"portfolio_category_id"`
	OverallHealthID     *uint `json:"overall_health_id"`
	PriorityLevelID     *uint `json:"priority_level_id"`

	// Joined Category Objects
	ProjectStatusCat       *Category  `json:"status,omitempty"`
	CurrentPhaseCat        *Category  `json:"phase,omitempty"`
	PortfolioCategoryCat   *Category  `json:"portfolio,omitempty"`
	OverallHealthCat       *Category  `json:"health,omitempty"`
	PriorityLevelCat       *Category  `json:"priority,omitempty"`
	ApprovedBudget         float64    `json:"approved_budget"`
	ActualCost             float64    `json:"actual_cost"`
	EAC                    float64    `json:"eac"`
	CapexOpexRatio         *string    `json:"capex_opex_ratio"`
	ExpectedROI            float64    `json:"expected_roi"`
	PaybackPeriod          int        `json:"payback_period"`
	BenefitRealizationDate *time.Time `json:"benefit_realization_date"`
	PlannedStartDate       *time.Time `json:"planned_start_date"`
	ActualStartDate        *time.Time `json:"actual_start_date"`
	PlannedEndDate         *time.Time `json:"planned_end_date"`
	ActualEndDate          *time.Time `json:"actual_end_date"`
	Progress               int        `json:"progress"`
	OverallHealth          *string    `json:"overall_health"`
	SPI                    float64    `json:"spi"`
	CPI                    float64    `json:"cpi"`
	LastExecutiveSummary   *string    `json:"last_executive_summary"`
	EstimatedEffort        int        `json:"estimated_effort"`
	ActualEffort           int        `json:"actual_effort"`
	ResourceRiskFlag       bool       `json:"resource_risk_flag"`
	MissingSkills          *string    `json:"missing_skills"`
	SystemicRiskLevel      *string    `json:"systemic_risk_level"`
	OpenCriticalRisks      int        `json:"open_critical_risks"`
	ComplianceImpact       *string    `json:"compliance_impact"`
	DependenciesSummary    *string    `json:"dependencies_summary"`
	LastReminderAt         *time.Time `json:"last_reminder_at"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}
