-- Add steps column to aanmeldingen table (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name = 'aanmeldingen'
                   AND column_name = 'steps') THEN
        ALTER TABLE aanmeldingen ADD COLUMN steps INTEGER DEFAULT 0;
    END IF;
END $$;