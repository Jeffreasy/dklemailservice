-- Migratie: V1_29__add_thumbnail_url_to_chat_messages.sql
-- Beschrijving: Add thumbnail_url column to chat_messages table for image thumbnails
-- Versie: 1.29.0

-- Add thumbnail_url column to chat_messages table
ALTER TABLE chat_messages
ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;

-- Add index for thumbnail_url if needed (optional, for performance)
-- CREATE INDEX IF NOT EXISTS idx_chat_messages_thumbnail_url ON chat_messages(thumbnail_url);

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast)
VALUES ('1.29.0', 'Add thumbnail_url to chat_messages', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;