-- ============================================================================
-- V27 CONSOLIDATED: Normalize All Status and Type Fields
-- ============================================================================
-- Consolidates V27_01 through V27_57
-- Purpose: Convert string status/type fields to normalized lookup tables
--          with foreign key constraints for data integrity
-- ============================================================================

-- ----------------------------------------------------------------------------
-- PART 1: CONTACT STATUS NORMALIZATION
-- ----------------------------------------------------------------------------

-- Create lookup table
CREATE TABLE IF NOT EXISTS contact_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE contact_status_types IS 'Lookup table voor contact formulier statussen (V27)';

-- Seed values
INSERT INTO contact_status_types (status, description) VALUES
('nieuw', 'Nieuw contactverzoek, nog niet bekeken'),
('in_behandeling', 'Contactverzoek wordt behandeld'),
('beantwoord', 'Contactverzoek is beantwoord'),
('gesloten', 'Contactverzoek is afgesloten')
ON CONFLICT (status) DO NOTHING;

-- Add temporary FK column
ALTER TABLE contact_formulieren ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES contact_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;

-- Migrate data
UPDATE contact_formulieren SET status_key = status WHERE status_key IS NULL;

-- Drop old column and rename
ALTER TABLE contact_formulieren DROP COLUMN IF EXISTS status CASCADE;
ALTER TABLE contact_formulieren RENAME COLUMN status_key TO status;

-- Create index
CREATE INDEX IF NOT EXISTS idx_contact_formulieren_status ON contact_formulieren(status);

-- ----------------------------------------------------------------------------
-- PART 2: REGISTRATION STATUS NORMALIZATION
-- ----------------------------------------------------------------------------

-- Create lookup table
CREATE TABLE IF NOT EXISTS registration_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE registration_status_types IS 'Lookup table voor aanmelding/registratie statussen (V27)';

-- Seed values
INSERT INTO registration_status_types (status, description) VALUES
('nieuw', 'Nieuwe aanmelding, nog niet bekeken'), -- ('registered', 'Geregistreerd, wacht op bevestiging'),
('confirmed', 'Bevestigd door admin'),
('waiting_list', 'Op wachtlijst geplaatst'),
('cancelled', 'Geannuleerd door deelnemer of admin'),
('attended', 'Heeft deelgenomen aan het evenement'),
('no_show', 'Niet verschenen bij het evenement'),
('beantwoord', 'Er is een antwoord gestuurd')
ON CONFLICT (status) DO NOTHING;

-- Add temporary FK column
ALTER TABLE aanmeldingen ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES registration_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;

-- Migrate data
UPDATE aanmeldingen SET status_key = status WHERE status_key IS NULL;

-- Drop old column and rename
ALTER TABLE aanmeldingen DROP COLUMN IF EXISTS status CASCADE;
ALTER TABLE aanmeldingen RENAME COLUMN status_key TO status;

-- Create index
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_status ON aanmeldingen(status);

-- ----------------------------------------------------------------------------
-- PART 3: EMAIL STATUS NORMALIZATION
-- ----------------------------------------------------------------------------

-- Create lookup table
CREATE TABLE IF NOT EXISTS email_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE email_status_types IS 'Lookup table voor email statussen (V27)';

-- Seed values
INSERT INTO email_status_types (status, description) VALUES
('pending', 'Email staat in de wachtrij'),
('sending', 'Email wordt verzonden'),
('verzonden', 'Email succesvol verzonden'),
('failed', 'Email verzenden mislukt'),
('bounced', 'Email teruggestuurd (bounce)'),
('delivered', 'Email is afgeleverd bij ontvanger')
ON CONFLICT (status) DO NOTHING;

-- Add temporary FK column
ALTER TABLE verzonden_emails ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES email_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;

-- Migrate data
UPDATE verzonden_emails SET status_key = status WHERE status_key IS NULL;

-- Drop old column and rename
ALTER TABLE verzonden_emails DROP COLUMN IF EXISTS status CASCADE;
ALTER TABLE verzonden_emails RENAME COLUMN status_key TO status;

-- Create index
CREATE INDEX IF NOT EXISTS idx_verzonden_emails_status ON verzonden_emails(status);

-- ----------------------------------------------------------------------------
-- PART 4: EVENT STATUS NORMALIZATION
-- ----------------------------------------------------------------------------

-- Create lookup table
CREATE TABLE IF NOT EXISTS event_status_types (
    status TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE event_status_types IS 'Lookup table voor event statussen (V27)';

-- Seed values
INSERT INTO event_status_types (status, description) VALUES
('draft', 'Event in concept fase'),
('upcoming', 'Aankomend event (gepland)'),
('active', 'Event is momenteel actief/bezig'),
('completed', 'Event is afgerond'),
('cancelled', 'Event is geannuleerd'),
('postponed', 'Event is uitgesteld')
ON CONFLICT (status) DO NOTHING;

-- Add temporary FK column
ALTER TABLE events ADD COLUMN IF NOT EXISTS status_key TEXT 
    REFERENCES event_status_types(status) ON UPDATE CASCADE ON DELETE RESTRICT;

-- Migrate data
UPDATE events SET status_key = status WHERE status_key IS NULL;

-- Drop old column and rename
ALTER TABLE events DROP COLUMN IF EXISTS status CASCADE;
ALTER TABLE events RENAME COLUMN status_key TO status;

-- Create index
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);

-- ----------------------------------------------------------------------------
-- PART 5: CHAT CHANNEL TYPE NORMALIZATION
-- ----------------------------------------------------------------------------

-- Create lookup table
CREATE TABLE IF NOT EXISTS chat_channel_types (
    type TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE chat_channel_types IS 'Lookup table voor chat channel types (V27)';

-- Seed values
INSERT INTO chat_channel_types (type, description) VALUES
('public', 'Publiek toegankelijk kanaal'),
('private', 'Privé kanaal (alleen op uitnodiging)'),
('direct', 'Direct bericht tussen twee gebruikers'),
('group', 'Groepskanaal (meerdere gebruikers)'),
('announcement', 'Aankondigingen kanaal (readonly voor meeste users)')
ON CONFLICT (type) DO NOTHING;

-- Add temporary FK column
ALTER TABLE chat_channels ADD COLUMN IF NOT EXISTS type_key TEXT 
    REFERENCES chat_channel_types(type) ON UPDATE CASCADE ON DELETE RESTRICT;

-- Migrate data
UPDATE chat_channels SET type_key = type WHERE type_key IS NULL;

-- Drop old column and rename
ALTER TABLE chat_channels DROP COLUMN IF EXISTS type CASCADE;
ALTER TABLE chat_channels RENAME COLUMN type_key TO type;

-- Create index
CREATE INDEX IF NOT EXISTS idx_chat_channels_type ON chat_channels(type);

-- ----------------------------------------------------------------------------
-- PART 6: NOTIFICATION TYPE AND PRIORITY NORMALIZATION
-- ----------------------------------------------------------------------------

-- Drop existing tables if they exist with wrong structure
DROP TABLE IF EXISTS notification_types CASCADE;
DROP TABLE IF EXISTS notification_priority_types CASCADE;

-- Create notification types lookup table
CREATE TABLE notification_types (
    type TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE notification_types IS 'Lookup table voor notificatie types (V27)';

-- Seed notification types
INSERT INTO notification_types (type, description) VALUES
('contact', 'Notificatie voor nieuw contactverzoek'),
('aanmelding', 'Notificatie voor nieuwe aanmelding/registratie'),
('auth', 'Notificatie voor authenticatie events'),
('system', 'System notificaties (startup, shutdown, errors)'),
('health', 'Health check notificaties')
ON CONFLICT (type) DO NOTHING;

-- Create notification priority types lookup table
CREATE TABLE notification_priority_types (
    priority TEXT PRIMARY KEY NOT NULL,
    description TEXT
);

COMMENT ON TABLE notification_priority_types IS 'Lookup table voor notificatie prioriteiten (V27)';

-- Seed notification priorities
INSERT INTO notification_priority_types (priority, description) VALUES
('low', 'Lage prioriteit - informationeel'),
('medium', 'Normale prioriteit'),
('high', 'Hoge prioriteit - vereist aandacht'),
('critical', 'Kritiek - vereist onmiddellijke aandacht')
ON CONFLICT (priority) DO NOTHING;

-- Migrate notification type field
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type_key TEXT
    REFERENCES notification_types(type) ON UPDATE CASCADE ON DELETE RESTRICT;

UPDATE notifications SET type_key = type WHERE type_key IS NULL;

ALTER TABLE notifications DROP COLUMN IF EXISTS type CASCADE;
ALTER TABLE notifications RENAME COLUMN type_key TO type;

CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- Migrate notification priority field
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS priority_key TEXT
    REFERENCES notification_priority_types(priority) ON UPDATE CASCADE ON DELETE RESTRICT;

UPDATE notifications SET priority_key = priority WHERE priority_key IS NULL;

ALTER TABLE notifications DROP COLUMN IF EXISTS priority CASCADE;
ALTER TABLE notifications RENAME COLUMN priority_key TO priority;

CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);

-- ----------------------------------------------------------------------------
-- COMPLETION LOG
-- ----------------------------------------------------------------------------

DO $$
BEGIN
    RAISE NOTICE '[V27 CONSOLIDATED] Normalisatie van statussen en types is voltooid.';
    RAISE NOTICE '=== BREAKING CHANGE ===';
    RAISE NOTICE '1. (Go Code) Pas je GORM-modellen aan voor de genormaliseerde velden.';
    RAISE NOTICE '   - contact_formulieren.status: FK naar contact_status_types';
    RAISE NOTICE '   - aanmeldingen.status: FK naar registration_status_types';
    RAISE NOTICE '   - verzonden_emails.status: FK naar email_status_types';
    RAISE NOTICE '   - events.status: FK naar event_status_types';
    RAISE NOTICE '   - chat_channels.type: FK naar chat_channel_types';
    RAISE NOTICE '   - notifications.type: FK naar notification_types';
    RAISE NOTICE '   - notifications.priority: FK naar notification_priority_types';
    RAISE NOTICE '2. Alle lookup tables zijn aangemaakt en data is gemigreerd.';
END $$;