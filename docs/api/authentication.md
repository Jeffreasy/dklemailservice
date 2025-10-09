# Authentication API

Complete authenticatie API referentie met daadwerkelijke implementatie details.

## Endpoints

### Login

#### POST /api/auth/login

Authenticeer een gebruiker en ontvang access + refresh tokens.

**Request Body:**
```json
{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "your-password"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
    "user": {
        "id": "7157f3f6-da85-4058-9d38-19133ec93b03",
        "email": "admin@dekoninklijkeloop.nl",
        "naam": "Admin",
        "rol": "admin",
        "permissions": [
            {
                "resource": "users",
                "action": "manage"
            },
            {
                "resource": "contacts",
                "action": "manage"
            }
        ],
        "is_actief": true
    }
}
```

**Response (401 Unauthorized):**
```json
{
    "error": "Ongeldige inloggegevens"
}
```

**Response (403 Forbidden):**
```json
{
    "error": "Gebruiker is inactief"
}
```

**Response (429 Too Many Requests):**
```json
{
    "error": "Te veel login pogingen, probeer het later opnieuw"
}
```

**Rate Limiting:**
- 5 pogingen per email adres per 5 minuten

**Implementatie:** [`handlers/auth_handler.go:29`](../../handlers/auth_handler.go:29)

**Code Voorbeeld:**
```go
// Login handler implementatie
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
    var loginData models.GebruikerLogin
    if err := c.BodyParser(&loginData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Ongeldige login data",
        })
    }
    
    // Rate limiting
    rateLimitKey := "login:" + loginData.Email
    if !h.rateLimiter.Allow(rateLimitKey) {
        return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
            "error": "Te veel login pogingen, probeer het later opnieuw",
        })
    }
    
    // Authenticeer
    token, refreshToken, err := h.authService.Login(c.Context(), loginData.Email, loginData.Wachtwoord)
    if err != nil {
        switch err {
        case services.ErrInvalidCredentials:
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Ongeldige inloggegevens",
            })
        case services.ErrUserInactive:
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Gebruiker is inactief",
            })
        }
    }
    
    // Haal gebruiker en permissies op
    gebruiker, _ := h.authService.GetUserFromToken(c.Context(), token)
    permissions, _ := h.permissionService.GetUserPermissions(c.Context(), gebruiker.ID)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "token": token,
        "refresh_token": refreshToken,
        "user": gebruiker,
    })
}
```

**cURL Voorbeeld:**
```bash
curl -X POST https://api.dekoninklijkeloop.nl/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "your-password"
  }'
```

**JavaScript Voorbeeld:**
```javascript
const login = async (email, wachtwoord) => {
    const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, wachtwoord }),
    });
    
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error);
    }
    
    const { token, refresh_token, user } = await response.json();
    
    // Store tokens
    localStorage.setItem('access_token', token);
    localStorage.setItem('refresh_token', refresh_token);
    localStorage.setItem('user', JSON.stringify(user));
    
    return { token, refresh_token, user };
};
```

### Logout

#### POST /api/auth/logout

Logt de gebruiker uit en verwijdert de auth cookie.

**Response (200 OK):**
```json
{
    "message": "Logout succesvol"
}
```

**Implementatie:** [`handlers/auth_handler.go:171`](../../handlers/auth_handler.go:171)

**Code Voorbeeld:**
```go
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
    // Verwijder cookie
    c.ClearCookie("auth_token")
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Logout succesvol",
    })
}
```

**JavaScript Voorbeeld:**
```javascript
const logout = async () => {
    await fetch('/api/auth/logout', {
        method: 'POST',
    });
    
    // Clear local storage
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    
    // Redirect to login
    window.location.href = '/login';
};
```

### Token Refresh

#### POST /api/auth/refresh

Vernieuwt een verlopen access token met een refresh token.

**Request Body:**
```json
{
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4..."
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "bmV3IHJlZnJlc2ggdG9rZW4..."
}
```

**Response (401 Unauthorized):**
```json
{
    "error": "Ongeldige of verlopen refresh token",
    "code": "REFRESH_TOKEN_INVALID"
}
```

**Implementatie:** [`handlers/auth_handler.go:133`](../../handlers/auth_handler.go:133)

**Code Voorbeeld:**
```go
func (h *AuthHandler) HandleRefreshToken(c *fiber.Ctx) error {
    var refreshData struct {
        RefreshToken string `json:"refresh_token"`
    }
    
    if err := c.BodyParser(&refreshData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Ongeldige refresh token data",
        })
    }
    
    // Refresh tokens
    accessToken, newRefreshToken, err := h.authService.RefreshAccessToken(c.Context(), refreshData.RefreshToken)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Ongeldige of verlopen refresh token",
            "code": "REFRESH_TOKEN_INVALID",
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "token": accessToken,
        "refresh_token": newRefreshToken,
    })
}
```

**Token Rotation:**
Bij elke refresh wordt een nieuwe refresh token gegenereerd en de oude ingetrokken voor extra security.

**Service Implementatie:** [`services/auth_service.go:358`](../../services/auth_service.go:358)

**JavaScript Voorbeeld met Interceptor:**
```javascript
// Axios interceptor voor automatische token refresh
axios.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;
        
        if (error.response?.status === 401 && error.response?.data?.code === 'TOKEN_EXPIRED') {
            if (!originalRequest._retry) {
                originalRequest._retry = true;
                
                try {
                    const refreshToken = localStorage.getItem('refresh_token');
                    const response = await axios.post('/api/auth/refresh', {
                        refresh_token: refreshToken
                    });
                    
                    const { token, refresh_token } = response.data;
                    localStorage.setItem('access_token', token);
                    localStorage.setItem('refresh_token', refresh_token);
                    
                    // Retry original request met nieuwe token
                    originalRequest.headers['Authorization'] = `Bearer ${token}`;
                    return axios(originalRequest);
                } catch (refreshError) {
                    // Redirect naar login
                    window.location.href = '/login';
                    return Promise.reject(refreshError);
                }
            }
        }
        
        return Promise.reject(error);
    }
);
```

### Get Profile

#### GET /api/auth/profile

Haalt het profiel van de ingelogde gebruiker op.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "id": "7157f3f6-da85-4058-9d38-19133ec93b03",
    "naam": "Admin",
    "email": "admin@dekoninklijkeloop.nl",
    "rol": "admin",
    "permissions": [
        {
            "resource": "users",
            "action": "manage"
        },
        {
            "resource": "contacts",
            "action": "read"
        },
        {
            "resource": "contacts",
            "action": "update"
        }
    ],
    "roles": [
        {
            "id": "role-uuid",
            "name": "Administrator",
            "description": "Full system access",
            "assigned_at": "2024-01-01T00:00:00Z",
            "is_active": true
        }
    ],
    "is_actief": true,
    "laatste_login": "2024-03-20T15:04:05Z",
    "created_at": "2024-01-01T00:00:00Z"
}
```

**Response (401 Unauthorized):**
```json
{
    "error": "Niet geautoriseerd",
    "code": "NO_AUTH_HEADER"
}
```

**Implementatie:** [`handlers/auth_handler.go:242`](../../handlers/auth_handler.go:242)

**Code Voorbeeld:**
```go
func (h *AuthHandler) HandleGetProfile(c *fiber.Ctx) error {
    userID, ok := c.Locals("userID").(string)
    if !ok || userID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Niet geautoriseerd",
        })
    }
    
    // Haal gebruiker op
    gebruiker, err := h.authService.GetUser(c.Context(), userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Kon gebruiker niet ophalen",
        })
    }
    
    // Haal permissies op via RBAC
    permissions, _ := h.permissionService.GetUserPermissions(c.Context(), userID)
    
    // Haal rollen op
    userRoles, _ := h.permissionService.GetUserRoles(c.Context(), userID)
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "id": gebruiker.ID,
        "naam": gebruiker.Naam,
        "email": gebruiker.Email,
        "rol": gebruiker.Rol,
        "permissions": permissions,
        "roles": userRoles,
        "is_actief": gebruiker.IsActief,
        "laatste_login": gebruiker.LaatsteLogin,
        "created_at": gebruiker.CreatedAt,
    })
}
```

**JavaScript Voorbeeld:**
```javascript
const getProfile = async () => {
    const token = localStorage.getItem('access_token');
    
    const response = await fetch('/api/auth/profile', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    
    if (!response.ok) {
        throw new Error('Failed to fetch profile');
    }
    
    return await response.json();
};
```

### Reset Password

#### POST /api/auth/reset-password

Wijzigt het wachtwoord van de ingelogde gebruiker.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "huidig_wachtwoord": "old-password",
    "nieuw_wachtwoord": "new-password"
}
```

**Response (200 OK):**
```json
{
    "message": "Wachtwoord succesvol gewijzigd"
}
```

**Response (400 Bad Request):**
```json
{
    "error": "Huidig wachtwoord en nieuw wachtwoord zijn verplicht"
}
```

**Response (401 Unauthorized):**
```json
{
    "error": "Ongeldig huidig wachtwoord"
}
```

**Implementatie:** [`handlers/auth_handler.go:181`](../../handlers/auth_handler.go:181)

**Code Voorbeeld:**
```go
func (h *AuthHandler) HandleResetPassword(c *fiber.Ctx) error {
    userID, ok := c.Locals("userID").(string)
    if !ok || userID == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Niet geautoriseerd",
        })
    }
    
    var resetData struct {
        HuidigWachtwoord string `json:"huidig_wachtwoord"`
        NieuwWachtwoord  string `json:"nieuw_wachtwoord"`
    }
    
    if err := c.BodyParser(&resetData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Ongeldige wachtwoord reset data",
        })
    }
    
    // Haal gebruiker op
    gebruiker, _ := h.authService.GetUserFromToken(c.Context(), c.Locals("token").(string))
    
    // Verifieer huidig wachtwoord
    if !h.authService.VerifyPassword(gebruiker.WachtwoordHash, resetData.HuidigWachtwoord) {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Ongeldig huidig wachtwoord",
        })
    }
    
    // Reset wachtwoord
    if err := h.authService.ResetPassword(c.Context(), gebruiker.Email, resetData.NieuwWachtwoord); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Er is een fout opgetreden bij het resetten van het wachtwoord",
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Wachtwoord succesvol gewijzigd",
    })
}
```

**JavaScript Voorbeeld:**
```javascript
const resetPassword = async (currentPassword, newPassword) => {
    const token = localStorage.getItem('access_token');
    
    const response = await fetch('/api/auth/reset-password', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            huidig_wachtwoord: currentPassword,
            nieuw_wachtwoord: newPassword
        }),
    });
    
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error);
    }
    
    return await response.json();
};
```

## JWT Token Details

### Token Structure

**Claims:** [`services/auth_service.go:36`](../../services/auth_service.go:36)
```go
type JWTClaims struct {
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.RegisteredClaims
}
```

**Decoded Token Example:**
```json
{
    "email": "admin@dekoninklijkeloop.nl",
    "role": "admin",
    "exp": 1640001200,
    "iat": 1640000000,
    "nbf": 1640000000,
    "iss": "dklemailservice",
    "sub": "7157f3f6-da85-4058-9d38-19133ec93b03"
}
```

### Token Generation

**Implementatie:** [`services/auth_service.go:263`](../../services/auth_service.go:263)
```go
func (s *AuthServiceImpl) generateToken(gebruiker *models.Gebruiker) (string, error) {
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
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}
```

### Token Validation

**Implementatie:** [`services/auth_service.go:132`](../../services/auth_service.go:132)
```go
func (s *AuthServiceImpl) ValidateToken(token string) (string, error) {
    token = strings.TrimPrefix(token, "Bearer ")
    
    parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("onverwachte signing methode: %v", token.Header["alg"])
        }
        return s.jwtSecret, nil
    })
    
    if err != nil {
        // Specifieke error handling
        if errors.Is(err, jwt.ErrTokenExpired) {
            return "", ErrInvalidToken
        }
        return "", ErrInvalidToken
    }
    
    claims, ok := parsedToken.Claims.(*JWTClaims)
    if !ok || claims.Subject == "" {
        return "", ErrInvalidToken
    }
    
    return claims.Subject, nil
}
```

## Refresh Token Details

### Token Generation

**Implementatie:** [`services/auth_service.go:331`](../../services/auth_service.go:331)
```go
func (s *AuthServiceImpl) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
    // Genereer random token (32 bytes)
    tokenBytes := make([]byte, 32)
    if _, err := rand.Read(tokenBytes); err != nil {
        return "", err
    }
    token := base64.URLEncoding.EncodeToString(tokenBytes)
    
    // Sla op in database met 7 dagen expiry
    refreshToken := &models.RefreshToken{
        UserID:    userID,
        Token:     token,
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
        IsRevoked: false,
    }
    
    if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
        return "", err
    }
    
    return token, nil
}
```

### Token Rotation

Bij elke refresh wordt de oude token ingetrokken en een nieuwe gegenereerd.

**Implementatie:** [`services/auth_service.go:358`](../../services/auth_service.go:358)
```go
func (s *AuthServiceImpl) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
    // Valideer refresh token
    token, err := s.refreshTokenRepo.GetByToken(ctx, refreshToken)
    if err != nil || token == nil || !token.IsValid() {
        return "", "", errors.New("ongeldige of verlopen refresh token")
    }
    
    // Haal gebruiker op
    gebruiker, err := s.gebruikerRepo.GetByID(ctx, token.UserID)
    if err != nil || gebruiker == nil || !gebruiker.IsActief {
        return "", "", errors.New("gebruiker niet gevonden of inactief")
    }
    
    // Genereer nieuwe tokens
    accessToken, err := s.generateToken(gebruiker)
    newRefreshToken, err := s.GenerateRefreshToken(ctx, gebruiker.ID)
    
    // Revoke oude token (token rotation)
    s.refreshTokenRepo.RevokeToken(ctx, refreshToken)
    
    return accessToken, newRefreshToken, nil
}
```

## Security Best Practices

### Token Storage

**Aanbevolen:**
- Access token: Memory only (React state, Vuex store)
- Refresh token: HttpOnly cookie of secure localStorage

**Niet Aanbevolen:**
- Tokens in sessionStorage
- Tokens in URL parameters
- Tokens in localStorage zonder encryption

### Token Expiry

**Configuratie:**
```bash
JWT_TOKEN_EXPIRY=20m    # Access token (kort voor security)
# Refresh token: 7 dagen (hardcoded in service)
```

**Strategie:**
- Korte access token expiry (20 minuten)
- Langere refresh token expiry (7 dagen)
- Automatische refresh bij expiry
- Token rotation bij refresh

### Password Requirements

**Minimum Requirements:**
- Minimaal 8 karakters (aanbevolen in frontend)
- Mix van letters en cijfers (aanbevolen)
- Speciale karakters (optioneel)

**Hashing:**
- Algoritme: bcrypt
- Cost: 10 (default)
- Salt: Automatisch per wachtwoord

**Implementatie:** [`services/auth_service.go:211`](../../services/auth_service.go:211)

## Error Handling

### Error Codes

| Code | HTTP Status | Beschrijving |
|------|-------------|--------------|
| `NO_AUTH_HEADER` | 401 | Geen Authorization header |
| `INVALID_AUTH_HEADER` | 401 | Ongeldige header format |
| `TOKEN_EXPIRED` | 401 | Token is verlopen |
| `TOKEN_MALFORMED` | 401 | Token structuur ongeldig |
| `TOKEN_SIGNATURE_INVALID` | 401 | Token signature ongeldig |
| `INVALID_TOKEN` | 401 | Algemene token fout |
| `REFRESH_TOKEN_INVALID` | 401 | Refresh token ongeldig |

### Error Response Format

```json
{
    "error": "Beschrijving van de fout",
    "code": "ERROR_CODE"
}
```

## Testing

### Test Mode

Voor testing zonder echte email verzending:

**Via Header:**
```http
X-Test-Mode: true
```

**Via Body:**
```json
{
    "test_mode": true
}
```

### Test Credentials

**Default Admin:**
- Email: `admin@dekoninklijkeloop.nl`
- Wachtwoord: Zie database seed data

## Monitoring

Alle authenticatie events worden gelogd en gemonitord:

- Login attempts (success/failure)
- Token validatie fouten
- Password reset attempts
- Rate limit violations

Zie [Monitoring Guide](../guides/monitoring.md) voor details.

## Zie Ook

- [REST API Overview](./rest-api.md) - Complete API overzicht
- [Email Endpoints](./email-endpoints.md) - Email API
- [Admin Endpoints](./admin-endpoints.md) - Admin API
- [Authentication Architecture](../architecture/authentication-and-authorization.md) - Technische details