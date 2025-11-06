-- GECONSOLIDEERDE V28 - DE GROTE REFACTOR (PERSOON vs. DEELNAME)
--
-- DOEL: Maakt de applicatie multi-event capabel (bv. DKL 2025, DKL 2026).
--       Scheidt 'wie een persoon is' van 'aan welk event ze meedoen'.
--
-- HOE:
-- 1. 'aanmeldingen' -> 'participants' (permanent persoon-record)
-- 2. 'event_participants' (V23) -> 'event_registrations' (deelname-record)
-- 3. Verplaatst alle event-specifieke data (steps, status, afstand)
--    van 'participants' naar 'event_registrations'.
--
-- IMPACT: ENORME BREAKING CHANGE
-- Vereist een volledige refactor van je Go-code (GORM-modellen en services)
-- die 'aanmeldingen' gebruikt voor event-data.
--

-- =====================================================
-- 1. HERNOEM DE KERN-TABELLEN VOOR DE DUIDELIJKHEID
-- =====================================================

-- 'aanmeldingen' (V1) wordt de 'participants' tabel.
ALTER TABLE IF EXISTS aanmeldingen RENAME TO participants;

-- 'event_participants' (V23) wordt de 'event_registrations' tabel.
ALTER TABLE IF EXISTS event_participants RENAME TO event_registrations;

-- Hernoem de bijbehorende view ook
DROP VIEW IF EXISTS event_participants_view; -- van V23
ALTER VIEW IF EXISTS user_participant_mapping RENAME TO user_participant_link; -- van V22

-- =====================================================
-- 2. VERPLAATS KOLOMMEN VAN 'participants' -> 'event_registrations'
-- =====================================================

-- We hebben deze kolommen nodig in de 'deelname'-tabel.
-- (We gebruiken de genormaliseerde kolomnamen uit V26 en V27)
ALTER TABLE event_registrations
    ADD COLUMN IF NOT EXISTS participant_role_name TEXT REFERENCES participant_roles(name) ON UPDATE CASCADE,
    ADD COLUMN IF NOT EXISTS distance_route TEXT REFERENCES distances(route) ON UPDATE CASCADE,
    ADD COLUMN IF NOT EXISTS status TEXT REFERENCES registration_status_types(status) ON UPDATE CASCADE,
    ADD COLUMN IF NOT EXISTS bijzonderheden TEXT,
    ADD COLUMN IF NOT EXISTS terms BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS notities TEXT,
    ADD COLUMN IF NOT EXISTS antwoorden_count INTEGER DEFAULT 0;

-- Hernoem de 'steps' kolom (uit V23 'current_steps') voor consistentie
ALTER TABLE event_registrations RENAME COLUMN IF EXISTS current_steps TO steps;

-- =====================================================
-- 3. MIGREER DE DATA
-- =====================================================
-- We moeten alle data uit de (oude) 'aanmeldingen' tabel
-- kopiëren naar de 'event_registrations' tabel, en
-- koppelen aan het HUIDIGE actieve event.

DO $$
DECLARE
    v_active_event_id UUID;
BEGIN
    -- 1. Vind het actieve event (functie uit V23)
    SELECT get_active_event() INTO v_active_event_id;
    
    IF v_active_event_id IS NULL THEN
        RAISE WARNING '[V28] Geen actief event gevonden (get_active_event() gaf NULL). Data migratie overgeslagen.';
        RETURN;
    END IF;

    RAISE NOTICE '[V28] Actief event ID: %. Migreren van data...', v_active_event_id;

    -- 2. Kopieer de data
    INSERT INTO event_registrations (
        event_id, 
        participant_id, 
        registered_at, 
        tracking_status,
        steps,
        status,
        distance_route,
        participant_role_name,
        bijzonderheden,
        terms,
        notities,
        antwoorden_count
    )
    SELECT
        v_active_event_id,      -- Koppel aan het actieve event
        p.id,                   -- De Participant ID
        p.created_at,           -- De originele aanmelddatum
        'registered',           -- Standaard tracking status
        p.steps,                -- De stappen (van de oude tabel)
        p.status,               -- De status (van de oude tabel)
        p.distance_route,       -- De afstand (van V26)
        p.participant_role_name,-- De rol (van V26)
        p.bijzonderheden,
        p.terms,
        p.notities,
        p.antwoorden_count
    FROM participants p
    -- Voorkom duplicaten als het script opnieuw draait
    WHERE NOT EXISTS (
        SELECT 1 FROM event_registrations er
        WHERE er.participant_id = p.id AND er.event_id = v_active_event_id
    );

    RAISE NOTICE '[V28] Data migratie voltooid.';
END $$;

-- =====================================================
-- 4. OPRAKELEN (Optioneel, pas uitvoeren na Go-code update)
-- =====================================================
-- Nadat je Go-code 100% is gemigreerd om 'event_registrations'
-- te gebruiken, kun je deze kolommen veilig uit de 'participants' tabel verwijderen.
-- Dit maakt de 'participants' tabel een schoon, permanent record.

/*
ALTER TABLE participants
    DROP COLUMN IF EXISTS steps,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS distance_route,
    DROP COLUMN IF EXISTS participant_role_name,
    DROP COLUMN IF EXISTS bijzonderheden,
    DROP COLUMN IF EXISTS terms,
    DROP COLUMN IF EXISTS notities,
    DROP COLUMN IF EXISTS antwoorden_count,
    DROP COLUMN IF EXISTS rol,        -- Legacy (V1)
    DROP COLUMN IF EXISTS afstand;    -- Legacy (V1)
*/

-- =====================================================
-- 5. FIX DE 'aanmelding_antwoorden' RELATIE
-- =====================================================
-- Deze tabel (V1) verwijst nog naar de oude 'aanmeldingen' tabel.
-- We moeten dit hernoemen.
ALTER TABLE IF EXISTS aanmelding_antwoorden RENAME TO participant_antwoorden;
ALTER TABLE participant_antwoorden
    DROP CONSTRAINT IF EXISTS fk_aanmelding_antwoorden_aanmelding_id,
    ADD CONSTRAINT fk_participant_antwoorden_participant_id
    FOREIGN KEY (aanmelding_id) REFERENCES participants(id) ON DELETE CASCADE;
COMMENT ON TABLE participant_antwoorden IS 'Antwoorden op participants (voorheen aanmeldingen) (V1)';

-- =====================================================
-- 6. LOG INSTRUCTIES
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE '[V28] De 'Persoon vs. Deelname' refactor is voorbereid.';
    RAISE NOTICE '=== ZEER GROTE BREAKING CHANGE ===';
    RAISE NOTICE '1. (Go Code) Hernoem je 'Aanmelding' GORM-model naar 'Participant'.';
    RAISE NOTICE '2. (Go Code) Hernoem je 'EventParticipant' (V23) GORM-model naar 'EventRegistration'.';
    RAISE NOTICE '3. (Go Code) Verplaats alle logica voor 'steps', 'status', 'afstand' etc.';
    RAISE NOTICE '   van de 'ParticipantService' naar de 'EventRegistrationService'.';
    RAISE NOTICE '4. (Optioneel) Voer Stap 4 (DROP COLUMNs) uit om de database definitief op te schonen.';
END $$;