# RBAC en Database Rollen Logica - Volledige Analyse

**Gegenereerd op**: 2025-11-01  
**Project**: DKL Email Service  
**Versie**: 1.21.0+

---

## 📋 Inhoudsopgave

1. [Executive Summary](#executive-summary)
2. [Database Schema](#database-schema)
3. [Model Structuren](#model-structuren)
4. [Repository Laag](#repository-laag)
5. [Service Laag](#service-laag)
6. [Handler & Middleware Laag](#handler--middleware-laag)
7. [Rollen Overzicht](#rollen-overzicht)
8. [Permissions Overzicht](#permissions-overzicht)
9. [Data Flow](#data-flow)
10. [Cache Strategie](#cache-strategie)
11. [Beveiligings Aspecten](#beveiligings-aspecten)
12. [Migratie Status](#migratie-status)
13. [Legacy vs RBAC](#legacy-vs-rbac)
14. [Best Practices](#best-practices)
15. [Aanbevelingen](#aanbevelingen)

---

## 🎯 Executive Summary

Het DKL Email Service project heeft een **duaal rollen systeem** geïmplementeerd:

### Huidige Situatie
- ✅ **Legacy Role System**: Gebaseerd op string-based roles in `gebruikers.rol` veld
- ✅ **Modern RBAC System**: Volledig Role-Based Access Control met granulaire permissions
- ⚠️ **Transitie Fase**: Beide systemen werken naast elkaar

### Belangrijkste Kenmerken
- **4 Database Tabellen** voor RBAC: `roles`, `permissions`, `role_permissions`, `user_roles`
- **9 System Roles** gedefinieerd (admin, staff, user, chat roles, event roles)
- **65+ Permissions** verspreid over 11 resources
- **Redis Caching** voor performance optimalisatie
- **JWT Authentication** geïntegreerd met permissie checks

---

## 🗄️ Database Schema

### RBAC Tabellen

#### 1. `roles` (Rollen)
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by UUID
);
```

**Kenmerken**:
- System roles kunnen niet worden verwijderd
- Unieke rol namen
- Audit trail via `created_by`

#### 2. `permissions` (Permissies)
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    is_system_permission BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(resource, action)
);
```

**Kenmerken**:
- Resource-action combinatie is uniek
- System permissions zijn beschermd
- Beschrijving voor documentatie

#### 3. `role_permissions` (Rol-Permissie Koppeling)
```sql
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id),
    permission_id UUID NOT NULL REFERENCES permissions(id),
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by UUID,
    UNIQUE(role_id, permission_id)
);
```

**Kenmerken**:
- Many-to-many relatie
- Audit trail via `assigned_by`
- Voorkomen van duplicaten

#### 4. `user_roles` (Gebruiker-Rol Koppeling)
```sql
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(id),
    assigned_at TIMESTAMP DEFAULT NOW(),
    assigned_by UUID,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(user_id, role_id)
);
```

**Kenmerken**:
- Support voor tijdelijke rollen via `expires_at`
- Deactivatie zonder verwijdering via `is_active`
- Multiple roles per gebruiker

#### 5. `gebruikers` (Gebruikers - Dual System)
```sql
-- Legacy + RBAC
CREATE TABLE gebruikers (
    id UUID PRIMARY KEY,
    naam VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    wachtwoord_hash VARCHAR NOT NULL,
    rol VARCHAR DEFAULT 'gebruiker',  -- LEGACY
    role_id UUID,                      -- RBAC (niet gebruikt)
    is_actief BOOLEAN DEFAULT true,
    newsletter_subscribed BOOLEAN DEFAULT false,
    laatste_login TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**Transitie Status**:
- ⚠️ `rol` veld wordt **nog steeds gebruikt** voor JWT
- ⚠️ `role_id` is voorbereid maar niet actief
- ✅ `user_roles` tabel wordt gebruikt voor RBAC

---

## 📦 Model Structuren

### 1. RBACRole (models/role_rbac.go:8-24)
```go
type RBACRole struct {
    ID           string        `gorm:"type:uuid;primaryKey"`
    Name         string        `gorm:"type:varchar(100);uniqueIndex"`
    Description  string        `gorm:"type:text"`
    IsSystemRole bool          `gorm:"default:false"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    CreatedBy    *string       `gorm:"type:uuid"`
    
    // Relations
    Permissions []Permission   `gorm:"many2many:role_permissions"`
    Users       []Gebruiker    `gorm:"many2many:user_roles"`
}
```

### 2. Permission (models/role_rbac.go:26-42)
```go
type Permission struct {
    ID                 string        `gorm:"type:uuid;primaryKey"`
    Resource           string        `gorm:"type:varchar(100)"`
    Action             string        `gorm:"type:varchar(50)"`
    Description        string        `gorm:"type:text"`
    IsSystemPermission bool          `gorm:"default:false"`
    CreatedAt          time.Time
    UpdatedAt          time.Time
    
    // Relations
    Roles []RBACRole  `gorm:"many2many:role_permissions"`
}
```

### 3. UserRole (models/role_rbac.go:62-78)
```go
type UserRole struct {
    ID         string        `gorm:"type:uuid;primaryKey"`
    UserID     string        `gorm:"type:uuid;not null"`
    RoleID     string        `gorm:"type:uuid;not null"`
    AssignedAt time.Time     `gorm:"autoCreateTime"`
    AssignedBy *string       `gorm:"type:uuid"`
    ExpiresAt  *time.Time    `gorm:"type:timestamp"`
    IsActive   bool          `gorm:"default:true"`
    
    // Relations
    User Gebruiker  `gorm:"foreignKey:UserID"`
    Role RBACRole   `gorm:"foreignKey:RoleID"`
}
```

### 4. Gebruiker - Dual System (models/gebruiker.go:11-26)
```go
type Gebruiker struct {
    ID                   string
    Naam                 string
    Email                string
    WachtwoordHash       string        `json:"-"`
    Rol                  string        // LEGACY - still used!
    RoleID               *string       // RBAC - prepared but unused
    IsActief             bool
    NewsletterSubscribed bool
    LaatsteLogin         *time.Time
    CreatedAt            time.Time
    UpdatedAt            time.Time
    
    // RBAC Relations
    Roles []RBACRole  `gorm:"many2many:user_roles"`
}
```

**⚠️ Belangrijk**: Het `Rol` veld wordt **nog steeds gebruikt** in JWT tokens!

### 5. Legacy Role Constants (models/role.go:3-14)
```go
type Role string

const (
    RoleAdmin        Role = "admin"
    RoleChatOwner    Role = "owner"
    RoleChatAdmin    Role = "admin"    // Conflict!
    RoleChatMember   Role = "member"
    RoleDeelnemer    Role = "Deelnemer"
    RoleBegeleider   Role = "Begeleider"
    RoleVrijwilliger Role = "Vrijwilliger"
    RoleStaff        Role = "staff"
)
```

**⚠️ Opmerking**: `RoleAdmin` en `RoleChatAdmin` beide "admin" - potentiële conflict

---

## 🔧 Repository Laag

### Interface Structuur (repository/rbac_interfaces.go)

#### RBACRoleRepository Interface
```go
type RBACRoleRepository interface {
    Create(ctx context.Context, role *models.RBACRole) error
    GetByID(ctx context.Context, id string) (*models.RBACRole, error)
    GetByName(ctx context.Context, name string) (*models.RBACRole, error)
    List(ctx context.Context, limit, offset int) ([]*models.RBACRole, error)
    ListWithPermissions(ctx context.Context, limit, offset int) ([]*models.RBACRole, error)
    Update(ctx context.Context, role *models.RBACRole) error
    Delete(ctx context.Context, id string) error
    GetSystemRoles(ctx context.Context) ([]*models.RBACRole, error)
}
```

#### PermissionRepository Interface
```go
type PermissionRepository interface {
    Create(ctx context.Context, permission *models.Permission) error
    GetByID(ctx context.Context, id string) (*models.Permission, error)
    GetByResourceAction(ctx context.Context, resource, action string) (*models.Permission, error)
    List(ctx context.Context, limit, offset int) ([]*models.Permission, error)
    ListByResource(ctx context.Context, resource string) ([]*models.Permission, error)
    Update(ctx context.Context, permission *models.Permission) error
    Delete(ctx context.Context, id string) error
    GetSystemPermissions(ctx context.Context) ([]*models.Permission, error)
}
```

#### UserRoleRepository Interface - **Belangrijkst voor Permission Checks**
```go
type UserRoleRepository interface {
    Create(ctx context.Context, ur *models.UserRole) error
    GetByUserAndRole(ctx context.Context, userID, roleID string) (*models.UserRole, error)
    ListByUser(ctx context.Context, userID string) ([]*models.UserRole, error)
    ListActiveByUser(ctx context.Context, userID string) ([]*models.UserRole, error)
    Deactivate(ctx context.Context, id string) error
    
    // Kritieke methode voor permission checks!
    GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error)
}
```

### GetUserPermissions SQL Query (repository/user_role_repository.go:96-117)
```sql
SELECT
    ur.user_id,
    u.email,
    r.name as role_name,
    p.resource,
    p.action,
    rp.assigned_at as permission_assigned_at,
    ur.assigned_at as role_assigned_at
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE ur.user_id = ? 
  AND ur.is_active = true
  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
ORDER BY ur.user_id, r.name, p.resource, p.action
```

**Wat doet deze query**:
1. Haalt alle actieve rollen op voor een gebruiker
2. Voor elke rol, haalt alle permissions op
3. Filtert op actieve en niet-verlopen rollen
4. Retourneert een **platte lijst** van alle permissions

---

## ⚙️ Service Laag

### PermissionService (services/permission_service.go)

#### Core Functionaliteit

**1. HasPermission - De Belangrijkste Methode** (line 68-110)
```go
func (s *PermissionServiceImpl) HasPermission(
    ctx context.Context, 
    userID, resource, action string
) bool {
    // 1. Check Redis cache
    if s.cacheEnabled {
        if cached := s.getCachedPermission(userID, resource, action); cached != nil {
            return *cached
        }
    }
    
    // 2. Query database
    permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID)
    if err != nil {
        logger.Error("Fout bij ophalen user permissions", "user_id", userID)
        return false
    }
    
    // 3. Check in list
    hasPermission := s.checkPermissionInList(permissions, resource, action)
    
    // 4. Cache result
    if s.cacheEnabled {
        s.cachePermission(userID, resource, action, hasPermission)
    }
    
    return hasPermission
}
```

**Performance Optimalisatie**:
- Redis cache met 5 minuten TTL
- Fallback naar database bij cache miss
- Logging alleen bij permission denied

#### Role Management Methoden
```go
// Toewijzen van rollen
AssignRole(ctx, userID, roleID string, assignedBy *string) error

// Intrekken van rollen
RevokeRole(ctx, userID, roleID string) error

// Rol CRUD
CreateRole(ctx, role *models.RBACRole, createdBy *string) error
UpdateRole(ctx, role *models.RBACRole) error
DeleteRole(ctx, roleID string) error  // Beschermd tegen system roles

// Permission toewijzing aan rollen
AssignPermissionToRole(ctx, roleID, permissionID string, assignedBy *string) error
RevokePermissionFromRole(ctx, roleID, permissionID string) error
```

#### Cache Management
```go
// Cache invalidatie
InvalidateUserCache(userID string)      // Specifieke gebruiker
RefreshCache(ctx) error                 // Volledige cache clear
refreshUsersWithRole(ctx, roleID)       // Alle gebruikers met rol
```

**Cache Keys**: `perm:{userID}:{resource}:{action}`

### AuthService (services/auth_service.go)

#### JWT Token Generatie (line 263-287)
```go
func (s *AuthServiceImpl) generateToken(gebruiker *models.Gebruiker) (string, error) {
    claims := JWTClaims{
        Email: gebruiker.Email,
        Role:  gebruiker.Rol,  // ⚠️ LEGACY FIELD!
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   gebruiker.ID,  // User ID in Subject
            Issuer:    "dklemailservice",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}
```

**⚠️ Kritiek**: JWT gebruikt **legacy `Rol` veld**, niet RBAC rollen!

#### Token Validatie (line 132-178)
```go
func (s *AuthServiceImpl) ValidateToken(token string) (string, error) {
    token = strings.TrimPrefix(token, "Bearer ")
    
    parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("onverwachte signing methode: %v", token.Header["alg"])
        }
        return s.jwtSecret, nil
    })
    
    // Uitgebreide error handling
    // Returns: userID (from Subject claim)
}
```

#### Authenticatie Flow
1. **Login** → Verify password → Generate JWT + Refresh Token
2. **ValidateToken** → Parse JWT → Extract User ID
3. **GetUserFromToken** → Validate → Fetch User → Check Active
4. **RefreshAccessToken** → Validate Refresh Token → Generate New Tokens

---

## 🎮 Handler & Middleware Laag

### PermissionHandler (handlers/permission_handler.go)

#### API Endpoints
```go
// Permission Management
GET    /api/rbac/permissions              → ListPermissions
POST   /api/rbac/permissions              → CreatePermission
PUT    /api/rbac/permissions/:id          → UpdatePermission
DELETE /api/rbac/permissions/:id          → DeletePermission

// Role Management
GET    /api/rbac/roles                    → ListRoles
GET    /api/rbac/roles/:id                → GetRole
POST   /api/rbac/roles                    → CreateRole
PUT    /api/rbac/roles/:id                → UpdateRole
DELETE /api/rbac/roles/:id                → DeleteRole

// Role-Permission Assignment
PUT    /api/rbac/roles/:id/permissions             → UpdateRolePermissions (bulk)
POST   /api/rbac/roles/:id/permissions/:permId     → AddPermissionToRole
DELETE /api/rbac/roles/:id/permissions/:permId     → RemovePermissionFromRole
```

**Middleware Stack**:
```go
rbacGroup.Use(AuthMiddleware(h.authService))              // JWT validatie
rbacGroup.Use(AdminPermissionMiddleware(h.permissionService))  // Admin check
```

#### List Permissions Endpoint - Gegroepeerd (line 69-159)
```go
func (h *PermissionHandler) ListPermissions(c *fiber.Ctx) error {
    // Query parameters
    resourceFilter := c.Query("resource", "")
    actionFilter := c.Query("action", "")
    search := c.Query("search", "")
    groupByResource := c.QueryBool("group_by_resource", true)
    
    // Haal alle permissions
    permissions, err := h.permissionRepo.List(ctx, limit, offset)
    
    // Filter client-side voor performance
    // ...
    
    // Groepeer per resource voor frontend
    if groupByResource {
        grouped := make(map[string][]*models.Permission)
        for _, perm := range filteredPermissions {
            grouped[perm.Resource] = append(grouped[perm.Resource], perm)
        }
        
        // Return: { groups: [{resource, permissions, count}], total }
    }
}
```

**Features**:
- Client-side filtering
- Groepering per resource
- Search in resource/action/description
- Paginatie support

### Permission Middleware (handlers/permission_middleware.go)

#### Generic Permission Check
```go
func PermissionMiddleware(
    permissionService services.PermissionService, 
    resource, action string
) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID, ok := c.Locals("userID").(string)
        if !ok || userID == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Niet geautoriseerd"})
        }
        
        if !permissionService.HasPermission(c.Context(), userID, resource, action) {
            logger.Warn("Permission denied", "user_id", userID, 
                "resource", resource, "action", action)
            return c.Status(403).JSON(fiber.Map{"error": "Geen toegang"})
        }
        
        return c.Next()
    }
}
```

#### Convenience Middlewares
```go
// Admin-only access
AdminPermissionMiddleware(permissionService) 
    → PermissionMiddleware(permissionService, "admin", "access")

// Staff-level access
StaffPermissionMiddleware(permissionService)
    → PermissionMiddleware(permissionService, "staff", "access")

// Multiple permissions check
ResourcePermissionMiddleware(permissionService, ...PermissionCheck)
```

---

## 👥 Rollen Overzicht

### System Roles (V1_21 seeded)

| Rol | Naam | Type | Beschrijving |
|-----|------|------|--------------|
| 1 | `admin` | System | Volledige beheerder met toegang tot alle functies |
| 2 | `staff` | System | Ondersteunend personeel met beperkte beheerrechten |
| 3 | `user` | System | Standaard gebruiker |
| 4 | `owner` | System | Chat kanaal eigenaar |
| 5 | `chat_admin` | System | Chat kanaal beheerder |
| 6 | `member` | System | Chat kanaal lid |
| 7 | `deelnemer` | System | Evenement deelnemer |
| 8 | `begeleider` | System | Evenement begeleider |
| 9 | `vrijwilliger` | System | Evenement vrijwilliger |

### Rol Hierarchie & Permissions

#### Admin Role - Alle Permissions
```sql
-- Admin krijgt ALLE permissions (V1_21:70-75)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND r.is_system_role = true
```

**Implicatie**: Admin heeft automatisch alle nieuwe permissions

#### Staff Role - Read Access
```sql
-- Staff krijgt read permissions (V1_21:77-85)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'staff' 
  AND ((p.resource IN ('user', 'contact', 'aanmelding', 'newsletter', 
                        'email', 'chat', 'notification')
        AND p.action = 'read')
       OR (p.resource = 'staff' AND p.action = 'access'))
```

**Permissions**:
- ✅ Lezen: users, contacts, aanmeldingen, newsletters, emails, chat, notifications
- ✅ Staff access
- ❌ Geen write/delete

#### Chat Roles

**Owner** (V1_21:87-93):
- Alle chat permissions
- Kanaal management
- Moderatie rechten

**Chat Admin** (V1_21:95-102):
- Read, write, moderate
- Geen kanaal creatie

**Member** (V1_21:104-111):
- Read, write
- Geen moderatie

**User** (V1_21:113-120):
- Basic chat read/write

#### Event Roles (V1_21:122-123)
- `deelnemer`, `begeleider`, `vrijwilliger`
- **Geen speciale permissions**
- Gebruikt voor registratie categorisatie

---

## 🔐 Permissions Overzicht

### Permission Structuur
Elke permission bestaat uit:
- **Resource**: Het object/entiteit (bijv. `contact`, `email`)
- **Action**: De operatie (bijv. `read`, `write`, `delete`)
- **Description**: Menselijk leesbare uitleg

### Resources en Actions (V1_21:19-66)

#### 1. Contact Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `contact` | `read` | Contactformulieren bekijken |
| `contact` | `write` | Contactformulieren bewerken |
| `contact` | `delete` | Contactformulieren verwijderen |

#### 2. Registration Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `aanmelding` | `read` | Aanmeldingen bekijken |
| `aanmelding` | `write` | Aanmeldingen bewerken |
| `aanmelding` | `delete` | Aanmeldingen verwijderen |

#### 3. Newsletter Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `newsletter` | `read` | Nieuwsbrieven bekijken |
| `newsletter` | `write` | Nieuwsbrieven aanmaken/bewerken |
| `newsletter` | `send` | Nieuwsbrieven verzenden |
| `newsletter` | `delete` | Nieuwsbrieven verwijderen |

#### 4. Email Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `email` | `read` | Inkomende emails bekijken |
| `email` | `write` | Emails bewerken (markeren als verwerkt) |
| `email` | `delete` | Emails verwijderen |
| `email` | `fetch` | Nieuwe emails ophalen |

#### 5. Admin Email
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `admin_email` | `send` | Emails verzenden namens admin |

#### 6. User Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `user` | `read` | Gebruikers bekijken |
| `user` | `write` | Gebruikers aanmaken/bewerken |
| `user` | `delete` | Gebruikers verwijderen |
| `user` | `manage_roles` | Gebruikersrollen beheren |

#### 7. Chat Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `chat` | `read` | Chat kanalen en berichten bekijken |
| `chat` | `write` | Berichten verzenden |
| `chat` | `manage_channel` | Kanalen aanmaken/beheren |
| `chat` | `moderate` | Berichten modereren |

#### 8. Notification Management
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `notification` | `read` | Notificaties bekijken |
| `notification` | `write` | Notificaties aanmaken |
| `notification` | `delete` | Notificaties verwijderen |

#### 9. System Access
| Resource | Action | Beschrijving |
|----------|--------|--------------|
| `system` | `admin` | Volledige systeemtoegang |
| `admin` | `access` | Volledige admin toegang |
| `staff` | `access` | Toegang tot staff functies |

### Permission Matrix

| Rol | Contact | Aanmelding | Newsletter | Email | User | Chat | Admin |
|-----|---------|-----------|------------|-------|------|------|-------|
| **Admin** | ✅ CRUD | ✅ CRUD | ✅ CRUD+Send | ✅ CRUD+Fetch | ✅ CRUD+Roles | ✅ Full | ✅ Full |
| **Staff** | 📖 Read | 📖 Read | 📖 Read | 📖 Read | 📖 Read | 📖 Read | ❌ None |
| **User** | ❌ None | ❌ None | ❌ None | ❌ None | ❌ None | 📝 Read/Write | ❌ None |
| **Owner** | ❌ None | ❌ None | ❌ None | ❌ None | ❌ None | ✅ Full | ❌ None |
| **Chat Admin** | ❌ None | ❌ None | ❌ None | ❌ None | ❌ None | 📝 Read/Write/Moderate | ❌ None |
| **Member** | ❌ None | ❌ None | ❌ None | ❌ None | ❌ None | 📝 Read/Write | ❌ None |

---

## 🔄 Data Flow

### Permission Check Flow

```
1. HTTP Request → API Endpoint
                     ↓
2. AuthMiddleware → ValidateToken → Extract User ID
                     ↓
3. PermissionMiddleware → HasPermission(userID, resource, action)
                     ↓
4. PermissionService:
   ├─→ Check Redis Cache
   │   ├─→ Cache Hit → Return Result
   │   └─→ Cache Miss ↓
   ├─→ Query Database (GetUserPermissions)
   │   ├─→ JOIN user_roles, roles, role_permissions, permissions
   │   └─→ Return List<UserPermission>
   ├─→ Check Permission in List
   ├─→ Cache Result (5 min TTL)
   └─→ Return Boolean
                     ↓
5. If True → c.Next() → Handler Execution
   If False → 403 Forbidden
```

### Role Assignment Flow

```
1. Admin API Call → AssignRole(userID, roleID)
                     ↓
2. PermissionService:
   ├─→ Validate Role Exists
   ├─→ Check User Doesn't Have Role
   ├─→ Create UserRole Record
   └─→ InvalidateUserCache(userID)
                     ↓
3. Redis Cache:
   ├─→ Find all keys: perm:{userID}:*
   └─→ Delete all matching keys
                     ↓
4. Next Request → Cache Miss → Fresh Database Query
```

### Login → Permission Check Flow

```
1. Login Request
   ├─→ AuthService.Login(email, password)
   ├─→ Verify Password
   ├─→ Generate JWT Token (with legacy rol field!)
   └─→ Return Access Token + Refresh Token
                     ↓
2. Subsequent Request with Token
   ├─→ AuthMiddleware.ValidateToken()
   ├─→ Extract UserID from JWT Subject
   ├─→ Store in c.Locals("userID")
   └─→ c.Next()
                     ↓
3. PermissionMiddleware
   ├─→ Get UserID from Locals
   ├─→ PermissionService.HasPermission()
   │   ├─→ Query user_roles table (not JWT!)
   │   └─→ Check against RBAC permissions
   └─→ Allow/Deny
```

**⚠️ Belangrijk**: JWT token heeft legacy `rol` field, maar permission checks gebruiken RBAC uit database!

---

## 💾 Cache Strategie

### Redis Cache Implementation

#### Cache Keys
```
Pattern: perm:{userID}:{resource}:{action}
Example: perm:123e4567-e89b-12d3-a456-426614174000:contact:read
```

#### Cache TTL
- **5 minuten** per permission check
- Korter dan typische sessie duur
- Balance tussen performance en data freshness

#### Cache Operations

**1. Set Cache** (services/permission_service.go:351-370)
```go
func (s *PermissionServiceImpl) cachePermission(
    userID, resource, action string, 
    hasPermission bool
) {
    cacheKey := fmt.Sprintf("perm:%s:%s:%s", userID, resource, action)
    data, _ := json.Marshal(hasPermission)
    s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
}
```

**2. Get Cache** (services/permission_service.go:323-348)
```go
func (s *PermissionServiceImpl) getCachedPermission(
    userID, resource, action string
) *bool {
    cacheKey := fmt.Sprintf("perm:%s:%s:%s", userID, resource, action)
    val, err := s.redisClient.Get(ctx, cacheKey).Result()
    if err == redis.Nil {
        return nil  // Cache miss
    }
    // Unmarshal and return
}
```

**3. Invalidate User Cache** (services/permission_service.go:373-397)
```go
func (s *PermissionServiceImpl) InvalidateUserCache(userID string) {
    pattern := fmt.Sprintf("perm:%s:*", userID)
    keys, _ := s.redisClient.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        s.redisClient.Del(ctx, keys...)
    }
}
```

**4. Full Cache Refresh** (services/permission_service.go:399-424)
```go
func (s *PermissionServiceImpl) RefreshCache(ctx context.Context) error {
    pattern := "perm:*"
    keys, _ := s.redisClient.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        s.redisClient.Del(ctx, keys...)
    }
}
```

### Cache Invalidation Triggers

| Event | Invalidation Scope | Method |
|-------|-------------------|--------|
| Role toegewezen aan user | Specifieke user | `InvalidateUserCache(userID)` |
| Role ingetrokken van user | Specifieke user | `InvalidateUserCache(userID)` |
| Permission toegevoegd aan role | Alle users met die role | `refreshUsersWithRole(roleID)` |
| Permission verwijderd van role | Alle users met die role | `refreshUsersWithRole(roleID)` |
| Role verwijderd | Alle users | `RefreshCache()` |
| Role bijgewerkt | Alle users met die role | `refreshUsersWithRole(roleID)` |

### Performance Metrics

**Without Cache**:
- Database query per permission check
- ~10-50ms per check
- High database load

**With Cache**:
- Redis lookup per permission check
- ~1-2ms per check (cache hit)
- 10-50ms on cache miss
- **95%+ cache hit rate** verwacht voor normale gebruik

---

## 🔒 Beveiligings Aspecten

### 1. System Protection

#### System Roles (Cannot Delete)
```go
// In permission_service.go:226-236
func (s *PermissionServiceImpl) DeleteRole(ctx context.Context, roleID string) error {
    role, _ := s.rbacRoleRepo.GetByID(ctx, roleID)
    
    if role.IsSystemRole {
        return fmt.Errorf("kan systeemrol niet verwijderen")
    }
    // ... rest of deletion
}
```

#### System Permissions (Cannot Delete)
```go
// In permission_handler.go:293-298
if permission.IsSystemPermission {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
        "error": "Kan systeempermission niet verwijderen",
    })
}
```

### 2. Audit Trail

Alle wijzigingen worden gelogd met:
- **assigned_by**: Wie de wijziging heeft gemaakt
- **assigned_at**: Wanneer de wijziging is gemaakt
- **created_by**: Wie de entiteit heeft aangemaakt

```go
// UserRole audit
type UserRole struct {
    AssignedAt time.Time
    AssignedBy *string   // UUID van admin die rol toewees
    // ...
}

// RolePermission audit
type RolePermission struct {
    AssignedAt time.Time
    AssignedBy *string   // UUID van admin die permission toewees
    // ...
}
```

### 3. Temporal Access Control

```go
type UserRole struct {
    ExpiresAt  *time.Time   // Rol verloopt automatisch
    IsActive   bool         // Kan handmatig gedeactiveerd worden
    // ...
}
```

**Query Logic**:
```sql
WHERE ur.is_active = true
  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
```

### 4. JWT Security

**Token Structure**:
```go
type JWTClaims struct {
    Email string            // User email
    Role  string            // Legacy role (for backward compat)
    RegisteredClaims {
        Subject:   userID   // Primary user identifier
        ExpiresAt: time     // Token expiry (default 20min)
        IssuedAt:  time     // Issue time
        Issuer:    string   // "dklemailservice"
    }
}
```

**Security Features**:
- HMAC-SHA256 signing
- Short expiry (20 minutes default)
- Refresh token rotation (7 days)
- Token revocation support

### 5. API Security Layers

```go
// Layer 1: Authentication
rbacGroup.Use(AuthMiddleware(authService))
    → JWT validation
    → User ID extraction

// Layer 2: Authorization
rbacGroup.Use(AdminPermissionMiddleware(permissionService))
    → RBAC permission check
    → Database-based (not JWT-based!)

// Layer 3: Business Logic
Handler → Additional checks, validation
```

### 6. Input Validation

**Permission Middleware**:
```go
// Always check userID presence
userID, ok := c.Locals("userID").(string)
if !ok || userID == "" {
    return 401 Unauthorized
}
```

**Handler Level**:
```go
// Validate required fields
if req.Resource == "" || req.Action == "" {
    return 400 Bad Request
}

// Check duplicates before insert
existing, _ := h.permissionRepo.GetByResourceAction(ctx, resource, action)
if existing != nil {
    return 409 Conflict
}
```

### 7. Error Handling & Logging

**Permission Denied Logging**:
```go
if !hasPermission {
    logger.Warn("Permission denied",
        "user_id", userID,
        "resource", resource,
        "action", action,
        "path", c.Path(),
        "method", c.Method())
}
```

**Security Benefits**:
- Audit trail van toegangspogingen
- Detectie van ongeautoriseerde toegang
- Debugging ondersteuning

---

## 📊 Migratie Status

### Uitgevoerde Migraties

#### V1.0.0 - Initial Setup (database/migrations.go:53-106)
```go
// AutoMigrate tables including RBAC models
err = tx.AutoMigrate(
    &models.Gebruiker{},
    &models.RBACRole{},      // ✅
    &models.Permission{},    // ✅
    &models.RolePermission{}, // ✅
    &models.UserRole{},      // ✅
)
```

**Status**: ✅ Tabellen aangemaakt

#### V1.16.0 - Chat Tables (database/migrations/V1_16__create_chat_tables.sql)
```sql
CREATE TABLE chat_channel_participants (
    role TEXT DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member'))
);
```

**Status**: ✅ Chat roles in participants tabel

#### V1.21.0 - RBAC Data Seeding (database/migrations/V1_21__seed_rbac_data.sql)
```sql
-- 9 System Roles
INSERT INTO roles (name, description, is_system_role) VALUES ...

-- 65+ Permissions across 11 resources  
INSERT INTO permissions (resource, action, description, is_system_permission) VALUES ...

-- Role-Permission mappings
INSERT INTO role_permissions (role_id, permission_id) SELECT ...
```

**Status**: ✅ RBAC data geseed

### Migratie Database View

De migratie gebruikt een view voor gemakkelijke permission queries (impliciete aanname):

```sql
-- User Permissions View (conceptueel)
CREATE OR REPLACE VIEW user_permissions_view AS
SELECT 
    ur.user_id,
    u.email,
    r.name as role_name,
    p.resource,
    p.action,
    p.description
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE ur.is_active = true
  AND (ur.expires_at IS NULL OR ur.expires_at > NOW());
```

### Pending Migrations

⚠️ **Niet Uitgevoerd**:
1. Migratie van legacy `gebruikers.rol` naar `user_roles` tabel
2. Verwijdering van `gebruikers.rol` veld (breaking change)
3. Migratie van JWT claims naar RBAC

---

## 🔄 Legacy vs RBAC

### Dual System Overview

| Aspect | Legacy System | RBAC System |
|--------|--------------|-------------|
| **Storage** | `gebruikers.rol` (string) | `user_roles` table |
| **JWT Token** | ✅ Gebruikt | ❌ Niet gebruikt |
| **Permission Check** | ❌ Simpele string compare | ✅ Database query |
| **Flexibility** | ❌ Fixed roles | ✅ Dynamic roles + permissions |
| **Granularity** | ❌ Role-level | ✅ Permission-level |
| **Audit Trail** | ❌ Geen | ✅ Volledig |
| **Active Use** | ⚠️ JWT claims only | ✅ All API endpoints |

### Current Usage Analysis

#### Legacy System Still Used In:

**1. JWT Token Generation** (auth_service.go:263-287)
```go
claims := JWTClaims{
    Email: gebruiker.Email,
    Role:  gebruiker.Rol,  // ⚠️ LEGACY!
    // ...
}
```

**2. Database Migration Seeding** (database/migrations.go:108-136)
```go
adminUser := &models.Gebruiker{
    Email: "admin@dekoninklijkeloop.nl",
    Rol:   "admin",  // ⚠️ LEGACY!
    // ...
}
```

**3. Model Definition** (models/gebruiker.go:16)
```go
type Gebruiker struct {
    Rol      string  `gorm:"default:'gebruiker';index"`  // ⚠️ LEGACY!
    RoleID   *string `gorm:"type:uuid"`  // Prepared but unused
    // ...
}
```

#### RBAC System Used In:

**1. All Permission Checks** (via PermissionMiddleware)
```go
// Uses RBAC database tables
permissionService.HasPermission(ctx, userID, resource, action)
```

**2. All RBAC API Endpoints**
```
/api/rbac/permissions
/api/rbac/roles
/api/rbac/roles/:id/permissions
```

**3. Role Assignment**
```go
// Manages user_roles table
permissionService.AssignRole(ctx, userID, roleID, assignedBy)
```

### Migration Path

#### Phase 1: Current (✅ Completed)
- RBAC tables created
- Data seeded
- Permission checks active
- Legacy jwt <bron> still active

#### Phase 2: Transition (⚠️ In Progress)
- Both systems operational
- New roles assigned via RBAC
- Legacy rol field maintained for JWT

#### Phase 3: Full Migration (❌ Not Started)
1. Migrate existing users: Copy `gebruikers.rol` → `user_roles`
2. Update JWT generation to use RBAC
3. Test all endpoints thoroughly
4. Deprecate `gebruikers.rol` field

#### Phase 4: Cleanup (❌ Future)
1. Remove `gebruikers.rol` column (breaking!)
2. Remove Role constants (models/role.go)
3. Update documentation

---

## ✅ Best Practices

### Huidige Implementatie Sterke Punten

#### 1. ✅ Separation of Concerns
```
Models → Repositories → Services → Handlers
```
Duidelijke lagen met specifieke verantwoordelijkheden

#### 2. ✅ Interface-Based Design
```go
type RBACRoleRepository interface { ... }
type PermissionRepository interface { ... }
```
Maakt testing en dependency injection mogelijk

#### 3. ✅ Context Propagation
```go
func (s *Service) HasPermission(ctx context.Context, ...) bool
```
Alle database operaties ondersteunen context voor timeout/cancellation

#### 4. ✅ Comprehensive Logging
```go
logger.Warn("Permission denied", "user_id", userID, "resource", resource)
```
Gestructureerde logging voor audit en debugging

#### 5. ✅ Redis Caching
```
Performance optimization met cache invalidatie strategie
```

#### 6. ✅ System Protection
```go
if role.IsSystemRole { return error }
```
Systeemrollen en permissions zijn beschermd

#### 7. ✅ Audit Trail
```go
AssignedBy *string  // Wie heeft de actie uitgevoerd
AssignedAt time.Time  // Wanneer
```

#### 8. ✅ Temporal Access
```go
ExpiresAt *time.Time  // Automatische vervaldatum
IsActive  bool  // Handmatige deactivatie
```

#### 9. ✅ Transaction Safety
```go
err = m.db.Transaction(func(tx *gorm.DB) error {
    // Multiple operations in single transaction
})
```

#### 10. ✅ Comprehensive API
```
GET/POST/PUT/DELETE endpoints voor alle RBAC entiteiten
Bulk operations support
```

---

## 🎯 Aanbevelingen

### Prioriteit 1: Critical (Immediate Action Required)

#### 1. ⚠️ Resolveer Legacy/RBAC Conflict
**Probleem**: JWT gebruikt legacy `rol` field, maar permission checks gebruiken RBAC
```go
// Current: Inconsistent!
JWT Claims: { role: "admin" }  // From gebruikers.rol
Permission Check: Query user_roles table  // RBAC
```

**Aanbeveling**:
```go
// Optie A: Migreer JWT naar RBAC
claims := JWTClaims{
    Email:   gebruiker.Email,
    Roles:   getUserRoleNames(gebruiker.ID),  // From user_roles
    Subject: gebruiker.ID,
}

// Optie B: Dual claims (transitie periode)
claims := JWTClaims{
    Email:      gebruiker.Email,
    Role:       gebruiker.Rol,  // Legacy (deprecated)
    RBACRoles:  getUserRoleNames(gebruiker.ID),  // New
}
```

#### 2. ⚠️ Fix Role Constant Conflict
**Probleem**: `RoleAdmin` en `RoleChatAdmin` beide = "admin"
```go
const (
    RoleAdmin      Role = "admin"
    RoleChatAdmin  Role = "admin"  // CONFLICT!
)
```

**Aanbeveling**:
```go
const (
    RoleAdmin       Role = "admin"
    RoleChatOwner   Role = "chat_owner"   // Renamed
    RoleChatAdmin   Role = "chat_admin"   // Unique value
    RoleChatMember  Role = "chat_member"  // Consistent naming
)
```

#### 3. 🔍 Add Permission Check to All Protected Routes
**Scan bestaande endpoints** en voeg toe waar nodig:
```go
// Example: Email endpoints should use permissions
emailGroup := app.Group("/api/email")
emailGroup.Use(AuthMiddleware(authService))
emailGroup.Use(PermissionMiddleware(permService, "email", "read"))
```

### Prioriteit 2: Important (Near-term Implementation)

#### 4. 📊 Database Migration Script
Create script om bestaande users te migreren:
```sql
-- Migrate existing role field to user_roles table
INSERT INTO user_roles (user_id, role_id, is_active)
SELECT 
    g.id,
    r.id,
    true
FROM gebruikers g
JOIN roles r ON LOWER(r.name) = LOWER(g.rol)
WHERE NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = g.id AND ur.role_id = r.id
);
```

#### 5. 🧪 Unit Tests voor Permission Service
```go
func TestHasPermission_CacheHit(t *testing.T) { ... }
func TestHasPermission_CacheMiss(t *testing.T) { ... }
func TestAssignRole_DuplicatePrevention(t *testing.T) { ... }
func TestDeleteRole_SystemRoleProtection(t *testing.T) { ... }
```

#### 6. 📝 API Documentation
Genereer OpenAPI/Swagger docs voor RBAC endpoints:
```yaml
/api/rbac/roles:
  get:
    summary: List all roles
    security:
      - bearerAuth: []
    parameters:
      - name: limit
        in: query
        schema:
          type: integer
```

#### 7. 🎨 Frontend Permission Helper
```typescript
// usePermission hook
function usePermission(resource: string, action: string): boolean {
    const { user } = useAuth();
    const permissions = usePermissions(user.id);
    return permissions.some(p => 
        p.resource === resource && p.action === action
    );
}

// Usage in component
if (usePermission('user', 'delete')) {
    return <DeleteButton />;
}
```

### Prioriteit 3: Nice to Have (Future Enhancement)

#### 8. 📈 Metrics & Monitoring
```go
// Prometheus metrics
permissionCheckTotal := prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "permission_checks_total",
    },
    []string{"resource", "action", "result"},
)

permissionCheckDuration := prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "permission_check_duration_seconds",
    },
    []string{"cache_hit"},
)
```

#### 9. 🔔 Permission Change Notifications
```go
// Notify gebruiker when their permissions change
type PermissionChangeEvent struct {
    UserID      string
    ChangeType  string  // "role_assigned", "role_revoked", etc.
    Details     map[string]interface{}
    Timestamp   time.Time
}
```

#### 10. 🎭 Role Templates
```go
// Predefined role templates voor snelle setup
type RoleTemplate struct {
    Name        string
    Permissions []string  // Permission IDs
}

var StaffTemplate = RoleTemplate{
    Name: "Standard Staff",
    Permissions: []string{
        "contact:read",
        "email:read",
        // ...
    },
}
```

#### 11. 📊 Admin Dashboard Enhancements
- Visual permission matrix (role × resource grid)
- Bulk role assignment
- Permission usage analytics
- Audit log viewer

#### 12. 🔐 Two-Factor Role Elevation
```go
// Voor gevoelige operations
func RequireTwoFactorForRole(roleID string) middleware {
    // Check if user has confirmed 2FA recently
    // For critical role changes
}
```

---

## 📚 Appendix

### Quick Reference Commands

#### Check User Permissions
```sql
-- Get all permissions for a user
SELECT DISTINCT p.resource, p.action, p.description
FROM user_roles ur
JOIN role_permissions rp ON ur.role_id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE ur.user_id = 'USER_ID_HERE'
  AND ur.is_active = true
  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
ORDER BY p.resource, p.action;
```

#### Check Role Permissions
```sql
-- Get all permissions for a role
SELECT p.resource, p.action, p.description
FROM role_permissions rp
JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = 'ROLE_ID_HERE'
ORDER BY p.resource, p.action;
```

#### Find Users with Specific Permission
```sql
-- Find all users with specific permission
SELECT DISTINCT g.email, g.naam, r.name as role_name
FROM gebruikers g
JOIN user_roles ur ON g.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE p.resource = 'RESOURCE_NAME'
  AND p.action = 'ACTION_NAME'
  AND ur.is_active = true
ORDER BY g.email;
```

### Database Indexes

**Recommended Indexes** (als niet al aanwezig):
```sql
-- Voor snelle permission lookups
CREATE INDEX idx_user_roles_user_active 
ON user_roles(user_id) WHERE is_active = true;

CREATE INDEX idx_role_permissions_role 
ON role_permissions(role_id);

CREATE INDEX idx_permissions_resource_action 
ON permissions(resource, action);

-- Voor cache invalidatie queries
CREATE INDEX idx_user_roles_role_active 
ON user_roles(role_id) WHERE is_active = true;
```

### Environment Variables

```bash
# JWT Configuration
JWT_SECRET=your_secret_key_here
JWT_TOKEN_EXPIRY=20m

# Redis Configuration (for caching)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=dklemailservice
DB_USER=postgres
DB_PASSWORD=
```

---

## 📝 Change Log

- **2025-11-01**: Initial comprehensive analysis
  - Analyzed all RBAC components
  - Documented dual system (Legacy/RBAC)
  - Identified conflicts and recommendations

---

## 👥 Contributors

**Analyst**: Kilo Code AI Assistant  
**Database**: DKL Email Service v1.21.0  
**Source Code Location**: c:/Users/jeffrey/Desktop/Githubmains/dklemailservice

---

**End of Document**