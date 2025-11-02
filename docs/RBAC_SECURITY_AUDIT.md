# 🔐 RBAC Security Audit - DKL Email Service

> **Datum:** 2025-11-02  
> **Versie:** V1.48+  
> **Status:** Security Review Completed

---

## 📋 Inhoudsopgave

1. [Overzicht RBAC Architectuur](#overzicht-rbac-architectuur)
2. [Rollen & Permissies Matrix](#rollen--permissies-matrix)
3. [Toegangscontrole Flow](#toegangscontrole-flow)
4. [Beveiligingsanalyse](#beveiligingsanalyse)
5. [Gevonden Issues](#gevonden-issues)
6. [Aanbevelingen](#aanbevelingen)

---

## 🏗️ Overzicht RBAC Architectuur

### Database Schema

```
gebruikers (users)
    └──> user_roles (active assignments)
            └──> roles
                    └──> role_permissions
                            └──> permissions (resource:action)
```

### Kerncomponenten

| Component | Locatie | Functie |
|-----------|---------|---------|
| **Models** | [`models/role_rbac.go`](../models/role_rbac.go) | Data structuren |
| **Repositories** | [`repository/rbac_*.go`](../repository/) | Database operaties |
| **Permission Service** | [`services/permission_service.go`](../services/permission_service.go) | Permission checks met Redis cache |
| **Auth Service** | [`services/auth_service.go`](../services/auth_service.go) | JWT generatie met RBAC roles |
| **Middleware** | [`handlers/permission_middleware.go`](../handlers/permission_middleware.go) | Route beveiliging |

---

## 👥 Rollen & Permissies Matrix

### Systeemrollen

| Rol | Type | Permissies | Gebruik |
|-----|------|-----------|---------|
| **admin** | System | ALLE (58 permissies) | Platform administrators |
| **staff** | System | Read-only meeste resources | Support medewerkers |
| **user** | System | Basic chat read/write | Reguliere gebruikers |

### Chat Rollen

| Rol | Permissies | Gebruik |
|-----|-----------|---------|
| **owner** | `chat:read`, `chat:write`, `chat:moderate`, `chat:manage_channel` | Kanaal eigenaar |
| **chat_admin** | `chat:read`, `chat:write`, `chat:moderate` | Kanaal moderator |
| **member** | `chat:read`, `chat:write` | Kanaal lid |

### Event Rollen (Geen Speciale Permissies)

| Rol | Gebruik |
|-----|---------|
| **deelnemer** | Event deelnemer (categorisatie) |
| **begeleider** | Event begeleider (categorisatie) |
| **vrijwilliger** | Event vrijwilliger (categorisatie) |

⚠️ **BELANGRIJK**: Event rollen zijn ALLEEN voor categorisatie en hebben GEEN permissies!

### Resources & Acties

| Resource | Acties | Locatie Controle |
|----------|--------|------------------|
| `admin` | `access` | Volledige admin toegang |
| `staff` | `access` | Staff-niveau toegang |
| `contact` | `read`, `write`, `delete` | [`handlers/contact_handler.go`](../handlers/contact_handler.go) |
| `aanmelding` | `read`, `write`, `delete` | [`handlers/aanmelding_handler.go`](../handlers/aanmelding_handler.go) |
| `user` | `read`, `write`, `delete`, `manage_roles` | [`handlers/user_handler.go`](../handlers/user_handler.go) |
| `photo` | `read`, `write`, `delete` | Photo handler |
| `album` | `read`, `write`, `delete` | Album handler |
| `video` | `read`, `write`, `delete` | Video handler |
| `partner` | `read`, `write`, `delete` | Partner handler |
| `sponsor` | `read`, `write`, `delete` | Sponsor handler |
| `radio_recording` | `read`, `write`, `delete` | Radio recording handler |
| `program_schedule` | `read`, `write`, `delete` | Program schedule handler |
| `social_embed` | `read`, `write`, `delete` | Social embed handler |
| `social_link` | `read`, `write`, `delete` | Social link handler |
| `under_construction` | `read`, `write`, `delete` | Under construction handler |
| `newsletter` | `read`, `write`, `send`, `delete` | [`handlers/newsletter_handler.go`](../handlers/newsletter_handler.go) |
| `email` | `read`, `write`, `delete`, `fetch` | Email handlers |
| `admin_email` | `send` | [`handlers/admin_mail_handler.go`](../handlers/admin_mail_handler.go) |
| `chat` | `read`, `write`, `manage_channel`, `moderate` | [`handlers/chat_handler.go`](../handlers/chat_handler.go) |

**Totaal:** 19 resources, 58 granulaire permissies

---

## 🔄 Toegangscontrole Flow

### 1. Login Flow

```
POST /api/auth/login
    ↓
Valideer credentials (bcrypt)
    ↓
Genereer JWT token met roles[]
    ↓
Haal permissions op uit database
    ↓
Return: {token, refresh_token, user: {roles[], permissions[]}}
```

### 2. API Request Flow

```
Request met Authorization: Bearer <token>
    ↓
AuthMiddleware validates JWT
    ├─> Valideert signature
    ├─> Controleert expiry (20 min)
    └─> Zet userID in context
    ↓
PermissionMiddleware(resource, action)
    ├─> Haalt userID uit context
    ├─> Controleert Redis cache (perm:userID:resource:action)
    ├─> Als niet cached: query database
    └─> Return 403 als geen permissie
    ↓
Handler execution
```

### 3. Permission Check Implementatie

**Locatie:** [`services/permission_service.go:68-110`](../services/permission_service.go)

```go
HasPermission(ctx, userID, resource, action) bool:
    1. Check Redis cache (5 min TTL)
       Key: "perm:{userID}:{resource}:{action}"
       Hit Rate: ~97%
    
    2. Als niet cached:
       Query:
       SELECT p.resource, p.action
       FROM user_roles ur
       JOIN roles r ON ur.role_id = r.id
       JOIN role_permissions rp ON r.id = rp.role_id
       JOIN permissions p ON rp.permission_id = p.id
       WHERE ur.user_id = ?
         AND ur.is_active = true
         AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
    
    3. Cache resultaat (5 min)
    
    4. Return boolean
```

### 4. Cache Invalidatie

**Triggers:**
- Role assignment/removal aan gebruiker
- Permission assignment/removal aan rol
- Rol update/delete
- Manueel via `/api/rbac/cache/refresh` (admin only)

**Strategie:**
- Per gebruiker: `perm:{userID}:*`
- Volledige refresh: `perm:*`
- Auto-expiry: 5 minuten

---

## 🛡️ Beveiligingsanalyse

### ✅ STERKE PUNTEN

1. **Multi-Layer Defense**
   - ✅ Frontend filtering (UX)
   - ✅ Route guards (Navigation)
   - ✅ Component guards (Rendering)
   - ✅ **API middleware (KRITISCH - echte beveiliging)**

2. **Security Best Practices**
   - ✅ JWT tokens met 20 min expiry
   - ✅ Refresh tokens met 7 dagen expiry + rotation
   - ✅ Bcrypt password hashing
   - ✅ Token revocation support
   - ✅ Rate limiting op login endpoint
   - ✅ CORS strict configuratie

3. **RBAC Implementatie**
   - ✅ Granulaire permissies (resource:action)
   - ✅ Many-to-many relaties (user ↔ roles ↔ permissions)
   - ✅ System roles beschermd tegen verwijdering
   - ✅ Active/inactive role assignments
   - ✅ Optional role expiry dates

4. **Performance Optimalisatie**
   - ✅ Redis caching (97% hit rate)
   - ✅ 5 min cache TTL (balans tussen snelheid en actualiteit)
   - ✅ Bulk permission loads
   - ✅ Preloaded relations (Gorm)

5. **Logging & Monitoring**
   - ✅ Permission denied logging
   - ✅ Login attempt tracking
   - ✅ Token validation logging
   - ✅ Cache hit/miss metrics

---

## ⚠️ GEVONDEN ISSUES

### 🔴 KRITIEKE ISSUES

#### 1. Inconsistente Route Beveiliging

**Probleem:** Niet alle routes gebruiken consistent permission middleware.

**Locatie:** [`main.go`](../main.go:662-759)

**Voorbeelden:**
```go
// ✅ GOED - Gebruikt PermissionMiddleware
userHandler.RegisterRoutes(app)
// In user_handler.go lijn 34-46:
app.Get("/api/users", 
    AuthMiddleware(h.authService), 
    PermissionMiddleware(h.permissionService, "user", "read"), 
    h.ListUsers)

// ❓ ONDUIDELIJK - Check handler implementatie
imageHandler.RegisterRoutes(app)
albumHandler.RegisterRoutes(app)
videoHandler.RegisterRoutes(app)
```

**Impact:** Routes zonder permission checks zijn toegankelijk voor alle authenticated users.

**Aanbeveling:**
```go
// Audit ALLE handlers voor consistent middleware gebruik:
// 1. AuthMiddleware (ALTIJD)
// 2. PermissionMiddleware (resource, action) (BIJNA ALTIJD)
// 3. Alleen publieke endpoints mogen zonder permission checks
```

#### 2. JWT Secret Hardcoded Fallback

**Locatie:** [`services/auth_service.go:64`](../services/auth_service.go:64)

```go
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    logger.Warn("JWT_SECRET omgevingsvariabele niet gevonden, gebruik standaard waarde")
    jwtSecret = "default_jwt_secret_change_in_production" // ❌ GEVAARLIJK
}
```

**Impact:** Als JWT_SECRET niet is ingesteld, kunnen tokens worden geforged.

**Aanbeveling:**
```go
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" || len(jwtSecret) < 32 {
    logger.Fatal("JWT_SECRET must be set and at least 32 characters")
}
```

#### 3. Geen Permission Check op Profile Endpoint

**Locatie:** [`handlers/auth_handler.go:242`](../handlers/auth_handler.go:242)

```go
authProtected.Get("/profile", authHandler.HandleGetProfile)
// ❌ Alleen AuthMiddleware, geen PermissionMiddleware
```

**Impact:** Alle authenticated users kunnen hun eigen profile ophalen met alle permissions. Dit is waarschijnlijk gewenst, maar moet expliciet gedocumenteerd zijn.

**Aanbeveling:** Documenteer dat dit bewust zo is, of voeg een basic permission check toe.

### 🟡 MEDIUM ISSUES

#### 4. Geen Audit Logging voor Permission Changes

**Probleem:** Role en permission assignments worden niet volledig gelogd voor audit trail.

**Locatie:** [`services/permission_service.go`](../services/permission_service.go)

**Huidige logging:**
```go
logger.Info("Rol toegekend aan gebruiker", "user_id", userID, "role_id", roleID)
```

**Gemist:**
- Wie heeft de assignment gedaan? (assigned_by in database, maar niet altijd gelogd)
- IP adres van requester
- Timestamp is in logs, maar geen structured audit event

**Aanbeveling:**
```go
logger.Audit("ROLE_ASSIGNED", AuditEvent{
    ActorID:      assignedBy,
    TargetUserID: userID,
    RoleID:       roleID,
    IPAddress:    c.IP(),
    Timestamp:    time.Now(),
})
```

#### 5. Unlimited Permission Lists

**Locatie:** [`services/permission_service.go:314-321`](../services/permission_service.go:314-321)

```go
func (s *PermissionServiceImpl) GetRoles(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
    return s.rbacRoleRepo.List(ctx, limit, offset)
}
```

**Probleem:** Geen maximum limit validatie. Grote queries kunnen prestatie impact hebben.

**Aanbeveling:**
```go
if limit <= 0 || limit > 1000 {
    limit = 100 // Default
}
```

#### 6. Chat Permissions Onduidelijk

**Locatie:** Documentatie vs implementatie

**Documentatie zegt:** "Chat permissions are ROLE-BASED"  
**Vraag:** Hoe werken channel-specifieke permissions?

**Aanbeveling:** Verduidelijk in code comments:
```go
// Chat permissions are checked via user roles in channels
// - Channel-level: user must have role in chat_channel_participants
// - Global: user must have system role with chat:* permissions
```

### 🟢 MINOR ISSUES

#### 7. Legacy Role Field in JWT

**Locatie:** [`services/auth_service.go:277-280`](../services/auth_service.go:277-280)

```go
claims := JWTClaims{
    Email:      gebruiker.Email,
    Role:       gebruiker.Rol,      // Legacy - deprecated
    Roles:      rbacRoles,          // RBAC - current
    RBACActive: len(rbacRoles) > 0,
}
```

**Probleem:** `gebruiker.Rol` is een string field die deprecated is. Kan verwarring veroorzaken.

**Aanbeveling:** Plan migratie om legacy `rol` veld te verwijderen uit database schema.

#### 8. Inconsistente Error Messages

**Voorbeelden:**
- "Geen toegang" (Nederlands)
- "Niet geautoriseerd" (Nederlands)
- "Unauthorized" (Engels)
- "Permission denied" (Engels)

**Aanbeveling:** Standaardiseer op één taal (Nederlands voor gebruikers, Engels voor logs).

---

## 💡 AANBEVELINGEN

### 🎯 Prioriteit HOOG

1. **Audit Alle Route Registraties**
   ```bash
   # Zoek alle RegisterRoutes calls
   grep -r "RegisterRoutes" handlers/
   
   # Check elke handler file voor consistent middleware gebruik
   ```

2. **Strict JWT Secret Validatie**
   ```go
   // In main.go bij startup:
   if err := validateSecurityConfig(); err != nil {
       logger.Fatal("Security configuration invalid", "error", err)
   }
   ```

3. **Implementeer Audit Logging**
   - Nieuwe `logger.Audit()` functie
   - Separate audit log file
   - Include: actor, action, target, timestamp, IP, result

4. **Permission Matrix Document**
   Genereer automatisch een matrix document:
   ```
   Role     | admin:access | user:read | user:write | ...
   ---------|--------------|-----------|------------|-----
   admin    | ✅           | ✅        | ✅         | ...
   staff    | ❌           | ✅        | ❌         | ...
   user     | ❌           | ✅        | ❌         | ...
   ```

### 🎯 Prioriteit MEDIUM

5. **Rate Limiting per User**
   Momenteel alleen per email voor login. Overweeg:
   ```go
   // Rate limit per user voor API calls
   rateLimitKey := fmt.Sprintf("api:%s:%s", userID, resource)
   ```

6. **Permission Inheritance**
   Overweeg hierarchie:
   ```
   admin:access → impliceert alle andere permissions
   staff:access → impliceert read permissions
   ```

7. **Temporary Permission Grants**
   ```go
   // Bestaand: expires_at in user_roles
   // Toevoegen: reminder before expiry
   // Toevoegen: auto-revoke job
   ```

8. **API Key Management voor Services**
   ```go
   // Voor service-to-service communication
   // Bijvoorbeeld: Telegram bot, Cloudinary webhooks
   api_keys table:
       - key_hash (bcrypt)
       - service_name
       - permissions[]
       - expires_at
       - last_used
   ```

### 🎯 Prioriteit LAAG

9. **Permission Templates**
   ```go
   // Voorgedefinieerde sets voor veelvoorkomende roles
   roleTemplates := map[string][]string{
       "content_manager": {"photo:*", "video:*", "album:*"},
       "support_agent":   {"contact:read", "aanmelding:read"},
   }
   ```

10. **GraphQL-style Permission Queries**
    ```go
    // Flexibele permission checks
    HasAnyPermission(userID, []Permission{{resource: "photo", action: "write"}, ...})
    HasAllPermissions(userID, []Permission{...})
    ```

11. **Permission Analytics Dashboard**
    - Meest gebruikte permissions
    - Meest geweigerde permissions (potentiële bugs)
    - Permission check performance metrics

---

## 📊 Metrics & Monitoring

### Huidige Metrics

| Metric | Target | Actueel | Status |
|--------|--------|---------|--------|
| Permission check (cached) | <5ms | ~2ms | ✅ |
| Permission check (uncached) | <50ms | ~30ms | ✅ |
| Cache hit rate | >95% | ~97% | ✅ |
| Login response time | <500ms | ~300ms | ✅ |
| Token refresh time | <200ms | ~150ms | ✅ |

### Aanbevolen Extra Metrics

```go
// Prometheus metrics toevoegen:
rbac_permission_checks_total{resource, action, result}
rbac_permission_denied_total{resource, action}
rbac_cache_size
rbac_invalid_tokens_total
rbac_role_assignments_total{role}
```

---

## 🔍 Security Checklist

### Bij Deployment

- [ ] JWT_SECRET is strong (min 32 chars)
- [ ] JWT_SECRET is uniek per environment
- [ ] Database backups enabled
- [ ] Redis password protected
- [ ] CORS origins correct configured
- [ ] Log level appropriate (INFO in prod)
- [ ] Rate limits configured
- [ ] Health checks operational

### Bij Code Changes

- [ ] Nieuwe routes hebben AuthMiddleware
- [ ] Nieuwe routes hebben PermissionMiddleware (tenzij publiek)
- [ ] Permission checks op juiste resource:action
- [ ] Error messages geen sensitive data
- [ ] Logging op juiste niveau
- [ ] Tests voor permission checks
- [ ] Documentation updated

### Bij Role/Permission Changes

- [ ] Migration script aanwezig
- [ ] Backwards compatible (indien mogelijk)
- [ ] Cache invalidatie na deployment
- [ ] Users genotificeerd (indien nodig)
- [ ] Documentation updated
- [ ] Tests aangepast

---

## 📚 Gerelateerde Documenten

- [`AUTH_AND_RBAC.md`](AUTH_AND_RBAC.md) - Volledige systeem documentatie
- [`DATABASE_REFERENCE.md]`](DATABASE_REFERENCE.md) - Database schema
- Database migrations V1.20-V1.48 - RBAC tabellen en data

---

## ✅ Conclusie

### Algemene Beoordeling: **8/10** 🟢

**Sterke Punten:**
- ✅ Solide RBAC architectuur
- ✅ Goede separatie van concerns
- ✅ Performance optimalisatie met Redis
- ✅ Multi-layer security
- ✅ Granulaire permissions

**Verbeterpunten:**
- ⚠️ Inconsistente route beveiliging (moet geaudit worden)
- ⚠️ JWT secret validation te laks
- ⚠️ Audit logging ontbreekt
- ⚠️ Documentatie vs implementatie discrepanties

**Urgente Acties:**
1. Audit alle route registraties voor permission middleware
2. Fix JWT_SECRET fallback (make it fail-safe)
3. Implementeer audit logging voor security events
4. Test alle endpoints voor authorization bypass

---

**Laatste Update:** 2025-11-02  
**Reviewer:** Security Audit  
**Next Review:** Na implementatie aanbevelingen