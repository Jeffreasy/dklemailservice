-- Migratie: V1_25__assign_new_permissions_to_admin.sql
-- Beschrijving: Assign new permissions for Photos, Albums, Partners, Sponsors, Videos to admin role
-- Versie: 1.25.0

-- Assign new permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
  AND p.resource IN ('photo', 'album', 'partner', 'sponsor', 'video')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.25.0', 'Assign new permissions for Photos, Albums, Partners, Sponsors, Videos to admin role', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;