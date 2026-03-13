-- 000008_create_category_management.up.sql

CREATE TABLE IF NOT EXISTS category_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    type_id INTEGER NOT NULL REFERENCES category_types(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(20),
    icon VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT categories_type_name_unique UNIQUE (type_id, name)
);

CREATE INDEX IF NOT EXISTS idx_categories_type_id ON categories(type_id);
