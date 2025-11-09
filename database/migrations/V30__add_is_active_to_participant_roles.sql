-- V30: Add is_active column to participant_roles
-- This field is required by the repository layer queries

ALTER TABLE participant_roles 
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Set existing rows to active
UPDATE participant_roles 
SET is_active = true 
WHERE is_active IS NULL;

-- Add comment
COMMENT ON COLUMN participant_roles.is_active IS 'Indicates if this role is currently active and available for selection';