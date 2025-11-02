-- Migratie: V1_50__add_aanmelding_permissions.sql
-- Beschrijving: Add aanmelding permissions and assign to staff and admin roles
-- Versie: 1.50.0

-- Voeg aanmelding permissions toe
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('aanmelding', 'read', 'Aanmeldingen bekijken', true),
('aanmelding', 'write', 'Aanmeldingen bewerken', true),
('aanmelding', 'delete', 'Aanmeldingen verwijderen', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Assign alle aanmelding permissions aan admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource = 'aanmelding'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read en write permissions aan staff role (geen delete)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource = 'aanmelding'
  AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.50.0', 'Add aanmelding permissions and assign to staff and admin roles', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;