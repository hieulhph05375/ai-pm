package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"strconv"
	"time"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *entity.User) error {
	query := `INSERT INTO users (email, hashed_password, full_name, role_id, is_admin, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	now := u.CreatedAt
	if now.IsZero() {
		now = u.UpdatedAt
	}

	return r.db.QueryRowContext(ctx, query, u.Email, u.HashedPassword, u.FullName, u.RoleID, u.IsAdmin, u.IsActive, now, now).Scan(&u.ID)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT id, email, hashed_password, full_name, role_id, is_admin, is_active, failed_login_attempts, locked_until, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.HashedPassword, &user.FullName, &user.RoleID, &user.IsAdmin, &user.IsActive, &user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT id, email, hashed_password, full_name, role_id, is_admin, is_active, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.HashedPassword, &user.FullName, &user.RoleID, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Update(ctx context.Context, u *entity.User) error {
	query := `UPDATE users SET 
				email = $1, 
				hashed_password = $2, 
				full_name = $3, 
				role_id = $4, 
				is_admin = $5,
				is_active = $6,
				failed_login_attempts = $7, 
				locked_until = $8, 
				updated_at = $9 
			  WHERE id = $10`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		u.Email, u.HashedPassword, u.FullName, u.RoleID, u.IsAdmin, u.IsActive,
		u.FailedLoginAttempts, u.LockedUntil, now, u.ID)
	return err
}

func (r *userRepo) List(ctx context.Context) ([]entity.User, error) {
	query := `SELECT id, email, full_name, role_id, is_admin, is_active, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entity.User{}
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.RoleID, &u.IsAdmin, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) Search(ctx context.Context, searchQuery string) ([]entity.User, error) {
	query := `SELECT id, email, full_name, role_id, is_admin, is_active, created_at, updated_at 
			  FROM users 
			  WHERE email ILIKE '%' || $1 || '%' OR full_name ILIKE '%' || $1 || '%'`

	rows, err := r.db.QueryContext(ctx, query, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entity.User{}
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.RoleID, &u.IsAdmin, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
func (r *userRepo) ListWithPagination(ctx context.Context, searchQuery string, offset, limit int) ([]entity.User, int, error) {
	query := `SELECT id, email, full_name, role_id, is_admin, is_active, created_at, updated_at FROM users`
	countQuery := `SELECT COUNT(*) FROM users`
	var args []interface{}

	if searchQuery != "" {
		filter := ` WHERE email ILIKE $1 OR full_name ILIKE $1`
		query += filter
		countQuery += filter
		args = append(args, "%"+searchQuery+"%")
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY full_name ASC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []entity.User{}
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.RoleID, &u.IsAdmin, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *userRepo) Count(ctx context.Context, searchQuery string) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var args []interface{}

	if searchQuery != "" {
		query += ` WHERE email ILIKE $1 OR full_name ILIKE $1`
		args = append(args, "%"+searchQuery+"%")
	}

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}
