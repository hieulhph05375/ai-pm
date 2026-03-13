package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepo{db: db}
}

// CategoryType Implementation

func (r *categoryRepo) CreateType(ctx context.Context, ct *entity.CategoryType) error {
	now := time.Now()
	ct.CreatedAt = now
	ct.UpdatedAt = now
	isActive := true
	if ct.IsActive != nil {
		isActive = *ct.IsActive
	} else {
		ct.IsActive = &isActive
	}
	query := `INSERT INTO category_types (name, code, description, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRowContext(ctx, query, ct.Name, ct.Code, ct.Description, isActive, ct.CreatedAt, ct.UpdatedAt).Scan(&ct.ID)
}

func (r *categoryRepo) GetTypeByID(ctx context.Context, id uint) (*entity.CategoryType, error) {
	ct := &entity.CategoryType{}
	query := `SELECT id, name, code, description, is_active, created_at, updated_at FROM category_types WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&ct.ID, &ct.Name, &ct.Code, &ct.Description, &ct.IsActive, &ct.CreatedAt, &ct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (r *categoryRepo) ListTypes(ctx context.Context, search string, limit, offset int) ([]entity.CategoryType, int, error) {
	baseQuery := `SELECT id, name, code, description, is_active, created_at, updated_at FROM category_types`
	countQuery := `SELECT COUNT(*) FROM category_types`
	whereClause := ""
	args := []interface{}{}

	if search != "" {
		whereClause = ` WHERE name ILIKE $1 OR code ILIKE $1 OR description ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := baseQuery + whereClause + ` ORDER BY name ASC`
	if limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
		args = append(args, limit, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	types := []entity.CategoryType{}
	for rows.Next() {
		var ct entity.CategoryType
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Code, &ct.Description, &ct.IsActive, &ct.CreatedAt, &ct.UpdatedAt); err != nil {
			return nil, 0, err
		}
		types = append(types, ct)
	}
	return types, total, nil
}

func (r *categoryRepo) UpdateType(ctx context.Context, ct *entity.CategoryType) error {
	ct.UpdatedAt = time.Now()
	isActive := true
	if ct.IsActive != nil {
		isActive = *ct.IsActive
	}
	query := `UPDATE category_types SET name = $1, code = $2, description = $3, is_active = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, ct.Name, ct.Code, ct.Description, isActive, ct.UpdatedAt, ct.ID)
	return err
}

func (r *categoryRepo) DeleteType(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM category_types WHERE id = $1`, id)
	return err
}

// Category Implementation

func (r *categoryRepo) CreateCategory(ctx context.Context, c *entity.Category) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	isActive := true
	if c.IsActive != nil {
		isActive = *c.IsActive
	} else {
		c.IsActive = &isActive
	}
	query := `INSERT INTO categories (type_id, parent_id, name, color, icon, description, is_active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	return r.db.QueryRowContext(ctx, query, c.TypeID, c.ParentID, c.Name, c.Color, c.Icon, c.Description, isActive, c.CreatedAt, c.UpdatedAt).Scan(&c.ID)
}

func (r *categoryRepo) GetCategoryByID(ctx context.Context, id uint) (*entity.Category, error) {
	c := &entity.Category{}
	query := `SELECT id, type_id, parent_id, name, color, icon, description, is_active, created_at, updated_at FROM categories WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.TypeID, &c.ParentID, &c.Name, &c.Color, &c.Icon, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *categoryRepo) ListCategories(ctx context.Context, typeID *uint, search string, limit, offset int) ([]entity.Category, int, error) {
	baseQuery := `SELECT c.id, c.type_id, c.parent_id, c.name, c.color, c.icon, c.description, c.is_active, c.created_at, c.updated_at,
	                 t.name as type_name, p.name as parent_name
			  FROM categories c
			  JOIN category_types t ON c.type_id = t.id
			  LEFT JOIN categories p ON c.parent_id = p.id`
	countQuery := `SELECT COUNT(*) FROM categories c JOIN category_types t ON c.type_id = t.id`

	whereClauses := []string{}
	args := []interface{}{}

	if typeID != nil {
		args = append(args, *typeID)
		whereClauses = append(whereClauses, fmt.Sprintf("c.type_id = $%d", len(args)))
	}

	if search != "" {
		args = append(args, "%"+search+"%")
		whereClauses = append(whereClauses, fmt.Sprintf("(c.name ILIKE $%d OR c.description ILIKE $%d)", len(args), len(args)))
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := baseQuery + whereClause + ` ORDER BY t.name ASC, COALESCE(p.name, ''), c.name ASC`
	if limit > 0 {
		args = append(args, limit, offset)
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)-1, len(args))
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	categories := []entity.Category{}
	for rows.Next() {
		var c entity.Category
		var typeName string
		var parentName sql.NullString
		var color, icon, description sql.NullString
		var isActive sql.NullBool
		var parentID sql.NullInt64

		err := rows.Scan(
			&c.ID, &c.TypeID, &parentID, &c.Name, &color, &icon, &description, &isActive, &c.CreatedAt, &c.UpdatedAt,
			&typeName, &parentName,
		)
		if err != nil {
			return nil, 0, err
		}

		if parentID.Valid {
			uID := uint(parentID.Int64)
			c.ParentID = &uID
		}
		if color.Valid {
			c.Color = &color.String
		}
		if icon.Valid {
			c.Icon = &icon.String
		}
		if description.Valid {
			c.Description = &description.String
		}
		if isActive.Valid {
			c.IsActive = &isActive.Bool
		}

		c.Type = &entity.CategoryType{ID: c.TypeID, Name: typeName}
		if c.ParentID != nil && parentName.Valid {
			c.Parent = &entity.Category{ID: *c.ParentID, Name: parentName.String}
		}
		categories = append(categories, c)
	}
	return categories, total, nil
}

func (r *categoryRepo) UpdateCategory(ctx context.Context, c *entity.Category) error {
	c.UpdatedAt = time.Now()
	isActive := true
	if c.IsActive != nil {
		isActive = *c.IsActive
	}
	query := `UPDATE categories SET type_id = $1, parent_id = $2, name = $3, color = $4, icon = $5, description = $6, is_active = $7, updated_at = $8 WHERE id = $9`
	_, err := r.db.ExecContext(ctx, query, c.TypeID, c.ParentID, c.Name, c.Color, c.Icon, c.Description, isActive, c.UpdatedAt, c.ID)
	return err
}

func (r *categoryRepo) DeleteCategory(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
