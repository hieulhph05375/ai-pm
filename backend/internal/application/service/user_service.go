package service

import (
	"context"
	"errors"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ListUsers(ctx context.Context, searchQuery string) ([]entity.User, error)
	ListUsersPaginated(ctx context.Context, searchQuery string, page, limit int) ([]entity.User, int, error)
	GetUser(ctx context.Context, id uint) (*entity.User, error)
	CreateUser(ctx context.Context, u *entity.User, password string) error
	UpdateUser(ctx context.Context, adminID uint, u *entity.User) error
	ResetPassword(ctx context.Context, adminID uint, id uint, newPassword string) error
	ToggleStatus(ctx context.Context, adminID uint, id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{userRepo: ur}
}

func (s *userService) ListUsers(ctx context.Context, searchQuery string) ([]entity.User, error) {
	if searchQuery != "" {
		return s.userRepo.Search(ctx, searchQuery)
	}
	return s.userRepo.List(ctx)
}

func (s *userService) ListUsersPaginated(ctx context.Context, searchQuery string, page, limit int) ([]entity.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.userRepo.ListWithPagination(ctx, searchQuery, offset, limit)
}

func (s *userService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, u *entity.User, password string) error {
	// Check email uniqueness
	existing, err := s.userRepo.GetByEmail(ctx, u.Email)
	if err == nil && existing != nil {
		return errors.New("email này đã được sử dụng")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("không thể mã hóa mật khẩu")
	}

	u.HashedPassword = string(hashed)
	u.IsActive = true
	return s.userRepo.Create(ctx, u)
}

func (s *userService) UpdateUser(ctx context.Context, adminID uint, u *entity.User) error {
	existingUser, err := s.userRepo.GetByID(ctx, u.ID)
	if err != nil {
		return err
	}

	requestingUser, err := s.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}

	// 1. Only existing admins can change the IsAdmin flag of any user
	if existingUser.IsAdmin != u.IsAdmin {
		if !requestingUser.IsAdmin {
			return errors.New("chỉ quản trị viên hệ thống mới có quyền thay đổi quyền Admin")
		}
	}

	// 2. Prevent an admin from removing their own admin status (safety)
	if adminID == u.ID && existingUser.IsAdmin && !u.IsAdmin {
		return errors.New("không thể tự gỡ quyền Admin của chính mình")
	}

	// 3. Prevent non-admin from setting themselves as admin
	if adminID == u.ID && !existingUser.IsAdmin && u.IsAdmin {
		return errors.New("không thể tự cấp quyền Admin cho chính mình")
	}

	// Prevent admin from removing their own admin role (from RoleID=1, if still using roles for something)
	if adminID == u.ID && existingUser.RoleID == 1 && u.RoleID != 1 {
		return errors.New("không thể tự xóa vai trò Admin của chính mình")
	}

	// Check email uniqueness when changed
	if existingUser.Email != u.Email {
		emailOwner, err := s.userRepo.GetByEmail(ctx, u.Email)
		if err == nil && emailOwner != nil && emailOwner.ID != u.ID {
			return errors.New("email này đã được sử dụng bởi tài khoản khác")
		}
	}

	existingUser.FullName = u.FullName
	existingUser.RoleID = u.RoleID
	existingUser.Email = u.Email
	existingUser.IsAdmin = u.IsAdmin // Update the new field

	return s.userRepo.Update(ctx, existingUser)
}

func (s *userService) ResetPassword(ctx context.Context, adminID uint, id uint, newPassword string) error {
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("không thể mã hóa mật khẩu")
	}

	existingUser.HashedPassword = string(hashed)
	return s.userRepo.Update(ctx, existingUser)
}

func (s *userService) ToggleStatus(ctx context.Context, adminID uint, id uint) error {
	if adminID == id {
		return errors.New("không thể tự khoá tài khoản của chính mình")
	}

	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existingUser.IsActive = !existingUser.IsActive
	return s.userRepo.Update(ctx, existingUser)
}
