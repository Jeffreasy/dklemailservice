-- V1_33: Voeg gebruiker_id toe aan aanmeldingen tabel en maak gebruikersaccounts voor deelnemers

-- Stap 1: Voeg gebruiker_id kolom toe aan aanmeldingen
ALTER TABLE aanmeldingen
ADD COLUMN gebruiker_id UUID REFERENCES gebruikers(id);

-- Stap 2: Maak gebruikersaccounts voor alle deelnemers die nog geen account hebben
-- Password hash voor tijdelijk wachtwoord "DKL2025!" (dit moeten gebruikers later wijzigen)
INSERT INTO gebruikers (id, naam, email, wachtwoord_hash, rol, is_actief, newsletter_subscribed, created_at, updated_at)
SELECT 
    gen_random_uuid() as id,
    a.naam,
    a.email,
    '$2a$10$YourDefaultHashHere' as wachtwoord_hash,  -- Dit wordt vervangen door een echte hash
    'deelnemer' as rol,
    true as is_actief,
    false as newsletter_subscribed,
    NOW() as created_at,
    NOW() as updated_at
FROM aanmeldingen a
WHERE NOT EXISTS (
    SELECT 1 FROM gebruikers g WHERE g.email = a.email
)
AND a.email IS NOT NULL 
AND a.email != '';

-- Stap 3: Link bestaande aanmeldingen aan hun gebruikersaccounts
UPDATE aanmeldingen a
SET gebruiker_id = g.id
FROM gebruikers g
WHERE a.email = g.email
AND a.gebruiker_id IS NULL;

-- Stap 4: Voeg index toe voor snellere queries
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);

-- Stap 5: Voeg comment toe voor documentatie
COMMENT ON COLUMN aanmeldingen.gebruiker_id IS 'Link naar gebruikersaccount voor authenticatie en step tracking';