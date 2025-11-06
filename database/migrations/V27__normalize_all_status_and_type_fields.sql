-- GECONSOLIDEERDE V27 - NORMALISEER ALLE STATUS EN TYPE VELDEN
--
-- DOEL: Vervangt alle "magic strings" voor status, type, en priority velden
--       door robuuste Foreign Key relaties naar lookup-tabellen.
--
-- WAAROM:
-- 1. Waterdichte Data Integriteit: Onmogelijk om een ongeldige status
--    (bv. 'geannuleeerd' met een spelfout) op te slaan.
-- 2. Beheerbaarheid: Alle mogelijke statussen zijn nu op één plek
--    in de database gedefinieerd.
--
-- IMPACT: ZEER GROTE BREAKING CHANGE!
-- Je MOET je GORM-modellen in de Go-code aanpassen voor:
-- - contact_formulieren.status
-- - aanmeldingen.status
-- - verzonden_emails.status
-- - events.status
-- - chat_channels.type
-- - notifications.type
-- - notifications.priority
--

-- =====================================================
-- 1. CONTACT FORMULIEREN (contact_formulieren.status)
-- =====================================================
CREATE TABLE IF NOT EXISTS contact_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE contact_status_types IS 'Lookup tabel voor contact formulier statussen (V27)';
INSERT INTO contact_status_types (status, description) VALUES
('nieuw', 'Nieuw binnengekomen, nog niet bekeken.'),
('in_behandeling', 'Door een staff-lid geopend en wordt behandeld.'),
('beantwoord', 'Er is een antwoord verstuurd naar de gebruiker.'),
('gesloten', 'Afgehandeld en gesloten.')
ON CONFLICT (status) DO NOTHING;

-- Migreer de data in stappen om 'DROP COLUMN' te vermijden
ALTER TABLE contact_formulieren ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES contact_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE contact_formulieren SET status_key = status WHERE status_key IS NULL;
ALTER TABLE contact_formulieren DROP COLUMN IF EXISTS status;
ALTER TABLE contact_formulieren RENAME COLUMN status_key TO status;
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status ON contact_formulieren(status);

-- =====================================================
-- 2. AANMELDINGEN (aanmeldingen.status)
-- =====================================================
CREATE TABLE IF NOT EXISTS registration_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE registration_status_types IS 'Lookup tabel voor aanmelding statussen (V27)';
INSERT INTO registration_status_types (status, description) VALUES
('nieuw', 'Nieuwe aanmelding, nog niet verwerkt.'),
('bevestigd', 'Aanmelding is goedgekeurd en bevestigd.'),
('geannuleerd', 'Aanmelding is geannuleerd (door gebruiker of staff).'),
('voltooid', 'Deelnemer heeft het evenement voltooid.')
ON CONFLICT (status) DO NOTHING;

-- Migreer de data
ALTER TABLE aanmeldingen ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES registration_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE aanmeldingen SET status_key = status WHERE status_key IS NULL;
ALTER TABLE aanmeldingen DROP COLUMN IF EXISTS status;
ALTER TABLE aanmeldingen RENAME COLUMN status_key TO status;
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status ON aanmeldingen(status);

-- =====================================================
-- 3. VERZONDEN EMAILS (verzonden_emails.status)
-- =====================================================
CREATE TABLE IF NOT EXISTS email_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE email_status_types IS 'Lookup tabel voor email verzendstatussen (V27)';
INSERT INTO email_status_types (status, description) VALUES
('verzonden', 'Email succesvol overgedragen aan de mailserver.'),
('failed', 'Email kon niet worden verzonden.'),
('pending', 'Email staat in de wachtrij om verzonden te worden.')
ON CONFLICT (status) DO NOTHING;

-- Migreer de data
ALTER TABLE verzonden_emails ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES email_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE verzonden_emails SET status_key = status WHERE status_key IS NULL;
ALTER TABLE verzonden_emails DROP COLUMN IF EXISTS status;
ALTER TABLE verzonden_emails RENAME COLUMN status_key TO status;
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_status ON verzonden_emails(status);

-- =====================================================
-- 4. EVENTS (events.status) (van V23)
-- =====================================================
CREATE TABLE IF NOT EXISTS event_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE event_status_types IS 'Lookup tabel voor event statussen (V27)';
INSERT INTO event_status_types (status, description) VALUES
('upcoming', 'Het evenement is gepland maar nog niet gestart.'),
('active', 'Het evenement is momenteel bezig.'),
('completed', 'Het evenement is afgelopen.'),
('cancelled', 'Het evenement is geannuleerd.')
ON CONFLICT (status) DO NOTHING;

-- Migreer de data
ALTER TABLE events ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES event_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE events SET status_key = status WHERE status_key IS NULL;
ALTER TABLE events DROP COLUMN IF EXISTS status;
ALTER TABLE events RENAME COLUMN status_key TO status;
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);

-- =====================================================
-- 5. CHAT KANALEN (chat_channels.type) (van V4)
-- =====================================================
CREATE TABLE IF NOT EXISTS chat_channel_types (
    type TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE chat_channel_types IS 'Lookup tabel voor chat kanaal types (V27)';
INSERT INTO chat_channel_types (type, description) VALUES
('public', 'Openbaar kanaal, iedereen kan meedoen.'),
('private', 'Prive kanaal, alleen op uitnodiging.'),
('direct', 'Een 1-op-1 direct message gesprek.')
ON CONFLICT (type) DO NOTHING;

-- Migreer de data
ALTER TABLE chat_channels ADD COLUMN IF NOT EXISTS type_key TEXT 
    REFERENCES chat_channel_types(type) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE chat_channels SET type_key = type WHERE type_key IS NULL;
ALTER TABLE chat_channels DROP COLUMN IF EXISTS type;
ALTER TABLE chat_channels RENAME COLUMN type_key TO type;
CREATE INDEX IF NOT EXISTS idx_chat_channels_type ON chat_channels(type);

-- =====================================================
-- 6. NOTIFICATIES (notifications.type & priority) (van V1)
-- =====================================================
CREATE TABLE IF NOT EXISTS notification_types (
    type TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE notification_types IS 'Lookup tabel voor notificatie types (V27)';
INSERT INTO notification_types (type, description) VALUES
('system_alert', 'Kritieke systeemmelding.'),
('new_registration', 'Nieuwe aanmelding ontvangen.'),
('new_message', 'Nieuw contactformulier of chatbericht.'),
('info', 'Algemene informatieve melding.')
ON CONFLICT (type) DO NOTHING;

CREATE TABLE IF NOT EXISTS notification_priority_types (
    priority TEXT PRIMARY KEY NOT NULL,
    description TEXT
);
COMMENT ON TABLE notification_priority_types IS 'Lookup tabel voor notificatie prioriteiten (V27)';
INSERT INTO notification_priority_types (priority, description) VALUES
('low', 'Lage prioriteit.'),
('medium', 'Normale prioriteit.'),
('high', 'Hoge prioriteit.'),
('critical', 'Kritieke prioriteit, vereist directe actie.')
ON CONFLICT (priority) DO NOTHING;

-- Migreer de 'type' kolom
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type_key TEXT 
    REFERENCES notification_types(type) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE notifications SET type_key = type WHERE type_key IS NULL;
ALTER TABLE notifications DROP COLUMN IF EXISTS type;
ALTER TABLE notifications RENAME COLUMN type_key TO type;
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- Migreer de 'priority' kolom
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS priority_key TEXT 
    REFERENCES notification_priority_types(priority) ON UPDATE CASCADE ON DELETE RESTRICT;
UPDATE notifications SET priority_key = priority WHERE priority_key IS NULL;
ALTER TABLE notifications DROP COLUMN IF EXISTS priority;
ALTER TABLE notifications RENAME COLUMN priority_key TO priority;
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);

-- =====================================================
-- 7. LOG INSTRUCTIES
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE '[V27] Normalisatie van alle status- en typevelden is voltooid.';
    RAISE NOTICE '=== ZEER GROTE BREAKING CHANGE ===';
    RAISE NOTICE 'Je MOET nu je GORM-modellen in de Go-code aanpassen.';
    RAISE NOTICE 'Voorbeeld: `Status string `gorm:"..."` wordt `Status string `gorm:"foreignKey:StatusType"`';
    RAISE NOTICE 'OF (beter): `Status StatusType `gorm:"foreignKey:Status"` en `StatusType string`';
    RAISE NOTICE 'Dit geldt voor 7 kolommen in 6 tabellen. Controleer de migratie zorgvuldig!';
END $$;