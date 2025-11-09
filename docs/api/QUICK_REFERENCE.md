# API Quick Reference

Snel overzicht van alle beschikbare API endpoints voor de DKL Email Service.

## Base URL

```
Development: http://localhost:8080
Production: https://api.dklemailservice.com
```

---

## Public Endpoints (No Auth Required)

### Health & Info
| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| GET | `/` | Service informatie |
| GET | `/api/health` | Health check |
| GET | `/favicon.ico` | Favicon |

### Email Verzending
| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| POST | `/api/contact-email` | Verstuur contactformulier |
| POST | `/api/register` | Aanmelding + email verzenden |

### Events (Public)
| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| GET | `/api/events` | Lijst van events |
| GET | `/api/events/active` | Huidig actief event |
| GET | `/api/events/:id` | Event details |

### Content (Public)
| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| GET | `/api/videos` | Zichtbare video's |
| GET | `/api/partners` | Zichtbare partners |
| GET | `/api/sponsors` | Zichtbare sponsors |
| GET | `/api/photos` | Zichtbare foto's |
| GET | `/api/albums` | Zichtbare albums |
| GET | `/api/albums/:id` | Album met foto's |
| GET | `/api/radio-recordings` | Zichtbare radio opnames |
| GET | `/api/program-schedule` | Programma schema |
| GET | `/api/social-links` | Social media links |
| GET | `/api/social-embeds` | Social media embeds |
| GET | `/api/title-sections` | Titel secties |
| GET | `/api/under-construction` | Onderhoudsstatus |

---

## Authentication Endpoints

### Auth
| Method | Endpoint | Auth | Beschrijving |
|--------|----------|------|--------------|
| POST | `/api/auth/login` | ❌ | Gebruiker inloggen |
| POST | `/api/auth/logout` | ✅ | Gebruiker uitloggen |
| POST | `/api/auth/refresh` | ❌ | Token vernieuwen |
| GET | `/api/auth/profile` | ✅ | Gebruikersprofiel ophalen |
| POST | `/api/auth/reset-password` | ✅ | Wachtwoord wijzigen |

**Note:** ✅ = JWT token vereist, ❌ = Geen auth vereist

---

## Admin Endpoints (Auth + Permissions Required)

### Contact Management
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/contact` | `contact:read` | Lijst contactformulieren |
| GET | `/api/contact/:id` | `contact:read` | Contact details |
| PUT | `/api/contact/:id` | `contact:write` | Contact bijwerken |
| DELETE | `/api/contact/:id` | `contact:delete` | Contact verwijderen |
| POST | `/api/contact/:id/antwoord` | `contact:write` | Antwoord toevoegen |
| GET | `/api/contact/status/:status` | `contact:read` | Filter op status |

### Participant Management
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/participant` | `participant:read` | Lijst deelnemers |
| GET | `/api/participant/:id` | `participant:read` | Deelnemer details |
| DELETE | `/api/participant/:id` | `participant:delete` | Deelnemer verwijderen |
| POST | `/api/participant/:id/antwoord` | `participant:write` | Antwoord toevoegen |

### Event Registrations
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/registration/:id` | `registration:read` | Registratie details |
| PUT | `/api/registration/:id` | `registration:write` | Registratie bijwerken |
| GET | `/api/registration/rol/:rol` | `registration:read` | Filter op rol |
| GET | `/api/events/:id/registrations` | `event:read` | Event registraties |

### Steps Tracking
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| POST | `/api/registration/:id/steps` | `steps:write` | Stappen bijwerken |
| GET | `/api/registration/:id/dashboard` | `steps:read` | Dashboard ophalen |
| GET | `/api/total-steps` | `steps:read` | Totaal stappen (jaar) |
| GET | `/api/funds-distribution` | `steps:read` | Fondsen verdeling |

### Route Funds
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/route-funds` | `route_fund:read` | Lijst fondsen |
| POST | `/api/route-funds` | `route_fund:write` | Fonds aanmaken |
| PUT | `/api/route-funds/:id` | `route_fund:write` | Fonds bijwerken |
| DELETE | `/api/route-funds/:id` | `route_fund:delete` | Fonds verwijderen |

### Gamification
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/achievements` | `achievement:read` | Lijst achievements |
| GET | `/api/achievements/user/:id` | `achievement:read` | User achievements |
| POST | `/api/achievements` | `achievement:write` | Achievement aanmaken |
| PUT | `/api/achievements/:id` | `achievement:write` | Achievement bijwerken |
| DELETE | `/api/achievements/:id` | `achievement:delete` | Achievement verwijderen |
| GET | `/api/badges` | `badge:read` | Lijst badges |
| GET | `/api/badges/user/:id` | `badge:read` | User badges |
| POST | `/api/badges` | `badge:write` | Badge aanmaken |
| PUT | `/api/badges/:id` | `badge:write` | Badge bijwerken |
| DELETE | `/api/badges/:id` | `badge:delete` | Badge verwijderen |
| POST | `/api/badges/:id/award` | `badge:write` | Badge toekennen |
| DELETE | `/api/badges/:id/remove` | `badge:delete` | Badge verwijderen |
| GET | `/api/leaderboard` | `leaderboard:read` | Ranglijst ophalen |
| GET | `/api/leaderboard/user/:id` | `leaderboard:read` | User statistieken |

### Notulen (Meeting Notes)
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/notulen` | `notulen:read` | Lijst notulen |
| GET | `/api/notulen/:id` | `notulen:read` | Notulen details |
| POST | `/api/notulen` | `notulen:write` | Notulen aanmaken |
| PUT | `/api/notulen/:id` | `notulen:write` | Notulen bijwerken |
| DELETE | `/api/notulen/:id` | `notulen:delete` | Notulen verwijderen |
| PUT | `/api/notulen/:id/finalize` | `notulen:write` | Notulen finaliseren |

### Newsletter
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/newsletter` | `newsletter:read` | Lijst nieuwsbrieven |
| GET | `/api/newsletter/:id` | `newsletter:read` | Nieuwsbrief details |
| POST | `/api/newsletter` | `newsletter:write` | Nieuwsbrief aanmaken |
| PUT | `/api/newsletter/:id` | `newsletter:write` | Nieuwsbrief bijwerken |
| DELETE | `/api/newsletter/:id` | `newsletter:delete` | Nieuwsbrief verwijderen |
| POST | `/api/newsletter/:id/send` | `newsletter:send` | Nieuwsbrief verzenden |

### Mail Management
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/mail` | `email:read` | Lijst inkomende emails |
| GET | `/api/mail/:id` | `email:read` | Email details |
| PUT | `/api/mail/:id/processed` | `email:write` | Markeer als verwerkt |
| DELETE | `/api/mail/:id` | `email:delete` | Email verwijderen |
| POST | `/api/mail/fetch` | `email:fetch` | Nieuwe emails ophalen |
| GET | `/api/mail/unprocessed` | `email:read` | Onverwerkte emails |
| GET | `/api/mail/account/:type` | `email:read` | Filter op account |

### Auto Response Management
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/mail/autoresponse` | Auth required | Lijst auto responses |
| GET | `/api/mail/autoresponse/:id` | Auth required | Auto response details |
| POST | `/api/mail/autoresponse` | Auth required | Auto response aanmaken |
| PUT | `/api/mail/autoresponse/:id` | Auth required | Auto response bijwerken |
| DELETE | `/api/mail/autoresponse/:id` | Auth required | Auto response verwijderen |

### Notifications
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/notifications` | Auth only | Lijst notificaties |
| GET | `/api/notifications/:id` | Auth only | Notificatie details |
| PUT | `/api/notifications/:id/read` | Auth only | Markeer als gelezen |
| PUT | `/api/notifications/read-all` | Auth only | Alle als gelezen |
| DELETE | `/api/notifications/:id` | Auth only | Notificatie verwijderen |
| DELETE | `/api/notifications/read` | Auth only | Gelezen verwijderen |
| POST | `/api/notifications` | `notification:write` | Notificatie aanmaken |
| POST | `/api/notifications/broadcast` | `notification:write` | Broadcast notificatie |

### Chat System
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/chat/channels` | `chat:read` | Gebruikers kanalen |
| GET | `/api/chat/public-channels` | `chat:read` | Publieke kanalen |
| POST | `/api/chat/channels` | `chat:write` | Kanaal aanmaken |
| POST | `/api/chat/direct` | `chat:write` | Direct chat aanmaken |
| POST | `/api/chat/channels/:id/join` | `chat:write` | Kanaal joinen |
| POST | `/api/chat/channels/:id/leave` | `chat:write` | Kanaal verlaten |
| GET | `/api/chat/channels/:id/messages` | `chat:read` | Berichten ophalen |
| POST | `/api/chat/channels/:id/messages` | `chat:write` | Bericht verzenden |
| PUT | `/api/chat/messages/:id` | `chat:write` | Bericht bewerken |
| DELETE | `/api/chat/messages/:id` | `chat:moderate` | Bericht verwijderen |
| POST | `/api/chat/messages/:id/reactions` | `chat:write` | Reactie toevoegen |
| DELETE | `/api/chat/messages/:id/reactions/:emoji` | `chat:write` | Reactie verwijderen |
| PUT | `/api/chat/presence` | `chat:write` | Presence bijwerken |
| GET | `/api/chat/online-users` | `chat:read` | Online gebruikers |

### User Management
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/users` | `user:read` | Lijst gebruikers |
| GET | `/api/users/:id` | `user:read` | Gebruiker details |
| POST | `/api/users` | `user:write` | Gebruiker aanmaken |
| PUT | `/api/users/:id` | `user:write` | Gebruiker bijwerken |
| DELETE | `/api/users/:id` | `user:delete` | Gebruiker verwijderen |
| GET | `/api/users/:id/roles` | `user:read` | Gebruikersrollen |
| POST | `/api/users/:id/roles` | `user:manage_roles` | Rol toewijzen |
| POST | `/api/users/:id/roles/batch` | `user:manage_roles` | Meerdere rollen toewijzen |
| DELETE | `/api/users/:id/roles/:role_id` | `user:manage_roles` | Rol verwijderen |

### RBAC Management
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/permissions` | `permission:read` | Lijst permissies |
| GET | `/api/permissions/:id` | `permission:read` | Permissie details |
| POST | `/api/permissions` | `permission:write` | Permissie aanmaken |
| PUT | `/api/permissions/:id` | `permission:write` | Permissie bijwerken |
| DELETE | `/api/permissions/:id` | `permission:delete` | Permissie verwijderen |
| GET | `/api/permissions/check` | Auth only | Permissie controleren |
| GET | `/api/roles` | `role:read` | Lijst rollen |
| GET | `/api/roles/:id` | `role:read` | Rol details |
| POST | `/api/roles` | `role:write` | Rol aanmaken |
| PUT | `/api/roles/:id` | `role:write` | Rol bijwerken |
| DELETE | `/api/roles/:id` | `role:delete` | Rol verwijderen |
| GET | `/api/roles/:id/permissions` | `role:read` | Rol permissies |
| POST | `/api/roles/:id/permissions` | `role:write` | Permissies toewijzen |
| PUT | `/api/roles/:id/permissions` | `role:write` | Batch permissie update |
| DELETE | `/api/roles/:id/permissions/:pid` | `role:write` | Permissie verwijderen |

### Events Management (Admin)
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| POST | `/api/events` | `event:write` | Event aanmaken |
| PUT | `/api/events/:id` | `event:write` | Event bijwerken |
| DELETE | `/api/events/:id` | `event:delete` | Event verwijderen |
| GET | `/api/events/:id/stats` | `event:read` | Event statistieken |

### CMS - Videos
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/videos/admin` | `video:read` | Alle video's (admin) |
| GET | `/api/videos/:id` | `video:read` | Video details |
| POST | `/api/videos` | `video:write` | Video aanmaken |
| PUT | `/api/videos/:id` | `video:write` | Video bijwerken |
| DELETE | `/api/videos/:id` | `video:delete` | Video verwijderen |

### CMS - Partners
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/partners/admin` | `partner:read` | Alle partners (admin) |
| GET | `/api/partners/:id` | `partner:read` | Partner details |
| POST | `/api/partners` | `partner:write` | Partner aanmaken |
| PUT | `/api/partners/:id` | `partner:write` | Partner bijwerken |
| DELETE | `/api/partners/:id` | `partner:delete` | Partner verwijderen |

### CMS - Sponsors
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/sponsors/admin` | `sponsor:read` | Alle sponsors (admin) |
| GET | `/api/sponsors/:id` | `sponsor:read` | Sponsor details |
| POST | `/api/sponsors` | `sponsor:write` | Sponsor aanmaken (+ upload) |
| PUT | `/api/sponsors/:id` | `sponsor:write` | Sponsor bijwerken |
| DELETE | `/api/sponsors/:id` | `sponsor:delete` | Sponsor verwijderen |

### CMS - Photos & Albums
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/photos/admin` | `photo:read` | Alle foto's (admin) |
| POST | `/api/photos` | `photo:write` | Foto aanmaken |
| PUT | `/api/photos/:id` | `photo:write` | Foto bijwerken |
| DELETE | `/api/photos/:id` | `photo:delete` | Foto verwijderen |
| GET | `/api/albums/admin` | `album:read` | Alle albums (admin) |
| POST | `/api/albums` | `album:write` | Album aanmaken |
| PUT | `/api/albums/:id` | `album:write` | Album bijwerken |
| DELETE | `/api/albums/:id` | `album:delete` | Album verwijderen |
| POST | `/api/albums/:id/photos` | `album:write` | Foto aan album toevoegen |
| PUT | `/api/albums/:id/photos/reorder` | `album:write` | Foto's herordenen |
| PUT | `/api/albums/reorder` | `album:write` | Albums herordenen |

### CMS - Other Content
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/radio-recordings/admin` | `radio_recording:read` | Alle radio opnames |
| POST | `/api/radio-recordings` | `radio_recording:write` | Radio opname aanmaken |
| PUT | `/api/radio-recordings/:id` | `radio_recording:write` | Radio opname bijwerken |
| DELETE | `/api/radio-recordings/:id` | `radio_recording:delete` | Radio opname verwijderen |
| POST | `/api/program-schedule` | `program_schedule:write` | Schema item aanmaken |
| PUT | `/api/program-schedule/:id` | `program_schedule:write` | Schema item bijwerken |
| DELETE | `/api/program-schedule/:id` | `program_schedule:delete` | Schema item verwijderen |
| POST | `/api/social-links` | `social_link:write` | Social link aanmaken |
| PUT | `/api/social-links/:id` | `social_link:write` | Social link bijwerken |
| DELETE | `/api/social-links/:id` | `social_link:delete` | Social link verwijderen |
| POST | `/api/social-embeds` | `social_embed:write` | Social embed aanmaken |
| PUT | `/api/social-embeds/:id` | `social_embed:write` | Social embed bijwerken |
| DELETE | `/api/social-embeds/:id` | `social_embed:delete` | Social embed verwijderen |
| POST | `/api/title-sections` | `title_section:write` | Titel sectie aanmaken |
| PUT | `/api/title-sections/:id` | `title_section:write` | Titel sectie bijwerken |
| DELETE | `/api/title-sections/:id` | `title_section:delete` | Titel sectie verwijderen |
| POST | `/api/under-construction` | `under_construction:write` | Onderhoud instellen |
| PUT | `/api/under-construction/:id` | `under_construction:write` | Onderhoud bijwerken |
| DELETE | `/api/under-construction/:id` | `under_construction:delete` | Onderhoud verwijderen |

### Image Upload
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| POST | `/api/upload/image` | `image:write` | Enkele afbeelding uploaden |
| POST | `/api/upload/images/batch` | `image:write` | Meerdere afbeeldingen |

---

## WebSocket Endpoints

### Connection URLs
| Endpoint | Auth | Beschrijving |
|----------|------|--------------|
| `/ws/steps?token=JWT` | JWT query param | Steps real-time updates |
| `/api/ws/notulen?token=JWT` | JWT query param | Notulen real-time updates |
| `/api/chat/ws/:channel_id?token=JWT` | JWT query param | Chat kanaal WebSocket |

### WebSocket Stats (Admin)
| Method | Endpoint | Permission | Beschrijving |
|--------|----------|-----------|--------------|
| GET | `/api/ws/stats` | `admin:read` | WebSocket statistieken |

---

## Metrics & Monitoring

### Metrics Endpoints (API Key Required)
| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| GET | `/api/metrics/email` | Email verzend statistieken |
| GET | `/api/metrics/rate-limits` | Rate limit status |
| GET | `/metrics` | Prometheus metrics |

**Header Required:**
```
X-API-Key: your_api_key
```

---

## Whisky for Charity

### WFC Orders (API Key Required)
| Method | Endpoint | Beschrijving |
|--------|----------|--------------|
| POST | `/api/wfc/order-email` | WFC order email verzenden |

**Header Required:**
```
X-API-Key: your_wfc_api_key
```

---

## Telegram Bot (Admin)

| Method | Endpoint | Auth | Beschrijving |
|--------|----------|------|--------------|
| GET | `/api/v1/telegrambot/config` | ✅ | Bot configuratie |
| POST | `/api/v1/telegrambot/send` | ✅ | Bericht verzenden |
| GET | `/api/v1/telegrambot/commands` | ✅ | Beschikbare commands |

---

## Request/Response Formats

### Standard Success Response

```json
{
  "success": true,
  "data": { },
  "message": "Optional message"
}
```

### Standard Error Response

```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

### Paginated Response

```json
{
  "success": true,
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## Common Query Parameters

| Parameter | Type | Beschrijving | Default |
|-----------|------|--------------|---------|
| `page` | integer | Pagina nummer | 1 |
| `limit` | integer | Resultaten per pagina | 10-20 |
| `offset` | integer | Aantal over te slaan | 0 |
| `sort` | string | Sorteer veld | varies |
| `order` | string | Sorteer richting (asc/desc) | desc |
| `search` | string | Zoek query | - |
| `status` | string | Filter op status | - |
| `type` | string | Filter op type | - |

---

## HTTP Status Codes

| Code | Betekenis |
|------|-----------|
| 200 | OK - Request succesvol |
| 201 | Created - Resource aangemaakt |
| 204 | No Content - Succesvol zonder response body |
| 400 | Bad Request - Ongeldige input |
| 401 | Unauthorized - Authenticatie vereist |
| 403 | Forbidden - Onvoldoende permissies |
| 404 | Not Found - Resource niet gevonden |
| 409 | Conflict - Resource conflict |
| 422 | Unprocessable Entity - Validatie fout |
| 429 | Too Many Requests - Rate limit overschreden |
| 500 | Internal Server Error - Server fout |

---

## Rate Limits

### Default Limits

| Endpoint Type | Requests | Window | Per |
|---------------|----------|--------|-----|
| Contact email | 5 | 1 uur | IP |
| Registration email | 3 | 1 uur | IP |
| Auth login | 5 | 15 min | IP |
| API default | 1000 | 1 uur | Token |
| Admin | 5000 | 1 uur | Token |

### Rate Limit Headers

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

---

## CORS Configuration

### Allowed Origins (Default)

```
https://www.dekoninklijkeloop.nl
https://dekoninklijkeloop.nl
https://admin.dekoninklijkeloop.nl
http://localhost:3000
http://localhost:5173
```

### Allowed Headers

```
Origin
Content-Type
Accept
Authorization
X-Test-Mode
X-API-Key
```

---

## Quick Start Examples

### Login & Authenticated Request

```bash
# 1. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Response bevat access_token

# 2. Authenticated request
curl http://localhost:8080/api/contact \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Create Event Registration

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "naam": "John Doe",
    "email": "john@example.com",
    "telefoon": "+31612345678",
    "rol": "deelnemer",
    "afstand": "10km",
    "privacy_akkoord": true
  }'
```

### Update Steps

```bash
curl -X POST http://localhost:8080/api/registration/UUID/steps \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"delta_steps": 1000}'
```

### Upload Image

```bash
curl -X POST http://localhost:8080/api/upload/image \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@photo.jpg" \
  -F "folder=event_photos"
```

---

## Detailed Documentation

Voor uitgebreide API documentatie, zie:

- [Authentication API](./AUTHENTICATION.md)
- [Permissions & RBAC](./PERMISSIONS.md)
- [Events & Registrations](./EVENTS.md)
- [Steps & Gamification](./STEPS_GAMIFICATION.md)
- [CMS APIs](./CMS.md)
- [Notifications](./NOTIFICATIONS.md)
- [WebSocket APIs](./WEBSOCKET.md)

---

**Last Updated:** 2025-01-08