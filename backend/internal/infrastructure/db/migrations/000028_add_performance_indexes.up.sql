-- Performance Index Migration
-- Adds additional missing indexes on high-traffic foreign key columns
-- NOTE: idx_wbs_nodes_project_id and idx_wbs_nodes_path_gist already exist from migration 000008

-- BTree index on wbs_nodes.assigned_to (filter by assignee)
CREATE INDEX IF NOT EXISTS idx_wbs_nodes_assigned_to ON wbs_nodes (assigned_to);

-- BTree index on wbs_nodes.type (filter by Phase/Milestone/Task)
CREATE INDEX IF NOT EXISTS idx_wbs_nodes_type ON wbs_nodes (type);

-- Compound index on wbs_nodes (project_id, type) — filter tree by type per project
CREATE INDEX IF NOT EXISTS idx_wbs_project_type ON wbs_nodes (project_id, type);

-- BTree index on project_members.user_id (used in RBAC lookups)
CREATE INDEX IF NOT EXISTS idx_project_members_user_id ON project_members (user_id);

-- BTree index on project_members.project_id (used in member resolution)
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members (project_id);

-- BTree index on timesheets.user_id (resource workload queries)
CREATE INDEX IF NOT EXISTS idx_timesheets_user_id ON timesheets (user_id);

-- BTree index on timesheets.task_id (workload by task)
CREATE INDEX IF NOT EXISTS idx_timesheets_task_id ON timesheets (task_id);

-- BTree index on risks.project_id (risk list per project)
CREATE INDEX IF NOT EXISTS idx_risks_project_id ON risks (project_id);

-- BTree index on issues.project_id (issue list per project)
CREATE INDEX IF NOT EXISTS idx_issues_project_id ON issues (project_id);
