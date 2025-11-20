-- V45 - Add Staff Users
-- Voeg staff gebruikers toe met standaard wachtwoord DKL2026!

-- Ginelly
INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
VALUES (
    'Ginelly',
    'ginelly@dekoninklijkeloop.nl',
    '$2a$10$o8X1vELAO2ZAnifk45MKfOcgoXMy3/Fh6X/xg8GZtJRH366eBbA2a',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING;

-- Jeffrey
INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
VALUES (
    'Jeffrey',
    'jeffrey@dekoninklijkeloop.nl',
    '$2a$10$o8X1vELAO2ZAnifk45MKfOcgoXMy3/Fh6X/xg8GZtJRH366eBbA2a',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING;

-- Lida
INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
VALUES (
    'Lida',
    'lida@dekoninklijkeloop.nl',
    '$2a$10$o8X1vELAO2ZAnifk45MKfOcgoXMy3/Fh6X/xg8GZtJRH366eBbA2a',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING;

-- Marieke
INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
VALUES (
    'Marieke',
    'marieke@dekoninklijkeloop.nl',
    '$2a$10$o8X1vELAO2ZAnifk45MKfOcgoXMy3/Fh6X/xg8GZtJRH366eBbA2a',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING;

-- Salih
INSERT INTO gebruikers (naam, email, wachtwoord_hash, is_actief, created_at, updated_at)
VALUES (
    'Salih',
    'salih@dekoninklijkeloop.nl',
    '$2a$10$o8X1vELAO2ZAnifk45MKfOcgoXMy3/Fh6X/xg8GZtJRH366eBbA2a',
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING;

-- Wijs 'staff' rol toe aan alle staff gebruikers
INSERT INTO user_roles (user_id, role_id, assigned_at, is_active)
SELECT u.id, r.id, CURRENT_TIMESTAMP, true
FROM gebruikers u
CROSS JOIN roles r
WHERE u.email IN (
    'ginelly@dekoninklijkeloop.nl',
    'jeffrey@dekoninklijkeloop.nl',
    'lida@dekoninklijkeloop.nl',
    'marieke@dekoninklijkeloop.nl',
    'salih@dekoninklijkeloop.nl'
)
  AND r.name = 'staff'
  AND r.is_system_role = true
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Log de toegevoegde gebruikers
DO $$
DECLARE
    user_record RECORD;
    total_added INTEGER := 0;
BEGIN
    RAISE NOTICE '=== Staff Users Migration Results ===';

    FOR user_record IN
        SELECT g.naam, g.email, COUNT(ur.id) as roles_count
        FROM gebruikers g
        LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
        WHERE g.email IN (
            'ginelly@dekoninklijkeloop.nl',
            'jeffrey@dekoninklijkeloop.nl',
            'lida@dekoninklijkeloop.nl',
            'marieke@dekoninklijkeloop.nl',
            'salih@dekoninklijkeloop.nl'
        )
        GROUP BY g.id, g.naam, g.email
        ORDER BY g.email
    LOOP
        RAISE NOTICE '✓ % (%) - Roles: %', user_record.naam, user_record.email, user_record.roles_count;
        total_added := total_added + 1;
    END LOOP;

    RAISE NOTICE 'Total staff users processed: %', total_added;
    RAISE NOTICE 'Default password for all: DKL2026!';
END $$;