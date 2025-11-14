-- Fix V25 Leaderboard Materialized View
-- This fixes the column references after V34 removed legacy columns from participants table
-- Steps and distance_route are now in event_registrations table

-- Drop the old materialized view
DROP MATERIALIZED VIEW IF EXISTS leaderboard_mv CASCADE;

-- Recreate with correct column references
CREATE MATERIALIZED VIEW leaderboard_mv AS
SELECT 
    p.id as participant_id,
    p.naam as display_name,
    p.email,
    er.distance_route as route,
    COALESCE(er.steps, 0) as steps,
    COALESCE(SUM(b.points), 0) as achievement_points,
    COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0) as total_score,
    RANK() OVER (ORDER BY (COALESCE(er.steps, 0) + COALESCE(SUM(b.points), 0)) DESC) as rank,
    COUNT(pa.id) as badge_count,
    p.created_at as joined_at,
    p.gebruiker_id,
    CASE 
        WHEN p.gebruiker_id IS NOT NULL THEN true 
        ELSE false 
    END as has_user_account,
    g.naam as user_naam,
    g.is_actief as user_is_actief
FROM participants p
LEFT JOIN event_registrations er ON p.id = er.participant_id 
    AND er.event_id = (SELECT get_active_event())
LEFT JOIN participant_achievements pa ON p.id = pa.participant_id
LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
LEFT JOIN gebruikers g ON p.gebruiker_id = g.id
GROUP BY p.id, er.id, g.id
ORDER BY total_score DESC, steps DESC;

COMMENT ON MATERIALIZED VIEW leaderboard_mv IS 'Snelle, pre-berekende cache van het unified leaderboard (V25, fixed in V34+)';

-- Recreate the unique index
CREATE UNIQUE INDEX idx_leaderboard_mv_participant_id 
ON leaderboard_mv(participant_id);

COMMENT ON INDEX idx_leaderboard_mv_participant_id IS 'Unieke index vereist voor CONCURRENTLY refresh van leaderboard_mv (V25)';

-- Verify the fix worked
SELECT 
    COUNT(*) as total_participants,
    COUNT(CASE WHEN steps > 0 THEN 1 END) as participants_with_steps
FROM leaderboard_mv;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '✅ Leaderboard materialized view recreated successfully!';
    RAISE NOTICE 'Now using event_registrations.steps and event_registrations.distance_route';
END $$;