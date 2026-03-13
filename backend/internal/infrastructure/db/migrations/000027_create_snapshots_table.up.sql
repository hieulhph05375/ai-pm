-- 000027_create_snapshots_table.up.sql

-- Project level snapshots for historical trend tracking
CREATE TABLE IF NOT EXISTS project_snapshots (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    spi NUMERIC(5, 2) DEFAULT 1.0,
    cpi NUMERIC(5, 2) DEFAULT 1.0,
    ev NUMERIC(15, 2) DEFAULT 0.0,
    ac NUMERIC(15, 2) DEFAULT 0.0,
    pv NUMERIC(15, 2) DEFAULT 0.0,
    progress INT DEFAULT 0,
    captured_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Milestone level snapshots for historical trend tracking
CREATE TABLE IF NOT EXISTS milestone_snapshots (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    node_id INT REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    milestone_name VARCHAR(255) NOT NULL,
    planned_date DATE,
    actual_date DATE,
    captured_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast trend queries
CREATE INDEX IF NOT EXISTS idx_project_snapshots_project_at ON project_snapshots(project_id, captured_at);
CREATE INDEX IF NOT EXISTS idx_milestone_snapshots_node_at ON milestone_snapshots(node_id, captured_at);
