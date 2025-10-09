# Authentication & Authorization

## Overzicht

De DKL Email Service gebruikt een moderne authenticatie en autorisatie architectuur gebaseerd op JWT tokens en een Role-Based Access Control (RBAC) systeem.

## Authenticatie Systeem

### JWT Token Authenticatie

De service gebruikt JWT (JSON Web Tokens) voor stateless authenticatie.

#### Token Structuur

**Claims** ([`services/auth_service.go:36`](../../services/auth_service.go:36)):
```go
type JWTClaims struct {
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.RegisteredClaims
}
```

**Token Generatie** ([`services/auth_service.go:263`](../../services/auth_service.go:263)):
```go
claims := JWTClaims{
    Email: gebruiker.Email,
    Role:  gebruiker.Rol,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        NotBefore: jwt.NewNumericDate(time.Now()),
        Issuer:    "dklemailservice",
        Subject:   gebruiker.ID,
    },
}
```

#### Token Configuratie

**Environment Variables:**
```bash
JWT_SECRET=your-secret-key-here          # Secret voor token signing
JWT_TOKEN_EXPIRY=20m                     # Access token expiry (default: 20 minuten)
```

### Refresh Token Mechanisme

Voor langdurige sessies gebruikt de service refresh tokens.

#### Refresh Token Flow

1. **Login** - Gebruiker ontvangt access token (20 min) + refresh token (7 dagen)
2. **Access Token Expired** - Frontend detecteert expiry
3. **Refresh Request** - Frontend stuurt refresh token naar `/api/auth/refresh`
4. **Token Rotation** - Nieuwe access + refresh tokens worden gegenereerd
5. **Old Token Revoked** - Oude refresh token wordt ingetrokken

**Implementatie** ([`services/auth_service.go:358`](../../services/auth_service.go:358)):
```go
func (s *AuthServiceImpl) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
    // Valideer refresh token
    token, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
    if err != nil || token == nil || !token.IsValid() {
        return "", "", errors.New("ongeldige of verlopen refresh token")
    }
    
    // Genereer nieuwe tokens
    accessToken, err := s.generateToken(gebruiker)
    newRefreshToken, err := s.GenerateRefreshToken(ctx, gebruiker.ID)
    
    // Revoke oude token (token rotation voor security)
    s.refreshTokenRepo.RevokeToken(ctx, refreshToken)
    
    return accessToken, newRefreshToken, nil
}
```

### Wachtwoord Beveiliging

**Hashing met bcrypt** ([`services/auth_service.go:211`](../../services/auth_service.go:211)):
```go
func (s *AuthServiceImpl) HashPassword(wachtwoord string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(wachtwoord), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func (s *AuthServiceImpl) VerifyPassword(hash, wachtwoord string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(wachtwoord))
    return err == nil
}
```

**Bcrypt Cost:** Default cost (10) - balans tussen security en performance

## Middleware

### Auth Middleware

Valideert JWT tokens en extraheert gebruiker informatie.

**Implementatie** ([`handlers/middleware.go:12`](../../handlers/middleware.go:12)):
```go
func AuthMiddleware(authService services.AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Haal token uit Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Niet geautoriseerd",
                "code":  "NO_AUTH_HEADER",
            })
        }
        
        // Valideer Bearer token format
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Ongeldige Authorization header",
                "code":  "INVALID_AUTH_HEADER",
            })
        }
        
        // Valideer token
        token := parts[1]
        userID, err := authService.ValidateToken(token)
        if err != nil {
            // Specifieke error codes voor frontend handling
            errorCode := "INVALID_TOKEN"
            if strings.Contains(err.Error(), "expired") {
                errorCode = "TOKEN_EXPIRED"
            }
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Ongeldig token",
                "code":  errorCode,
            })
        }
        
        // Sla user ID en token op in context
        c.Locals("userID", userID)
        c.Locals("token", token)
        
        return c.Next()
    }
}
```

**Error Codes:**
- `NO_AUTH_HEADER` - Geen Authorization header aanwezig
- `INVALID_AUTH_HEADER` - Geen Bearer token format
- `TOKEN_EXPIRED` - Token is verlopen
- `TOKEN_MALFORMED` - Token heeft ongeldige structuur
- `TOKEN_SIGNATURE_INVALID` - Token signature is ongeldig
- `INVALID_TOKEN` - Algemene token validatie fout

### Admin Middleware

Controleert of gebruiker admin rechten heeft.

**Implementatie** ([`handlers/middleware.go:111`](../../handlers/middleware.go:111)):
```go
func AdminMiddleware(authService services.AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token, ok := c.Locals("token").(string)
        if !ok || token == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Niet geautoriseerd",
            })
        }
        
        gebruiker, err := authService.GetUserFromToken(ctx, token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Niet geautoriseerd",
            })
        }
        
        if gebruiker.Rol != "admin" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Geen toegang",
            })
        }
        
        c.Locals("gebruiker", gebruiker)
        return c.Next()
    }
}
```

### Staff Middleware

Staat zowel admin als staff rollen toe.

**Implementatie** ([`handlers/middleware.go:73`](../../handlers/middleware.go:73)):
```go
func StaffMiddleware(authService services.AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // ... token validatie ...
        
        if gebruiker.Rol != "admin" && gebruiker.Rol != "staff" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Geen toegang",
            })
        }
        
        return c.Next()
    }
}
```

## RBAC Systeem

### Architectuur

Het RBAC systeem bestaat uit vier hoofdcomponenten:

1. **Roles** - Groepen van permissies
2. **Permissions** - Granulaire toegangsrechten
3. **User-Role Assignments** - Koppeling gebruikers aan rollen
4. **Role-Permission Assignments** - Koppeling rollen aan permissies

### Database Schema

**Roles Table** ([`models/role_rbac.go:8`](../../models/role_rbac.go:8)):
```go
type RBACRole struct {
    ID           string    `gorm:"type:uuid;primaryKey"`
    Name         string    `gorm:"type:varchar(100);not null;uniqueIndex"`
    Description  string    `gorm:"type:text"`
    IsSystemRole bool      `gorm:"default:false"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    CreatedBy    *string   `gorm:"type:uuid"`
    
    Permissions []Permission `gorm:"many2many:role_permissions"`
    Users       []Gebruiker  `gorm:"many2many:user_roles"`
}
```

**Permissions Table** ([`models/role_rbac.go:27`](../../models/role_rbac.go:27)):
```go
type Permission struct {
    ID                 string    `gorm:"type:uuid;primaryKey"`
    Resource           string    `gorm:"type:varchar(100);not null"`
    Action             string    `gorm:"type:varchar(50);not null"`
    Description        string    `gorm:"type:text"`
    IsSystemPermission bool      `gorm:"default:false"`
    CreatedAt          time.Time
    UpdatedAt          time.Time
    
    Roles []RBACRole `gorm:"many2many:role_permissions"`
}
```

**User Roles Table** ([`models/role_rbac.go:62`](../../models/role_rbac.go:62)):
```go
type UserRole struct {
    ID         string     `gorm:"type:uuid;primaryKey"`
    UserID     string     `gorm:"type:uuid;not null"`
    RoleID     string     `gorm:"type:uuid;not null"`
    AssignedAt time.Time
    AssignedBy *string    `gorm:"type:uuid"`
    ExpiresAt  *time.Time
    IsActive   bool       `gorm:"default:true"`
    
    User Gebruiker `gorm:"foreignKey:UserID"`
    Role RBACRole  `gorm:"foreignKey:RoleID"`
}
```

### Permission Service

**Permission Check** ([`services/permission_service.go:68`](../../services/permission_service.go:68)):
```go
func (s *PermissionServiceImpl) HasPermission(ctx context.Context, userID, resource, action string) bool {
    // Probeer eerst uit cache
    if s.cacheEnabled {
        if cached := s.getCachedPermission(userID, resource, action); cached != nil {
            return *cached
        }
    }
    
    // Haal permissies uit database
    permissions, err := s.userRoleRepo.GetUserPermissions(ctx, userID)
    if err != nil {
        return false
    }
    
    hasPermission := s.checkPermissionInList(permissions, resource, action)
    
    // Cache resultaat (5 minuten)
    if s.cacheEnabled {
        s.cachePermission(userID, resource, action, hasPermission)
    }
    
    return hasPermission
}
```

### Redis Caching

Voor performance worden permission checks gecached in Redis.

**Cache Strategy:**
- **TTL:** 5 minuten
- **Key Format:** `perm:{userID}:{resource}:{action}`
- **Invalidatie:** Bij rol/permissie wijzigingen

**Cache Implementatie** ([`services/permission_service.go:324`](../../services/permission_service.go:324)):
```go
func (s *PermissionServiceImpl) cachePermission(userID, resource, action string, hasPermission bool) {
    cacheKey := fmt.Sprintf("perm:%s:%s:%s", userID, resource, action)
    data, _ := json.Marshal(hasPermission)
    s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
}

func (s *PermissionServiceImpl) InvalidateUserCache(userID string) {
    pattern := fmt.Sprintf("perm:%s:*", userID)
    keys, _ := s.redisClient.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        s.redisClient.Del(ctx, keys...)
    }
}
```

### Permission Middleware

**Resource-Based Access Control** ([`handlers/permission_middleware.go`](../../handlers/permission_middleware.go:1)):
```go
func PermissionMiddleware(permissionService services.PermissionService, resource, action string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID, ok := c.Locals("userID").(string)
        if !ok || userID == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Niet geautoriseerd",
            })
        }
        
        if !permissionService.HasPermission(c.Context(), userID, resource, action) {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Geen toegang tot deze resource",
                "required_permission": map[string]string{
                    "resource": resource,
                    "action":   action,
                },
            })
        }
        
        return c.Next()
    }
}
```

## Standaard Rollen

### Admin Role

**Permissies:**
- Volledige toegang tot alle resources
- Gebruikersbeheer
- Rol en permissie beheer
- Systeem configuratie

### Staff Role

**Permissies:**
- Lees toegang tot meeste resources
- Beperkte schrijf toegang
- Geen gebruikersbeheer
- Geen systeem configuratie

### User Role

**Permissies:**
- Basis toegang
- Eigen profiel beheer
- Lees toegang tot publieke resources

## Permission Resources

### Beschikbare Resources

| Resource | Beschrijving |
|----------|--------------|
| `users` | Gebruikersbeheer |
| `roles` | Rol beheer |
| `permissions` | Permissie beheer |
| `contacts` | Contact formulieren |
| `aanmeldingen` | Aanmeldingen |
| `emails` | Email beheer |
| `mail` | Inkomende emails |
| `chat` | Chat functionaliteit |
| `newsletters` | Nieuwsbrieven |
| `notifications` | Notificaties |
| `staff_access` | Staff toegang |

### Beschikbare Actions

| Action | Beschrijving |
|--------|--------------|
| `read` | Lees toegang |
| `create` | Aanmaken |
| `update` | Bijwerken |
| `delete` | Verwijderen |
| `manage` | Volledig beheer |

## API Endpoints

### Authentication Endpoints

```
POST   /api/auth/login           - Login met email/wachtwoord
POST   /api/auth/logout          - Logout
POST   /api/auth/refresh         - Refresh access token
GET    /api/auth/profile         - Haal gebruikersprofiel op (auth required)
POST   /api/auth/reset-password  - Wijzig wachtwoord (auth required)
```

### RBAC Endpoints

```
GET    /api/rbac/roles                    - Lijst van rollen (admin)
POST   /api/rbac/roles                    - Nieuwe rol aanmaken (admin)
GET    /api/rbac/roles/:id                - Rol details (admin)
PUT    /api/rbac/roles/:id                - Rol bijwerken (admin)
DELETE /api/rbac/roles/:id                - Rol verwijderen (admin)

GET    /api/rbac/permissions              - Lijst van permissies (admin)
POST   /api/rbac/roles/:id/permissions    - Permissie toekennen (admin)
DELETE /api/rbac/roles/:id/permissions/:permId - Permissie verwijderen (admin)

POST   /api/users/:id/roles               - Rol toekennen aan gebruiker (admin)
DELETE /api/users/:id/roles/:roleId       - Rol verwijderen van gebruiker (admin)
```

## Security Best Practices

### Token Security

1. **HTTPS Only** - Tokens alleen over HTTPS verzenden
2. **HttpOnly Cookies** - Optioneel voor extra security
3. **Short Expiry** - Access tokens 20 minuten
4. **Token Rotation** - Refresh tokens worden geroteerd
5. **Revocation** - Oude tokens worden ingetrokken

### Password Security

1. **Bcrypt Hashing** - Industry standard
2. **Salt per Password** - Automatisch door bcrypt
3. **Minimum Length** - Enforce in frontend
4. **Complexity Requirements** - Aanbevolen in frontend

### Rate Limiting

**Login Endpoint** ([`handlers/auth_handler.go:49`](../../handlers/auth_handler.go:49)):
```go
rateLimitKey := "login:" + loginData.Email
if !h.rateLimiter.Allow(rateLimitKey) {
    return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
        "error": "Te veel login pogingen, probeer het later opnieuw",
    })
}
```

**Configuratie:**
```bash
LOGIN_LIMIT_COUNT=5          # Max pogingen
LOGIN_LIMIT_PERIOD=300       # Periode in seconden (5 minuten)
```

## Frontend Integratie

### Login Flow

```typescript
// 1. Login request
const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, wachtwoord })
});

const { token, refresh_token, user } = await response.json();

// 2. Store tokens
localStorage.setItem('access_token', token);
localStorage.setItem('refresh_token', refresh_token);

// 3. Use token in requests
fetch('/api/protected-endpoint', {
    headers: {
        'Authorization': `Bearer ${token}`
    }
});
```

### Token Refresh Flow

```typescript
// Interceptor voor expired tokens
async function refreshAccessToken() {
    const refreshToken = localStorage.getItem('refresh_token');
    
    const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken })
    });
    
    if (response.ok) {
        const { token, refresh_token } = await response.json();
        localStorage.setItem('access_token', token);
        localStorage.setItem('refresh_token', refresh_token);
        return token;
    } else {
        // Redirect to login
        window.location.href = '/login';
    }
}
```

### Permission Check

```typescript
// Check permission in frontend
function hasPermission(user: User, resource: string, action: string): boolean {
    return user.permissions.some(p => 
        p.resource === resource && p.action === action
    );
}

// Conditional rendering
{hasPermission(user, 'users', 'manage') && (
    <AdminPanel />
)}
```

## Troubleshooting

### Common Issues

**Token Expired:**
```json
{
    "error": "Ongeldig token",
    "code": "TOKEN_EXPIRED"
}
```
**Oplossing:** Gebruik refresh token endpoint

**Invalid Signature:**
```json
{
    "error": "Ongeldig token",
    "code": "TOKEN_SIGNATURE_INVALID"
}
```
**Oplossing:** Controleer JWT_SECRET configuratie

**Permission Denied:**
```json
{
    "error": "Geen toegang tot deze resource",
    "required_permission": {
        "resource": "users",
        "action": "manage"
    }
}
```
**Oplossing:** Controleer gebruiker rollen en permissies

## Monitoring

### Metrics

- Login attempts (success/failure)
- Token refresh rate
- Permission check latency
- Cache hit rate
- Failed authorization attempts

### Logging

Alle authenticatie events worden gelogd:
- Login pogingen
- Token validatie fouten
- Permission denials
- Rol wijzigingen

Zie [Monitoring Guide](../guides/monitoring.md) voor details.