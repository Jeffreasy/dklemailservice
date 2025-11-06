-- Migratie: V1_49__assign_staff_roles_to_dekoninklijkeloop_users.sql
-- Beschrijving: Assign staff role to all @dekoninklijkeloop.nl users (except admin)
-- Versie: 1.49.0
-- Datum: 2025-11-02

-- =====================================================
-- STAP 1: Assign staff role to all @dekoninklijkeloop.nl users
-- =====================================================

INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    NOW() as assigned_at
FROM gebruikers g
CROSS JOIN roles r
WHERE 
    -- Alleen @dekoninklijkeloop.nl email adressen
    g.email LIKE '%@dekoninklijkeloop.nl'
    -- Staff role
    AND r.name = 'staff' 
    AND r.is_system_role = true
    -- Alleen actieve gebruikers
    AND g.is_actief = true
    -- Voorkom duplicaten - gebruiker heeft nog geen staff role
    AND NOT EXISTS (
        SELECT 1 FROM user_roles ur 
        WHERE ur.user_id = g.id 
          AND ur.role_id = r.id
    )
    -- Exclusief: gebruikers die al admin role hebben krijgen geen staff
    -- (admin heeft al alle staff permissies en meer)
    AND NOT EXISTS (
        SELECT 1 FROM user_roles ur2
        JOIN roles r2 ON ur2.role_id = r2.id
        WHERE ur2.user_id = g.id
          AND r2.name = 'admin'
          AND ur2.is_active = true
    )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- =====================================================
-- STAP 2: Update legacy rol field naar 'staff' voor consistency
-- =====================================================

-- Update gebruikers.rol naar 'staff' als ze staff role hebben via RBAC
-- maar ALLEEN als hun huidige rol 'user' of 'gebruiker' is
UPDATE gebruikers g
SET rol = 'staff'
WHERE 
    g.email LIKE '%@dekoninklijkeloop.nl'
    AND g.is_actief = true
    AND g.rol IN ('user', 'gebruiker', '')
    AND EXISTS (
        SELECT 1 FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        WHERE ur.user_id = g.id
          AND r.name = 'staff'
          AND ur.is_active = true
    );

-- =====================================================
-- STAP 3: Log resultaten
-- =====================================================

DO $$
DECLARE
    total_dekoninklijkeloop INTEGER;
    staff_assigned INTEGER;
    admin_users INTEGER;
BEGIN
    -- Tel @dekoninklijkeloop.nl gebruikers
    SELECT COUNT(*) INTO total_dekoninklijkeloop
    FROM gebruikers
    WHERE email LIKE '%@dekoninklijkeloop.nl' AND is_actief = true;
    
    -- Tel gebruikers met staff role
    SELECT COUNT(DISTINCT ur.user_id) INTO staff_assigned
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers g ON ur.user_id = g.id
    WHERE r.name = 'staff' 
      AND ur.is_active = true
      AND g.email LIKE '%@dekoninklijkeloop.nl';
    
    -- Tel gebruikers met admin role
    SELECT COUNT(DISTINCT ur.user_id) INTO admin_users
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN gebruikers g ON ur.user_id = g.id
    WHERE r.name = 'admin' 
      AND ur.is_active = true
      AND g.email LIKE '%@dekoninklijkeloop.nl';
    
    RAISE NOTICE '=== Staff Role Assignment Results ===';
    RAISE NOTICE 'Total @dekoninklijkeloop.nl users: %', total_dekoninklijkeloop;
    RAISE NOTICE 'Users with admin role: %', admin_users;
    RAISE NOTICE 'Users with staff role: %', staff_assigned;
    RAISE NOTICE 'Users without role: %', (total_dekoninklijkeloop - admin_users - staff_assigned);
    
    IF (total_dekoninklijkeloop - admin_users - staff_assigned) > 0 THEN
        RAISE WARNING '% users from @dekoninklijkeloop.nl do not have staff or admin role!', 
            (total_dekoninklijkeloop - admin_users - staff_assigned);
    ELSE
        RAISE NOTICE '✓ All @dekoninklijkeloop.nl users have appropriate roles';
    END IF;
END $$;

-- =====================================================
-- STAP 4: Toon overzicht van alle @dekoninklijkeloop.nl users
-- =====================================================

DO $$
DECLARE
    user_record RECORD;
    counter INTEGER := 0;
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '=== @dekoninklijkeloop.nl Users Overview ===';
    
    FOR user_record IN 
        SELECT 
            g.email,
            g.naam,
            STRING_AGG(DISTINCT r.name, ', ' ORDER BY r.name) as roles,
            COUNT(DISTINCT p.id) as permission_count
        FROM gebruikers g
        LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
        LEFT JOIN roles r ON ur.role_id = r.id
        LEFT JOIN role_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id
        WHERE g.email LIKE '%@dekoninklijkeloop.nl' AND g.is_actief = true
        GROUP BY g.email, g.naam
        ORDER BY g.email
    LOOP
        counter := counter + 1;
        RAISE NOTICE '  % - % (%): Roles: %, Permissions: %', 
            counter,
            user_record.email,
            user_record.naam,
            COALESCE(user_record.roles, 'NONE'),
            user_record.permission_count;
    END LOOP;
    
    IF counter = 0 THEN
        RAISE NOTICE '  No active @dekoninklijkeloop.nl users found';
    END IF;
END $$;

-- =====================================================
-- STAP 5: Registreer de migratie
-- =====================================================

INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.49.0', 'Assign staff roles to @dekoninklijkeloop.nl users', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;

-- =====================================================
-- STAP 6: Instructies voor toekomstige users
-- =====================================================

DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '=== Future User Creation ===';
    RAISE NOTICE 'For new @dekoninklijkeloop.nl users, use the API:';
    RAISE NOTICE 'POST /api/users/:id/roles';
    RAISE NOTICE 'Body: {"role_id": "<staff_role_id>"}';
    RAISE NOTICE '';
    RAISE NOTICE 'Or manually:';
    RAISE NOTICE 'INSERT INTO user_roles (user_id, role_id, is_active)';
    RAISE NOTICE 'SELECT user_id, (SELECT id FROM roles WHERE name = ''staff''), true';
    RAISE NOTICE 'FROM gebruikers WHERE email = ''new@dekoninklijkeloop.nl'';';
    RAISE NOTICE '';
END $$;