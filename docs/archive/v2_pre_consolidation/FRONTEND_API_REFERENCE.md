# Frontend API Reference

This document provides comprehensive API documentation for the DKL Email Service frontend integration.

## Base URL
```
https://your-domain.com/api
```

## Authentication
Most admin endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Content Types
- Request: `application/json`
- Response: `application/json`

---

## 📸 Albums API

### Public Endpoints

#### Get Visible Albums
```http
GET /api/albums
```

**Response:**
```json
[
  {
    "id": "uuid",
    "title": "Album Title",
    "description": "Album description",
    "cover_photo_id": "uuid",
    "visible": true,
    "order_number": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

### Admin Endpoints

#### List All Albums
```http
GET /api/albums/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Album by ID
```http
GET /api/albums/{id}
Authorization: Bearer <token>
```

#### Create Album
```http
POST /api/albums
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "New Album",
  "description": "Album description",
  "cover_photo_id": "uuid",
  "visible": true,
  "order_number": 1
}
```

#### Update Album
```http
PUT /api/albums/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Album",
  "description": "Updated description",
  "visible": false
}
```

#### Delete Album
```http
DELETE /api/albums/{id}
Authorization: Bearer <token>
```

---

## 🎥 Videos API

### Public Endpoints

#### Get Visible Videos
```http
GET /api/videos
```

**Response:**
```json
[
  {
    "id": "uuid",
    "video_id": "streamable_id",
    "url": "https://streamable.com/e/...",
    "title": "Video Title",
    "description": "Video description",
    "thumbnail_url": null,
    "visible": true,
    "order_number": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

### Admin Endpoints

#### List All Videos
```http
GET /api/videos/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Video by ID
```http
GET /api/videos/{id}
Authorization: Bearer <token>
```

#### Create Video
```http
POST /api/videos
Authorization: Bearer <token>
Content-Type: application/json

{
  "video_id": "q9ngqu",
  "url": "https://streamable.com/e/q9ngqu",
  "title": "New Video",
  "description": "Video description",
  "thumbnail_url": "https://...",
  "visible": true,
  "order_number": 1
}
```

#### Update Video
```http
PUT /api/videos/{id}
Authorization: Bearer <token>
```

#### Delete Video
```http
DELETE /api/videos/{id}
Authorization: Bearer <token>
```

---

## 🤝 Sponsors API

### Public Endpoints

#### Get Visible Sponsors
```http
GET /api/sponsors
```

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "Sponsor Name",
    "description": "Sponsor description",
    "logo_url": "https://...",
    "website_url": "https://...",
    "order_number": 1,
    "is_active": true,
    "visible": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

### Admin Endpoints

#### List All Sponsors
```http
GET /api/sponsors/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Sponsor by ID
```http
GET /api/sponsors/{id}
Authorization: Bearer <token>
```

#### Create Sponsor
```http
POST /api/sponsors
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Sponsor",
  "description": "Sponsor description",
  "logo_url": "https://...",
  "website_url": "https://...",
  "order_number": 1,
  "is_active": true,
  "visible": true
}
```

#### Update Sponsor
```http
PUT /api/sponsors/{id}
Authorization: Bearer <token>
```

#### Delete Sponsor
```http
DELETE /api/sponsors/{id}
Authorization: Bearer <token>
```

---

## 📅 Program Schedule API

### Public Endpoints

#### Get Visible Program Schedule
```http
GET /api/program-schedule
```

**Response:**
```json
[
  {
    "id": "uuid",
    "time": "10:00u",
    "event_description": "Event description",
    "category": "Start",
    "icon_name": "start",
    "order_number": 1,
    "visible": true,
    "latitude": 52.123456,
    "longitude": 5.123456,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

### Admin Endpoints

#### List All Program Schedule
```http
GET /api/program-schedule/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Program Schedule by ID
```http
GET /api/program-schedule/{id}
Authorization: Bearer <token>
```

#### Create Program Schedule Item
```http
POST /api/program-schedule
Authorization: Bearer <token>
Content-Type: application/json

{
  "time": "10:00u",
  "event_description": "New event",
  "category": "Start",
  "icon_name": "start",
  "order_number": 1,
  "visible": true,
  "latitude": 52.123456,
  "longitude": 5.123456
}
```

#### Update Program Schedule Item
```http
PUT /api/program-schedule/{id}
Authorization: Bearer <token>
```

#### Delete Program Schedule Item
```http
DELETE /api/program-schedule/{id}
Authorization: Bearer <token>
```

---

## 📱 Social Embeds API

### Public Endpoints

#### Get Visible Social Embeds
```http
GET /api/social-embeds
```

**Response:**
```json
[
  {
    "id": "uuid",
    "platform": "instagram",
    "embed_code": "<iframe>...</iframe>",
    "order_number": 1,
    "visible": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

### Admin Endpoints

#### List All Social Embeds
```http
GET /api/social-embeds/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Social Embed by ID
```http
GET /api/social-embeds/{id}
Authorization: Bearer <token>
```

#### Create Social Embed
```http
POST /api/social-embeds
Authorization: Bearer <token>
Content-Type: application/json

{
  "platform": "instagram",
  "embed_code": "<iframe>...</iframe>",
  "order_number": 1,
  "visible": true
}
```

#### Update Social Embed
```http
PUT /api/social-embeds/{id}
Authorization: Bearer <token>
```

#### Delete Social Embed
```http
DELETE /api/social-embeds/{id}
Authorization: Bearer <token>
```

---

## 🔗 Social Links API

### Public Endpoints

#### Get Visible Social Links
```http
GET /api/social-links
```

**Response:**
```json
[
  {
    "id": "uuid",
    "platform": "instagram",
    "url": "https://instagram.com/...",
    "bg_color_class": null,
    "icon_color_class": null,
    "order_number": 1,
    "visible": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

### Admin Endpoints

#### List All Social Links
```http
GET /api/social-links/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Social Link by ID
```http
GET /api/social-links/{id}
Authorization: Bearer <token>
```

#### Create Social Link
```http
POST /api/social-links
Authorization: Bearer <token>
Content-Type: application/json

{
  "platform": "instagram",
  "url": "https://instagram.com/...",
  "bg_color_class": "bg-pink-500",
  "icon_color_class": "text-white",
  "order_number": 1,
  "visible": true
}
```

#### Update Social Link
```http
PUT /api/social-links/{id}
Authorization: Bearer <token>
```

#### Delete Social Link
```http
DELETE /api/social-links/{id}
Authorization: Bearer <token>
```

---

## 🚧 Under Construction API

### Public Endpoints

#### Get Active Under Construction
```http
GET /api/under-construction/active
```

**Response:**
```json
{
  "id": 1,
  "is_active": true,
  "title": "Website Under Maintenance",
  "message": "We're updating our website...",
  "footer_text": "Thank you for your patience",
  "logo_url": "https://...",
  "expected_date": "2025-12-31T23:59:59Z",
  "social_links": "[{\"url\": \"...\", \"platform\": \"...\"}]",
  "progress_percentage": 85,
  "contact_email": "info@example.com",
  "newsletter_enabled": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

### Admin Endpoints

#### List All Under Construction Records
```http
GET /api/under-construction/admin?limit=10&offset=0
Authorization: Bearer <token>
```

#### Get Under Construction by ID
```http
GET /api/under-construction/{id}
Authorization: Bearer <token>
```

#### Create Under Construction Record
```http
POST /api/under-construction
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_active": true,
  "title": "Maintenance Mode",
  "message": "Website is under maintenance",
  "footer_text": "Thank you",
  "logo_url": "https://...",
  "expected_date": "2025-12-31T23:59:59Z",
  "social_links": "[]",
  "progress_percentage": 50,
  "contact_email": "info@example.com",
  "newsletter_enabled": false
}
```

#### Update Under Construction Record
```http
PUT /api/under-construction/{id}
Authorization: Bearer <token>
```

#### Delete Under Construction Record
```http
DELETE /api/under-construction/{id}
Authorization: Bearer <token>
```

---

## 🔐 Authentication

### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "password"
}
```

**Response:**
```json
{
  "token": "jwt-token-here",
  "refresh_token": "refresh-token-here",
  "user": {
    "id": "uuid",
    "email": "admin@example.com",
    "naam": "Admin User"
  }
}
```

### Refresh Token
```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh-token-here"
}
```

---

## 📊 Common Response Patterns

### Success Response
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

### Paginated Response
```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "pages": 10
  }
}
```

---

## 🔧 Error Codes

- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

---

## 📝 Notes for Frontend Developers

1. **CORS**: The API supports CORS for configured origins
2. **Rate Limiting**: Authentication endpoints have rate limiting
3. **File Uploads**: Use `/api/images/upload` for image uploads
4. **WebSocket**: Real-time features available via WebSocket
5. **Pagination**: Use `limit` and `offset` query parameters
6. **Sorting**: Results are ordered by `order_number` and `created_at`

## 🧪 Testing

Use the following test credentials for development:
- Email: `admin@example.com`
- Password: Check your `.env` file or database

## 📞 Support

For API issues, check the server logs or contact the development team.