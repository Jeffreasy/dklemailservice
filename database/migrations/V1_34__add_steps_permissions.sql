-- Migratie: V1_34__add_steps_permissions.sql
-- Beschrijving: Voeg steps permissions toe en wijs ze toe aan deelnemer/begeleider/vrijwilliger rollen
-- Versie: 1.34.0

-- ========================================
-- STAP 1: Maak steps permissions aan
-- ========================================

INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('steps', 'read', 'Eigen stappen en dashboard bekijken', true),
('steps', 'write', 'Eigen stappen bijwerken', true),
('steps', 'read_all', 'Alle deelnemers stappen bekijken (admin/staff)', true),
('steps', 'write_all', 'Alle deelnemers stappen bijwerken (admin/staff)', true),
('steps', 'manage', 'Volledige steps beheer (route funds, etc.)', true)
ON CONFLICT (resource, action) DO NOTHING;

-- ========================================
-- STAP 2: Wijs read/write permissions toe aan deelnemer rol
-- ========================================

-- Deelnemers kunnen hun eigen stappen lezen en schrijven
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'deelnemer'
AND p.resource = 'steps'
AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 3: Wijs read/write permissions toe aan begeleider rol
-- ========================================

-- Begeleiders kunnen hun eigen stappen lezen en schrijven
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'begeleider'
AND p.resource = 'steps'
AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 4: Wijs read/write permissions toe aan vrijwilliger rol
-- ========================================

-- Vrijwilligers kunnen hun eigen stappen lezen en schrijven
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'vrijwilliger'
AND p.resource = 'steps'
AND p.action IN ('read', 'write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 5: Wijs read_all/write_all permissions toe aan staff rol
-- ========================================

-- Staff kan alle deelnemers stappen lezen en schrijven
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'staff'
AND p.resource = 'steps'
AND p.action IN ('read', 'write', 'read_all', 'write_all')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 6: Admin krijgt alle steps permissions
-- ========================================

-- Admin krijgt automatisch alle permissions (al geconfigureerd in V1_21)
-- Maar voor zekerheid expliciet toewijzen:
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'admin'
AND p.resource = 'steps'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 7: Verificatie
-- ========================================

-- Toon alle steps permissions
SELECT 
    r.name as rol,
    p.resource,
    p.action,
    p.description,
    rp.assigned_at
FROM role_permissions rp
JOIN roles r ON r.id = rp.role_id
JOIN permissions p ON p.id = rp.permission_id
WHERE p.resource = 'steps'
ORDER BY r.name, p.action;

-- Toon overzicht per rol
SELECT 
    r.name as rol,
    COUNT(p.id) as aantal_steps_permissions
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id AND p.resource = 'steps'
WHERE r.name IN ('deelnemer', 'begeleider', 'vrijwilliger', 'staff', 'admin')
GROUP BY r.id, r.name
ORDER BY r.name;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.34.0', 'Add steps permissions for deelnemer/begeleider/vrijwilliger roles', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;