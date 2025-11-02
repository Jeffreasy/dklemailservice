-- HOTFIX: Add aanmelding permissions immediately to production database
-- Run this script directly on production database to fix staff 403 errors

-- Step 1: Add aanmelding permissions if they don't exist
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken', true),
('aanmelding', 'delete', 'Aanmeldingen verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Step 2: Assign all aanmelding permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'aanmelding'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Step 3: Assign read and write permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'aanmelding'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Step 4: Verify the permissions were added
SELECT 
    r.name as role_name,
    p.resource,
    p.action,
    p.description
FROM role_permissions rp
JOIN roles r ON rp.role_id = r.id
JOIN permissions p ON rp.permission_id = p.id
WHERE p.resource = 'aanmelding'
ORDER BY r.name, p.action;

-- Step 5: Check which users have the staff role
SELECT 
    u.email,
    u.naam,
    r.name as role_name,
    ur.is_active,
    ur.assigned_at
FROM user_roles ur
JOIN gebruikers u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE r.name IN ('staff', 'admin')
  AND ur.is_active = true
ORDER BY r.name, u.email;