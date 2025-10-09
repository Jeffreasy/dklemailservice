# DKL Email Service - Documentatie

Welkom bij de documentatie van de DKL Email Service. Deze service biedt een complete oplossing voor email management, inclusief contact formulieren, aanmeldingen, automatische email ophaling, en notificaties.

## 📚 Documentatie Structuur

### [Architecture](./architecture/)
Technische architectuur en systeemontwerp documentatie:
- **[Components](./architecture/components.md)** - Overzicht van alle systeemcomponenten
- **[Authentication & Authorization](./architecture/authentication-and-authorization.md)** - JWT authenticatie en RBAC systeem
- **[Database Schema](./architecture/database-schema.md)** - Database structuur en relaties
- **[Email System](./architecture/email-system.md)** - Email service architectuur

### [API Documentation](./api/)
Complete API referentie met voorbeelden:
- **[REST API](./api/rest-api.md)** - Alle REST endpoints met request/response voorbeelden
- **[Authentication](./api/authentication.md)** - Login, logout, token refresh
- **[Email Endpoints](./api/email-endpoints.md)** - Contact en aanmelding emails
- **[Admin Endpoints](./api/admin-endpoints.md)** - Beheer endpoints
- **[WebSocket API](./api/websocket-api.md)** - Real-time chat functionaliteit

### [Guides](./guides/)
Praktische handleidingen voor ontwikkelaars:
- **[Getting Started](./guides/getting-started.md)** - Snelstart gids
- **[Development](./guides/development.md)** - Ontwikkelomgeving setup
- **[Deployment](./guides/deployment.md)** - Productie deployment
- **[Testing](./guides/testing.md)** - Test procedures en strategieën
- **[Security](./guides/security.md)** - Security best practices
- **[Monitoring](./guides/monitoring.md)** - Monitoring en logging setup

### [Reports](./reports/)
Technische rapporten en analyses:
- **[Features Audit](./reports/features-audit.md)** - Overzicht van alle features
- **[Performance Analysis](./reports/performance-analysis.md)** - Performance metrics
- **[Security Audit](./reports/security-audit.md)** - Security assessment

## 🚀 Quick Start

```bash
# Clone repository
git clone https://github.com/Jeffreasy/dklemailservice.git
cd dklemailservice

# Configureer environment
cp .env.example .env
# Bewerk .env met jouw configuratie

# Start de service
go run main.go
```

## 🔑 Belangrijkste Features

- **Email Management**: Contact formulieren en aanmeldingen
- **Automatische Email Ophaling**: IMAP integratie voor inkomende emails
- **JWT Authenticatie**: Veilige gebruikersauthenticatie
- **RBAC Systeem**: Role-Based Access Control
- **Real-time Chat**: WebSocket-based chat functionaliteit
- **Newsletter System**: Geautomatiseerde nieuwsbrieven
- **Telegram Notificaties**: Real-time notificaties via Telegram
- **Prometheus Metrics**: Uitgebreide monitoring
- **Rate Limiting**: Bescherming tegen misbruik

## 📊 Technologie Stack

- **Backend**: Go 1.21+ met Fiber framework
- **Database**: PostgreSQL met GORM
- **Cache**: Redis voor rate limiting en sessions
- **Email**: SMTP met HTML templates
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logging met ELK stack support

## 🔗 Gerelateerde Projecten

- [DKL Admin Panel](https://github.com/Jeffreasy/DKL25AdminPanel) - Frontend admin interface

## 📝 Versie Informatie

- **Huidige Versie**: 1.0.0
- **Go Versie**: 1.21+
- **API Versie**: v1

## 🤝 Contributing

Zie [Development Guide](./guides/development.md) voor informatie over bijdragen aan dit project.

## 📄 Licentie

Copyright © 2024 De Koninklijke Loop. Alle rechten voorbehouden.