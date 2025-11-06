-- GECONSOLIDEERDE V22 - VERBETER USER/PARTICIPANT LINKING
-- Logica van V1_52.
-- OPTIMALISATIE: Sectie 1 (toevoegen van FK) is verwijderd.
-- Onze V12__add_gebruiker_id_to_aanmeldingen.sql heeft deze FK al aangemaakt.
-- De rest van de views en functies is nuttig en wordt behouden.

-- =====================================================
-- 1. VIEW VOOR GEBRUIKER-PARTICIPANT MAPPING
-- =====================================================
CREATE OR REPLACE VIEW user_participant_mapping AS
SELECT 
    g.id as gebruiker_id,
    g.naam as gebruiker_naam,
    g.email as gebruiker_email,
    g.rol as legacy_rol,
    a.id as aanmelding_id,
    a.naam as participant_naam,
    a.email as participant_email,
    a.afstand as route,
    a.steps,
    a.status as participant_status,
    CASE 
        WHEN a.id IS NOT NULL THEN true 
        ELSE false 
    END as is_participant
FROM gebruikers g
LEFT JOIN aanmeldingen a ON g.id = a.gebruiker_id;

COMMENT ON VIEW user_participant_mapping IS 'Toont welke users ook participants zijn (hebben aanmelding)';

-- =====================================================
-- 2. FUNCTIE OM PARTICIPANT ID TE VINDEN VOOR USER
-- =====================================================
CREATE OR REPLACE FUNCTION get_participant_id_for_user(p_user_id UUID)
RETURNS UUID AS $$
DECLARE
    v_participant_id UUID;
BEGIN
    SELECT id INTO v_participant_id
    FROM aanmeldingen
    WHERE gebruiker_id = p_user_id
    LIMIT 1;
    
    RETURN v_participant_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_participant_id_for_user IS 'Vindt participant_id voor een gebruiker (retourneert NULL als user geen participant is)';

-- =====================================================
-- 3. FUNCTIE OM GEBRUIKER ID TE VINDEN VOOR PARTICIPANT
-- =====================================================
CREATE OR REPLACE FUNCTION get_user_id_for_participant(p_participant_id UUID)
RETURNS UUID AS $$
DECLARE
    v_user_id UUID;
BEGIN
    SELECT gebruiker_id INTO v_user_id
    FROM aanmeldingen
    WHERE id = p_participant_id;
    
    RETURN v_user_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_id_for_participant IS 'Vindt user_id voor een participant (retourneert NULL als participant geen user account heeft)';

-- =====================================================
-- 4. VIEW VOOR UNIFIED USER LEADERBOARD
-- =====================================================
CREATE OR REPLACE VIEW unified_leaderboard AS
SELECT 
    a.id as participant_id,
    a.naam as display_name,
    a.email,
    a.afstand as route,
    a.steps,
    COALESCE(SUM(b.points), 0) as achievement_points,
    a.steps + COALESCE(SUM(b.points), 0) as total_score,
    RANK() OVER (ORDER BY (a.steps + COALESCE(SUM(b.points), 0)) DESC) as rank,
    COUNT(pa.id) as badge_count,
    a.created_at as joined_at,
    a.gebruiker_id,
    CASE 
        WHEN a.gebruiker_id IS NOT NULL THEN true 
        ELSE false 
    END as has_user_account,
    g.naam as user_naam,
    g.is_actief as user_is_actief
FROM aanmeldingen a
LEFT JOIN participant_achievements pa ON a.id = pa.participant_id
LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
LEFT JOIN gebruikers g ON a.gebruiker_id = g.id
GROUP BY a.id, a.naam, a.email, a.afstand, a.steps, a.created_at, a.gebruiker_id, g.naam, g.is_actief
ORDER BY total_score DESC, a.steps DESC;

COMMENT ON VIEW unified_leaderboard IS 'Leaderboard met gebruikers info waar beschikbaar';

-- =====================================================
-- 5. TRIGGER VOOR AUTOMATISCHE EMAIL SYNC
-- =====================================================
CREATE OR REPLACE FUNCTION sync_user_email_to_aanmelding()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE aanmeldingen
    SET email = NEW.email
    WHERE gebruiker_id = NEW.id
    AND email != NEW.email;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- De trigger zelf is uit-commentarieerd, zoals in het originele script
/*
DROP TRIGGER IF EXISTS sync_user_email_trigger ON gebruikers;
CREATE TRIGGER sync_user_email_trigger
    AFTER UPDATE OF email ON gebruikers
    FOR EACH ROW
    WHEN (OLD.email IS DISTINCT FROM NEW.email)
    EXECUTE FUNCTION sync_user_email_to_aanmelding();
*/

-- =====================================================
-- 6. HELPER VIEW: USERS WITHOUT PARTICIPATION
-- =====================================================
CREATE OR REPLACE VIEW users_without_participation AS
SELECT 
    g.id,
    g.naam,
    g.email,
    g.rol as legacy_rol,
    g.is_actief,
    g.created_at,
    array_agg(r.name) as actual_roles
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN aanmeldingen a ON g.id = a.gebruiker_id
WHERE a.id IS NULL
GROUP BY g.id, g.naam, g.email, g.rol, g.is_actief, g.created_at;

COMMENT ON VIEW users_without_participation IS 'Users die nog geen participant/aanmelding hebben';

-- =====================================================
-- 7. ADD COMMENT TO CONSTRAINT (IDEMPOTENT)
-- =====================================================
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'aanmeldingen_gebruiker_id_fkey' -- Naam aangepast naar die van V12
        AND table_name = 'aanmeldingen'
    ) THEN
        COMMENT ON CONSTRAINT aanmeldingen_gebruiker_id_fkey ON aanmeldingen IS 
        'Foreign key naar gebruikers - een participant kan een gebruikersaccount hebben voor inloggen (V12, V22)';
    END IF;
END $$;