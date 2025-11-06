-- V1_36__add_sponsor_permissions.sql
-- Add permissions for sponsors management

-- Insert permissions for sponsors
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('sponsor', 'read', 'Sponsors bekijken', true),
('sponsor', 'write', 'Sponsors aanmaken/bewerken', true),
('sponsor', 'delete', 'Sponsors verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign sponsor permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'sponsor'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'sponsor'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;