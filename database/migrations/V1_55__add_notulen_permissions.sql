-- V1_55__add_notulen_permissions.sql
-- Add RBAC permissions for notulen management

-- Insert notulen permissions (idempotent)
INSERT INTO permissions (resource, action, description) VALUES
('notulen', 'read', 'Kan notulen lezen en bekijken'),
('notulen', 'write', 'Kan notulen aanmaken en bijwerken'),
('notulen', 'delete', 'Kan notulen verwijderen'),
('notulen', 'finalize', 'Kan notulen finaliseren'),
('notulen', 'archive', 'Kan notulen archiveren')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to admin role (idempotent)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND p.resource = 'notulen'
ON CONFLICT DO NOTHING;

-- Assign read and write permissions to staff role (idempotent)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND p.resource = 'notulen' AND p.action IN ('read', 'write')
ON CONFLICT DO NOTHING;