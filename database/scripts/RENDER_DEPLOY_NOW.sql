-- ========================================
-- RENDER.COM PRODUCTION DATABASE HOTFIX
-- ========================================
-- Voer dit script uit via Render.com Dashboard > Database > Shell
-- Deze fix lost ALLE staff permission problemen op
-- SAFE: Idempotent, kan meerdere keren uitgevoerd worden
-- ========================================

-- STEP 1: Ensure staff role exists
INSERT INTO roles (name, description, is_system_role) VALUES
('staff', 'Ondersteunend personeel met beperkte beheerrechten', true)
ON CONFLICT (name) DO NOTHING;

-- STEP 2: Ensure aanmelding permissions exist
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken', true),
('aanmelding', 'delete', 'Aanmeldingen verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- STEP 3: Assign aanmelding permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'staff' 
  AND r.is_system_role = true
  AND p.resource = 'aanmelding' 
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- STEP 4: Assign aanmelding permissions to admin role (full access)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'admin' 
  AND r.is_system_role = true
  AND p.resource = 'aanmelding'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- STEP 5: Migrate legacy staff users to RBAC
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
ON CONFLICT (user_id, role_id) DO NOTHING;

-- STEP 6: Auto-assign staff role to @dekoninklijkeloop.com users
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
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ========================================
-- VERIFICATION QUERIES
-- ========================================

-- Verify staff role has aanmelding permissions
SELECT 
    '✓ Staff Permissions' as check_name,
    r.name as role_name,
    p.resource,
    p.action,
    'GRANTED' as status
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'staff' AND p.resource = 'aanmelding'
ORDER BY p.action;

-- Verify which users have staff role
SELECT 
    '✓ Users with Staff Role' as check_name,
    u.email,
    u.naam,
    ur.is_active,
    ur.assigned_at
FROM user_roles ur
JOIN gebruikers u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE r.name = 'staff' 
  AND ur.is_active = true
ORDER BY u.email;

-- ========================================
-- SUCCESS MESSAGE
-- ========================================
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
    RAISE NOTICE '✓✓✓ HOTFIX SCRIPT COMPLETED SUCCESSFULLY ✓✓✓';
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
    RAISE NOTICE '';
    RAISE NOTICE 'NEXT STEPS:';
    RAISE NOTICE '1. Staff users should LOGOUT and LOGIN again';
    RAISE NOTICE '2. New JWT token will include staff role and permissions';
    RAISE NOTICE '3. Staff can now access /api/aanmeldingen without 403 errors';
    RAISE NOTICE '';
    RAISE NOTICE 'Review the verification queries above to confirm setup.';
    RAISE NOTICE '';
END $$;