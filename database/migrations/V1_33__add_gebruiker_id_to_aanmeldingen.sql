-- V1_33: Voeg gebruiker_id toe aan aanmeldingen tabel en link aan bestaande gebruikers

-- Stap 1: Voeg gebruiker_id kolom toe aan aanmeldingen (als deze nog niet bestaat)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'aanmeldingen'
        AND column_name = 'gebruiker_id'
    ) THEN
        ALTER TABLE aanmeldingen ADD COLUMN gebruiker_id UUID;
    END IF;
END $$;

-- Stap 2: Voeg foreign key constraint toe (als deze nog niet bestaat)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'aanmeldingen_gebruiker_id_fkey'
        AND table_name = 'aanmeldingen'
    ) THEN
        ALTER TABLE aanmeldingen
        ADD CONSTRAINT aanmeldingen_gebruiker_id_fkey
        FOREIGN KEY (gebruiker_id) REFERENCES gebruikers(id);
    END IF;
END $$;

-- Stap 3: Link bestaande aanmeldingen aan hun gebruikersaccounts (alleen voor bestaande emails)
UPDATE aanmeldingen a
SET gebruiker_id = g.id
FROM gebruikers g
WHERE a.email = g.email
AND a.gebruiker_id IS NULL
AND a.email IS NOT NULL
AND a.email != '';

-- Stap 4: Voeg index toe voor snellere queries
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);

-- Stap 5: Voeg comment toe voor documentatie
COMMENT ON COLUMN aanmeldingen.gebruiker_id IS 'Link naar gebruikersaccount voor authenticatie en step tracking';

-- Note: Gebruik scripts/create_participant_accounts.go om gebruikersaccounts aan te maken voor deelnemers zonder account