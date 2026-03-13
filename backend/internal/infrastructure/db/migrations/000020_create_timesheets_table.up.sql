CREATE TABLE IF NOT EXISTS timesheets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    node_id INT REFERENCES wbs_nodes(id) ON DELETE CASCADE,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    work_date DATE NOT NULL,
    hours NUMERIC(5, 2) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'DRAFT',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_timesheets_user_id ON timesheets(user_id);
CREATE INDEX idx_timesheets_project_id ON timesheets(project_id);
CREATE INDEX idx_timesheets_node_id ON timesheets(node_id);
CREATE INDEX idx_timesheets_task_id ON timesheets(task_id);
CREATE INDEX idx_timesheets_work_date ON timesheets(work_date);
