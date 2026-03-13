-- 000019_add_pmi_fields.down.sql

ALTER TABLE wbs_nodes DROP COLUMN IF EXISTS planned_value;
ALTER TABLE wbs_nodes DROP COLUMN IF EXISTS actual_cost;

ALTER TABLE wbs_baseline_nodes DROP COLUMN IF EXISTS planned_value;
ALTER TABLE wbs_baseline_nodes DROP COLUMN IF EXISTS actual_cost;

DROP INDEX IF EXISTS idx_projects_last_reminder;
ALTER TABLE projects DROP COLUMN IF EXISTS last_reminder_at;
