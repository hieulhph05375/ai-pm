package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *entity.Notification) error
	ListByUser(ctx context.Context, userID, limit, offset int) ([]entity.Notification, int, error)
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	MarkRead(ctx context.Context, id int) error
	MarkAllRead(ctx context.Context, userID int) error
}
