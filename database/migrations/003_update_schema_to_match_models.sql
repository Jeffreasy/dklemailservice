-- Migratie: 003_update_schema_to_match_models.sql
-- Beschrijving: Update database schema om overeen te komen met Go models
-- Versie: 1.0.2

-- Aanmeldingen tabel aanpassen
ALTER TABLE aanmeldingen 
    -- Verwijder kolommen die niet in het model zitten
    DROP COLUMN IF EXISTS evenement,
    DROP COLUMN IF EXISTS ip_adres,
    DROP COLUMN IF EXISTS extra_info,
    
    -- Voeg nieuwe kolommen toe die in het model zitten
    ADD COLUMN IF NOT EXISTS rol VARCHAR(255),
    ADD COLUMN IF NOT EXISTS afstand VARCHAR(255),
    ADD COLUMN IF NOT EXISTS ondersteuning VARCHAR(255),
    ADD COLUMN IF NOT EXISTS bijzonderheden TEXT,
    ADD COLUMN IF NOT EXISTS terms BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS email_verzonden_op TIMESTAMP,
    ADD COLUMN IF NOT EXISTS behandeld_door VARCHAR(255),
    ADD COLUMN IF NOT EXISTS behandeld_op TIMESTAMP,
    ADD COLUMN IF NOT EXISTS notities TEXT;

-- Contact formulieren tabel aanpassen
ALTER TABLE contact_formulieren
    -- Verwijder kolommen die niet in het model zitten
    DROP COLUMN IF EXISTS onderwerp,
    DROP COLUMN IF EXISTS ip_adres,
    
    -- Voeg nieuwe kolommen toe die in het model zitten
    ADD COLUMN IF NOT EXISTS bericht TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS email_verzonden_op TIMESTAMP,
    ADD COLUMN IF NOT EXISTS privacy_akkoord BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS behandeld_door VARCHAR(255),
    ADD COLUMN IF NOT EXISTS behandeld_op TIMESTAMP,
    ADD COLUMN IF NOT EXISTS notities TEXT,
    ADD COLUMN IF NOT EXISTS beantwoord BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS antwoord_tekst TEXT,
    ADD COLUMN IF NOT EXISTS antwoord_datum TIMESTAMP,
    ADD COLUMN IF NOT EXISTS antwoord_door VARCHAR(255);

-- Contact antwoorden tabel aanpassen
ALTER TABLE contact_antwoorden
    -- Verwijder kolommen die niet in het model zitten
    DROP COLUMN IF EXISTS onderwerp,
    DROP COLUMN IF EXISTS verzonden_door,
    
    -- Voeg nieuwe kolommen toe die in het model zitten
    ADD COLUMN IF NOT EXISTS tekst TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS verzond_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS verzond_door VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT FALSE;

-- Aanmelding antwoorden tabel aanpassen
ALTER TABLE aanmelding_antwoorden
    -- Verwijder kolommen die niet in het model zitten
    DROP COLUMN IF EXISTS onderwerp,
    DROP COLUMN IF EXISTS verzonden_door,
    
    -- Voeg nieuwe kolommen toe die in het model zitten
    ADD COLUMN IF NOT EXISTS tekst TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS verzond_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS verzond_door VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS email_verzonden BOOLEAN NOT NULL DEFAULT FALSE;

-- Verzonden emails tabel aanpassen
ALTER TABLE verzonden_emails
    -- Voeg nieuwe kolommen toe die in het model zitten
    ADD COLUMN IF NOT EXISTS fout_bericht TEXT;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.0.2', 'Update schema om overeen te komen met Go models', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING; 