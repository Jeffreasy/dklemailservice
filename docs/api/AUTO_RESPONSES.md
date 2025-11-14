# Auto Response API

**Base Path:** `/api/mail/autoresponse`  
**Authentication:** Required  
**Permissions:** `mail:manage` (admin only)

---

## Overview

The Auto Response API allows managing automated email replies for specified email accounts. This is useful for:
- Out-of-office messages
- Vacation responders
- Automated acknowledgments
- Service maintenance notifications

### Database Table (V33)

Created in migration [`V33__create_auto_responses_table.sql`](../../database/migrations/V33__create_auto_responses_table.sql)

```sql
CREATE TABLE auto_responses (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    subject VARCHAR(255),
    message TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## Endpoints

### 1. List All Auto Responses

**GET** `/api/mail/autoresponse`

Retrieves all configured auto responses.

#### Request

```http
GET /api/mail/autoresponse HTTP/1.1
Host: api.dekoninklijkeloop.nl
Authorization: Bearer <jwt_token>
```

#### Response

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "email": "info@dekoninklijkeloop.nl",
      "is_active": true,
      "subject": "We hebben je bericht ontvangen",
      "message": "Bedankt voor je bericht. We nemen binnen 24 uur contact met je op.",
      "start_date": "2025-11-01T00:00:00Z",
      "end_date": "2025-11-30T23:59:59Z",
      "created_at": "2025-11-01T10:00:00Z",
      "updated_at": "2025-11-01T10:00:00Z"
    },
    {
      "id": 2,
      "email": "inschrijving@dekoninklijkeloop.nl",
      "is_active": false,
      "subject": "Inschrijving bevestiging",
      "message": "Je inschrijving is ontvangen en wordt verwerkt.",
      "start_date": null,
      "end_date": null,
      "created_at": "2025-10-15T14:30:00Z",
      "updated_at": "2025-10-20T09:15:00Z"
    }
  ],
  "count": 2
}
```

#### Status Codes

- `200 OK` - Successfully retrieved auto responses
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - User doesn't have `mail:manage` permission
- `500 Internal Server Error` - Database error

---

### 2. Get Auto Response by ID

**GET** `/api/mail/autoresponse/:id`

Retrieves a specific auto response configuration.

#### Request

```http
GET /api/mail/autoresponse/1 HTTP/1.1
Host: api.dekoninklijkeloop.nl
Authorization: Bearer <jwt_token>
```

#### Response

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "info@dekoninklijkeloop.nl",
    "is_active": true,
    "subject": "We hebben je bericht ontvangen",
    "message": "Bedankt voor je bericht. We nemen binnen 24 uur contact met je op.",
    "start_date": "2025-11-01T00:00:00Z",
    "end_date": "2025-11-30T23:59:59Z",
    "created_at": "2025-11-01T10:00:00Z",
    "updated_at": "2025-11-01T10:00:00Z"
  }
}
```

#### Status Codes

- `200 OK` - Successfully retrieved auto response
- `404 Not Found` - Auto response with specified ID doesn't exist
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - User doesn't have `mail:manage` permission
- `500 Internal Server Error` - Database error

---

### 3. Create Auto Response

**POST** `/api/mail/autoresponse`

Creates a new auto response configuration.

#### Request

```http
POST /api/mail/autoresponse HTTP/1.1
Host: api.dekoninklijkeloop.nl
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "email": "info@dekoninklijkeloop.nl",
  "is_active": true,
  "subject": "Out of Office",
  "message": "We zijn op vakantie van 1-15 december. We reageren daarna op je bericht.",
  "start_date": "2025-12-01T00:00:00Z",
  "end_date": "2025-12-15T23:59:59Z"
}
```

#### Request Body Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | ✅ Yes | Email address this auto response applies to |
| `is_active` | boolean | No | Whether auto response is active (default: false) |
| `subject` | string | No | Subject line for auto response email |
| `message` | text | No | Body content of auto response |
| `start_date` | timestamp | No | When to start sending auto responses |
| `end_date` | timestamp | No | When to stop sending auto responses |

#### Response

```json
{
  "success": true,
  "message": "Auto response created successfully",
  "data": {
    "id": 3,
    "email": "info@dekoninklijkeloop.nl",
    "is_active": true,
    "subject": "Out of Office",
    "message": "We zijn op vakantie van 1-15 december. We reageren daarna op je bericht.",
    "start_date": "2025-12-01T00:00:00Z",
    "end_date": "2025-12-15T23:59:59Z",
    "created_at": "2025-11-10T18:06:00Z",
    "updated_at": "2025-11-10T18:06:00Z"
  }
}
```

#### Status Codes

- `201 Created` - Auto response created successfully
- `400 Bad Request` - Invalid request body (missing email, invalid format)
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - User doesn't have `mail:manage` permission
- `409 Conflict` - Auto response already exists for this email
- `500 Internal Server Error` - Database error

---

### 4. Update Auto Response

**PUT** `/api/mail/autoresponse/:id`

Updates an existing auto response configuration.

#### Request

```http
PUT /api/mail/autoresponse/1 HTTP/1.1
Host: api.dekoninklijkeloop.nl
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "is_active": false,
  "subject": "We zijn weer bereikbaar",
  "message": "Bedankt voor je geduld. We zijn weer bereikbaar."
}
```

#### Request Body

All fields are optional. Only included fields will be updated.

| Field | Type | Description |
|-------|------|-------------|
| `email` | string | Email address (usually not changed) |
| `is_active` | boolean | Toggle auto response on/off |
| `subject` | string | Updated subject line |
| `message` | text | Updated message content |
| `start_date` | timestamp | Updated start date |
| `end_date` | timestamp | Updated end date |

#### Response

```json
{
  "success": true,
  "message": "Auto response updated successfully",
  "data": {
    "id": 1,
    "email": "info@dekoninklijkeloop.nl",
    "is_active": false,
    "subject": "We zijn weer bereikbaar",
    "message": "Bedankt voor je geduld. We zijn weer bereikbaar.",
    "start_date": "2025-11-01T00:00:00Z",
    "end_date": "2025-11-30T23:59:59Z",
    "created_at": "2025-11-01T10:00:00Z",
    "updated_at": "2025-11-10T18:07:00Z"
  }
}
```

#### Status Codes

- `200 OK` - Auto response updated successfully
- `400 Bad Request` - Invalid request body
- `404 Not Found` - Auto response with specified ID doesn't exist
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - User doesn't have `mail:manage` permission
- `500 Internal Server Error` - Database error

---

### 5. Delete Auto Response

**DELETE** `/api/mail/autoresponse/:id`

Permanently deletes an auto response configuration.

#### Request

```http
DELETE /api/mail/autoresponse/1 HTTP/1.1
Host: api.dekoninklijkeloop.nl
Authorization: Bearer <jwt_token>
```

#### Response

```json
{
  "success": true,
  "message": "Auto response deleted successfully"
}
```

#### Status Codes

- `200 OK` - Auto response deleted successfully
- `404 Not Found` - Auto response with specified ID doesn't exist
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - User doesn't have `mail:manage` permission
- `500 Internal Server Error` - Database error

---

## Use Cases

### 1. Vacation Responder

```javascript
// Create vacation auto response
const response = await fetch('/api/mail/autoresponse', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'info@dekoninklijkeloop.nl',
    is_active: true,
    subject: 'Out of Office - Vakantie',
    message: `
      Beste correspondent,
      
      Wij zijn met vakantie van 1 t/m 15 december 2025.
      We reageren op je bericht zodra we terug zijn.
      
      Voor urgente zaken kun je contact opnemen met:
      - Tel: +31 6 12345678
      - Email: urgent@dekoninklijkeloop.nl
      
      Met vriendelijke groet,
      Team De Koninklijke Loop
    `,
    start_date: '2025-12-01T00:00:00Z',
    end_date: '2025-12-15T23:59:59Z'
  })
});
```

### 2. Event Registration Confirmation

```javascript
// Set up automatic registration confirmation
const response = await fetch('/api/mail/autoresponse', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'inschrijving@dekoninklijkeloop.nl',
    is_active: true,
    subject: 'Bevestiging inschrijving DKL 2026',
    message: `
      Beste deelnemer,
      
      Hartelijk dank voor je inschrijving voor De Koninklijke Loop 2026!
      
      We hebben je aanmelding ontvangen en deze wordt nu verwerkt.
      Je ontvangt binnen 2 werkdagen een definitieve bevestiging per email.
      
      Heb je vragen? Stuur een email naar info@dekoninklijkeloop.nl
      
      Met sportieve groet,
      Organisatie DKL 2026
    `
  })
});
```

### 3. Maintenance Notice

```javascript
// Temporary maintenance message
const response = await fetch('/api/mail/autoresponse', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'support@dekoninklijkeloop.nl',
    is_active: true,
    subject: 'Onderhoud Support Systeem',
    message: `
      Beste gebruiker,
      
      Ons support systeem ondergaat momenteel onderhoud.
      We zijn tussen 22:00 en 02:00 niet bereikbaar.
      
      Je bericht wordt morgenochtend als eerste behandeld.
      
      Met vriendelijke groet,
      Support Team
    `,
    start_date: '2025-11-15T22:00:00Z',
    end_date: '2025-11-16T02:00:00Z'
  })
});
```

---

## Frontend Integration

### React Hook Example

```typescript
// hooks/useAutoResponses.ts
import { useState, useEffect } from 'react';

interface AutoResponse {
  id: number;
  email: string;
  is_active: boolean;
  subject: string;
  message: string;
  start_date: string | null;
  end_date: string | null;
  created_at: string;
  updated_at: string;
}

export function useAutoResponses() {
  const [autoResponses, setAutoResponses] = useState<AutoResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchAutoResponses = async () => {
    try {
      const response = await fetch('/api/mail/autoresponse', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) throw new Error('Failed to fetch auto responses');
      
      const data = await response.json();
      setAutoResponses(data.data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createAutoResponse = async (data: Omit<AutoResponse, 'id' | 'created_at' | 'updated_at'>) => {
    const response = await fetch('/api/mail/autoresponse', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });
    
    if (!response.ok) throw new Error('Failed to create auto response');
    
    await fetchAutoResponses(); // Refresh list
    return response.json();
  };

  const updateAutoResponse = async (id: number, updates: Partial<AutoResponse>) => {
    const response = await fetch(`/api/mail/autoresponse/${id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updates)
    });
    
    if (!response.ok) throw new Error('Failed to update auto response');
    
    await fetchAutoResponses(); // Refresh list
    return response.json();
  };

  const deleteAutoResponse = async (id: number) => {
    const response = await fetch(`/api/mail/autoresponse/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) throw new Error('Failed to delete auto response');
    
    await fetchAutoResponses(); // Refresh list
  };

  useEffect(() => {
    fetchAutoResponses();
  }, []);

  return {
    autoResponses,
    loading,
    error,
    createAutoResponse,
    updateAutoResponse,
    deleteAutoResponse,
    refresh: fetchAutoResponses
  };
}
```

### Vue Composable Example

```typescript
// composables/useAutoResponses.ts
import { ref, onMounted } from 'vue';

export function useAutoResponses() {
  const autoResponses = ref([]);
  const loading = ref(true);
  const error = ref(null);

  const fetchAutoResponses = async () => {
    try {
      loading.value = true;
      const response = await fetch('/api/mail/autoresponse', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) throw new Error('Failed to fetch');
      
      const data = await response.json();
      autoResponses.value = data.data;
    } catch (err) {
      error.value = err.message;
    } finally {
      loading.value = false;
    }
  };

  const toggleActive = async (id, isActive) => {
    try {
      await fetch(`/api/mail/autoresponse/${id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ is_active: isActive })
      });
      
      await fetchAutoResponses();
    } catch (err) {
      error.value = err.message;
    }
  };

  onMounted(() => {
    fetchAutoResponses();
  });

  return {
    autoResponses,
    loading,
    error,
    fetchAutoResponses,
    toggleActive
  };
}
```

---

## Implementation Details

### Handler

**File:** [`handlers/auto_response_handler.go`](../../handlers/auto_response_handler.go)

```go
type AutoResponseHandler struct {
    repo repository.AutoResponseRepository
}

// GetAll lists all auto responses
// GET /api/mail/autoresponse
func (h *AutoResponseHandler) GetAll(c *gin.Context) { ... }

// GetByID retrieves specific auto response
// GET /api/mail/autoresponse/:id
func (h *AutoResponseHandler) GetByID(c *gin.Context) { ... }

// Create adds new auto response
// POST /api/mail/autoresponse
func (h *AutoResponseHandler) Create(c *gin.Context) { ... }

// Update modifies existing auto response
// PUT /api/mail/autoresponse/:id
func (h *AutoResponseHandler) Update(c *gin.Context) { ... }

// Delete removes auto response
// DELETE /api/mail/autoresponse/:id
func (h *AutoResponseHandler) Delete(c *gin.Context) { ... }
```

### Repository

**File:** [`repository/auto_response_repository.go`](../../repository/auto_response_repository.go)

Provides CRUD operations for auto_responses table.

### Model

**File:** [`models/auto_response.go`](../../models/auto_response.go)

```go
type AutoResponse struct {
    ID         int        `json:"id" gorm:"primaryKey"`
    Email      string     `json:"email" gorm:"type:varchar(255);not null"`
    IsActive   bool       `json:"is_active" gorm:"default:false"`
    Subject    string     `json:"subject" gorm:"type:varchar(255)"`
    Message    string     `json:"message" gorm:"type:text"`
    StartDate  *time.Time `json:"start_date"`
    EndDate    *time.Time `json:"end_date"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}
```

---

## Related Documentation

- [Migration V33](../../database/migrations/V33__create_auto_responses_table.sql) - Database schema
- [CORRECTIONS_APPLIED.md](../CORRECTIONS_APPLIED.md) - Implementation details
- [Mail Management API](./QUICK_REFERENCE.md#mail-management) - Related endpoints

---

**Created:** V33 Migration (2025-11-10)  
**Last Updated:** 2025-11-10  
**Status:** ✅ Production Ready