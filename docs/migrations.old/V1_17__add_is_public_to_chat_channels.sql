-- Migratie: V1_17__add_is_public_to_chat_channels.sql
-- Beschrijving: Add is_public column to chat_channels
-- Versie: 1.17.0

-- Add is_public column
ALTER TABLE chat_channels
ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT false;

-- Registreer de migratie
INSERT INTO migraties (versie, naam, toegepast) 
VALUES ('1.17.0', 'Add is_public to chat_channels', CURRENT_TIMESTAMP)
ON CONFLICT (versie) DO NOTHING;
