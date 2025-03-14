-- Migratie: 001_initial_schema.sql
-- Beschrijving: Initiële database setup
-- Versie: 1.0.0

-- Maak migraties tabel aan
CREATE TABLE IF NOT EXISTS migraties (
    id SERIAL PRIMARY KEY,
    versie VARCHAR(50) NOT NULL UNIQUE,
    naam VARCHAR(255) NOT NULL,
    toegepast TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak gebruikers tabel aan
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

-- Maak contact formulieren tabel aan
CREATE TABLE IF NOT EXISTS contact_formulieren (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    onderwerp VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw',
    ip_adres VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak contact antwoorden tabel aan
CREATE TABLE IF NOT EXISTS contact_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL REFERENCES contact_formulieren(id) ON DELETE CASCADE,
    onderwerp VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzonden_door UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak aanmeldingen tabel aan
CREATE TABLE IF NOT EXISTS aanmeldingen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    evenement VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw',
    extra_info JSONB,
    ip_adres VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak aanmelding antwoorden tabel aan
CREATE TABLE IF NOT EXISTS aanmelding_antwoorden (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aanmelding_id UUID NOT NULL REFERENCES aanmeldingen(id) ON DELETE CASCADE,
    onderwerp VARCHAR(255) NOT NULL,
    bericht TEXT NOT NULL,
    verzonden_op TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verzonden_door UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Maak email templates tabel aan
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

-- Maak verzonden emails tabel aan
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
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.0.0', 'Initiële database setup', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING; 