-- =====================================================
-- RBAC Database Verification Script
-- Version: 1.49
-- Datum: 2025-11-02
-- =====================================================
-- 
-- Dit script controleert de integriteit van de RBAC tables
-- in zowel productie als docker environments
--
-- Gebruik: 
--   psql -h localhost -U user -d dbname -f verify_rbac_tables.sql
--   OF via docker:
--   docker exec -i postgres_container psql -U user dbname < verify_rbac_tables.sql
-- =====================================================

\echo '=== RBAC DATABASE VERIFICATION ==='
\echo ''

-- =====================================================
-- 1. TABLE EXISTENCE CHECK
-- =====================================================
\echo '1. CHECKING TABLE EXISTENCE...'
\echo ''

SELECT 
    table_name,
    CASE 
        WHEN table_name IN (
            SELECT tablename FROM pg_tables 
            WHERE schemaname = 'public'
        ) THEN '✓ EXISTS'
        ELSE '✗ MISSING'
    END as status
FROM (VALUES 
    ('roles'),
    ('permissions'),
    ('role_permissions'),
    ('user_roles'),
    ('refresh_tokens'),
    ('gebruikers')
) AS required_tables(table_name)
ORDER BY table_name;

\echo ''

-- =====================================================
-- 2. TABLE STRUCTURE VERIFICATION
-- =====================================================
\echo '2. VERIFYING TABLE STRUCTURES...'
\echo ''

\echo '--- ROLES TABLE ---'
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_schema = 'public' 
  AND table_name = 'roles'
ORDER BY ordinal_position;

\echo ''
\echo '--- PERMISSIONS TABLE ---'
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_schema = 'public' 
  AND table_name = 'permissions'
ORDER BY ordinal_position;

\echo ''
\echo '--- USER_ROLES TABLE ---'
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_schema = 'public' 
  AND table_name = 'user_roles'
ORDER BY ordinal_position;

\echo ''
\echo '--- REFRESH_TOKENS TABLE ---'
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_schema = 'public' 
  AND table_name = 'refresh_tokens'
ORDER BY ordinal_position;

\echo ''

-- =====================================================
-- 3. INDEX VERIFICATION
-- =====================================================
\echo '3. VERIFYING INDEXES...'
\echo ''

SELECT 
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename IN ('roles', 'permissions', 'role_permissions', 'user_roles', 'refresh_tokens')
ORDER BY tablename, indexname;

\echo ''

-- =====================================================
-- 4. DATA INTEGRITY CHECK
-- =====================================================
\echo '4. CHECKING DATA INTEGRITY...'
\echo ''

-- Count records per table
\echo '--- RECORD COUNTS ---'
SELECT 
    'roles' as table_name,
    COUNT(*) as total_records,
    COUNT(*) FILTER (WHERE is_system_role = true) as system_roles,
    COUNT(*) FILTER (WHERE is_system_role = false) as custom_roles
FROM roles
UNION ALL
SELECT 
    'permissions',
    COUNT(*),
    COUNT(*) FILTER (WHERE is_system_permission = true),
    COUNT(*) FILTER (WHERE is_system_permission = false)
FROM permissions
UNION ALL
SELECT 
    'role_permissions',
    COUNT(*),
    COUNT(DISTINCT role_id),
    COUNT(DISTINCT permission_id)
FROM role_permissions
UNION ALL
SELECT 
    'user_roles',
    COUNT(*),
    COUNT(*) FILTER (WHERE is_active = true),
    COUNT(*) FILTER (WHERE is_active = false)
FROM user_roles
UNION ALL
SELECT 
    'refresh_tokens',
    COUNT(*),
    COUNT(*) FILTER (WHERE is_revoked = false AND expires_at > NOW()),
    COUNT(*) FILTER (WHERE is_revoked = true OR expires_at <= NOW())
FROM refresh_tokens;

\echo ''

-- =====================================================
-- 5. SYSTEM ROLES VERIFICATION
-- =====================================================
\echo '5. VERIFYING SYSTEM ROLES...'
\echo ''

WITH expected_roles AS (
    SELECT unnest(ARRAY[
        'admin', 'staff', 'user', 
        'owner', 'chat_admin', 'member',
        'deelnemer', 'begeleider', 'vrijwilliger'
    ]) as expected_name
)
SELECT 
    er.expected_name,
    CASE 
        WHEN r.id IS NOT NULL THEN '✓ EXISTS'
        ELSE '✗ MISSING'
    END as status,
    r.description,
    r.is_system_role
FROM expected_roles er
LEFT JOIN roles r ON r.name = er.expected_name
ORDER BY 
    CASE 
        WHEN r.id IS NULL THEN 1 
        ELSE 2 
    END,
    er.expected_name;

\echo ''

-- =====================================================
-- 6. PERMISSIONS PER ROLE
-- =====================================================
\echo '6. PERMISSIONS PER ROLE...'
\echo ''

SELECT 
    r.name as role_name,
    r.is_system_role,
    COUNT(p.id) as permission_count,
    STRING_AGG(p.resource || ':' || p.action, ', ' ORDER BY p.resource, p.action) as permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
GROUP BY r.id, r.name, r.is_system_role
ORDER BY permission_count DESC, r.name;

\echo ''

-- =====================================================
-- 7. USER ROLES DISTRIBUTION
-- =====================================================
\echo '7. USER ROLES DISTRIBUTION...'
\echo ''

SELECT 
    r.name as role_name,
    COUNT(DISTINCT ur.user_id) as user_count,
    COUNT(DISTINCT ur.user_id) FILTER (WHERE ur.is_active = true) as active_users,
    COUNT(DISTINCT ur.user_id) FILTER (WHERE ur.expires_at IS NOT NULL AND ur.expires_at > NOW()) as with_expiry
FROM roles r
LEFT JOIN user_roles ur ON r.id = ur.role_id
GROUP BY r.id, r.name
ORDER BY user_count DESC, r.name;

\echo ''

-- =====================================================
-- 8. LEGACY VS RBAC COMPARISON
-- =====================================================
\echo '8. LEGACY VS RBAC ROLE COMPARISON...'
\echo ''

SELECT 
    g.email,
    g.rol as legacy_role,
    STRING_AGG(r.name, ', ' ORDER BY r.name) as rbac_roles,
    COUNT(ur.id) as rbac_role_count,
    CASE 
        WHEN COUNT(ur.id) = 0 THEN '✗ NO RBAC ROLES'
        WHEN g.rol IS NULL OR g.rol = '' THEN '✓ RBAC ONLY'
        WHEN EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = g.id 
            AND LOWER(r2.name) = LOWER(g.rol)
            AND ur2.is_active = true
        ) THEN '✓ MIGRATED'
        ELSE '⚠ MISMATCH'
    END as status
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
GROUP BY g.id, g.email, g.rol
ORDER BY 
    CASE 
        WHEN COUNT(ur.id) = 0 THEN 1
        WHEN NOT EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = g.id 
            AND LOWER(r2.name) = LOWER(COALESCE(g.rol, ''))
            AND ur2.is_active = true
        ) THEN 2
        ELSE 3
    END,
    g.email;

\echo ''

-- =====================================================
-- 9. ORPHANED RECORDS CHECK
-- =====================================================
\echo '9. CHECKING FOR ORPHANED RECORDS...'
\echo ''

-- Role permissions without valid role
SELECT 
    'Orphaned Role Permissions' as check_type,
    COUNT(*) as count
FROM role_permissions rp
LEFT JOIN roles r ON rp.role_id = r.id
WHERE r.id IS NULL

UNION ALL

-- Role permissions without valid permission
SELECT 
    'Orphaned Permission References',
    COUNT(*)
FROM role_permissions rp
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE p.id IS NULL

UNION ALL

-- User roles without valid user
SELECT 
    'Orphaned User Roles (User)',
    COUNT(*)
FROM user_roles ur
LEFT JOIN gebruikers g ON ur.user_id = g.id
WHERE g.id IS NULL

UNION ALL

-- User roles without valid role
SELECT 
    'Orphaned User Roles (Role)',
    COUNT(*)
FROM user_roles ur
LEFT JOIN roles r ON ur.role_id = r.id
WHERE r.id IS NULL

UNION ALL

-- Refresh tokens without valid user
SELECT 
    'Orphaned Refresh Tokens',
    COUNT(*)
FROM refresh_tokens rt
LEFT JOIN gebruikers g ON rt.user_id = g.id
WHERE g.id IS NULL;

\echo ''

-- =====================================================
-- 10. REFRESH TOKENS STATUS
-- =====================================================
\echo '10. REFRESH TOKENS STATUS...'
\echo ''

SELECT 
    COUNT(*) as total_tokens,
    COUNT(*) FILTER (WHERE is_revoked = false) as active_tokens,
    COUNT(*) FILTER (WHERE is_revoked = true) as revoked_tokens,
    COUNT(*) FILTER (WHERE expires_at > NOW() AND is_revoked = false) as valid_tokens,
    COUNT(*) FILTER (WHERE expires_at <= NOW()) as expired_tokens,
    COUNT(DISTINCT user_id) as users_with_tokens
FROM refresh_tokens;

\echo ''

-- =====================================================
-- 11. PERMISSION COVERAGE CHECK
-- =====================================================
\echo '11. CHECKING PERMISSION COVERAGE...'
\echo ''

-- Resources with permissions
SELECT 
    resource,
    COUNT(*) as action_count,
    STRING_AGG(action, ', ' ORDER BY action) as actions
FROM permissions
GROUP BY resource
ORDER BY resource;

\echo ''

-- =====================================================
-- 12. PROBLEMATIC CASES
-- =====================================================
\echo '12. IDENTIFYING PROBLEMATIC CASES...'
\echo ''

-- Users without any roles
\echo '--- Users WITHOUT any RBAC roles ---'
SELECT 
    g.id,
    g.email,
    g.naam,
    g.rol as legacy_role,
    g.is_actief
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
WHERE ur.id IS NULL
ORDER BY g.email
LIMIT 10;

\echo ''

-- Users with expired roles
\echo '--- Users with EXPIRED roles ---'
SELECT 
    g.email,
    r.name as role_name,
    ur.expires_at,
    ur.assigned_at,
    EXTRACT(DAY FROM (NOW() - ur.expires_at)) as days_expired
FROM user_roles ur
JOIN gebruikers g ON ur.user_id = g.id
JOIN roles r ON ur.role_id = r.id
WHERE ur.expires_at IS NOT NULL 
  AND ur.expires_at <= NOW()
  AND ur.is_active = true
ORDER BY ur.expires_at DESC
LIMIT 10;

\echo ''

-- Users with inactive roles
\echo '--- Users with INACTIVE roles ---'
SELECT 
    g.email,
    r.name as role_name,
    ur.assigned_at,
    ur.is_active
FROM user_roles ur
JOIN gebruikers g ON ur.user_id = g.id
JOIN roles r ON ur.role_id = r.id
WHERE ur.is_active = false
ORDER BY ur.assigned_at DESC
LIMIT 10;

\echo ''

-- =====================================================
-- 13. SUMMARY STATISTICS
-- =====================================================
\echo '13. SUMMARY STATISTICS...'
\echo ''

WITH stats AS (
    SELECT 
        (SELECT COUNT(*) FROM roles) as total_roles,
        (SELECT COUNT(*) FROM roles WHERE is_system_role = true) as system_roles,
        (SELECT COUNT(*) FROM permissions) as total_permissions,
        (SELECT COUNT(*) FROM permissions WHERE is_system_permission = true) as system_permissions,
        (SELECT COUNT(DISTINCT user_id) FROM user_roles WHERE is_active = true) as users_with_roles,
        (SELECT COUNT(*) FROM gebruikers) as total_users,
        (SELECT COUNT(*) FROM gebruikers WHERE is_actief = true) as active_users,
        (SELECT COUNT(*) FROM refresh_tokens WHERE is_revoked = false AND expires_at > NOW()) as valid_refresh_tokens
)
SELECT 
    '9 System Roles Expected' as check_item,
    system_roles as actual,
    CASE WHEN system_roles = 9 THEN '✓ PASS' ELSE '✗ FAIL' END as status
FROM stats
UNION ALL
SELECT 
    'At least 25 Permissions Expected',
    system_permissions,
    CASE WHEN system_permissions >= 25 THEN '✓ PASS' ELSE '✗ FAIL' END
FROM stats
UNION ALL
SELECT 
    'All Active Users Have RBAC Roles',
    users_with_roles || ' of ' || active_users,
    CASE WHEN users_with_roles >= active_users THEN '✓ PASS' ELSE '⚠ WARNING' END
FROM stats
UNION ALL
SELECT 
    'Roles Table Populated',
    total_roles,
    CASE WHEN total_roles > 0 THEN '✓ PASS' ELSE '✗ FAIL' END
FROM stats
UNION ALL
SELECT 
    'Permissions Table Populated',
    total_permissions,
    CASE WHEN total_permissions > 0 THEN '✓ PASS' ELSE '✗ FAIL' END
FROM stats
UNION ALL
SELECT 
    'Valid Refresh Tokens',
    valid_refresh_tokens,
    '✓ INFO'
FROM stats;

\echo ''

-- =====================================================
-- 14. ADMIN USER VERIFICATION
-- =====================================================
\echo '14. ADMIN USER VERIFICATION...'
\echo ''

SELECT 
    g.email,
    g.naam,
    g.is_actief,
    g.rol as legacy_role,
    STRING_AGG(DISTINCT r.name, ', ' ORDER BY r.name) as rbac_roles,
    COUNT(DISTINCT p.id) as total_permissions
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE g.email LIKE '%admin%' OR g.rol = 'admin' OR r.name = 'admin'
GROUP BY g.id, g.email, g.naam, g.is_actief, g.rol
ORDER BY g.email;

\echo ''

-- =====================================================
-- 15. CONSTRAINT VERIFICATION
-- =====================================================
\echo '15. VERIFYING FOREIGN KEY CONSTRAINTS...'
\echo ''

SELECT 
    tc.table_name,
    tc.constraint_name,
    tc.constraint_type,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu 
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu 
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.table_schema = 'public'
  AND tc.table_name IN ('roles', 'permissions', 'role_permissions', 'user_roles', 'refresh_tokens')
  AND tc.constraint_type = 'FOREIGN KEY'
ORDER BY tc.table_name, tc.constraint_name;

\echo ''

-- =====================================================
-- 16. MIGRATION STATUS
-- =====================================================
\echo '16. MIGRATION STATUS...'
\echo ''

SELECT 
    versie,
    naam,
    toegepast
FROM migraties
WHERE versie IN ('1.20.0', '1.21.0', '1.22.0', '1.28.0')
  OR naam LIKE '%RBAC%'
  OR naam LIKE '%refresh%'
ORDER BY versie;

\echo ''

-- =====================================================
-- 17. POTENTIAL ISSUES
-- =====================================================
\echo '17. SCANNING FOR POTENTIAL ISSUES...'
\echo ''

-- Issue 1: Duplicate user-role assignments
WITH duplicate_roles AS (
    SELECT user_id, role_id, COUNT(*) as count
    FROM user_roles
    WHERE is_active = true
    GROUP BY user_id, role_id
    HAVING COUNT(*) > 1
)
SELECT 
    'Duplicate Active User-Role Assignments' as issue,
    COUNT(*) as occurrences,
    CASE WHEN COUNT(*) = 0 THEN '✓ NONE' ELSE '✗ FOUND' END as status
FROM duplicate_roles;

-- Issue 2: Permissions without roles
SELECT 
    'Permissions Not Assigned to Any Role',
    COUNT(*),
    CASE WHEN COUNT(*) = 0 THEN '✓ NONE' ELSE '⚠ FOUND' END
FROM permissions p
LEFT JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.id IS NULL;

-- Issue 3: Active users without permissions
SELECT 
    'Active Users Without Any Permissions',
    COUNT(*),
    CASE WHEN COUNT(*) = 0 THEN '✓ NONE' ELSE '⚠ FOUND' END
FROM gebruikers g
WHERE g.is_actief = true
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur
      JOIN role_permissions rp ON ur.role_id = rp.role_id
      WHERE ur.user_id = g.id AND ur.is_active = true
  );

-- Issue 4: System roles that can be deleted
SELECT 
    'System Roles with is_system_role=false',
    COUNT(*),
    CASE WHEN COUNT(*) = 0 THEN '✓ NONE' ELSE '✗ FOUND' END
FROM roles
WHERE name IN ('admin', 'staff', 'user', 'owner', 'chat_admin', 'member')
  AND is_system_role = false;

\echo ''

-- =====================================================
-- VERIFICATION COMPLETE
-- =====================================================
\echo '=== VERIFICATION COMPLETE ==='
\echo ''
\echo 'Review the output above for any ✗ FAIL or ⚠ WARNING statuses'
\echo 'Expected: All checks should show ✓ PASS or ✓ NONE'
\echo ''