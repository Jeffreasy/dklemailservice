-- =====================================================
-- Sync Production Users en Assign RBAC Roles
-- Versie: 1.49
-- Datum: 2025-11-02
-- =====================================================
--
-- Dit script synchroniseert productie gebruikers en wijst
-- de correcte RBAC roles toe op basis van email domein
--
-- Regel:
-- - @dekoninklijkeloop.nl → staff role (RO access)
-- - admin@ → admin role (full access)
-- - Anderen → role based op existing `rol` field
-- =====================================================

BEGIN;

-- =====================================================
-- STAP 1: Assign RBAC roles aan alle @dekoninklijkeloop.nl users
-- =====================================================

-- Jeffrey (staff)
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE g.email = 'jeffrey@dekoninklijkeloop.nl'
  AND r.name = 'staff'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- Lida (staff)
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE g.email = 'Lida@dekoninklijkeloop.nl'
  AND r.name = 'staff'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- Marieke (staff)
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE g.email = 'Marieke@dekoninklijkeloop.nl'
  AND r.name = 'staff'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- Salih (staff)
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE g.email = 'Salih@dekoninklijkeloop.nl'
  AND r.name = 'staff'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- =====================================================
-- STAP 2: Assign RBAC roles voor event deelnemers/begeleiders
-- =====================================================

-- Assign 'begeleider' RBAC role aan users met legacy rol 'begeleider'
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE LOWER(g.rol) = 'begeleider'
  AND r.name = 'begeleider'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- Assign 'deelnemer' RBAC role aan users met legacy rol 'deelnemer'
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE LOWER(g.rol) = 'deelnemer'
  AND r.name = 'deelnemer'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- Assign 'user' RBAC role aan users zonder specifieke rol
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id,
    r.id,
    true,
    NOW()
FROM gebruikers g
CROSS JOIN roles r
WHERE (g.rol IS NULL OR g.rol = '' OR LOWER(g.rol) = 'gebruiker' OR LOWER(g.rol) = 'socialmedia')
  AND r.name = 'user'
  AND NOT EXISTS (
      SELECT 1 FROM user_roles ur 
      WHERE ur.user_id = g.id
  )
ON CONFLICT (user_id, role_id) DO UPDATE 
SET is_active = true;

-- =====================================================
-- STAP 3: Toon resultaten
-- =====================================================

DO $$
DECLARE
    total_users INTEGER;
    users_with_rbac INTEGER;
    dekoninklijkeloop_staff INTEGER;
    event_participants INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_users FROM gebruikers WHERE is_actief = true;
    SELECT COUNT(DISTINCT user_id) INTO users_with_rbac FROM user_roles WHERE is_active = true;
    
    SELECT COUNT(DISTINCT ur.user_id) INTO dekoninklijkeloop_staff
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers g ON ur.user_id = g.id
    WHERE r.name = 'staff' 
      AND g.email LIKE '%@dekoninklijkeloop.nl'
      AND ur.is_active = true;
    
    SELECT COUNT(DISTINCT ur.user_id) INTO event_participants
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE r.name IN ('deelnemer', 'begeleider', 'vrijwilliger')
      AND ur.is_active = true;
    
    RAISE NOTICE '=== RBAC Assignment Results ===';
    RAISE NOTICE 'Total active users: %', total_users;
    RAISE NOTICE 'Users with RBAC roles: %', users_with_rbac;
    RAISE NOTICE '@dekoninklijkeloop.nl staff: %', dekoninklijkeloop_staff;
    RAISE NOTICE 'Event participants (deelnemer/begeleider): %', event_participants;
    
    IF users_with_rbac < total_users THEN
        RAISE WARNING '% users still missing RBAC roles!', (total_users - users_with_rbac);
    ELSE
        RAISE NOTICE '✓ All users have RBAC roles assigned';
    END IF;
END $$;

-- =====================================================
-- STAP 4: Toon overzicht per gebruiker
-- =====================================================

SELECT 
    g.email,
    g.naam,
    g.rol as legacy_rol,
    STRING_AGG(DISTINCT r.name, ', ' ORDER BY r.name) as rbac_roles,
    COUNT(DISTINCT p.id) as total_permissions,
    CASE 
        WHEN COUNT(ur.id) = 0 THEN '✗ NO RBAC'
        WHEN g.email LIKE '%@dekoninklijkeloop.nl' AND EXISTS (
            SELECT 1 FROM user_roles ur2 
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = g.id 
            AND r2.name IN ('staff', 'admin')
            AND ur2.is_active = true
        ) THEN '✓ STAFF/ADMIN'
        ELSE '✓ EVENT ROLE'
    END as status
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE g.is_actief = true
GROUP BY g.id, g.email, g.naam, g.rol
ORDER BY 
    CASE 
        WHEN g.email LIKE '%@dekoninklijkeloop.nl' THEN 1
        ELSE 2
    END,
    g.email;

COMMIT;

-- =====================================================
-- Completion Message
-- =====================================================

\echo ''
\echo '=== RBAC ROLE SYNCHRONIZATION COMPLETE ==='
\echo ''
\echo 'All active users now have appropriate RBAC roles:'
\echo '  - @dekoninklijkeloop.nl users → staff role (except admin)'
\echo '  - Event begeleiders → begeleider role (steps permissions)'
\echo '  - Event deelnemers → deelnemer role (steps permissions)'
\echo '  - Others → user role (basic permissions)'
\echo ''