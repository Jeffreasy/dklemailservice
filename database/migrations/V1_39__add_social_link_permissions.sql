-- V1_39__add_social_link_permissions.sql
-- Add permissions for social links management

-- Insert permissions for social_link
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('social_link', 'read', 'Social links bekijken', true),
('social_link', 'write', 'Social links aanmaken/bewerken', true),
('social_link', 'delete', 'Social links verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign social_link permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'social_link'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'social_link'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;