-- =====================================================
-- RBAC User Maintenance Script
-- Version: 2.0 (Consolidated)
-- Date: 2025-11-08
-- =====================================================
--
-- This script consolidates best practices from:
-- - comprehensive_staff_fix.sql
-- - sync_production_users_with_rbac.sql
-- - diagnose_staff_permissions.sql
--
-- PURPOSE:
-- - Sync legacy roles to RBAC
-- - Assign domain-based roles (@dekoninklijkeloop.nl → staff)
-- - Verify user-role assignments
-- - Diagnose permission issues
--
-- SAFE: Idempotent, can be run multiple times
-- =====================================================

\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo 'RBAC USER MAINTENANCE & SYNCHRONIZATION'
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo ''

BEGIN;

-- =====================================================
-- SECTION 1: Diagnostic Pre-Check
-- =====================================================
\echo 'SECTION 1: Pre-Check Diagnostics...'
\echo ''

DO $$
DECLARE
    total_users INTEGER;
    users_with_rbac INTEGER;
    users_without_rbac INTEGER;
    system_roles_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_users FROM gebruikers WHERE is_actief = true;
    SELECT COUNT(DISTINCT user_id) INTO users_with_rbac FROM user_roles WHERE is_active = true;
    SELECT COUNT(*) INTO system_roles_count FROM roles WHERE is_system_role = true;
    users_without_rbac := total_users - users_with_rbac;
    
    RAISE NOTICE '=== Current State ===';
    RAISE NOTICE 'Total active users: %', total_users;
    RAISE NOTICE 'Users with RBAC roles: %', users_with_rbac;
    RAISE NOTICE 'Users without RBAC: %', users_without_rbac;
    RAISE NOTICE 'System roles available: % (expected: 9)', system_roles_count;
    RAISE NOTICE '';
    
    IF users_without_rbac > 0 THEN
        RAISE NOTICE '⚠ % users need RBAC role assignment', users_without_rbac;
    ELSE
        RAISE NOTICE '✓ All active users have RBAC roles';
    END IF;
    RAISE NOTICE '';
END $$;

-- =====================================================
-- SECTION 2: Ensure Core System Roles Exist
-- =====================================================
\echo 'SECTION 2: Verifying system roles...'
\echo ''

INSERT INTO roles (name, description, is_system_role) VALUES
('staff', 'Ondersteunend personeel met beperkte beheerrechten', true)
ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (name, description, is_system_role) VALUES
('user', 'Standaard gebruiker met basis rechten', true)
ON CONFLICT (name) DO NOTHING;

-- Verify roles exist
SELECT 
    name,
    '✓ EXISTS' as status,
    description
FROM roles
WHERE name IN ('admin', 'staff', 'user', 'deelnemer', 'begeleider')
ORDER BY name;

\echo ''

-- =====================================================
-- SECTION 3: Ensure Staff Permissions Exist
-- =====================================================
\echo 'SECTION 3: Verifying staff permissions...'
\echo ''

-- Ensure participant permissions exist
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('participant', 'read', 'Participants bekijken', true),
('participant', 'write', 'Participants bewerken', true),
('participant', 'delete', 'Participants verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign to staff role (read + write only)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'staff'
  AND r.is_system_role = true
  AND p.resource = 'participant'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign all participant permissions to admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'admin'
  AND r.is_system_role = true
  AND p.resource = 'participant'
ON CONFLICT (role_id, permission_id) DO NOTHING;

\echo 'Staff permissions verified'
\echo ''

-- =====================================================
-- SECTION 4: Migrate Legacy Staff Users
-- =====================================================
\echo 'SECTION 4: Migrating legacy staff users...'
\echo ''

INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    COALESCE(g.created_at, NOW()) as assigned_at
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'staff' 
  AND r.is_system_role = true
  AND g.rol = 'staff'
  AND g.is_actief = true
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = true;

-- Count migrated users
DO $$
DECLARE
    migrated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers g ON ur.user_id = g.id
    WHERE r.name = 'staff' 
      AND g.rol = 'staff'
      AND ur.is_active = true;
    
    IF migrated_count > 0 THEN
        RAISE NOTICE '✓ Migrated % legacy staff users to RBAC', migrated_count;
    ELSE
        RAISE NOTICE 'ℹ No legacy staff users to migrate';
    END IF;
END $$;

\echo ''

-- =====================================================
-- SECTION 5: Domain-Based Role Assignment
-- =====================================================
\echo 'SECTION 5: Assigning domain-based roles...'
\echo ''

-- Assign staff role to all @dekoninklijkeloop.nl users
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    NOW() as assigned_at
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'staff' 
  AND r.is_system_role = true
  AND g.email LIKE '%@dekoninklijkeloop.nl'
  AND g.is_actief = true
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = true;

-- Report domain assignments
DO $$
DECLARE
    dkl_staff_count INTEGER;
BEGIN
    SELECT COUNT(DISTINCT ur.user_id) INTO dkl_staff_count
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers g ON ur.user_id = g.id
    WHERE r.name = 'staff' 
      AND g.email LIKE '%@dekoninklijkeloop.nl'
      AND ur.is_active = true;
    
    RAISE NOTICE '✓ @dekoninklijkeloop.nl users with staff role: %', dkl_staff_count;
END $$;

\echo ''

-- =====================================================
-- SECTION 6: Assign Event Participant Roles
-- =====================================================
\echo 'SECTION 6: Syncing event participant roles...'
\echo ''

-- Begeleiders
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT g.id, r.id, true, NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE LOWER(g.rol) = 'begeleider'
  AND r.name = 'begeleider'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = true;

-- Deelnemers
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT g.id, r.id, true, NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE LOWER(g.rol) = 'deelnemer'
  AND r.name = 'deelnemer'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = true;

-- Default user role for users without specific role
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT g.id, r.id, true, NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE (g.rol IS NULL OR g.rol = '' OR LOWER(g.rol) IN ('gebruiker', 'socialmedia'))
  AND r.name = 'user'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur WHERE ur.user_id = g.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = true;

\echo 'Event participant roles synced'
\echo ''

-- =====================================================
-- SECTION 7: Verification Report
-- =====================================================
\echo 'SECTION 7: Verification Report'
\echo ''

-- Users with staff role
\echo '=== Staff Users ==='
SELECT 
    u.email,
    u.naam,
    ur.assigned_at,
    CASE 
        WHEN ur.expires_at IS NULL THEN 'PERMANENT'
        WHEN ur.expires_at > NOW() THEN 'ACTIVE'
        ELSE 'EXPIRED'
    END as status
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE r.name = 'staff' AND ur.is_active = true
ORDER BY u.email;

\echo ''
\echo '=== Staff Permissions ==='
SELECT 
    p.resource,
    p.action,
    p.description
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff'
ORDER BY p.resource, p.action;

\echo ''
\echo '=== Role Distribution ==='
SELECT 
    r.name as role_name,
    COUNT(DISTINCT ur.user_id) as user_count,
    STRING_AGG(DISTINCT g.email, ', ' ORDER BY g.email) as sample_users
FROM roles r
LEFT JOIN user_roles ur ON r.id = ur.role_id AND ur.is_active = true
LEFT JOIN gebruikers g ON ur.user_id = g.id
WHERE r.is_system_role = true
GROUP BY r.name
ORDER BY user_count DESC, r.name;

\echo ''

-- =====================================================
-- SECTION 8: Final Status Report
-- =====================================================
\echo ''
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo 'FINAL STATUS REPORT'
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo ''

DO $$
DECLARE
    total_users INTEGER;
    users_with_rbac INTEGER;
    staff_users INTEGER;
    dkl_staff INTEGER;
    event_roles INTEGER;
    staff_has_read BOOLEAN;
    staff_has_write BOOLEAN;
BEGIN
    SELECT COUNT(*) INTO total_users FROM gebruikers WHERE is_actief = true;
    SELECT COUNT(DISTINCT user_id) INTO users_with_rbac FROM user_roles WHERE is_active = true;
    
    SELECT COUNT(DISTINCT ur.user_id) INTO staff_users
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE r.name = 'staff' AND ur.is_active = true;
    
    SELECT COUNT(DISTINCT ur.user_id) INTO dkl_staff
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers g ON ur.user_id = g.id
    WHERE r.name = 'staff' 
      AND g.email LIKE '%@dekoninklijkeloop.nl'
      AND ur.is_active = true;
    
    SELECT COUNT(DISTINCT ur.user_id) INTO event_roles
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE r.name IN ('deelnemer', 'begeleider', 'vrijwilliger')
      AND ur.is_active = true;
    
    SELECT EXISTS(
        SELECT 1 FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE r.name = 'staff' AND p.resource = 'participant' AND p.action = 'read'
    ) INTO staff_has_read;
    
    SELECT EXISTS(
        SELECT 1 FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE r.name = 'staff' AND p.resource = 'participant' AND p.action = 'write'
    ) INTO staff_has_write;
    
    RAISE NOTICE '=== RESULTS ===';
    RAISE NOTICE '';
    RAISE NOTICE 'Total active users: %', total_users;
    RAISE NOTICE 'Users with RBAC roles: %', users_with_rbac;
    RAISE NOTICE 'Coverage: %%%', ROUND(100.0 * users_with_rbac / NULLIF(total_users, 0), 1);
    RAISE NOTICE '';
    RAISE NOTICE 'Staff role users: %', staff_users;
    RAISE NOTICE '@dekoninklijkeloop.nl staff: %', dkl_staff;
    RAISE NOTICE 'Event participants: %', event_roles;
    RAISE NOTICE '';
    
    IF staff_has_read THEN
        RAISE NOTICE '✓ Staff has participant:read permission';
    ELSE
        RAISE NOTICE '✗ Staff MISSING participant:read permission';
    END IF;
    
    IF staff_has_write THEN
        RAISE NOTICE '✓ Staff has participant:write permission';
    ELSE
        RAISE NOTICE '✗ Staff MISSING participant:write permission';
    END IF;
    
    RAISE NOTICE '';
    
    IF users_with_rbac >= total_users AND staff_has_read AND staff_has_write THEN
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE '✓✓✓ ALL CHECKS PASSED ✓✓✓';
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE '';
        RAISE NOTICE 'NEXT STEPS:';
        RAISE NOTICE '1. Staff users should logout/login to refresh JWT tokens';
        RAISE NOTICE '2. Test access to /api/participant endpoints';
        RAISE NOTICE '3. Run verify_rbac_tables.sql for detailed verification';
    ELSE
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE '⚠ ISSUES DETECTED - REVIEW NEEDED';
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        
        IF users_with_rbac < total_users THEN
            RAISE NOTICE '• % users still need RBAC roles', (total_users - users_with_rbac);
        END IF;
        
        IF NOT staff_has_read OR NOT staff_has_write THEN
            RAISE NOTICE '• Staff role missing required permissions';
        END IF;
    END IF;
    
    RAISE NOTICE '';
END $$;

COMMIT;

\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo 'RBAC USER MAINTENANCE COMPLETED'
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo ''