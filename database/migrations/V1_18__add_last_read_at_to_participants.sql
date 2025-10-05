-- Migratie: V1_18__add_last_read_at_to_participants.sql
-- Beschrijving: Add last_read_at column to chat_channel_participants
-- Versie: 1.18.0

ALTER TABLE chat_channel_participants
ADD COLUMN IF NOT EXISTS last_read_at TIMESTAMP WITH TIME ZONE;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.18.0', 'Add last_read_at to chat_channel_participants', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
