-- 000030_global_category_refactor.up.sql

-- 1. Create Category Types
INSERT INTO category_types (name, code, description)
SELECT 'Project Status', 'PROJECT_STATUS', 'Standard statuses for project lifecycle' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'PROJECT_STATUS' OR name = 'Project Status');
INSERT INTO category_types (name, code, description)
SELECT 'Project Phase', 'PROJECT_PHASE', 'Standard phases for project management' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'PROJECT_PHASE' OR name = 'Project Phase');
INSERT INTO category_types (name, code, description)
SELECT 'Portfolio Category', 'PORTFOLIO_CATEGORY', 'Strategic classification for projects' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'PORTFOLIO_CATEGORY' OR name = 'Portfolio Category');
INSERT INTO category_types (name, code, description)
SELECT 'Project Health', 'PROJECT_HEALTH', 'Visual health indicators (RAG)' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'PROJECT_HEALTH' OR name = 'Project Health');
INSERT INTO category_types (name, code, description)
SELECT 'Task Status', 'TASK_STATUS', 'Simple workflow for personal tasks' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'TASK_STATUS' OR name = 'Task Status');
INSERT INTO category_types (name, code, description)
SELECT 'Task Priority', 'TASK_PRIORITY', 'Priority levels for personal tasks' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'TASK_PRIORITY' OR name = 'Task Priority');
INSERT INTO category_types (name, code, description)
SELECT 'Risk Status', 'RISK_STATUS', 'Workflow for project risks' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'RISK_STATUS' OR name = 'Risk Status');
INSERT INTO category_types (name, code, description)
SELECT 'Stakeholder Role', 'STAKEHOLDER_ROLE', 'Organizational roles for stakeholders' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'STAKEHOLDER_ROLE' OR name = 'Stakeholder Role');
INSERT INTO category_types (name, code, description)
SELECT 'Holiday Type', 'HOLIDAY_TYPE', 'Classification for holidays' WHERE NOT EXISTS (SELECT 1 FROM category_types WHERE code = 'HOLIDAY_TYPE' OR name = 'Holiday Type');

-- 2. Seed Categories
-- Project Status
INSERT INTO categories (type_id, name, color) 
SELECT id, 'Planned', 'slate' FROM category_types WHERE code = 'PROJECT_STATUS' UNION ALL
SELECT id, 'Active', 'blue' FROM category_types WHERE code = 'PROJECT_STATUS' UNION ALL
SELECT id, 'On Hold', 'amber' FROM category_types WHERE code = 'PROJECT_STATUS' UNION ALL
SELECT id, 'Completed', 'emerald' FROM category_types WHERE code = 'PROJECT_STATUS' UNION ALL
SELECT id, 'Cancelled', 'rose' FROM category_types WHERE code = 'PROJECT_STATUS'
ON CONFLICT (type_id, name) DO NOTHING;

-- Project Phase
INSERT INTO categories (type_id, name)
SELECT id, 'Initiation' FROM category_types WHERE code = 'PROJECT_PHASE' UNION ALL
SELECT id, 'Planning' FROM category_types WHERE code = 'PROJECT_PHASE' UNION ALL
SELECT id, 'Execution' FROM category_types WHERE code = 'PROJECT_PHASE' UNION ALL
SELECT id, 'Monitoring' FROM category_types WHERE code = 'PROJECT_PHASE' UNION ALL
SELECT id, 'Closing' FROM category_types WHERE code = 'PROJECT_PHASE'
ON CONFLICT (type_id, name) DO NOTHING;

-- Portfolio Category
INSERT INTO categories (type_id, name)
SELECT id, 'Strategic' FROM category_types WHERE code = 'PORTFOLIO_CATEGORY' UNION ALL
SELECT id, 'Operational' FROM category_types WHERE code = 'PORTFOLIO_CATEGORY' UNION ALL
SELECT id, 'Compliance' FROM category_types WHERE code = 'PORTFOLIO_CATEGORY'
ON CONFLICT (type_id, name) DO NOTHING;

-- Project Health
INSERT INTO categories (type_id, name, color)
SELECT id, 'Green', 'emerald' FROM category_types WHERE code = 'PROJECT_HEALTH' UNION ALL
SELECT id, 'Yellow', 'amber' FROM category_types WHERE code = 'PROJECT_HEALTH' UNION ALL
SELECT id, 'Red', 'rose' FROM category_types WHERE code = 'PROJECT_HEALTH'
ON CONFLICT (type_id, name) DO NOTHING;

-- Task Status
INSERT INTO categories (type_id, name, color)
SELECT id, 'Todo', 'slate' FROM category_types WHERE code = 'TASK_STATUS' UNION ALL
SELECT id, 'In Progress', 'blue' FROM category_types WHERE code = 'TASK_STATUS' UNION ALL
SELECT id, 'Done', 'emerald' FROM category_types WHERE code = 'TASK_STATUS'
ON CONFLICT (type_id, name) DO NOTHING;

-- Task Priority
INSERT INTO categories (type_id, name, color)
SELECT id, 'Low', 'slate' FROM category_types WHERE code = 'TASK_PRIORITY' UNION ALL
SELECT id, 'Medium', 'blue' FROM category_types WHERE code = 'TASK_PRIORITY' UNION ALL
SELECT id, 'High', 'amber' FROM category_types WHERE code = 'TASK_PRIORITY' UNION ALL
SELECT id, 'Urgent', 'rose' FROM category_types WHERE code = 'TASK_PRIORITY'
ON CONFLICT (type_id, name) DO NOTHING;

-- Risk Status
INSERT INTO categories (type_id, name, color)
SELECT id, 'Open', 'rose' FROM category_types WHERE code = 'RISK_STATUS' UNION ALL
SELECT id, 'Mitigated', 'blue' FROM category_types WHERE code = 'RISK_STATUS' UNION ALL
SELECT id, 'Closed', 'emerald' FROM category_types WHERE code = 'RISK_STATUS'
ON CONFLICT (type_id, name) DO NOTHING;

-- Stakeholder Role
INSERT INTO categories (type_id, name)
SELECT id, 'Internal' FROM category_types WHERE code = 'STAKEHOLDER_ROLE' UNION ALL
SELECT id, 'External' FROM category_types WHERE code = 'STAKEHOLDER_ROLE' UNION ALL
SELECT id, 'Vendor' FROM category_types WHERE code = 'STAKEHOLDER_ROLE' UNION ALL
SELECT id, 'Government' FROM category_types WHERE code = 'STAKEHOLDER_ROLE'
ON CONFLICT (type_id, name) DO NOTHING;

-- Holiday Type
INSERT INTO categories (type_id, name)
SELECT id, 'State' FROM category_types WHERE code = 'HOLIDAY_TYPE' UNION ALL
SELECT id, 'Company' FROM category_types WHERE code = 'HOLIDAY_TYPE'
ON CONFLICT (type_id, name) DO NOTHING;

-- 3. Migration (Adding ID columns)
-- Project
ALTER TABLE projects 
ADD COLUMN IF NOT EXISTS project_status_id INTEGER REFERENCES categories(id),
ADD COLUMN IF NOT EXISTS current_phase_id INTEGER REFERENCES categories(id),
ADD COLUMN IF NOT EXISTS portfolio_category_id INTEGER REFERENCES categories(id),
ADD COLUMN IF NOT EXISTS overall_health_id INTEGER REFERENCES categories(id),
ADD COLUMN IF NOT EXISTS priority_level_id INTEGER REFERENCES categories(id);

-- Task
ALTER TABLE tasks
ADD COLUMN IF NOT EXISTS status_id INTEGER REFERENCES categories(id),
ADD COLUMN IF NOT EXISTS priority_id INTEGER REFERENCES categories(id);

-- Risk
ALTER TABLE risks
ADD COLUMN IF NOT EXISTS status_id INTEGER REFERENCES categories(id);

-- Category
ALTER TABLE categories
ADD COLUMN IF NOT EXISTS parent_id INTEGER REFERENCES categories(id) ON DELETE CASCADE;

ALTER TABLE stakeholders
ADD COLUMN IF NOT EXISTS role_id INTEGER REFERENCES categories(id);

-- Holiday
ALTER TABLE holidays
ADD COLUMN IF NOT EXISTS type_id INTEGER REFERENCES categories(id);

-- 4. Initial Data Mapping (Attempt to map existing values)
-- This is best-effort. In a production env, we'd do precise mapping.
UPDATE projects p SET project_status_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'PROJECT_STATUS') AND LOWER(name) = LOWER(p.project_status)) WHERE project_status IS NOT NULL;
UPDATE tasks t SET status_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'TASK_STATUS') AND LOWER(name) = LOWER(REPLACE(t.status, '_', ' '))) WHERE status IS NOT NULL;
UPDATE tasks t SET priority_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'TASK_PRIORITY') AND LOWER(name) = LOWER(t.priority)) WHERE priority IS NOT NULL;
UPDATE risks r SET status_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'RISK_STATUS') AND LOWER(name) = LOWER(r.status)) WHERE status IS NOT NULL;
UPDATE stakeholders s SET role_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'STAKEHOLDER_ROLE') AND LOWER(name) = LOWER(s.role)) WHERE role IS NOT NULL;
UPDATE holidays h SET type_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'HOLIDAY_TYPE') AND LOWER(name) = LOWER(h.type)) WHERE type IS NOT NULL;

-- 5. Drop old columns (Optional: can be done in a later step for safety, but we'll do it here for clean refactor)
-- ALTER TABLE projects DROP COLUMN project_status, DROP COLUMN current_phase, DROP COLUMN portfolio_category, DROP COLUMN overall_health;
-- ALTER TABLE tasks DROP COLUMN status, DROP COLUMN priority;
-- ALTER TABLE risks DROP COLUMN status;
-- ALTER TABLE stakeholders DROP COLUMN role;
-- ALTER TABLE holidays DROP COLUMN type;
