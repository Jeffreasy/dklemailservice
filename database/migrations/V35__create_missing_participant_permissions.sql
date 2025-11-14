-- V35: Maak de ontbrekende permissies aan die hardgecodeerd worden gebruikt
-- door de V34 participant authenticatie (Oplossing 2).
-- Dit zorgt ervoor dat het permissiesysteem consistent is.

INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('app', 'access', 'Toegang tot de DKL Step App (voor participants)', true),
('leaderboard', 'view', 'Bekijken van het leaderboard (voor participants)', true),
('events', 'view', 'Bekijken van evenementen (voor participants)', true),
('events', 'register', 'Registreren voor evenementen (voor participants)', true),
('profile', 'read', 'Eigen profiel bekijken (alias voor participant:read)', true),
('profile', 'update', 'Eigen profiel bijwerken (alias voor participant:write)', true)
ON CONFLICT (resource, action) DO NOTHING;

-- Wijs deze nieuwe permissies OOK toe aan de 'admin' rol,
-- zodat beheerders ze ook hebben en het systeem consistent blijft.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
  AND p.resource IN ('app', 'leaderboard', 'events', 'profile')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Log de voltooiing
DO $$
BEGIN
    RAISE NOTICE '[V35] Noodzakelijke permissies voor participant-app (app, leaderboard, events, profile) aangemaakt en toegekend aan admin.';
END $$;