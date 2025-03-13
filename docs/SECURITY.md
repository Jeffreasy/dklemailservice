# Security Handleiding

## Overzicht

Deze handleiding beschrijft de security maatregelen en best practices voor de DKL Email Service, inclusief:
- Input validatie
- Rate limiting
- SMTP security
- API security
- Data bescherming
- Monitoring & logging
- Incident response

## Input Validatie

### Email Validatie
```go
// EmailValidator valideert email adressen
type EmailValidator struct {
    maxLength int
    pattern   *regexp.Regexp
}

func NewEmailValidator() *EmailValidator {
    return &EmailValidator{
        maxLength: 254, // RFC 5321
        pattern: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
    }
}

func (v *EmailValidator) Validate(email string) error {
    if len(email) > v.maxLength {
        return fmt.Errorf("email te lang (max %d karakters)", v.maxLength)
    }
    if !v.pattern.MatchString(email) {
        return errors.New("ongeldig email formaat")
    }
    return nil
}
```

### Form Validatie
```go
// ContactForm validatie
type ContactForm struct {
    Naam           string `validate:"required,min=2,max=100"`
    Email          string `validate:"required,email"`
    Bericht        string `validate:"required,min=10,max=5000"`
    PrivacyAkkoord bool   `validate:"required,eq=true"`
}

// AanmeldingForm validatie
type AanmeldingForm struct {
    Naam           string `validate:"required,min=2,max=100"`
    Email          string `validate:"required,email"`
    Telefoon       string `validate:"required,e164"`
    Rol            string `validate:"required,oneof=loper vrijwilliger"`
    Afstand        string `validate:"required,oneof=5km 10km 21.1km"`
    Ondersteuning  string `validate:"omitempty,max=500"`
    Bijzonderheden string `validate:"omitempty,max=1000"`
    Terms          bool   `validate:"required,eq=true"`
}
```

## Rate Limiting

### Configuratie
```go
// RateLimiter configuratie
type RateLimiterConfig struct {
    GlobalLimit struct {
        Contact    int           `env:"GLOBAL_RATE_LIMIT_CONTACT" default:"100"`
        Aanmelding int           `env:"GLOBAL_RATE_LIMIT_AANMELDING" default:"200"`
        Window     time.Duration `env:"RATE_LIMIT_WINDOW" default:"1h"`
    }
    IPLimit struct {
        Contact    int           `env:"IP_RATE_LIMIT_CONTACT" default:"5"`
        Aanmelding int           `env:"IP_RATE_LIMIT_AANMELDING" default:"10"`
        Window     time.Duration `env:"RATE_LIMIT_WINDOW" default:"1h"`
    }
}

// Rate limiter middleware
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        endpoint := c.FullPath()
        
        if err := limiter.Allow(endpoint, ip); err != nil {
            c.JSON(429, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": limiter.RetryAfter(endpoint, ip),
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## SMTP Security

### TLS Configuratie
```go
// SMTP TLS configuratie
type SMTPConfig struct {
    Host       string        `env:"SMTP_HOST,required"`
    Port       int           `env:"SMTP_PORT" default:"587"`
    Username   string        `env:"SMTP_USER,required"`
    Password   string        `env:"SMTP_PASSWORD,required"`
    TLSEnabled bool         `env:"SMTP_TLS_ENABLED" default:"true"`
    Timeout    time.Duration `env:"SMTP_TIMEOUT" default:"10s"`
}

func NewSMTPClient(config SMTPConfig) (*smtp.Client, error) {
    tlsConfig := &tls.Config{
        ServerName:         config.Host,
        MinVersion:        tls.VersionTLS12,
        CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }
    
    // SMTP client setup met TLS
    return smtp.NewClient(config, tlsConfig)
}
```

## API Security

### CORS Configuratie
```go
// CORS middleware configuratie
func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

### Security Headers
```go
// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

## Template Security

### XSS Preventie
```go
// Template sanitization
func sanitizeTemplate(tmpl *template.Template) *template.Template {
    // HTML escape by default
    tmpl.Option("missingkey=error")
    tmpl.Funcs(template.FuncMap{
        "safeHTML": func(s string) template.HTML {
            return template.HTML(bluemonday.UGCPolicy().Sanitize(s))
        },
        "safeURL": func(s string) template.URL {
            return template.URL(bluemonday.UGCPolicy().Sanitize(s))
        },
    })
    return tmpl
}
```

## Data Protection

### Sensitive Data Handling
```go
// Sensitive data masking
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

// Secure logging
func secureLog(logger *zap.Logger, email string) {
    logger.Info("email processed",
        zap.String("email", maskEmail(email)),
    )
}
```

## Error Handling

### Secure Error Responses
```go
// ErrorResponse struct
type ErrorResponse struct {
    Error       string `json:"error"`
    Code        string `json:"code,omitempty"`
    RequestID   string `json:"request_id,omitempty"`
}

// Error handler middleware
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *ValidationError:
                c.JSON(400, ErrorResponse{
                    Error: "Validatie fout",
                    Code:  e.Code,
                })
            default:
                // Log internal error details
                logger.Error("internal error",
                    zap.Error(err.Err),
                    zap.String("request_id", c.GetString("request_id")),
                )
                
                // Return safe error to client
                c.JSON(500, ErrorResponse{
                    Error:     "Er is een fout opgetreden",
                    RequestID: c.GetString("request_id"),
                })
            }
        }
    }
}
```

## Security Monitoring

### Audit Logging
```go
// AuditLogger voor security events
type AuditLogger struct {
    logger *zap.Logger
}

func (a *AuditLogger) LogSecurityEvent(event string, details map[string]interface{}) {
    fields := []zap.Field{
        zap.String("event", event),
        zap.Time("timestamp", time.Now()),
    }
    
    for k, v := range details {
        fields = append(fields, zap.Any(k, v))
    }
    
    a.logger.Info("security_audit", fields...)
}
```

### Security Metrics
```go
// Security metrics
var (
    securityEvents = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "security_events_total",
            Help: "Aantal security events per type",
        },
        []string{"event_type"},
    )
    
    rateLimitExceeded = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rate_limit_exceeded_total",
            Help: "Aantal rate limit overschrijdingen per endpoint",
        },
        []string{"endpoint", "ip"},
    )
    
    invalidAuthAttempts = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "invalid_auth_attempts_total",
            Help: "Aantal ongeldige authenticatie pogingen",
        },
    )
)
```

## Incident Response

### Security Incident Procedure

1. Detectie & Analyse
```go
// Incident detectie
type SecurityIncident struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Type      string    `json:"type"`
    Severity  string    `json:"severity"`
    Details   string    `json:"details"`
    IPAddress string    `json:"ip_address"`
}

// Incident handler
func HandleSecurityIncident(incident SecurityIncident) {
    // Log incident
    logger.Error("security incident detected",
        zap.String("id", incident.ID),
        zap.String("type", incident.Type),
        zap.String("severity", incident.Severity),
    )
    
    // Notify administrators
    notifyAdmins(incident)
    
    // Take immediate action if needed
    switch incident.Type {
    case "brute_force":
        blockIP(incident.IPAddress)
    case "spam":
        updateRateLimits(incident.IPAddress)
    }
}
```

2. Containment & Recovery
```go
// IP blocking
func blockIP(ip string) error {
    return firewall.BlockIP(ip, time.Hour*24)
}

// Rate limit adjustment
func updateRateLimits(ip string) {
    rateLimiter.SetLimit(ip, 1, time.Hour)
}
```

3. Post-Incident
```go
// Incident report
type IncidentReport struct {
    Incident   SecurityIncident `json:"incident"`
    Resolution string          `json:"resolution"`
    Timeline   []TimelineEvent `json:"timeline"`
    Actions    []Action        `json:"actions"`
}

// Generate report
func GenerateIncidentReport(incident SecurityIncident) *IncidentReport {
    return &IncidentReport{
        Incident:   incident,
        Resolution: "Issue resolved by blocking malicious IP",
        Timeline:   getIncidentTimeline(incident.ID),
        Actions:    getRemediationActions(incident.Type),
    }
}
```

## Security Checklist

### Pre-Deployment
- [ ] Input validatie geïmplementeerd
- [ ] Rate limiting geconfigureerd
- [ ] CORS correct ingesteld
- [ ] TLS/SSL geconfigureerd
- [ ] Security headers actief
- [ ] Logging & monitoring opgezet
- [ ] Error handling geïmplementeerd
- [ ] Template sanitization actief
- [ ] Audit logging ingeschakeld
- [ ] Incident response plan gereed

### Regular Checks
- [ ] SSL certificaten geldig
- [ ] Dependencies up-to-date
- [ ] Security patches geïnstalleerd
- [ ] Firewall regels correct
- [ ] Rate limits effectief
- [ ] Logging functioneel
- [ ] Backups recent
- [ ] Monitoring actief
- [ ] Alerts functioneel
- [ ] Incident response getest 