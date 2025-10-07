-- Migratie: V1_26__assign_new_permissions_to_staff.sql
-- Beschrijving: Assign read permissions for Photos, Albums, Partners, Sponsors, Videos to staff role
-- Versie: 1.26.0

-- Assign read permissions for new resources to staff role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND r.is_system_role = true
  AND p.resource IN ('photo', 'album', 'partner', 'sponsor', 'video')
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.26.0', 'Assign read permissions for Photos, Albums, Partners, Sponsors, Videos to staff role', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;