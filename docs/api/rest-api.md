# REST API Reference

Complete API referentie voor de DKL Email Service met daadwerkelijke code voorbeelden uit de codebase.

## Base URL

```
Production: https://api.dekoninklijkeloop.nl
Development: http://localhost:8080
```

## Authenticatie

De meeste endpoints vereisen authenticatie via JWT Bearer tokens.

**Header Format:**
```http
Authorization: Bearer <your-jwt-token>
```

Zie [Authentication API](./authentication.md) voor login en token management.

## Response Formats

### Success Response

```json
{
    "success": true,
    "message": "Operatie succesvol",
    "data": { ... }
}
```

### Error Response

```json
{
    "success": false,
    "error": "Foutmelding",
    "code": "ERROR_CODE"
}
```

## Rate Limiting

**Headers in Response:**
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

**Rate Limit Exceeded (429):**
```json
{
    "success": false,
    "error": "Te veel verzoeken, probeer het later opnieuw"
}
```

## Endpoints Overview

### Public Endpoints

| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| GET | `/` | Service informatie |
| GET | `/api/health` | Health check |
| POST | `/api/contact-email` | Contact formulier |
| POST | `/api/aanmelding-email` | Aanmelding formulier |

### Authentication Endpoints

| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/logout` | Logout |
| POST | `/api/auth/refresh` | Token refresh |
| GET | `/api/auth/profile` | Gebruikersprofiel |
| POST | `/api/auth/reset-password` | Wachtwoord wijzigen |

### Admin Endpoints

| Method | Endpoint | Beschrijving | Permission |
|--------|----------|--------------|------------|
| GET | `/api/contact` | Contact lijst | `contacts:read` |
| GET | `/api/contact/:id` | Contact details | `contacts:read` |
| PUT | `/api/contact/:id` | Contact bijwerken | `contacts:update` |
| DELETE | `/api/contact/:id` | Contact verwijderen | `contacts:delete` |
| POST | `/api/contact/:id/antwoord` | Antwoord toevoegen | `contacts:update` |
| GET | `/api/aanmelding` | Aanmelding lijst | `aanmeldingen:read` |
| GET | `/api/aanmelding/:id` | Aanmelding details | `aanmeldingen:read` |
| PUT | `/api/aanmelding/:id` | Aanmelding bijwerken | `aanmeldingen:update` |
| DELETE | `/api/aanmelding/:id` | Aanmelding verwijderen | `aanmeldingen:delete` |

### Mail Management Endpoints

| Method | Endpoint | Beschrijving | Permission |
|--------|----------|--------------|------------|
| GET | `/api/mail` | Inkomende emails | `mail:read` |
| GET | `/api/mail/:id` | Email details | `mail:read` |
| PUT | `/api/mail/:id/processed` | Markeer als verwerkt | `mail:update` |
| DELETE | `/api/mail/:id` | Email verwijderen | `mail:delete` |
| POST | `/api/mail/fetch` | Handmatig ophalen | `mail:manage` |
| GET | `/api/mail/unprocessed` | Onverwerkte emails | `mail:read` |

## Detailed Endpoints

### Root Endpoint

#### GET /

Service informatie en beschikbare endpoints.

**Response (200 OK):**
```json
{
    "service": "DKL Email Service API",
    "version": "1.0.0",
    "status": "running",
    "environment": "production",
    "timestamp": "2024-03-20T15:04:05Z",
    "endpoints": [
        {
            "path": "/api/health",
            "method": "GET",
            "description": "Service health status"
        },
        {
            "path": "/api/contact-email",
            "method": "POST",
            "description": "Send contact form email"
        }
        // ... meer endpoints
    ]
}
```

**Implementatie:** [`main.go:337`](../../main.go:337)

### Health Check

#### GET /api/health

Controleert de status van de service en dependencies.

**Response (200 OK):**
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
        "email_auto_fetcher": "running"
    }
}
```

**Response (503 Service Unavailable):**
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

**Implementatie:** [`handlers/health_handler.go`](../../handlers/health_handler.go:1)

### Contact Email

#### POST /api/contact-email

Verzendt een contact formulier email naar admin en bevestiging naar gebruiker.

**Request Body:**
```json
{
    "naam": "John Doe",
    "email": "john@example.com",
    "bericht": "Hallo, ik heb een vraag over het evenement.",
    "privacy_akkoord": true,
    "test_mode": false
}
```

**Validatie Regels:**
- `naam`: Verplicht, min 2 karakters
- `email`: Verplicht, geldig email formaat
- `bericht`: Verplicht, min 10 karakters
- `privacy_akkoord`: Moet `true` zijn
- `test_mode`: Optioneel, voor testing (geen echte emails)

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Je bericht is verzonden! Je ontvangt ook een bevestiging per email."
}
```

**Response (200 OK - Test Mode):**
```json
{
    "success": true,
    "message": "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
    "test_mode": true
}
```

**Response (400 Bad Request):**
```json
{
    "success": false,
    "error": "Naam, email en bericht zijn verplicht"
}
```

**Response (429 Too Many Requests):**
```json
{
    "success": false,
    "error": "Te veel emails in korte tijd, probeer het later opnieuw"
}
```

**Rate Limits:**
- Global: 100 emails per uur
- Per IP: 5 emails per uur

**Implementatie:** [`handlers/email_handler.go:43`](../../handlers/email_handler.go:43)

**cURL Voorbeeld:**
```bash
curl -X POST https://api.dekoninklijkeloop.nl/api/contact-email \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "John Doe",
    "email": "john@example.com",
    "bericht": "Hallo, ik heb een vraag over het evenement.",
    "privacy_akkoord": true
  }'
```

**JavaScript Voorbeeld:**
```javascript
const response = await fetch('https://api.dekoninklijkeloop.nl/api/contact-email', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        naam: 'John Doe',
        email: 'john@example.com',
        bericht: 'Hallo, ik heb een vraag over het evenement.',
        privacy_akkoord: true
    })
});

const data = await response.json();
console.log(data);
```

### Aanmelding Email

#### POST /api/aanmelding-email

Verzendt een aanmelding email en slaat de aanmelding op in de database.

**Request Body:**
```json
{
    "naam": "John Doe",
    "email": "john@example.com",
    "telefoon": "0612345678",
    "rol": "loper",
    "afstand": "10km",
    "ondersteuning": "geen",
    "bijzonderheden": "Vegetarisch eten graag",
    "terms": true,
    "test_mode": false
}
```

**Validatie Regels:**
- `naam`: Verplicht
- `email`: Verplicht, geldig email formaat
- `telefoon`: Optioneel
- `rol`: Verplicht, een van: `loper`, `vrijwilliger`
- `afstand`: Verplicht voor lopers, een van: `5km`, `10km`, `21.1km`
- `ondersteuning`: Optioneel
- `bijzonderheden`: Optioneel
- `terms`: Moet `true` zijn
- `test_mode`: Optioneel, voor testing

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Je aanmelding is verzonden! Je ontvangt ook een bevestiging per email."
}
```

**Response (400 Bad Request):**
```json
{
    "success": false,
    "error": "Naam is verplicht"
}
```

**Rate Limits:**
- Global: 200 emails per uur
- Per IP: 10 emails per uur

**Implementatie:** [`handlers/email_handler.go:205`](../../handlers/email_handler.go:205)

**cURL Voorbeeld:**
```bash
curl -X POST https://api.dekoninklijkeloop.nl/api/aanmelding-email \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "John Doe",
    "email": "john@example.com",
    "telefoon": "0612345678",
    "rol": "loper",
    "afstand": "10km",
    "terms": true
  }'
```

### Contact Beheer

#### GET /api/contact

Haalt een lijst van contact formulieren op (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 50)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)
- `status` (optioneel): Filter op status (`nieuw`, `in_behandeling`, `afgehandeld`)

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "naam": "John Doe",
            "email": "john@example.com",
            "bericht": "Vraag over het evenement",
            "privacy_akkoord": true,
            "status": "nieuw",
            "created_at": "2024-03-20T15:04:05Z",
            "updated_at": "2024-03-20T15:04:05Z"
        }
    ],
    "total": 1,
    "limit": 50,
    "offset": 0
}
```

**Implementatie:** [`handlers/contact_handler.go`](../../handlers/contact_handler.go:1)

#### GET /api/contact/:id

Haalt details van een specifiek contact formulier op.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "naam": "John Doe",
        "email": "john@example.com",
        "bericht": "Vraag over het evenement",
        "privacy_akkoord": true,
        "status": "nieuw",
        "antwoorden": [
            {
                "id": "660e8400-e29b-41d4-a716-446655440001",
                "bericht": "Bedankt voor je vraag...",
                "created_by": "admin@example.com",
                "created_at": "2024-03-20T16:00:00Z"
            }
        ],
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T16:00:00Z"
    }
}
```

**Response (404 Not Found):**
```json
{
    "success": false,
    "error": "Contact formulier niet gevonden"
}
```

#### PUT /api/contact/:id

Werkt een contact formulier bij.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "status": "in_behandeling",
    "notities": "Klant teruggebeld"
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Contact formulier bijgewerkt",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "status": "in_behandeling",
        "updated_at": "2024-03-20T16:30:00Z"
    }
}
```

#### POST /api/contact/:id/antwoord

Voegt een antwoord toe aan een contact formulier.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "bericht": "Bedankt voor je vraag. Het evenement vindt plaats op...",
    "send_email": true
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Antwoord toegevoegd",
    "data": {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "bericht": "Bedankt voor je vraag...",
        "created_by": "admin@example.com",
        "created_at": "2024-03-20T16:00:00Z",
        "email_sent": true
    }
}
```

### Mail Management

#### GET /api/mail

Haalt inkomende emails op die door EmailAutoFetcher zijn opgehaald.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 10)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)
- `account_type` (optioneel): Filter op account (`info`, `inschrijving`)
- `is_processed` (optioneel): Filter op verwerkt status (`true`, `false`)

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "770e8400-e29b-41d4-a716-446655440000",
            "message_id": "<message123@example.com>",
            "from": "sender@example.com",
            "to": "info@dekoninklijkeloop.nl",
            "subject": "Vraag over het evenement",
            "body": "Hallo, ik heb een vraag...",
            "content_type": "text/plain",
            "received_at": "2024-04-01T09:30:00Z",
            "uid": "AAABBCCC123",
            "account_type": "info",
            "is_processed": false,
            "processed_at": null,
            "created_at": "2024-04-01T09:35:00Z",
            "updated_at": "2024-04-01T09:35:00Z"
        }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
}
```

**Implementatie:** [`handlers/mail_handler.go`](../../handlers/mail_handler.go:1)

#### POST /api/mail/fetch

Haalt handmatig nieuwe emails op van de mailserver.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "emails_found": 5,
    "emails_saved": 3,
    "last_run": "2024-04-02T14:25:00Z",
    "message": "Emails succesvol opgehaald"
}
```

**Implementatie:** [`handlers/mail_handler.go`](../../handlers/mail_handler.go:1)

### Metrics Endpoints

#### GET /api/metrics/email

Haalt email verzend statistieken op (admin only).

**Headers:**
```http
X-API-Key: <admin-api-key>
```

**Response (200 OK):**
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

#### GET /metrics

Prometheus metrics endpoint.

**Response (200 OK):**
```
# HELP email_sent_total Total number of emails sent
# TYPE email_sent_total counter
email_sent_total{type="contact",status="success"} 50
email_sent_total{type="aanmelding",status="success"} 100

# HELP email_failed_total Total number of failed emails
# TYPE email_failed_total counter
email_failed_total{type="contact",reason="smtp_error"} 1
email_failed_total{type="aanmelding",reason="rate_limited"} 2

# HELP email_latency_seconds Email sending latency
# TYPE email_latency_seconds histogram
email_latency_seconds_bucket{type="contact",le="0.1"} 10
email_latency_seconds_bucket{type="contact",le="0.5"} 45
email_latency_seconds_bucket{type="contact",le="1"} 50
```

## Error Codes

| HTTP Status | Code | Beschrijving |
|-------------|------|--------------|
| 400 | `INVALID_INPUT` | Ongeldige invoer data |
| 401 | `NO_AUTH_HEADER` | Geen Authorization header |
| 401 | `INVALID_AUTH_HEADER` | Ongeldige Authorization header |
| 401 | `TOKEN_EXPIRED` | Token is verlopen |
| 401 | `TOKEN_MALFORMED` | Token heeft ongeldige structuur |
| 401 | `TOKEN_SIGNATURE_INVALID` | Token signature is ongeldig |
| 401 | `INVALID_TOKEN` | Algemene token validatie fout |
| 403 | `FORBIDDEN` | Geen toegang tot resource |
| 404 | `NOT_FOUND` | Resource niet gevonden |
| 429 | `RATE_LIMIT_EXCEEDED` | Rate limit overschreden |
| 500 | `INTERNAL_ERROR` | Interne server fout |

## CORS

**Toegestane Origins:**
```
https://www.dekoninklijkeloop.nl
https://dekoninklijkeloop.nl
https://admin.dekoninklijkeloop.nl
http://localhost:3000
http://localhost:5173
```

**Toegestane Headers:**
```
Origin, Content-Type, Accept, Authorization, X-Test-Mode
```

**Toegestane Methods:**
```
GET, POST, PUT, DELETE, OPTIONS
```

**Configuratie:** [`main.go:303`](../../main.go:303)

## Test Mode

Voor testing kunnen endpoints worden aangeroepen in test mode. Dit voorkomt dat echte emails worden verzonden.

**Activatie via Header:**
```http
X-Test-Mode: true
```

**Activatie via Body:**
```json
{
    "test_mode": true,
    ...
}
```

**Response in Test Mode:**
```json
{
    "success": true,
    "message": "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
    "test_mode": true
}
```

**Implementatie:** [`handlers/middleware.go:174`](../../handlers/middleware.go:174)

## Pagination

Endpoints die lijsten teruggeven ondersteunen paginatie.

**Query Parameters:**
```
?limit=50&offset=0
```

**Response Format:**
```json
{
    "success": true,
    "data": [...],
    "total": 150,
    "limit": 50,
    "offset": 0,
    "has_more": true
}
```

## Filtering

Veel endpoints ondersteunen filtering via query parameters.

**Voorbeelden:**
```
GET /api/contact?status=nieuw
GET /api/aanmelding?rol=loper&afstand=10km
GET /api/mail?account_type=info&is_processed=false
```

## Sorting

Sorteer resultaten met de `sort` en `order` parameters.

**Voorbeelden:**
```
GET /api/contact?sort=created_at&order=desc
GET /api/aanmelding?sort=naam&order=asc
```

## Versioning

De API gebruikt semantic versioning.

**Version Header:**
```http
X-API-Version: 1.0.0
```

**Deprecated Features:**
```http
X-API-Deprecated-Feature: "oude_endpoint"
X-API-Alternative: "nieuwe_endpoint"
```

## SDK Examples

### Go Client

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type ContactRequest struct {
    Naam           string `json:"naam"`
    Email          string `json:"email"`
    Bericht        string `json:"bericht"`
    PrivacyAkkoord bool   `json:"privacy_akkoord"`
}

func sendContactEmail(contact ContactRequest) error {
    jsonData, _ := json.Marshal(contact)
    
    resp, err := http.Post(
        "https://api.dekoninklijkeloop.nl/api/contact-email",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var response struct {
        Success bool   `json:"success"`
        Message string `json:"message,omitempty"`
        Error   string `json:"error,omitempty"`
    }
    
    json.NewDecoder(resp.Body).Decode(&response)
    
    if !response.Success {
        return fmt.Errorf(response.Error)
    }
    
    return nil
}
```

### TypeScript Client

```typescript
interface ContactRequest {
    naam: string;
    email: string;
    bericht: string;
    privacy_akkoord: boolean;
}

interface ApiResponse<T = any> {
    success: boolean;
    message?: string;
    error?: string;
    data?: T;
}

async function sendContactEmail(contact: ContactRequest): Promise<ApiResponse> {
    const response = await fetch('https://api.dekoninklijkeloop.nl/api/contact-email', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(contact),
    });
    
    return await response.json();
}
```

## Zie Ook

- [Authentication API](./authentication.md) - Authenticatie endpoints
- [Email Endpoints](./email-endpoints.md) - Email specifieke endpoints
- [Admin Endpoints](./admin-endpoints.md) - Admin beheer endpoints
- [WebSocket API](./websocket-api.md) - Real-time chat API