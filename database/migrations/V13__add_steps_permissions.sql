-- GECONSOLIDEERDE V13 - STEPS PERMISSIES
-- Logica van het tweede V1_34 script.
-- De 'INSERT INTO migraties' is verwijderd.

-- ========================================
-- STAP 1: Maak steps permissions aan
-- ========================================
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES
('steps', 'read', 'Eigen stappen en dashboard bekijken', true),
('steps', 'write', 'Eigen stappen bijwerken', true),
('steps', 'read_total', 'Totaal aantal stappen van alle deelnemers bekijken', true),
('steps', 'read_all', 'Alle deelnemers stappen bekijken (admin/staff)', true),
('steps', 'write_all', 'Alle deelnemers stappen bijwerken (admin/staff)', true),
('steps', 'manage', 'Volledige steps beheer (route funds, etc.)', true)
ON CONFLICT (resource, action) DO NOTHING;

-- ========================================
-- STAP 2: Wijs permissions toe aan deelnemer rol
-- ========================================
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'deelnemer'
AND p.resource = 'steps'
AND p.action IN ('read', 'write', 'read_total')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 3: Wijs permissions toe aan begeleider rol
-- ========================================
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'begeleider'
AND p.resource = 'steps'
AND p.action IN ('read', 'write', 'read_total')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 4: Wijs permissions toe aan vrijwilliger rol
-- ========================================
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'vrijwilliger'
AND p.resource = 'steps'
AND p.action IN ('read', 'write', 'read_total')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 5: Wijs permissions toe aan staff rol
-- ========================================
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
INSERT INTO role_permissions (role_id, permission_id, assigned_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'admin'
AND p.resource = 'steps'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ========================================
-- STAP 7: Verificatie (Logging)
-- ========================================
DO $$
DECLARE
    status_record RECORD;
BEGIN
    RAISE NOTICE '=== Steps Migration Status per Category ===';
    FOR status_record IN 
        SELECT 
            r.name as rol,
            COUNT(p.id) as aantal_steps_permissions
        FROM roles r
        LEFT JOIN role_permissions rp ON rp.role_id = r.id
        LEFT JOIN permissions p ON p.id = rp.permission_id AND p.resource = 'steps'
        WHERE r.name IN ('deelnemer', 'begeleider', 'vrijwilliger', 'staff', 'admin')
        GROUP BY r.id, r.name
        ORDER BY r.name
    LOOP
        RAISE NOTICE '  %: % steps permissions', status_record.rol, status_record.aantal_steps_permissions;
    END LOOP;
END $$;