package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://devops:bNyb3cTNLxeM@dev-rds-postgres-aurora-common-apps-ekyc-db.cluster-ccjec86kuais.ap-southeast-1.rds.amazonaws.com:5432/location_service_test?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-2026"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
