-- V1_34__add_album_permissions.sql
-- Add permissions for albums management

-- Insert permissions for albums
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('album', 'read', 'Albums bekijken', true),
('album', 'write', 'Albums aanmaken/bewerken', true),
('album', 'delete', 'Albums verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign album permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'album'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'album'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;