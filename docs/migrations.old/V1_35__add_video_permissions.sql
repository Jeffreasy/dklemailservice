-- V1_35__add_video_permissions.sql
-- Add permissions for videos management

-- Insert permissions for videos
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('video', 'read', 'Videos bekijken', true),
('video', 'write', 'Videos aanmaken/bewerken', true),
('video', 'delete', 'Videos verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign video permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'video'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'video'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;