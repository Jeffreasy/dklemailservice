-- =====================================================
-- V34 FIX: Refresh Tokens Foreign Key Constraint
-- =====================================================
-- Issue: refresh_tokens.user_id has FK to gebruikers.id
--        but now also needs to accept participant.id
-- Solution: Remove FK constraint, rename column for clarity
-- Date: 2025-11-10
-- =====================================================

BEGIN;

-- Step 1: Remove the foreign key constraint
ALTER TABLE refresh_tokens 
DROP CONSTRAINT IF EXISTS refresh_tokens_user_id_fkey;

-- Step 2: Rename column for clarity
ALTER TABLE refresh_tokens 
RENAME COLUMN user_id TO owner_id;

-- Step 3: Add comment explaining the change
COMMENT ON COLUMN refresh_tokens.owner_id IS 
'Can reference either gebruikers.id (admin/staff) or participants.id (deelnemers). No FK constraint to support both types.';

-- Step 4: Add index for performance
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_owner_id 
ON refresh_tokens(owner_id);

-- Step 5: Verify the change
DO $$
BEGIN
    -- Check if column exists and FK is removed
    IF EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'refresh_tokens' 
        AND column_name = 'owner_id'
    ) THEN
        RAISE NOTICE '✅ Column renamed to owner_id';
    ELSE
        RAISE WARNING '❌ Column rename failed';
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.table_constraints 
        WHERE constraint_name = 'refresh_tokens_user_id_fkey'
    ) THEN
        RAISE NOTICE '✅ Foreign key constraint removed';
    ELSE
        RAISE WARNING '⚠️ Foreign key constraint still exists';
    END IF;
END $$;

COMMIT;

-- Display table structure for verification
\d refresh_tokens;

-- Show sample data (if any)
SELECT 
    id,
    owner_id,
    LEFT(token, 20) || '...' as token_preview,
    expires_at,
    created_at,
    is_revoked
FROM refresh_tokens
LIMIT 5;

-- Success message
SELECT '✅ V34 Refresh Tokens Fix Applied Successfully!' as status;