package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"strconv"
	"time"
)

type holidayRepo struct {
	db *sql.DB
}

func NewHolidayRepository(db *sql.DB) repository.HolidayRepository {
	return &holidayRepo{db: db}
}

func (r *holidayRepo) scanHoliday(rows interface {
	Scan(dest ...interface{}) error
}, h *entity.Holiday) error {
	var typeID sql.NullInt64
	var typeName, typeColor, typeIcon sql.NullString

	err := rows.Scan(
		&h.ID, &h.Name, &h.Date, &h.Type, &typeID, &h.IsRecurring, &h.CreatedAt, &h.UpdatedAt,
		&typeName, &typeColor, &typeIcon,
	)
	if err != nil {
		return err
	}

	if typeID.Valid {
		h.TypeID = uint(typeID.Int64)
		h.TypeCat = &entity.Category{
			ID:    h.TypeID,
			Name:  typeName.String,
			Color: nullStringPtr(typeColor),
			Icon:  nullStringPtr(typeIcon),
		}
	}
	return nil
}

func (r *holidayRepo) Create(ctx context.Context, h *entity.Holiday) error {
	if h.TypeID == 0 && h.Type != "" {
		// Lookup typeID from categories for HOLIDAY_TYPE
		err := r.db.QueryRowContext(ctx, `
			SELECT id FROM categories 
			WHERE type_id = (SELECT id FROM category_types WHERE code = 'HOLIDAY_TYPE') 
			AND name = $1`, h.Type).Scan(&h.TypeID)
		if err != nil && err != sql.ErrNoRows {
			// Not a fatal error if not found, but FK might still fail later
			// We can decide to error here or let the INSERT fail
		}
	}
	query := `INSERT INTO holidays (name, date, type, type_id, is_recurring) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, h.Name, h.Date, h.Type, h.TypeID, h.IsRecurring).Scan(&h.ID, &h.CreatedAt, &h.UpdatedAt)
}

func (r *holidayRepo) GetByID(ctx context.Context, id int) (*entity.Holiday, error) {
	h := &entity.Holiday{}
	query := `
		SELECT h.id, h.name, h.date, h.type, h.type_id, h.is_recurring, h.created_at, h.updated_at,
		       c.name as type_name, c.color as type_color, c.icon as type_icon
		FROM holidays h
		LEFT JOIN categories c ON h.type_id = c.id
		WHERE h.id = $1`
	if err := r.scanHoliday(r.db.QueryRowContext(ctx, query, id), h); err != nil {
		return nil, err
	}
	return h, nil
}

func (r *holidayRepo) GetByDate(ctx context.Context, date time.Time) (*entity.Holiday, error) {
	h := &entity.Holiday{}
	query := `
		SELECT h.id, h.name, h.date, h.type, h.type_id, h.is_recurring, h.created_at, h.updated_at,
		       c.name as type_name, c.color as type_color, c.icon as type_icon
		FROM holidays h
		LEFT JOIN categories c ON h.type_id = c.id
		WHERE h.date = $1`
	if err := r.scanHoliday(r.db.QueryRowContext(ctx, query, date.Format("2006-01-02")), h); err != nil {
		return nil, err
	}
	return h, nil
}

func (r *holidayRepo) Update(ctx context.Context, h *entity.Holiday) error {
	query := `UPDATE holidays SET name=$1, date=$2, type=$3, type_id=$4, is_recurring=$5, updated_at=CURRENT_TIMESTAMP WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query, h.Name, h.Date, h.Type, h.TypeID, h.IsRecurring, h.ID)
	return err
}

func (r *holidayRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM holidays WHERE id=$1", id)
	return err
}

func (r *holidayRepo) List(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error) {
	query := `
		SELECT h.id, h.name, h.date, h.type, h.type_id, h.is_recurring, h.created_at, h.updated_at,
		       c.name as type_name, c.color as type_color, c.icon as type_icon
		FROM holidays h
		LEFT JOIN categories c ON h.type_id = c.id`
	var args []interface{}
	if !start.IsZero() && !end.IsZero() {
		query += " WHERE date BETWEEN $1 AND $2"
		args = append(args, start.Format("2006-01-02"), end.Format("2006-01-02"))
	}
	query += " ORDER BY date ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*entity.Holiday{}
	for rows.Next() {
		h := &entity.Holiday{}
		if err := r.scanHoliday(rows, h); err != nil {
			return nil, err
		}
		results = append(results, h)
	}
	return results, nil
}

func (r *holidayRepo) ListWithPagination(ctx context.Context, start, end time.Time, offset, limit int) ([]*entity.Holiday, int, error) {
	query := `
		SELECT h.id, h.name, h.date, h.type, h.type_id, h.is_recurring, h.created_at, h.updated_at,
		       c.name as type_name, c.color as type_color, c.icon as type_icon
		FROM holidays h
		LEFT JOIN categories c ON h.type_id = c.id`
	countQuery := `SELECT COUNT(*) FROM holidays h`
	var args []interface{}

	if !start.IsZero() && !end.IsZero() {
		filter := " WHERE date BETWEEN $1 AND $2"
		query += filter
		countQuery += filter
		args = append(args, start.Format("2006-01-02"), end.Format("2006-01-02"))
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY date ASC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := []*entity.Holiday{}
	for rows.Next() {
		h := &entity.Holiday{}
		if err := r.scanHoliday(rows, h); err != nil {
			return nil, 0, err
		}
		results = append(results, h)
	}
	return results, total, nil
}

func (r *holidayRepo) Count(ctx context.Context, start, end time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM holidays`
	var args []interface{}
	if !start.IsZero() && !end.IsZero() {
		query += " WHERE date BETWEEN $1 AND $2"
		args = append(args, start.Format("2006-01-02"), end.Format("2006-01-02"))
	}

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}
