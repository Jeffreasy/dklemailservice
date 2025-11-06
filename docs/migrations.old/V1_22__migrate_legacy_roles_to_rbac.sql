-- Migratie: V1_22__migrate_legacy_roles_to_rbac.sql
-- Beschrijving: Migreer bestaande legacy roles (gebruikers.rol) naar RBAC user_roles tabel
-- Versie: 1.22.0
-- Datum: 2025-11-01

-- Stap 1: Voeg user_roles toe voor alle bestaande gebruikers met legacy roles
-- Deze migratie behoudt de bestaande gebruikers.rol velden voor backward compatibility
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    COALESCE(g.created_at, NOW()) as assigned_at
FROM gebruikers g
JOIN roles r ON LOWER(r.name) = LOWER(g.rol)
WHERE g.rol IS NOT NULL 
  AND g.rol != ''
  -- Voorkom duplicaten
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Stap 2: Voeg standaard 'user' role toe voor gebruikers zonder specifieke rol
INSERT INTO user_roles (user_id, role_id, is_active, assigned_at)
SELECT 
    g.id as user_id,
    r.id as role_id,
    true as is_active,
    COALESCE(g.created_at, NOW()) as assigned_at
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'user'
  AND r.is_system_role = true
  AND (g.rol IS NULL OR g.rol = '' OR g.rol = 'gebruiker')
  -- Voorkom duplicaten
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Stap 3: Log de migratie resultaten
DO $$
DECLARE
    migrated_count INTEGER;
    total_users INTEGER;
    users_without_rbac INTEGER;
BEGIN
    -- Tel totaal aantal gebruikers
    SELECT COUNT(*) INTO total_users FROM gebruikers;
    
    -- Tel gebruikers met RBAC roles
    SELECT COUNT(DISTINCT user_id) INTO migrated_count FROM user_roles WHERE is_active = true;
    
    -- Tel gebruikers zonder RBAC roles
    users_without_rbac := total_users - migrated_count;
    
    -- Log resultaten
    RAISE NOTICE 'Legacy to RBAC Migration Results:';
    RAISE NOTICE '  Total users: %', total_users;
    RAISE NOTICE '  Users with RBAC roles: %', migrated_count;
    RAISE NOTICE '  Users without RBAC roles: %', users_without_rbac;
    
    -- Waarschuwing als er gebruikers zonder RBAC rollen zijn
    IF users_without_rbac > 0 THEN
        RAISE WARNING 'Er zijn % gebruikers zonder RBAC rollen! Controleer de rol mapping.', users_without_rbac;
    END IF;
END $$;

-- Stap 4: Maak een view voor gemakkelijke controle van de migratie
CREATE OR REPLACE VIEW v_user_role_migration_status AS
SELECT 
    g.id as user_id,
    g.email,
    g.naam,
    g.rol as legacy_role,
    COALESCE(
        STRING_AGG(r.name, ', ' ORDER BY r.name), 
        'GEEN RBAC ROL'
    ) as rbac_roles,
    COUNT(ur.id) as rbac_role_count,
    CASE 
        WHEN COUNT(ur.id) = 0 THEN 'MISSING RBAC'
        WHEN g.rol IS NULL OR g.rol = '' THEN 'NO LEGACY'
        WHEN EXISTS (
            SELECT 1 FROM user_roles ur2
            JOIN roles r2 ON ur2.role_id = r2.id
            WHERE ur2.user_id = g.id 
            AND LOWER(r2.name) = LOWER(g.rol)
            AND ur2.is_active = true
        ) THEN 'MIGRATED'
        ELSE 'MISMATCH'
    END as migration_status
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
GROUP BY g.id, g.email, g.naam, g.rol
ORDER BY 
    CASE 
        WHEN COUNT(ur.id) = 0 THEN 1
        ELSE 2
    END,
    g.email;

-- Stap 5: Toon migratie status overzicht
DO $$
DECLARE
    status_record RECORD;
BEGIN
    RAISE NOTICE '=== Migration Status per Category ===';
    FOR status_record IN 
        SELECT 
            migration_status,
            COUNT(*) as count
        FROM v_user_role_migration_status
        GROUP BY migration_status
        ORDER BY 
            CASE migration_status
                WHEN 'MIGRATED' THEN 1
                WHEN 'NO LEGACY' THEN 2
                WHEN 'MISMATCH' THEN 3
                WHEN 'MISSING RBAC' THEN 4
            END
    LOOP
        RAISE NOTICE '  %: % users', status_record.migration_status, status_record.count;
    END LOOP;
END $$;

-- Stap 6: Toon problematische gevallen
DO $$
DECLARE
    problem_record RECORD;
    problem_count INTEGER := 0;
BEGIN
    -- Check voor gebruikers met mismatches of missing RBAC
    FOR problem_record IN 
        SELECT user_id, email, naam, legacy_role, rbac_roles, migration_status
        FROM v_user_role_migration_status
        WHERE migration_status IN ('MISMATCH', 'MISSING RBAC')
        LIMIT 10
    LOOP
        IF problem_count = 0 THEN
            RAISE NOTICE '=== Problematische Gebruikers (max 10) ===';
        END IF;
        problem_count := problem_count + 1;
        RAISE NOTICE '  % - % (%) | Legacy: % | RBAC: % | Status: %', 
            problem_count,
            problem_record.email,
            problem_record.naam,
            COALESCE(problem_record.legacy_role, 'NULL'),
            problem_record.rbac_roles,
            problem_record.migration_status;
    END LOOP;
    
    IF problem_count = 0 THEN
        RAISE NOTICE '✓ Geen problematische gebruikers gevonden!';
    END IF;
END $$;

-- Stap 7: Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.22.0', 'Migrate legacy roles to RBAC user_roles', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;

-- Stap 8: Toon instructies voor vervolgstappen
DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '=== NEXT STEPS ===';
    RAISE NOTICE '1. Review migration status: SELECT * FROM v_user_role_migration_status;';
    RAISE NOTICE '2. Fix any mismatches or missing RBAC roles manually';
    RAISE NOTICE '3. Test RBAC permissions thoroughly';
    RAISE NOTICE '4. Update JWT generation to use RBAC roles (code change required)';
    RAISE NOTICE '5. After thorough testing, consider deprecating gebruikers.rol field';
    RAISE NOTICE '';
    RAISE NOTICE '⚠️  BELANGRIJK: gebruikers.rol field is NIET verwijderd voor backward compatibility';
    RAISE NOTICE '⚠️  Beide systemen werken nu naast elkaar';
END $$;