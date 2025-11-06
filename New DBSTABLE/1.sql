-- Gecombineerde Migratie: Geïntegreerde initiële database setup
-- Beschrijving: Combinatie van migraties 001 t/m 004, V1_05, V1_6, V1_8, V1_9, V1_10 voor een enkelvoudig initiëel schema
-- Versie: 1.0.0 (geconsolideerd)

-- Zorg ervoor dat de pgcrypto extensie beschikbaar is voor gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Maak migraties tabel aan
CREATE TABLE IF NOT EXISTS migraties (
    id SERIAL PRIMARY KEY,
    versie VARCHAR(50) NOT NULL UNIQUE,
    naam VARCHAR(255) NOT NULL,
    toegepast TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak gebruikers tabel aan (ongewijzigd)
CREATE TABLE IF NOT EXISTS gebruikers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    wachtwoord_hash VARCHAR(255) NOT NULL,
    rol VARCHAR(50) NOT NULL DEFAULT 'gebruiker',
    is_actief BOOLEAN NOT NULL DEFAULT TRUE,
    laatste_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak contact formulieren tabel aan (geïntegreerd met alters uit 003, V1_05, V1_8: verwijderd onderwerp, ip_adres; toegevoegd bericht default, email_verzonden, email_verzonden_op, privacy_akkoord, behandeld_door, behandeld_op, notities, beantwoord, antwoord_tekst, antwoord_datum, antwoord_door, test_mode)
CREATE TABLE IF NOT EXISTS contact_formulieren (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMP,
    privacy_akkoord BOOLEAN NOT NULL DEFAULT FALSE,
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP,
    notities TEXT,
    beantwoord BOOLEAN NOT NULL DEFAULT FALSE,
    antwoord_tekst TEXT,
    antwoord_datum TIMESTAMP,
    antwoord_door VARCHAR(255),
    test_mode BOOLEAN NOT NULL DEFAULT FALSE
);

-- Indices voor contact_formulieren (uit V1_8)
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_email ON contact_formulieren(email);
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status ON contact_formulieren(status);

-- Commentaren voor contact_formulieren (uit V1_8)
COMMENT ON TABLE contact_formulieren IS 'Contactformulieren van de website';
COMMENT ON COLUMN contact_formulieren.test_mode IS 'Geeft aan of dit een testbericht is (geen echte email verzenden)';
COMMENT ON COLUMN contact_formulieren.email_verzonden IS 'Geeft aan of er een email is verzonden naar de afzender';

-- Maak contact antwoorden tabel aan (geïntegreerd met alters uit 003, V1_10: verwijderd onderwerp, verzonden_door; toegevoegd tekst, verzond_op, verzond_door, email_verzonden; hernoemd bericht naar tekst)
CREATE TABLE IF NOT EXISTS contact_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL REFERENCES contact_formulieren(id) ON DELETE CASCADE,
    tekst TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzond_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzond_door VARCHAR(255) NOT NULL DEFAULT '',
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    verzonden_door VARCHAR(255)
);

-- Indices voor contact_antwoorden (uit V1_10)
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_contact_id ON contact_antwoorden(contact_id);
CREATE INDEX IF NOT EXISTS idx_contact_antwoorden_verzonden_door ON contact_antwoorden(verzonden_door);

-- Commentaren voor contact_antwoorden (uit V1_10)
COMMENT ON TABLE contact_antwoorden IS 'Antwoorden op contactformulieren';

-- Maak aanmeldingen tabel aan (geïntegreerd met alters uit 003, V1_05, V1_9: verwijderd evenement, ip_adres, extra_info; toegevoegd rol, afstand, ondersteuning, bijzonderheden, terms, email_verzonden, email_verzonden_op, behandeld_door, behandeld_op, notities, test_mode)
CREATE TABLE IF NOT EXISTS aanmeldingen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    rol VARCHAR(50),
    afstand VARCHAR(50),
    ondersteuning VARCHAR(255),
    bijzonderheden TEXT,
    terms BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMP,
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP,
    notities TEXT,
    test_mode BOOLEAN NOT NULL DEFAULT FALSE
);

-- Indices voor aanmeldingen (uit V1_9)
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_email ON aanmeldingen(email);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status ON aanmeldingen(status);
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_rol ON aanmeldingen(rol);

-- Commentaren voor aanmeldingen (uit V1_9)
COMMENT ON TABLE aanmeldingen IS 'Aanmeldingen voor De Koninklijke Loop';
COMMENT ON COLUMN aanmeldingen.rol IS 'Rol van de deelnemer (deelnemer, vrijwilliger, sponsor)';
COMMENT ON COLUMN aanmeldingen.afstand IS 'Gekozen afstand voor hardlopers';
COMMENT ON COLUMN aanmeldingen.test_mode IS 'Geeft aan of dit een testaanmelding is (geen echte email verzenden)';

-- Maak aanmelding antwoorden tabel aan (geïntegreerd met alters uit 003, V1_10: verwijderd onderwerp, verzonden_door; toegevoegd tekst, verzond_op, verzond_door, email_verzonden; hernoemd bericht naar tekst)
CREATE TABLE IF NOT EXISTS aanmelding_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aanmelding_id UUID NOT NULL REFERENCES aanmeldingen(id) ON DELETE CASCADE,
    tekst TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzond_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzond_door VARCHAR(255) NOT NULL DEFAULT '',
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    verzonden_door VARCHAR(255)
);

-- Indices voor aanmelding_antwoorden (uit V1_10)
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_aanmelding_id ON aanmelding_antwoorden(aanmelding_id);
CREATE INDEX IF NOT EXISTS idx_aanmelding_antwoorden_verzonden_door ON aanmelding_antwoorden(verzonden_door);

-- Commentaren voor aanmelding_antwoorden (uit V1_10)
COMMENT ON TABLE aanmelding_antwoorden IS 'Antwoorden op aanmeldingen';

-- Maak email templates tabel aan (ongewijzigd)
CREATE TABLE IF NOT EXISTS email_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL UNIQUE,
    onderwerp VARCHAR(255) NOT NULL,
    inhoud TEXT NOT NULL,
    beschrijving TEXT,
    is_actief BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES gebruikers(id)
);

-- Maak verzonden emails tabel aan (geïntegreerd met alter uit 003: toegevoegd fout_bericht)
CREATE TABLE IF NOT EXISTS verzonden_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ontvanger VARCHAR(255) NOT NULL,
    onderwerp VARCHAR(255) NOT NULL,
    inhoud TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'verzonden',
    contact_id UUID REFERENCES contact_formulieren(id),
    aanmelding_id UUID REFERENCES aanmeldingen(id),
    template_id UUID REFERENCES email_templates(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fout_bericht TEXT
);

-- Maak incoming_emails tabel aan (uit 004)
CREATE TABLE IF NOT EXISTS incoming_emails (
    id VARCHAR(255) PRIMARY KEY,
    message_id VARCHAR(255),
    "from" VARCHAR(255) NOT NULL,
    "to" VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    body TEXT,
    content_type VARCHAR(255),
    received_at TIMESTAMP NOT NULL,
    uid VARCHAR(255) UNIQUE,
    account_type VARCHAR(50),
    is_processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indices voor incoming_emails (uit 004)
CREATE INDEX IF NOT EXISTS idx_incoming_emails_message_id ON incoming_emails(message_id);
CREATE INDEX IF NOT EXISTS idx_incoming_emails_account_type ON incoming_emails(account_type);
CREATE INDEX IF NOT EXISTS idx_incoming_emails_is_processed ON incoming_emails(is_processed);

-- Maak notifications tabel aan (uit V1_6)
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    sent BOOLEAN NOT NULL DEFAULT FALSE,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indices voor notifications (uit V1_6)
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);
CREATE INDEX IF NOT EXISTS idx_notifications_sent ON notifications(sent);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Commentaar voor notifications (uit V1_6)
COMMENT ON TABLE notifications IS 'Stores notifications to be sent via Telegram';

-- Seed data (uit 002: admin gebruiker en email templates)
DO $$
DECLARE
    admin_id UUID;
BEGIN
    -- Maak admin gebruiker aan als niet bestaat (wachtwoord: admin)
    IF NOT EXISTS (SELECT 1 FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl') THEN
        INSERT INTO gebruikers (naam, email, wachtwoord_hash, rol, is_actief, created_at, updated_at)
        VALUES (
            'Admin',
            'admin@dekoninklijkeloop.nl',
            '$2a$10$5Yse5i2BJV.bwTzbmywa9e/3G.XxzQPayGPlTsut/nBrZr05pKMCK',
            'admin',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        );
    END IF;

    -- Haal admin ID op
    SELECT id INTO admin_id FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl';

    -- Maak standaard email templates aan als ze nog niet bestaan
    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'contact_admin_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'contact_admin_email',
            'Nieuw contactformulier',
            '<p>Er is een nieuw contactformulier ingevuld door {{.Contact.Naam}}.</p><p>Email: {{.Contact.Email}}</p><p>Bericht: {{.Contact.Bericht}}</p>',
            'Email die naar de admin wordt gestuurd bij een nieuw contactformulier',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'contact_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'contact_email',
            'Bedankt voor je bericht',
            '<p>Beste {{.Contact.Naam}},</p><p>Bedankt voor je bericht. We nemen zo snel mogelijk contact met je op.</p>',
            'Bevestigingsemail die naar de gebruiker wordt gestuurd bij een contactformulier',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'aanmelding_admin_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'aanmelding_admin_email',
            'Nieuwe aanmelding ontvangen',
            '<p>Er is een nieuwe aanmelding ontvangen van {{.Aanmelding.Naam}}.</p><p>Email: {{.Aanmelding.Email}}</p>',
            'Email die naar de admin wordt gestuurd bij een nieuwe aanmelding',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM email_templates WHERE naam = 'aanmelding_email') THEN
        INSERT INTO email_templates (naam, onderwerp, inhoud, beschrijving, is_actief, created_at, updated_at, created_by)
        VALUES (
            'aanmelding_email',
            'Bedankt voor je aanmelding',
            '<p>Beste {{.Aanmelding.Naam}},</p><p>Bedankt voor je aanmelding. We hebben je aanmelding ontvangen en zullen deze zo snel mogelijk verwerken.</p>',
            'Bevestigingsemail die naar de gebruiker wordt gestuurd bij een aanmelding',
            TRUE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP,
            admin_id
        );
    END IF;
END $$;

-- Registreer de gecombineerde migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.0.0', 'Geïntegreerde initiële database setup', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;