-- V38: Remove legacy rol column from gebruikers table
-- Date: 2025-11-17
-- Purpose: Complete legacy authorization system removal - drop deprecated rol column
-- Status: BREAKING CHANGE - Legacy role field no longer supported

-- IMPORTANT: This migration removes the legacy 'rol' column from gebruikers table
-- All authorization now uses RBAC system (user_roles, roles, role_permissions tables)

-- Check if migration is already complete (column doesn't exist) and execute entire migration conditionally
DO $$
DECLARE
    rol_column_exists BOOLEAN;
    rbac_users_count INTEGER;
    legacy_users_count INTEGER;
BEGIN
    -- Check if rol column exists
    SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'gebruikers'
        AND column_name = 'rol'
    ) INTO rol_column_exists;

    IF NOT rol_column_exists THEN
        RAISE NOTICE 'V38: Migration already completed - rol column does not exist';
        RAISE NOTICE 'V38: Legacy authorization system removal completed';
        RAISE NOTICE 'V38: System now uses RBAC-only authorization';
        RETURN;
    END IF;

    -- Count users with RBAC roles
    SELECT COUNT(DISTINCT ur.user_id) INTO rbac_users_count
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE ur.is_active = true;

    -- Count users with legacy rol values (non-empty)
    SELECT COUNT(*) INTO legacy_users_count
    FROM gebruikers
    WHERE rol IS NOT NULL AND rol != '';

    RAISE NOTICE 'V38: Users with RBAC roles: %, Users with legacy rol: %', rbac_users_count, legacy_users_count;

    -- Warning if there are users without RBAC roles but with legacy roles
    IF legacy_users_count > 0 AND rbac_users_count = 0 THEN
        RAISE EXCEPTION 'V38: CRITICAL - Found % users with legacy rol but no RBAC roles. Migration cannot proceed safely.', legacy_users_count;
    END IF;

    -- Log the migration status
    RAISE NOTICE 'V38: Legacy rol column removal proceeding - RBAC system verified';
END $$;

-- Execute the rest of the migration only if column exists
DO $$
DECLARE
    rol_column_exists BOOLEAN;
BEGIN
    -- Check if rol column exists
    SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'gebruikers'
        AND column_name = 'rol'
    ) INTO rol_column_exists;

    IF rol_column_exists THEN
        -- CRITICAL FIX: Ensure admin role has ALL permissions before removing legacy system
        -- This addresses the issue where admin doesn't have access everywhere
        INSERT INTO role_permissions (role_id, permission_id)
        SELECT r.id, p.id
        FROM roles r, permissions p
        WHERE r.name = 'admin' AND r.is_system_role = true
        ON CONFLICT (role_id, permission_id) DO NOTHING;

        -- Log admin permissions fix
        DECLARE
            admin_perm_count INTEGER;
            total_perm_count INTEGER;
        BEGIN
            SELECT COUNT(*) INTO admin_perm_count
            FROM roles r
            JOIN role_permissions rp ON r.id = rp.role_id
            WHERE r.name = 'admin';

            SELECT COUNT(*) INTO total_perm_count
            FROM permissions;

            RAISE NOTICE '[V38] Admin role permissions fix: %/% permissions assigned', admin_perm_count, total_perm_count;

            IF admin_perm_count < total_perm_count THEN
                RAISE WARNING '[V38] WARNING: Admin role has fewer permissions than total system permissions!';
            END IF;
        END;

        -- Drop views that depend on the rol column first
        DROP VIEW IF EXISTS users_without_participation;
        DROP VIEW IF EXISTS v_user_role_migration_status;

        -- Drop the index that depends on the rol column
        DROP INDEX IF EXISTS idx_gebruikers_role_id;

        -- Drop the legacy rol column
        -- This is a BREAKING CHANGE - all code should use RBAC system now
        ALTER TABLE gebruikers DROP COLUMN IF EXISTS rol;

        -- Final verification
        DECLARE
            column_still_exists BOOLEAN;
        BEGIN
            SELECT EXISTS (
                SELECT 1
                FROM information_schema.columns
                WHERE table_name = 'gebruikers'
                AND column_name = 'rol'
            ) INTO column_still_exists;

            IF NOT column_still_exists THEN
                RAISE NOTICE 'V38: rol column successfully removed from gebruikers table';
                RAISE NOTICE 'V38: Legacy authorization system removal completed';
                RAISE NOTICE 'V38: System now uses RBAC-only authorization';
            ELSE
                RAISE EXCEPTION 'V38: Failed to remove rol column from gebruikers table';
            END IF;
        END;
    END IF;
END $$;