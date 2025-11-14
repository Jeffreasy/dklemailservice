# Leaderboard Fix - 2025-11-10

## Problem Identified

The application was experiencing errors when fetching leaderboard data:
```
ERROR: column a.afstand does not exist (SQLSTATE 42703)
```

This error was causing:
- Users getting logged out frequently
- Application going into offline mode
- Leaderboard functionality completely broken

## Root Cause

During the V28 migration ([`database/migrations/V28_CONSOLIDATED__rename_tables_participant_refactor.sql`](database/migrations/V28_CONSOLIDATED__rename_tables_participant_refactor.sql:1)), the `afstand` column (renamed to `distance_route`) and `steps` column were moved from the `participants` table to the `event_registrations` table as part of the participant refactor.

However, the leaderboard queries in [`repository/leaderboard_repository.go`](repository/leaderboard_repository.go:1) were not updated to reflect this schema change and continued trying to query `a.afstand` from the `participants` table alias.

## Solution Applied

Updated all queries in [`repository/leaderboard_repository.go`](repository/leaderboard_repository.go:1) to:

1. **Join with event_registrations table**: Added LEFT JOIN to connect participants with their event registrations
2. **Use active event function**: Filter event_registrations by active event using `(SELECT get_active_event())`
3. **Update column references**: Changed all references from:
   - `a.afstand` → `er.distance_route`
   - `a.steps` → `er.steps`
4. **Add COALESCE for safety**: Added `COALESCE(er.steps, 0)` to handle participants without event registrations

### Changes Made

#### GetLeaderboard Query
```sql
-- OLD (BROKEN):
FROM participants a
WHERE ...
SELECT a.afstand as route, a.steps

-- NEW (FIXED):  
FROM participants p
LEFT JOIN event_registrations er ON p.id = er.participant_id 
    AND er.event_id = (SELECT get_active_event())
WHERE ...
SELECT COALESCE(er.distance_route, 'Unknown') as route, 
       COALESCE(er.steps, 0) as steps
```

#### GetParticipantRank Query
Similarly updated to join with event_registrations and use proper column references.

## Files Modified

1. [`repository/leaderboard_repository.go`](repository/leaderboard_repository.go:1)
   - Line 30-90: GetLeaderboard method
   - Line 95-117: Count query in GetLeaderboard
   - Line 137-175: GetParticipantRank method  
   - Line 177-204: Above me query
   - Line 206-234: Below me query

## Additional Fix Required - V25 Migration

During testing, discovered that V25 migration was blocking deployment:

**Problem**: [`V25__create_leaderboard_materialized_view.sql`](database/migrations/V25__create_leaderboard_materialized_view.sql:1) tried to create a materialized view using the `participants` and `event_registrations` tables, but those tables don't exist until V28. The migration was trying to use schema from the future.

**Solution**: Completely replaced V25 migration content with a simple NOTICE statement. The materialized view approach is obsolete anyway since the leaderboard now uses direct queries via [`repository/leaderboard_repository.go`](repository/leaderboard_repository.go:1), which is:
- Always up-to-date (no refresh needed)
- More flexible for filtering and pagination
- Works properly with V28+ schema

## Testing Status

✅ **COMPLETED SUCCESSFULLY** (2025-11-10 21:13 CET)

### Test Results

1. ✅ All migrations pass successfully (V01-V34)
2. ✅ Service starts without errors
3. ✅ Leaderboard API responds correctly: `GET /api/leaderboard`
4. ✅ Returns proper JSON with participant data
5. ✅ Route field correctly shows distance (e.g., "2.5 KM", "6 KM", "10 KM", "15 KM")
6. ✅ Steps field present and working (currently 0 for all participants)
7. ✅ No "column a.afstand does not exist" errors in logs
8. ✅ Service health check passes

### Example Leaderboard Response
```json
{
    "entries": [
        {
            "id": "1ca80f61-f5c1-431f-b224-e6557150b65b",
            "naam": "Han van Doornik",
            "route": "2.5 KM",
            "steps": 0,
            "achievement_points": 0,
            "total_score": 0,
            "rank": 1,
            "badge_count": 0,
            "joined_at": "2025-03-30T08:07:45.334762Z"
        }
        // ... more entries
    ],
    "total_entries": 18,
    "current_page": 1,
    "total_pages": 1,
    "limit": 50
}
```

## Deployment Status

✅ **READY FOR PRODUCTION**
- All code fixes applied
- All migrations pass (V01-V34)
- Leaderboard functionality verified
- No breaking errors
- Service running stable

## Summary of Changes

### Code Changes
1. [`repository/leaderboard_repository.go`](repository/leaderboard_repository.go:1) - Updated all queries to use `event_registrations` table
2. [`database/migrations/V25__create_leaderboard_materialized_view.sql`](database/migrations/V25__create_leaderboard_materialized_view.sql:1) - Replaced with simple NOTICE (obsolete materialized view approach)

### Database Schema
- Leaderboard now properly queries from `participants` and `event_registrations` tables
- Uses `get_active_event()` function to filter for current event
- Handles participants without event registrations gracefully using COALESCE

### Impact
- ✅ Leaderboard functionality fully restored
- ✅ Users will no longer be logged out due to leaderboard errors
- ✅ Application stays online consistently
- ✅ All migrations execute successfully
- ✅ Ready for production deployment

**Fix verified and tested: 2025-11-10 21:13 CET**

## Technical Details

**Migration Context**: This fix addresses a breaking change introduced in V28 where the data model was refactored to separate:
- `participants` table: Person data only (naam, email, telefoon)
- `event_registrations` table: Event-specific data (status, rol, distance_route, steps)

The leaderboard needs to aggregate data from both tables, joining them via the active event.

## Impact

Once deployed and tested, this fix will:
- ✅ Restore leaderboard functionality
- ✅ Prevent users from being logged out unexpectedly  
- ✅ Keep application online consistently
- ✅ Allow proper tracking of steps and rankings per event