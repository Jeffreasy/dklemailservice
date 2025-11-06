-- V1_8__sync_contact_formulieren.sql
-- Synchroniseren van contact_formulieren tabel met Go model

-- Verwijder eerst oude constraint en velden indien nodig
ALTER TABLE contact_formulieren 
DROP COLUMN IF EXISTS onderwerp,
DROP COLUMN IF EXISTS ip_adres;

-- Toevoegen van ontbrekende velden voor email verwerking
ALTER TABLE contact_formulieren
ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS email_verzonden_op TIMESTAMP NULL;

-- Toevoegen van privacy gerelateerd veld
ALTER TABLE contact_formulieren
ADD COLUMN IF NOT EXISTS privacy_akkoord BOOLEAN NOT NULL DEFAULT false;

-- Toevoegen van behaneling gerelateerde velden
ALTER TABLE contact_formulieren
ADD COLUMN IF NOT EXISTS behandeld_door VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS behandeld_op TIMESTAMP NULL,
ADD COLUMN IF NOT EXISTS notities TEXT NULL;

-- Toevoegen van antwoord gerelateerde velden
ALTER TABLE contact_formulieren
ADD COLUMN IF NOT EXISTS beantwoord BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS antwoord_tekst TEXT NULL,
ADD COLUMN IF NOT EXISTS antwoord_datum TIMESTAMP NULL,
ADD COLUMN IF NOT EXISTS antwoord_door VARCHAR(255) NULL;

-- Controleren of de test_mode kolom al is toegevoegd (in V1_7)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='contact_formulieren' AND column_name='test_mode'
    ) THEN
        ALTER TABLE contact_formulieren
        ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT false;
    END IF;
END$$;

-- Indices toevoegen voor betere performance
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_email ON contact_formulieren(email);
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status ON contact_formulieren(status);

-- Opmerkingen toevoegen voor documentatie
COMMENT ON TABLE contact_formulieren IS 'Contactformulieren van de website';
COMMENT ON COLUMN contact_formulieren.test_mode IS 'Geeft aan of dit een testbericht is (geen echte email verzenden)';
COMMENT ON COLUMN contact_formulieren.email_verzonden IS 'Geeft aan of er een email is verzonden naar de afzender'; 