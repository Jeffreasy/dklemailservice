-- Database Schema Synchronisatie Script
-- Uitvoeren op: dklemailservice database
-- Datum: $(Get-Date -Format "yyyy-MM-dd")

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

-- V1_9__sync_aanmeldingen.sql
-- Synchroniseren van aanmeldingen tabel met Go model

-- Verwijder eerst oude constraint en velden indien nodig
ALTER TABLE aanmeldingen 
DROP COLUMN IF EXISTS evenement,
DROP COLUMN IF EXISTS extra_info,
DROP COLUMN IF EXISTS ip_adres;

-- Toevoegen van rol en gerelateerde velden voor hardlopers
ALTER TABLE aanmeldingen
ADD COLUMN IF NOT EXISTS rol VARCHAR(50) NULL,
ADD COLUMN IF NOT EXISTS afstand VARCHAR(50) NULL,
ADD COLUMN IF NOT EXISTS ondersteuning VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS bijzonderheden TEXT NULL;

-- Toevoegen van voorwaarden acceptatie
ALTER TABLE aanmeldingen
ADD COLUMN IF NOT EXISTS terms BOOLEAN NOT NULL DEFAULT false;

-- Toevoegen van velden voor email verwerking
ALTER TABLE aanmeldingen
ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS email_verzonden_op TIMESTAMP NULL;

-- Toevoegen van behandeling gerelateerde velden
ALTER TABLE aanmeldingen
ADD COLUMN IF NOT EXISTS behandeld_door VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS behandeld_op TIMESTAMP NULL,
ADD COLUMN IF NOT EXISTS notities TEXT NULL;

-- Controleren of de test_mode kolom al is toegevoegd (in V1_7)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name='aanmeldingen' AND column_name='test_mode'
    ) THEN
        ALTER TABLE aanmeldingen
        ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT false;
    END IF;
END$$;

-- Indices toevoegen voor betere performance
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_email ON aanmeldingen(email);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status ON aanmeldingen(status);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_rol ON aanmeldingen(rol);

-- Opmerkingen toevoegen voor documentatie
COMMENT ON TABLE aanmeldingen IS 'Aanmeldingen voor De Koninklijke Loop';
COMMENT ON COLUMN aanmeldingen.rol IS 'Rol van de deelnemer (deelnemer, vrijwilliger, sponsor)';
COMMENT ON COLUMN aanmeldingen.afstand IS 'Gekozen afstand voor hardlopers';
COMMENT ON COLUMN aanmeldingen.test_mode IS 'Geeft aan of dit een testaanmelding is (geen echte email verzenden)';

-- V1_10__fix_antwoord_tables.sql
-- Repareren van de contact_antwoorden en aanmelding_antwoorden tabellen

-- Eerst de contact_antwoorden tabel bijwerken
ALTER TABLE contact_antwoorden
DROP COLUMN IF EXISTS onderwerp,
DROP COLUMN IF EXISTS bericht,
ADD COLUMN IF NOT EXISTS tekst TEXT NULL,
ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT false;

-- Indices aanmaken voor betere performance
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_verzonden_door ON contact_antwoorden(verzonden_door);

-- Daarna de aanmelding_antwoorden tabel bijwerken
ALTER TABLE aanmelding_antwoorden
DROP COLUMN IF EXISTS onderwerp,
DROP COLUMN IF EXISTS bericht,
ADD COLUMN IF NOT EXISTS tekst TEXT NULL,
ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT false;

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

-- Registreer de migraties
INSERT INTO migraties (versie, naam, toegepast) 
VALUES 
  ('1.8', 'Synchronisatie van contact_formulieren', CURRENT_TIMESTAMP),
  ('1.9', 'Synchronisatie van aanmeldingen', CURRENT_TIMESTAMP),
  ('1.10', 'Reparatie van antwoord tabellen', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;

-- Meld voltooiing
DO $$
BEGIN
  RAISE NOTICE 'Database schema is nu gesynchroniseerd met de Go models!';
END $$; 