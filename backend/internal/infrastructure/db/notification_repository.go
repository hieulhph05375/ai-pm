package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
)

type notificationRepo struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *notificationRepo {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(ctx context.Context, n *entity.Notification) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO notifications (user_id, type, title, body, ref_id, ref_type)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		n.UserID, n.Type, n.Title, n.Body, n.RefID, n.RefType,
	).Scan(&n.ID, &n.CreatedAt)
}

func (r *notificationRepo) ListByUser(ctx context.Context, userID, limit, offset int) ([]entity.Notification, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id=$1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, type, title, body, ref_id, ref_type, is_read, created_at
		 FROM notifications WHERE user_id=$1
		 ORDER BY is_read ASC, created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []entity.Notification
	for rows.Next() {
		var n entity.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.RefID, &n.RefType, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}
	return notifications, total, nil
}

func (r *notificationRepo) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=FALSE`,
		userID,
	).Scan(&count)
	return count, err
}

func (r *notificationRepo) MarkRead(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read=TRUE WHERE id=$1`, id)
	return err
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET is_read=TRUE WHERE user_id=$1`, userID)
	return err
}
