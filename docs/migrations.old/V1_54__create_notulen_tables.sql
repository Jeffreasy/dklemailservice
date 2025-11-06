-- V1_54__create_notulen_tables.sql
-- Create notulen tables with versioning support

-- Main notulen table
CREATE TABLE IF NOT EXISTS notulen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titel VARCHAR(255) NOT NULL UNIQUE,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[], -- Array of names
    afwezigen TEXT[],  -- Array of names
    agenda_items JSONB, -- Structured agenda items
    besluiten JSONB,   -- Structured decisions
    actiepunten JSONB, -- Structured action points
    notities TEXT,     -- Free text notes
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'finalized', 'archived')),
    versie INTEGER DEFAULT 1,
    created_by UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finalized_at TIMESTAMP WITH TIME ZONE,
    finalized_by UUID REFERENCES gebruikers(id),
    updated_by UUID REFERENCES gebruikers(id)  -- Added: Needed for the version trigger
);

-- Notulen versions table for history tracking
CREATE TABLE IF NOT EXISTS notulen_versies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notulen_id UUID NOT NULL REFERENCES notulen(id) ON DELETE CASCADE,
    versie INTEGER NOT NULL,
    titel VARCHAR(255) NOT NULL,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[],
    afwezigen TEXT[],
    agenda_items JSONB,
    besluiten JSONB,
    actiepunten JSONB,
    notities TEXT,
    status VARCHAR(50),
    gewijzigd_door UUID REFERENCES gebruikers(id),
    gewijzigd_op TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    wijziging_reden TEXT
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_notulen_datum ON notulen(vergadering_datum);
CREATE INDEX IF NOT EXISTS idx_notulen_status ON notulen(status);
CREATE INDEX IF NOT EXISTS idx_notulen_created_by ON notulen(created_by);
CREATE INDEX IF NOT EXISTS idx_notulen_finalized_by ON notulen(finalized_by);
CREATE INDEX IF NOT EXISTS idx_notulen_titel_gin ON notulen USING gin(to_tsvector('dutch', titel));
CREATE INDEX IF NOT EXISTS idx_notulen_notities_gin ON notulen USING gin(to_tsvector('dutch', notities));

-- Indexes for versions table
CREATE INDEX IF NOT EXISTS idx_notulen_versies_notulen_id ON notulen_versies(notulen_id);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_versie ON notulen_versies(notulen_id, versie);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_gewijzigd_door ON notulen_versies(gewijzigd_door);

-- Function to automatically create version snapshot on update
CREATE OR REPLACE FUNCTION create_notulen_version()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create version if this is an update (not insert)
    IF TG_OP = 'UPDATE' THEN
        -- Insert current version into history table
        INSERT INTO notulen_versies (
            notulen_id, versie, titel, vergadering_datum, locatie,
            voorzitter, notulist, aanwezigen, afwezigen,
            agenda_items, besluiten, actiepunten, notities, status,
            gewijzigd_door, wijziging_reden
        ) VALUES (
            OLD.id, OLD.versie, OLD.titel, OLD.vergadering_datum, OLD.locatie,
            OLD.voorzitter, OLD.notulist, OLD.aanwezigen, OLD.afwezigen,
            OLD.agenda_items, OLD.besluiten, OLD.actiepunten, OLD.notities, OLD.status,
            NEW.updated_by, 'Automatic version snapshot'
        );

        -- Increment version number
        NEW.versie = OLD.versie + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically create versions
DROP TRIGGER IF EXISTS notulen_version_trigger ON notulen;
CREATE OR REPLACE TRIGGER notulen_version_trigger
BEFORE UPDATE ON notulen
FOR EACH ROW
EXECUTE FUNCTION create_notulen_version();

-- Drop the trigger before the function to avoid dependency errors
DROP TRIGGER IF EXISTS notulen_updated_at_trigger ON notulen;
DROP FUNCTION IF EXISTS update_notulen_updated_at() CASCADE;  -- Added CASCADE for safety in case of other dependencies

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_notulen_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update timestamp
CREATE OR REPLACE TRIGGER notulen_updated_at_trigger
BEFORE UPDATE ON notulen
FOR EACH ROW
EXECUTE FUNCTION update_notulen_updated_at();