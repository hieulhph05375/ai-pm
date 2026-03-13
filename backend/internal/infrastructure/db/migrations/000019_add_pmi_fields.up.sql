-- 000019_add_pmi_fields.up.sql

-- Add budget and actual cost fields to WBS nodes and baseline nodes for EVM
ALTER TABLE wbs_nodes ADD COLUMN IF NOT EXISTS planned_value DECIMAL(15, 2) DEFAULT 0;
ALTER TABLE wbs_nodes ADD COLUMN IF NOT EXISTS actual_cost DECIMAL(15, 2) DEFAULT 0;

ALTER TABLE wbs_baseline_nodes ADD COLUMN IF NOT EXISTS planned_value DECIMAL(15, 2) DEFAULT 0;
ALTER TABLE wbs_baseline_nodes ADD COLUMN IF NOT EXISTS actual_cost DECIMAL(15, 2) DEFAULT 0;

-- Add reminder tracking to projects
ALTER TABLE projects ADD COLUMN IF NOT EXISTS last_reminder_at TIMESTAMP WITH TIME ZONE;

-- Add index for reminder checks
CREATE INDEX IF NOT EXISTS idx_projects_last_reminder ON projects(last_reminder_at);
