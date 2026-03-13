package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type snapshotRepository struct {
	db *sql.DB
}

func NewSnapshotRepository(db *sql.DB) repository.SnapshotRepository {
	return &snapshotRepository{db: db}
}

func (r *snapshotRepository) CreateProjectSnapshot(ctx context.Context, s *entity.ProjectSnapshot) error {
	query := `INSERT INTO project_snapshots (project_id, spi, cpi, ev, ac, pv, progress, captured_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id, captured_at`
	return r.db.QueryRowContext(ctx, query, s.ProjectID, s.SPI, s.CPI, s.EV, s.AC, s.PV, s.Progress).Scan(&s.ID, &s.CapturedAt)
}

func (r *snapshotRepository) CreateMilestoneSnapshot(ctx context.Context, s *entity.MilestoneSnapshot) error {
	query := `INSERT INTO milestone_snapshots (project_id, node_id, milestone_name, planned_date, actual_date, captured_at)
	          VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id, captured_at`
	return r.db.QueryRowContext(ctx, query, s.ProjectID, s.NodeID, s.MilestoneName, s.PlannedDate, s.ActualDate).Scan(&s.ID, &s.CapturedAt)
}

func (r *snapshotRepository) GetProjectSnapshots(ctx context.Context, projectID int) ([]entity.ProjectSnapshot, error) {
	query := `SELECT id, project_id, spi, cpi, ev, ac, pv, progress, captured_at 
	          FROM project_snapshots WHERE project_id = $1 ORDER BY captured_at ASC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []entity.ProjectSnapshot
	for rows.Next() {
		var s entity.ProjectSnapshot
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.SPI, &s.CPI, &s.EV, &s.AC, &s.PV, &s.Progress, &s.CapturedAt); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (r *snapshotRepository) GetMilestoneSnapshots(ctx context.Context, projectID int) ([]entity.MilestoneSnapshot, error) {
	query := `SELECT id, project_id, node_id, milestone_name, planned_date, actual_date, captured_at 
	          FROM milestone_snapshots WHERE project_id = $1 ORDER BY captured_at ASC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []entity.MilestoneSnapshot
	for rows.Next() {
		var s entity.MilestoneSnapshot
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.NodeID, &s.MilestoneName, &s.PlannedDate, &s.ActualDate, &s.CapturedAt); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}
