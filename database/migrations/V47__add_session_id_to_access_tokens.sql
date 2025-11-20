-- V47: Add session_id column to access_tokens table for session management
-- This migration adds session tracking to access tokens

-- Add session_id column to access_tokens table (nullable initially, idempotent)
DO $$
BEGIN
    -- Add session_id column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'access_tokens' AND column_name = 'session_id') THEN
        ALTER TABLE access_tokens ADD COLUMN session_id UUID NULL;
    END IF;

    -- Create index on session_id for performance (only if it doesn't exist)
    IF NOT EXISTS (SELECT 1 FROM pg_indexes
                   WHERE tablename = 'access_tokens' AND indexname = 'idx_access_tokens_session_id') THEN
        CREATE INDEX idx_access_tokens_session_id ON access_tokens(session_id);
    END IF;

    -- Add comment to document the column (only if column exists)
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'access_tokens' AND column_name = 'session_id') THEN
        COMMENT ON COLUMN access_tokens.session_id IS 'References the session this access token belongs to for session management';
    END IF;
END $$;

-- Add foreign key constraint to sessions table
DO $$
BEGIN
    -- Add foreign key constraint if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                   WHERE table_name = 'access_tokens' AND constraint_name = 'fk_access_tokens_session_id') THEN
        ALTER TABLE access_tokens ADD CONSTRAINT fk_access_tokens_session_id
        FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE;
    END IF;
END $$;