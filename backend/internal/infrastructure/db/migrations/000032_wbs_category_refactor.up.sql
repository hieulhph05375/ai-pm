-- 000032_wbs_category_refactor.up.sql

-- 1. Create Category Type for WBS Node Types
INSERT INTO category_types (name, code, description) VALUES
('WBS Node Type', 'WBS_NODE_TYPE', 'Classification for WBS elements (Phase, Milestone, Task)')
ON CONFLICT (code) DO NOTHING;

-- 2. Seed Categories for WBS Node Types
INSERT INTO categories (type_id, name, color)
SELECT id, 'Phase', 'slate' FROM category_types WHERE code = 'WBS_NODE_TYPE' UNION ALL
SELECT id, 'Milestone', 'amber' FROM category_types WHERE code = 'WBS_NODE_TYPE' UNION ALL
SELECT id, 'Work Package', 'blue' FROM category_types WHERE code = 'WBS_NODE_TYPE'
ON CONFLICT (type_id, name) DO NOTHING;

-- 3. Add type_id column to wbs_nodes
ALTER TABLE wbs_nodes ADD COLUMN IF NOT EXISTS type_id INTEGER REFERENCES categories(id);

-- 4. Map existing data
-- Note: 'Task' and 'Sub-task' are mapped to 'Work Package' for standard WBS consistency
UPDATE wbs_nodes SET type_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'WBS_NODE_TYPE') AND name = 'Phase') WHERE type = 'Phase';
UPDATE wbs_nodes SET type_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'WBS_NODE_TYPE') AND name = 'Milestone') WHERE type = 'Milestone';
UPDATE wbs_nodes SET type_id = (SELECT id FROM categories WHERE type_id = (SELECT id FROM category_types WHERE code = 'WBS_NODE_TYPE') AND name = 'Work Package') WHERE type IN ('Task', 'Sub-task');
