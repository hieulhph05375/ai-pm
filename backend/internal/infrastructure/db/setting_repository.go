package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"project-mgmt/backend/internal/domain/entity"
)

type settingRepo struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) *settingRepo {
	return &settingRepo{db: db}
}

func (r *settingRepo) scanSetting(rows interface {
	Scan(dest ...interface{}) error
}, s *entity.SystemSetting) error {
	var rawData []byte
	if err := rows.Scan(&s.Key, &rawData, &s.UpdatedAt); err != nil {
		return err
	}
	return json.Unmarshal(rawData, &s.Value)
}

func (r *settingRepo) Get(ctx context.Context, key string) (*entity.SystemSetting, error) {
	s := &entity.SystemSetting{}
	query := `SELECT key, value, updated_at FROM system_settings WHERE key = $1`
	if err := r.scanSetting(r.db.QueryRowContext(ctx, query, key), s); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *settingRepo) Set(ctx context.Context, s *entity.SystemSetting) error {
	rawData, err := json.Marshal(s.Value)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO system_settings (key, value, updated_at) 
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP`
	_, err = r.db.ExecContext(ctx, query, s.Key, rawData)
	return err
}

func (r *settingRepo) GetAll(ctx context.Context) ([]*entity.SystemSetting, error) {
	query := `SELECT key, value, updated_at FROM system_settings`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.SystemSetting
	for rows.Next() {
		s := &entity.SystemSetting{}
		if err := r.scanSetting(rows, s); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}
