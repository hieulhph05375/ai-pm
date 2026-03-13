-- Migration: 000010_create_system_settings.up.sql
-- Table: system_settings
-- Stores organization-wide global settings and configurations
CREATE TABLE IF NOT EXISTS system_settings (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed defaults: Saturday(6) and Sunday(0) are rest days
-- The value is stored as an array of day numbers (0-6, where 0=Sunday)
INSERT INTO system_settings (key, value)
VALUES ('rest_days', '[6, 0]'::jsonb)
ON CONFLICT (key) DO NOTHING;
