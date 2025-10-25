-- SQL Script om gebruikersaccounts aan te maken voor alle deelnemers
-- Wachtwoord voor alle nieuwe accounts: DKL2025!
-- Bcrypt hash (cost 10): $2a$10$YPFzRKvJe5vE0H0mxPqHq.VfZ8KQqX0YXJKxJK0fYdJH3LV.qhL8K

-- ========================================
-- STAP 1: Maak gebruikersaccounts aan
-- ========================================

-- Insert nieuwe gebruikers voor deelnemers die nog geen account hebben
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, newsletter_subscribed, created_at, updated_at)
SELECT 
    gen_random_uuid() as id,
    a.naam,
    a.email,
    '$2a$10$YPFzRKvJe5vE0H0mxPqHq.VfZ8KQqX0YXJKxJK0fYdJH3LV.qhL8K' as wachtwoord_hash,
    'deelnemer' as rol,
    true as is_actief,
    false as newsletter_subscribed,
    NOW() as created_at,
    NOW() as updated_at
FROM aanmeldingen a
WHERE NOT EXISTS (
    SELECT 1 FROM gebruikers g WHERE LOWER(g.email) = LOWER(a.email)
)
AND a.email IS NOT NULL 
AND a.email != ''
AND TRIM(a.email) != '';

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
-- STAP 3: Verificatie
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
-- STAP 5: Check voor problemen
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