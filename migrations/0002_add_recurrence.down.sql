-- Rollback: Remove recurrence columns from tasks table
DROP TRIGGER IF EXISTS trg_validate_recurrence_type ON tasks;
DROP FUNCTION IF EXISTS validate_recurrence_type();

ALTER TABLE tasks DROP COLUMN IF EXISTS recurrence_type;
ALTER TABLE tasks DROP COLUMN IF EXISTS recurrence_config;
ALTER TABLE tasks DROP COLUMN IF EXISTS recurrence_end_date;
