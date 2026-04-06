-- Add recurrence columns to tasks table
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS recurrence_type TEXT DEFAULT 'none';
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS recurrence_config JSONB;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS recurrence_end_date TIMESTAMPTZ;

-- Create index for recurrence type for faster queries
CREATE INDEX IF NOT EXISTS idx_tasks_recurrence_type ON tasks (recurrence_type);

-- Create a trigger function to validate recurrence_type
CREATE OR REPLACE FUNCTION validate_recurrence_type()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.recurrence_type NOT IN ('none', 'daily', 'monthly', 'weekly', 'odd_even', 'specific') THEN
        RAISE EXCEPTION 'Invalid recurrence_type: %', NEW.recurrence_type;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to check recurrence_type on insert/update
DROP TRIGGER IF EXISTS trg_validate_recurrence_type ON tasks;
CREATE TRIGGER trg_validate_recurrence_type
BEFORE INSERT OR UPDATE ON tasks
FOR EACH ROW
EXECUTE FUNCTION validate_recurrence_type();
