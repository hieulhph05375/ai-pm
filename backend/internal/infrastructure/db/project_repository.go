package db

import (
	"context"
	"database/sql"
	"fmt"
	"project-mgmt/backend/internal/domain/entity"
	"strings"
)

type projectRepo struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *projectRepo {
	return &projectRepo{db: db}
}

func (r *projectRepo) scanProject(rows interface {
	Scan(dest ...interface{}) error
}, p *entity.Project) error {
	var projectStatusID, currentPhaseID, portfolioCategoryID, overallHealthID, priorityLevelID sql.NullInt64
	var statusName, statusColor, statusIcon sql.NullString
	var phaseName, phaseColor, phaseIcon sql.NullString
	var portfolioName, portfolioColor, portfolioIcon sql.NullString
	var healthName, healthColor, healthIcon sql.NullString
	var priorityName, priorityColor, priorityIcon sql.NullString

	err := rows.Scan(
		&p.ID, &p.ProjectID, &p.ProjectName, &p.Description, &p.ProjectManager, &p.Sponsor, &p.RequestingDepartment,
		&p.CurrentPhase, &p.ProjectStatus, &p.StrategicGoal, &p.PortfolioCategory, &p.StrategicScore, &p.PriorityLevel,
		&projectStatusID, &currentPhaseID, &portfolioCategoryID, &overallHealthID, &priorityLevelID,
		&p.ApprovedBudget, &p.ActualCost, &p.EAC, &p.CapexOpexRatio, &p.ExpectedROI, &p.PaybackPeriod, &p.BenefitRealizationDate,
		&p.PlannedStartDate, &p.ActualStartDate, &p.PlannedEndDate, &p.ActualEndDate, &p.Progress, &p.OverallHealth,
		&p.SPI, &p.CPI, &p.LastExecutiveSummary, &p.EstimatedEffort, &p.ActualEffort, &p.ResourceRiskFlag,
		&p.MissingSkills, &p.SystemicRiskLevel, &p.OpenCriticalRisks, &p.ComplianceImpact, &p.DependenciesSummary,
		&p.CreatedAt, &p.UpdatedAt, &p.LastReminderAt,
		&statusName, &statusColor, &statusIcon,
		&phaseName, &phaseColor, &phaseIcon,
		&portfolioName, &portfolioColor, &portfolioIcon,
		&healthName, &healthColor, &healthIcon,
		&priorityName, &priorityColor, &priorityIcon,
	)
	if err != nil {
		return err
	}

	if projectStatusID.Valid {
		p.ProjectStatusID = uintPtr(uint(projectStatusID.Int64))
		p.ProjectStatusCat = &entity.Category{ID: *p.ProjectStatusID, Name: statusName.String, Color: nullStringPtr(statusColor), Icon: nullStringPtr(statusIcon)}
	}
	if currentPhaseID.Valid {
		p.CurrentPhaseID = uintPtr(uint(currentPhaseID.Int64))
		p.CurrentPhaseCat = &entity.Category{ID: *p.CurrentPhaseID, Name: phaseName.String, Color: nullStringPtr(phaseColor), Icon: nullStringPtr(phaseIcon)}
	}
	if portfolioCategoryID.Valid {
		p.PortfolioCategoryID = uintPtr(uint(portfolioCategoryID.Int64))
		p.PortfolioCategoryCat = &entity.Category{ID: *p.PortfolioCategoryID, Name: portfolioName.String, Color: nullStringPtr(portfolioColor), Icon: nullStringPtr(portfolioIcon)}
	}
	if overallHealthID.Valid {
		p.OverallHealthID = uintPtr(uint(overallHealthID.Int64))
		p.OverallHealthCat = &entity.Category{ID: *p.OverallHealthID, Name: healthName.String, Color: nullStringPtr(healthColor), Icon: nullStringPtr(healthIcon)}
	}
	if priorityLevelID.Valid {
		p.PriorityLevelID = uintPtr(uint(priorityLevelID.Int64))
		p.PriorityLevelCat = &entity.Category{ID: *p.PriorityLevelID, Name: priorityName.String, Color: nullStringPtr(priorityColor), Icon: nullStringPtr(priorityIcon)}
	}

	return nil
}

func uintPtr(u uint) *uint {
	return &u
}

func (r *projectRepo) Create(ctx context.Context, p *entity.Project) error {
	query := `
		INSERT INTO projects (
			project_id, project_name, description, project_manager, sponsor, requesting_department,
			current_phase, project_status, strategic_goal, portfolio_category, strategic_score, priority_level,
			project_status_id, current_phase_id, portfolio_category_id, overall_health_id, priority_level_id,
			approved_budget, actual_cost, eac, capex_opex_ratio, expected_roi, payback_period, benefit_realization_date,
			planned_start_date, actual_start_date, planned_end_date, actual_end_date, progress, overall_health,
			spi, cpi, last_executive_summary, estimated_effort, actual_effort, resource_risk_flag,
			missing_skills, systemic_risk_level, open_critical_risks, compliance_impact, dependencies_summary, last_reminder_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		p.ProjectID, p.ProjectName, p.Description, p.ProjectManager, p.Sponsor, p.RequestingDepartment,
		p.CurrentPhase, p.ProjectStatus, p.StrategicGoal, p.PortfolioCategory, p.StrategicScore, p.PriorityLevel,
		p.ProjectStatusID, p.CurrentPhaseID, p.PortfolioCategoryID, p.OverallHealthID, p.PriorityLevelID,
		p.ApprovedBudget, p.ActualCost, p.EAC, p.CapexOpexRatio, p.ExpectedROI, p.PaybackPeriod, p.BenefitRealizationDate,
		p.PlannedStartDate, p.ActualStartDate, p.PlannedEndDate, p.ActualEndDate, p.Progress, p.OverallHealth,
		p.SPI, p.CPI, p.LastExecutiveSummary, p.EstimatedEffort, p.ActualEffort, p.ResourceRiskFlag,
		p.MissingSkills, p.SystemicRiskLevel, p.OpenCriticalRisks, p.ComplianceImpact, p.DependenciesSummary, p.LastReminderAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *projectRepo) GetByID(ctx context.Context, id int, userID int, isAdmin bool) (*entity.Project, error) {
	p := &entity.Project{}
	baseQuery := `
		SELECT 
			p.id, p.project_id, p.project_name, p.description, p.project_manager, p.sponsor, p.requesting_department,
			p.current_phase, p.project_status, p.strategic_goal, p.portfolio_category, p.strategic_score, p.priority_level,
			p.project_status_id, p.current_phase_id, p.portfolio_category_id, p.overall_health_id, p.priority_level_id,
			p.approved_budget, p.actual_cost, p.eac, p.capex_opex_ratio, p.expected_roi, p.payback_period, p.benefit_realization_date,
			p.planned_start_date, p.actual_start_date, p.planned_end_date, p.actual_end_date, p.progress, p.overall_health,
			p.spi, p.cpi, p.last_executive_summary, p.estimated_effort, p.actual_effort, p.resource_risk_flag,
			p.missing_skills, p.systemic_risk_level, p.open_critical_risks, p.compliance_impact, p.dependencies_summary,
			p.created_at, p.updated_at, p.last_reminder_at,
			s.name as status_name, s.color as status_color, s.icon as status_icon,
			ph.name as phase_name, ph.color as phase_color, ph.icon as phase_icon,
			pc.name as portfolio_name, pc.color as portfolio_color, pc.icon as portfolio_icon,
			h.name as health_name, h.color as health_color, h.icon as health_icon,
			pr.name as priority_name, pr.color as priority_color, pr.icon as priority_icon
		FROM projects p
		LEFT JOIN categories s ON p.project_status_id = s.id
		LEFT JOIN categories ph ON p.current_phase_id = ph.id
		LEFT JOIN categories pc ON p.portfolio_category_id = pc.id
		LEFT JOIN categories h ON p.overall_health_id = h.id
		LEFT JOIN categories pr ON p.priority_level_id = pr.id`

	var query string
	var args []interface{}
	if isAdmin {
		query = baseQuery + ` WHERE p.id = $1`
		args = []interface{}{id}
	} else {
		query = baseQuery + `
			JOIN project_members pm ON p.id = pm.project_id
			WHERE p.id = $1 AND pm.user_id = $2`
		args = []interface{}{id, userID}
	}
	if err := r.scanProject(r.db.QueryRowContext(ctx, query, args...), p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *projectRepo) GetByProjectID(ctx context.Context, projectID string, userID int, isAdmin bool) (*entity.Project, error) {
	p := &entity.Project{}
	baseQuery := `
		SELECT 
			p.id, p.project_id, p.project_name, p.description, p.project_manager, p.sponsor, p.requesting_department,
			p.current_phase, p.project_status, p.strategic_goal, p.portfolio_category, p.strategic_score, p.priority_level,
			p.project_status_id, p.current_phase_id, p.portfolio_category_id, p.overall_health_id, p.priority_level_id,
			p.approved_budget, p.actual_cost, p.eac, p.capex_opex_ratio, p.expected_roi, p.payback_period, p.benefit_realization_date,
			p.planned_start_date, p.actual_start_date, p.planned_end_date, p.actual_end_date, p.progress, p.overall_health,
			p.spi, p.cpi, p.last_executive_summary, p.estimated_effort, p.actual_effort, p.resource_risk_flag,
			p.missing_skills, p.systemic_risk_level, p.open_critical_risks, p.compliance_impact, p.dependencies_summary,
			p.created_at, p.updated_at, p.last_reminder_at,
			s.name as status_name, s.color as status_color, s.icon as status_icon,
			ph.name as phase_name, ph.color as phase_color, ph.icon as phase_icon,
			pc.name as portfolio_name, pc.color as portfolio_color, pc.icon as portfolio_icon,
			h.name as health_name, h.color as health_color, h.icon as health_icon,
			pr.name as priority_name, pr.color as priority_color, pr.icon as priority_icon
		FROM projects p
		LEFT JOIN categories s ON p.project_status_id = s.id
		LEFT JOIN categories ph ON p.current_phase_id = ph.id
		LEFT JOIN categories pc ON p.portfolio_category_id = pc.id
		LEFT JOIN categories h ON p.overall_health_id = h.id
		LEFT JOIN categories pr ON p.priority_level_id = pr.id`

	var query string
	var args []interface{}
	if isAdmin {
		query = baseQuery + ` WHERE p.project_id = $1`
		args = []interface{}{projectID}
	} else {
		query = baseQuery + `
			JOIN project_members pm ON p.id = pm.project_id
			WHERE p.project_id = $1 AND pm.user_id = $2`
		args = []interface{}{projectID, userID}
	}
	if err := r.scanProject(r.db.QueryRowContext(ctx, query, args...), p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *projectRepo) Update(ctx context.Context, p *entity.Project) error {
	query := `
		UPDATE projects SET
			project_id=$1, project_name=$2, description=$3, project_manager=$4, sponsor=$5, requesting_department=$6,
			current_phase=$7, project_status=$8, strategic_goal=$9, portfolio_category=$10, strategic_score=$11, priority_level=$12,
			project_status_id=$13, current_phase_id=$14, portfolio_category_id=$15, overall_health_id=$16, priority_level_id=$17,
			approved_budget=$18, actual_cost=$19, eac=$20, capex_opex_ratio=$21, expected_roi=$22, payback_period=$23, benefit_realization_date=$24,
			planned_start_date=$25, actual_start_date=$26, planned_end_date=$27, actual_end_date=$28, progress=$29, overall_health=$30,
			spi=$31, cpi=$32, last_executive_summary=$33, estimated_effort=$34, actual_effort=$35, resource_risk_flag=$36,
			missing_skills=$37, systemic_risk_level=$38, open_critical_risks=$39, compliance_impact=$40, dependencies_summary=$41, last_reminder_at=$42,
			updated_at=CURRENT_TIMESTAMP
		WHERE id=$43`

	_, err := r.db.ExecContext(ctx, query,
		p.ProjectID, p.ProjectName, p.Description, p.ProjectManager, p.Sponsor, p.RequestingDepartment,
		p.CurrentPhase, p.ProjectStatus, p.StrategicGoal, p.PortfolioCategory, p.StrategicScore, p.PriorityLevel,
		p.ProjectStatusID, p.CurrentPhaseID, p.PortfolioCategoryID, p.OverallHealthID, p.PriorityLevelID,
		p.ApprovedBudget, p.ActualCost, p.EAC, p.CapexOpexRatio, p.ExpectedROI, p.PaybackPeriod, p.BenefitRealizationDate,
		p.PlannedStartDate, p.ActualStartDate, p.PlannedEndDate, p.ActualEndDate, p.Progress, p.OverallHealth,
		p.SPI, p.CPI, p.LastExecutiveSummary, p.EstimatedEffort, p.ActualEffort, p.ResourceRiskFlag,
		p.MissingSkills, p.SystemicRiskLevel, p.OpenCriticalRisks, p.ComplianceImpact, p.DependenciesSummary, p.LastReminderAt,
		p.ID,
	)
	return err
}

func (r *projectRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id=$1", id)
	return err
}

func (r *projectRepo) buildListQuery(search, status string) (string, []interface{}) {
	var where []string
	var args []interface{}
	argIdx := 1
	if search != "" {
		where = append(where, fmt.Sprintf("(p.project_name ILIKE $%d OR p.project_id ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if status != "" {
		where = append(where, fmt.Sprintf("p.project_status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}
	return whereClause, args
}

func (r *projectRepo) List(ctx context.Context, offset, limit int, search string, status string, userID int, isAdmin bool) ([]entity.Project, int, error) {
	whereClause, args := r.buildListQuery(search, status)

	countQuery := ""
	listQuery := ""

	if isAdmin {
		countQuery = fmt.Sprintf("SELECT COUNT(*) FROM projects p %s", whereClause)
	} else {
		membershipFilter := fmt.Sprintf("p.id IN (SELECT project_id FROM project_members WHERE user_id = $%d)", len(args)+1)
		if whereClause == "" {
			whereClause = "WHERE " + membershipFilter
		} else {
			whereClause += " AND " + membershipFilter
		}
		countQuery = fmt.Sprintf("SELECT COUNT(*) FROM projects p %s", whereClause)
		args = append(args, userID)
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	argIdx := len(args) + 1
	listQuery = fmt.Sprintf(`
		SELECT 
			p.id, p.project_id, p.project_name, p.description, p.project_manager, p.sponsor, p.requesting_department,
			p.current_phase, p.project_status, p.strategic_goal, p.portfolio_category, p.strategic_score, p.priority_level,
			p.project_status_id, p.current_phase_id, p.portfolio_category_id, p.overall_health_id, p.priority_level_id,
			p.approved_budget, p.actual_cost, p.eac, p.capex_opex_ratio, p.expected_roi, p.payback_period, p.benefit_realization_date,
			p.planned_start_date, p.actual_start_date, p.planned_end_date, p.actual_end_date, p.progress, p.overall_health,
			p.spi, p.cpi, p.last_executive_summary, p.estimated_effort, p.actual_effort, p.resource_risk_flag,
			p.missing_skills, p.systemic_risk_level, p.open_critical_risks, p.compliance_impact, p.dependencies_summary,
			p.created_at, p.updated_at, p.last_reminder_at,
			s.name as status_name, s.color as status_color, s.icon as status_icon,
			ph.name as phase_name, ph.color as phase_color, ph.icon as phase_icon,
			pc.name as portfolio_name, pc.color as portfolio_color, pc.icon as portfolio_icon,
			h.name as health_name, h.color as health_color, h.icon as health_icon,
			pr.name as priority_name, pr.color as priority_color, pr.icon as priority_icon
		FROM projects p
		LEFT JOIN categories s ON p.project_status_id = s.id
		LEFT JOIN categories ph ON p.current_phase_id = ph.id
		LEFT JOIN categories pc ON p.portfolio_category_id = pc.id
		LEFT JOIN categories h ON p.overall_health_id = h.id
		LEFT JOIN categories pr ON p.priority_level_id = pr.id
		%s ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)
	rows, err := r.db.QueryContext(ctx, listQuery, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	projects := []entity.Project{}
	for rows.Next() {
		var p entity.Project
		if err := r.scanProject(rows, &p); err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}
	return projects, total, nil
}

func (r *projectRepo) GetPortfolioOverview(ctx context.Context, userID int, isAdmin bool) (*entity.PortfolioOverview, error) {
	overview := &entity.PortfolioOverview{}

	query := `
		SELECT
			COUNT(*) as total_projects,
			COUNT(*) FILTER (WHERE s.name = 'Active') as active_projects,
			COUNT(*) FILTER (WHERE s.name = 'Completed') as completed_projects,
			COUNT(*) FILTER (WHERE s.name = 'On Hold') as on_hold_projects,
			COALESCE(SUM(p.approved_budget), 0) as total_budget,
			COUNT(*) FILTER (WHERE h.name = 'Green') as green_projects,
			COUNT(*) FILTER (WHERE h.name = 'Yellow') as yellow_projects,
			COUNT(*) FILTER (WHERE h.name = 'Red') as red_projects
		FROM projects p
		LEFT JOIN categories s ON p.project_status_id = s.id
		LEFT JOIN categories h ON p.overall_health_id = h.id
	`

	if !isAdmin {
		query += " WHERE p.id IN (SELECT project_id FROM project_members WHERE user_id = $1)"
		err := r.db.QueryRowContext(ctx, query, userID).Scan(
			&overview.TotalProjects, &overview.ActiveProjects, &overview.CompletedProjects, &overview.OnHoldProjects,
			&overview.TotalBudget, &overview.GreenProjects, &overview.YellowProjects, &overview.RedProjects,
		)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.QueryRowContext(ctx, query).Scan(
			&overview.TotalProjects, &overview.ActiveProjects, &overview.CompletedProjects, &overview.OnHoldProjects,
			&overview.TotalBudget, &overview.GreenProjects, &overview.YellowProjects, &overview.RedProjects,
		)
		if err != nil {
			return nil, err
		}
	}

	highRiskQuery := `
		SELECT 
			p.id, p.project_id, p.project_name, p.description, p.project_manager, p.sponsor, p.requesting_department,
			p.current_phase, p.project_status, p.strategic_goal, p.portfolio_category, p.strategic_score, p.priority_level,
			p.project_status_id, p.current_phase_id, p.portfolio_category_id, p.overall_health_id, p.priority_level_id,
			p.approved_budget, p.actual_cost, p.eac, p.capex_opex_ratio, p.expected_roi, p.payback_period, p.benefit_realization_date,
			p.planned_start_date, p.actual_start_date, p.planned_end_date, p.actual_end_date, p.progress, p.overall_health,
			p.spi, p.cpi, p.last_executive_summary, p.estimated_effort, p.actual_effort, p.resource_risk_flag,
			p.missing_skills, p.systemic_risk_level, p.open_critical_risks, p.compliance_impact, p.dependencies_summary,
			p.created_at, p.updated_at, p.last_reminder_at,
			s.name as status_name, s.color as status_color, s.icon as status_icon,
			ph.name as phase_name, ph.color as phase_color, ph.icon as phase_icon,
			pc.name as portfolio_name, pc.color as portfolio_color, pc.icon as portfolio_icon,
			h.name as health_name, h.color as health_color, h.icon as health_icon,
			pr.name as priority_name, pr.color as priority_color, pr.icon as priority_icon
		FROM projects p
		LEFT JOIN categories s ON p.project_status_id = s.id
		LEFT JOIN categories ph ON p.current_phase_id = ph.id
		LEFT JOIN categories pc ON p.portfolio_category_id = pc.id
		LEFT JOIN categories h ON p.overall_health_id = h.id
		LEFT JOIN categories pr ON p.priority_level_id = pr.id
		WHERE h.name = 'Red'`

	var hArgs []interface{}
	if !isAdmin {
		highRiskQuery += " AND p.id IN (SELECT project_id FROM project_members WHERE user_id = $1)"
		hArgs = append(hArgs, userID)
	}
	highRiskQuery += " ORDER BY p.created_at DESC LIMIT 10"

	rows, err := r.db.QueryContext(ctx, highRiskQuery, hArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	overview.HighRiskProjects = []entity.Project{}
	for rows.Next() {
		var p entity.Project
		if err := r.scanProject(rows, &p); err != nil {
			return nil, err
		}
		overview.HighRiskProjects = append(overview.HighRiskProjects, p)
	}

	return overview, nil
}
