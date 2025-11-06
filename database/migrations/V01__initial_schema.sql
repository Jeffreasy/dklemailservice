-- GECONSOLIDEERDE V1 - INITIEEL SCHEMA (DEFINITIEVE VERSIE)
-- Dit bestand combineert de logica van 001, 003, 004, V1_05, V1_6, V1_8, V1_9, en V1_10.
-- FIX: Alle TIMESTAMP omgezet naar TIMESTAMPTZ voor GORM-compatibiliteit.
-- FIX 2: Alle VARCHAR omgezet naar TEXT voor GORM-compatibiliteit (verhelpt view-lock).
-- FIX 3: Typo 'TIMESTAMT_Z' gecorrigeerd naar 'TIMESTAMPTZ'.

-- Maak gebruikers tabel aan
CREATE TABLE IF NOT EXISTS gebruikers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    wachtwoord_hash TEXT NOT NULL,
    rol TEXT NOT NULL DEFAULT 'gebruiker',
    is_actief BOOLEAN NOT NULL DEFAULT TRUE,
    laatste_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak contact formulieren tabel aan (gecombineerde versie)
CREATE TABLE IF NOT EXISTS contact_formulieren (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam TEXT NOT NULL,
    email TEXT NOT NULL,
    bericht TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'nieuw',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Kolommen toegevoegd in 003 / V1_8
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMPTZ,
    privacy_akkoord BOOLEAN NOT NULL DEFAULT TRUE,
    behandeld_door TEXT,
    behandeld_op TIMESTAMPTZ,
    notities TEXT,
    beantwoord BOOLEAN NOT NULL DEFAULT FALSE,
    antwoord_tekst TEXT,
    antwoord_datum TIMESTAMPTZ,
    antwoord_door TEXT,

    -- Kolom toegevoegd in V1_05
    test_mode BOOLEAN NOT NULL DEFAULT false
);

-- Maak aanmeldingen tabel aan (gecombineerde versie)
CREATE TABLE IF NOT EXISTS aanmeldingen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam TEXT NOT NULL,
    email TEXT NOT NULL,
    telefoon TEXT,
    status TEXT NOT NULL DEFAULT 'nieuw',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Kolommen toegevoegd in 003 / V1_9
    rol TEXT NULL,
    afstand TEXT NULL,
    ondersteuning TEXT NULL,
    bijzonderheden TEXT NULL,
    terms BOOLEAN NOT NULL DEFAULT false,
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMPTZ NULL,
    behandeld_door TEXT NULL,
    behandeld_op TIMESTAMPTZ NULL,
    notities TEXT NULL,

    -- Kolom toegevoegd in V1_05
    test_mode BOOLEAN NOT NULL DEFAULT false
);

-- Maak contact antwoorden tabel aan (gecombineerde/gerepareerde versie)
CREATE TABLE IF NOT EXISTS contact_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL,
    verzonden_op TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Kolommen toegevoegd in V1_10 (repareert 001/003)
    tekst TEXT NOT NULL,
    email_verzonden BOOLEAN NOT NULL DEFAULT false,
    verzonden_door TEXT,

    -- Foreign Key uit V1_10
    CONSTRAINT fk_contact_antwoorden_contact_id
        FOREIGN KEY (contact_id) REFERENCES contact_formulieren (id)
        ON DELETE CASCADE
);

-- Maak aanmelding antwoorden tabel aan (gecombineerde/gerepareerde versie)
CREATE TABLE IF NOT EXISTS aanmelding_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aanmelding_id UUID NOT NULL,
    verzonden_op TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Kolommen toegevoegd in V1_10 (repareert 001/003)
    tekst TEXT NOT NULL,
    email_verzonden BOOLEAN NOT NULL DEFAULT false,
    verzonden_door TEXT,

    -- Foreign Key uit V1_10
    CONSTRAINT fk_aanmelding_antwoorden_aanmelding_id
        FOREIGN KEY (aanmelding_id) REFERENCES aanmeldingen (id)
        ON DELETE CASCADE
);

-- Maak email templates tabel aan (van 001)
CREATE TABLE IF NOT EXISTS email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam TEXT NOT NULL UNIQUE,
    onderwerp TEXT NOT NULL,
    inhoud TEXT NOT NULL,
    beschrijving TEXT,
    is_actief BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES gebruikers(id)
);

-- Maak verzonden emails tabel aan (gecombineerde versie)
CREATE TABLE IF NOT EXISTS verzonden_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ontvanger TEXT NOT NULL,
    onderwerp TEXT NOT NULL,
    inhoud TEXT NOT NULL,
    verzonden_op TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'verzonden',
    contact_id UUID REFERENCES contact_formulieren(id),
    aanmelding_id UUID REFERENCES aanmeldingen(id),
    template_id UUID REFERENCES email_templates(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Kolom toegevoegd in 003
    fout_bericht TEXT
);

-- Maak de incoming_emails tabel aan (van 004)
CREATE TABLE IF NOT EXISTS incoming_emails (
    id TEXT PRIMARY KEY,
    message_id TEXT,
    "from" TEXT NOT NULL,
    "to" TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT,
    content_type TEXT,
    received_at TIMESTAMPTZ NOT NULL, -- *** HIER WAS DE TYPO ***
    uid TEXT UNIQUE,
    account_type TEXT,
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak notifications tabel aan (van V1_6)
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL,
    priority TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    sent BOOLEAN NOT NULL DEFAULT FALSE,
    sent_at TIMESTAMPTZ, -- *** HIER WAS DE TYPO ***
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- --- GECONSOLIDEERDE INDEXEN ---
-- (van 004, V1_6, V1_8, V1_9, V1_10)

-- 004
CREATE INDEX IF NOT EXISTS idx_incoming_emails_message_id ON incoming_emails(message_id);
CREATE INDEX IF NOT EXISTS idx_incoming_emails_account_type ON incoming_emails(account_type);
CREATE INDEX IF NOT EXISTS idx_incoming_emails_is_processed ON incoming_emails(is_processed);

-- V1_6
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);
CREATE INDEX IF NOT EXISTS idx_notifications_sent ON notifications(sent);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- V1_8
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_email ON contact_formulieren(email);
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status ON contact_formulieren(status);

-- V1_9
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_email ON aanmeldingen(email);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status ON aanmeldingen(status);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_rol ON aanmeldingen(rol);

-- V1_10
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_verzonden_door ON contact_antwoorden(verzonden_door);
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_aanmelding_id ON aanmelding_antwoorden(aanmelding_id);
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_verzonden_door ON aanmelding_antwoorden(verzonden_door);


-- --- GECONSOLIDEERDE COMMENTAREN ---
-- (van V1_6, V1_8, V1_9, V1_10)

COMMENT ON TABLE notifications IS 'Stores notifications to be sent via Telegram';
COMMENT ON TABLE contact_formulieren IS 'Contactformulieren van de website';
COMMENT ON COLUMN contact_formulieren.test_mode IS 'Geeft aan of dit een testbericht is (geen echte email verzenden)';
COMMENT ON COLUMN contact_formulieren.email_verzonden IS 'Geeft aan of er een email is verzonden naar de afzender';
COMMENT ON TABLE aanmeldingen IS 'Aanmeldingen voor De Koninklijke Loop';
COMMENT ON COLUMN aanmeldingen.rol IS 'Rol van de deelnemer (deelnemer, vrijwilliger, sponsor)';
COMMENT ON COLUMN aanmeldingen.afstand IS 'Gekozen afstand voor hardlopers';
COMMENT ON COLUMN aanmeldingen.test_mode IS 'Geeft aan of dit een testaanmelding is (geen echte email verzenden)';
COMMENT ON TABLE contact_antwoorden IS 'Antwoorden op contactformulieren';
COMMENT ON TABLE aanmelding_antwoorden IS 'Antwoorden op aanmeldingen';