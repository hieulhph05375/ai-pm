-- 000014_add_task_progress.up.sql
-- Add progress column to tasks for Phase 2 enhancement

ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS progress INTEGER DEFAULT 0 CHECK (progress >= 0 AND progress <= 100);
