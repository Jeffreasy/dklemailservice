-- V1_38__add_social_embed_permissions.sql
-- Add permissions for social embeds management

-- Insert permissions for social_embed
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('social_embed', 'read', 'Social embeds bekijken', true),
('social_embed', 'write', 'Social embeds aanmaken/bewerken', true),
('social_embed', 'delete', 'Social embeds verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign social_embed permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'social_embed'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'social_embed'
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;