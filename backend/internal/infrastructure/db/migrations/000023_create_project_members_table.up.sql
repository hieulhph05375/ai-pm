-- 000023_create_project_members_table.up.sql

CREATE TABLE IF NOT EXISTS project_members (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    project_role VARCHAR(50) DEFAULT 'Member', -- e.g., 'Lead', 'Member', 'Observer'
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, user_id)
);

-- Index for fast lookups by user (most common query)
CREATE INDEX IF NOT EXISTS idx_project_members_user ON project_members(user_id);
