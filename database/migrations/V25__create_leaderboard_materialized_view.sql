-- GECONSOLIDEERDE V25 - LEADERBOARD OPTIMALISATIE (MV)
--
-- DOEL: Vervangt de zware, dynamische 'unified_leaderboard' (V22) 
--       door een 'Materialized View' (MV).
--
-- WAAROM: Een dynamische view berekent de HELE ranglijst (joins, group by, sum, rank)
--       ELKE KEER dat iemand de pagina opvraagt. Dit is traag en schaalt niet.
--
--       Een Materialized View (MV) is een tabel die de resultaten opslaat.
--       De berekening gebeurt maar 1x (tijdens 'refresh').
--       Het opvragen is daarna bliksemsnel (een simpele SELECT).
--
-- IMPACT: Zeer hoge performancewinst op de leaderboard-feature.
--
-- VEREISTE ACTIE (in de Go-code):
-- 1. Pas de API-handler aan om `SELECT * FROM leaderboard_mv` te bevragen
--    in plaats van `SELECT * FROM unified_leaderboard`.
-- 2. Roep `SELECT refresh_leaderboard_mv();` aan nadat stappen zijn
--    geüpdatet, of draai het als een achtergrondtaak (bv. elke 5 min).
--

-- =====================================================
-- 1. MAAK DE MATERIALIZED VIEW AAN
-- =====================================================
-- Dit voert de zware query één keer uit en slaat het resultaat op.
CREATE MATERIALIZED VIEW IF NOT EXISTS leaderboard_mv AS
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
GROUP BY a.id, g.id -- Groepeer op de primary keys
ORDER BY total_score DESC, a.steps DESC;

COMMENT ON MATERIALIZED VIEW leaderboard_mv IS 'Snelle, pre-berekende cache van het unified leaderboard (V25)';

-- =====================================================
-- 2. MAAK EEN UNIEKE INDEX AAN
-- =====================================================
-- Dit is ESSENTIEEL om de view CONCURRENTLY (non-blocking) te kunnen verversen.
-- Zonder dit lockt de 'refresh' de hele view voor lezers.
CREATE UNIQUE INDEX IF NOT EXISTS idx_leaderboard_mv_participant_id 
ON leaderboard_mv(participant_id);

COMMENT ON INDEX idx_leaderboard_mv_participant_id IS 'Unieke index vereist voor CONCURRENTLY refresh van leaderboard_mv (V25)';

-- =====================================================
-- 3. MAAK DE REFRESH-FUNCTIE AAN
-- =====================================================
-- Dit is een helper-functie die je vanuit je Go-code kunt aanroepen.
CREATE OR REPLACE FUNCTION refresh_leaderboard_mv()
RETURNS void AS $$
BEGIN
    -- Ververs de view CONCURRENTLY (non-blocking)
    REFRESH MATERIALIZED VIEW CONCURRENTLY leaderboard_mv;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION refresh_leaderboard_mv IS 'Ververst het leaderboard (Materialized View) non-blocking (V25).';

-- =====================================================
-- 4. LOG INSTRUCTIES
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE '[V25] Materialized View `leaderboard_mv` is aangemaakt.';
    RAISE NOTICE '=== VOLGENDE STAPPEN ===';
    RAISE NOTICE '1. (Go Code) Pas je API aan om te lezen van `leaderboard_mv` i.p.v. `unified_leaderboard`.';
    RAISE NOTICE '2. (Go Code) Roep `SELECT refresh_leaderboard_mv();` aan nadat stappen zijn geüpdatet OF als een periodieke (cron) job.';
    RAISE NOTICE '3. (Optioneel) Na migratie, overweeg: `DROP VIEW IF EXISTS unified_leaderboard;` en `DROP VIEW IF EXISTS leaderboard_view;` om op te ruimen.';
END $$;