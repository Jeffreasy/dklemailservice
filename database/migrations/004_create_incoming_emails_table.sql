-- Migratie: 004_create_incoming_emails_table.sql
-- Beschrijving: Maakt de incoming_emails tabel aan voor email fetching functionaliteit
-- Versie: 1.0.3

-- Maak de incoming_emails tabel aan
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

-- Indexen aanmaken voor betere performance
CREATE INDEX IF NOT EXISTS idx_incoming_emails_message_id ON incoming_emails(message_id);
CREATE INDEX IF NOT EXISTS idx_incoming_emails_account_type ON incoming_emails(account_type);
CREATE INDEX IF NOT EXISTS idx_incoming_emails_is_processed ON incoming_emails(is_processed);

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.0.3', 'Create incoming_emails table', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING; 