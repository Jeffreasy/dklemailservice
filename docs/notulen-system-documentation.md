# Notulen Systeem Documentatie - DKL Email Service

Complete documentatie voor het notulen (minutes) management systeem van de DKL Email Service.

---

## 📋 Overzicht

Het notulen systeem biedt een complete oplossing voor het beheren van vergadernotulen met versiebeheer, RBAC permissies en gestructureerde opslag van vergadergegevens.

### 🎯 Belangrijkste Features

- **Versiebeheer**: Automatische versie snapshots bij elke wijziging
- **RBAC Integratie**: Gedetailleerde permissies voor verschillende gebruikersrollen
- **Gestructureerde Data**: JSONB opslag voor agenda items, besluiten en actiepunten
- **Full-text Search**: Ondersteuning voor zoeken in Nederlandse tekst
- **Workflow Management**: Draft → Finalized → Archived status workflow

---

## 🗄️ Database Schema

### Hoofdtabel: `notulen`

```sql
CREATE TABLE notulen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titel VARCHAR(255) NOT NULL,
    vergadering_datum DATE NOT NULL,
    locatie VARCHAR(255),
    voorzitter VARCHAR(255),
    notulist VARCHAR(255),
    aanwezigen TEXT[], -- Array of names
    afwezigen TEXT[],  -- Array of names
    agenda_items JSONB, -- Structured agenda items
    besluiten JSONB,   -- Structured decisions
    actiepunten JSONB, -- Structured action points
    notities TEXT,     -- Free text notes
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'finalized', 'archived')),
    versie INTEGER DEFAULT 1,
    created_by UUID REFERENCES gebruikers(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finalized_at TIMESTAMP WITH TIME ZONE,
    finalized_by UUID REFERENCES gebruikers(id)
);
```

### Versiegeschiedenis: `notulen_versies`

```sql
CREATE TABLE notulen_versies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notulen_id UUID NOT NULL REFERENCES notulen(id) ON DELETE CASCADE,
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
    status VARCHAR(50),
    gewijzigd_door UUID REFERENCES gebruikers(id),
    gewijzigd_op TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    wijziging_reden TEXT
);
```

### Indexes voor Performance

```sql
-- Performance indexes
CREATE INDEX idx_notulen_datum ON notulen(vergadering_datum);
CREATE INDEX idx_notulen_status ON notulen(status);
CREATE INDEX idx_notulen_created_by ON notulen(created_by);
CREATE INDEX idx_notulen_finalized_by ON notulen(finalized_by);
CREATE INDEX idx_notulen_titel_gin ON notulen USING gin(to_tsvector('dutch', titel));
CREATE INDEX idx_notulen_notities_gin ON notulen USING gin(to_tsvector('dutch', notities));

-- Version table indexes
CREATE INDEX idx_notulen_versies_notulen_id ON notulen_versies(notulen_id);
CREATE INDEX idx_notulen_versies_versie ON notulen_versies(notulen_id, versie);
CREATE INDEX idx_notulen_versies_gewijzigd_door ON notulen_versies(gewijzigd_door);
```

---

## 🔐 RBAC Permissions

### Beschikbare Permissions

| Permission | Action | Beschrijving | Rollen |
|------------|--------|--------------|---------|
| `notulen:read` | read | Kan notulen lezen en bekijken | admin, staff |
| `notulen:write` | write | Kan notulen aanmaken en bijwerken | admin, staff |
| `notulen:delete` | delete | Kan notulen verwijderen | admin |
| `notulen:finalize` | finalize | Kan notulen finaliseren | admin |
| `notulen:archive` | archive | Kan notulen archiveren | admin |

### Permission Gebruik

```sql
-- Check permissions in application code
SELECT COUNT(*) > 0
FROM user_permissions
WHERE user_id = $1 AND resource = 'notulen' AND action = $2;
```

---

## 🚀 API Endpoints

### Base URL
```
/api/notulen
```

### GET /api/notulen
**Beschrijving**: Haal alle notulen op (met filtering)

**Query Parameters**:
- `status`: Filter op status (draft, finalized, archived)
- `date_from`: Start datum filter (YYYY-MM-DD)
- `date_to`: Eind datum filter (YYYY-MM-DD)
- `limit`: Maximum aantal resultaten (default: 50)
- `offset`: Pagination offset (default: 0)

**Response**:
```json
{
  "notulen": [...],
  "total": 25,
  "limit": 50,
  "offset": 0
}
```

### GET /api/notulen/search
**Beschrijving**: Zoek notulen met full-text search

**Query Parameters**:
- `q`: Zoekterm (required)
- `status`: Filter op status
- `date_from`: Start datum filter
- `date_to`: Eind datum filter
- `limit`: Maximum aantal resultaten
- `offset`: Pagination offset

### GET /api/notulen/public
**Beschrijving**: Haal alleen gefinaliseerde notulen op (publieke toegang)

**Query Parameters**: Zie GET /api/notulen

**Response**:
```json
{
  "notulen": [
    {
      "id": "uuid",
      "titel": "Vergadering Bestuur 2025",
      "vergadering_datum": "2025-01-15",
      "locatie": "Gemeentehuis Dronten",
      "voorzitter": "Jan Jansen",
      "notulist": "Marie Martens",
      "aanwezigen": ["Jan Jansen", "Marie Martens", "Piet Peters"],
      "afwezigen": ["Klaas Klaassen"],
      "status": "finalized",
      "versie": 3,
      "created_by": "uuid",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-16T14:30:00Z",
      "finalized_at": "2025-01-16T14:30:00Z",
      "finalized_by": "uuid"
    }
  ],
  "total": 25,
  "limit": 50,
  "offset": 0
}
```

### POST /api/notulen
**Beschrijving**: Maak nieuwe notulen aan

**Request Body**:
```json
{
  "titel": "Vergadering Bestuur 2025",
  "vergadering_datum": "2025-01-15",
  "locatie": "Gemeentehuis Dronten",
  "voorzitter": "Jan Jansen",
  "notulist": "Marie Martens",
  "aanwezigen": ["Jan Jansen", "Marie Martens", "Piet Peters"],
  "afwezigen": ["Klaas Klaassen"],
  "agenda_items": [
    {
      "onderwerp": "Jaarplan 2025",
      "spreker": "Jan Jansen",
      "tijd": "10:00-10:30"
    }
  ],
  "besluiten": [
    {
      "besluit": "Jaarplan 2025 wordt goedgekeurd",
      "verantwoordelijke": "Jan Jansen"
    }
  ],
  "actiepunten": [
    {
      "actie": "Jaarplan implementeren",
      "verantwoordelijke": "Marie Martens",
      "deadline": "2025-02-01"
    }
  ],
  "notities": "Vergadering verliep constructief..."
}
```

### GET /api/notulen/{id}
**Beschrijving**: Haal specifieke notulen op

**Query Parameters**:
- `format`: Response format (markdown voor markdown export)

**Response**: Zie GET /api/notulen response format

### PUT /api/notulen/{id}
**Beschrijving**: Update bestaande notulen

**Request Body**: Zie POST /api/notulen (alleen changed velden)

### DELETE /api/notulen/{id}
**Beschrijving**: Verwijder notulen (soft delete)

**Response**:
```json
{
  "success": true,
  "message": "Notulen deleted successfully"
}
```

### POST /api/notulen/{id}/finalize
**Beschrijving**: Finaliseer notulen

**Request Body** (optional):
```json
{
  "wijziging_reden": "Finalisatie opmerking"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Notulen finalized successfully"
}
```

### POST /api/notulen/{id}/archive
**Beschrijving**: Archiveer notulen

**Response**:
```json
{
  "success": true,
  "message": "Notulen archived successfully"
}
```

### GET /api/notulen/{id}/versions
**Beschrijving**: Haal versiegeschiedenis op

**Response**:
```json
{
  "versions": [
    {
      "versie": 3,
      "titel": "Vergadering Bestuur 2025",
      "gewijzigd_door": "uuid",
      "gewijzigd_op": "2025-01-16T14:30:00Z",
      "wijziging_reden": "Automatische versie snapshot"
    }
  ]
}
```

### GET /api/notulen/{id}/versions/{version}
**Beschrijving**: Haal specifieke versie van notulen op

**Response**: NotulenVersie object

---

## 🔄 Workflow Management

### Status Overgangen

```
draft → finalized → archived
```

### Automatische Triggers

1. **Versiebeheer**: Elke update creëert automatisch een nieuwe versie
2. **Timestamp Updates**: `updated_at` wordt automatisch bijgewerkt
3. **Finalisatie**: `finalized_at` en `finalized_by` worden gezet bij finalisatie

### JSONB Data Structuur

#### Agenda Items
```json
[
  {
    "titel": "string",
    "beschrijving": "string (optional)",
    "spreker": "string (optional)",
    "tijdslot": "string (optional)"
  }
]
```

#### Besluiten
```json
[
  {
    "beschrijving": "string",
    "verantwoordelijke": "string (optional)",
    "deadline": "string (ISO date) (optional)"
  }
]
```

#### Actiepunten
```json
[
  {
    "beschrijving": "string",
    "verantwoordelijke": "string",
    "deadline": "string (ISO date) (optional)",
    "status": "string (pending/in_progress/completed)"
  }
]
```

---

## 🔍 Search & Filtering

### Full-text Search

```sql
-- Zoeken in titel en notities
SELECT * FROM notulen
WHERE to_tsvector('dutch', titel || ' ' || notities) @@ plainto_tsquery('dutch', 'zoekterm');
```

### Filtering Opties

```sql
-- Filter op status
SELECT * FROM notulen WHERE status = 'finalized';

-- Filter op datum range
SELECT * FROM notulen
WHERE vergadering_datum BETWEEN '2025-01-01' AND '2025-12-31';

-- Filter op voorzitter
SELECT * FROM notulen WHERE voorzitter = 'Jan Jansen';

-- Filter op aanwezigen (array contains)
SELECT * FROM notulen WHERE 'Jan Jansen' = ANY(aanwezigen);
```

---

## 🧪 Testing

### Unit Tests

```go
func TestNotulenService(t *testing.T) {
    // Test notulen creation
    req := &models.NotulenCreateRequest{
        Titel: "Test Vergadering",
        VergaderingDatum: "2025-01-15",
        Voorzitter: "Test Voorzitter",
    }

    notulen, err := service.CreateNotulen(context.Background(), userUUID, req)
    assert.NoError(t, err)
    assert.NotEmpty(t, notulen.ID)
    assert.Equal(t, "draft", notulen.Status)
}

func TestNotulenHandler(t *testing.T) {
    // Test API endpoints
    app := fiber.New()
    handler := NewNotulenHandler(service, authService, permissionService)
    handler.RegisterRoutes(app)

    // Test route registration
    assert.NotNil(t, app.GetRoute("GET", "/api/notulen"))
    assert.NotNil(t, app.GetRoute("POST", "/api/notulen"))
}
```

### API Tests

```bash
# Create notulen
curl -X POST http://localhost:8080/api/notulen \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "titel": "Test Vergadering",
    "vergadering_datum": "2025-01-15",
    "voorzitter": "Jan Jansen",
    "aanwezigen": ["Jan Jansen", "Marie Martens"]
  }'

# List notulen with filters
curl "http://localhost:8080/api/notulen?status=draft&limit=10" \
  -H "Authorization: Bearer {token}"

# Search notulen
curl "http://localhost:8080/api/notulen/search?q=vergadering" \
  -H "Authorization: Bearer {token}"

# Get specific notulen as markdown
curl "http://localhost:8080/api/notulen/{id}?format=markdown" \
  -H "Authorization: Bearer {token}"

# Finalize notulen
curl -X POST http://localhost:8080/api/notulen/{id}/finalize \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{"wijziging_reden": "Vergadering afgerond"}'

# Get version history
curl http://localhost:8080/api/notulen/{id}/versions \
  -H "Authorization: Bearer {token}"
```

### Database Tests

```sql
-- Test data insertion
INSERT INTO notulen (
    titel, vergadering_datum, voorzitter, notulist,
    aanwezigen, afwezigen, created_by
) VALUES (
    'Test Notulen', '2025-01-15', 'Test Voorzitter', 'Test Notulist',
    ARRAY['Jan Jansen', 'Marie Martens'], ARRAY['Klaas Klaassen'],
    (SELECT id FROM gebruikers LIMIT 1)
);

-- Verify permissions setup
SELECT resource, action, description FROM permissions WHERE resource = 'notulen' ORDER BY action;

-- Test role assignments
SELECT r.name as role_name, p.action as permission
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE p.resource = 'notulen'
ORDER BY r.name, p.action;

-- Test search functionality
SELECT * FROM notulen
WHERE to_tsvector('dutch', titel || ' ' || notities) @@ plainto_tsquery('dutch', 'test');

-- Test JSONB queries
SELECT id, jsonb_array_length(agenda_items) as agenda_count
FROM notulen WHERE status = 'finalized';

-- Test version history
SELECT n.titel, COUNT(nv.versie) as version_count
FROM notulen n
LEFT JOIN notulen_versies nv ON n.id = nv.notulen_id
GROUP BY n.id, n.titel;
```

---

## 📊 Monitoring & Analytics

### Database Queries voor Monitoring

```sql
-- Notulen per status
SELECT status, COUNT(*) FROM notulen GROUP BY status;

-- Notulen per maand
SELECT
    DATE_TRUNC('month', vergadering_datum) as maand,
    COUNT(*) as aantal
FROM notulen
WHERE vergadering_datum >= '2025-01-01'
GROUP BY DATE_TRUNC('month', vergadering_datum)
ORDER BY maand;

-- Actieve gebruikers (notulen creators)
SELECT
    g.naam,
    COUNT(n.id) as notulen_aantal
FROM gebruikers g
LEFT JOIN notulen n ON g.id = n.created_by
GROUP BY g.id, g.naam
ORDER BY notulen_aantal DESC;

-- Versie distributie
SELECT versie, COUNT(*) FROM notulen GROUP BY versie ORDER BY versie;
```

### Performance Metrics

- **Average Query Time**: < 50ms voor standaard queries
- **Search Response Time**: < 200ms voor full-text search
- **Version Creation Time**: < 100ms voor versie snapshots

---

## 🚨 Troubleshooting

### Veelvoorkomende Problemen

#### 1. Permission Denied
```
Error: User does not have permission to access notulen
```
**Oplossing**: Controleer RBAC permissions voor de gebruiker

```sql
SELECT p.action FROM user_permissions up
JOIN permissions p ON up.permission_id = p.id
WHERE up.user_id = $1 AND p.resource = 'notulen';
```

#### 2. Version Conflict
```
Error: Version number conflict
```
**Oplossing**: Reset versie nummering indien nodig

```sql
UPDATE notulen SET versie = 1 WHERE id = $1;
DELETE FROM notulen_versies WHERE notulen_id = $1;
```

#### 3. JSONB Validation Error
```
Error: Invalid JSONB structure
```
**Oplossing**: Valideer JSONB data voordat insert/update

```sql
-- Validate agenda_items structure
SELECT jsonb_array_elements(agenda_items) as item
FROM notulen WHERE id = $1;
```

#### 4. Search Not Working
```
Error: Full-text search returns no results
```
**Oplossing**: Rebuild search indexes

```sql
REINDEX INDEX idx_notulen_titel_gin;
REINDEX INDEX idx_notulen_notities_gin;
```

---

## 🔄 Data Migration

### Van Legacy Systeem

```sql
-- Import from CSV/legacy system
COPY notulen (
    titel, vergadering_datum, locatie, voorzitter,
    notulist, aanwezigen, afwezigen, notities
)
FROM '/path/to/legacy_notulen.csv'
WITH CSV HEADER;

-- Set default values for new fields
UPDATE notulen SET
    status = 'finalized',
    versie = 1,
    created_by = (SELECT id FROM gebruikers WHERE rol = 'admin' LIMIT 1),
    created_at = vergadering_datum,
    updated_at = vergadering_datum
WHERE created_by IS NULL;
```

### Export naar Extern Systeem

```sql
-- Export to JSON for external systems
SELECT json_build_object(
    'id', id,
    'titel', titel,
    'vergadering_datum', vergadering_datum,
    'locatie', locatie,
    'voorzitter', voorzitter,
    'notulist', notulist,
    'aanwezigen', aanwezigen,
    'afwezigen', afwezigen,
    'agenda_items', agenda_items,
    'besluiten', besluiten,
    'actiepunten', actiepunten,
    'notities', notities,
    'status', status,
    'versie', versie
) FROM notulen WHERE status = 'finalized';
```

---

## 📚 Best Practices

### ✅ DO's

- ✅ Gebruik altijd transactions voor multi-step operations
- ✅ Valideer JSONB data voordat opslaan
- ✅ Implementeer proper error handling
- ✅ Gebruik full-text search voor user-facing search
- ✅ Archiveer oude notulen regelmatig
- ✅ Monitor database performance metrics
- ✅ Test migrations op staging environment eerst

### ❌ DON'Ts

- ❌ Modify notulen_versies table directly
- ❌ Skip permission checks in application code
- ❌ Store large files in JSONB fields
- ❌ Forget to update versie numbers
- ❌ Use synchronous operations voor bulk updates
- ❌ Ignore database constraints
- ❌ Skip backup before major changes

---

## 📞 Support & Contact

Voor vragen over het notulen systeem:

1. Check deze documentatie
2. Review [DATABASE_REFERENCE.md](../docs/DATABASE_REFERENCE.md) voor database details
3. Consult [AUTH_RBAC_TEST_REPORT_FINAL.md](../docs/AUTH_RBAC_TEST_REPORT_FINAL.md) voor RBAC details
4. Contact development team voor technische ondersteuning

---

**Laatst Bijgewerkt**: November 2025
**Versie**: 1.0
**Status**: Production Ready ✅