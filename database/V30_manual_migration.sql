-- V30 Handmatige Migratie - Direct Uitvoerbaar
-- Kopieer en plak dit in pgAdmin of psql

-- Stap 1: Voeg kolommen toe
ALTER TABLE participants ADD COLUMN IF NOT EXISTS account_type TEXT DEFAULT 'temporary';
ALTER TABLE participants ADD COLUMN IF NOT EXISTS registration_year INTEGER DEFAULT 2026;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS wachtwoord_hash TEXT;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS upgraded_to_gebruiker_id UUID;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS upgraded_at TIMESTAMPTZ;
ALTER TABLE participants ADD COLUMN IF NOT EXISTS has_app_access BOOLEAN DEFAULT FALSE;

-- Stap 2: Update bestaande NULL waarden
UPDATE participants SET account_type = 'temporary' WHERE account_type IS NULL;
UPDATE participants SET has_app_access = FALSE WHERE has_app_access IS NULL;
UPDATE participants SET registration_year = 2026 WHERE registration_year IS NULL;

-- Stap 3: Maak NOT NULL
ALTER TABLE participants ALTER COLUMN account_type SET NOT NULL;
ALTER TABLE participants ALTER COLUMN has_app_access SET NOT NULL;

-- Stap 4: Constraints
ALTER TABLE participants ADD CONSTRAINT check_account_type CHECK (account_type IN ('full', 'temporary'));

ALTER TABLE participants ADD CONSTRAINT fk_participants_upgraded_to_gebruiker 
FOREIGN KEY (upgraded_to_gebruiker_id) REFERENCES gebruikers(id) ON DELETE SET NULL;

-- Stap 5: Indexes
CREATE INDEX idx_participants_account_type ON participants(account_type);
CREATE INDEX idx_participants_registration_year ON participants(registration_year);
CREATE INDEX idx_participants_has_app_access ON participants(has_app_access);

-- Stap 6: Audit tabel
CREATE TABLE participant_upgrades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    gebruiker_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    upgraded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

CREATE INDEX idx_participant_upgrades_participant_id ON participant_upgrades(participant_id);
CREATE INDEX idx_participant_upgrades_gebruiker_id ON participant_upgrades(gebruiker_id);

-- Verificatie
SELECT 
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE account_type = 'full') as full_accounts,
    COUNT(*) FILTER (WHERE account_type = 'temporary') as temp_accounts
FROM participants;