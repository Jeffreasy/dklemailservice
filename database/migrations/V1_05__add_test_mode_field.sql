-- Add test_mode field to contact_formulieren table if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'contact_formulieren' AND column_name = 'test_mode'
    ) THEN
        ALTER TABLE contact_formulieren
        ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT false;
    END IF;
END $$;

-- Add test_mode field to aanmeldingen table if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'aanmeldingen' AND column_name = 'test_mode'
    ) THEN
        ALTER TABLE aanmeldingen
        ADD COLUMN test_mode BOOLEAN NOT NULL DEFAULT false;
    END IF;
END $$; 