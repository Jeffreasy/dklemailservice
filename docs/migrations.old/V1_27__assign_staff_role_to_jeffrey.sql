-- Migratie: V1_27__assign_staff_role_to_jeffrey.sql
-- Beschrijving: Assign staff role to jeffrey@dekoninklijkeloop.nl
-- Versie: 1.27.0

-- Assign staff role to jeffrey
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT u.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers u
CROSS JOIN roles r
WHERE u.email = 'jeffrey@dekoninklijkeloop.nl'
  AND r.name = 'staff'
  AND r.is_system_role = true
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur
    WHERE ur.user_id = u.id AND ur.role_id = r.id AND ur.is_active = true
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.27.0', 'Assign staff role to jeffrey@dekoninklijkeloop.nl', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;