-- DIAGNOSTIC: Check complete staff role and permission setup
-- Run this to understand why staff users get 403 on aanmeldingen

-- ================================================
-- PART 1: Check if staff role exists
-- ================================================
\echo '=== PART 1: Staff Role Existence ==='
SELECT 
    id,
    name,
    description,
    is_system_role,
    created_at
FROM roles
WHERE name = 'staff';

-- ================================================
-- PART 2: Check all permissions for staff role
-- ================================================
\echo ''
\echo '=== PART 2: All Permissions for Staff Role ==='
SELECT 
    r.name as role_name,
    p.resource,
    p.action,
    p.description,
    rp.assigned_at
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff'
ORDER BY p.resource, p.action;

-- ================================================
-- PART 3: Check specifically for aanmelding permissions
-- ================================================
\echo ''
\echo '=== PART 3: Aanmelding Permissions for Staff ==='
SELECT 
    r.name as role_name,
    p.resource,
    p.action,
    p.description,
    CASE 
        WHEN rp.id IS NOT NULL THEN '✓ TOEGEWEZEN'
        ELSE '✗ NIET TOEGEWEZEN'
    END as status
FROM roles r
CROSS JOIN permissions p
LEFT JOIN role_permissions rp ON r.id = rp.role_id AND p.id = rp.permission_id
WHERE r.name = 'staff' 
  AND p.resource = 'aanmelding'
ORDER BY p.action;

-- ================================================
-- PART 4: Check which users have staff role
-- ================================================
\echo ''
\echo '=== PART 4: Users with Staff Role ==='
SELECT 
    u.email,
    u.naam,
    r.name as role_name,
    ur.is_active,
    ur.assigned_at,
    CASE 
        WHEN ur.expires_at IS NULL THEN 'PERMANENT'
        WHEN ur.expires_at > NOW() THEN 'ACTIVE (expires: ' || ur.expires_at::text || ')'
        ELSE 'EXPIRED'
    END as expiry_status
FROM user_roles ur
JOIN gebruikers u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE r.name = 'staff'
ORDER BY ur.is_active DESC, u.email;

-- ================================================
-- PART 5: Check legacy vs RBAC roles
-- ================================================
\echo ''
\echo '=== PART 5: Legacy Role vs RBAC Role Mismatch ==='
SELECT 
    u.email,
    u.naam,
    u.rol as legacy_role,
    STRING_AGG(r.name, ', ' ORDER BY r.name) as rbac_roles,
    CASE 
        WHEN u.rol = 'staff' AND EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = u.id 
              AND r2.name = 'staff'
              AND ur2.is_active = true
        ) THEN '✓ MATCH'
        WHEN u.rol = 'staff' THEN '✗ LEGACY STAFF BUT NO RBAC'
        WHEN EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2 .id
            WHERE ur2.user_id = u.id 
              AND r2.name = 'staff'
              AND ur2.is_active = true
        ) THEN '✗ RBAC STAFF BUT LEGACY IS: ' || COALESCE(u.rol, 'NULL')
        ELSE 'NO STAFF ROLE'
    END as status
FROM gebruikers u
LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.rol = 'staff' 
   OR EXISTS (
       SELECT 1 FROM user_roles ur3
       JOIN roles r3 ON ur3.role_id = r3.id
       WHERE ur3.user_id = u.id 
         AND r3.name = 'staff'
         AND ur3.is_active = true
   )
GROUP BY u.id, u.email, u.naam, u.rol
ORDER BY u.email;

-- ================================================
-- PART 6: Check dekoninklijkeloop.com users specifically
-- ================================================
\echo ''
\echo '=== PART 6: dekoninklijkeloop.com Domain Users ==='
SELECT 
    u.email,
    u.naam,
    u.rol as legacy_role,
    STRING_AGG(DISTINCT r.name, ', ' ORDER BY r.name) as rbac_roles,
    CASE 
        WHEN COUNT(ur.id) = 0 THEN '✗ NO RBAC ROLES'
        WHEN EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = u.id 
              AND r2.name IN ('staff', 'admin')
              AND ur2.is_active = true
        ) THEN '✓ HAS STAFF/ADMIN ROLE'
        ELSE '✗ HAS RBAC BUT NOT STAFF/ADMIN'
    END as status
FROM gebruikers u
LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.email LIKE '%@dekoninklijkeloop.com'
GROUP BY u.id, u.email, u.naam, u.rol
ORDER BY u.email;

-- ================================================
-- PART 7: Check effective permissions for a specific user
-- ================================================
\echo ''
\echo '=== PART 7: Effective Permissions for User (replace email) ==='
\echo 'To check specific user, run: '
\echo 'SELECT p.resource, p.action FROM permissions p'
\echo 'JOIN role_permissions rp ON p.id = rp.permission_id'
\echo 'JOIN user_roles ur ON rp.role_id = ur.role_id'
\echo 'JOIN gebruikers u ON ur.user_id = u.id'
\echo 'WHERE u.email = ''your-email@dekoninklijkeloop.com'''
\echo 'AND ur.is_active = true'
\echo 'ORDER BY p.resource, p.action;'

-- ================================================
-- SUMMARY & RECOMMENDATIONS
-- ================================================
\echo ''
\echo '=== SUMMARY & RECOMMENDATIONS ==='
DO $$
DECLARE
    staff_role_exists BOOLEAN;
    aanmelding_read_assigned BOOLEAN;
    aanmelding_write_assigned BOOLEAN;
    users_with_staff_role INTEGER;
    users_with_legacy_staff INTEGER;
BEGIN
    -- Check staff role exists
    SELECT EXISTS(SELECT 1 FROM roles WHERE name = 'staff') INTO staff_role_exists;
    
    -- Check aanmelding permissions
    SELECT EXISTS(
        SELECT 1 FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE r.name = 'staff' AND p.resource = 'aanmelding' AND p.action = 'read'
    ) INTO aanmelding_read_assigned;
    
    SELECT EXISTS(
        SELECT 1 FROM role_permissions rp
        JOIN roles r ON rp.role_id = r.id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE r.name = 'staff' AND p.resource = 'aanmelding' AND p.action = 'write'
    ) INTO aanmelding_write_assigned;
    
    -- Count users
    SELECT COUNT(DISTINCT ur.user_id) INTO users_with_staff_role
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE r.name = 'staff' AND ur.is_active = true;
    
    SELECT COUNT(*) INTO users_with_legacy_staff
    FROM gebruikers WHERE rol = 'staff';
    
    -- Print summary
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
    RAISE NOTICE 'DIAGNOSTIC SUMMARY';
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
    
    IF staff_role_exists THEN
        RAISE NOTICE '✓ Staff role exists in roles table';
    ELSE
        RAISE NOTICE '✗ Staff role MISSING from roles table!';
    END IF;
    
    IF aanmelding_read_assigned THEN
        RAISE NOTICE '✓ Staff has aanmelding:read permission';
    ELSE
        RAISE NOTICE '✗ Staff MISSING aanmelding:read permission!';
    END IF;
    
    IF aanmelding_write_assigned THEN
        RAISE NOTICE '✓ Staff has aanmelding:write permission';
    ELSE
        RAISE NOTICE '⚠ Staff MISSING aanmelding:write permission (read-only)';
    END IF;
    
    RAISE NOTICE 'ℹ Users with RBAC staff role: %', users_with_staff_role;
    RAISE NOTICE 'ℹ Users with legacy staff role: %', users_with_legacy_staff;
    
    IF users_with_legacy_staff > users_with_staff_role THEN
        RAISE NOTICE '⚠ % users have legacy staff but no RBAC staff role!', 
            (users_with_legacy_staff - users_with_staff_role);
    END IF;
    
    RAISE NOTICE '';
    RAISE NOTICE 'RECOMMENDATIONS:';
    
    IF NOT aanmelding_read_assigned THEN
        RAISE NOTICE '1. Run hotfix script to add aanmelding permissions';
    END IF;
    
    IF users_with_legacy_staff > users_with_staff_role THEN
        RAISE NOTICE '2. Run V1_22 migration to migrate legacy staff users';
    END IF;
    
    IF NOT aanmelding_write_assigned THEN
        RAISE NOTICE '3. Staff only has READ access - run V1_50 for WRITE access';
    END IF;
    
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
END $$;