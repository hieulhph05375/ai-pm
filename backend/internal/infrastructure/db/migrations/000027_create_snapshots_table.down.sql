-- 000027_create_snapshots_table.down.sql

DROP INDEX IF EXISTS idx_milestone_snapshots_node_at;
DROP INDEX IF EXISTS idx_project_snapshots_project_at;
DROP TABLE IF EXISTS milestone_snapshots;
DROP TABLE IF EXISTS project_snapshots;
