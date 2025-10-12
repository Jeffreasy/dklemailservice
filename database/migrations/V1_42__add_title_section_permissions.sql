-- V1_42__add_title_section_permissions.sql
-- Add permissions for title section management

-- Insert permissions for title_section
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('title_section', 'read', 'Title section bekijken', true),
('title_section', 'write', 'Title section aanmaken/bewerken', true),
('title_section', 'delete', 'Title section verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign title_section permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'title_section'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'title_section'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;