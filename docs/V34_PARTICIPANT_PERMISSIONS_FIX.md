# V34 Participant Permissions Fix

## Issue Summary
After V34 login fix, participants could authenticate successfully but received `permissions_count=0` errors when accessing protected endpoints like `/api/participant/dashboard`.

## Root Cause Analysis

### The Problem
1. **Login Handler** ([`auth_service.go:238`](../services/auth_service.go:238))
   - Generated JWT with hardcoded `participant_user` role
   - Token contained: `Roles: []string{"participant_user"}`
   - This worked for authentication but NOT for authorization

2. **Permission Middleware** ([`permission_service.go:98`](../services/permission_service.go:98))
   - Called `userRoleRepo.GetUserPermissions(userID)`
   - Query joined: `user_roles` → `roles` → `role_permissions` → `permissions` → **`gebruikers`**
   - **CRITICAL**: Participant ID exists in `participants` table, NOT `gebruikers` table
   - Result: Query returned 0 permissions

3. **Mismatch**
   - JWT claims said: "You have participant_user role"
   - Database said: "No permissions found for this user_id"

## Solution Implemented

### Modified [`HasPermission()`](../services/permission_service.go:85) in PermissionService

Added participant detection before RBAC database query:

```go
// V34 FIX: Check if this is a participant ID
if s.participantRepo != nil {
    participant, err := s.participantRepo.GetByID(ctx, userID)
    if err == nil && participant != nil {
        // Use participant permissions instead of RBAC
        hasPermission := s.checkParticipantPermission(participant, resource, action)
        return hasPermission
    }
}

// Fallback to RBAC for regular gebruikers
permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID)
// ... rest of RBAC logic
```

### Added [`checkParticipantPermission()`](../services/permission_service.go:172)

Hardcoded permissions for participants with full accounts:

```go
participantPermissions := map[string][]string{
    "steps":       {"view_own", "create"},
    "participant": {"view_own", "update_own"},
    "events":      {"view", "register"},
    "app":         {"access"},
}
```

**Requirements:**
- Must have `account_type = "full"`
- Must have `has_app_access = true`

## Verification

### Before Fix
```json
{
  "niveau": "WARN",
  "bericht": "Permission denied - permissions_count=0",
  "resource": "steps",
  "action": "view_own"
}
{
  "niveau": "WARN",
  "bericht": "Toegang geweigerd - onvoldoende rechten"
}
```

### After Fix
```json
{
  "niveau": "INFO",
  "bericht": "HasPermission called",
  "participant_repo_available": true
}
{
  "niveau": "INFO",
  "bericht": "Token gevalideerd"
}
// NO permission denied errors!
```

## Impact

✅ **FIXED**: Participant authentication and authorization
✅ **FIXED**: Permission middleware now recognizes participants
✅ **FIXED**: Participants can access protected endpoints

The original error "geen deelnemersregistratie gevonden" is a **different issue** related to dashboard data requirements, NOT permissions.

## Technical Details

### Files Modified
- [`services/permission_service.go`](../services/permission_service.go) - Added participant permission logic

### Dependencies
- Participant repository must be injected into PermissionService
- Already configured in [`services/factory.go:98`](../services/factory.go:98)

## Testing Results

```
✅ Participant login successful
✅ Tokens generated without FK errors
✅ Token validation passing
✅ Permission middleware allowing access
✅ No more "permissions_count=0" errors
```

## Next Steps (Optional)

1. Consider storing participant permissions in database for consistency
2. Add RBAC role for participants if needed
3. Review other endpoints that may need participant permission support

## Version
- **Fix Date**: 2025-11-10
- **Version**: V34 Post-Login Fix
- **Status**: ✅ RESOLVED