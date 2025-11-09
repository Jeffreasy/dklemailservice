# API Documentation

Complete API reference for the DKL Email Service.

## 🚀 Quick Start

**New to the API?** Start here: [API Quick Reference](./QUICK_REFERENCE.md) - Overzicht van alle endpoints in één document.

## Base URL

```
Production: https://api.dklemailservice.com
Development: http://localhost:8080
```

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

### Authentication Endpoints

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and receive JWT token
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout and invalidate refresh token
- `GET /api/auth/me` - Get current user information

See [Authentication API](./AUTHENTICATION.md) for detailed documentation.

## API Categories

### 🔐 [Authentication & Authorization](./AUTHENTICATION.md)
- User registration and login
- JWT token management
- Session handling
- Password management

### 🛡️ [Permissions & RBAC](./PERMISSIONS.md)
- Role management
- Permission assignments
- Access control
- User role management

### 🎉 [Events & Registrations](./EVENTS.md)
- Event management (CRUD)
- Event registrations
- Geofencing & location tracking
- Event statistics

### 🏃 [Steps & Gamification](./STEPS_GAMIFICATION.md)
- Steps tracking and updates
- Achievements & badges
- Leaderboards
- Real-time WebSocket updates

### 🎨 [Content Management System (CMS)](./CMS.md)
- Videos (YouTube integration)
- Partners & Sponsors
- Photo Albums & Photos
- Radio Recordings
- Program Schedule
- Social Links & Embeds
- Title Sections
- Image uploads (Cloudinary)

### 🔔 [Notifications](./NOTIFICATIONS.md)
- User notifications
- Real-time updates (SSE)
- Notification preferences
- Broadcast messaging

### 📡 [WebSocket APIs](./WEBSOCKET.md)
- Notulen (Meeting Notes) real-time updates
- Steps application real-time updates
- Connection management
- Message types

## Response Format

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "message": "Optional success message"
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    // Optional error details
  }
}
```

## HTTP Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `204 No Content` - Successful request with no response body
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Request conflicts with current state
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

## Rate Limiting

API requests are rate-limited:
- **Anonymous users**: 100 requests per hour
- **Authenticated users**: 1000 requests per hour
- **Admin users**: 5000 requests per hour

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

## Pagination

List endpoints support pagination using query parameters:

```
GET /api/resource?page=1&limit=20&sort=created_at&order=desc
```

Response includes pagination metadata:

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

## Filtering & Sorting

Most list endpoints support filtering and sorting:

```
GET /api/resource?filter[status]=active&filter[type]=email&sort=created_at&order=desc
```

## WebSocket Endpoints

Real-time updates are available via WebSocket:

- `ws://localhost:8080/ws/steps` - Steps updates
- `ws://localhost:8080/api/ws/notulen` - Notulen updates
- `ws://localhost:8080/api/chat/ws/:channel_id` - Chat channel updates

See [WebSocket Documentation](./WEBSOCKET.md) for details.

## API Versioning

The API uses path-based versioning:
- Current version: `/api/v1/...` (default)
- Legacy support: `/api/...` (maps to v1)

## CORS

CORS is enabled for configured frontend origins. Contact your administrator to add additional origins.

## Security

- All endpoints use HTTPS in production
- JWT tokens expire after 24 hours
- Refresh tokens expire after 7 days
- Rate limiting protects against abuse
- SQL injection protection via parameterized queries
- XSS protection via input sanitization

## Testing

Use the provided Postman collection or cURL examples:

```bash
# Health check
curl http://localhost:8080/api/health

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Protected endpoint
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

For detailed endpoint documentation, see the specific API category pages listed above.