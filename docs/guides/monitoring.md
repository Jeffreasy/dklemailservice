# Monitoring Guide

Complete monitoring handleiding voor de DKL Email Service.

## Overzicht

De service biedt uitgebreide monitoring via:
- Prometheus metrics
- Health checks
- Structured logging
- ELK Stack integratie
- Real-time alerts

## Health Checks

### Health Endpoint

**Endpoint:** `GET /api/health`

**Implementatie:** [`handlers/health_handler.go`](../../handlers/health_handler.go:1)

**Response (Healthy):**
```json
{
    "status": "healthy",
    "version": "1.0.0",
    "timestamp": "2024-03-20T15:04:05Z",
    "uptime": "24h3m12s",
    "services": {
        "database": "connected",
        "redis": "connected",
        "email_service": "operational",
        "email_auto_fetcher": "running",
        "rate_limiter": "active"
    }
}
```

**Response (Unhealthy):**
```json
{
    "status": "unhealthy",
    "version": "1.0.0",
    "timestamp": "2024-03-20T15:04:05Z",
    "services": {
        "database": "disconnected",
        "redis": "connected",
        "email_service": "operational"
    },
    "error": "Database connection failed"
}
```

### Monitoring Setup

**Uptime Monitoring:**
```bash
# Curl check
curl -f http://localhost:8080/api/health || exit 1

# Docker health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1
```

**Kubernetes Liveness Probe:**
```yaml
livenessProbe:
  httpGet:
    path: /api/health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3
```

## Prometheus Metrics

### Metrics Endpoint

**Endpoint:** `GET /metrics`

**Implementatie:** [`main.go:525`](../../main.go:525)

```go
app.Get("/metrics", func(c *fiber.Ctx) error {
    registry := prometheus.DefaultRegisterer.(*prometheus.Registry)
    handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
    
    recorder := httptest.NewRecorder()
    request := httptest.NewRequest("GET", "/metrics", nil)
    
    handler.ServeHTTP(recorder, request)
    
    for k, v := range recorder.Header() {
        for _, val := range v {
            c.Set(k, val)
        }
    }
    
    return c.Status(recorder.Code).Send(recorder.Body.Bytes())
})
```

### Available Metrics

**Email Metrics:**
```prometheus
# Total emails sent
email_sent_total{type="contact",status="success"} 150
email_sent_total{type="aanmelding",status="success"} 200

# Failed emails
email_failed_total{type="contact",reason="smtp_error"} 2
email_failed_total{type="aanmelding",reason="rate_limited"} 1

# Email latency
email_latency_seconds_bucket{type="contact",le="0.1"} 50
email_latency_seconds_bucket{type="contact",le="0.5"} 140
email_latency_seconds_bucket{type="contact",le="1"} 150
email_latency_seconds_sum{type="contact"} 45.2
email_latency_seconds_count{type="contact"} 150
```

**Rate Limit Metrics:**
```prometheus
# Rate limit exceeded
rate_limit_exceeded_total{operation="login",type="per_ip"} 5
rate_limit_exceeded_total{operation="email_generic",type="global"} 2
```

**System Metrics:**
```prometheus
# Go runtime
go_goroutines 42
go_gc_duration_seconds 0.002
go_memstats_alloc_bytes 8388608

# Process
process_cpu_seconds_total 120.5
process_open_fds 15
process_max_fds 1024
```

**Email Auto Fetcher Metrics:**
```prometheus
# Fetch operations
email_fetcher_runs_total 96
email_fetcher_errors_total 1

# Fetch duration
email_fetcher_duration_seconds_bucket{le="1"} 80
email_fetcher_duration_seconds_bucket{le="2"} 92
email_fetcher_duration_seconds_sum 104.2
email_fetcher_duration_seconds_count 96

# Emails processed
emails_fetched_total{account="info"} 100
emails_fetched_total{account="inschrijving"} 50
emails_saved_total{account="info"} 85
emails_duplicates_total{account="info"} 15
```

### Prometheus Configuration

**prometheus.yml:**
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'dklemailservice'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scheme: 'http'
```

**Production:**
```yaml
scrape_configs:
  - job_name: 'dklemailservice'
    static_configs:
      - targets: ['api.dekoninklijkeloop.nl:8080']
    metrics_path: '/metrics'
    scheme: 'https'
    tls_config:
      insecure_skip_verify: false
```

## Custom Metrics Endpoints

### Email Metrics

**Endpoint:** `GET /api/metrics/email`

**Headers:**
```http
X-API-Key: your-admin-api-key
```

**Response:**
```json
{
    "total_emails": 150,
    "success_rate": 98.5,
    "emails_by_type": {
        "contact": {
            "sent": 50,
            "failed": 1
        },
        "aanmelding": {
            "sent": 100,
            "failed": 2
        }
    },
    "generated_at": "2024-03-20T15:04:05Z"
}
```

**Implementatie:** [`handlers/metrics_handler.go`](../../handlers/metrics_handler.go:1)

### Rate Limit Metrics

**Endpoint:** `GET /api/metrics/rate-limits`

**Response:**
```json
{
    "rate_limits": {
        "contact_email": {
            "global_count": 45,
            "limit": 100,
            "period": "1h"
        },
        "aanmelding_email": {
            "global_count": 120,
            "limit": 200,
            "period": "1h"
        }
    },
    "generated_at": "2024-03-20T15:04:05Z"
}
```

## Grafana Dashboards

### Email Service Dashboard

**Panels:**

**1. Email Sending Rate:**
```promql
rate(email_sent_total[5m])
```

**2. Error Rate:**
```promql
rate(email_failed_total[5m])
```

**3. Success Rate:**
```promql
sum(rate(email_sent_total{status="success"}[5m])) / 
sum(rate(email_sent_total[5m])) * 100
```

**4. Latency (95th percentile):**
```promql
histogram_quantile(0.95, 
  sum(rate(email_latency_seconds_bucket[5m])) by (le, type)
)
```

**5. Rate Limit Status:**
```promql
rate_limit_exceeded_total
```

### System Dashboard

**Panels:**

**1. Memory Usage:**
```promql
go_memstats_alloc_bytes
```

**2. CPU Usage:**
```promql
rate(process_cpu_seconds_total[5m])
```

**3. Goroutines:**
```promql
go_goroutines
```

**4. GC Duration:**
```promql
rate(go_gc_duration_seconds_sum[5m])
```

### Email Auto Fetcher Dashboard

**Panels:**

**1. Fetch Success Rate:**
```promql
sum(rate(email_fetcher_runs_total[5m])) - 
sum(rate(email_fetcher_errors_total[5m]))
```

**2. Emails Fetched:**
```promql
rate(emails_fetched_total[5m])
```

**3. Duplicate Rate:**
```promql
rate(emails_duplicates_total[5m]) / 
rate(emails_fetched_total[5m]) * 100
```

**4. Fetch Duration:**
```promql
histogram_quantile(0.95, 
  sum(rate(email_fetcher_duration_seconds_bucket[5m])) by (le)
)
```

## Logging

### Log Levels

**Configuratie:**
```bash
LOG_LEVEL=info  # debug, info, warn, error, fatal
```

**Implementatie:** [`main.go:116`](../../main.go:116)

```go
logLevel := os.Getenv("LOG_LEVEL")
if logLevel == "" {
    logLevel = logger.InfoLevel
}
logger.Setup(logLevel)
```

### Structured Logging

**Voorbeeld:**
```go
logger.Info("email sent",
    "type", emailType,
    "recipient", recipient,
    "duration", duration,
    "template", templateName,
)

logger.Error("failed to send email",
    "type", emailType,
    "error", err,
    "retry_count", retryCount,
    "recipient", recipient,
)
```

### Log Output

**Console Output:**
```json
{
    "level": "info",
    "timestamp": "2024-03-20T15:04:05Z",
    "message": "email sent",
    "type": "contact",
    "recipient": "user@example.com",
    "duration": "0.5s"
}
```

## ELK Stack Integration

### Configuration

**Environment:**
```bash
ELK_ENDPOINT=http://elk:9200
ELK_INDEX=dkl-emails
ELK_BATCH_SIZE=100
```

**Implementatie:** [`main.go:152`](../../main.go:152)

```go
if elkEndpoint := os.Getenv("ELK_ENDPOINT"); elkEndpoint != "" {
    logger.SetupELK(logger.ELKConfig{
        Endpoint:      elkEndpoint,
        BatchSize:     100,
        FlushInterval: 5 * time.Second,
        AppName:       "dklemailservice",
        Environment:   os.Getenv("ENVIRONMENT"),
    })
    logger.Info("ELK logging enabled", "endpoint", elkEndpoint)
}
```

### Elasticsearch Index

**Index Pattern:**
```
dkl-emails-2024.03.20
dkl-emails-2024.03.21
```

**Mapping:**
```json
{
    "mappings": {
        "properties": {
            "@timestamp": { "type": "date" },
            "level": { "type": "keyword" },
            "message": { "type": "text" },
            "email_type": { "type": "keyword" },
            "template": { "type": "keyword" },
            "duration_ms": { "type": "long" },
            "success": { "type": "boolean" },
            "error": { "type": "text" }
        }
    }
}
```

### Kibana Queries

**Failed Emails:**
```
level: "error" AND message: "failed to send email"
```

**Rate Limit Violations:**
```
message: "rate limit" AND level: "warn"
```

**Slow Emails:**
```
duration_ms: >1000
```

## Alerting

### Alert Rules

**Prometheus Alerts:**
```yaml
groups:
  - name: email_service_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(email_failed_total[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High email error rate"
          description: "Error rate is {{ $value }}% over 5 minutes"
      
      - alert: EmailLatencyHigh
        expr: histogram_quantile(0.95, sum(rate(email_latency_seconds_bucket[5m])) by (le)) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Email latency is high"
          description: "95th percentile latency is {{ $value }}s"
      
      - alert: RateLimitNearlyExhausted
        expr: rate_limit_exceeded_total > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Rate limit frequently exceeded"
      
      - alert: EmailFetcherDown
        expr: up{job="dklemailservice"} == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Email Auto Fetcher is down"
```

### Alert Channels

**Email Notifications:**
```yaml
receivers:
  - name: email-alerts
    email_configs:
      - to: admin@dekoninklijkeloop.nl
        from: alerts@dekoninklijkeloop.nl
        smarthost: smtp.example.com:587
        auth_username: alerts@dekoninklijkeloop.nl
        auth_password: secret
```

**Slack Notifications:**
```yaml
receivers:
  - name: slack-alerts
    slack_configs:
      - api_url: https://hooks.slack.com/services/...
        channel: '#monitoring'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ .CommonAnnotations.description }}'
```

**Telegram Notifications:**

De service heeft ingebouwde Telegram notificaties.

**Configuratie:**
```bash
ENABLE_NOTIFICATIONS=true
TELEGRAM_BOT_TOKEN=your-bot-token
TELEGRAM_CHAT_ID=your-chat-id
NOTIFICATION_MIN_PRIORITY=medium
```

**Implementatie:** [`services/notification_service.go`](../../services/notification_service.go:1)

## Application Metrics

### Email Metrics

**Implementatie:** [`services/email_metrics.go`](../../services/email_metrics.go:1)

```go
type EmailMetrics struct {
    mu           sync.RWMutex
    sentCount    map[string]int
    failedCount  map[string]int
    lastReset    time.Time
    resetPeriod  time.Duration
}

func (m *EmailMetrics) RecordEmailSent(emailType string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.sentCount[emailType]++
}

func (m *EmailMetrics) RecordEmailFailed(emailType string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.failedCount[emailType]++
}
```

### Prometheus Metrics

**Implementatie:** [`services/prometheus_metrics.go`](../../services/prometheus_metrics.go:1)

```go
var (
    emailsSent = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "email_sent_total",
            Help: "Total number of emails sent",
        },
        []string{"type", "status"},
    )
    
    emailsFailed = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "email_failed_total",
            Help: "Total number of failed emails",
        },
        []string{"type", "reason"},
    )
    
    emailLatency = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "email_latency_seconds",
            Help:    "Email sending latency in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"type"},
    )
)
```

**Usage:**
```go
func (s *EmailService) SendEmail(to, subject, body string) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        s.prometheusMetrics.ObserveEmailLatency("email_generic", duration.Seconds())
    }()
    
    err := s.smtpClient.Send(msg)
    
    if err != nil {
        s.prometheusMetrics.RecordEmailFailed("email_generic", "smtp_error")
        return err
    }
    
    s.prometheusMetrics.RecordEmailSent("email_generic", "success")
    return nil
}
```

## Log Monitoring

### Log Aggregation

**ELK Writer:** [`logger/elk_writer.go`](../../logger/elk_writer.go:1)

```go
type ELKWriter struct {
    endpoint      string
    index         string
    batchSize     int
    flushInterval time.Duration
    buffer        [][]byte
    mutex         sync.Mutex
}

func (w *ELKWriter) Write(p []byte) (n int, err error) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    w.buffer = append(w.buffer, p)
    
    if len(w.buffer) >= w.batchSize {
        return len(p), w.flush()
    }
    
    return len(p), nil
}
```

**Batch Configuration:**
```bash
ELK_BATCH_SIZE=100
ELK_FLUSH_INTERVAL=5s
```

### Log Queries

**Kibana Queries:**

**Failed Emails:**
```
level: "error" AND message: "failed to send email"
```

**Authentication Failures:**
```
level: "warn" AND message: "Ongeldige inloggegevens"
```

**Rate Limit Violations:**
```
message: "rate limit overschreden"
```

**Slow Operations:**
```
duration_ms: >1000
```

## Performance Monitoring

### Response Time

**Track Latency:**
```go
func trackLatency(operation string) func() {
    start := time.Now()
    return func() {
        duration := time.Since(start)
        logger.Info("operation completed",
            "operation", operation,
            "duration", duration,
        )
        prometheusMetrics.ObserveLatency(operation, duration.Seconds())
    }
}

// Usage
defer trackLatency("send_email")()
```

### Database Performance

**Slow Query Logging:**
```go
// GORM logger configuratie
gormLog := gormlogger.New(
    &dbLogger{},
    gormlogger.Config{
        SlowThreshold: time.Second,  // Log queries > 1s
        LogLevel:      gormlogger.Warn,
    },
)
```

**Query Metrics:**
```promql
# Slow queries
db_query_duration_seconds{quantile="0.95"} > 1

# Query count
rate(db_queries_total[5m])
```

## Email Auto Fetcher Monitoring

### Status Monitoring

**Health Check Response:**
```json
{
    "services": {
        "email_auto_fetcher": "running",
        "last_email_fetch": "2024-04-02T14:25:00Z"
    }
}
```

### Fetch Metrics

**Endpoint:** `GET /api/metrics/email-fetch`

**Response:**
```json
{
    "auto_fetcher_status": "running",
    "fetch_interval_minutes": 15,
    "last_run": "2024-04-02T14:25:00Z",
    "next_scheduled_run": "2024-04-02T14:40:00Z",
    "total_fetches": 96,
    "successful_fetches": 95,
    "failed_fetches": 1,
    "fetch_stats": {
        "total_emails_found": 150,
        "total_emails_saved": 120,
        "duplicate_emails": 30
    },
    "accounts": {
        "info@dekoninklijkeloop.nl": {
            "emails_found": 100,
            "emails_saved": 85,
            "last_fetch_status": "success"
        },
        "inschrijving@dekoninklijkeloop.nl": {
            "emails_found": 50,
            "emails_saved": 35,
            "last_fetch_status": "success"
        }
    }
}
```

### Log Monitoring

**Filter Logs:**
```bash
# EmailAutoFetcher activiteit
grep "EmailAutoFetcher" server.log

# Fetch operaties
grep "Fetching emails" server.log

# Fouten
grep "Error fetching emails" server.log
```

**Log Output:**
```
2024-04-02T14:25:00Z INFO EmailAutoFetcher: Starting email fetch operation
2024-04-02T14:25:01Z INFO EmailAutoFetcher: Fetching emails from info@dekoninklijkeloop.nl
2024-04-02T14:25:02Z INFO EmailAutoFetcher: Found 5 emails, 3 new emails saved
2024-04-02T14:25:03Z INFO EmailAutoFetcher: Email fetch operation completed successfully
```

## Alerting Configuration

### Alert Manager

**alertmanager.yml:**
```yaml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'email-alerts'
  
  routes:
    - match:
        severity: critical
      receiver: 'critical-alerts'
    
    - match:
        severity: warning
      receiver: 'warning-alerts'

receivers:
  - name: 'email-alerts'
    email_configs:
      - to: 'admin@dekoninklijkeloop.nl'
        from: 'alerts@dekoninklijkeloop.nl'
        smarthost: 'smtp.example.com:587'
        auth_username: 'alerts@dekoninklijkeloop.nl'
        auth_password: 'secret'
  
  - name: 'critical-alerts'
    email_configs:
      - to: 'admin@dekoninklijkeloop.nl'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/...'
        channel: '#critical-alerts'
  
  - name: 'warning-alerts'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/...'
        channel: '#warnings'
```

## Monitoring Best Practices

### Metrics Collection

**What to Monitor:**
- Request rate
- Error rate
- Response time (p50, p95, p99)
- Resource usage (CPU, memory)
- Database connections
- Email queue size
- Rate limit usage

**What NOT to Monitor:**
- Sensitive user data
- Password hashes
- API keys
- Personal information

### Log Retention

**Retention Policy:**
- **Debug logs:** 7 dagen
- **Info logs:** 30 dagen
- **Warn logs:** 90 dagen
- **Error logs:** 1 jaar
- **Audit logs:** 2 jaar

**Elasticsearch ILM:**
```json
{
    "policy": {
        "phases": {
            "hot": {
                "actions": {
                    "rollover": {
                        "max_age": "1d",
                        "max_size": "50gb"
                    }
                }
            },
            "delete": {
                "min_age": "30d",
                "actions": {
                    "delete": {}
                }
            }
        }
    }
}
```

### Dashboard Organization

**Recommended Dashboards:**
1. **Overview** - High-level metrics
2. **Email Service** - Email-specific metrics
3. **System** - Resource usage
4. **Security** - Auth and rate limiting
5. **Email Fetcher** - IMAP operations

## Troubleshooting

### High Error Rate

**Check:**
```bash
# Recent errors
curl http://localhost:8080/api/metrics/email | jq '.emails_by_type'

# Prometheus query
rate(email_failed_total[5m])
```

**Common Causes:**
- SMTP server down
- Invalid credentials
- Rate limiting
- Network issues

### High Latency

**Check:**
```promql
histogram_quantile(0.95, sum(rate(email_latency_seconds_bucket[5m])) by (le))
```

**Common Causes:**
- Slow SMTP server
- Network latency
- Database slow queries
- High concurrent load

### Memory Leaks

**Check:**
```promql
go_memstats_alloc_bytes
go_goroutines
```

**Debug:**
```bash
# Memory profile
go test -memprofile=mem.prof
go tool pprof mem.prof

# Check goroutines
curl http://localhost:8080/debug/pprof/goroutine
```

## Zie Ook

- [Deployment Guide](./deployment.md) - Production deployment
- [Security Guide](./security.md) - Security best practices
- [Development Guide](./development.md) - Development setup
- [API Documentation](../api/rest-api.md) - API reference