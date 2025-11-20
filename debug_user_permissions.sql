-- Diagnostic SQL to check user permissions for ADMIN user: f426a6c0-7cff-46cb-a78b-7cbc6ca57831

-- 1. Check user roles
SELECT
    ur.user_id,
    r.name as role_name,
    r.description,
    ur.is_active,
    ur.assigned_at
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = 'f426a6c0-7cff-46cb-a78b-7cbc6ca57831'
    AND ur.is_active = true;

-- 2. Check all permissions for this user
SELECT
    r.name as role_name,
    p.resource,
    p.action,
    p.description
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE ur.user_id = 'f426a6c0-7cff-46cb-a78b-7cbc6ca57831'
    AND ur.is_active = true
ORDER BY p.resource, p.action;

-- 3. Check if 'steps:read' permission exists in system
SELECT id, resource, action, description, is_system_permission
FROM permissions
WHERE resource = 'steps';

-- 4. Check what 'steps' permissions participant_user role has
SELECT 
    r.name as role_name,
    p.resource,
    p.action,
    p.description
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name IN ('participant_user', 'participant_guide')
    AND p.resource = 'steps'
ORDER BY r.name, p.action;

-- 5. Check participant account details
SELECT
    p.id,
    p.email,
    p.naam,
    p.account_type,
    p.has_app_access,
    p.gebruiker_id
FROM participants p
WHERE p.gebruiker_id = 'f426a6c0-7cff-46cb-a78b-7cbc6ca57831';

-- 6. Check ALL permissions in system
SELECT resource, action, description
FROM permissions
ORDER BY resource, action;

-- 7. Check which permissions ADMIN role has
SELECT p.resource, p.action, p.description
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'admin'
ORDER BY p.resource, p.action;

-- 8. Check which permissions exist but are NOT assigned to admin
SELECT p.resource, p.action, p.description
FROM permissions p
WHERE NOT EXISTS (
    SELECT 1 FROM roles r
    JOIN role_permissions rp ON r.id = rp.role_id
    WHERE r.name = 'admin' AND rp.permission_id = p.id
)
ORDER BY p.resource, p.action;

-- 9. Check admin user details
SELECT
    g.id,
    g.naam,
    g.email,
    g.is_actief,
    g.created_at
FROM gebruikers g
WHERE g.email = 'admin@dekoninklijkeloop.nl';