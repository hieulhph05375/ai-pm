-- Create enum types for WBS nodes and dependencies
DO $$ BEGIN
    CREATE TYPE wbs_node_type AS ENUM ('Phase', 'Milestone', 'Task', 'Sub-task');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE dependency_type AS ENUM ('FS', 'SF', 'SS', 'FF');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Enable ltree extension for hierarchical queries (re-run is safe)
CREATE EXTENSION IF NOT EXISTS ltree;

-- We need to drop existing tables from migration 000006 if they conflict, 
-- but since 000006 is already committed, let's ALTER the existing tables 
-- or drop and recreate if we want a fresh start for wbs_nodes and wbs_dependencies.
-- Given ADR-009 states wbs_nodes is completely independent of tasks and uses ltree differently,
-- let's drop the ones from 000006 and recreate them exactly as requested in 4.1 Plan.

DROP TABLE IF EXISTS wbs_dependencies CASCADE;
DROP TABLE IF EXISTS project_baselines CASCADE;
DROP TABLE IF EXISTS wbs_nodes CASCADE;

-- WBS Nodes Table (Hierarchical with Ltree, Independent of Tasks)
CREATE TABLE wbs_nodes (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    type wbs_node_type NOT NULL,
    path LTREE NOT NULL CHECK (path::text != ''),
    order_index INT DEFAULT 0,
    planned_start_date DATE,
    planned_end_date DATE,
    actual_start_date DATE,
    actual_end_date DATE,
    progress DECIMAL(5,2) DEFAULT 0.00 CHECK (progress >= 0 AND progress <= 100),
    assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for fast searching and Ltree operations
CREATE INDEX idx_wbs_nodes_project_id ON wbs_nodes(project_id);
CREATE INDEX idx_wbs_nodes_path_gist ON wbs_nodes USING GIST(path);

-- WBS Dependencies Table
CREATE TABLE wbs_dependencies (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    predecessor_id INT NOT NULL REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    successor_id INT NOT NULL REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    type dependency_type NOT NULL DEFAULT 'FS',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_wbs_dependency UNIQUE (predecessor_id, successor_id)
);

CREATE INDEX idx_wbs_deps_project_id ON wbs_dependencies(project_id);
