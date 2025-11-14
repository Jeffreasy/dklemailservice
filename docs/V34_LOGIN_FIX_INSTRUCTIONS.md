# V34 Login Fix: Refresh Token Foreign Key Resolution

**Status:** ✅ READY TO APPLY  
**Date:** 2025-11-10  
**Issue:** Login succeeds but fails on refresh token creation due to FK constraint

---

## Quick Fix (5 Minutes)

### The Problem

```
✅ Participant login succesvol (participant_id: 08e6be3e...)
❌ ERROR: insert or update on table "refresh_tokens" violates foreign key constraint
   "refresh_tokens_user_id_fkey"
```

**Root Cause:** `refresh_tokens.user_id` has a foreign key to `gebruikers.id`, but after V34 we're storing **participant IDs** which don't exist in the `gebruikers` table.

### The Solution

Remove the FK constraint and rename `user_id` → `owner_id` to clarify it can hold both types of IDs.

---

## Apply The Fix

### Local (Docker)

```bash
# Stop if running
docker-compose down

# Run the fix script
docker exec -i dkl-postgres psql -U postgres -d dkl_db \
  < database/scripts/fix_v34_refresh_tokens.sql

# Rebuild and start
docker-compose up --build -d

# Watch logs
docker-compose logs app --tail 50 --follow
```

### Production (Render)

```bash
# 1. Go to Render Dashboard
# 2. Click your service: dklemailservice
# 3. Click "Shell" tab
# 4. In the shell:

psql $DATABASE_URL

# 5. Copy-paste the ENTIRE content of:
#    database/scripts/fix_v34_refresh_tokens.sql

# 6. Exit
\q

# 7. Redeploy service (if needed)
```

---

## Verify It Works

### Test Login

```bash
# Mobile app: Login met diesbosje@hotmail.com / DKL2025!

# Expected logs in backend:
✅ Participant login succesvol
✅ Refresh token gegenereerd
✅ Participant tokens gegenereerd

# Should NOT see:
❌ ERROR: foreign key constraint violated
```

### Check Database

```sql
-- Verify column renamed
\d refresh_tokens

-- Should show:
-- owner_id | uuid | not null (NO user_id column!)

-- Check active tokens
SELECT 
    owner_id,
    LEFT(token, 20) || '...' as token,
    expires_at,
    is_revoked
FROM refresh_tokens
WHERE is_revoked = false
LIMIT 5;
```

---

## What Changed

### Database Schema

**Before:**
```sql
refresh_tokens.user_id UUID NOT NULL
  REFERENCES gebruikers(id)
```

**After:**
```sql
refresh_tokens.owner_id UUID NOT NULL
  -- No FK constraint
  -- Can reference EITHER gebruikers.id OR participants.id
```

### Backend Code

**Model Updated:** [`models/refresh_token.go`](../models/refresh_token.go:8)
```go
// Before
UserID string `json:"user_id" gorm:"not null;type:uuid;index"`

// After
OwnerID string `json:"owner_id" gorm:"column:owner_id;not null;type:uuid;index"`
```

**Service Updated:** [`services/auth_service.go`](../services/auth_service.go:510)
```go
// Before
refreshToken := &models.RefreshToken{
    UserID: userID,  // Would fail for participant IDs
    ...
}

// After  
refreshToken := &models.RefreshToken{
    OwnerID: userID,  // Works for both gebruiker and participant IDs
    ...
}
```

---

## Rollback (If Needed)

⚠️ **WARNING:** Rollback deletes ALL participant refresh tokens!

```bash
# Local
docker exec -i dkl-postgres psql -U postgres -d dkl_db \
  < database/scripts/rollback_v34_refresh_tokens.sql

# Production (via Render Shell)
psql $DATABASE_URL < rollback_v34_refresh_tokens.sql
```

---

## Files Changed

### SQL Scripts (New)
- ✅ [`database/scripts/fix_v34_refresh_tokens.sql`](../database/scripts/fix_v34_refresh_tokens.sql) - Main fix
- ✅ [`database/scripts/rollback_v34_refresh_tokens.sql`](../database/scripts/rollback_v34_refresh_tokens.sql) - Rollback

### Backend Code (Modified)
- ✅ [`models/refresh_token.go`](../models/refresh_token.go) - UserID → OwnerID
- ✅ [`services/auth_service.go`](../services/auth_service.go) - Updated token generation & refresh

### Documentation (New)
- ✅ [`docs/V34_LOGIN_FIX_INSTRUCTIONS.md`](V34_LOGIN_FIX_INSTRUCTIONS.md) - This file
- ✅ [`database/scripts/README.md`](../database/scripts/README.md) - Updated with fix section

---

## FAQ

**Q: Will this break existing gebruiker logins?**  
A: No! The code handles both cases - gebruiker tokens work exactly as before.

**Q: Do I need to update the frontend?**  
A: No! Frontend is unchanged - the token flow remains identical.

**Q: What happens to existing refresh tokens?**  
A: They keep working! The database just renames the column, existing data is preserved.

**Q: Can I still use gebruikers table?**  
A: Yes! Admin/staff accounts continue using `gebruikers` table normally.

---

## Success Checklist

After applying the fix:

- [ ] SQL script ran without errors
- [ ] `\d refresh_tokens` shows `owner_id` column (not `user_id`)
- [ ] No FK constraint `refresh_tokens_user_id_fkey` in `\d refresh_tokens`
- [ ] Backend rebuilt and restarted
- [ ] Participant can login successfully
- [ ] No FK constraint errors in logs
- [ ] Refresh token is generated and saved
- [ ] Mobile app receives token and can access protected routes

---

## Related Documentation

- [`BACKEND_FIX_V34_LOGIN_ISSUE.md`](BACKEND_FIX_V34_LOGIN_ISSUE.md) - Detailed analysis
- [`V30_DUAL_REGISTRATION_SYSTEM.md`](V30_DUAL_REGISTRATION_SYSTEM.md) - System overview
- [`docs/migrations/V34_BREAKING_CHANGES.md`](migrations/V34_BREAKING_CHANGES.md) - All V34 changes

---

**Status:** ✅ Fix ready to apply  
**Impact:** High (blocks all participant logins)  
**Risk:** Low (backward compatible, rollback available)  
**Time:** ~5 minutes to apply