-- 000003_create_tasks_table.down.sql

DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP TABLE IF EXISTS tasks;
