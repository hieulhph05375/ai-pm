-- 000003_create_tasks_table.up.sql

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    wbs_code VARCHAR(50),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'Backlog',
    priority VARCHAR(50) DEFAULT 'Medium',
    assignee_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    project_id INTEGER, -- Simplified for now, can point to a projects table later
    due_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
