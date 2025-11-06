-- Fix V1_58 Migration Issue
-- This script cleans up any conflicting views before re-running migrations

-- Drop any existing notulen views that might conflict
DROP VIEW IF EXISTS notulen_with_participants CASCADE;
DROP VIEW IF EXISTS notulen_versies_with_participants CASCADE;

-- Verify the updated_by column exists (it should from V1_54)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'notulen' 
        AND column_name = 'updated_by'
    ) THEN
        ALTER TABLE notulen ADD COLUMN updated_by UUID REFERENCES gebruikers(id);
        RAISE NOTICE 'Added updated_by column to notulen table';
    ELSE
        RAISE NOTICE 'Column updated_by already exists in notulen table';
    END IF;
END $$;

-- Now the database is ready for V1_58 to run cleanly
SELECT 'Database prepared for V1_58 migration' as status;