CREATE TABLE IF NOT EXISTS issues (
    id            SERIAL PRIMARY KEY,
    project_id    INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    issue_type    VARCHAR(50) NOT NULL DEFAULT 'Issue' CHECK (issue_type IN ('Bug', 'Issue')),
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    status        VARCHAR(50) NOT NULL DEFAULT 'Open' CHECK (status IN ('Open', 'In Progress', 'Resolved', 'Closed')),
    priority      VARCHAR(50) NOT NULL DEFAULT 'Medium' CHECK (priority IN ('Low', 'Medium', 'High', 'Critical')),
    assignee_id   INTEGER REFERENCES users(id) ON DELETE SET NULL,
    reporter_id   INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_issues_project_id ON issues(project_id);
CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
