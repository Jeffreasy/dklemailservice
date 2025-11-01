# RBAC Implementatie Guide

**Versie**: 1.22.0  
**Datum**: 2025-11-01  
**Status**: Ready for Deployment

---

## 📋 Overzicht van Wijzigingen

### 🎯 Uitgevoerde Fixes

#### 1. ✅ Role Constants Conflict Opgelost
**Bestand**: [`models/role.go`](../models/role.go)

**Probleem**: 
- `RoleAdmin` en `RoleChatAdmin` hadden beide waarde "admin"

**Oplossing**:
```go
const (
    RoleAdmin       Role = "admin"
    RoleChatOwner   Role = "owner"
    RoleChatAdmin   Role = "chat_admin"  // ✅ Changed from "admin"
    RoleChatMember  Role = "member"
    RoleDeelnemer   Role = "deelnemer"   // ✅ Lowercase voor consistency
    RoleBegeleider  Role = "begeleider"  // ✅ Lowercase voor consistency  
    RoleVrijwilliger Role = "vrijwilliger" // ✅ Lowercase voor consistency
)
```

**Impact**: ⚠️ **Breaking change** voor code die `RoleChatAdmin` gebruikt!

#### 2. ✅ Database Migratie Script
**Bestand**: [`database/migrations/V1_22__migrate_legacy_roles_to_rbac.sql`](../database/migrations/V1_22__migrate_legacy_roles_to_rbac.sql)

**Functionaliteit**:
- Migreert alle bestaande `gebruikers.rol` values naar `user_roles` tabel
- Behoudt legacy field voor backward compatibility
- Voegt automatisch 'user' role toe aan gebruikers zonder rol
- Maakt monitoring view: `v_user_role_migration_status`
- Uitgebreide logging en status reports

**Features**:
- ✅ Duplicate prevention
- ✅ Migration status tracking
- ✅ Detailed logging
- ✅ Problem detection

#### 3. ✅ JWT RBAC Support
**Bestand**: [`services/auth_service.go`](../services/auth_service.go)

**Wijzigingen**:

**JWTClaims Extended**:
```go
type JWTClaims struct {
    Email      string   `json:"email"`
    Role       string   `json:"role"`        // Legacy - deprecated
    Roles      []string `json:"roles"`       // ✅ NEW: RBAC roles
    RBACActive bool     `json:"rbac_active"` // ✅ NEW: RBAC indicator
    jwt.RegisteredClaims
}
```

**Token Generation Enhanced**:
```go
func (s *AuthServiceImpl) generateToken(gebruiker *models.Gebruiker) (string, error) {
    rbacRoles := s.getUserRBACRoles(gebruiker.ID)  // ✅ NEW
    
    claims := JWTClaims{
        Email:      gebruiker.Email,
        Role:       gebruiker.Rol,      // Legacy
        Roles:      rbacRoles,          // ✅ RBAC
        RBACActive: len(rbacRoles) > 0, // ✅ Indicator
        // ...
    }
}
```

**New Helper Method**:
```go
func (s *AuthServiceImpl) getUserRBACRoles(userID string) []string {
    // Queries user_roles table via UserRoleRepository
    // Returns array of role names: ["admin", "staff"]
}
```

#### 4. ✅ Service Factory Update
**Bestand**: [`services/factory.go`](../services/factory.go)

**Wijziging**:
```go
// Before:
authService := NewAuthService(repoFactory.Gebruiker, repoFactory.RefreshToken)

// After:
authService := NewAuthServiceWithRBAC(
    repoFactory.Gebruiker, 
    repoFactory.RefreshToken,
    repoFactory.UserRole,  // ✅ RBAC support
)
```

---

## 🚀 Deployment Stappen

### Stap 1: Pre-Deployment Checklist

```bash
# 1. Backup database
pg_dump -h localhost -U postgres dklemailservice > backup_pre_v1.22.sql

# 2. Verifieer huidige migratie status
psql -h localhost -U postgres -d dklemailservice -c "SELECT * FROM migraties ORDER BY toegepast DESC LIMIT 5;"

# 3. Tel huidige gebruikers
psql -h localhost -U postgres -d dklemailservice -c "SELECT COUNT(*) as total_users FROM gebruikers;"
```

### Stap 2: Deploy Code Changes

```bash
# 1. Pull latest code (of deploy via git)
git pull origin main

# 2. Build nieuwe versie
go build -o dklemailservice

# 3. Stop huidige service
# (Afhankelijk van deployment methode)
```

### Stap 3: Run Database Migration

**Optie A - Via Application Startup** (Aanbevolen)
```bash
# Migratie gebeurt automatisch bij opstarten
./dklemailservice
```

**Optie B - Handmatig via PostgreSQL**
```bash
psql -h localhost -U postgres -d dklemailservice -f database/migrations/V1_22__migrate_legacy_roles_to_rbac.sql
```

### Stap 4: Verificatie

```sql
-- 1. Check migratie status
SELECT * FROM v_user_role_migration_status 
ORDER BY migration_status, email;

-- 2. Count migration results
SELECT 
    migration_status,
    COUNT(*) as count
FROM v_user_role_migration_status
GROUP BY migration_status;

-- Expected output:
-- | migration_status | count |
-- |-----------------|-------|
-- | MIGRATED        | XX    |
-- | NO LEGACY       | XX    |

-- 3. Check for problems
SELECT * FROM v_user_role_migration_status
WHERE migration_status IN ('MISMATCH', 'MISSING RBAC')
LIMIT 10;

-- 4. Verify een specifieke gebruiker
SELECT 
    g.email,
    g.rol as legacy_role,
    STRING_AGG(r.name, ', ') as rbac_roles
FROM gebruikers g
LEFT JOIN user_roles ur ON g.id = ur.user_id AND ur.is_active = true
LEFT JOIN roles r ON ur.role_id = r.id
WHERE g.email = 'admin@dekoninklijkeloop.nl'
GROUP BY g.id, g.email, g.rol;
```

### Stap 5: Test JWT Tokens

```bash
# 1. Login en krijg nieuw token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"admin"}'

# 2. Decode JWT token (gebruik jwt.io of een JWT decoder)
# Verifieer dat token bevat:
# - "role": "admin" (legacy)
# - "roles": ["admin"] (RBAC)
# - "rbac_active": true
```

### Stap 6: Test Permission Checks

```bash
# Test admin endpoint met nieuwe token
curl -X GET http://localhost:8080/api/rbac/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Verwachte response: 200 OK met lijst van roles
```

### Stap 7: Monitor Logs

```bash
# Check voor errors
tail -f logs/application.log | grep -i "error\|warn\|rbac"

# Check voor permission checks
tail -f logs/application.log | grep "Permission"
```

---

## 🔧 Post-Deployment Stappen

### 1. Verify RBAC Functionality

```sql
-- Test permission check voor admin user
SELECT 
    u.email,
    r.name as role_name,
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

-- Verwacht: Alle permissions omdat admin alle permissions heeft
```

### 2. Clean Up Migration View (Optioneel)

Na succesvolle migratie en verificatie:
```sql
-- De view is handig voor monitoring, maar kan verwijderd worden indien gewenst
-- DROP VIEW IF EXISTS v_user_role_migration_status;
```

### 3. Setup Monitoring

```sql
-- Create functie om RBAC status te monitoren
CREATE OR REPLACE FUNCTION check_rbac_health() 
RETURNS TABLE (
    metric TEXT,
    value BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 'total_users'::TEXT, COUNT(*)::BIGINT FROM gebruikers
    UNION ALL
    SELECT 'users_with_rbac_roles'::TEXT, COUNT(DISTINCT user_id)::BIGINT FROM user_roles WHERE is_active = true
    UNION ALL
    SELECT 'total_roles'::TEXT, COUNT(*)::BIGINT FROM roles
    UNION ALL
    SELECT 'system_roles'::TEXT, COUNT(*)::BIGINT FROM roles WHERE is_system_role = true
    UNION ALL
    SELECT 'total_permissions'::TEXT, COUNT(*)::BIGINT FROM permissions
    UNION ALL
    SELECT 'active_role_assignments'::TEXT, COUNT(*)::BIGINT FROM user_roles WHERE is_active = true;
END;
$$ LANGUAGE plpgsql;

-- Run periodically
SELECT * FROM check_rbac_health();
```

---

## 🐛 Troubleshooting

### Problem: Users zonder RBAC roles na migratie

**Diagnose**:
```sql
SELECT * FROM v_user_role_migration_status
WHERE migration_status = 'MISSING RBAC';
```

**Fix**:
```sql
-- Voeg handmatig 'user' role toe
INSERT INTO user_roles (user_id, role_id, is_active)
SELECT 
    g.id,
    r.id,
    true
FROM gebruikers g
CROSS JOIN roles r
WHERE r.name = 'user'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur WHERE ur.user_id = g.id
  );
```

### Problem: JWT tokens zonder RBAC roles

**Diagnose**:
- Login en decode JWT
- Check of `roles` array leeg is terwijl `rbac_active` true is

**Oorzaak**: UserRoleRepository niet correct geinitializeerd

**Fix**:
```go
// In services/factory.go - verify initialization:
authService := NewAuthServiceWithRBAC(
    repoFactory.Gebruiker, 
    repoFactory.RefreshToken,
    repoFactory.UserRole,  // ← Moet niet nil zijn!
)
```

### Problem: Permission denied na migratie

**Diagnose**:
```sql
-- Check of user permissions heeft
SELECT * FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = 'USER_ID_HERE'
  AND ur.is_active = true;

-- Check of role permissions heeft
SELECT * FROM role_permissions rp
JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = 'ROLE_ID_HERE';
```

**Fix**: Assign correct role of permission

### Problem: Cache issues

**Diagnose**:
```bash
# Check Redis
redis-cli KEYS "perm:*"
```

**Fix**:
```bash
# Clear permission cache
redis-cli --scan --pattern "perm:*" | xargs redis-cli DEL

# Of via API (als admin):
curl -X POST http://localhost:8080/api/rbac/cache/refresh \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📊 Backwards Compatibility

### Wat blijft werken?

✅ **Oude JWT tokens**:
- Tokens met alleen `role` field werken nog
- Nieuwe tokens hebben beide `role` EN `roles`

✅ **Legacy `gebruikers.rol` field**:
- Blijft bestaan in database
- Wordt nog steeds gebruikt voor JWT `role` claim
- Wordt NIET gebruikt voor permission checks

✅ **Bestaande API endpoints**:
- Alle bestaande endpoints blijven werken
- Permission checks gebruiken automatisch RBAC

### Wat verandert?

⚠️ **Permission checks**:
- Gebruiken nu `user_roles` tabel
- Legacy `gebruikers.rol` wordt genegeerd voor permissions

⚠️ **Nieuwe JWT tokens**:
- Bevatten extra `roles` array
- Bevatten `rbac_active` boolean

⚠️ **Role constants**:
- `RoleChatAdmin` = "chat_admin" (was "admin")
- Event roles zijn lowercase (was Capitalized)

---

## 🧪 Testing Checklist

### Pre-Deployment Tests

```bash
# 1. Run Go tests
go test ./... -v

# 2. Specifiek test auth service
go test ./tests -v -run TestAuth

# 3. Test migratie in test database
psql -h localhost -U postgres -d dklemailservice_test \
  -f database/migrations/V1_22__migrate_legacy_roles_to_rbac.sql
```

### Post-Deployment Tests

#### Test 1: Login met Admin
```bash
# Request
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "admin@dekoninklijkeloop.nl",
  "wachtwoord": "admin"
}

# Expected Response: 200 OK
{
  "access_token": "eyJ...",
  "refresh_token": "..."
}

# Decode JWT en verify:
# - "role": "admin"
# - "roles": ["admin"]
# - "rbac_active": true
```

#### Test 2: Permission Check
```bash
# Request (met JWT van Test 1)
GET http://localhost:8080/api/rbac/roles
Authorization: Bearer eyJ...

# Expected Response: 200 OK
[
  {
    "id": "...",
    "name": "admin",
    "description": "...",
    "permissions": [...]
  },
  ...
]
```

#### Test 3: Staff User (Read-Only)
```bash
# 1. Maak staff user aan (als admin)
POST http://localhost:8080/api/users
Authorization: Bearer ADMIN_TOKEN
Content-Type: application/json

{
  "naam": "Staff User",
  "email": "staff@test.nl",
  "wachtwoord": "test123",
  "rol": "staff"
}

# 2. Login als staff
POST http://localhost:8080/api/auth/login
{
  "email": "staff@test.nl",
  "wachtwoord": "test123"
}

# 3. Test read permission (should work)
GET http://localhost:8080/api/contacts
Authorization: Bearer STAFF_TOKEN

# Expected: 200 OK

# 4. Test write permission (should fail)
POST http://localhost:8080/api/rbac/roles
Authorization: Bearer STAFF_TOKEN
Content-Type: application/json

{
  "name": "test_role",
  "description": "Test"
}

# Expected: 403 Forbidden
{
  "error": "Geen toegang"
}
```

#### Test 4: Migration Status View
```sql
-- Query migration status
SELECT 
    migration_status,
    COUNT(*) as count
FROM v_user_role_migration_status
GROUP BY migration_status;

-- Expected:
-- | migration_status | count |
-- |-----------------|-------|
-- | MIGRATED        | X     |
-- | NO LEGACY       | Y     |
-- (geen MISMATCH of MISSING RBAC)
```

#### Test 5: Permission Caching
```bash
# 1. Eerste request (cache miss - slower)
time curl -X GET http://localhost:8080/api/rbac/roles \
  -H "Authorization: Bearer TOKEN"

# 2. Tweede request (cache hit - faster)
time curl -X GET http://localhost:8080/api/rbac/roles \
  -H "Authorization: Bearer TOKEN"

# Expected: Tweede request significant sneller (1-2ms vs 10-50ms)
```

---

## 📝 Rollback Plan

Als er problemen ontstaan:

### Optie 1: Quick Rollback (Code Only)

```bash
# 1. Revert naar vorige versie
git revert HEAD
git push

# 2. Rebuild en deploy
go build -o dklemailservice
# Deploy nieuwe build

# 3. Database blijft onveranderd (user_roles blijven bestaan)
# Legacy systeem blijft werken
```

### Optie 2: Full Rollback (Code + Database)

```bash
# 1. Stop service
systemctl stop dklemailservice  # of equivalent

# 2. Restore database backup
psql -h localhost -U postgres -d dklemailservice < backup_pre_v1.22.sql

# 3. Revert code
git revert HEAD
git push

# 4. Rebuild en start
go build -o dklemailservice
systemctl start dklemailservice
```

### Optie 3: Remove Migration Only

```sql
-- Verwijder user_role assignments gemaakt door migratie
DELETE FROM user_roles
WHERE assigned_at >= (SELECT toegepast FROM migraties WHERE versie = '1.22.0');

-- Verwijder migratie record
DELETE FROM migraties WHERE versie = '1.22.0';

-- Drop migration view
DROP VIEW IF EXISTS v_user_role_migration_status;
```

---

## 🔍 Monitoring Queries

### Daily Health Checks

```sql
-- 1. RBAC System Health
SELECT * FROM check_rbac_health();

-- 2. Users zonder active roles
SELECT 
    g.email, 
    g.naam,
    g.rol as legacy_role
FROM gebruikers g
WHERE NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id 
    AND ur.is_active = true
)
AND g.is_actief = true
ORDER BY g.email;

-- 3. Expired roles die nog active zijn
SELECT 
    u.email,
    r.name as role_name,
    ur.expires_at,
    ur.is_active
FROM user_roles ur
JOIN gebruikers u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE ur.expires_at IS NOT NULL 
  AND ur.expires_at < NOW()
  AND ur.is_active = true;

-- 4. Permission distribution
SELECT 
    p.resource,
    p.action,
    COUNT(DISTINCT ur.user_id) as user_count
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.is_active = true
GROUP BY p.resource, p.action
ORDER BY user_count DESC, p.resource, p.action;
```

### Redis Cache Metrics

```bash
# Check cache hit rate (via Redis)
redis-cli INFO stats | grep keyspace

# Count permission cache keys
redis-cli KEYS "perm:*" | wc -l

# Sample cache key
redis-cli GET "perm:USER_ID:contact:read"
```

---

## 🎓 Frontend Integration Guide

### Using RBAC in Frontend

#### TypeScript Types

```typescript
// JWT Token Type
interface JWTClaims {
  email: string;
  role: string;        // Legacy - deprecated
  roles: string[];     // RBAC roles
  rbac_active: boolean; // RBAC indicator
  exp: number;
  iat: number;
  sub: string;         // User ID
}

// Permission Check
interface Permission {
  resource: string;
  action: string;
}

// User Permission dari API
interface UserPermission {
  user_id: string;
  email: string;
  role_name: string;
  resource: string;
  action: string;
  permission_assigned_at: string;
  role_assigned_at: string;
}
```

#### React Hook Voorbeeld

```typescript
// usePermission.ts
import { useMemo } from 'react';
import { useAuth } from './useAuth';

export function usePermission(resource: string, action: string): boolean {
  const { user, permissions } = useAuth();
  
  return useMemo(() => {
    if (!user || !permissions) return false;
    
    return permissions.some(p => 
      p.resource === resource && p.action === action
    );
  }, [user, permissions, resource, action]);
}

// Usage in component
function DeleteButton() {
  const canDelete = usePermission('user', 'delete');
  
  if (!canDelete) return null;
  
  return <button onClick={handleDelete}>Delete</button>;
}
```

#### Fetch Permissions

```typescript
// API call om user permissions op te halen
async function fetchUserPermissions(userId: string): Promise<UserPermission[]> {
  const response = await fetch(`/api/users/${userId}/permissions`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`
    }
  });
  
  if (!response.ok) throw new Error('Failed to fetch permissions');
  
  return response.json();
}
```

#### Permission-Based Routing

```typescript
// ProtectedRoute.tsx
import { Navigate } from 'react-router-dom';
import { usePermission } from './usePermission';

interface ProtectedRouteProps {
  children: React.ReactNode;
  resource: string;
  action: string;
}

export function ProtectedRoute({ children, resource, action }: ProtectedRouteProps) {
  const hasPermission = usePermission(resource, action);
  
  if (!hasPermission) {
    return <Navigate to="/unauthorized" replace />;
  }
  
  return <>{children}</>;
}

// Usage in router
<Route path="/admin/users" element={
  <ProtectedRoute resource="user" action="read">
    <UserManagementPage />
  </ProtectedRoute>
} />
```

---

## 📚 API Reference voor Frontend

### GET /api/users/:id/permissions

Haalt alle permissions op voor een gebruiker.

**Response**:
```json
[
  {
    "user_id": "uuid",
    "email": "admin@example.nl",
    "role_name": "admin",
    "resource": "contact",
    "action": "read",
    "permission_assigned_at": "2025-11-01T10:00:00Z",
    "role_assigned_at": "2025-11-01T09:00:00Z"
  },
  ...
]
```

### GET /api/users/:id/roles

Haalt alle roles op voor een gebruiker.

**Response**:
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "role_id": "uuid",
    "assigned_at": "2025-11-01T09:00:00Z",
    "is_active": true,
    "role": {
      "id": "uuid",
      "name": "admin",
      "description": "...",
      "is_system_role": true
    }
  }
]
```

### GET /api/rbac/roles

Haalt alle roles met hun permissions op (Admin only).

**Response**:
```json
[
  {
    "id": "uuid",
    "name": "admin",
    "description": "Volledige beheerder",
    "is_system_role": true,
    "permissions": [
      {
        "id": "uuid",
        "resource": "contact",
        "action": "read",
        "description": "Contactformulieren bekijken"
      },
      ...
    ]
  }
]
```

### GET /api/rbac/permissions?group_by_resource=true

Haalt alle permissions gegroepeerd op (Admin only).

**Response**:
```json
{
  "groups": [
    {
      "resource": "contact",
      "permissions": [
        {
          "id": "uuid",
          "resource": "contact",
          "action": "read",
          "description": "..."
        },
        ...
      ],
      "count": 3
    },
    ...
  ],
  "total": 65
}
```

---

## 🔐 Security Considerations

### 1. Token Security

- ✅ JWT tokens zijn short-lived (20 min default)
- ✅ Refresh tokens hebben langere TTL (7 dagen)
- ✅ Refresh token rotation bij elke refresh
- ✅ Revocation support via database

### 2. Permission Caching

- ✅ Cache TTL is kort (5 min) voor snelle updates
- ✅ Cache wordt geinvalideerd bij role changes
- ✅ Fallback naar database bij cache failures

### 3. System Protection

- ✅ System roles kunnen niet worden verwijderd
- ✅ System permissions kunnen niet worden verwijderd
- ✅ Audit trail voor alle wijzigingen

### 4. Rate Limiting

Blijft actief voor:
- Login attempts (5 per 5 min)
- Contact forms (5 per hour)
- Registrations (3 per day)

---

## 📖 Best Practices voor Developers

### 1. Nieuwe Permission Toevoegen

```go
// 1. Voeg permission toe via API (als admin)
POST /api/rbac/permissions
{
  "resource": "report",
  "action": "generate",
  "description": "Rapporten genereren"
}

// 2. Assign aan roles
POST /api/rbac/roles/:role_id/permissions/:permission_id

// 3. Use in code
app.Post("/api/reports", 
  AuthMiddleware(authService),
  PermissionMiddleware(permService, "report", "generate"),
  reportHandler.Generate)
```

### 2. Nieuwe Endpoint Beveiligen

```go
// Optie A: Specific permission
app.Get("/api/sensitive-data",
    AuthMiddleware(authService),
    PermissionMiddleware(permService, "data", "read"),
    handler.GetSensitiveData)

// Optie B: Admin only
app.Delete("/api/users/:id",
    AuthMiddleware(authService),
    AdminPermissionMiddleware(permService),
    handler.DeleteUser)

// Optie C: Multiple permissions
app.Post("/api/complex-operation",
    AuthMiddleware(authService),
    ResourcePermissionMiddleware(permService,
        models.PermissionCheck{Resource: "data", Action: "read"},
        models.PermissionCheck{Resource: "data", Action: "write"},
    ),
    handler.ComplexOperation)
```

### 3. Permission Check in Handler

```go
func (h *Handler) SomeAction(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
    
    // Runtime permission check
    if !h.permissionService.HasPermission(c.Context(), userID, "resource", "action") {
        return c.Status(403).JSON(fiber.Map{"error": "Geen toegang"})
    }
    
    // Continue with action
}
```

### 4. Assign Role Programmatically

```go
// In een service of handler
func (s *Service) PromoteToAdmin(ctx context.Context, userID, adminID string) error {
    // Get admin role
    adminRole, err := s.rbacRoleRepo.GetByName(ctx, "admin")
    if err != nil {
        return err
    }
    
    // Assign role
    return s.permissionService.AssignRole(ctx, userID, adminRole.ID, &adminID)
}
```

---

## 🎯 Next Steps (Phase 2)

Na succesvolle deployment van V1.22:

### Short Term (1-2 weken)
1. ✅ Monitor logs voor issues
2. ✅ Verify all users have RBAC roles
3. ✅ Collect frontend integration feedback
4. ⬜ Add permission checks to remaining unprotected endpoints

### Medium Term (1-2 maanden)
5. ⬜ Create admin dashboard voor RBAC management
6. ⬜ Add bulk role assignment feature
7. ⬜ Implement role templates
8. ⬜ Add permission usage analytics

### Long Term (3-6 maanden)
9. ⬜ Deprecate legacy `gebruikers.rol` field
10. ⬜ Remove legacy role constants
11. ⬜ Full RBAC in JWT (remove legacy `role` claim)
12. ⬜ Advanced features (2FA for role elevation, etc.)

---

## 📞 Support

### Common Questions

**Q: Moet ik alle users opnieuw laten inloggen?**  
A: Nee, oude tokens blijven werken. Nieuwe tokens krijgen automatisch RBAC roles.

**Q: Wat als Redis down gaat?**  
A: Permission checks vallen automatisch terug op database. Werkt normaal maar langzamer.

**Q: Hoe voeg ik een nieuwe role toe?**  
A: Via API: `POST /api/rbac/roles` met admin token, of direct in database.

**Q: Kan ik system roles bewerken?**  
A: Nee, system roles zijn read-only voor veiligheid. Maak custom roles voor wijzigingen.

**Q: Hoe weet ik welke permissions een user heeft?**  
A: `GET /api/users/:id/permissions` of query `v_user_role_migration_status` view.

### Debugging

```bash
# Enable debug logging
export LOG_LEVEL=DEBUG

# Check permission service logs
grep "Permission" logs/application.log | tail -n 50

# Check Redis cache
redis-cli --scan --pattern "perm:*" | head -n 10
redis-cli GET "perm:USER_ID:contact:read"

# Check database connections
psql -h localhost -U postgres -d dklemailservice -c "SELECT COUNT(*) FROM user_roles WHERE is_active = true;"
```

---

## ✅ Sign-Off Checklist

Voor je de wijzigingen deployed:

- [ ] Database backup gemaakt
- [ ] Alle tests draaien groen
- [ ] Code review afgerond
- [ ] Migration script getest in test environment
- [ ] Rollback plan gereview
- [ ] Monitoring setup geverifieerd
- [ ] Team geïnformeerd over wijzigingen
- [ ] Documentation bijgewerkt

Na deployment:

- [ ] Verify migration status view
- [ ] Test admin login en permissions
- [ ] Test staff login en permissions
- [ ] Monitor logs voor 30 minuten
- [ ] Verify Redis cache functionaliteit
- [ ] Test frontend integration (indien applicable)

---

**Document Version**: 1.0  
**Last Updated**: 2025-11-01  
**Author**: Kilo Code AI Assistant  
**Review Status**: Ready for Implementation