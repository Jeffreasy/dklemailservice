# Deployment Guide

Production deployment guide for the DKL Email Service.

## Overview

This guide covers deployment to production environments, focusing on:
- Render.com deployment (recommended)
- Docker deployment
- Security hardening
- Environment configuration
- Monitoring setup

## Pre-Deployment Checklist

### Code Preparation
- [ ] All tests passing
- [ ] Code reviewed and merged to main branch
- [ ] Database migrations tested
- [ ] Environment variables documented
- [ ] Secrets configured securely
- [ ] CORS settings updated for production domains
- [ ] Rate limiting configured appropriately

### Infrastructure
- [ ] PostgreSQL database provisioned
- [ ] Redis instance provisioned
- [ ] Cloudinary account configured
- [ ] Email service configured
- [ ] SSL certificates obtained
- [ ] Domain DNS configured
- [ ] Backup system configured

## Render.com Deployment

### 1. Create Render Account

Sign up at [Render.com](https://render.com) and create a new project.

### 2. Create PostgreSQL Database

```yaml
# In Render Dashboard:
1. Click "New +" → "PostgreSQL"
2. Name: dkl-email-service-db
3. Database: dklemailservice
4. User: dkl_admin
5. Region: Frankfurt (or closest to users)
6. Plan: Choose based on needs (Starter for development)
```

Save the internal connection string:
```
postgresql://user:password@hostname:5432/database
```

### 3. Create Redis Instance

```yaml
# In Render Dashboard:
1. Click "New +" → "Redis"
2. Name: dkl-email-service-redis  
3. Region: Same as PostgreSQL
4. Plan: Choose based on needs
```

Save the internal connection string:
```
redis://hostname:6379
```

### 4. Create Web Service

```yaml
# render.yaml (place in project root)
services:
  - type: web
    name: dkl-email-service
    env: go
    region: frankfurt
    plan: starter
    buildCommand: go build -o main .
    startCommand: ./main
    
    envVars:
      - key: PORT
        value: 8080
      
      - key: ENV
        value: production
      
      - key: DB_HOST
        fromDatabase:
          name: dkl-email-service-db
          property: host
      
      - key: DB_PORT
        fromDatabase:
          name: dkl-email-service-db
          property: port
      
      - key: DB_USER
        fromDatabase:
          name: dkl-email-service-db
          property: user
      
      - key: DB_PASSWORD
        fromDatabase:
          name: dkl-email-service-db
          property: password
      
      - key: DB_NAME
        fromDatabase:
          name: dkl-email-service-db
          property: database
      
      - key: DB_SSLMODE
        value: require
      
      - key: REDIS_HOST
        fromService:
          name: dkl-email-service-redis
          type: redis
          property: host
      
      - key: REDIS_PORT
        fromService:
          name: dkl-email-service-redis
          type: redis
          property: port
      
      - key: JWT_SECRET
        generateValue: true
      
      - key: JWT_EXPIRATION
        value: 24h
      
      - key: JWT_REFRESH_EXPIRATION
        value: 168h
      
      - key: FRONTEND_URL
        value: https://your-frontend-domain.com
      
      - key: SMTP_HOST
        sync: false
      
      - key: SMTP_PORT
        value: 587
      
      - key: SMTP_USER
        sync: false
      
      - key: SMTP_PASSWORD
        sync: false
      
      - key: CLOUDINARY_CLOUD_NAME
        sync: false
      
      - key: CLOUDINARY_API_KEY
        sync: false
      
      - key: CLOUDINARY_API_SECRET
        sync: false
      
      - key: RATE_LIMIT_REQUESTS
        value: 1000
      
      - key: RATE_LIMIT_WINDOW
        value: 1h
      
      - key: LOG_LEVEL
        value: info
      
      - key: LOG_FORMAT
        value: json

databases:
  - name: dkl-email-service-db
    databaseName: dklemailservice
    plan: starter
    region: frankfurt

redis:
  - name: dkl-email-service-redis
    plan: starter
    region: frankfurt
```

### 5. Deploy via GitHub

```bash
# In Render Dashboard:
1. Click "New +" → "Web Service"
2. Connect your GitHub repository
3. Select branch: main
4. Render will use render.yaml automatically

# Or deploy manually:
git push origin main
# Render auto-deploys on push to main
```

### 6. Run Migrations

```bash
# SSH into Render instance
render ssh dkl-email-service

# Run migrations
./main migrate

# Or use Render shell
go run database/migrations/run_migrations.go
```

### 7. Configure Custom Domain

```yaml
# In Render Dashboard:
1. Go to your service → "Settings" → "Custom Domain"
2. Add your domain: api.yourdomain.com
3. Update DNS records as instructed:
   
   Type: CNAME
   Name: api
   Value: your-service.onrender.com
```

### 8. Enable Health Checks

Render automatically uses `/api/health` endpoint for health checks.

## Docker Deployment

### 1. Build Docker Image

```bash
# Build production image
docker build -t dkl-email-service:latest .

# Or with specific version
docker build -t dkl-email-service:v1.0.0 .
```

### 2. Docker Compose Production

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    image: dkl-email-service:latest
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=dklemailservice
      - DB_SSLMODE=require
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    networks:
      - dkl-network

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=dklemailservice
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - dkl-network

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    networks:
      - dkl-network

volumes:
  postgres_data:

networks:
  dkl-network:
    driver: bridge
```

### 3. Deploy with Docker Compose

```bash
# Start services
docker-compose up -d

# Check logs
docker-compose logs -f

# Run migrations
docker-compose exec app ./main migrate

# Stop services
docker-compose down
```

### 4. Docker Swarm Deployment

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml dkl

# Check services
docker service ls

# Scale service
docker service scale dkl_app=3

# Update service
docker service update --image dkl-email-service:v1.0.1 dkl_app

# Remove stack
docker stack rm dkl
```

## Environment Configuration

### Production Environment Variables

```bash
# Server
ENV=production
PORT=8080

# Database (use connection string from Render)
DB_HOST=dpg-xxxxx-a.frankfurt-postgres.render.com
DB_PORT=5432
DB_USER=dkl_admin
DB_PASSWORD=<strong-random-password>
DB_NAME=dklemailservice
DB_SSLMODE=require

# Redis (use connection string from Render)
REDIS_HOST=red-xxxxx-a.frankfurt-redis.render.com
REDIS_PORT=6379
REDIS_PASSWORD=<redis-password>
REDIS_DB=0

# JWT (generate secure random strings)
JWT_SECRET=<minimum-32-char-random-string>
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=notifications@yourdomain.com
SMTP_PASSWORD=<app-specific-password>
SMTP_FROM=noreply@yourdomain.com

# Cloudinary
CLOUDINARY_CLOUD_NAME=<your-cloud-name>
CLOUDINARY_API_KEY=<your-api-key>
CLOUDINARY_API_SECRET=<your-api-secret>

# CORS
FRONTEND_URL=https://yourdomain.com,https://www.yourdomain.com

# Rate Limiting
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=1h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Generating Secure Secrets

```bash
# Generate JWT secret (32+ characters)
openssl rand -base64 32

# Generate password
openssl rand -base64 24

# Or use Go
go run -c 'import "crypto/rand"; import "encoding/base64"; b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b))'
```

## Security Hardening

### 1. SSL/TLS Configuration

Render handles SSL automatically. For custom deployments:

```go
// main.go - Force HTTPS in production
if os.Getenv("ENV") == "production" {
    srv := &http.Server{
        Addr:         ":443",
        Handler:      router,
        TLSConfig:    tlsConfig,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

### 2. Security Headers

```go
// Add security middleware
router.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        next.ServeHTTP(w, r)
    })
})
```

### 3. Database Security

```sql
-- Create limited database user
CREATE USER dkl_app WITH PASSWORD 'strong-password';

-- Grant only necessary permissions
GRANT CONNECT ON DATABASE dklemailservice TO dkl_app;
GRANT USAGE ON SCHEMA public TO dkl_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO dkl_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO dkl_app;

-- Revoke dangerous permissions
REVOKE CREATE ON SCHEMA public FROM dkl_app;
REVOKE ALL ON SCHEMA pg_catalog FROM dkl_app;
```

### 4. Rate Limiting

Configure appropriate rate limits:

```go
// Production rate limits
rateLimits := map[string]int{
    "auth.login":    5,   // 5 per 15 minutes
    "auth.register": 3,   // 3 per hour
    "api.default":   1000, // 1000 per hour
}
```

## Monitoring & Logging

### 1. Health Checks

```bash
# Health check endpoint
curl https://api.yourdomain.com/api/health

# Expected response
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "version": "1.0.0",
  "uptime": "2h45m"
}
```

### 2. Application Logs

Render automatically captures logs:

```bash
# View logs in Render Dashboard
# Or via CLI
render logs dkl-email-service --tail

# Filter by level
render logs dkl-email-service --tail | grep ERROR
```

### 3. Error Tracking

Consider integrating Sentry:

```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
    Dsn: os.Getenv("SENTRY_DSN"),
    Environment: os.Getenv("ENV"),
})
```

### 4. Performance Monitoring

Monitor key metrics:
- Response times
- Error rates
- Database query performance
- Memory usage
- CPU usage

## Database Maintenance

### Automated Backups

Render provides automatic daily backups. For manual backups:

```bash
# Create backup
pg_dump -U dkl_admin -h hostname -d dklemailservice > backup.sql

# Restore backup
psql -U dkl_admin -h hostname -d dklemailservice < backup.sql
```

### Database Optimization

```sql
-- Run weekly
VACUUM ANALYZE;

-- Update statistics
ANALYZE;

-- Reindex if needed
REINDEX DATABASE dklemailservice;
```

## Rollback Procedures

### Application Rollback

```bash
# Render: Rollback to previous deployment
# In Dashboard → Deployments → Click "Redeploy" on previous version

# Docker: Rollback image
docker service update --image dkl-email-service:v1.0.0 dkl_app
```

### Database Rollback

```bash
# Restore from backup
psql -U dkl_admin -d dklemailservice < backup_before_migration.sql

# Or manual rollback of migration
psql -U dkl_admin -d dklemailservice
DELETE FROM schema_migrations WHERE version = 'V16__problem_migration';
# Run compensating SQL
```

## Scaling

### Vertical Scaling

Upgrade Render plan or increase Docker resources:

```yaml
# render.yaml
services:
  - type: web
    plan: standard  # or pro
```

### Horizontal Scaling

```bash
# Render: Scale instances in Dashboard
# Settings → Autoscaling → Enable

# Docker Swarm
docker service scale dkl_app=3
```

### Database Scaling

- Enable read replicas
- Connection pooling (already implemented)
- Caching strategy (Redis)

## Troubleshooting

### Application Won't Start

```bash
# Check logs
render logs dkl-email-service --tail

# Common issues:
# 1. Missing environment variables
# 2. Database connection failed
# 3. Port already in use
# 4. Migration failures
```

### Database Connection Issues

```bash
# Test connection
psql "postgresql://user:pass@host:5432/db?sslmode=require"

# Verify SSL mode
# Production must use sslmode=require
```

### High Memory Usage

```bash
# Check container stats
docker stats

# Optimize Go application
GOMEMLIMIT=500MiB go buil