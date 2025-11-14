-- Diagnostic SQL to check user permissions for user ID: 19ff57b7-bb7a-4658-b3c4-28df672cdee1

-- 1. Check user roles
SELECT 
    ur.user_id,
    r.name as role_name,
    r.description,
    ur.is_active,
    ur.assigned_at
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = '19ff57b7-bb7a-4658-b3c4-28df672cdee1'
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
WHERE ur.user_id = '19ff57b7-bb7a-4658-b3c4-28df672cdee1'
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
WHERE p.gebruiker_id = '19ff57b7-bb7a-4658-b3c4-28df672cdee1';