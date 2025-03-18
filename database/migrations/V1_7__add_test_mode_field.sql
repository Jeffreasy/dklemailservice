-- Add test_mode field to contact_formulieren table
ALTER TABLE contact_formulieren
ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT false;

-- Add test_mode field to aanmeldingen table
ALTER TABLE aanmeldingen
ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT false; 