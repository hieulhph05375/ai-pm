-- 000013_upgrade_tasks_personal.up.sql
-- Upgrade tasks table for personal task management (ADR-011)
-- Drop old columns tied to WBS/Project, add new personal task fields

ALTER TABLE tasks
    DROP COLUMN IF EXISTS wbs_code,
    DROP COLUMN IF EXISTS project_id,
    ADD COLUMN IF NOT EXISTS start_date TIMESTAMP,
    ADD COLUMN IF NOT EXISTS labels JSONB DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS created_by INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Normalize status values to match new constants
UPDATE tasks SET status = 'TODO' WHERE status = 'Backlog' OR status NOT IN ('TODO', 'IN_PROGRESS', 'DONE');
UPDATE tasks SET priority = 'MEDIUM' WHERE priority = 'Medium' OR priority NOT IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT');

-- Create task_activities table for audit log
CREATE TABLE IF NOT EXISTS task_activities (
    id          SERIAL PRIMARY KEY,
    task_id     INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    actor_id    INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action      VARCHAR(100) NOT NULL,  -- e.g. 'created', 'status_changed', 'commented'
    old_value   TEXT,
    new_value   TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_task_activities_task_id ON task_activities(task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks(created_by);
CREATE INDEX IF NOT EXISTS idx_tasks_start_date ON tasks(start_date);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
