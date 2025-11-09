# CMS API Documentation

Content Management System API endpoints voor het beheren van media, partners, sponsors en andere content.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

## Authentication

Alle admin endpoints vereisen JWT authenticatie en specifieke permissies:

```
Authorization: Bearer <your_jwt_token>
```

---

## Videos

Beheer video content van YouTube en andere bronnen.

### List Visible Videos (Public)

**Endpoint:** `GET /api/videos`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "video_id": "dQw4w9WgXcQ",
    "url": "https://youtube.com/watch?v=dQw4w9WgXcQ",
    "title": "Event Highlights 2024",
    "description": "Hoogtepunten van het evenement",
    "thumbnail_url": "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
    "visible": true,
    "order_number": 1,
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### List All Videos (Admin)

**Endpoint:** `GET /api/videos/admin`

**Authentication:** Vereist - `video:read` permissie

**Query Parameters:**
- `limit` (optional): Aantal resultaten (default: 10, max: 100)
- `offset` (optional): Aantal over te slaan (default: 0)

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "video_id": "dQw4w9WgXcQ",
    "url": "https://youtube.com/watch?v=dQw4w9WgXcQ",
    "title": "Event Highlights 2024",
    "description": "Hoogtepunten van het evenement",
    "thumbnail_url": "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
    "visible": false,
    "order_number": 1,
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### Get Video

**Endpoint:** `GET /api/videos/:id`

**Authentication:** Vereist - `video:read` permissie

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "video_id": "dQw4w9WgXcQ",
  "url": "https://youtube.com/watch?v=dQw4w9WgXcQ",
  "title": "Event Highlights 2024",
  "description": "Hoogtepunten van het evenement",
  "thumbnail_url": "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
  "visible": true,
  "order_number": 1,
  "created_at": "2025-01-08T10:00:00Z",
  "updated_at": "2025-01-08T10:00:00Z"
}
```

---

### Create Video

**Endpoint:** `POST /api/videos`

**Authentication:** Vereist - `video:write` permissie

**Request Body:**

```json
{
  "video_id": "dQw4w9WgXcQ",
  "url": "https://youtube.com/watch?v=dQw4w9WgXcQ",
  "title": "Event Highlights 2024",
  "description": "Hoogtepunten van het evenement",
  "thumbnail_url": "https://i.ytimg.com/vi/dQw4w9WgXcQ/maxresdefault.jpg",
  "visible": true,
  "order_number": 1
}
```

**Response:** `201 Created`

---

### Update Video

**Endpoint:** `PUT /api/videos/:id`

**Authentication:** Vereist - `video:write` permissie

**Request Body:**

```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "visible": false,
  "order_number": 2
}
```

**Response:** `200 OK`

---

### Delete Video

**Endpoint:** `DELETE /api/videos/:id`

**Authentication:** Vereist - `video:delete` permissie

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Video deleted successfully"
}
```

---

## Partners

Beheer partnerorganisaties.

### List Visible Partners (Public)

**Endpoint:** `GET /api/partners`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "Partner Name",
    "description": "Partner description",
    "logo": "https://res.cloudinary.com/...",
    "website": "https://partner.com",
    "tier": "gold",
    "since": "2020-01-01T00:00:00Z",
    "visible": true,
    "order_number": 1,
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### List All Partners (Admin)

**Endpoint:** `GET /api/partners/admin`

**Authentication:** Vereist - `partner:read` permissie

**Query Parameters:**
- `limit` (optional): Aantal resultaten (default: 10, max: 100)
- `offset` (optional): Aantal over te slaan (default: 0)

---

### Get Partner

**Endpoint:** `GET /api/partners/:id`

**Authentication:** Vereist - `partner:read` permissie

---

### Create Partner

**Endpoint:** `POST /api/partners`

**Authentication:** Vereist - `partner:write` permissie

**Request Body:**

```json
{
  "name": "New Partner",
  "description": "Partner description",
  "logo": "https://res.cloudinary.com/...",
  "website": "https://partner.com",
  "tier": "gold",
  "since": "2020-01-01T00:00:00Z",
  "visible": true,
  "order_number": 1
}
```

**Response:** `201 Created`

---

### Update Partner

**Endpoint:** `PUT /api/partners/:id`

**Authentication:** Vereist - `partner:write` permissie

---

### Delete Partner

**Endpoint:** `DELETE /api/partners/:id`

**Authentication:** Vereist - `partner:delete` permissie

---

## Sponsors

Beheer sponsors met logo upload functionaliteit.

### List Visible Sponsors (Public)

**Endpoint:** `GET /api/sponsors`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "Sponsor Name",
    "description": "Sponsor description",
    "logo_url": "https://res.cloudinary.com/...",
    "website_url": "https://sponsor.com",
    "order_number": 1,
    "is_active": true,
    "visible": true,
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### List All Sponsors (Admin)

**Endpoint:** `GET /api/sponsors/admin`

**Authentication:** Vereist - `sponsor:read` permissie

**Query Parameters:**
- `limit` (optional): Aantal resultaten (default: 10, max: 100)
- `offset` (optional): Aantal over te slaan (default: 0)

---

### Get Sponsor

**Endpoint:** `GET /api/sponsors/:id`

**Authentication:** Vereist - `sponsor:read` permissie

---

### Create Sponsor

**Endpoint:** `POST /api/sponsors`

**Authentication:** Vereist - `sponsor:write` permissie

**Content-Type:** `application/json` of `multipart/form-data`

**JSON Request Body:**

```json
{
  "name": "New Sponsor",
  "description": "Sponsor description",
  "logo_url": "https://res.cloudinary.com/...",
  "website_url": "https://sponsor.com",
  "order_number": 1,
  "is_active": true,
  "visible": true
}
```

**Multipart Form-Data:**

```
name: "New Sponsor"
description: "Sponsor description"
website_url: "https://sponsor.com"
order_number: 1
is_active: true
visible: true
logo: <file> (JPEG, PNG, GIF, WebP)
```

**Response:** `201 Created`

---

### Update Sponsor

**Endpoint:** `PUT /api/sponsors/:id`

**Authentication:** Vereist - `sponsor:write` permissie

---

### Delete Sponsor

**Endpoint:** `DELETE /api/sponsors/:id`

**Authentication:** Vereist - `sponsor:delete` permissie

---

## Radio Recordings

Beheer radio opnames en podcasts.

### List Visible Radio Recordings (Public)

**Endpoint:** `GET /api/radio-recordings`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "title": "Episode 1",
    "description": "First episode description",
    "date": "2025-01-08T00:00:00Z",
    "audio_url": "https://cdn.example.com/audio.mp3",
    "thumbnail_url": "https://res.cloudinary.com/...",
    "visible": true,
    "order_number": 1,
    "created_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### List All Radio Recordings (Admin)

**Endpoint:** `GET /api/radio-recordings/admin`

**Authentication:** Vereist - `radio_recording:read` permissie

**Query Parameters:**
- `limit` (optional): Aantal resultaten (default: 10, max: 100)
- `offset` (optional): Aantal over te slaan (default: 0)

---

### Get Radio Recording

**Endpoint:** `GET /api/radio-recordings/:id`

**Authentication:** Vereist - `radio_recording:read` permissie

---

### Create Radio Recording

**Endpoint:** `POST /api/radio-recordings`

**Authentication:** Vereist - `radio_recording:write` permissie

**Request Body:**

```json
{
  "title": "Episode 1",
  "description": "First episode description",
  "date": "2025-01-08T00:00:00Z",
  "audio_url": "https://cdn.example.com/audio.mp3",
  "thumbnail_url": "https://res.cloudinary.com/...",
  "visible": true,
  "order_number": 1
}
```

**Response:** `201 Created`

---

### Update Radio Recording

**Endpoint:** `PUT /api/radio-recordings/:id`

**Authentication:** Vereist - `radio_recording:write` permissie

---

### Delete Radio Recording

**Endpoint:** `DELETE /api/radio-recordings/:id`

**Authentication:** Vereist - `radio_recording:delete` permissie

---

## Photos

Beheer individuele foto's (los van albums).

### List Visible Photos (Public)

**Endpoint:** `GET /api/photos`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "afbeelding_url": "https://res.cloudinary.com/...",
    "cloudinary_public_id": "dkl_images/photo123",
    "title": "Photo Title",
    "description": "Photo description",
    "visible": true,
    "order_number": 1
  }
]
```

---

### Create Photo

**Endpoint:** `POST /api/photos`

**Authentication:** Vereist - `photo:write` permissie

**Request Body:**

```json
{
  "afbeelding_url": "https://res.cloudinary.com/...",
  "cloudinary_public_id": "dkl_images/photo123",
  "title": "Photo Title",
  "description": "Photo description",
  "visible": true,
  "order_number": 1
}
```

---

### Update Photo

**Endpoint:** `PUT /api/photos/:id`

**Authentication:** Vereist - `photo:write` permissie

---

### Delete Photo

**Endpoint:** `DELETE /api/photos/:id`

**Authentication:** Vereist - `photo:delete` permissie

---

## Albums

Beheer foto albums met ondersteuning voor meerdere foto's.

### List Visible Albums (Public)

**Endpoint:** `GET /api/albums`

**Authentication:** Niet vereist

---

### Get Album with Photos (Public)

**Endpoint:** `GET /api/albums/:id`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "titel": "Summer 2024",
  "beschrijving": "Foto's van de zomer editie",
  "cover_foto_url": "https://res.cloudinary.com/...",
  "zichtbaar": true,
  "order_number": 1,
  "fotos": [
    {
      "id": "uuid",
      "afbeelding_url": "https://res.cloudinary.com/...",
      "title": "Start Line",
      "order_number": 1
    }
  ]
}
```

---

### Create Album

**Endpoint:** `POST /api/albums`

**Authentication:** Vereist - `album:write` permissie

**Request Body:**

```json
{
  "titel": "Summer 2024",
  "beschrijving": "Foto's van de zomer editie",
  "cover_foto_url": "https://res.cloudinary.com/...",
  "zichtbaar": true,
  "order_number": 1
}
```

---

### Update Album

**Endpoint:** `PUT /api/albums/:id`

**Authentication:** Vereist - `album:write` permissie

---

### Delete Album

**Endpoint:** `DELETE /api/albums/:id`

**Authentication:** Vereist - `album:delete` permissie

---

### Add Photo to Album

**Endpoint:** `POST /api/albums/:id/photos`

**Authentication:** Vereist - `album:write` permissie

**Request Body:**

```json
{
  "photo_id": "uuid",
  "order_number": 1
}
```

**Response:** `201 Created`

---

### Reorder Album Photos

**Endpoint:** `PUT /api/albums/:id/photos/reorder`

**Authentication:** Vereist - `album:write` permissie

**Request Body:**

```json
{
  "photo_ids": ["uuid1", "uuid2", "uuid3"]
}
```

---

### Reorder Albums

**Endpoint:** `PUT /api/albums/reorder`

**Authentication:** Vereist - `album:write` permissie

**Request Body:**

```json
{
  "album_ids": ["uuid1", "uuid2", "uuid3"]
}
```

---

## Program Schedule

Beheer programma items en tijdschema's.

### List Visible Program Schedule (Public)

**Endpoint:** `GET /api/program-schedule`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "title": "Opening Ceremony",
    "description": "Event opening",
    "start_time": "2025-06-15T09:00:00Z",
    "end_time": "2025-06-15T09:30:00Z",
    "location": "Start Area",
    "visible": true,
    "order_number": 1
  }
]
```

---

### Create Program Schedule Item

**Endpoint:** `POST /api/program-schedule`

**Authentication:** Vereist - `program_schedule:write` permissie

**Request Body:**

```json
{
  "title": "Opening Ceremony",
  "description": "Event opening",
  "start_time": "2025-06-15T09:00:00Z",
  "end_time": "2025-06-15T09:30:00Z",
  "location": "Start Area",
  "visible": true,
  "order_number": 1
}
```

---

### Update Program Schedule Item

**Endpoint:** `PUT /api/program-schedule/:id`

**Authentication:** Vereist - `program_schedule:write` permissie

---

### Delete Program Schedule Item

**Endpoint:** `DELETE /api/program-schedule/:id`

**Authentication:** Vereist - `program_schedule:delete` permissie

---

## Social Links

Beheer social media links.

### List Visible Social Links (Public)

**Endpoint:** `GET /api/social-links`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "platform": "facebook",
    "url": "https://facebook.com/dkl",
    "icon_url": "https://res.cloudinary.com/...",
    "visible": true,
    "order_number": 1
  }
]
```

---

### Create Social Link

**Endpoint:** `POST /api/social-links`

**Authentication:** Vereist - `social_link:write` permissie

**Request Body:**

```json
{
  "platform": "facebook",
  "url": "https://facebook.com/dkl",
  "icon_url": "https://res.cloudinary.com/...",
  "visible": true,
  "order_number": 1
}
```

---

## Social Embeds

Beheer embedded social media content.

### List Visible Social Embeds (Public)

**Endpoint:** `GET /api/social-embeds`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "platform": "instagram",
    "embed_code": "<blockquote class='instagram-media'>...</blockquote>",
    "url": "https://instagram.com/p/...",
    "visible": true,
    "order_number": 1
  }
]
```

---

### Create Social Embed

**Endpoint:** `POST /api/social-embeds`

**Authentication:** Vereist - `social_embed:write` permissie

---

## Title Sections

Beheer titel secties voor de website.

### List Visible Title Sections (Public)

**Endpoint:** `GET /api/title-sections`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "section_key": "home_hero",
    "title": "De Koninklijke Loop",
    "subtitle": "Hardlopen voor KWF",
    "background_image": "https://res.cloudinary.com/...",
    "visible": true,
    "order_number": 1
  }
]
```

---

### Create Title Section

**Endpoint:** `POST /api/title-sections`

**Authentication:** Vereist - `title_section:write` permissie

**Request Body:**

```json
{
  "section_key": "home_hero",
  "title": "De Koninklijke Loop",
  "subtitle": "Hardlopen voor KWF",
  "background_image": "https://res.cloudinary.com/...",
  "visible": true,
  "order_number": 1
}
```

---

## Under Construction

Beheer onderhoudsmeldingen.

### Get Under Construction Status

**Endpoint:** `GET /api/under-construction`

**Authentication:** Niet vereist

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "is_active": false,
  "title": "Website in Onderhoud",
  "message": "We zijn bezig met updates",
  "estimated_end": "2025-01-08T12:00:00Z",
  "show_countdown": true
}
```

---

### Create/Update Under Construction

**Endpoint:** `POST /api/under-construction`

**Authentication:** Vereist - `under_construction:write` permissie

**Request Body:**

```json
{
  "is_active": true,
  "title": "Website in Onderhoud",
  "message": "We zijn bezig met updates",
  "estimated_end": "2025-01-08T12:00:00Z",
  "show_countdown": true
}
```

---

## Image Upload

Upload afbeeldingen naar Cloudinary.

### Upload Single Image

**Endpoint:** `POST /api/upload/image`

**Authentication:** Vereist - `image:write` permissie

**Content-Type:** `multipart/form-data`

**Form Data:**
- `file`: Image file (JPEG, PNG, GIF, WebP)
- `folder` (optional): Cloudinary folder naam

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "url": "https://res.cloudinary.com/...",
    "public_id": "dkl_images/abc123",
    "width": 1920,
    "height": 1080,
    "format": "jpg",
    "resource_type": "image"
  }
}
```

---

### Upload Batch Images

**Endpoint:** `POST /api/upload/images/batch`

**Authentication:** Vereist - `image:write` permissie

**Content-Type:** `multipart/form-data`

**Form Data:**
- `files`: Multiple image files
- `folder` (optional): Cloudinary folder naam

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "uploaded": [
      {
        "url": "https://res.cloudinary.com/...",
        "public_id": "dkl_images/abc123"
      }
    ],
    "failed": []
  }
}
```

---

## Error Codes

| Code | HTTP Status | Beschrijving |
|------|-------------|--------------|
| `INVALID_REQUEST` | 400 | Ongeldige request body |
| `UNAUTHORIZED` | 401 | Authenticatie vereist |
| `FORBIDDEN` | 403 | Onvoldoende permissies |
| `NOT_FOUND` | 404 | Resource niet gevonden |
| `VALIDATION_ERROR` | 422 | Validatie fout |
| `SERVER_ERROR` | 500 | Server fout |

---

## Best Practices

### Ordering Content

Gebruik `order_number` om de volgorde van content items te bepalen:

```json
{
  "order_number": 1
}
```

Lagere nummers verschijnen eerst in lijsten.

### Visibility Control

Gebruik `visible` of `zichtbaar` veld om content te verbergen zonder te verwijderen:

```json
{
  "visible": false
}
```

### Image Optimization

Bij het uploaden van images naar Cloudinary:
- Maximale bestandsgrootte: 10MB
- Ondersteunde formaten: JPEG, PNG, GIF, WebP
- Automatische optimalisatie wordt toegepast
- Thumbnail generatie voor performance

---

Voor meer informatie:
- [API Overview](./README.md)
- [Authentication](./AUTHENTICATION.md)
- [Frontend Integration](../guides/FRONTEND_INTEGRATION.md)