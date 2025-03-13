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
```prometheus
# Counter voor verzonden emails
email_sent_total{type="contact|aanmelding",template="admin|user"}

# Counter voor gefaalde emails
email_failed_total{type="contact|aanmelding",error="smtp|template|validation"}

# Histogram voor email verzend latency
email_latency_seconds{type="contact|aanmelding"}
histogram_quantile(0.95, sum(rate(email_latency_bucket[5m])) by (le))

# Gauge voor huidige batch grootte
email_batch_size{type="contact|aanmelding"}

# Histogram voor template render tijd
email_template_render_duration_seconds
```

#### Rate Limiting Metrics
```prometheus
# Counter voor rate limit overschrijdingen
rate_limit_exceeded_total{type="ip|global"}

# Gauge voor resterende requests
rate_limit_remaining{type="ip|global"}

# Gauge voor tijd tot reset
rate_limit_reset_seconds
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
