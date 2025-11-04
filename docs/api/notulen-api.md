# Notulen API Documentation

## Overview

The Notulen API provides comprehensive endpoints for managing meeting minutes (notulen) in the DKL Email Service system. This API handles the complete lifecycle of meeting minutes with support for complex data structures including attendees, agenda items, decisions, and action items.

## Authentication

All protected endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

Public endpoints (like `/api/notulen/public`) do not require authentication.

## Data Structures

### Core Types

#### AgendaItem
```typescript
interface AgendaItem {
  title: string;      // Agenda item title (required)
  details?: string;   // Optional details/description
}
```

#### Besluit (Decision)
```typescript
interface Besluit {
  besluit: string;           // Decision text (required)
  teams?: {                  // Optional team assignments
    [teamName: string]: string[]
  };
}
```

#### Actiepunt (Action Item)
```typescript
interface Actiepunt {
  actie: string;                    // Action description (required)
  verantwoordelijke: string | string[]; // Responsible person(s) (required)
}
```

### Notulen Structure

#### Notulen (Meeting Minutes)
```typescript
interface Notulen {
  id: string;                    // UUID (auto-generated)
  titel: string;                 // Meeting title (required, 3-255 chars)
  vergadering_datum: string;     // Meeting date (ISO format, required)
  locatie?: string;              // Optional location
  voorzitter?: string;           // Optional chairperson
  notulist?: string;             // Optional note taker
  aanwezigen: string[];          // Array of attendees (PostgreSQL TEXT[])
  afwezigen: string[];           // Array of absentees (PostgreSQL TEXT[])
  agenda_items: {                // Agenda items (JSONB wrapper)
    items: AgendaItem[]
  };
  besluiten: {                   // Decisions (JSONB wrapper)
    besluiten: Besluit[]
  };
  actiepunten: {                 // Action items (JSONB wrapper)
    acties: Actiepunt[]
  };
  notities?: string;             // Optional notes
  status: 'draft' | 'finalized' | 'archived'; // Default: 'draft'
  versie: number;                // Version number (auto-incremented)
  created_by: string;            // UUID of creator (required)
  created_at: string;            // ISO timestamp (auto-generated)
  updated_at: string;            // ISO timestamp (auto-updated)
  finalized_at?: string;         // ISO timestamp (set when finalized)
  finalized_by?: string;         // UUID of finalizer (set when finalized)
}
```

#### NotulenVersie (Version History)
```typescript
interface NotulenVersie {
  id: string;                    // UUID (auto-generated)
  notulen_id: string;            // Reference to parent notulen
  versie: number;                // Version number
  titel: string;                 // Meeting title at this version
  vergadering_datum: string;     // Meeting date at this version
  locatie?: string;              // Location at this version
  voorzitter?: string;           // Chairperson at this version
  notulist?: string;             // Note taker at this version
  aanwezigen: string[];          // Attendees at this version
  afwezigen: string[];           // Absentees at this version
  agenda_items: {                // Agenda items at this version
    items: AgendaItem[]
  };
  besluiten: {                   // Decisions at this version
    besluiten: Besluit[]
  };
  actiepunten: {                 // Action items at this version
    acties: Actiepunt[]
  };
  notities?: string;             // Notes at this version
  status: string;                // Status at this version
  gewijzigd_door: string;        // UUID of person who made this version
  gewijzigd_op: string;          // ISO timestamp of version creation
  wijziging_reden?: string;      // Optional reason for change
}
```

### Request/Response Types

#### NotulenCreateRequest
```typescript
interface NotulenCreateRequest {
  titel: string;              // Required, 3-255 characters
  vergadering_datum: string;  // Required, format: "YYYY-MM-DD"
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  aanwezigen?: string[];
  afwezigen?: string[];
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
}
```

#### NotulenUpdateRequest
```typescript
interface NotulenUpdateRequest {
  titel?: string;             // Optional, 3-255 characters if provided
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  aanwezigen?: string[];
  afwezigen?: string[];
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
}
```

#### NotulenFinalizeRequest
```typescript
interface NotulenFinalizeRequest {
  wijziging_reden?: string; // Optional reason for finalization
}
```

#### NotulenListResponse
```typescript
interface NotulenListResponse {
  notulen: Notulen[];
  total: number;
  limit: number;
  offset: number;
}
```

#### NotulenVersionsResponse
```typescript
interface NotulenVersionsResponse {
  versions: NotulenVersie[];
}
```

#### ErrorResponse
```typescript
interface ErrorResponse {
  error: string; // Error message
}
```

#### SuccessResponse
```typescript
interface SuccessResponse {
  success: boolean;
  message: string;
}
```

## API Endpoints

### Public Endpoints

#### GET /api/notulen/public
Get finalized meeting minutes (public access, no authentication required)

**Query Parameters:**
- `date_from` (optional): Start date filter (YYYY-MM-DD)
- `date_to` (optional): End date filter (YYYY-MM-DD)
- `limit` (optional): Maximum number of results (default: 50, max: 50)
- `offset` (optional): Pagination offset (default: 0)

**Response:**
```typescript
interface NotulenListResponse {
  notulen: Notulen[];
  total: number;
  limit: number;
  offset: number;
}
```

**Example:**
```bash
GET /api/notulen/public?date_from=2025-01-01&limit=10
```

### Protected Endpoints

#### GET /api/notulen
List all meeting minutes (authenticated users only)

**Query Parameters:**
- `date_from` (optional): Start date filter (YYYY-MM-DD)
- `date_to` (optional): End date filter (YYYY-MM-DD)
- `status` (optional): Filter by status ('draft', 'finalized', 'archived')
- `limit` (optional): Maximum number of results
- `offset` (optional): Pagination offset

**Permissions Required:** `notulen:read`

**Response:** Same as public endpoint but includes all statuses

#### GET /api/notulen/:id
Get a specific meeting minutes by ID

**Path Parameters:**
- `id`: UUID of the notulen

**Permissions Required:** `notulen:read`

**Response:**
```typescript
Notulen // Single notulen object
```

#### POST /api/notulen
Create new meeting minutes

**Request Body:**
```typescript
interface NotulenCreateRequest {
  titel: string;              // Required, 3-255 characters
  vergadering_datum: string;  // Required, ISO date string
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  aanwezigen?: string[];
  afwezigen?: string[];
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
}
```

**Permissions Required:** `notulen:write`

**Response:** Created `Notulen` object

#### PUT /api/notulen/:id
Update existing meeting minutes

**Path Parameters:**
- `id`: UUID of the notulen

**Request Body:** Same as create request but all fields optional

**Permissions Required:** `notulen:write`

**Response:** Updated `Notulen` object

#### POST /api/notulen/:id/finalize
Finalize a meeting minutes document

**Path Parameters:**
- `id`: UUID of the notulen

**Request Body (optional):**
```typescript
interface NotulenFinalizeRequest {
  wijziging_reden?: string; // Optional change reason
}
```

**Permissions Required:** `notulen:write`

**Response:**
```typescript
{
  success: true;
  message: "Notulen finalized successfully";
}
```

#### POST /api/notulen/:id/archive
Archive a meeting minutes document

**Path Parameters:**
- `id`: UUID of the notulen

**Permissions Required:** `notulen:write`

**Response:**
```typescript
{
  success: true;
  message: "Notulen archived successfully";
}
```

#### DELETE /api/notulen/:id
Delete a meeting minutes document (soft delete)

**Path Parameters:**
- `id`: UUID of the notulen

**Permissions Required:** `notulen:delete`

**Response:**
```typescript
{
  success: true;
  message: "Notulen deleted successfully";
}
```

### Search and Versioning

#### GET /api/notulen/search
Search meeting minutes with full-text search

**Query Parameters:**
- `q`: Search query (required)
- `date_from` (optional): Start date filter
- `date_to` (optional): End date filter
- `status` (optional): Status filter
- `limit` (optional): Result limit
- `offset` (optional): Pagination offset

**Permissions Required:** `notulen:read`

**Response:** Same as list endpoint

#### GET /api/notulen/:id/versions
Get all versions of a meeting minutes

**Path Parameters:**
- `id`: UUID of the notulen

**Permissions Required:** `notulen:read`

**Response:**
```typescript
interface NotulenVersionsResponse {
  versions: NotulenVersie[];
}
```

#### GET /api/notulen/:id/versions/:version
Get a specific version of meeting minutes

**Path Parameters:**
- `id`: UUID of the notulen
- `version`: Version number

**Permissions Required:** `notulen:read`

**Response:** `NotulenVersie` object

## Data Format Notes

### JSONB Field Structure

The API returns JSONB fields in a structured format. The backend stores these as wrapper objects:

**Agenda Items:**
```json
{
  "agenda_items": {
    "items": [
      {
        "title": "Meeting agenda item",
        "details": "Additional details"
      }
    ]
  }
}
```

**Decisions:**
```json
{
  "besluiten": {
    "besluiten": [
      {
        "besluit": "Decision text",
        "teams": {
          "Team A": ["person1", "person2"]
        }
      }
    ]
  }
}
```

**Action Items:**
```json
{
  "actiepunten": {
    "acties": [
      {
        "actie": "Action description",
        "verantwoordelijke": "responsible_person"
      }
    ]
  }
}
```

### Date Handling

- All dates are returned in ISO 8601 format
- Input dates should be in `YYYY-MM-DD` format for date-only fields
- Timestamps include full timezone information

### Arrays

- `aanwezigen` and `afwezigen` are string arrays
- `verantwoordelijke` in action items can be either a string or string array
- All JSONB arrays are properly structured objects

## Error Handling

The API returns standard HTTP status codes:

- `200`: Success
- `400`: Bad Request (validation errors)
- `401`: Unauthorized (missing/invalid token)
- `403`: Forbidden (insufficient permissions)
- `404`: Not Found
- `500`: Internal Server Error

Error responses include:
```typescript
interface ErrorResponse {
  error: string; // Error message
}
```

## Frontend Integration Examples

### Fetching Meeting Minutes List
```typescript
const fetchNotulen = async (filters?: NotulenFilters) => {
  const params = new URLSearchParams();

  if (filters?.dateFrom) params.append('date_from', filters.dateFrom);
  if (filters?.dateTo) params.append('date_to', filters.dateTo);
  if (filters?.status) params.append('status', filters.status);
  if (filters?.limit) params.append('limit', filters.limit.toString());
  if (filters?.offset) params.append('offset', filters.offset.toString());

  const response = await fetch(`/api/notulen?${params}`, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });

  if (!response.ok) {
    throw new Error('Failed to fetch notulen');
  }

  return response.json() as NotulenListResponse;
};
```

### Creating New Meeting Minutes
```typescript
const createNotulen = async (data: NotulenCreateRequest) => {
  const response = await fetch('/api/notulen', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });

  if (!response.ok) {
    throw new Error('Failed to create notulen');
  }

  return response.json() as Notulen;
};
```

### Handling JSONB Data
```typescript
// Extract agenda items from response
const agendaItems = notulen.agenda_items?.items || [];

// Extract decisions
const besluiten = notulen.besluiten?.besluiten || [];

// Extract action items
const actiepunten = notulen.actiepunten?.acties || [];
```

## Permissions

The following permissions are required for different operations:

- `notulen:read`: View meeting minutes
- `notulen:write`: Create, update, finalize, archive meeting minutes
- `notulen:delete`: Delete meeting minutes

## Rate Limiting

API endpoints are subject to rate limiting. Public endpoints have higher limits than authenticated endpoints.

## Validation Rules

### NotulenCreateRequest Validation
- `titel`: Required, 3-255 characters
- `vergadering_datum`: Required, valid date in YYYY-MM-DD format
- `locatie`: Optional, string
- `voorzitter`: Optional, string
- `notulist`: Optional, string
- `aanwezigen`: Optional, array of strings
- `afwezigen`: Optional, array of strings
- `agenda_items`: Optional, array of AgendaItem objects
- `besluiten`: Optional, array of Besluit objects
- `actiepunten`: Optional, array of Actiepunt objects
- `notities`: Optional, string

### AgendaItem Validation
- `title`: Required, non-empty string
- `details`: Optional, string

### Besluit Validation
- `besluit`: Required, non-empty string
- `teams`: Optional, object with string keys and string array values

### Actiepunt Validation
- `actie`: Required, non-empty string
- `verantwoordelijke`: Required, string or array of strings

## Status Transitions

Meeting minutes can have the following statuses:

- **`draft`**: Initial status, can be edited
- **`finalized`**: Cannot be edited, publicly visible
- **`archived`**: Hidden from normal views, preserved for history

### Status Transition Rules
- `draft` → `finalized`: Via POST /api/notulen/:id/finalize
- `finalized` → `archived`: Via POST /api/notulen/:id/archive
- `draft` → `archived`: Via POST /api/notulen/:id/archive
- No transitions allowed from `archived` status

### Permissions by Status
- `draft`: Only creator and admins can view/edit
- `finalized`: All authenticated users can view
- `archived`: Only admins can view

## Full-Text Search

The search endpoint uses PostgreSQL's full-text search capabilities:

### Search Query Syntax
- Natural language search: `"meeting agenda"`
- Multiple terms: `"vergadering besluiten"`
- Dutch language support with stemming

### Search Fields
- `titel` (title)
- `notities` (notes)
- Weighted search with title having higher priority

### Search Filters
Can be combined with date ranges and status filters for precise results.

## Pagination

All list endpoints support pagination:

### Parameters
- `limit`: Number of items per page (default: 50, max: 100)
- `offset`: Number of items to skip (default: 0)

### Response Format
```typescript
{
  notulen: Notulen[],
  total: number,    // Total items available
  limit: number,    // Items per page
  offset: number    // Items skipped
}
```

### Calculation Examples
- Page 1: `limit=10&offset=0`
- Page 2: `limit=10&offset=10`
- Page 3: `limit=10&offset=20`

## Middleware and Security

### Authentication Middleware
- JWT token validation
- User context injection
- Automatic token refresh handling

### Permission Middleware
- RBAC-based access control
- Permission checking per endpoint
- Granular permission system

### Rate Limiting
- Public endpoints: Higher limits
- Authenticated endpoints: Standard limits
- Burst protection enabled

## Database Schema

### Tables

#### notulen
```sql
CREATE TABLE notulen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titel VARCHAR(255) NOT NULL,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[],
    afwezigen TEXT[],
    agenda_items JSONB,
    besluiten JSONB,
    actiepunten JSONB,
    notities TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    versie INTEGER NOT NULL DEFAULT 1,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    finalized_at TIMESTAMP WITH TIME ZONE,
    finalized_by UUID
);
```

#### notulen_versies (Version History)
```sql
CREATE TABLE notulen_versies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notulen_id UUID NOT NULL REFERENCES notulen(id),
    versie INTEGER NOT NULL,
    titel VARCHAR(255) NOT NULL,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[],
    afwezigen TEXT[],
    agenda_items JSONB,
    besluiten JSONB,
    actiepunten JSONB,
    notities TEXT,
    status VARCHAR(50) NOT NULL,
    gewijzigd_door UUID NOT NULL,
    gewijzigd_op TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    wijziging_reden TEXT
);
```

### Indexes
- Primary key indexes on ID fields
- Foreign key indexes on notulen_id
- GIN indexes on JSONB fields for efficient querying
- Date indexes for range queries

## Migration Information

### V1_54__create_notulen_tables.sql
- Creates notulen and notulen_versies tables
- Sets up proper constraints and indexes
- Includes JSONB fields for complex data

### V1_56__add_sample_notulen_data.sql
- Inserts sample data for testing
- Demonstrates proper JSONB structure
- Includes realistic meeting minutes content

## Testing Examples

### Create Test Notulen
```bash
curl -X POST http://localhost:8082/api/notulen \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "titel": "Test Vergadering",
    "vergadering_datum": "2025-01-15",
    "locatie": "Test Room",
    "aanwezigen": ["Alice", "Bob"],
    "agenda_items": [
      {"title": "Opening", "details": "Welcome everyone"}
    ],
    "besluiten": [
      {"besluit": "Approve budget"}
    ],
    "actiepunten": [
      {"actie": "Send invites", "verantwoordelijke": "Alice"}
    ]
  }'
```

### Search Test
```bash
curl "http://localhost:8082/api/notulen/search?q=vergadering&status=draft" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Version History Test
```bash
curl http://localhost:8082/api/notulen/YOUR_ID/versions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Troubleshooting

### Common Issues

#### JSONB Unmarshal Errors
- Ensure JSONB wrapper structure is correct
- Check that nested objects match expected schema
- Validate JSON structure before sending

#### Date Parsing Errors
- Use YYYY-MM-DD format for date inputs
- Ensure dates are valid calendar dates
- Check timezone handling in frontend

#### Permission Errors
- Verify JWT token is valid and not expired
- Check user has required permissions
- Confirm RBAC roles are properly assigned

#### Validation Errors
- Check all required fields are provided
- Verify string lengths are within limits
- Ensure array elements are correct types

### Debug Endpoints

#### Health Check
```bash
GET /api/health
```

#### Database Connection Test
```bash
GET /api/test-db
```

### Logging

All API requests are logged with:
- Request ID for tracing
- User ID (if authenticated)
- Operation type
- Response status
- Error details (if applicable)

## Change Log

### Version 1.0.0 (Current)
- Initial release of Notulen API
- Full CRUD operations
- Version history support
- Full-text search
- RBAC permissions
- JSONB field handling
- PostgreSQL TEXT[] arrays

### Planned Features
- Bulk operations
- Export to PDF/Markdown
- Email notifications
- Advanced filtering
- Audit trails

## Support

For API support:
- Check this documentation first
- Review error messages for specific guidance
- Test with provided examples
- Check application logs for detailed error information

## API Versioning

The API follows semantic versioning:
- Major version in URL path (currently v1)
- Backward compatible changes in minor versions
- Breaking changes require new major version

Current version: **v1** (stable)
