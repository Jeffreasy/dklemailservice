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