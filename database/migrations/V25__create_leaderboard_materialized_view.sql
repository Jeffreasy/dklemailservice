-- GECONSOLIDEERDE V25 - LEADERBOARD OPTIMALISATIE (SKIPPED)
--
-- ORIGINAL PURPOSE: Create materialized view for leaderboard performance
-- CURRENT STATUS: OBSOLETE - Skipped in favor of direct queries
--
-- WAAROM SKIPPED:
-- 1. V28 introduces new schema (participants + event_registrations tables)
-- 2. Materialized view would need to be recreated after V28 anyway
-- 3. Direct queries in leaderboard_repository.go are more flexible and always up-to-date
-- 4. No refresh mechanism needed with direct queries
--
-- THE FIX:
-- The leaderboard_repository.go now queries directly from participants and 
-- event_registrations tables with proper JOINs. This approach is:
-- - Always up-to-date (no refresh needed)
-- - Works with V28+ schema
-- - More flexible for filtering and pagination
-- - Simpler to maintain
--
-- This migration is kept for version continuity but does nothing.

DO $$
BEGIN
    RAISE NOTICE '[V25] SKIPPED - Materialized view approach replaced by direct queries';
    RAISE NOTICE '[V25] Leaderboard now uses repository pattern with direct SQL queries';
    RAISE NOTICE '[V25] See repository/leaderboard_repository.go for implementation';
END $$;