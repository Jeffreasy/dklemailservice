-- Fix schema issues identified in the logs
-- This script addresses the column mismatches causing SQL errors

-- 1. Add file_path column to uploaded_images table (if needed for compatibility)
-- Note: The table already has the necessary columns, so this query might be outdated

-- 2. The events table uses 'name' instead of 'title' - this is correct in schema
-- The query should be updated to use 'name' instead of 'title'

-- 3. The email_templates table uses 'naam' instead of 'name' - this is correct in schema
-- The query should be updated to use 'naam' instead of 'name'

-- 4. Fix ambiguous column reference in JOIN query
-- The query should prefix columns with table aliases

-- Since we can't modify the source code directly, let's add the missing columns
-- to maintain backward compatibility

-- Add file_path column to uploaded_images if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'uploaded_images'
                   AND column_name = 'file_path') THEN
        ALTER TABLE uploaded_images ADD COLUMN file_path TEXT;
    END IF;
END $$;

-- Add title column to events if it doesn't exist (for backward compatibility)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'events'
                   AND column_name = 'title') THEN
        ALTER TABLE events ADD COLUMN title TEXT;
        -- Copy data from name to title for compatibility
        UPDATE events SET title = name WHERE title IS NULL;
    END IF;
END $$;

-- Add name column to email_templates if it doesn't exist (for backward compatibility)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'email_templates'
                   AND column_name = 'name') THEN
        ALTER TABLE email_templates ADD COLUMN name TEXT;
        -- Copy data from naam to name for compatibility
        UPDATE email_templates SET name = naam WHERE name IS NULL;
    END IF;
END $$;

-- Add event_date column to events if it doesn't exist (for backward compatibility)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'events'
                   AND column_name = 'event_date') THEN
        ALTER TABLE events ADD COLUMN event_date TIMESTAMP WITH TIME ZONE;
        -- Copy data from start_time to event_date for compatibility
        UPDATE events SET event_date = start_time WHERE event_date IS NULL;
    END IF;
END $$;

-- Add subject column to email_templates if it doesn't exist (for backward compatibility)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'email_templates'
                   AND column_name = 'subject') THEN
        ALTER TABLE email_templates ADD COLUMN subject TEXT;
        -- Copy data from onderwerp to subject for compatibility
        UPDATE email_templates SET subject = onderwerp WHERE subject IS NULL;
    END IF;
END $$;