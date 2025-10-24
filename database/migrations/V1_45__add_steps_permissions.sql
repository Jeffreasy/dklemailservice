-- Add steps permissions to the RBAC system
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('steps', 'read', 'Stappen data bekijken (dashboard, totaal, verdeling)', true),
('steps', 'write', 'Stappen bijwerken voor deelnemers', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign steps permissions to admin and staff roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name IN ('admin', 'staff') AND r.is_system_role = true
  AND p.resource = 'steps'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permission to deelnemer role for their own dashboard
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'deelnemer' AND r.is_system_role = true
  AND p.resource = 'steps' AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.45.0', 'Add steps permissions for RBAC system', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;