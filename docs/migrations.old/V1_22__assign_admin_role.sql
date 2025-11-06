-- Migratie: V1_22__assign_admin_role.sql
-- Beschrijving: Assign admin role to existing admin user
-- Versie: 1.22.0

-- Assign admin role to the admin user
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT u.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers u
CROSS JOIN roles r
WHERE u.email = 'admin@dekoninklijkeloop.nl'
  AND r.name = 'admin'
  AND r.is_system_role = true
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur
    WHERE ur.user_id = u.id AND ur.role_id = r.id AND ur.is_active = true
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.22.0', 'Assign admin role to existing admin user', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;