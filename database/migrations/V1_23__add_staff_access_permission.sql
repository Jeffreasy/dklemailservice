-- Migratie: V1_23__add_staff_access_permission.sql
-- Beschrijving: Add staff access permission and assign to roles
-- Versie: 1.23.0

-- Add staff access permission if it doesn't exist
INSERT INTO permissions (resource, action, description, is_system_permission)
VALUES ('staff', 'access', 'Toegang tot staff functies', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign staff access permission to admin and staff roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name IN ('admin', 'staff') AND r.is_system_role = true
  AND p.resource = 'staff' AND p.action = 'access'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.23.0', 'Add staff access permission and assign to admin/staff roles', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;