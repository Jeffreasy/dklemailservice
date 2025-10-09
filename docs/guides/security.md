# Security Guide

Complete security handleiding voor de DKL Email Service.

## Security Overzicht

De DKL Email Service implementeert meerdere beveiligingslagen:
- Input validatie en sanitization
- JWT authenticatie met refresh tokens
- RBAC (Role-Based Access Control)
- Rate limiting
- SMTP/IMAP TLS encryption
- Secure password hashing
- CORS configuratie
- Security headers
- Audit logging

## Input Validatie

### Email Validatie

**Implementatie:** [`handlers/email_handler.go:268`](../../handlers/email_handler.go:268)

```go
// Validate email format
if !strings.Contains(aanmelding.Email, "@") || !strings.Contains(aanmelding.Email, ".") {
    logger.Warn("Ongeldig email formaat",
        "email", aanmelding.Email,
        "remote_ip", c.IP())
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Ongeldig email adres",
    })
}
```

**Best Practices:**
- RFC 5321 compliant (max 254 karakters)
- Regex validatie voor format
- Domain verificatie (optioneel)
- Disposable email detectie (optioneel)

### Form Validatie

**Contact Formulier:** [`handlers/email_handler.go:85`](../../handlers/email_handler.go:85)

```go
if request.Naam == "" || request.Email == "" || request.Bericht == "" {
    logger.Warn("Onvolledig contact formulier",
        "naam", request.Naam,
        "email", request.Email,
        "bericht_empty", request.Bericht == "",
        "remote_ip", c.IP())
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Naam, email en bericht zijn verplicht",
    })
}

if !request.PrivacyAkkoord {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Je moet akkoord gaan met het privacybeleid",
    })
}
```

**Aanmelding Formulier:** [`handlers/email_handler.go:247`](../../handlers/email_handler.go:247)

```go
if aanmelding.Naam == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Naam is verplicht",
    })
}

if !aanmelding.Terms {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Je moet akkoord gaan met de voorwaarden",
    })
}
```

## Authenticatie Security

### Password Hashing

**Bcrypt Implementation:** [`services/auth_service.go:211`](../../services/auth_service.go:211)

```go
func (s *AuthServiceImpl) HashPassword(wachtwoord string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(wachtwoord), bcrypt.DefaultCost)
    if err != nil {
        logger.Error("Fout bij hashen wachtwoord", "error", err)
        return "", err
    }
    return string(hash), nil
}

func (s *AuthServiceImpl) VerifyPassword(hash, wachtwoord string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(wachtwoord))
    return err == nil
}
```

**Security Features:**
- **Algorithm:** bcrypt
- **Cost Factor:** 10 (default) - balans tussen security en performance
- **Salt:** Automatisch per wachtwoord
- **Rainbow Table Resistant:** Ja

### JWT Token Security

**Token Generation:** [`services/auth_service.go:263`](../../services/auth_service.go:263)

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

**Security Features:**
- **Algorithm:** HS256 (HMAC-SHA256)
- **Expiry:** 20 minuten (configureerbaar)
- **Issuer Validation:** "dklemailservice"
- **Subject:** User ID voor tracking

**Token Validation:** [`services/auth_service.go:132`](../../services/auth_service.go:132)

```go
func (s *AuthServiceImpl) ValidateToken(token string) (string, error) {
    token = strings.TrimPrefix(token, "Bearer ")
    
    parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("onverwachte signing methode: %v", token.Header["alg"])
        }
        return s.jwtSecret, nil
    })
    
    if err != nil {
        // Specifieke error handling voor verschillende JWT fouten
        if errors.Is(err, jwt.ErrTokenExpired) {
            logger.Warn("Token expired", "error", err)
        } else if errors.Is(err, jwt.ErrTokenMalformed) {
            logger.Error("Token malformed", "error", err)
        } else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
            logger.Error("Invalid signature", "error", err)
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

### Refresh Token Security

**Token Rotation:** [`services/auth_service.go:358`](../../services/auth_service.go:358)

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
    
    // Revoke oude token (token rotation voor extra security)
    s.refreshTokenRepo.RevokeToken(ctx, refreshToken)
    
    return accessToken, newRefreshToken, nil
}
```

**Security Features:**
- **Token Rotation:** Oude token wordt ingetrokken bij refresh
- **Expiry:** 7 dagen
- **Revocation:** Tokens kunnen worden ingetrokken
- **Random Generation:** Cryptographically secure random

## Rate Limiting

### Implementation

**Rate Limiter:** [`services/rate_limiter.go:49`](../../services/rate_limiter.go:49)

```go
type RateLimiter struct {
    mutex             sync.Mutex
    globalCounts      map[string]*counter
    ipCounts          map[string]map[string]*counter
    limits            map[string]RateLimit
    prometheusMetrics *PrometheusMetrics
    redisClient       *redis.Client
    useRedis          bool
}

func (rl *RateLimiter) AllowEmail(emailType, userID string) bool {
    if limit, exists := rl.limits[emailType]; exists {
        key := emailType
        if userID != "" && limit.PerIP {
            key = fmt.Sprintf("%s:%s", emailType, userID)
        }
        
        if rl.useRedis {
            return rl.allowWithRedis(key, limit, emailType, userID)
        } else {
            return rl.allowWithMemory(key, limit, emailType, userID)
        }
    }
    return true
}
```

### Redis-Based Rate Limiting

**Sliding Window Algorithm:** [`services/rate_limiter.go:121`](../../services/rate_limiter.go:121)

```go
func (r *RateLimiter) allowWithRedis(key string, limit RateLimit, emailType, userID string) bool {
    ctx := context.Background()
    
    // Sliding window met Redis sorted sets
    now := time.Now().Unix()
    windowStart := now - int64(limit.Period.Seconds())
    
    // Verwijder oude entries
    r.redisClient.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
    
    // Tel huidige entries
    count, err := r.redisClient.ZCard(ctx, key).Result()
    if err != nil {
        logger.Error("Redis ZCard failed", "error", err)
        return r.allowWithMemory(key, limit, emailType, userID)
    }
    
    // Check limiet
    if count >= int64(limit.Count) {
        logger.Warn("Rate limit exceeded",
            "operation", emailType,
            "limit", limit.Count,
            "current", count)
        
        if r.prometheusMetrics != nil {
            r.prometheusMetrics.RecordRateLimitExceeded(emailType, "per_user")
        }
        return false
    }
    
    // Voeg nieuwe entry toe
    r.redisClient.ZAdd(ctx, key, redis.Z{
        Score:  float64(now),
        Member: fmt.Sprintf("%d:%d", now, time.Now().Nanosecond()),
    })
    
    // Set TTL
    r.redisClient.Expire(ctx, key, limit.Period+time.Minute)
    
    return true
}
```

### Login Rate Limiting

**Implementation:** [`handlers/auth_handler.go:49`](../../handlers/auth_handler.go:49)

```go
// Rate limiting voor login pogingen
rateLimitKey := "login:" + loginData.Email
if !h.rateLimiter.Allow(rateLimitKey) {
    logger.Warn("Rate limit overschreden voor login", "email", loginData.Email)
    return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
        "error": "Te veel login pogingen, probeer het later opnieuw",
    })
}
```

**Configuratie:**
```bash
LOGIN_LIMIT_COUNT=5          # Max 5 pogingen
LOGIN_LIMIT_PERIOD=300       # Per 5 minuten
```

## SMTP/IMAP Security

### TLS Configuration

**SMTP Client:** [`services/smtp_client.go`](../../services/smtp_client.go:1)

```go
tlsConfig := &tls.Config{
    ServerName:         config.Host,
    MinVersion:         tls.VersionTLS12,
    CurvePreferences:   []tls.CurveID{tls.X25519, tls.CurveP256},
    CipherSuites: []uint16{
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
    },
}
```

**Security Features:**
- **Minimum TLS:** 1.2
- **Strong Ciphers:** ECDHE + AES-256-GCM
- **Perfect Forward Secrecy:** Ja
- **Certificate Validation:** Ja

### IMAP Security

**Mail Fetcher:** [`services/mail_fetcher.go`](../../services/mail_fetcher.go:1)

```go
func DialIMAP(host string, port int) (*client.Client, error) {
    imapAddr := fmt.Sprintf("%s:%d", host, port)
    c, err := client.DialTLS(imapAddr, &tls.Config{
        MinVersion: tls.VersionTLS12,
    })
    if err != nil {
        return nil, fmt.Errorf("kan niet verbinden met IMAP server: %w", err)
    }
    return c, nil
}
```

## CORS Security

### Configuration

**Implementatie:** [`main.go:296`](../../main.go:296)

```go
allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
    allowedOrigins = []string{
        "https://www.dekoninklijkeloop.nl",
        "https://dekoninklijkeloop.nl",
        "https://admin.dekoninklijkeloop.nl",
        "http://localhost:3000",
        "http://localhost:5173",
    }
}

app.Use(cors.New(cors.Config{
    AllowOrigins:     strings.Join(allowedOrigins, ","),
    AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Test-Mode",
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
    AllowCredentials: true,
    ExposeHeaders:    "Content-Length, Content-Type",
}))
```

**Production:**
```bash
ALLOWED_ORIGINS=https://www.dekoninklijkeloop.nl,https://admin.dekoninklijkeloop.nl
```

## Middleware Security

### Auth Middleware

**Token Validation:** [`handlers/middleware.go:12`](../../handlers/middleware.go:12)

```go
func AuthMiddleware(authService services.AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Niet geautoriseerd",
                "code":  "NO_AUTH_HEADER",
            })
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Ongeldige Authorization header",
                "code":  "INVALID_AUTH_HEADER",
            })
        }
        
        token := parts[1]
        userID, err := authService.ValidateToken(token)
        if err != nil {
            errorCode := "INVALID_TOKEN"
            if strings.Contains(err.Error(), "expired") {
                errorCode = "TOKEN_EXPIRED"
            }
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Ongeldig token",
                "code":  errorCode,
            })
        }
        
        c.Locals("userID", userID)
        c.Locals("token", token)
        
        return c.Next()
    }
}
```

### Permission Middleware

**RBAC Check:** [`handlers/permission_middleware.go`](../../handlers/permission_middleware.go:1)

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
            logger.Warn("Permission denied",
                "user_id", userID,
                "resource", resource,
                "action", action)
            
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

## Email Content Security

### Inkomende Email Sanitization

**HTML Sanitization:**
```go
import "github.com/microcosm-cc/bluemonday"

func sanitizeEmailContent(body string, contentType string) string {
    if strings.Contains(contentType, "text/html") {
        // Strict policy voor HTML cleaning
        p := bluemonday.StrictPolicy()
        return p.Sanitize(body)
    }
    return body
}
```

**XSS Prevention:**
- HTML tags worden gestript
- JavaScript wordt verwijderd
- Inline styles worden verwijderd
- Alleen veilige HTML tags toegestaan

### Template Security

**Template Rendering:**
```go
templateFuncs := template.FuncMap{
    "safeHTML": func(s string) template.HTML {
        return template.HTML(bluemonday.UGCPolicy().Sanitize(s))
    },
}
```

## RBAC Security

### Permission Caching

**Redis Cache:** [`services/permission_service.go:324`](../../services/permission_service.go:324)

```go
func (s *PermissionServiceImpl) cachePermission(userID, resource, action string, hasPermission bool) {
    if !s.cacheEnabled {
        return
    }
    
    ctx := context.Background()
    cacheKey := fmt.Sprintf("perm:%s:%s:%s", userID, resource, action)
    
    data, _ := json.Marshal(hasPermission)
    
    // Cache voor 5 minuten
    err := s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute).Err()
    if err != nil {
        logger.Error("Redis cache set error", "error", err)
    }
}
```

**Cache Invalidation:**
```go
func (s *PermissionServiceImpl) InvalidateUserCache(userID string) {
    if !s.cacheEnabled {
        return
    }
    
    ctx := context.Background()
    pattern := fmt.Sprintf("perm:%s:*", userID)
    
    keys, _ := s.redisClient.Keys(ctx, pattern).Result()
    if len(keys) > 0 {
        s.redisClient.Del(ctx, keys...)
    }
}
```

**Security Benefits:**
- Snellere permission checks
- Verminderde database load
- Automatische invalidatie bij wijzigingen
- TTL van 5 minuten voor freshness

## Logging & Monitoring

### Security Event Logging

**Login Attempts:**
```go
logger.Info("Login poging", "email", email)
logger.Warn("Ongeldige inloggegevens", "email", email)
logger.Warn("Inactieve gebruiker", "email", email)
```

**Permission Denials:**
```go
logger.Warn("Permission denied",
    "user_id", userID,
    "resource", resource,
    "action", action,
    "permissions_count", len(permissions))
```

**Rate Limit Violations:**
```go
logger.Warn("Rate limit overschreden",
    "operation", emailType,
    "limit", limit.Count,
    "current_count", count)
```

### Sensitive Data Masking

**Email Masking:**
```go
func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "invalid-email"
    }
    username := parts[0]
    domain := parts[1]
    
    if len(username) <= 3 {
        return username[:1] + "***@" + domain
    }
    return username[:2] + "***@" + domain
}

// Gebruik in logging
logger.Info("email processed",
    "email", maskEmail(email),
)
```

## Environment Security

### Secret Management

**Nooit Committen:**
```gitignore
.env
.env.local
.env.production
*.pem
*.key
*.crt
```

**Environment Variables:**
```bash
# ✅ Correct: Via environment
export JWT_SECRET="very-secure-random-string"

# ❌ Incorrect: Hardcoded
const jwtSecret = "my-secret"
```

### Production Secrets

**Render:**
- Gebruik Render's environment variables
- Enable "Sync" voor shared secrets
- Gebruik secret groups

**Docker:**
```bash
# Via environment file
docker run --env-file .env.production app

# Via secrets
docker secret create jwt_secret jwt_secret.txt
docker service create --secret jwt_secret app
```

## Security Headers

### HTTP Security Headers

```go
func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("Content-Security-Policy", "default-src 'self'")
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

## Vulnerability Prevention

### SQL Injection

**GORM Parameterized Queries:**
```go
// ✅ Safe: Parameterized
db.Where("email = ?", email).First(&user)

// ❌ Unsafe: String concatenation
db.Where("email = '" + email + "'").First(&user)
```

### XSS Prevention

**Template Escaping:**
```go
// Auto-escaped in templates
{{ .UserInput }}

// Manual escaping
{{ .UserInput | html }}
```

### CSRF Protection

**Token-Based:**
```go
// Generate CSRF token
csrfToken := generateCSRFToken()

// Validate CSRF token
if !validateCSRFToken(token) {
    return errors.New("invalid CSRF token")
}
```

## Incident Response

### Security Incident Detection

**Brute Force Detection:**
```go
func detectBruteForce(email string, attempts int) bool {
    if attempts > 5 {
        logger.Warn("Possible brute force attack",
            "email", email,
            "attempts", attempts)
        return true
    }
    return false
}
```

**Anomaly Detection:**
```go
func detectAnomalies(emails []*models.IncomingEmail) []SecurityAlert {
    var alerts []SecurityAlert
    
    // Volume spike
    if len(emails) > normalThreshold {
        alerts = append(alerts, SecurityAlert{
            Type: "volume_spike",
            Details: fmt.Sprintf("Unusual volume: %d emails", len(emails)),
        })
    }
    
    // Sender volume
    senderCounts := make(map[string]int)
    for _, email := range emails {
        senderCounts[email.From]++
    }
    
    for sender, count := range senderCounts {
        if count > senderThreshold {
            alerts = append(alerts, SecurityAlert{
                Type: "sender_volume",
                Details: fmt.Sprintf("High volume from: %s (%d)", sender, count),
            })
        }
    }
    
    return alerts
}
```

### Incident Response Procedure

**1. Detection:**
- Monitor logs voor verdachte activiteit
- Check Prometheus alerts
- Review rate limit violations

**2. Containment:**
```go
// Block IP
func blockIP(ip string) error {
    return firewall.BlockIP(ip, 24*time.Hour)
}

// Revoke tokens
func revokeUserTokens(userID string) error {
    return authService.RevokeAllUserRefreshTokens(ctx, userID)
}
```

**3. Recovery:**
- Reset affected user passwords
- Clear rate limit counters
- Review and update security rules

**4. Post-Incident:**
- Document incident
- Update security measures
- Train team on new threats

## Security Checklist

### Pre-Deployment

- [ ] JWT_SECRET is sterk en uniek
- [ ] Alle wachtwoorden zijn gehashed
- [ ] CORS origins zijn correct
- [ ] Rate limiting is geconfigureerd
- [ ] TLS is enabled voor SMTP/IMAP
- [ ] SSL certificaten zijn geldig
- [ ] Environment variabelen zijn beveiligd
- [ ] Database credentials zijn veilig
- [ ] Logging is geconfigureerd
- [ ] Security headers zijn actief

### Regular Audits

- [ ] Review access logs
- [ ] Check failed login attempts
- [ ] Monitor rate limit violations
- [ ] Verify SSL certificate expiry
- [ ] Update dependencies
- [ ] Review user permissions
- [ ] Check for security patches
- [ ] Audit database access
- [ ] Review API usage patterns
- [ ] Test backup/restore procedures

## Security Best Practices

### Password Policy

**Requirements:**
- Minimum 8 karakters
- Mix van letters en cijfers
- Speciale karakters aanbevolen
- Geen common passwords
- Geen persoonlijke informatie

**Implementation:**
```go
func validatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("wachtwoord moet minimaal 8 karakters zijn")
    }
    
    hasLetter := false
    hasNumber := false
    
    for _, char := range password {
        if unicode.IsLetter(char) {
            hasLetter = true
        }
        if unicode.IsNumber(char) {
            hasNumber = true
        }
    }
    
    if !hasLetter || !hasNumber {
        return errors.New("wachtwoord moet letters en cijfers bevatten")
    }
    
    return nil
}
```

### Token Storage

**Frontend Best Practices:**

**✅ Aanbevolen:**
```javascript
// Access token in memory (React state)
const [accessToken, setAccessToken] = useState(null);

// Refresh token in HttpOnly cookie (server-side)
// Of encrypted localStorage
```

**❌ Niet Aanbevolen:**
```javascript
// Tokens in plain localStorage
localStorage.setItem('token', token);

// Tokens in URL
window.location.href = `/dashboard?token=${token}`;

// Tokens in sessionStorage
sessionStorage.setItem('token', token);
```

### API Key Security

**WFC API Key:** [`handlers/wfc_order_handler.go`](../../handlers/wfc_order_handler.go:1)

```go
func validateAPIKey(c *fiber.Ctx) error {
    apiKey := c.Get("X-API-Key")
    expectedKey := os.Getenv("WFC_API_KEY")
    
    if apiKey == "" || apiKey != expectedKey {
        logger.Warn("Invalid API key attempt",
            "ip", c.IP(),
            "path", c.Path())
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid API key",
        })
    }
    
    return c.Next()
}
```

## Monitoring & Alerts

### Security Metrics

**Prometheus Metrics:**
```
# Failed login attempts
auth_login_failed_total{reason="invalid_credentials"} 5

# Rate limit violations
rate_limit_exceeded_total{operation="login",type="per_ip"} 3

# Permission denials
permission_denied_total{resource="users",action="manage"} 2
```

### Alert Rules

**Prometheus Alerts:**
```yaml
groups:
  - name: security_alerts
    rules:
      - alert: HighFailedLoginRate
        expr: rate(auth_login_failed_total[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High failed login rate detected"
          
      - alert: RateLimitViolations
        expr: rate(rate_limit_exceeded_total[5m]) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Multiple rate limit violations"
```

## Compliance

### GDPR

**Data Protection:**
- User consent voor email opslag
- Privacy policy acceptance
- Data retention policies
- Right to deletion
- Data export capabilities

**Implementation:**
```go
// Privacy akkoord check
if !request.PrivacyAkkoord {
    return errors.New("privacy policy must be accepted")
}

// Data deletion
func (r *Repository) DeleteUserData(ctx context.Context, userID string) error {
    // Delete all user data
    db.Where("user_id = ?", userID).Delete(&models.ContactFormulier{})
    db.Where("user_id = ?", userID).Delete(&models.Aanmelding{})
    db.Where("id = ?", userID).Delete(&models.Gebruiker{})
    return nil
}
```

## Zie Ook

- [Authentication Architecture](../architecture/authentication-and-authorization.md)
- [Deployment Guide](./deployment.md)
- [Monitoring Guide](./monitoring.md)
- [API Documentation](../api/rest-api.md)