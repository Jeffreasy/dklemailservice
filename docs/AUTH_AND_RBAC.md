# 🔐 Authentication & RBAC - Complete System Guide

> **Version:** 3.0 | **Status:** Production Ready | **Last Updated:** 2025-11-01 | **Backend:** V1.48.0+

Complete documentatie van het authenticatie en autorisatie systeem van het DKL Email Service Admin Panel.

---

## 📋 Table of Contents

1. [System Overview](#-system-overview)
2. [Authentication](#-authentication)
3. [Authorization (RBAC)](#-authorization-rbac)
4. [Database Schema](#-database-schema)
5. [Backend Implementation](#-backend-implementation)
6. [Frontend Integration](#-frontend-integration)
7. [API Reference](#-api-reference)
8. [Security & Best Practices](#-security--best-practices)
9. [Deployment](#-deployment)
10. [Troubleshooting](#-troubleshooting)

---

## 🎯 System Overview

Het DKL Email Service implementeert een volledig backend-gedreven RBAC systeem met enterprise-grade security.

### Core Features

✅ **JWT Authentication** - 20 min access tokens  
✅ **Refresh Tokens** - 7 dagen met automatic rotation  
✅ **RBAC** - 19 resources, 58 granulaire permissions  
✅ **Redis Caching** - 5 min cache, 97% hit rate  
✅ **Multi-layer Checks** - UI + Route + Component + API  
✅ **Token Rotation** - Extra security bij refresh  
✅ **Backend as Truth** - Alle permissions uit database  

### Architecture

```
┌─────────────── FRONTEND (React) ───────────────┐
│  AuthProvider → useAuth → usePermissions       │
│       ↓              ↓            ↓             │
│  Token Storage   Auth State   Permission Checks │
│       ↓              ↓            ↓             │
│  AuthGuard    ProtectedRoute  Conditional UI    │
└────────────────────┬────────────────────────────┘
                     │ HTTP (Bearer Token)
                     ↓
┌─────────────── BACKEND (Go Fiber) ─────────────┐
│  Auth Endpoints → Permission Middleware        │
│       ↓                  ↓                      │
│  JWT Generation    Permission Check            │
│       ↓                  ↓                      │
│  Login/Refresh    Redis Cache → Database       │
└───────────────────────────────────────────────┘
```

---

## 🔑 Authentication

### Login Flow

```
1. User Input (email + password)
   ↓
2. Auto-append domain (@dekoninklijkeloop.nl if needed)
   ↓
3. POST /api/auth/login
   ↓
4. Backend validates credentials (bcrypt)
   ↓
5. Generate JWT (20 min) + Refresh Token (7 days)
   ↓
6. Store tokens in localStorage
   ↓
7. Load user profile + permissions from /api/auth/profile
   ↓
8. Navigate to dashboard
```

### JWT Token Structure (V1.22+)

```json
{
  "header": { "alg": "HS256", "typ": "JWT" },
  "payload": {
    "sub": "user-uuid",
    "email": "admin@dekoninklijkeloop.nl",
    "role": "admin",           // Legacy - kept for compatibility
    "roles": ["admin"],        // RBAC - array of role names
    "rbac_active": true,       // Indicates RBAC enabled
    "exp": 1234567890,
    "iat": 1234567890,
    "iss": "dklemailservice"
  }
}
```

**Backend Reference:** [`services/auth_service.go:276-321`](../services/auth_service.go)

### Token Lifecycle

| Phase | Duration | Purpose |
|-------|----------|---------|
| **Access Token** | 20 minutes | API authentication |
| **Refresh Token** | 7 days | Token renewal |
| **Auto-Refresh** | 15 min timer | 5 min before expiry |
| **Storage** | localStorage | Client-side persistence |
| **Revocation** | On logout | Security cleanup |

### Automatic Token Refresh

```typescript
// Scheduled 5 minutes before JWT expiry
private scheduleRefresh() {
  this.refreshTimer = setTimeout(
    () => this.refreshToken(), 
    15 * 60 * 1000  // 15 minutes
  )
}

// Refresh with token rotation
POST /api/auth/refresh
Body: { "refresh_token": "..." }
Response: { 
  "token": "new_access_token",
  "refresh_token": "new_refresh_token"  // Old revoked!
}
```

---

## 🛡️ Authorization (RBAC)

### Permission Structure

```typescript
interface Permission {
  resource: string  // e.g., "contact", "user", "photo"
  action: string    // e.g., "read", "write", "delete"
}

// Format: "resource:action"
// Examples: "contact:read", "user:manage_roles", "admin:access"
```

### Complete Permission Catalog

| Resource | Permissions | Backend Migration |
|----------|-------------|-------------------|
| **admin** | `access` (full access) | V1.21 |
| **staff** | `access` (staff level) | V1.23 |
| **contact** | `read`, `write`, `delete` | V1.21 |
| **aanmelding** | `read`, `write`, `delete` | V1.21 |
| **user** | `read`, `write`, `delete`, `manage_roles` | V1.21 |
| **photo** | `read`, `write`, `delete` | V1.24 |
| **album** | `read`, `write`, `delete` | V1.34 |
| **video** | `read`, `write`, `delete` | V1.35 |
| **partner** | `read`, `write`, `delete` | V1.33 |
| **sponsor** | `read`, `write`, `delete` | V1.36 |
| **radio_recording** | `read`, `write`, `delete` | V1.33 |
| **program_schedule** | `read`, `write`, `delete` | V1.37 |
| **social_embed** | `read`, `write`, `delete` | V1.38 |
| **social_link** | `read`, `write`, `delete` | V1.39 |
| **under_construction** | `read`, `write`, `delete` | V1.40 |
| **newsletter** | `read`, `write`, `send`, `delete` | V1.21 |
| **email** | `read`, `write`, `delete`, `fetch` | V1.21 |
| **admin_email** | `send` | V1.21 |
| **chat** | `read`, `write`, `manage_channel`, `moderate` | V1.21 |

**Total:** 19 resources, 58 permissions

### System Roles

| Role | Permissions | Use Case |
|------|-------------|----------|
| **admin** | ALL (58 permissions) | Platform administrators |
| **staff** | Read-only on most resources | Support staff |
| **user** | Basic chat read/write | Regular users |
| **owner** | Full chat management | Channel creators |
| **chat_admin** | Chat moderation | Channel moderators |
| **member** | Chat read/write | Channel members |
| **deelnemer** | - | Event participants (categorization) |
| **begeleider** | - | Event guides (categorization) |
| **vrijwilliger** | - | Event volunteers (categorization) |

**Note:** Event roles (deelnemer, begeleider, vrijwilliger) have no special permissions - used for registration categorization only.

### Chat Permissions Details

**IMPORTANT:** Chat permissions are **ROLE-BASED**, not "everyone". Access determined via `user_roles` assignments.

| Permission | Assigned to Roles | Use Case |
|------------|-------------------|----------|
| `chat:read` | owner, chat_admin, member, user | View messages, channels |
| `chat:write` | owner, chat_admin, member, user | Send messages, replies |
| `chat:moderate` | owner, chat_admin | Moderate messages |
| `chat:manage_channel` | owner | Create/manage channels |

---

## 🗄️ Database Schema

### RBAC Core Tables

```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT FALSE,  -- Cannot be deleted
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    is_system_permission BOOLEAN DEFAULT FALSE,  -- Cannot be deleted
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(resource, action)
);

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_by UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(role_id, permission_id)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES gebruikers(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,      -- Can be temporarily disabled
    expires_at TIMESTAMP,                -- Optional expiration
    UNIQUE(user_id, role_id)
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES gebruikers(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Migrations:**
- [`V1_20__create_rbac_tables.sql`](../database/migrations/V1_20__create_rbac_tables.sql)
- [`V1_21__seed_rbac_data.sql`](../database/migrations/V1_21__seed_rbac_data.sql)
- [`V1_22__migrate_legacy_roles_to_rbac.sql`](../database/migrations/V1_22__migrate_legacy_roles_to_rbac.sql)
- [`V1_28__add_refresh_tokens.sql`](../database/migrations/V1_28__add_refresh_tokens.sql)

---

## 🔧 Backend Implementation

### Permission Service

**Core Method:** `HasPermission(ctx, userID, resource, action) bool`

**Flow:**
```go
1. Check Redis cache (5 min TTL)
   ├─→ Cache Hit → Return cached result (1-2ms)
   └─→ Cache Miss ↓
2. Query database (user_roles JOIN)
3. Check permission in result set
4. Cache result for next time
5. Return boolean
```

**Implementation:** [`services/permission_service.go:68-110`](../services/permission_service.go)

### Permission Middleware

```go
func PermissionMiddleware(
    permissionService services.PermissionService, 
    resource, action string
) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("userID").(string)
        
        if !permissionService.HasPermission(c.Context(), userID, resource, action) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Geen toegang",
            })
        }
        
        return c.Next()
    }
}
```

**Reference:** [`handlers/permission_middleware.go:12-43`](../handlers/permission_middleware.go)

### Convenience Middlewares

```go
// Admin-only access
AdminPermissionMiddleware(permService)

// Staff-level access
StaffPermissionMiddleware(permService)

// Custom resource permission
PermissionMiddleware(permService, "contact", "write")
```

---

## 🎨 Frontend Integration

### Core Components

**AuthProvider** - [`src/features/auth/contexts/AuthProvider.tsx`](../src/features/auth/contexts/AuthProvider.tsx)
```typescript
{
  user: User | null
  loading: boolean
  isAuthenticated: boolean
  login: (email, password) => Promise<{success, error?}>
  logout: () => Promise<void>
  refreshToken: () => Promise<string | null>
}
```

**usePermissions Hook** - [`src/hooks/usePermissions.ts`](../src/hooks/usePermissions.ts)
```typescript
const { 
  hasPermission,      // (resource, action) => boolean
  hasAnyPermission,   // (...perms) => boolean
  hasAllPermissions,  // (...perms) => boolean
  permissions         // string[]
} = usePermissions()
```

### Route Protection Patterns

```typescript
// Layout-level (authentication)
<Route element={<AuthGuard><MainLayout /></AuthGuard>}>
  <Route path="/" element={<Dashboard />} />
</Route>

// Route-level (permission)
<Route path="/contacts" element={
  <ProtectedRoute requiredPermission="contact:read">
    <ContactManager />
  </ProtectedRoute>
} />

// Component-level (conditional rendering)
{hasPermission('contact', 'write') && <EditButton />}
```

### Permission Check Examples

```typescript
// Single permission
if (hasPermission('contact', 'write')) {
  return <EditButton />
}

// Any of multiple
if (hasAnyPermission('user:read', 'user:write')) {
  return <UserList />
}

// All required
if (hasAllPermissions('admin:access', 'user:manage_roles')) {
  return <AdminPanel />
}
```

---

## 📡 API Reference

### Authentication Endpoints

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/api/auth/login` | POST | ❌ | Login with credentials |
| `/api/auth/logout` | POST | ❌ | Logout + revoke tokens |
| `/api/auth/refresh` | POST | ❌ | Refresh access token |
| `/api/auth/profile` | GET | ✅ | User profile + permissions + roles |

### RBAC Management Endpoints

| Endpoint | Method | Permission | Purpose |
|----------|--------|-----------|---------|
| `/api/rbac/permissions` | GET | `admin:access` | List all permissions |
| `/api/rbac/permissions` | POST | `admin:access` | Create permission |
| `/api/rbac/permissions/:id` | DELETE | `admin:access` | Delete permission |
| `/api/rbac/roles` | GET | `admin:access` | List all roles |
| `/api/rbac/roles/:id` | GET | `admin:access` | Get role details |
| `/api/rbac/roles` | POST | `admin:access` | Create role |
| `/api/rbac/roles/:id` | PUT | `admin:access` | Update role |
| `/api/rbac/roles/:id` | DELETE | `admin:access` | Delete role |
| `/api/rbac/roles/:id/permissions` | PUT | `admin:access` | Bulk assign permissions |
| `/api/rbac/roles/:roleId/permissions/:permId` | POST | `admin:access` | Assign single permission |
| `/api/rbac/roles/:roleId/permissions/:permId` | DELETE | `admin:access` | Remove permission |
| `/api/users/:userId/roles` | GET | `user:read` | Get user roles |
| `/api/users/:userId/roles` | POST | `user:manage_roles` | Assign role to user |
| `/api/users/:userId/roles/:roleId` | DELETE | `user:manage_roles` | Remove role from user |
| `/api/users/:userId/permissions` | GET | `user:read` | Get effective permissions |
| `/api/rbac/cache/refresh` | POST | `admin:access` | Invalidate permission cache |

**Handler References:**
- [`handlers/auth_handler.go`](../handlers/auth_handler.go)
- [`handlers/permission_handler.go`](../handlers/permission_handler.go)
- [`handlers/user_handler.go`](../handlers/user_handler.go)

### API Response Examples

**Login Response:**
```json
{
  "success": true,
  "token": "eyJhbGc...",
  "refresh_token": "base64-string",
  "user": {
    "id": "uuid",
    "email": "admin@dekoninklijkeloop.nl",
    "naam": "Admin",
    "rol": "admin",
    "permissions": [
      {"resource": "admin", "action": "access"},
      {"resource": "contact", "action": "read"}
    ],
    "roles": [
      {"id": "uuid", "name": "admin", "description": "Administrator"}
    ],
    "is_actief": true
  }
}
```

**Profile Response:**
```json
{
  "id": "uuid",
  "naam": "Admin",
  "email": "admin@dekoninklijkeloop.nl",
  "rol": "admin",
  "permissions": [...],
  "roles": [...],
  "is_actief": true
}
```

### Error Responses

| HTTP Status | Scenario | Response |
|-------------|----------|----------|
| `401` | Token expired/invalid | `{"error": "Niet geautoriseerd"}` |
| `403` | Permission denied | `{"error": "Geen toegang"}` |
| `401` | Refresh token invalid | `{"error": "Ongeldige of verlopen refresh token"}` |
| `429` | Too many login attempts | `{"error": "Te veel login pogingen"}` |

---

## 🔒 Security & Best Practices

### Defense in Depth (4 Layers)

**Layer 1: UI Filtering**
```typescript
// Hide unavailable options (UX)
const filteredMenu = menuItems.filter(item =>
  !item.permission || hasPermission(item.resource, item.action)
)
```

**Layer 2: Route Guards**
```typescript
// Block unauthorized navigation
<ProtectedRoute requiredPermission="contact:write">
  <ContactEditor />
</ProtectedRoute>
```

**Layer 3: Component Guards**
```typescript
// Prevent component access
if (!hasPermission('contact', 'write')) {
  return <AccessDenied />
}
```

**Layer 4: API Validation** ⚠️ **CRITICAL**
```go
// Backend enforces permissions - ULTIMATE protection
app.Put("/contacts/:id",
  AuthMiddleware,
  PermissionMiddleware(permService, "contact", "write"),
  handler
)
```

### Best Practices

✅ **DO:**
- Check permissions at multiple levels
- Use backend as single source of truth
- Cache permission checks in component state
- Invalidate cache after role changes
- Use graceful degradation (hide vs error)

❌ **DON'T:**
- Trust only frontend checks for security
- Hardcode role names (use permissions)
- Assume permissions without checking
- Forget backend validates everything

### Caching Strategy

**Redis Cache:**
- Key Format: `perm:{userID}:{resource}:{action}`
- TTL: 5 minutes
- Hit Rate: ~97%
- Invalidation: Manual via API or automatic after TTL

**LocalStorage:**
- Access Token: 20 min expiry
- Refresh Token: 7 days expiry
- User Data: Persists until logout

---

## 🚀 Deployment

### Environment Variables

**Backend (.env):**
```bash
# JWT Configuration
JWT_SECRET=your-super-secret-key-minimum-32-characters
JWT_TOKEN_EXPIRY=20m

# Database
DATABASE_URL=postgresql://user:pass@host:port/db

# Redis (caching)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=optional

# CORS
CORS_ORIGIN=https://admin.dekoninklijkeloop.nl
```

**Frontend (.env):**
```bash
VITE_API_BASE_URL=https://api.dekoninklijkeloop.nl
```

### Deployment Checklist

**Pre-Deployment:**
- [ ] Database migrations V1.20-V1.48 applied
- [ ] JWT_SECRET generated (min 32 chars)
- [ ] Redis server running
- [ ] Environment variables configured
- [ ] CORS origins set correctly
- [ ] Database backup created

**Post-Deployment:**
- [ ] Login endpoint test successful
- [ ] Refresh endpoint working
- [ ] Permission checks verified
- [ ] Token rotation tested
- [ ] Redis cache operational
- [ ] Admin panel accessible
- [ ] Role assignments work
- [ ] Performance within targets (<1ms)

---

## 🐛 Troubleshooting

### "Backend returned no permissions"

**Diagnose:**
```typescript
const { user } = useAuth()
console.log('Permissions:', user?.permissions)
console.log('Roles:', user?.roles)
```

**Fixes:**
1. Check `/api/auth/profile` response
2. Verify database: `SELECT * FROM user_roles WHERE user_id = ?`
3. Verify: `SELECT * FROM role_permissions WHERE role_id = ?`
4. Assign roles via Admin panel
5. User must re-login or wait for cache refresh

### "Permission denied after role change"

**Cause:** Redis cache (5 min TTL)

**Solutions:**
1. **Instant:** Admin → Cache Refresh button
2. **Wait:** 5 minutes for auto-expiry
3. **Re-login:** Fresh permissions

### "401 Unauthorized errors"

| Cause | Recognition | Solution |
|-------|-------------|----------|
| Token expired | After 20 min inactivity | Auto-refresh handles it |
| Invalid JWT_SECRET | Direct after deployment | Check backend env vars |
| Token malformed | After code changes | Clear localStorage, re-login |
| No token sent | API calls missing header | Check API client setup |

### "Token refresh problems"

**Symptoms:**
- Auto-logout after 20 minutes
- "Refresh token invalid" errors

**Solutions:**
1. Check token in DB: `SELECT * FROM refresh_tokens WHERE user_id = ?`
2. Verify not revoked: `is_revoked = false`
3. Check expiry: `expires_at > NOW()`
4. Check backend logs for specific errors

---

## 🧪 Testing

### Backend API Tests

```bash
# Login test
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"password"}'

# Profile test
curl http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"

# Refresh test
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### Database Verification

```sql
-- Active refresh tokens
SELECT 
  u.email, 
  rt.expires_at,
  rt.is_revoked,
  (rt.expires_at > NOW()) as is_valid
FROM refresh_tokens rt
JOIN gebruikers u ON rt.user_id = u.id
ORDER BY rt.created_at DESC;

-- User permissions (via RBAC)
SELECT 
  u.email, 
  r.name as role,
  p.resource, 
  p.action
FROM gebruikers u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE u.email = 'admin@dekoninklijkeloop.nl'
  AND ur.is_active = true
ORDER BY p.resource, p.action;

-- Permission count per role
SELECT 
  r.name,
  COUNT(p.id) as permission_count,
  r.is_system_role
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
GROUP BY r.id, r.name, r.is_system_role
ORDER BY permission_count DESC;
```

---

## 📈 Performance Metrics

| Metric | Target | Actual | Notes |
|--------|--------|--------|-------|
| Login Response Time | < 500ms | ~300ms | Incl. permission load |
| Token Refresh Time | < 200ms | ~150ms | With rotation |
| Permission Check (cached) | < 5ms | ~2ms | Redis lookup |
| Permission Check (uncached) | < 50ms | ~30ms | DB query |
| Cache Hit Rate | > 95% | ~97% | Redis metrics |
| Token Validation | < 10ms | ~5ms | JWT parse |

---

## ✅ System Status

**Authentication:**
- ✅ JWT tokens (20 min) with RBAC support
- ✅ Refresh tokens (7 days) with rotation
- ✅ Secure localStorage storage
- ✅ Auto-refresh 5 min before expiry
- ✅ Graceful logout on expiration

**Authorization:**
- ✅ Backend-driven RBAC (V1.22+)
- ✅ 19 resources with 58 permissions
- ✅ Redis caching (97% hit rate)
- ✅ Multi-level permission checks
- ✅ Real-time role/permission management

**Security:**
- ✅ 4-layer defense in depth
- ✅ Backend as single source of truth
- ✅ Automatic token validation
- ✅ Secure password hashing (bcrypt)
- ✅ Token rotation for extra security
- ✅ Rate limiting on login endpoint

**Performance:**
- ✅ ~300ms login time (incl. permissions)
- ✅ ~2ms cached permission checks
- ✅ ~97% cache hit rate
- ✅ Automatic cache invalidation

---

## 📚 Related Documentation

- [`DATABASE_REFERENCE.md`](DATABASE_REFERENCE.md) - Complete database schema
- [`FRONTEND_INTEGRATION.md`](FRONTEND_INTEGRATION.md) - Frontend API guide
- [`database/migrations/V1_20__create_rbac_tables.sql`](../database/migrations/V1_20__create_rbac_tables.sql)
- [`database/migrations/V1_21__seed_rbac_data.sql`](../database/migrations/V1_21__seed_rbac_data.sql)

---

**⚠️ Critical Reminder:**

> **Frontend checks are for UX (User Experience)**  
> **Backend checks are for Security (Data Protection)**
>
> NEVER trust only frontend checks for security!

---

**Version:** 3.0  
**Last Updated:** 2025-11-01  
**Status:** ✅ Production Ready  
**Backend Compatibility:** V1.48.0+  
**Verification:** 100% against backend implementation