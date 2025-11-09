-- GECONSOLIDEERDE V24 - COMPLETE NOTULEN MODULE
-- Dit bestand vervangt V1_54, V1_55, V1_56, V1_57, V1_58, V1_59, en V1_60.
-- Het combineert het schema, de triggers (met fixes), permissies, data en helpers.

-- =====================================================
-- 1. SCHEMA (Gecoördineerd van V1_54, V1_58, V1_60)
-- =====================================================

-- Main notulen table (met alle kolommen)
CREATE TABLE IF NOT EXISTS notulen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titel VARCHAR(255) NOT NULL UNIQUE,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[], -- Legacy
    afwezigen TEXT[],  -- Legacy
    agenda_items JSONB,
    besluiten JSONB,
    actiepunten JSONB,
    notities TEXT,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'finalized', 'archived')),
    versie INTEGER DEFAULT 1,
    created_by UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finalized_at TIMESTAMP WITH TIME ZONE,
    finalized_by UUID REFERENCES gebruikers(id),
    updated_by UUID REFERENCES gebruikers(id), -- Van V1_60 (en V1_54)
    -- Van V1_58
    aanwezigen_gebruikers UUID[],
    afwezigen_gebruikers UUID[],
    aanwezigen_gasten TEXT[],
    afwezigen_gasten TEXT[]
);

-- Notulen versions table (met alle kolommen)
CREATE TABLE IF NOT EXISTS notulen_versies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notulen_id UUID NOT NULL REFERENCES notulen(id) ON DELETE CASCADE,
    versie INTEGER NOT NULL,
    titel VARCHAR(255) NOT NULL,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[], -- Legacy
    afwezigen TEXT[], -- Legacy
    agenda_items JSONB,
    besluiten JSONB,
    actiepunten JSONB,
    notities TEXT,
    status VARCHAR(50),
    gewijzigd_door UUID REFERENCES gebruikers(id),
    gewijzigd_op TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    wijziging_reden TEXT,
    -- Van V1_58
    aanwezigen_gebruikers UUID[],
    afwezigen_gebruikers UUID[],
    aanwezigen_gasten TEXT[],
    afwezigen_gasten TEXT[]
);

-- =====================================================
-- 2. INDEXEN (Gecoördineerd van V1_54, V1_58, V1_60)
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_notulen_datum ON notulen(vergadering_datum);
CREATE INDEX IF NOT EXISTS idx_notulen_status ON notulen(status);
CREATE INDEX IF NOT EXISTS idx_notulen_created_by ON notulen(created_by);
CREATE INDEX IF NOT EXISTS idx_notulen_finalized_by ON notulen(finalized_by);
CREATE INDEX IF NOT EXISTS idx_notulen_titel_gin ON notulen USING gin(to_tsvector('dutch', titel));
CREATE INDEX IF NOT EXISTS idx_notulen_notities_gin ON notulen USING gin(to_tsvector('dutch', notities));
CREATE INDEX IF NOT EXISTS idx_notulen_updated_by ON notulen(updated_by); -- Van V1_60

CREATE INDEX IF NOT EXISTS idx_notulen_versies_notulen_id ON notulen_versies(notulen_id);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_versie ON notulen_versies(notulen_id, versie);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_gewijzigd_door ON notulen_versies(gewijzigd_door);

-- GIN indexen voor participant UUIDs (van V1_58)
CREATE INDEX IF NOT EXISTS idx_notulen_aanwezigen_gebruikers ON notulen USING GIN(aanwezigen_gebruikers);
CREATE INDEX IF NOT EXISTS idx_notulen_afwezigen_gebruikers ON notulen USING GIN(afwezigen_gebruikers);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_aanwezigen_gebruikers ON notulen_versies USING GIN(aanwezigen_gebruikers);
CREATE INDEX IF NOT EXISTS idx_notulen_versies_afwezigen_gebruikers ON notulen_versies USING GIN(afwezigen_gebruikers);

-- =====================================================
-- 3. TRIGGERS (Definitieve versie van V1_59 + V19)
-- =====================================================

-- De definitieve, correcte versie-functie uit V1_59
CREATE OR REPLACE FUNCTION create_notulen_version()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create version if this is an update (not insert)
    IF TG_OP = 'UPDATE' THEN
        -- Insert current version into history table with ALL fields
        INSERT INTO notulen_versies (
            notulen_id, versie, titel, vergadering_datum, locatie,
            voorzitter, notulist, 
            aanwezigen, afwezigen,  -- Legacy fields
            aanwezigen_gebruikers, afwezigen_gebruikers,  -- NEW: UUID arrays
            aanwezigen_gasten, afwezigen_gasten,          -- NEW: Guest text arrays
            agenda_items, besluiten, actiepunten, notities, status,
            gewijzigd_door, wijziging_reden
        ) VALUES (
            OLD.id, OLD.versie, OLD.titel, OLD.vergadering_datum, OLD.locatie,
            OLD.voorzitter, OLD.notulist, 
            OLD.aanwezigen, OLD.afwezigen,  -- Legacy fields
            OLD.aanwezigen_gebruikers, OLD.afwezigen_gebruikers,  -- NEW: UUID arrays
            OLD.aanwezigen_gasten, OLD.afwezigen_gasten,          -- NEW: Guest text arrays
            OLD.agenda_items, OLD.besluiten, OLD.actiepunten, OLD.notities, OLD.status,
            COALESCE(NEW.updated_by, NEW.created_by), 'Automatic version snapshot'
        );

        -- Increment version number
        NEW.versie = OLD.versie + 1;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION create_notulen_version IS 'Creates version snapshot with complete participant data (UUIDs + guests)';

-- Koppel de versie-trigger
DROP TRIGGER IF EXISTS notulen_version_trigger ON notulen;
CREATE TRIGGER notulen_version_trigger
BEFORE UPDATE ON notulen
FOR EACH ROW
EXECUTE FUNCTION create_notulen_version();

-- Koppel de generieke updated_at trigger (van V19)
DROP TRIGGER IF EXISTS notulen_updated_at_trigger ON notulen;
CREATE TRIGGER notulen_updated_at_trigger
BEFORE UPDATE ON notulen
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 4. PERMISSIES (van V1_55)
-- =====================================================
INSERT INTO permissions (resource, action, description) VALUES
('notulen', 'read', 'Kan notulen lezen en bekijken'),
('notulen', 'write', 'Kan notulen aanmaken en bijwerken'),
('notulen', 'delete', 'Kan notulen verwijderen'),
('notulen', 'finalize', 'Kan notulen finaliseren'),
('notulen', 'archive', 'Kan notulen archiveren')
ON CONFLICT (resource, action) DO NOTHING;

-- Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND p.resource = 'notulen'
ON CONFLICT DO NOTHING;

-- Staff
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' AND p.resource = 'notulen' AND p.action IN ('read', 'write')
ON CONFLICT DO NOTHING;

-- =====================================================
-- 5. SAMPLE DATA (van V1_56)
-- =====================================================
DO $$
DECLARE
    admin_user_id UUID;
BEGIN
    SELECT id INTO admin_user_id FROM gebruikers WHERE email = 'admin@dekoninklijkeloop.nl' LIMIT 1;

    IF admin_user_id IS NOT NULL THEN
        IF NOT EXISTS (SELECT 1 FROM notulen WHERE titel = 'Notulen 30 oktober 2025') THEN
            INSERT INTO notulen (
                titel,
                vergadering_datum,
                locatie,
                voorzitter,
                aanwezigen, -- Legacy
                afwezigen, -- Legacy
                agenda_items,
                besluiten,
                actiepunten,
                notities,
                status,
                created_by
            ) VALUES (
                'Notulen 30 oktober 2025',
                '2025-10-30',
                'Teams',
                'Salih',
                ARRAY['Salih', 'Angelique', 'Jeffrey', 'Marieke', 'Lida', 'Ginelli'],
                ARRAY[]::TEXT[],
                -- GEFIXT: Dit is nu een array en gebruikt "titel" en "beschrijving"
                '[{"titel": "Afspraken voor DKL 2026", "beschrijving": "Datum, Doel, Streefdoel, Organisatie teams, Vaste meetings, Partners, Website & online"}, {"titel": "Ideeën voor DKL26", "beschrijving": "Social media, Mascotte, Merch, Betrokkenheid doelgroep"}]'::JSONB,
                -- GEFIXT: Dit is nu een array en gebruikt "beschrijving"
                '[{"beschrijving": "Datum: zaterdag 16 mei 2026"}, {"beschrijving": "Doel: Only Friends[](https://onlyfriends.nl/), stichting in Apeldoorn, goed doel omdat het lokaal is en mooi aansluit bij samenwerking met Alessandro Bistolfi."}, {"beschrijving": "Streefdoel aantal deelnemers: 100"}, {"beschrijving": "Verdeling van de organisatie in miniteams"}, {"beschrijving": "Vaste meetings: Eenmaal in de 6 weken op de maandag in Teams"}, {"beschrijving": "Partners: Accres (beheer sportfaciliteiten en evenementen voor de gemeente Apeldoorn), Alessandro Bistolfi van Nedarg handbikes[](https://www.nedarg.com)"}, {"beschrijving": "Website & online: Stappenteller ontwikkeld waaraan donatie optie gekoppeld, iedere deelnemer krijgt account voor prestaties en gezamenlijke teller, opties via site of app"}]'::JSONB,
                -- GEFIXT: Dit is nu een array en gebruikt "beschrijving" en "status"
                '[{"beschrijving": "Salih gaat een meeting cyclus op Teams aanmaken.", "verantwoordelijke": "Salih", "status": "pending"}, {"beschrijving": "Salih neemt contact op met Alessandro.", "verantwoordelijke": "Salih", "status": "pending"}, {"beschrijving": "Lida kijkt naar scholen in Harderwijk.", "verantwoordelijke": "Lida", "status": "pending"}, {"beschrijving": "Salih kijkt naar optie bij Nijmegen: https://www.maartenschool.nl/home.", "verantwoordelijke": "Salih", "status": "pending"}, {"beschrijving": "Jeffrey deelt eerste versie van de app.", "verantwoordelijke": "Jeffrey", "status": "pending"}, {"beschrijving": "Angelique benadert collega''s via intranet.", "verantwoordelijke": "Angelique", "status": "pending"}, {"beschrijving": "Lida en Ginelli gaan aan de slag met social media content (vlogs, live, intro videos).", "verantwoordelijke": "Lida", "status": "pending"}, {"beschrijving": "Brainstorm over eigen mascotte (met doelgroep bedenken of zelf ontwikkelen).", "verantwoordelijke": "Team", "status": "pending"}, {"beschrijving": "Angelique kijkt naar merch opties (draagtasjes, bekers, later T-shirts) met Wesleys en collega grafisch ontwerper.", "verantwoordelijke": "Angelique", "status": "pending"}, {"beschrijving": "Salih informeert bij Kathelijn over betrokkenheid doelgroep met cliëntenraad (pas vanaf januari).", "verantwoordelijke": "Salih", "status": "pending"}]'::JSONB,
                'Pas vanaf januari doelgroep betrekken, anders te vroeg.',
                'draft',
                admin_user_id
            );
        END IF;
    END IF;
END $$;


-- =====================================================
-- 6. HELPER FUNCTIES & VIEWS (van V1_58)
-- =====================================================
CREATE OR REPLACE FUNCTION get_user_names_from_uuids(user_uuids UUID[])
RETURNS TEXT[] AS $$
DECLARE
    user_names TEXT[];
BEGIN
    IF user_uuids IS NULL OR array_length(user_uuids, 1) = 0 THEN
        RETURN ARRAY[]::TEXT[];
    END IF;

    SELECT array_agg(naam ORDER BY naam)
    INTO user_names
    FROM gebruikers
    WHERE id = ANY(user_uuids)
    AND is_actief = true;

    RETURN COALESCE(user_names, ARRAY[]::TEXT[]);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_names_from_uuids IS 'Converteert array van user UUIDs naar array van user namen';

CREATE OR REPLACE FUNCTION get_user_uuids_from_names(user_names TEXT[])
RETURNS UUID[] AS $$
DECLARE
    user_uuids UUID[];
BEGIN
    IF user_names IS NULL OR array_length(user_names, 1) = 0 THEN
        RETURN ARRAY[]::UUID[];
    END IF;

    SELECT array_agg(id ORDER BY naam)
    INTO user_uuids
    FROM gebruikers
    WHERE naam = ANY(user_names)
    AND is_actief = true;

    RETURN COALESCE(user_uuids, ARRAY[]::UUID[]);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_uuids_from_names IS 'Converteert array van user namen naar array van user UUIDs (best effort)';

-- Views
DROP VIEW IF EXISTS notulen_with_participants CASCADE;
DROP VIEW IF EXISTS notulen_versies_with_participants CASCADE;

CREATE VIEW notulen_with_participants AS
SELECT
    n.*,
    -- Combined aanwezigen (users + guests)
    CASE
        WHEN n.aanwezigen_gebruikers IS NOT NULL AND n.aanwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(n.aanwezigen_gebruikers) || n.aanwezigen_gasten
        WHEN n.aanwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(n.aanwezigen_gebruikers)
        WHEN n.aanwezigen_gasten IS NOT NULL THEN
            n.aanwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as aanwezigen_combined,
    -- Combined afwezigen (users + guests)
    CASE
        WHEN n.afwezigen_gebruikers IS NOT NULL AND n.afwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(n.afwezigen_gebruikers) || n.afwezigen_gasten
        WHEN n.afwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(n.afwezigen_gebruikers)
        WHEN n.afwezigen_gasten IS NOT NULL THEN
            n.afwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as afwezigen_combined
FROM notulen n;

COMMENT ON VIEW notulen_with_participants IS 'Notulen view met gecombineerde participant lijsten (voor API backwards compatibility)';

CREATE VIEW notulen_versies_with_participants AS
SELECT
    nv.*,
    -- Combined aanwezigen (users + guests)
    CASE
        WHEN nv.aanwezigen_gebruikers IS NOT NULL AND nv.aanwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(nv.aanwezigen_gebruikers) || nv.aanwezigen_gasten
        WHEN nv.aanwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(nv.aanwezigen_gebruikers)
        WHEN nv.aanwezigen_gasten IS NOT NULL THEN
            nv.aanwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as aanwezigen_combined,
    -- Combined afwezigen (users + guests)
    CASE
        WHEN nv.afwezigen_gebruikers IS NOT NULL AND nv.afwezigen_gasten IS NOT NULL THEN
            get_user_names_from_uuids(nv.afwezigen_gebruikers) || nv.afwezigen_gasten
        WHEN nv.afwezigen_gebruikers IS NOT NULL THEN
            get_user_names_from_uuids(nv.afwezigen_gebruikers)
        WHEN nv.afwezigen_gasten IS NOT NULL THEN
            nv.afwezigen_gasten
        ELSE
            ARRAY[]::TEXT[]
    END as afwezigen_combined
FROM notulen_versies nv;

COMMENT ON VIEW notulen_versies_with_participants IS 'Notulen versies view met gecombineerde participant lijsten';