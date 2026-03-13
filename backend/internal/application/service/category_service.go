package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type CategoryService interface {
	// CategoryType
	CreateType(ctx context.Context, ct *entity.CategoryType) error
	GetType(ctx context.Context, id uint) (*entity.CategoryType, error)
	ListTypes(ctx context.Context, search string, limit, offset int) ([]entity.CategoryType, int, error)
	UpdateType(ctx context.Context, ct *entity.CategoryType) error
	DeleteType(ctx context.Context, id uint) error

	// Category
	CreateCategory(ctx context.Context, c *entity.Category) error
	GetCategory(ctx context.Context, id uint) (*entity.Category, error)
	ListCategories(ctx context.Context, typeID *uint, search string, limit, offset int) ([]entity.Category, int, error)
	UpdateCategory(ctx context.Context, c *entity.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateType(ctx context.Context, ct *entity.CategoryType) error {
	return s.repo.CreateType(ctx, ct)
}

func (s *categoryService) GetType(ctx context.Context, id uint) (*entity.CategoryType, error) {
	return s.repo.GetTypeByID(ctx, id)
}

func (s *categoryService) ListTypes(ctx context.Context, search string, limit, offset int) ([]entity.CategoryType, int, error) {
	return s.repo.ListTypes(ctx, search, limit, offset)
}

func (s *categoryService) UpdateType(ctx context.Context, ct *entity.CategoryType) error {
	return s.repo.UpdateType(ctx, ct)
}

func (s *categoryService) DeleteType(ctx context.Context, id uint) error {
	return s.repo.DeleteType(ctx, id)
}

func (s *categoryService) CreateCategory(ctx context.Context, c *entity.Category) error {
	return s.repo.CreateCategory(ctx, c)
}

func (s *categoryService) GetCategory(ctx context.Context, id uint) (*entity.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *categoryService) ListCategories(ctx context.Context, typeID *uint, search string, limit, offset int) ([]entity.Category, int, error) {
	return s.repo.ListCategories(ctx, typeID, search, limit, offset)
}

func (s *categoryService) UpdateCategory(ctx context.Context, c *entity.Category) error {
	return s.repo.UpdateCategory(ctx, c)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.DeleteCategory(ctx, id)
}
