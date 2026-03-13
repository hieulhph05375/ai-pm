-- Rollback: Drop all performance indexes added in 000028

DROP INDEX IF EXISTS idx_wbs_node_type;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP INDEX IF EXISTS idx_tasks_assigned_to;
DROP INDEX IF EXISTS idx_project_members_user_id;
DROP INDEX IF EXISTS idx_project_members_project_id;
DROP INDEX IF EXISTS idx_timesheets_user_id;
DROP INDEX IF EXISTS idx_timesheets_task_id;
DROP INDEX IF EXISTS idx_risks_project_id;
DROP INDEX IF EXISTS idx_issues_project_id;
DROP INDEX IF EXISTS idx_wbs_project_path;
