-- V1_60__add_notulen_updated_by_column.sql
-- Add updated_by column to notulen table if it doesn't exist

DO $$ 
BEGIN
    -- Check if updated_by column exists in notulen table
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'notulen' 
        AND column_name = 'updated_by'
    ) THEN
        -- Add the column
        ALTER TABLE notulen 
        ADD COLUMN updated_by UUID REFERENCES gebruikers(id);
        
        -- Add comment for documentation
        COMMENT ON COLUMN notulen.updated_by IS 'UUID of the user who last updated this notulen';
        
        -- Create index for better query performance
        CREATE INDEX IF NOT EXISTS idx_notulen_updated_by ON notulen(updated_by);
        
        RAISE NOTICE 'Added updated_by column to notulen table';
    ELSE
        RAISE NOTICE 'Column updated_by already exists in notulen table';
    END IF;
END $$;