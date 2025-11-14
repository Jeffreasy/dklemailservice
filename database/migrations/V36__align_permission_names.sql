-- V36: Corrigeer de permissienamen zodat ze overeenkomen met de applicatielogica
-- Dit is de IDEMPOTENTE versie die conflicten oplost.

BEGIN;

-- ====================================================================
-- STAP 1: Corrigeer 'steps' permissies (van V13)
-- ====================================================================

-- Verwijder eerst de foute 'read' en 'write' permissies (als ze bestaan)
-- We moeten ook de koppelingen verwijderen.
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE resource = 'steps' AND action IN ('read', 'write')
);
DELETE FROM permissions WHERE resource = 'steps' AND action IN ('read', 'write');

-- Zorg ervoor dat de JUISTE permissies ('view_own' en 'create') bestaan
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('steps', 'view_own', 'Eigen stappen en dashboard bekijken', true),
('steps', 'create', 'Eigen stappen aanmaken/bijwerken', true)
ON CONFLICT (resource, action) DO NOTHING;

-- ====================================================================
-- STAP 2: Corrigeer 'participant' permissies (van V29)
-- ====================================================================

-- Verwijder eerst de foute 'read' en 'write' permissies (als ze bestaan)
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE resource = 'participant' AND action IN ('read', 'write')
);
DELETE FROM permissions WHERE resource = 'participant' AND action IN ('read', 'write');

-- Zorg ervoor dat de JUISTE permissies ('view_own' en 'update_own') bestaan
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('participant', 'view_own', 'Eigen deelnemer-gegevens bekijken', true),
('participant', 'update_own', 'Eigen deelnemer-gegevens bijwerken', true)
ON CONFLICT (resource, action) DO NOTHING;

-- ====================================================================
-- STAP 3: Corrigeer 'leaderboard' en 'events' (van V35)
-- ====================================================================

-- Verwijder de foute 'read' permissies (als ze bestaan)
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE resource IN ('leaderboard', 'events') AND action = 'read'
);
DELETE FROM permissions WHERE resource IN ('leaderboard', 'events') AND action = 'read';

-- Zorg ervoor dat de JUISTE 'view' permissies bestaan
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('leaderboard', 'view', 'Bekijken van het leaderboard (voor participants)', true),
('events', 'view', 'Bekijken van evenementen (voor participants)', true)
ON CONFLICT (resource, action) DO NOTHING;

-- ====================================================================
-- STAP 4: Zorg dat de Admin-rol alle gecorrigeerde permissies heeft
-- ====================================================================

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
  AND p.resource IN ('steps', 'participant', 'leaderboard', 'events')
  AND p.action IN ('view_own', 'create', 'update_own', 'view')
ON CONFLICT (role_id, permission_id) DO NOTHING;

COMMIT;

-- Log de voltooiing
DO $$
BEGIN
    RAISE NOTICE '[V36] Permissienamen voor steps, participant, en leaderboard gecorrigeerd (idempotente versie).';
END $$;