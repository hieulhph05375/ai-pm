-- Enable ltree extension for hierarchical queries
CREATE EXTENSION IF NOT EXISTS ltree;

-- Projects Table
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    project_id VARCHAR(50) UNIQUE NOT NULL, -- Business ID (e.g., PRJ-2026-001)
    project_name VARCHAR(255) NOT NULL,
    description TEXT,
    project_manager VARCHAR(255),
    sponsor VARCHAR(255),
    requesting_department VARCHAR(255),
    current_phase VARCHAR(50) DEFAULT 'Initiation',
    project_status VARCHAR(50) DEFAULT 'Running',
    strategic_goal TEXT,
    portfolio_category VARCHAR(100),
    strategic_score INT DEFAULT 0,
    priority_level VARCHAR(50) DEFAULT 'Medium',
    approved_budget NUMERIC(15, 2) DEFAULT 0.0,
    actual_cost NUMERIC(15, 2) DEFAULT 0.0,
    eac NUMERIC(15, 2) DEFAULT 0.0,
    capex_opex_ratio VARCHAR(50),
    expected_roi NUMERIC(15, 2) DEFAULT 0.0,
    payback_period INT DEFAULT 0, -- months
    benefit_realization_date DATE,
    planned_start_date DATE,
    actual_start_date DATE,
    planned_end_date DATE,
    actual_end_date DATE,
    progress INT DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    overall_health VARCHAR(50) DEFAULT 'Green',
    spi NUMERIC(5, 2) DEFAULT 1.0,
    cpi NUMERIC(5, 2) DEFAULT 1.0,
    last_executive_summary TEXT,
    estimated_effort INT DEFAULT 0, -- hours
    actual_effort INT DEFAULT 0, -- hours
    resource_risk_flag BOOLEAN DEFAULT FALSE,
    missing_skills TEXT,
    systemic_risk_level VARCHAR(50) DEFAULT 'Low',
    open_critical_risks INT DEFAULT 0,
    compliance_impact VARCHAR(100) DEFAULT 'No Impact',
    dependencies_summary TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- WBS Nodes Table (Hierarchical)
CREATE TABLE IF NOT EXISTS wbs_nodes (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    parent_id INT REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    wbs_code VARCHAR(50), -- e.g., 1.2.1
    depth INT DEFAULT 1 CHECK (depth >= 1 AND depth <= 5),
    path LTREE,
    node_type VARCHAR(50) NOT NULL, -- phase, milestone, task
    task_id INT, -- Link to tasks table from Phase 2
    order_index INT DEFAULT 0,
    start_date DATE,
    end_date DATE,
    progress INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast path-based searching
CREATE INDEX idx_wbs_path ON wbs_nodes USING GIST (path);
CREATE INDEX idx_wbs_project ON wbs_nodes(project_id);

-- Project Baselines Table
CREATE TABLE IF NOT EXISTS project_baselines (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    version_name VARCHAR(100) NOT NULL,
    description TEXT,
    snapshot_data JSONB NOT NULL, -- Full project+wbs state
    created_by INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Project Dependencies Table (Inter-task/node dependencies)
CREATE TABLE IF NOT EXISTS wbs_dependencies (
    id SERIAL PRIMARY KEY,
    predecessor_id INT REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    successor_id INT REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    dependency_type VARCHAR(10) DEFAULT 'FS', -- FS, SS, FF, SF
    lag_days INT DEFAULT 0
);
