# REST API Reference

Complete API referentie voor de DKL Email Service met daadwerkelijke code voorbeelden uit de codebase.

## Base URL

```
Production: https://dklemailservice.onrender.com
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
| GET | `/api/contact` | Contact lijst | `contact:read` |
| GET | `/api/contact/:id` | Contact details | `contact:read` |
| GET | `/api/contact/status/:status` | Contact filter op status | `contact:read` |
| PUT | `/api/contact/:id` | Contact bijwerken | `contact:write` |
| DELETE | `/api/contact/:id` | Contact verwijderen | `contact:delete` |
| POST | `/api/contact/:id/antwoord` | Antwoord toevoegen | `contact:write` |
| GET | `/api/aanmelding` | Aanmelding lijst | `aanmelding:read` |
| GET | `/api/aanmelding/:id` | Aanmelding details | `aanmelding:read` |
| GET | `/api/aanmelding/rol/:rol` | Aanmelding filter op rol | `aanmelding:read` |
| PUT | `/api/aanmelding/:id` | Aanmelding bijwerken | `aanmelding:write` |
| DELETE | `/api/aanmelding/:id` | Aanmelding verwijderen | `aanmelding:delete` |
| POST | `/api/aanmelding/:id/antwoord` | Antwoord toevoegen | `aanmelding:write` |
| GET | `/api/mail` | Inkomende emails | API Key of `admin:access` |
| GET | `/api/mail/:id` | Email details | API Key of `admin:access` |
| PUT | `/api/mail/:id/processed` | Markeer als verwerkt | API Key of `admin:access` |
| DELETE | `/api/mail/:id` | Email verwijderen | API Key of `admin:access` |
| POST | `/api/mail/fetch` | Handmatig ophalen | API Key of `admin:access` |
| GET | `/api/mail/unprocessed` | Onverwerkte emails | API Key of `admin:access` |
| GET | `/api/mail/account/:type` | Emails per account type | API Key of `admin:access` |
| GET | `/api/users` | Gebruikers lijst | `user:read` |
| GET | `/api/users/:id` | Gebruiker details | `user:read` |
| POST | `/api/users` | Gebruiker aanmaken | `user:write` |
| PUT | `/api/users/:id` | Gebruiker bijwerken | `user:write` |
| PUT | `/api/users/:id/roles` | Roles toewijzen | Admin |
| DELETE | `/api/users/:id` | Gebruiker verwijderen | `user:delete` |
| GET | `/api/newsletter` | Nieuwsbrief lijst | `newsletter:read` |
| GET | `/api/newsletter/:id` | Nieuwsbrief details | `newsletter:read` |
| POST | `/api/newsletter` | Nieuwsbrief aanmaken | `newsletter:write` |
| PUT | `/api/newsletter/:id` | Nieuwsbrief bijwerken | `newsletter:write` |
| DELETE | `/api/newsletter/:id` | Nieuwsbrief verwijderen | `newsletter:delete` |
| POST | `/api/newsletter/:id/send` | Nieuwsbrief verzenden | `newsletter:send` |
| GET | `/api/rbac/permissions` | Permissions lijst | Admin |
| POST | `/api/rbac/permissions` | Permission aanmaken | Admin |
| GET | `/api/rbac/roles` | Roles lijst | Admin |
| POST | `/api/rbac/roles` | Role aanmaken | Admin |
| PUT | `/api/rbac/roles/:id/permissions` | Permissions toewijzen | Admin |
| POST | `/api/rbac/roles/:id/permissions/:permissionId` | Permission toevoegen | Admin |
| DELETE | `/api/rbac/roles/:id/permissions/:permissionId` | Permission verwijderen | Admin |
| POST | `/api/admin/mail/send` | Admin email verzenden | `admin_email:send` |
| GET | `/api/v1/notifications` | Notificaties lijst | Auth |
| POST | `/api/v1/notifications` | Notificatie aanmaken | Auth |
| GET | `/api/v1/notifications/:id` | Notificatie details | Auth |
| DELETE | `/api/v1/notifications/:id` | Notificatie verwijderen | Auth |
| POST | `/api/v1/notifications/reprocess-all` | Notificaties herverwerken | Auth |
| GET | `/api/albums` | Zichtbare albums lijst | Public |
| GET | `/api/albums/:id/photos` | Foto's van album | Public |
| GET | `/api/albums/admin` | Alle albums (admin) | `album:read` |
| GET | `/api/albums/:id` | Album details | `album:read` |
| POST | `/api/albums` | Album aanmaken | `album:write` |
| PUT | `/api/albums/:id` | Album bijwerken | `album:write` |
| PUT | `/api/albums/reorder` | Albums herschikken | `album:write` |
| DELETE | `/api/albums/:id` | Album verwijderen | `album:delete` |
| POST | `/api/albums/:id/photos` | Foto toevoegen aan album | `album:write` |
| PUT | `/api/albums/:id/photos/reorder` | Foto's in album herschikken | `album:write` |
| DELETE | `/api/albums/:id/photos/:photoId` | Foto uit album verwijderen | `album:delete` |
| GET | `/api/photos` | Zichtbare foto's lijst | Public |
| GET | `/api/photos/admin` | Alle foto's (admin) | `photo:read` |
| GET | `/api/photos/:id` | Foto details | `photo:read` |
| POST | `/api/photos` | Foto aanmaken | `photo:write` |
| PUT | `/api/photos/:id` | Foto bijwerken | `photo:write` |
| DELETE | `/api/photos/:id` | Foto verwijderen | `photo:delete` |
| GET | `/api/partners` | Partners lijst | Public |
| GET | `/api/partners/admin` | Alle partners (admin) | `partner:read` |
| GET | `/api/partners/:id` | Partner details | `partner:read` |
| POST | `/api/partners` | Partner aanmaken | `partner:write` |
| PUT | `/api/partners/:id` | Partner bijwerken | `partner:write` |
| DELETE | `/api/partners/:id` | Partner verwijderen | `partner:delete` |
| GET | `/api/radio-recordings` | Radio opnames lijst | Public |
| GET | `/api/radio-recordings/admin` | Alle radio opnames (admin) | `radio_recording:read` |
| GET | `/api/radio-recordings/:id` | Radio opname details | `radio_recording:read` |
| POST | `/api/radio-recordings` | Radio opname aanmaken | `radio_recording:write` |
| PUT | `/api/radio-recordings/:id` | Radio opname bijwerken | `radio_recording:write` |
| DELETE | `/api/radio-recordings/:id` | Radio opname verwijderen | `radio_recording:delete` |
| GET | `/api/videos` | Video's lijst | Public |
| GET | `/api/videos/admin` | Alle video's (admin) | `video:read` |
| GET | `/api/videos/:id` | Video details | `video:read` |
| POST | `/api/videos` | Video aanmaken | `video:write` |
| PUT | `/api/videos/:id` | Video bijwerken | `video:write` |
| DELETE | `/api/videos/:id` | Video verwijderen | `video:delete` |
| GET | `/api/sponsors` | Sponsors lijst | Public |
| GET | `/api/sponsors/admin` | Alle sponsors (admin) | `sponsor:read` |
| GET | `/api/sponsors/:id` | Sponsor details | `sponsor:read` |
| POST | `/api/sponsors` | Sponsor aanmaken | `sponsor:write` |
| PUT | `/api/sponsors/:id` | Sponsor bijwerken | `sponsor:write` |
| DELETE | `/api/sponsors/:id` | Sponsor verwijderen | `sponsor:delete` |
| GET | `/api/program-schedule` | Programma schema lijst | Public |
| GET | `/api/program-schedule/admin` | Alle programma schema's (admin) | `program_schedule:read` |
| GET | `/api/program-schedule/:id` | Programma schema details | `program_schedule:read` |
| POST | `/api/program-schedule` | Programma schema aanmaken | `program_schedule:write` |
| PUT | `/api/program-schedule/:id` | Programma schema bijwerken | `program_schedule:write` |
| DELETE | `/api/program-schedule/:id` | Programma schema verwijderen | `program_schedule:delete` |
| GET | `/api/social-embeds` | Social embeds lijst | Public |
| GET | `/api/social-embeds/admin` | Alle social embeds (admin) | `social_embed:read` |
| GET | `/api/social-embeds/:id` | Social embed details | `social_embed:read` |
| POST | `/api/social-embeds` | Social embed aanmaken | `social_embed:write` |
| PUT | `/api/social-embeds/:id` | Social embed bijwerken | `social_embed:write` |
| DELETE | `/api/social-embeds/:id` | Social embed verwijderen | `social_embed:delete` |
| GET | `/api/social-links` | Social links lijst | Public |
| GET | `/api/social-links/admin` | Alle social links (admin) | `social_link:read` |
| GET | `/api/social-links/:id` | Social link details | `social_link:read` |
| POST | `/api/social-links` | Social link aanmaken | `social_link:write` |
| PUT | `/api/social-links/:id` | Social link bijwerken | `social_link:write` |
| DELETE | `/api/social-links/:id` | Social link verwijderen | `social_link:delete` |
| GET | `/api/under-construction` | Under construction status | Public |
| PUT | `/api/under-construction` | Under construction bijwerken | `under_construction:write` |
| GET | `/api/images/upload` | Upload instellingen | Auth |
| POST | `/api/images/upload` | Afbeelding uploaden | Auth |
| GET | `/api/images/:id` | Afbeelding details | Auth |
| DELETE | `/api/images/:id` | Afbeelding verwijderen | Auth |

### Chat Endpoints

| Method | Endpoint | Beschrijving | Permission |
|--------|----------|--------------|------------|
| GET | `/api/chat/channels` | Gebruikers channels | Auth |
| GET | `/api/chat/channels/:id/participants` | Channel deelnemers | Auth |
| GET | `/api/chat/public-channels` | Publieke channels | Auth |
| POST | `/api/chat/direct` | Direct channel aanmaken | Auth |
| POST | `/api/chat/channels` | Channel aanmaken | Auth |
| POST | `/api/chat/channels/:id/join` | Channel joinen | Auth |
| POST | `/api/chat/channels/:id/leave` | Channel verlaten | Auth |
| GET | `/api/chat/users` | Gebruikers lijst | Auth |
| GET | `/api/chat/channels/:channel_id/messages` | Berichten ophalen | Auth |
| POST | `/api/chat/channels/:channel_id/messages` | Bericht verzenden | Auth |
| PUT | `/api/chat/messages/:id` | Bericht bewerken | Auth |
| DELETE | `/api/chat/messages/:id` | Bericht verwijderen | Auth |
| POST | `/api/chat/messages/:id/reactions` | Reactie toevoegen | Auth |
| DELETE | `/api/chat/messages/:id/reactions/:emoji` | Reactie verwijderen | Auth |
| PUT | `/api/chat/presence` | Presence bijwerken | Auth |
| GET | `/api/chat/online-users` | Online gebruikers | Auth |
| POST | `/api/chat/channels/:channel_id/typing/start` | Typing starten | Auth |
| POST | `/api/chat/channels/:channel_id/typing/stop` | Typing stoppen | Auth |
| GET | `/api/chat/channels/:channel_id/typing` | Typing gebruikers | Auth |
| POST | `/api/chat/channels/:id/read` | Als gelezen markeren | Auth |
| GET | `/api/chat/unread` | Ongelezen berichten | Auth |
| GET | `/api/chat/ws/:channel_id` | WebSocket verbinding | Auth |
| GET | `/api/chat/ws` | WebSocket verbinding | Auth |

### Whisky for Charity Endpoints

| Method | Endpoint | Beschrijving | Auth |
|--------|----------|--------------|------|
| POST | `/api/wfc/order-email` | WFC order emails verzenden | API Key |

### Mail Management Endpoints

| Method | Endpoint | Beschrijving | Permission |
|--------|----------|--------------|------------|
| GET | `/api/mail` | Inkomende emails | API Key of `admin:access` |
| GET | `/api/mail/:id` | Email details | API Key of `admin:access` |
| PUT | `/api/mail/:id/processed` | Markeer als verwerkt | API Key of `admin:access` |
| DELETE | `/api/mail/:id` | Email verwijderen | API Key of `admin:access` |
| POST | `/api/mail/fetch` | Handmatig ophalen | API Key of `admin:access` |
| GET | `/api/mail/unprocessed` | Onverwerkte emails | API Key of `admin:access` |
| GET | `/api/mail/account/:type` | Emails per account type | API Key of `admin:access` |

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
    	},
    	{
    		"path": "/api/aanmelding-email",
    		"method": "POST",
    		"description": "Send registration form email"
    	},
    	{
    		"path": "/api/auth/login",
    		"method": "POST",
    		"description": "User login"
    	},
    	{
    		"path": "/api/auth/logout",
    		"method": "POST",
    		"description": "User logout"
    	},
    	{
    		"path": "/api/auth/refresh",
    		"method": "POST",
    		"description": "Token refresh"
    	},
    	{
    		"path": "/api/auth/profile",
    		"method": "GET",
    		"description": "Get user profile (requires auth)"
    	},
    	{
    		"path": "/api/auth/reset-password",
    		"method": "POST",
    		"description": "Reset password (requires auth)"
    	},
    	{
    		"path": "/api/contact",
    		"method": "GET",
    		"description": "List contact forms (requires admin auth)"
    	},
    	{
    		"path": "/api/contact/:id",
    		"method": "GET",
    		"description": "Get contact form details (requires admin auth)"
    	},
    	{
    		"path": "/api/contact/:id",
    		"method": "PUT",
    		"description": "Update contact form (requires admin auth)"
    	},
    	{
    		"path": "/api/contact/:id",
    		"method": "DELETE",
    		"description": "Delete contact form (requires admin auth)"
    	},
    	{
    		"path": "/api/contact/:id/antwoord",
    		"method": "POST",
    		"description": "Add reply to contact form (requires admin auth)"
    	},
    	{
    		"path": "/api/contact/status/:status",
    		"method": "GET",
    		"description": "Filter contact forms by status (requires admin auth)"
    	},
    	{
    		"path": "/api/aanmelding",
    		"method": "GET",
    		"description": "List registrations (requires admin auth)"
    	},
    	{
    		"path": "/api/aanmelding/:id",
    		"method": "GET",
    		"description": "Get registration details (requires admin auth)"
    	},
    	{
    		"path": "/api/aanmelding/:id",
    		"method": "PUT",
    		"description": "Update registration (requires admin auth)"
    	},
    	{
    		"path": "/api/aanmelding/:id",
    		"method": "DELETE",
    		"description": "Delete registration (requires admin auth)"
    	},
    	{
    		"path": "/api/aanmelding/:id/antwoord",
    		"method": "POST",
    		"description": "Add reply to registration (requires admin auth)"
    	},
    	{
    		"path": "/api/aanmelding/rol/:rol",
    		"method": "GET",
    		"description": "Filter registrations by role (requires admin auth)"
    	},
    	{
    		"path": "/api/wfc/order-email",
    		"method": "POST",
    		"description": "Send Whisky for Charity order emails (requires API key)"
    	},
    	{
    		"path": "/metrics",
    		"method": "GET",
    		"description": "Prometheus metrics"
    	},
    	{
    		"path": "/api/metrics/email",
    		"method": "GET",
    		"description": "Email metrics (requires API key)"
    	},
    	{
    		"path": "/api/metrics/rate-limits",
    		"method": "GET",
    		"description": "Rate limit metrics (requires API key)"
    	},
    	{
    		"path": "/api/mail",
    		"method": "GET",
    		"description": "List incoming emails (requires admin auth)"
    	},
    	{
    		"path": "/api/mail/:id",
    		"method": "GET",
    		"description": "Get email details (requires admin auth)"
    	},
    	{
    		"path": "/api/mail/:id/processed",
    		"method": "PUT",
    		"description": "Mark email as processed (requires admin auth)"
    	},
    	{
    		"path": "/api/mail/:id",
    		"method": "DELETE",
    		"description": "Delete email (requires admin auth)"
    	},
    	{
    		"path": "/api/mail/fetch",
    		"method": "POST",
    		"description": "Manually fetch emails (requires admin auth)"
    	},
    	{
    		"path": "/api/mail/unprocessed",
    		"method": "GET",
    		"description": "List unprocessed emails (requires admin auth)"
    	},
    	{
    		"path": "/api/mail/account/:type",
    		"method": "GET",
    		"description": "List emails by account type (requires admin auth)"
    	},
    	{
    		"path": "/api/users",
    		"method": "GET",
    		"description": "List users (requires admin auth)"
    	},
    	{
    		"path": "/api/users/:id",
    		"method": "GET",
    		"description": "Get user details (requires admin auth)"
    	},
    	{
    		"path": "/api/users",
    		"method": "POST",
    		"description": "Create user (requires admin auth)"
    	},
    	{
    		"path": "/api/users/:id",
    		"method": "PUT",
    		"description": "Update user (requires admin auth)"
    	},
    	{
    		"path": "/api/users/:id/roles",
    		"method": "PUT",
    		"description": "Assign roles to user (requires admin auth)"
    	},
    	{
    		"path": "/api/users/:id",
    		"method": "DELETE",
    		"description": "Delete user (requires admin auth)"
    	},
    	{
    		"path": "/api/newsletter",
    		"method": "GET",
    		"description": "List newsletters (requires admin auth)"
    	},
    	{
    		"path": "/api/newsletter/:id",
    		"method": "GET",
    		"description": "Get newsletter details (requires admin auth)"
    	},
    	{
    		"path": "/api/newsletter",
    		"method": "POST",
    		"description": "Create newsletter (requires admin auth)"
    	},
    	{
    		"path": "/api/newsletter/:id",
    		"method": "PUT",
    		"description": "Update newsletter (requires admin auth)"
    	},
    	{
    		"path": "/api/newsletter/:id",
    		"method": "DELETE",
    		"description": "Delete newsletter (requires admin auth)"
    	},
    	{
    		"path": "/api/newsletter/:id/send",
    		"method": "POST",
    		"description": "Send newsletter (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/permissions",
    		"method": "GET",
    		"description": "List permissions (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/permissions",
    		"method": "POST",
    		"description": "Create permission (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/roles",
    		"method": "GET",
    		"description": "List roles (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/roles",
    		"method": "POST",
    		"description": "Create role (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/roles/:id/permissions",
    		"method": "PUT",
    		"description": "Update role permissions (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/roles/:id/permissions/:permissionId",
    		"method": "POST",
    		"description": "Add permission to role (requires admin auth)"
    	},
    	{
    		"path": "/api/rbac/roles/:id/permissions/:permissionId",
    		"method": "DELETE",
    		"description": "Remove permission from role (requires admin auth)"
    	},
    	{
    		"path": "/api/admin/mail/send",
    		"method": "POST",
    		"description": "Send admin email (requires admin auth)"
    	},
    	{
    		"path": "/api/v1/notifications",
    		"method": "GET",
    		"description": "List notifications (requires auth)"
    	},
    	{
    		"path": "/api/v1/notifications",
    		"method": "POST",
    		"description": "Create notification (requires auth)"
    	},
    	{
    		"path": "/api/v1/notifications/:id",
    		"method": "GET",
    		"description": "Get notification details (requires auth)"
    	},
    	{
    		"path": "/api/v1/notifications/:id",
    		"method": "DELETE",
    		"description": "Delete notification (requires auth)"
    	},
    	{
    		"path": "/api/v1/notifications/reprocess-all",
    		"method": "POST",
    		"description": "Reprocess all notifications (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels",
    		"method": "GET",
    		"description": "List user channels (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:id/participants",
    		"method": "GET",
    		"description": "List channel participants (requires auth)"
    	},
    	{
    		"path": "/api/chat/public-channels",
    		"method": "GET",
    		"description": "List public channels (requires auth)"
    	},
    	{
    		"path": "/api/chat/direct",
    		"method": "POST",
    		"description": "Create direct channel (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels",
    		"method": "POST",
    		"description": "Create channel (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:id/join",
    		"method": "POST",
    		"description": "Join channel (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:id/leave",
    		"method": "POST",
    		"description": "Leave channel (requires auth)"
    	},
    	{
    		"path": "/api/chat/users",
    		"method": "GET",
    		"description": "List users for chat (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:channel_id/messages",
    		"method": "GET",
    		"description": "Get channel messages (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:channel_id/messages",
    		"method": "POST",
    		"description": "Send message (requires auth)"
    	},
    	{
    		"path": "/api/chat/messages/:id",
    		"method": "PUT",
    		"description": "Edit message (requires auth)"
    	},
    	{
    		"path": "/api/chat/messages/:id",
    		"method": "DELETE",
    		"description": "Delete message (requires auth)"
    	},
    	{
    		"path": "/api/chat/messages/:id/reactions",
    		"method": "POST",
    		"description": "Add reaction (requires auth)"
    	},
    	{
    		"path": "/api/chat/messages/:id/reactions/:emoji",
    		"method": "DELETE",
    		"description": "Remove reaction (requires auth)"
    	},
    	{
    		"path": "/api/chat/presence",
    		"method": "PUT",
    		"description": "Update presence (requires auth)"
    	},
    	{
    		"path": "/api/chat/online-users",
    		"method": "GET",
    		"description": "List online users (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:channel_id/typing/start",
    		"method": "POST",
    		"description": "Start typing (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:channel_id/typing/stop",
    		"method": "POST",
    		"description": "Stop typing (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:channel_id/typing",
    		"method": "GET",
    		"description": "Get typing users (requires auth)"
    	},
    	{
    		"path": "/api/chat/channels/:id/read",
    		"method": "POST",
    		"description": "Mark as read (requires auth)"
    	},
    	{
    		"path": "/api/chat/unread",
    		"method": "GET",
    		"description": "Get unread count (requires auth)"
    	},
    	{
    		"path": "/api/chat/ws/:channel_id",
    		"method": "GET",
    		"description": "WebSocket connection (requires auth)"
    	},
    	{
    		"path": "/api/chat/ws",
    		"method": "GET",
    		"description": "WebSocket connection (requires auth)"
    	}
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
curl -X POST https://dklemailservice.onrender.com/api/contact-email \
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
const response = await fetch('https://dklemailservice.onrender.com/api/contact-email', {
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

### User Management

#### GET /api/users

Haalt een lijst van gebruikers op (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 50)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "naam": "John Doe",
            "email": "john@example.com",
            "rol": "admin",
            "is_actief": true,
            "newsletter_subscribed": true,
            "created_at": "2024-03-20T15:04:05Z",
            "updated_at": "2024-03-20T15:04:05Z"
        }
    ],
    "total": 1,
    "limit": 50,
    "offset": 0
}
```

**Implementatie:** [`handlers/user_handler.go:27`](../../handlers/user_handler.go:27)

#### POST /api/users

Maakt een nieuwe gebruiker aan (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "email": "john@example.com",
    "naam": "John Doe",
    "rol": "admin",
    "password": "securepassword",
    "is_actief": true,
    "newsletter_subscribed": false
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "naam": "John Doe",
    "email": "john@example.com",
    "rol": "admin",
    "is_actief": true,
    "newsletter_subscribed": false,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T15:04:05Z"
}
```

#### PUT /api/users/:id/roles

Werkt de roles van een gebruiker bij (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "role_ids": ["role-uuid-1", "role-uuid-2"]
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Roles toegewezen aan user",
    "assigned_roles": 2,
    "total_requested": 2
}
```

**cURL Voorbeeld:**
```bash
curl -X POST https://dklemailservice.onrender.com/api/aanmelding-email \
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

#### GET /api/mail/account/:type

Haalt emails op per account type met paginering.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 10)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
{
    "emails": [
        {
            "id": "770e8400-e29b-41d4-a716-446655440000",
            "message_id": "<message123@example.com>",
            "sender": "sender@example.com",
            "to": "info@dekoninklijkeloop.nl",
            "subject": "Vraag over het evenement",
            "html": "Hallo, ik heb een vraag...",
            "content_type": "text/plain",
            "received_at": "2024-04-01T09:30:00Z",
            "uid": "AAABBCCC123",
            "account_type": "info",
            "read": false,
            "processed_at": null,
            "created_at": "2024-04-01T09:35:00Z",
            "updated_at": "2024-04-01T09:35:00Z"
        }
    ],
    "totalCount": 1
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

### Newsletter Management

#### GET /api/newsletter

Haalt een lijst van nieuwsbrieven op (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 10)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "subject": "Nieuwsbrief onderwerp",
    "content": "Nieuwsbrief inhoud...",
    "sent_at": null,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T15:04:05Z"
}
```

**Implementatie:** [`handlers/newsletter_handler.go`](../../handlers/newsletter_handler.go:1)

#### POST /api/newsletter

Maakt een nieuwe nieuwsbrief aan (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "subject": "Nieuwsbrief onderwerp",
    "content": "Nieuwsbrief inhoud..."
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "subject": "Nieuwsbrief onderwerp",
    "content": "Nieuwsbrief inhoud...",
    "sent_at": null,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T15:04:05Z"
}
```

#### POST /api/newsletter/:id/send

Verzendt een nieuwsbrief naar subscribers (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Nieuwsbrief wordt verzonden naar subscribers"
}
```

**Implementatie:** [`handlers/newsletter_handler.go`](../../handlers/newsletter_handler.go:1)

### Notification Management

#### GET /api/v1/notifications

Haalt een lijst van notificaties op.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `type` (optioneel): Filter op type
- `priority` (optioneel): Filter op prioriteit

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "type": "contact",
        "priority": "medium",
        "title": "Nieuwe contact aanvraag",
        "message": "Nieuwe contact aanvraag ontvangen",
        "sent": false,
        "created_at": "2024-03-20T15:04:05Z"
    }
]
```

**Implementatie:** [`handlers/notification_handler.go`](../../handlers/notification_handler.go:1)

#### POST /api/v1/notifications

Maakt een nieuwe notificatie aan.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "type": "contact",
    "priority": "medium",
    "title": "Nieuwe contact aanvraag",
    "message": "Nieuwe contact aanvraag ontvangen"
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "contact",
    "priority": "medium",
    "title": "Nieuwe contact aanvraag",
    "message": "Nieuwe contact aanvraag ontvangen",
    "sent": false,
    "created_at": "2024-03-20T15:04:05Z"
}
```

### RBAC Management

#### GET /api/rbac/permissions

Haalt een lijst van permissions op (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "resource": "contact",
        "action": "read",
        "description": "Kan contact formulieren lezen"
    }
]
```

**Implementatie:** [`handlers/permission_handler.go`](../../handlers/permission_handler.go:1)

#### GET /api/rbac/roles

Haalt een lijst van roles met hun permissions op (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "admin",
        "description": "Administrator rol",
        "permissions": [
            {
                "id": "550e8400-e29b-41d4-a716-446655440001",
                "resource": "contact",
                "action": "read"
            }
        ]
    }
]
```

#### PUT /api/rbac/roles/:id/permissions

Werkt permissions voor een role bij (admin only).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "permission_ids": ["perm-uuid-1", "perm-uuid-2"]
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Role permissions bijgewerkt",
    "added_count": 1,
    "removed_count": 0,
    "total_requested": 2
}
```

**Implementatie:** [`handlers/permission_handler.go`](../../handlers/permission_handler.go:1)

### Chat API

#### GET /api/chat/channels

Haalt channels op waar de gebruiker lid van is.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 50)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Algemeen",
        "type": "public",
        "created_by": "user-uuid",
        "created_at": "2024-03-20T15:04:05Z"
    }
]
```

**Implementatie:** [`handlers/chat_handler.go`](../../handlers/chat_handler.go:1)

#### POST /api/chat/direct

Maakt een direct channel aan tussen twee gebruikers.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "user_id": "target-user-uuid"
}
```

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Chat between User A and User B",
    "type": "direct",
    "created_by": "current-user-uuid",
    "created_at": "2024-03-20T15:04:05Z"
}
```

#### GET /api/chat/channels/:channel_id/messages

Haalt berichten op voor een channel.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 50)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "channel_id": "channel-uuid",
        "user_id": "user-uuid",
        "content": "Hallo allemaal!",
        "created_at": "2024-03-20T15:04:05Z",
        "edited_at": null
    }
]
```

#### POST /api/chat/channels/:channel_id/messages

Verzendt een nieuw bericht naar een channel.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "content": "Hallo allemaal!"
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "channel_id": "channel-uuid",
    "user_id": "user-uuid",
    "content": "Hallo allemaal!",
    "created_at": "2024-03-20T15:04:05Z",
    "edited_at": null
}
```

#### GET /api/chat/ws/:channel_id

WebSocket verbinding voor real-time chat in een specifiek channel.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**WebSocket Protocol:**
- Client stuurt: `{"type": "message", "content": "Hallo!"}`
- Server stuurt: `{"type": "message", "data": {...}}`

**Implementatie:** [`handlers/chat_handler.go`](../../handlers/chat_handler.go:1)

### Whisky for Charity API

#### POST /api/wfc/order-email

Verzendt order emails voor Whisky for Charity bestellingen.

**Headers:**
```http
X-API-Key: <wfc-api-key>
Content-Type: application/json
```

**Request Body:**
```json
{
    "order_id": "WFC-12345",
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_address": "Straat 123",
    "customer_city": "Amsterdam",
    "customer_postal": "1234AB",
    "customer_country": "NL",
    "total_amount": 150.00,
    "items": [
        {
            "name": "Whisky Glas Set",
            "quantity": 2,
            "price": 75.00
        }
    ]
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "customer_email_sent": true,
    "admin_email_sent": true,
    "order_id": "WFC-12345"
}
```

**Implementatie:** [`handlers/wfc_order_handler.go`](../../handlers/wfc_order_handler.go:1)

### Admin Email API

#### POST /api/admin/mail/send

Verzendt emails namens admin gebruikers.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "to": "recipient@example.com",
    "subject": "Onderwerp",
    "body": "Email inhoud...",
    "template_name": "optional_template",
    "template_variables": {
        "name": "John"
    }
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Email succesvol verzonden."
}
```

**Implementatie:** [`handlers/admin_mail_handler.go`](../../handlers/admin_mail_handler.go:1)

### Album Management

#### GET /api/albums

Haalt een lijst van zichtbare albums op voor publiek gebruik.

**Query Parameters:**
- `include_covers` (optioneel): Inclusief cover foto informatie (`true`/`false`)

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Album Titel",
        "description": "Album beschrijving",
        "cover_photo_id": "photo-uuid",
        "visible": true,
        "order_number": 1,
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z"
    }
]
```

**Implementatie:** [`handlers/album_handler.go:77`](../../handlers/album_handler.go:77)

#### GET /api/albums/:id/photos

Haalt alle zichtbare foto's van een specifiek album op.

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "url": "https://example.com/photo.jpg",
        "alt_text": "Foto beschrijving",
        "visible": true,
        "thumbnail_url": "https://example.com/thumb.jpg",
        "title": "Foto titel",
        "description": "Foto beschrijving",
        "year": 2024,
        "cloudinary_folder": "albums/2024",
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z",
        "album_id": "album-uuid",
        "order_number": 1
    }
]
```

**Implementatie:** [`handlers/album_handler.go:350`](../../handlers/album_handler.go:350)

#### GET /api/albums/admin

Haalt een lijst van alle albums op voor admin beheer.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 10)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "title": "Album Titel",
            "description": "Album beschrijving",
            "cover_photo_id": "photo-uuid",
            "visible": true,
            "order_number": 1,
            "created_at": "2024-03-20T15:04:05Z",
            "updated_at": "2024-03-20T15:04:05Z"
        }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
}
```

**Implementatie:** [`handlers/album_handler.go:117`](../../handlers/album_handler.go:117)

#### GET /api/albums/:id

Haalt details van een specifiek album op.

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
        "title": "Album Titel",
        "description": "Album beschrijving",
        "cover_photo_id": "photo-uuid",
        "visible": true,
        "order_number": 1,
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z"
    }
}
```

**Implementatie:** [`handlers/album_handler.go:159`](../../handlers/album_handler.go:159)

#### POST /api/albums

Maakt een nieuw album aan.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "title": "Nieuw Album",
    "description": "Album beschrijving",
    "cover_photo_id": "photo-uuid",
    "visible": true,
    "order_number": 1
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Nieuw Album",
    "description": "Album beschrijving",
    "cover_photo_id": "photo-uuid",
    "visible": true,
    "order_number": 1,
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T15:04:05Z"
}
```

**Implementatie:** [`handlers/album_handler.go:198`](../../handlers/album_handler.go:198)

#### PUT /api/albums/:id

Werkt een bestaand album bij.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "title": "Bijgewerkte Titel",
    "description": "Bijgewerkte beschrijving",
    "visible": false
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Album bijgewerkt",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Bijgewerkte Titel",
        "description": "Bijgewerkte beschrijving",
        "visible": false,
        "updated_at": "2024-03-20T16:00:00Z"
    }
}
```

**Implementatie:** [`handlers/album_handler.go:239`](../../handlers/album_handler.go:239)

#### PUT /api/albums/reorder

Herschikt de volgorde van meerdere albums (bulk operatie).

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "album_order": [
        {"id": "album-1", "order_number": 1},
        {"id": "album-2", "order_number": 2},
        {"id": "album-3", "order_number": 3}
    ]
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Albums reordered successfully"
}
```

**Implementatie:** [`handlers/album_handler.go:657`](../../handlers/album_handler.go:657)

#### DELETE /api/albums/:id

Verwijdert een album.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Album deleted successfully"
}
```

**Implementatie:** [`handlers/album_handler.go:302`](../../handlers/album_handler.go:302)

#### POST /api/albums/:id/photos

Voegt een foto toe aan een album.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "photo_id": "photo-uuid",
    "order_number": 1
}
```

**Response (201 Created):**
```json
{
    "id": "album-photo-uuid",
    "album_id": "album-uuid",
    "photo_id": "photo-uuid",
    "order_number": 1,
    "created_at": "2024-03-20T15:04:05Z"
}
```

**Implementatie:** [`handlers/album_handler.go:400`](../../handlers/album_handler.go:400)

#### PUT /api/albums/:id/photos/reorder

Herschikt de volgorde van foto's in een album.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "photo_order": [
        {"photo_id": "photo-1", "order_number": 1},
        {"photo_id": "photo-2", "order_number": 2}
    ]
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Photos reordered successfully"
}
```

**Implementatie:** [`handlers/album_handler.go:586`](../../handlers/album_handler.go:586)

#### DELETE /api/albums/:id/photos/:photoId

Verwijdert een foto uit een album.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Photo removed from album successfully"
}
```

**Implementatie:** [`handlers/album_handler.go:500`](../../handlers/album_handler.go:500)

### Photo Management

#### GET /api/photos

Haalt een lijst van zichtbare foto's op.

**Query Parameters:**
- `year` (optioneel): Filter op jaar
- `title` (optioneel): Filter op titel (gedeeltelijke match, case-insensitive)
- `description` (optioneel): Filter op beschrijving (gedeeltelijke match, case-insensitive)
- `cloudinary_folder` (optioneel): Filter op Cloudinary folder (exact match)

**Response (200 OK):**
```json
[
    {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "url": "https://example.com/photo.jpg",
        "alt_text": "Foto beschrijving",
        "visible": true,
        "thumbnail_url": "https://example.com/thumb.jpg",
        "title": "Foto titel",
        "description": "Foto beschrijving",
        "year": 2024,
        "cloudinary_folder": "albums/2024",
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z"
    }
]
```

**Implementatie:** [`handlers/photo_handler.go:69`](../../handlers/photo_handler.go:69)

#### GET /api/photos/admin

Haalt een lijst van alle foto's op voor admin beheer.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `limit` (optioneel): Maximum aantal resultaten (default: 10)
- `offset` (optioneel): Aantal resultaten om over te slaan (default: 0)

**Response (200 OK):**
```json
{
    "success": true,
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "url": "https://example.com/photo.jpg",
            "alt_text": "Foto beschrijving",
            "visible": true,
            "thumbnail_url": "https://example.com/thumb.jpg",
            "title": "Foto titel",
            "description": "Foto beschrijving",
            "year": 2024,
            "cloudinary_folder": "albums/2024",
            "created_at": "2024-03-20T15:04:05Z",
            "updated_at": "2024-03-20T15:04:05Z"
        }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
}
```

**Implementatie:** [`handlers/photo_handler.go:126`](../../handlers/photo_handler.go:126)

#### GET /api/photos/:id

Haalt details van een specifieke foto op.

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
        "url": "https://example.com/photo.jpg",
        "alt_text": "Foto beschrijving",
        "visible": true,
        "thumbnail_url": "https://example.com/thumb.jpg",
        "title": "Foto titel",
        "description": "Foto beschrijving",
        "year": 2024,
        "cloudinary_folder": "albums/2024",
        "created_at": "2024-03-20T15:04:05Z",
        "updated_at": "2024-03-20T15:04:05Z"
    }
}
```

**Implementatie:** [`handlers/photo_handler.go:168`](../../handlers/photo_handler.go:168)

#### POST /api/photos

Maakt een nieuwe foto aan.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "url": "https://example.com/photo.jpg",
    "alt_text": "Foto beschrijving",
    "visible": true,
    "thumbnail_url": "https://example.com/thumb.jpg",
    "title": "Foto titel",
    "description": "Foto beschrijving",
    "year": 2024,
    "cloudinary_folder": "albums/2024"
}
```

**Response (201 Created):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "url": "https://example.com/photo.jpg",
    "alt_text": "Foto beschrijving",
    "visible": true,
    "thumbnail_url": "https://example.com/thumb.jpg",
    "title": "Foto titel",
    "description": "Foto beschrijving",
    "year": 2024,
    "cloudinary_folder": "albums/2024",
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T15:04:05Z"
}
```

**Implementatie:** [`handlers/photo_handler.go:207`](../../handlers/photo_handler.go:207)

#### PUT /api/photos/:id

Werkt een bestaande foto bij.

**Headers:**
```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "title": "Bijgewerkte titel",
    "description": "Bijgewerkte beschrijving",
    "visible": false
}
```

**Response (200 OK):**
```json
{
    "success": true,
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Bijgewerkte titel",
        "description": "Bijgewerkte beschrijving",
        "visible": false,
        "updated_at": "2024-03-20T16:00:00Z"
    }
}
```

**Implementatie:** [`handlers/photo_handler.go:248`](../../handlers/photo_handler.go:248)

#### DELETE /api/photos/:id

Verwijdert een foto.

**Headers:**
```http
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
    "success": true,
    "message": "Photo deleted successfully"
}
```

**Implementatie:** [`handlers/photo_handler.go:314`](../../handlers/photo_handler.go:314)

## Error Codes
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
        "https://dklemailservice.onrender.com/api/contact-email",
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
    const response = await fetch('https://dklemailservice.onrender.com/api/contact-email', {
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
- [RBAC Frontend Guide](../../docs/RBAC_FRONTEND.md) - RBAC permissies voor frontend
- [README](../../README.md) - Project documentatie