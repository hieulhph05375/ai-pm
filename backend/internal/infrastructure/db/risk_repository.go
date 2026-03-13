package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
)

type riskRepo struct {
	db *sql.DB
}

func NewRiskRepository(db *sql.DB) *riskRepo {
	return &riskRepo{db: db}
}

func (r *riskRepo) Create(ctx context.Context, risk *entity.Risk) error {
	query := `
		INSERT INTO risks (project_id, title, description, probability, impact, status, status_id, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		risk.ProjectID, risk.Title, risk.Description, risk.Probability, risk.Impact, risk.Status, risk.StatusID, risk.OwnerID,
	).Scan(&risk.ID, &risk.CreatedAt, &risk.UpdatedAt)
}

func (r *riskRepo) ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Risk, error) {
	query := `
		SELECT r.id, r.project_id, r.title, r.description, r.probability, r.impact,
		       r.probability * r.impact AS risk_score, r.status, r.status_id, r.owner_id, r.created_at, r.updated_at,
			   s.name as status_name, s.color as status_color, s.icon as status_icon
		FROM risks r
		LEFT JOIN categories s ON r.status_id = s.id
		WHERE r.project_id = $1
		ORDER BY r.probability * r.impact DESC, r.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	risks := []entity.Risk{}
	for rows.Next() {
		var risk entity.Risk
		var statusID sql.NullInt64
		var statusName, statusColor, statusIcon sql.NullString

		if err := rows.Scan(
			&risk.ID, &risk.ProjectID, &risk.Title, &risk.Description,
			&risk.Probability, &risk.Impact, &risk.RiskScore,
			&risk.Status, &statusID, &risk.OwnerID, &risk.CreatedAt, &risk.UpdatedAt,
			&statusName, &statusColor, &statusIcon,
		); err != nil {
			return nil, err
		}

		if statusID.Valid {
			risk.StatusID = uint(statusID.Int64)
			risk.StatusCat = &entity.Category{
				ID:    risk.StatusID,
				Name:  statusName.String,
				Color: nullStringPtr(statusColor),
				Icon:  nullStringPtr(statusIcon),
			}
		}

		risks = append(risks, risk)
	}
	return risks, nil
}

func (r *riskRepo) CountByProject(ctx context.Context, projectID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM risks WHERE project_id = $1", projectID).Scan(&count)
	return count, err
}

func (r *riskRepo) Update(ctx context.Context, risk *entity.Risk) error {
	query := `
		UPDATE risks
		SET title=$1, description=$2, probability=$3, impact=$4, status=$5, status_id=$6, owner_id=$7, updated_at=NOW()
		WHERE id=$8`
	_, err := r.db.ExecContext(ctx, query,
		risk.Title, risk.Description, risk.Probability, risk.Impact, risk.Status, risk.StatusID, risk.OwnerID, risk.ID,
	)
	return err
}

func (r *riskRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM risks WHERE id=$1", id)
	return err
}
