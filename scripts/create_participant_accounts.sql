-- SQL Script om gebruikersaccounts aan te maken voor alle deelnemers met RBAC rollen
-- Wachtwoord voor alle nieuwe accounts: DKL2025!
-- Bcrypt hash (cost 10): $2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku

-- ========================================
-- STAP 1: Maak gebruikersaccounts aan met juiste rol
-- ========================================

-- Insert nieuwe gebruikers voor deelnemers die nog geen account hebben
-- Gebruikt DISTINCT ON om duplicaat emails te vermijden (meerdere aanmeldingen met zelfde email)
-- Rol wordt bepaald obv de "rol" kolom in aanmeldingen: Deelnemer/Begeleider -> deelnemer/begeleider
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, newsletter_subscribed, created_at, updated_at)
SELECT
    gen_random_uuid() as id,
    naam,
    email,
    '$2a$10$/kWEPOMqYfcy5hNYne8J5.oJQgfaBYMDE9tClXKiHlCBd/l78Dmku' as wachtwoord_hash,
    rol,
    true as is_actief,
    false as newsletter_subscribed,
    NOW() as created_at,
    NOW() as updated_at
FROM (
    -- Subquery die duplicaat emails elimineert en de meest recente aanmelding gebruikt
    SELECT DISTINCT ON (LOWER(email))
        naam,
        email,
        CASE
            WHEN LOWER(rol) = 'begeleider' THEN 'begeleider'
            WHEN LOWER(rol) = 'vrijwilliger' THEN 'vrijwilliger'
            ELSE 'deelnemer'
        END as rol,
        created_at
    FROM aanmeldingen
    WHERE email IS NOT NULL
    AND email != ''
    AND TRIM(email) != ''
    AND NOT EXISTS (
        SELECT 1 FROM gebruikers g WHERE LOWER(g.email) = LOWER(aanmeldingen.email)
    )
    ORDER BY LOWER(email), created_at DESC  -- Laatste aanmelding per email
) a;

-- ========================================
-- STAP 2: Link aanmeldingen aan gebruikers
-- ========================================

-- Update alle aanmeldingen om ze te linken aan hun gebruikersaccounts
UPDATE aanmeldingen a
SET gebruiker_id = g.id,
    updated_at = NOW()
FROM gebruikers g
WHERE LOWER(a.email) = LOWER(g.email)
AND a.gebruiker_id IS NULL
AND a.email IS NOT NULL 
AND a.email != '';

-- ========================================
-- STAP 3: Wijs RBAC rollen toe aan gebruikers
-- ========================================

-- Wijs RBAC rol toe aan alle nieuwe gebruikers op basis van hun rol
-- Dit maakt een entry in de user_roles tabel voor het RBAC systeem
INSERT INTO user_roles (id, user_id, role_id, assigned_at, is_active)
SELECT
    gen_random_uuid() as id,
    g.id as user_id,
    r.id as role_id,
    NOW() as assigned_at,
    true as is_active
FROM gebruikers g
JOIN roles r ON r.name = g.rol
WHERE g.rol IN ('deelnemer', 'begeleider', 'vrijwilliger')
AND NOT EXISTS (
    SELECT 1 FROM user_roles ur
    WHERE ur.user_id = g.id AND ur.role_id = r.id
)
ORDER BY g.created_at DESC;

-- ========================================
-- STAP 4: Verificatie
-- ========================================

-- Tel hoeveel deelnemers een account hebben
SELECT 
    'Deelnemers met account' as category,
    COUNT(*) as count
FROM gebruikers 
WHERE rol = 'deelnemer'

UNION ALL

-- Tel hoeveel aanmeldingen gelinkt zijn
SELECT 
    'Aanmeldingen gelinkt' as category,
    COUNT(*) as count
FROM aanmeldingen 
WHERE gebruiker_id IS NOT NULL

UNION ALL

-- Tel hoeveel aanmeldingen NIET gelinkt zijn
SELECT 
    'Aanmeldingen NIET gelinkt' as category,
    COUNT(*) as count
FROM aanmeldingen 
WHERE gebruiker_id IS NULL
AND email IS NOT NULL 
AND email != '';

-- ========================================
-- STAP 4: Toon nieuwe accounts
-- ========================================

-- Laat alle nieuwe deelnemer accounts zien
SELECT 
    g.id,
    g.naam,
    g.email,
    g.rol,
    g.is_actief,
    g.created_at,
    (SELECT COUNT(*) FROM aanmeldingen WHERE gebruiker_id = g.id) as aantal_aanmeldingen
FROM gebruikers g
WHERE g.rol = 'deelnemer'
ORDER BY g.created_at DESC;

-- ========================================
-- STAP 5: Toon RBAC rollen en permissies
-- ========================================

-- Toon alle deelnemers met hun RBAC rollen
SELECT
    g.naam,
    g.email,
    g.rol as legacy_rol,
    r.name as rbac_rol,
    r.description as rol_beschrijving,
    ur.assigned_at as rol_toegewezen,
    (SELECT COUNT(*) FROM aanmeldingen WHERE gebruiker_id = g.id) as aantal_aanmeldingen
FROM gebruikers g
LEFT JOIN user_roles ur ON ur.user_id = g.id AND ur.is_active = true
LEFT JOIN roles r ON r.id = ur.role_id
WHERE g.rol IN ('deelnemer', 'begeleider', 'vrijwilliger')
ORDER BY g.created_at DESC;

-- Toon overzicht van rollen en hun aantal gebruikers
SELECT
    r.name as rol,
    r.description,
    COUNT(DISTINCT ur.user_id) as aantal_gebruikers
FROM roles r
LEFT JOIN user_roles ur ON ur.role_id = r.id AND ur.is_active = true
WHERE r.name IN ('deelnemer', 'begeleider', 'vrijwilliger')
GROUP BY r.id, r.name, r.description
ORDER BY aantal_gebruikers DESC;

-- ========================================
-- STAP 6: Check voor problemen
-- ========================================

-- Toon aanmeldingen zonder gebruikersaccount
SELECT 
    'WAARSCHUWING: Aanmeldingen zonder gebruikersaccount' as waarschuwing,
    a.id,
    a.naam,
    a.email,
    a.created_at
FROM aanmeldingen a
WHERE a.gebruiker_id IS NULL
AND a.email IS NOT NULL 
AND a.email != ''
ORDER BY a.created_at DESC;

-- ========================================
-- NOTITIES
-- ========================================

-- Standaard wachtwoord voor ALLE nieuwe accounts: DKL2025!
-- 
-- Gebruikers kunnen hun wachtwoord wijzigen via:
-- POST /api/auth/reset-password
-- Body: {
--   "huidig_wachtwoord": "DKL2025!",
--   "nieuw_wachtwoord": "hun_nieuwe_wachtwoord"
-- }
--
-- Login voorbeeld:
-- POST /api/auth/login
-- Body: {
--   "email": "diesbosje@hotmail.com",
--   "wachtwoord": "DKL2025!"
-- }