-- V48: Remove legacy role_id column from gebruikers table
-- Date: 2025-11-20
-- Purpose: Complete legacy authorization system removal - drop deprecated role_id column
-- Status: BREAKING CHANGE - Legacy role_id field no longer supported

-- IMPORTANT: This migration removes the legacy 'role_id' column from gebruikers table
-- All authorization now uses RBAC system (user_roles, roles, role_permissions tables)

-- Check if migration is already complete (column doesn't exist) and execute entire migration conditionally
DO $$
DECLARE
    role_id_column_exists BOOLEAN;
    rbac_users_count INTEGER;
    legacy_users_count INTEGER;
BEGIN
    -- Check if role_id column exists
    SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'gebruikers'
        AND column_name = 'role_id'
    ) INTO role_id_column_exists;

    IF NOT role_id_column_exists THEN
        RAISE NOTICE 'V48: Migration already completed - role_id column does not exist';
        RAISE NOTICE 'V48: Legacy authorization system removal completed';
        RAISE NOTICE 'V48: System now uses RBAC-only authorization';
        RETURN;
    END IF;

    -- Count users with RBAC roles
    SELECT COUNT(DISTINCT ur.user_id) INTO rbac_users_count
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE ur.is_active = true;

    -- Count users with legacy role_id values (non-null)
    SELECT COUNT(*) INTO legacy_users_count
    FROM gebruikers
    WHERE role_id IS NOT NULL;

    RAISE NOTICE 'V48: Users with RBAC roles: %, Users with legacy role_id: %', rbac_users_count, legacy_users_count;

    -- Warning if there are users without RBAC roles but with legacy role_id
    IF legacy_users_count > 0 AND rbac_users_count = 0 THEN
        RAISE EXCEPTION 'V48: CRITICAL - Found % users with legacy role_id but no RBAC roles. Migration cannot proceed safely.', legacy_users_count;
    END IF;

    -- Log the migration status
    RAISE NOTICE 'V48: Legacy role_id column removal proceeding - RBAC system verified';
END $$;

-- Execute the rest of the migration only if column exists
DO $$
DECLARE
    role_id_column_exists BOOLEAN;
BEGIN
    -- Check if role_id column exists
    SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'gebruikers'
        AND column_name = 'role_id'
    ) INTO role_id_column_exists;

    IF role_id_column_exists THEN
        -- Drop the index that depends on the role_id column
        DROP INDEX IF EXISTS idx_gebruikers_role_id;

        -- Drop the legacy role_id column
        -- This is a BREAKING CHANGE - all code should use RBAC system now
        ALTER TABLE gebruikers DROP COLUMN IF EXISTS role_id;

        -- Final verification
        DECLARE
            column_still_exists BOOLEAN;
        BEGIN
            SELECT EXISTS (
                SELECT 1
                FROM information_schema.columns
                WHERE table_name = 'gebruikers'
                AND column_name = 'role_id'
            ) INTO column_still_exists;

            IF NOT column_still_exists THEN
                RAISE NOTICE 'V48: role_id column successfully removed from gebruikers table';
                RAISE NOTICE 'V48: Legacy authorization system removal completed';
                RAISE NOTICE 'V48: System now uses RBAC-only authorization';
            ELSE
                RAISE EXCEPTION 'V48: Failed to remove role_id column from gebruikers table';
            END IF;
        END;
    END IF;
END $$;