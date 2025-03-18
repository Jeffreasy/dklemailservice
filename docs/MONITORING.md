# Monitoring Handleiding

## Overzicht

Deze handleiding beschrijft de monitoring setup voor de DKL Email Service, inclusief:
- Prometheus metrics
- ELK logging
- Health checks
- Alerting
- Dashboards

## Prometheus Metrics

### Metrics Endpoints
- `/metrics` - Standaard Prometheus metrics endpoint
- `/api/metrics/email` - Email specifieke metrics
- `/api/metrics/rate-limits` - Rate limiting metrics

### Core Metrics

#### Email Metrics

Monitor verzend statistieken en email gerelateerde informatie:

```http
GET /api/metrics/email
```

**Headers:**
```
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

#### Rate Limiting Metrics

Rate limiting statistieken:

```http
GET /api/metrics/rate-limits
```

**Headers:**
```
X-API-Key: your-admin-api-key
```

**Response:**
```json
{
    "rate_limits": {
        "contact_email": {
            "global_count": 45
        },
        "aanmelding_email": {
            "global_count": 120
        }
    },
    "generated_at": "2024-03-20T15:04:05Z"
}
```

#### System Metrics
```prometheus
# Go runtime metrics
go_goroutines
go_gc_duration_seconds
go_memory_alloc_bytes
go_memory_heap_alloc_bytes

# Process metrics
process_cpu_seconds_total
process_open_fds
process_max_fds
```

### Prometheus Configuratie

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'dklemailservice'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scheme: 'https'
    tls_config:
      insecure_skip_verify: false
    basic_auth:
      username: 'prometheus'
      password: 'secret'
```

## Grafana Dashboards

### Email Service Dashboard

#### Dashboard Panels

1. Email Verzending
```grafana
Panel: Email Verzending Rate
Query: rate(email_sent_total[5m])
Type: Graph
Legend: {{type}} - {{template}}
```

2. Error Rate
```grafana
Panel: Email Fouten
Query: rate(email_failed_total[5m])
Type: Graph
Legend: {{type}} - {{error}}
```

3. Latency
```grafana
Panel: Verzend Latency
Query: histogram_quantile(0.95, sum(rate(email_latency_bucket[5m])) by (le))
Type: Graph
Legend: 95th percentile
```

4. Rate Limiting
```grafana
Panel: Rate Limit Status
Query: rate_limit_remaining
Type: Gauge
Thresholds: 
  - 0-20: red
  - 21-50: yellow
  - 51-100: green
```

### System Dashboard

#### Dashboard Panels

1. Resource Gebruik
```grafana
Panel: Memory Usage
Query: go_memory_alloc_bytes
Type: Graph

Panel: CPU Usage
Query: rate(process_cpu_seconds_total[5m])
Type: Graph

Panel: Goroutines
Query: go_goroutines
Type: Graph
```

2. GC Metrics
```grafana
Panel: GC Duration
Query: rate(go_gc_duration_seconds_sum[5m])
Type: Graph
```

## ELK Stack Logging

### Logstash Configuratie

```ruby
# logstash.conf
input {
  tcp {
    port => 5000
    codec => json
  }
}

filter {
  if [type] == "email_service" {
    date {
      match => [ "timestamp", "ISO8601" ]
    }
    mutate {
      add_field => {
        "service" => "dklemailservice"
      }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "dkl-emails-%{+YYYY.MM.dd}"
  }
}
```

### Elasticsearch Index Template

```json
{
  "index_patterns": ["dkl-emails-*"],
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "@timestamp": { "type": "date" },
      "level": { "type": "keyword" },
      "message": { "type": "text" },
      "email_type": { "type": "keyword" },
      "template": { "type": "keyword" },
      "duration_ms": { "type": "long" },
      "success": { "type": "boolean" }
    }
  }
}
```

### Kibana Visualisaties

1. Email Verzending Dashboard
   - Line chart: Emails per uur per type
   - Pie chart: Verdeling van email types
   - Data table: Laatste fouten
   - Metrics: Success rate

2. Performance Dashboard
   - Line chart: Gemiddelde verzendtijd
   - Heat map: Verzendtijd distributie
   - Bar chart: Fouten per type

## Health Checks

### Endpoints

1. `/api/health` - Basis health check
```json
{
  "status": "up",
  "timestamp": "2024-03-20T15:04:05Z",
  "version": "1.2.0",
  "checks": {
    "smtp": {
      "status": "up",
      "latency_ms": 150
    },
    "templates": {
      "status": "up",
      "loaded": 4
    },
    "rate_limiter": {
      "status": "up",
      "current_rate": 45
    }
  }
}
```

2. `/api/health/smtp` - SMTP specifieke check
```json
{
  "status": "up",
  "last_check": "2024-03-20T15:04:05Z",
  "latency_ms": 150,
  "tls_enabled": true,
  "connection_pool": {
    "active": 5,
    "idle": 3,
    "max": 10
  }
}
```

### Monitoring Configuratie

```yaml
# health-check.yml
endpoints:
  - url: https://api.dekoninklijkeloop.nl/api/health
    method: GET
    interval: 30s
    timeout: 5s
    conditions:
      - type: status
        status: 200
      - type: json
        expression: "$.status"
        equals: "up"
    alerts:
      - type: email
        to: admin@dekoninklijkeloop.nl
      - type: slack
        webhook: https://hooks.slack.com/...
```

## Alerting

### Alert Regels

```yaml
# alerting.yml
groups:
  - name: email_service_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(email_failed_total[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Hoge error rate in email verzending"
          description: "Error rate is {{ $value }}% over 5 minuten"

      - alert: SMTPLatencyHigh
        expr: histogram_quantile(0.95, sum(rate(email_latency_bucket[5m])) by (le)) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "SMTP latency is hoog"
          description: "95th percentile latency is {{ $value }}s"

      - alert: RateLimitNearlyExhausted
        expr: rate_limit_remaining < 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Rate limit bijna bereikt"
          description: "Nog {{ $value }} requests beschikbaar"
```

### Alert Notificaties

1. Email Notificaties
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

2. Slack Notificaties
```yaml
receivers:
  - name: slack-alerts
    slack_configs:
      - api_url: https://hooks.slack.com/services/...
        channel: '#monitoring'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ .CommonAnnotations.description }}'
```

## Backup & Recovery

### Metrics Data

1. Prometheus Data Backup
```bash
# Snapshot van Prometheus data
curl -XPOST http://localhost:9090/-/snapshot
```

2. Elasticsearch Backup
```bash
# Create repository
PUT /_snapshot/backup_repo
{
  "type": "fs",
  "settings": {
    "location": "/backups/elasticsearch"
  }
}

# Create snapshot
PUT /_snapshot/backup_repo/snapshot_1
```

### Recovery Procedures

1. Prometheus Recovery
```bash
# Restore van snapshot
prometheus --storage.tsdb.path=/path/to/snapshot
```

2. Elasticsearch Recovery
```bash
# Restore van snapshot
POST /_snapshot/backup_repo/snapshot_1/_restore
```

### Email Auto Fetcher Monitoring

Voor het monitoren van de automatische email ophaal functionaliteit zijn de volgende mogelijkheden beschikbaar:

#### Status Monitoring

De huidige status van de EmailAutoFetcher is beschikbaar via de health endpoint, die aangeeft of de service actief is:

```http
GET /api/health
```

**Response:**
```json
{
    "status": "healthy",
    "version": "1.2.3",
    "services": {
        "database": "connected",
        "email_service": "operational",
        "email_auto_fetcher": "running", // of "stopped" als niet actief
        "rate_limiter": "active"
    },
    "last_email_fetch": "2024-04-02T14:25:00Z"
}
```

#### Email Fetch Metrics

Gedetailleerde metrics over de email fetch operaties:

```http
GET /api/metrics/email-fetch
```

**Headers:**
```
X-API-Key: your-admin-api-key
```

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
    },
    "generated_at": "2024-04-02T14:30:00Z"
}
```

#### Log Monitoring

De EmailAutoFetcher logt alle activiteiten en fouten. Relevante log entries kunnen gefilterd worden:

```bash
# Filter logs op EmailAutoFetcher activiteit
grep "EmailAutoFetcher" server.log

# Filter logs op fetch operaties
grep "Fetching emails" server.log

# Filter logs op fouten
grep "Error fetching emails" server.log
```

Voorbeeld log output:
```
2024-04-02T14:25:00Z INFO EmailAutoFetcher: Starting email fetch operation
2024-04-02T14:25:01Z INFO EmailAutoFetcher: Fetching emails from info@dekoninklijkeloop.nl
2024-04-02T14:25:02Z INFO EmailAutoFetcher: Found 5 emails, 3 new emails saved
2024-04-02T14:25:03Z INFO EmailAutoFetcher: Fetching emails from inschrijving@dekoninklijkeloop.nl
2024-04-02T14:25:04Z INFO EmailAutoFetcher: Found 2 emails, 2 new emails saved
2024-04-02T14:25:05Z INFO EmailAutoFetcher: Email fetch operation completed successfully
```

#### Prometheus Metrics

De volgende Prometheus metrics zijn beschikbaar op `/metrics`:

```
# HELP dkl_email_fetcher_runs_total Totaal aantal email fetch operaties
# TYPE dkl_email_fetcher_runs_total counter
dkl_email_fetcher_runs_total 96

# HELP dkl_email_fetcher_errors_total Aantal mislukte email fetch operaties
# TYPE dkl_email_fetcher_errors_total counter
dkl_email_fetcher_errors_total 1

# HELP dkl_email_fetcher_duration_seconds Tijd besteed aan email fetch operaties
# TYPE dkl_email_fetcher_duration_seconds histogram
dkl_email_fetcher_duration_seconds_bucket{le="0.1"} 3
dkl_email_fetcher_duration_seconds_bucket{le="0.5"} 56
dkl_email_fetcher_duration_seconds_bucket{le="1"} 80
dkl_email_fetcher_duration_seconds_bucket{le="2"} 92
dkl_email_fetcher_duration_seconds_bucket{le="5"} 96
dkl_email_fetcher_duration_seconds_bucket{le="+Inf"} 96
dkl_email_fetcher_duration_seconds_sum 104.2
dkl_email_fetcher_duration_seconds_count 96

# HELP dkl_emails_fetched_total Totaal aantal opgehaalde emails
# TYPE dkl_emails_fetched_total counter
dkl_emails_fetched_total{account="info"} 100
dkl_emails_fetched_total{account="inschrijving"} 50

# HELP dkl_emails_saved_total Totaal aantal opgeslagen emails
# TYPE dkl_emails_saved_total counter
dkl_emails_saved_total{account="info"} 85
dkl_emails_saved_total{account="inschrijving"} 35

# HELP dkl_emails_duplicates_total Totaal aantal gedetecteerde duplicate emails
# TYPE dkl_emails_duplicates_total counter
dkl_emails_duplicates_total{account="info"} 15
dkl_emails_duplicates_total{account="inschrijving"} 15

# HELP dkl_email_fetcher_last_run_timestamp_seconds Timestamp van laatste email fetch operatie
# TYPE dkl_email_fetcher_last_run_timestamp_seconds gauge
dkl_email_fetcher_last_run_timestamp_seconds 1712072700
```

## Logging en Alerts
