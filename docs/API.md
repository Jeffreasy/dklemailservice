# API Documentatie

## Overzicht

De DKL Email Service biedt een REST API voor het verzenden van emails voor contact formulieren en aanmeldingen.

## Base URL

```
https://api.dekoninklijkeloop.nl
```

## Endpoints

### Root Endpoint

#### GET /
Basis informatie over de service.

**Response**
```json
{
    "message": "DKL Email Service API",
    "version": "1.0.0"
}
```

### Contact Email

#### POST /api/contact-email
Verzendt een email voor een contact formulier bericht.

**Request**
```json
{
    "naam": "John Doe",
    "email": "john@example.com",
    "bericht": "Hallo, ik heb een vraag over het evenement.",
    "privacy_akkoord": true
}
```

**Response Success (200 OK)**
```json
{
    "success": true,
    "message": "Je bericht is verzonden! Je ontvangt ook een bevestiging per email."
}
```

**Response Error (400 Bad Request)**
```json
{
    "success": false,
    "error": "Naam, email en bericht zijn verplicht"
}
```

**Validatie Regels**
- `naam`: Verplicht
- `email`: Verplicht, moet geldig email formaat zijn
- `bericht`: Verplicht
- `privacy_akkoord`: Moet true zijn

### Aanmelding Email

#### POST /api/aanmelding-email
Verzendt een email voor een nieuwe aanmelding.

**Request**
```json
{
    "naam": "John Doe",
    "email": "john@example.com",
    "telefoon": "0612345678",
    "rol": "deelnemer",
    "afstand": "10km",
    "ondersteuning": "geen",
    "bijzonderheden": "geen",
    "terms": true
}
```

**Response Success (200 OK)**
```json
{
    "success": true,
    "message": "Je aanmelding is verzonden! Je ontvangt ook een bevestiging per email."
}
```

**Response Error (400 Bad Request)**
```json
{
    "success": false,
    "error": "Naam is verplicht"
}
```

**Validatie Regels**
- `naam`: Verplicht
- `email`: Verplicht, moet geldig email formaat zijn
- `terms`: Moet true zijn
- Andere velden zijn optioneel

### Health Check

#### GET /api/health
Controleert de status van de service.

**Response Success (200 OK)**
```json
{
    "status": "ok",
    "version": "1.0.0",
    "timestamp": "2024-03-20T15:04:05Z",
    "uptime": "24h3m12s"
}
```

### Metrics

#### GET /api/metrics/email
Haalt email verzend statistieken op. Vereist authenticatie met API key.

**Headers**
```http
X-API-Key: your-admin-api-key
```

**Response Success (200 OK)**
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

#### GET /api/metrics/rate-limits
Haalt rate limit statistieken op. Vereist authenticatie met API key.

**Headers**
```http
X-API-Key: your-admin-api-key
```

**Response Success (200 OK)**
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

#### GET /metrics
Prometheus metrics endpoint.

### Mail Management

De volgende endpoints zijn beschikbaar voor het beheren van inkomende e-mails. Alle endpoints vereisen administratieve rechten (JWT authenticatie + admin rol).

#### GET /api/mail
Haalt een lijst van inkomende e-mails op.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Parameters**
- `limit` (query, optioneel): Maximum aantal resultaten (default: 10)
- `offset` (query, optioneel): Aantal resultaten om over te slaan (default: 0)

**Response Success (200 OK)**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "message_id": "<message123@example.com>",
    "from": "sender@example.com",
    "to": "info@dekoninklijkeloop.nl",
    "subject": "Vraag over het evenement",
    "body": "Hallo, ik heb een vraag over het evenement...",
    "content_type": "text/plain",
    "received_at": "2024-04-01T09:30:00Z",
    "uid": "AAABBCCC123",
    "account_type": "info",
    "is_processed": false,
    "processed_at": null,
    "created_at": "2024-04-01T09:35:00Z",
    "updated_at": "2024-04-01T09:35:00Z"
  },
  // ... meer e-mails
]
```

#### GET /api/mail/:id
Haalt details van een specifieke e-mail op.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Parameters**
- `id` (path): ID van de e-mail

**Response Success (200 OK)**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "message_id": "<message123@example.com>",
  "from": "sender@example.com",
  "to": "info@dekoninklijkeloop.nl",
  "subject": "Vraag over het evenement",
  "body": "Hallo, ik heb een vraag over het evenement...",
  "content_type": "text/plain",
  "received_at": "2024-04-01T09:30:00Z",
  "uid": "AAABBCCC123",
  "account_type": "info",
  "is_processed": false,
  "processed_at": null,
  "created_at": "2024-04-01T09:35:00Z",
  "updated_at": "2024-04-01T09:35:00Z"
}
```

#### PUT /api/mail/:id/processed
Markeert een e-mail als verwerkt.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Parameters**
- `id` (path): ID van de e-mail

**Response Success (200 OK)**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "message_id": "<message123@example.com>",
  "from": "sender@example.com",
  "to": "info@dekoninklijkeloop.nl",
  "subject": "Vraag over het evenement",
  "body": "Hallo, ik heb een vraag over het evenement...",
  "content_type": "text/plain",
  "received_at": "2024-04-01T09:30:00Z",
  "uid": "AAABBCCC123",
  "account_type": "info",
  "is_processed": true,
  "processed_at": "2024-04-02T14:22:00Z",
  "created_at": "2024-04-01T09:35:00Z",
  "updated_at": "2024-04-02T14:22:00Z"
}
```

#### DELETE /api/mail/:id
Verwijdert een e-mail.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Parameters**
- `id` (path): ID van de e-mail

**Response Success (200 OK)**
```json
{
  "success": true,
  "message": "Email succesvol verwijderd"
}
```

#### POST /api/mail/fetch
Haalt handmatig nieuwe e-mails op van de mailserver.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Response Success (200 OK)**
```json
{
  "success": true,
  "emails_found": 5,
  "emails_saved": 3,
  "last_run": "2024-04-02T14:25:00Z",
  "message": "Emails succesvol opgehaald"
}
```

#### GET /api/mail/unprocessed
Haalt een lijst van onverwerkte e-mails op.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Response Success (200 OK)**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "message_id": "<message456@example.com>",
    "from": "another@example.com",
    "to": "info@dekoninklijkeloop.nl",
    "subject": "Nieuwe vraag",
    "body": "Nieuwe onverwerkte vraag...",
    "content_type": "text/plain",
    "received_at": "2024-04-02T10:15:00Z",
    "uid": "AAABBCCC124",
    "account_type": "info",
    "is_processed": false,
    "processed_at": null,
    "created_at": "2024-04-02T10:20:00Z",
    "updated_at": "2024-04-02T10:20:00Z"
  },
  // ... meer onverwerkte e-mails
]
```

#### GET /api/mail/account/:type
Haalt e-mails gefilterd op account type (info, inschrijving) op.

**Headers**
```http
Authorization: Bearer your-jwt-token
```

**Parameters**
- `type` (path): Type account ("info" of "inschrijving")

**Response Success (200 OK)**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "message_id": "<message789@example.com>",
    "from": "user@example.com",
    "to": "inschrijving@dekoninklijkeloop.nl",
    "subject": "Inschrijving vraag",
    "body": "Vraag over inschrijving...",
    "content_type": "text/plain",
    "received_at": "2024-04-02T11:45:00Z",
    "uid": "AAABBCCC125",
    "account_type": "inschrijving",
    "is_processed": false,
    "processed_at": null,
    "created_at": "2024-04-02T11:50:00Z",
    "updated_at": "2024-04-02T11:50:00Z"
  },
  // ... meer e-mails van het specifieke account type
]
```

## Rate Limiting

### Limieten
- Contact emails:
  - 100 emails per uur globaal
  - 5 emails per uur per IP
- Aanmelding emails:
  - 200 emails per uur globaal
  - 10 emails per uur per IP

### Response (429 Too Many Requests)
```json
{
    "success": false,
    "error": "Te veel emails in korte tijd, probeer het later opnieuw"
}
```

## Error Codes

| HTTP Status | Beschrijving |
|-------------|--------------|
| 400 | Ongeldige invoer data |
| 401 | Ongeautoriseerd (voor metrics endpoints) |
| 429 | Rate limit overschreden |
| 500 | Interne server fout |

## CORS

### Toegestane Origins
```
https://www.dekoninklijkeloop.nl
https://dekoninklijkeloop.nl
```

### Headers
```http
Access-Control-Allow-Origin: [configured-origins]
Access-Control-Allow-Headers: Origin, Content-Type, Accept
Access-Control-Allow-Methods: GET,POST,OPTIONS
```

## Voorbeelden

### cURL

#### Contact Formulier
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

#### Aanmelding
```bash
curl -X POST https://api.dekoninklijkeloop.nl/api/aanmelding-email \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "John Doe",
    "email": "john@example.com",
    "telefoon": "0612345678",
    "rol": "deelnemer",
    "afstand": "10km",
    "ondersteuning": "geen",
    "bijzonderheden": "geen",
    "terms": true
  }'
```

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
    jsonData, err := json.Marshal(contact)
    if err != nil {
        return err
    }
    
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
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return err
    }
    
    if !response.Success {
        return fmt.Errorf(response.Error)
    }
    
    return nil
}
```

## Versioning

De API gebruikt semantic versioning (MAJOR.MINOR.PATCH).

### Versie Headers
```http
X-API-Version: 1.0.0
```

### Breaking Changes
- Major version updates kunnen breaking changes bevatten
- Minor en patch updates zijn backwards compatible
- Deprecated features worden aangekondigd in de headers:
  ```http
  X-API-Deprecated-Feature: "oude_endpoint"
  X-API-Alternative: "nieuwe_endpoint"
  ```

## Security

### TLS
- Minimaal TLS 1.2
- Sterke cipher suites
- HSTS enabled
- Perfect Forward Secrecy

### Input Validatie
- Alle input wordt gevalideerd
- XSS preventie
- SQL injection preventie
- Input length limits

### Rate Limiting
- IP-based limits
- Global limits
- Sliding window implementatie
- Duidelijke error messages

### Error Handling
- Geen stack traces in productie
- Generieke error messages
- Logging van security events
- Monitoring van error rates 