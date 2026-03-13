package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context) ([]entity.User, error)
	Search(ctx context.Context, query string) ([]entity.User, error)
	ListWithPagination(ctx context.Context, query string, offset, limit int) ([]entity.User, int, error)
	Count(ctx context.Context, query string) (int, error)
}
