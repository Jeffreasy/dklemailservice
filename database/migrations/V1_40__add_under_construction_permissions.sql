-- V1_40__add_under_construction_permissions.sql
-- Add permissions for under construction management

-- Insert permissions for under_construction
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('under_construction', 'read', 'Under construction bekijken', true),
('under_construction', 'write', 'Under construction aanmaken/bewerken', true),
('under_construction', 'delete', 'Under construction verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign under_construction permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'under_construction'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'under_construction'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;