-- 000030_global_category_refactor.down.sql

ALTER TABLE holidays DROP COLUMN IF EXISTS type_id;
ALTER TABLE stakeholders DROP COLUMN IF EXISTS role_id;
ALTER TABLE risks DROP COLUMN IF EXISTS status_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS status_id, DROP COLUMN IF EXISTS priority_id;
ALTER TABLE projects DROP COLUMN IF EXISTS project_status_id, DROP COLUMN IF EXISTS current_phase_id, DROP COLUMN IF EXISTS portfolio_category_id, DROP COLUMN IF EXISTS overall_health_id;

-- Note: We don't delete the category types or categories as they might be used by other things now.
