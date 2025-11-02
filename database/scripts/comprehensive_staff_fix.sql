-- COMPREHENSIVE FIX: Complete staff role and permissions setup
-- This script checks and fixes ALL issues with staff role access to aanmeldingen
-- Safe to run multiple times (idempotent)

\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo 'COMPREHENSIVE STAFF ROLE FIX'
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'

-- ================================================
-- STEP 1: Ensure staff role exists
-- ================================================
\echo ''
\echo 'STEP 1: Ensuring staff role exists...'

INSERT INTO roles (name, description, is_system_role) VALUES
('staff', 'Ondersteunend personeel met beperkte beheerrechten', true)
ON CONFLICT (name) DO NOTHING;

DO $$
BEGIN
    IF FOUND THEN
        RAISE NOTICE '✓ Staff role created';
    ELSE
        RAISE NOTICE '✓ Staff role already exists';
    END IF;
END $$;

-- ================================================
-- STEP 2: Ensure aanmelding permissions exist
-- ================================================
\echo ''
\echo 'STEP 2: Ensuring aanmelding permissions exist...'

INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken', true),
('aanmelding', 'delete', 'Aanmeldingen verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

SELECT 
    p.resource,
    p.action,
    '✓ Exists' as status
FROM permissions p
WHERE p.resource = 'aanmelding'
ORDER BY p.action;

-- ================================================
-- STEP 3: Assign aanmelding permissions to staff role
-- ================================================
\echo ''
\echo 'STEP 3: Assigning aanmelding permissions to staff role...'

-- Give staff read AND write permissions (not delete)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' 
  AND r.is_system_role = true
  AND p.resource = 'aanmelding'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

\echo 'Staff role aanmelding permissions:'
SELECT 
    r.name as role_name,
    p.resource,
    p.action,
    '✓ Assigned' as status
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff' 
  AND p.resource = 'aanmelding'
ORDER BY p.action;

-- ================================================
-- STEP 4: Migrate legacy staff users to RBAC
-- ================================================
\echo ''
\echo 'STEP 4: Migrating legacy staff users to RBAC...'

-- Migrate users with gebruikers.rol = 'staff' to RBAC
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
  -- Prevent duplicates
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id 
      AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ================================================
-- STEP 5: Assign staff role to @dekoninklijkeloop.com users
-- ================================================
\echo ''
\echo 'STEP 5: Assigning staff role to @dekoninklijkeloop.com domain users...'

-- Assign staff role to all @dekoninklijkeloop.com users who don't have it
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
  AND g.email LIKE '%@dekoninklijkeloop.com'
  AND g.is_actief = true
  -- Prevent duplicates
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id 
      AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ================================================
-- STEP 6: Verification - Show all staff users
-- ================================================
\echo ''
\echo 'STEP 6: Verification - All users with staff role:'

SELECT 
    u.email,
    u.naam,
    u.rol as legacy_role,
    u.is_actief,
    ur.is_active as rbac_active,
    ur.assigned_at,
    CASE 
        WHEN ur.expires_at IS NULL THEN 'PERMANENT'
        WHEN ur.expires_at > NOW() THEN 'ACTIVE'
        ELSE 'EXPIRED'
    END as expiry_status
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE r.name = 'staff'
ORDER BY u.email;

-- ================================================
-- STEP 7: Verification - Staff permissions summary
-- ================================================
\echo ''
\echo 'STEP 7: Staff permissions summary:'

SELECT 
    p.resource,
    p.action,
    p.description,
    '✓ GRANTED' as status
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff'
ORDER BY p.resource, p.action;

-- ================================================
-- STEP 8: Final Status Report
-- ================================================
\echo ''
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo 'FINAL STATUS REPORT'
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'

DO $$
DECLARE
    staff_role_exists BOOLEAN;
    aanmelding_read_exists BOOLEAN;
    aanmelding_write_exists BOOLEAN;
    staff_has_read BOOLEAN;
    staff_has_write BOOLEAN;
    users_with_staff INTEGER;
    dkl_users_with_staff INTEGER;
    dkl_users_total INTEGER;
BEGIN
    -- Check staff role
    SELECT EXISTS(SELECT 1 FROM roles WHERE name = 'staff') INTO staff_role_exists;
    
    -- Check aanmelding permissions exist
    SELECT EXISTS(SELECT 1 FROM permissions WHERE resource = 'aanmelding' AND action = 'read') 
    INTO aanmelding_read_exists;
    
    SELECT EXISTS(SELECT 1 FROM permissions WHERE resource = 'aanmelding' AND action = 'write') 
    INTO aanmelding_write_exists;
    
    -- Check staff has permissions
    SELECT EXISTS(
        SELECT 1 FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE r.name = 'staff' AND p.resource = 'aanmelding' AND p.action = 'read'
    ) INTO staff_has_read;
    
    SELECT EXISTS(
        SELECT 1 FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE r.name = 'staff' AND p.resource = 'aanmelding' AND p.action = 'write'
    ) INTO staff_has_write;
    
    -- Count users
    SELECT COUNT(DISTINCT ur.user_id) INTO users_with_staff
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE r.name = 'staff' AND ur.is_active = true;
    
    SELECT COUNT(DISTINCT ur.user_id) INTO dkl_users_with_staff
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers u ON ur.user_id = u.id
    WHERE r.name = 'staff' 
      AND ur.is_active = true
      AND u.email LIKE '%@dekoninklijkeloop.com';
    
    SELECT COUNT(*) INTO dkl_users_total
    FROM gebruikers 
    WHERE email LIKE '%@dekoninklijkeloop.com' 
      AND is_actief = true;
    
    -- Print report
    RAISE NOTICE '';
    RAISE NOTICE '=== ROLE & PERMISSIONS ===';
    
    IF staff_role_exists THEN
        RAISE NOTICE '✓ Staff role exists';
    ELSE
        RAISE NOTICE '✗ Staff role MISSING';
    END IF;
    
    IF aanmelding_read_exists THEN
        RAISE NOTICE '✓ aanmelding:read permission exists';
    ELSE
        RAISE NOTICE '✗ aanmelding:read permission MISSING';
    END IF;
    
    IF aanmelding_write_exists THEN
        RAISE NOTICE '✓ aanmelding:write permission exists';
    ELSE
        RAISE NOTICE '✗ aanmelding:write permission MISSING';
    END IF;
    
    IF staff_has_read THEN
        RAISE NOTICE '✓ Staff role HAS aanmelding:read';
    ELSE
        RAISE NOTICE '✗ Staff role MISSING aanmelding:read';
    END IF;
    
    IF staff_has_write THEN
        RAISE NOTICE '✓ Staff role HAS aanmelding:write';
    ELSE
        RAISE NOTICE '✗ Staff role MISSING aanmelding:write';
    END IF;
    
    RAISE NOTICE '';
    RAISE NOTICE '=== USER ASSIGNMENTS ===';
    RAISE NOTICE '• Total users with staff role: %', users_with_staff;
    RAISE NOTICE '• @dekoninklijkeloop.com users with staff: %', dkl_users_with_staff;
    RAISE NOTICE '• Total active @dekoninklijkeloop.com users: %', dkl_users_total;
    
    IF dkl_users_with_staff < dkl_users_total THEN
        RAISE NOTICE '⚠ % @dekoninklijkeloop.com users still need staff role', 
            (dkl_users_total - dkl_users_with_staff);
    END IF;
    
    RAISE NOTICE '';
    
    IF staff_role_exists AND staff_has_read AND staff_has_write AND users_with_staff > 0 THEN
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE '✓✓✓ ALL CHECKS PASSED ✓✓✓';
        RAISE NOTICE 'Staff users should now have access to aanmeldingen!';
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE '';
        RAISE NOTICE 'NEXT STEPS:';
        RAISE NOTICE '1. Staff users may need to logout/login to refresh permissions';
        RAISE NOTICE '2. Test by accessing /api/aanmeldingen as a staff user';
        RAISE NOTICE '3. If still 403, check JWT token contains staff role';
    ELSE
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE '⚠ ISSUES FOUND - MANUAL INTERVENTION NEEDED';
        RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
        RAISE NOTICE 'Please review the diagnostic output above';
    END IF;
    
    RAISE NOTICE '';
END $$;

\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
\echo 'Script completed. Review output above for results.'
\echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'