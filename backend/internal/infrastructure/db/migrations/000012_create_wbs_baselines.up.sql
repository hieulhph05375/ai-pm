CREATE TABLE IF NOT EXISTS wbs_baselines (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wbs_baseline_nodes (
    baseline_id INT NOT NULL REFERENCES wbs_baselines(id) ON DELETE CASCADE,
    node_id INT NOT NULL, /* deliberately not fk referenced to allow node deletion but keep baseline */
    path ltree NOT NULL,
    planned_start_date TIMESTAMP WITH TIME ZONE,
    planned_end_date TIMESTAMP WITH TIME ZONE,
    progress NUMERIC(5,2) DEFAULT 0,
    PRIMARY KEY (baseline_id, node_id)
);

CREATE INDEX IF NOT EXISTS idx_wbs_baselines_project_id ON wbs_baselines(project_id);
CREATE INDEX IF NOT EXISTS idx_wbs_baseline_nodes_baseline_id ON wbs_baseline_nodes(baseline_id);
