package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"project-mgmt/backend/internal/infrastructure/config"
	"time"
)

func NewConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil
}
