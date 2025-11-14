-- =====================================================
-- V34 ROLLBACK: Refresh Tokens Foreign Key Constraint
-- =====================================================
-- WARNING: This will delete all participant refresh tokens!
-- Only execute if you need to rollback the V34 fix
-- Date: 2025-11-10
-- =====================================================

BEGIN;

-- Step 1: Delete all refresh tokens that reference participants
-- (Cannot have FK constraint with participant IDs in the table)
DELETE FROM refresh_tokens 
WHERE owner_id IN (
    SELECT id FROM participants
);

-- Step 2: Rename column back to user_id
ALTER TABLE refresh_tokens 
RENAME COLUMN owner_id TO user_id;

-- Step 3: Re-add foreign key constraint
ALTER TABLE refresh_tokens
ADD CONSTRAINT refresh_tokens_user_id_fkey
FOREIGN KEY (user_id) REFERENCES gebruikers(id) ON DELETE CASCADE;

-- Step 4: Remove the owner_id index
DROP INDEX IF EXISTS idx_refresh_tokens_owner_id;

-- Step 5: Verify rollback
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'refresh_tokens' 
        AND column_name = 'user_id'
    ) THEN
        RAISE NOTICE '✅ Column renamed back to user_id';
    ELSE
        RAISE WARNING '❌ Column rename rollback failed';
    END IF;
    
    IF EXISTS (
        SELECT 1 
        FROM information_schema.table_constraints 
        WHERE constraint_name = 'refresh_tokens_user_id_fkey'
    ) THEN
        RAISE NOTICE '✅ Foreign key constraint restored';
    ELSE
        RAISE WARNING '⚠️ Foreign key constraint restoration failed';
    END IF;
END $$;

COMMIT;

-- Display table structure for verification
\d refresh_tokens;

-- Show remaining tokens
SELECT 
    COUNT(*) as token_count,
    'Only gebruiker tokens remain' as note
FROM refresh_tokens;

-- Warning message
SELECT '⚠️ V34 Refresh Tokens Rollback Complete - All participant tokens deleted!' as status;