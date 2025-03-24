-- V1_10__fix_antwoord_tables.sql
-- Repareren van de contact_antwoorden en aanmelding_antwoorden tabellen

-- Eerst de contact_antwoorden tabel bijwerken
ALTER TABLE contact_antwoorden
DROP COLUMN IF EXISTS onderwerp,
DROP COLUMN IF EXISTS bericht,
ADD COLUMN IF NOT EXISTS tekst TEXT NOT NULL,
ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS verzonden_door VARCHAR(255);

-- Indices aanmaken voor betere performance
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_verzonden_door ON contact_antwoorden(verzonden_door);

-- Daarna de aanmelding_antwoorden tabel bijwerken
ALTER TABLE aanmelding_antwoorden
DROP COLUMN IF EXISTS onderwerp,
DROP COLUMN IF EXISTS bericht,
ADD COLUMN IF NOT EXISTS tekst TEXT NOT NULL,
ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS verzonden_door VARCHAR(255);

-- Indices aanmaken voor betere performance
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_aanmelding_id ON aanmelding_antwoorden(aanmelding_id);
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_verzonden_door ON aanmelding_antwoorden(verzonden_door);

-- Opmerkingen toevoegen voor documentatie
COMMENT ON TABLE contact_antwoorden IS 'Antwoorden op contactformulieren';
COMMENT ON TABLE aanmelding_antwoorden IS 'Antwoorden op aanmeldingen';

-- Controleer en repareer de foreign keys indien nodig
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_contact_antwoorden_contact_id'
    ) THEN
        ALTER TABLE contact_antwoorden
        ADD CONSTRAINT fk_contact_antwoorden_contact_id
        FOREIGN KEY (contact_id) REFERENCES contact_formulieren (id)
        ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_aanmelding_antwoorden_aanmelding_id'
    ) THEN
        ALTER TABLE aanmelding_antwoorden
        ADD CONSTRAINT fk_aanmelding_antwoorden_aanmelding_id
        FOREIGN KEY (aanmelding_id) REFERENCES aanmeldingen (id)
        ON DELETE CASCADE;
    END IF;
END $$; 