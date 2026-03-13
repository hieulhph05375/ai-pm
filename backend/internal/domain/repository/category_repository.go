package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type CategoryRepository interface {
	// CategoryType CRUD
	CreateType(ctx context.Context, ct *entity.CategoryType) error
	GetTypeByID(ctx context.Context, id uint) (*entity.CategoryType, error)
	ListTypes(ctx context.Context, search string, limit, offset int) ([]entity.CategoryType, int, error)
	UpdateType(ctx context.Context, ct *entity.CategoryType) error
	DeleteType(ctx context.Context, id uint) error

	// Category CRUD
	CreateCategory(ctx context.Context, c *entity.Category) error
	GetCategoryByID(ctx context.Context, id uint) (*entity.Category, error)
	ListCategories(ctx context.Context, typeID *uint, search string, limit, offset int) ([]entity.Category, int, error)
	UpdateCategory(ctx context.Context, c *entity.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}
