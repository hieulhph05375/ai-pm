package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
)

type resourceRepo struct {
	db *sql.DB
}

func NewResourceRepository(db *sql.DB) *resourceRepo {
	return &resourceRepo{db: db}
}

func (r *resourceRepo) GetWorkload(ctx context.Context, startDate, endDate string) (*entity.WorkloadOverview, error) {
	// Query tasks assigned to users that overlap with the requested date range
	// Distribute hours evenly across the task duration days
	query := `
		WITH date_series AS (
			SELECT generate_series(
				GREATEST(t.start_date, $1::date),
				LEAST(COALESCE(t.due_date, t.start_date), $2::date),
				INTERVAL '1 day'
			)::date AS work_date
			FROM tasks t
			WHERE t.assignee_id IS NOT NULL
			  AND t.start_date IS NOT NULL
			  AND t.start_date <= $2::date
			  AND COALESCE(t.due_date, t.start_date) >= $1::date
		),
		task_days AS (
			SELECT
				t.assignee_id,
				generate_series(
					GREATEST(t.start_date, $1::date),
					LEAST(COALESCE(t.due_date, t.start_date), $2::date),
					INTERVAL '1 day'
				)::date AS work_date,
				8.0 / GREATEST(
					DATE_PART('day', COALESCE(t.due_date, t.start_date) - t.start_date) + 1,
					1
				) AS hours_per_day,
				t.id AS task_id
			FROM tasks t
			WHERE t.assignee_id IS NOT NULL
			  AND t.start_date IS NOT NULL
			  AND t.start_date <= $2::date
			  AND COALESCE(t.due_date, t.start_date) >= $1::date
		)
		SELECT
			u.id AS user_id,
			u.full_name,
			u.email,
			r.name AS role_name,
			td.work_date::text AS date,
			COUNT(DISTINCT td.task_id) AS task_count,
			ROUND(SUM(td.hours_per_day)::numeric, 2) AS total_hours,
			ROUND((SUM(td.hours_per_day) / 8.0 * 100.0)::numeric, 2) AS load_percentage
		FROM task_days td
		JOIN users u ON u.id = td.assignee_id
		LEFT JOIN roles r ON u.role_id = r.id
		GROUP BY u.id, u.full_name, u.email, r.name, td.work_date
		ORDER BY u.full_name, td.work_date
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userMap := map[int]*entity.ResourceWorkload{}
	var userOrder []int

	for rows.Next() {
		var entry entity.ResourceWorkloadEntry
		var roleName sql.NullString
		if err := rows.Scan(&entry.UserID, &entry.FullName, &entry.Email, &roleName, &entry.Date, &entry.TaskCount, &entry.TotalHours, &entry.LoadPercentage); err != nil {
			return nil, err
		}
		if roleName.Valid {
			entry.Role = roleName.String
		}

		if _, exists := userMap[entry.UserID]; !exists {
			userMap[entry.UserID] = &entity.ResourceWorkload{
				UserID:   entry.UserID,
				FullName: entry.FullName,
				Email:    entry.Email,
				Role:     entry.Role,
				Entries:  []entity.ResourceWorkloadEntry{},
			}
			userOrder = append(userOrder, entry.UserID)
		}
		userMap[entry.UserID].Entries = append(userMap[entry.UserID].Entries, entry)
	}

	users := make([]entity.ResourceWorkload, 0, len(userOrder))
	for _, uid := range userOrder {
		users = append(users, *userMap[uid])
	}

	return &entity.WorkloadOverview{
		StartDate: startDate,
		EndDate:   endDate,
		Users:     users,
	}, nil
}
