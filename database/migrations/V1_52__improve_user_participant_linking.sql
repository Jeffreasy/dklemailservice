-- Migration: V1_52__improve_user_participant_linking.sql
-- Description: Verbetert de koppeling tussen gebruikers en aanmeldingen voor gamification
-- Date: 2026-01-02

-- =====================================================
-- 1. VOEG FOREIGN KEY TOE VOOR GEBRUIKER_ID
-- =====================================================
-- Zorg dat GebruikerID correct verwijst naar gebruikers tabel
ALTER TABLE aanmeldingen 
ADD CONSTRAINT fk_aanmelding_gebruiker 
FOREIGN KEY (gebruiker_id) REFERENCES gebruikers(id) ON DELETE SET NULL;

COMMENT ON COLUMN aanmeldingen.gebruiker_id IS 'Optionele link naar gebruikersaccount - voor participants die ook een user account hebben';

-- Index voor snelle lookups
CREATE INDEX IF NOT EXISTS idx_aanmeldingen_gebruiker_id ON aanmeldingen(gebruiker_id);

-- =====================================================
-- 2. VIEW VOOR GEBRUIKER-PARTICIPANT MAPPING
-- =====================================================
-- Handige view die laat zien welke users ook participants zijn
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
-- 3. FUNCTIE OM PARTICIPANT ID TE VINDEN VOOR USER
-- =====================================================
-- Hulpfunctie om participant_id te vinden voor een gebruiker
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
-- 4. FUNCTIE OM GEBRUIKER ID TE VINDEN VOOR PARTICIPANT
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
-- 5. VIEW VOOR UNIFIED USER LEADERBOARD
-- =====================================================
-- Leaderboard die ook user info toont waar beschikbaar
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
-- 6. TRIGGER VOOR AUTOMATISCHE EMAIL SYNC
-- =====================================================
-- Optioneel: Als gebruiker email wijzigt, update aanmelding email
CREATE OR REPLACE FUNCTION sync_user_email_to_aanmelding()
RETURNS TRIGGER AS $$
BEGIN
    -- Update aanmelding email als gebruiker email wijzigt
    UPDATE aanmeldingen
    SET email = NEW.email
    WHERE gebruiker_id = NEW.id
    AND email != NEW.email; -- Alleen als email anders is
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger aanmaken (alleen als je automatische sync wilt)
-- UNCOMMENT onderstaande regels als je dit wilt:
/*
CREATE TRIGGER sync_user_email_trigger
    AFTER UPDATE OF email ON gebruikers
    FOR EACH ROW
    WHEN (OLD.email IS DISTINCT FROM NEW.email)
    EXECUTE FUNCTION sync_user_email_to_aanmelding();
*/

-- =====================================================
-- 7. HELPER VIEW: USERS WITHOUT PARTICIPATION
-- =====================================================
-- Toon users die nog geen participant zijn (nuttig voor admin UI)
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
WHERE a.id IS NULL -- Heeft geen aanmelding
GROUP BY g.id, g.naam, g.email, g.rol, g.is_actief, g.created_at;

COMMENT ON VIEW users_without_participation IS 'Users die nog geen participant/aanmelding hebben';

-- =====================================================
-- 8. DOCUMENTATIE COMMENTS
-- =====================================================
COMMENT ON CONSTRAINT fk_aanmelding_gebruiker ON aanmeldingen IS 
'Foreign key naar gebruikers - een participant kan een gebruikersaccount hebben voor inloggen';

-- =====================================================
-- MIGRATION NOTES
-- =====================================================
/*
GEBRUIKERS vs PARTICIPANTS:

1. TWEE TYPES PARTICIPANTS:
   - Participants ZONDER account: Alleen entry in aanmeldingen (gebruiker_id = NULL)
   - Participants MET account: Entry in beide tabellen, gekoppeld via gebruiker_id

2. ADMIN/STAFF PARTICIPATIE:
   - Admin/Staff die ALLEEN beheren: Alleen entry in gebruikers (geen aanmelding)
   - Admin/Staff die OOK deelnemen: Entry in beide tabellen

3. GAMIFICATION:
   - Badges/Achievements zijn altijd gekoppeld aan aanmeldingen (participant_id)
   - Een gebruiker moet een aanmelding hebben om achievements te krijgen
   - Check via: get_participant_id_for_user(user_id)

4. AUTHENTICATION:
   - Login is via gebruikers tabel
   - Na login kan je via gebruiker_id in aanmeldingen je participant data vinden
   - Handler moet beide IDs kunnen gebruiken (userID voor auth, participantID voor gamification)

5. BEST PRACTICES:
   - Gebruik user_participant_mapping view voor UI die beide moet tonen
   - Gebruik get_participant_id_for_user() in handlers die participant_id nodig hebben
   - Bij stappen updates: check eerst of user een participant is

VOORBEELD QUERIES:

-- Check of user een participant is:
SELECT get_participant_id_for_user('user-uuid-here');

-- Haal leaderboard op met user info:
SELECT * FROM unified_leaderboard LIMIT 10;

-- Vind alle users die kunnen deelnemen maar nog niet doen:
SELECT * FROM users_without_participation;
*/

-- =====================================================
-- MIGRATION COMPLETE
-- =====================================================