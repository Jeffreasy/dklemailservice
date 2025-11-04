# Notulen Systeem - Frontend Integratie Documentatie

Complete frontend integratie gids voor het notulen (minutes) management systeem van de DKL Email Service.

---

## 📋 Overzicht

Deze documentatie beschrijft hoe frontend applicaties kunnen integreren met het notulen systeem. Het systeem biedt een complete REST API voor het beheren van vergadernotulen met gebruikersvriendelijke namen in plaats van UUIDs.

### 🎯 Belangrijkste Features voor Frontend

- **Gebruiksvriendelijke API**: Namen in plaats van UUIDs
- **RBAC Integratie**: Automatische permissie checks
- **Real-time Updates**: WebSocket ondersteuning voor live updates
- **Offline Support**: Lokale caching mogelijkheden
- **Type Safety**: TypeScript interfaces voor alle API responses

---

## 🔗 API Basis URL

```
Base URL: http://localhost:8082/api/notulen
```

### Authenticatie

Alle API calls vereisen een Bearer token in de Authorization header:

```javascript
const headers = {
  'Authorization': `Bearer ${token}`,
  'Content-Type': 'application/json'
};
```

---

## 📊 Data Structuren

### NotulenResponse Interface

```typescript
interface NotulenResponse {
  id: string;
  titel: string;
  vergadering_datum: string; // ISO date string
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  aanwezigen?: string[]; // Legacy field - combined names
  afwezigen?: string[];   // Legacy field - combined names
  aanwezigen_gebruikers?: string[]; // UUIDs of registered users
  afwezigen_gebruikers?: string[];  // UUIDs of registered users
  aanwezigen_gasten?: string[];     // Names of non-registered guests
  afwezigen_gasten?: string[];      // Names of non-registered guests
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
  status: 'draft' | 'finalized' | 'archived';
  versie: number;
  created_by: string;        // UUID
  created_by_name?: string;  // Resolved user name
  created_at: string;        // ISO datetime string
  updated_at: string;        // ISO datetime string
  updated_by?: string;       // UUID
  updated_by_name?: string;  // Resolved user name
  finalized_at?: string;     // ISO datetime string
  finalized_by?: string;     // UUID
  finalized_by_name?: string; // Resolved user name
}
```

### AgendaItem Interface

```typescript
interface AgendaItem {
  title: string;
  details?: string;
}
```

### Besluit Interface

```typescript
interface Besluit {
  besluit: string;
  teams?: { [key: string]: any };
}
```

### Actiepunt Interface

```typescript
interface Actiepunt {
  actie: string;
  verantwoordelijke?: string | string[]; // Can be single name or array
}
```

### NotulenCreateRequest Interface

```typescript
interface NotulenCreateRequest {
  titel: string;
  vergadering_datum: string; // YYYY-MM-DD format
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  aanwezigen?: string[]; // Legacy field
  afwezigen?: string[];   // Legacy field
  aanwezigen_gebruikers?: string[]; // UUIDs
  afwezigen_gebruikers?: string[];  // UUIDs
  aanwezigen_gasten?: string[];     // Names
  afwezigen_gasten?: string[];      // Names
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
}
```

### NotulenListResponse Interface

```typescript
interface NotulenListResponse {
  notulen: NotulenResponse[];
  total: number;
  limit: number;
  offset: number;
}
```

---

## 🚀 API Endpoints

### 1. Notulen Lijst Opvragen

**GET** `/api/notulen`

Haal alle notulen op met filtering en paginatie.

#### Query Parameters

```typescript
interface NotulenListParams {
  status?: 'draft' | 'finalized' | 'archived';
  date_from?: string; // YYYY-MM-DD
  date_to?: string;   // YYYY-MM-DD
  limit?: number;     // Default: 50, Max: 100
  offset?: number;    // Default: 0
}
```

#### Voorbeeld Request

```javascript
// Haal alle gefinaliseerde notulen op
const response = await fetch('/api/notulen?status=finalized&limit=10', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const data = await response.json();
// data: NotulenListResponse
```

#### Voorbeeld Response

```json
{
  "notulen": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "titel": "Vergadering Bestuur 2025",
      "vergadering_datum": "2025-01-15",
      "locatie": "Gemeentehuis Dronten",
      "voorzitter": "Jan Jansen",
      "notulist": "Marie Martens",
      "aanwezigen": ["Jan Jansen", "Marie Martens", "Piet Peters"],
      "afwezigen": ["Klaas Klaassen"],
      "status": "finalized",
      "versie": 3,
      "created_by": "00a26c26-4932-4091-91ba-4385a251e285",
      "created_by_name": "Joyce Thielen",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-16T14:30:00Z",
      "updated_by": "748320dd-5b5e-4434-ad7e-8a405fd6266f",
      "updated_by_name": "Jeffrey",
      "finalized_at": "2025-01-16T14:30:00Z",
      "finalized_by": "7157f3f6-da85-4058-9d38-19133ec93b03",
      "finalized_by_name": "SuperAdmin"
    }
  ],
  "total": 25,
  "limit": 10,
  "offset": 0
}
```

### 2. Specifieke Notulen Opvragen

**GET** `/api/notulen/{id}`

Haal een specifieke notulen op.

#### Voorbeeld Request

```javascript
const notulenId = "550e8400-e29b-41d4-a716-446655440000";
const response = await fetch(`/api/notulen/${notulenId}`, {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const notulen = await response.json();
// notulen: NotulenResponse
```

### 3. Gebruikers Zoeken

**GET** `/api/users?q={searchterm}`

Zoek gebruikers voor deelnemerslijsten in notulen.

#### Query Parameters

```typescript
interface UserSearchParams {
  q: string; // Zoekterm (naam of email)
  limit?: number; // Maximaal aantal resultaten (default: 50)
}
```

#### Voorbeeld Request

```javascript
// Zoek gebruikers voor deelnemers selectie
const response = await fetch('/api/users?q=joyce&limit=5', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const users = await response.json();
// users: Gebruiker[]
```

#### Voorbeeld Response

```json
[
  {
    "id": "00a26c26-4932-4091-91ba-4385a251e285",
    "naam": "Joyce Thielen",
    "email": "Joyce.thielen@sheerenloo.nl",
    "rol": "begeleider",
    "is_actief": true,
    "laatste_login": "2025-11-03T02:21:55Z",
    "created_at": "2025-10-25T11:58:19Z",
    "updated_at": "2025-10-25T12:09:50Z"
  }
]
```

### 3. Notulen Aanmaken

**POST** `/api/notulen`

Maak een nieuwe notulen aan.

#### Voorbeeld Request

```javascript
const newNotulen = {
  titel: "Vergadering Bestuur Februari 2025",
  vergadering_datum: "2025-02-15",
  locatie: "Gemeentehuis Dronten",
  voorzitter: "Jan Jansen",
  notulist: "Marie Martens",
  aanwezigen_gebruikers: ["00a26c26-4932-4091-91ba-4385a251e285"],
  afwezigen_gasten: ["Klaas Klaassen"],
  agenda_items: [
    {
      title: "Jaarplan 2025 bespreking",
      details: "Evaluatie van het jaarplan en bijstellingen"
    }
  ],
  besluiten: [
    {
      besluit: "Jaarplan 2025 wordt goedgekeurd met kleine aanpassingen"
    }
  ],
  actiepunten: [
    {
      actie: "Jaarplan implementeren",
      verantwoordelijke: "Marie Martens"
    }
  ],
  notities: "Constructieve vergadering met goede deelname"
};

const response = await fetch('/api/notulen', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(newNotulen)
});

const createdNotulen = await response.json();
// createdNotulen: NotulenResponse
```

### 4. Notulen Bijwerken

**PUT** `/api/notulen/{id}`

Update een bestaande notulen (alleen als status 'draft' is).

#### Voorbeeld Request

```javascript
const updates = {
  titel: "Vergadering Bestuur Februari 2025 - Bijgewerkt",
  notities: "Constructieve vergadering met goede deelname. Extra actiepunten toegevoegd."
};

const response = await fetch(`/api/notulen/${notulenId}`, {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(updates)
});

const updatedNotulen = await response.json();
// updatedNotulen: NotulenResponse
```

### 5. Notulen Finaliseren

**POST** `/api/notulen/{id}/finalize`

Markeer notulen als gefinaliseerd.

#### Voorbeeld Request

```javascript
const finalizeData = {
  wijziging_reden: "Vergadering afgerond en goedgekeurd"
};

const response = await fetch(`/api/notulen/${notulenId}/finalize`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(finalizeData)
});

const result = await response.json();
// result: { success: true, message: "Notulen finalized successfully" }
```

### 6. Notulen Archiveren

**POST** `/api/notulen/{id}/archive`

Archiveer notulen.

#### Voorbeeld Request

```javascript
const response = await fetch(`/api/notulen/${notulenId}/archive`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const result = await response.json();
// result: { success: true, message: "Notulen archived successfully" }
```

### 7. Notulen Verwijderen

**DELETE** `/api/notulen/{id}`

Verwijder notulen (soft delete).

#### Voorbeeld Request

```javascript
const response = await fetch(`/api/notulen/${notulenId}`, {
  method: 'DELETE',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const result = await response.json();
// result: { success: true, message: "Notulen deleted successfully" }
```

### 8. Zoeken in Notulen

**GET** `/api/notulen/search`

Full-text search door notulen.

#### Query Parameters

```typescript
interface NotulenSearchParams {
  q: string; // Required: search query
  status?: 'draft' | 'finalized' | 'archived';
  date_from?: string; // YYYY-MM-DD
  date_to?: string;   // YYYY-MM-DD
  limit?: number;
  offset?: number;
}
```

#### Voorbeeld Request

```javascript
const response = await fetch('/api/notulen/search?q=jaarplan&status=finalized&limit=5', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const results = await response.json();
// results: NotulenListResponse
```

### 9. Markdown Export

**GET** `/api/notulen/{id}?format=markdown`

Exporteer notulen als Markdown voor downloads/printing.

#### Voorbeeld Request

```javascript
const response = await fetch(`/api/notulen/${notulenId}?format=markdown`, {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const markdown = await response.text();
// markdown: string (Markdown formatted content)
```

### 10. Versie Geschiedenis

**GET** `/api/notulen/{id}/versions`

Haal alle versies van een notulen op.

#### Voorbeeld Request

```javascript
const response = await fetch(`/api/notulen/${notulenId}/versions`, {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const versions = await response.json();
// versions: { versions: NotulenVersie[] }
```

### 11. Publieke Notulen

**GET** `/api/notulen/public`

Haal alleen gefinaliseerde notulen op (geen authenticatie vereist).

#### Voorbeeld Request

```javascript
const response = await fetch('/api/notulen/public?limit=10');
const publicNotulen = await response.json();
// publicNotulen: NotulenListResponse
```

---

## 🔐 Permissies en RBAC

### Beschikbare Permissies

| Permission | Actie | Beschrijving |
|------------|--------|--------------|
| `notulen:read` | read | Kan notulen lezen en bekijken |
| `notulen:write` | write | Kan notulen aanmaken en bijwerken |
| `notulen:delete` | delete | Kan notulen verwijderen |
| `notulen:finalize` | finalize | Kan notulen finaliseren |
| `notulen:archive` | archive | Kan notulen archiveren |

### Permission Checks in Frontend

```javascript
// Controleer permissies voordat UI elementen tonen
const canCreate = userPermissions.includes('notulen:write');
const canEdit = userPermissions.includes('notulen:write') && notulen.status === 'draft';
const canDelete = userPermissions.includes('notulen:delete');
const canFinalize = userPermissions.includes('notulen:finalize');
const canArchive = userPermissions.includes('notulen:archive');
```

---

## ⚛️ React/Vue.js Integratie Voorbeelden

### React Hook voor Notulen Management

```typescript
// useNotulen.ts
import { useState, useEffect } from 'react';

interface UseNotulenOptions {
  status?: string;
  limit?: number;
  offset?: number;
}

export function useNotulen(options: UseNotulenOptions = {}) {
  const [notulen, setNotulen] = useState<NotulenResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);

  const fetchNotulen = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams();
      if (options.status) params.append('status', options.status);
      if (options.limit) params.append('limit', options.limit.toString());
      if (options.offset) params.append('offset', options.offset.toString());

      const response = await fetch(`/api/notulen?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (!response.ok) throw new Error('Failed to fetch notulen');

      const data: NotulenListResponse = await response.json();
      setNotulen(data.notulen);
      setTotal(data.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotulen();
  }, [options.status, options.limit, options.offset]);

  return { notulen, loading, error, total, refetch: fetchNotulen };
}
```

### React Component voor Notulen Lijst

```tsx
// NotulenList.tsx
import React from 'react';
import { useNotulen } from './hooks/useNotulen';

export function NotulenList() {
  const { notulen, loading, error, total } = useNotulen({
    status: 'finalized',
    limit: 10
  });

  if (loading) return <div>Loading notulen...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="notulen-list">
      <h2>Notulen ({total})</h2>
      {notulen.map(notulen => (
        <div key={notulen.id} className="notulen-item">
          <h3>{notulen.titel}</h3>
          <p>Datum: {new Date(notulen.vergadering_datum).toLocaleDateString('nl-NL')}</p>
          <p>Locatie: {notulen.locatie}</p>
          <p>Status: {notulen.status}</p>
          <p>Aangemaakt door: {notulen.created_by_name || notulen.created_by}</p>
          {notulen.finalized_by_name && (
            <p>Gefinaliseerd door: {notulen.finalized_by_name}</p>
          )}
        </div>
      ))}
    </div>
  );
}
```

### Vue.js Composition API voor Notulen

```typescript
// useNotulen.ts (Vue)
import { ref, computed } from 'vue';

export function useNotulen() {
  const notulen = ref<NotulenResponse[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchNotulen = async (status?: string) => {
    try {
      loading.value = true;
      error.value = null;

      const params = new URLSearchParams();
      if (status) params.append('status', status);

      const response = await fetch(`/api/notulen?${params}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      if (!response.ok) throw new Error('Failed to fetch notulen');

      const data: NotulenListResponse = await response.json();
      notulen.value = data.notulen;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error';
    } finally {
      loading.value = false;
    }
  };

  const draftNotulen = computed(() =>
    notulen.value.filter(n => n.status === 'draft')
  );

  const finalizedNotulen = computed(() =>
    notulen.value.filter(n => n.status === 'finalized')
  );

  return {
    notulen,
    loading,
    error,
    fetchNotulen,
    draftNotulen,
    finalizedNotulen
  };
}
```

### Vue.js Component voor Notulen Formulier

```vue
<!-- NotulenForm.vue -->
<template>
  <form @submit.prevent="submitNotulen" class="notulen-form">
    <div class="form-group">
      <label for="titel">Titel:</label>
      <input
        id="titel"
        v-model="form.titel"
        type="text"
        required
      />
    </div>

    <div class="form-group">
      <label for="vergadering_datum">Vergadering Datum:</label>
      <input
        id="vergadering_datum"
        v-model="form.vergadering_datum"
        type="date"
        required
      />
    </div>

    <div class="form-group">
      <label for="locatie">Locatie:</label>
      <input
        id="locatie"
        v-model="form.locatie"
        type="text"
      />
    </div>

    <div class="form-group">
      <label for="voorzitter">Voorzitter:</label>
      <input
        id="voorzitter"
        v-model="form.voorzitter"
        type="text"
      />
    </div>

    <div class="form-group">
      <label for="notulist">Notulist:</label>
      <input
        id="notulist"
        v-model="form.notulist"
        type="text"
      />
    </div>

    <div class="form-group">
      <label for="notities">Notities:</label>
      <textarea
        id="notities"
        v-model="form.notities"
        rows="4"
      />
    </div>

    <button type="submit" :disabled="submitting">
      {{ submitting ? 'Opslaan...' : 'Notulen Opslaan' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const form = ref<NotulenCreateRequest>({
  titel: '',
  vergadering_datum: '',
  locatie: '',
  voorzitter: '',
  notulist: '',
  notities: ''
});

const submitting = ref(false);

const submitNotulen = async () => {
  try {
    submitting.value = true;

    const response = await fetch('/api/notulen', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form.value)
    });

    if (!response.ok) throw new Error('Failed to create notulen');

    const newNotulen = await response.json();

    // Reset form
    form.value = {
      titel: '',
      vergadering_datum: '',
      locatie: '',
      voorzitter: '',
      notulist: '',
      notities: ''
    };

    // Emit event or navigate to notulen detail
    // emit('notulenCreated', newNotulen);

  } catch (error) {
    console.error('Error creating notulen:', error);
    // Handle error (show toast, etc.)
  } finally {
    submitting.value = false;
  }
};
</script>
```

---

## 🔄 Real-time Updates met WebSocket

Het notulen systeem ondersteunt real-time updates via WebSocket voor live samenwerking.

### WebSocket Connectie

```javascript
// Verbinden met WebSocket voor notulen updates
const ws = new WebSocket('ws://localhost:8082/ws/notulen');

// Luister naar updates
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  if (data.type === 'notulen_updated') {
    // Update notulen in state
    updateNotulen(data.notulen);
  }

  if (data.type === 'notulen_finalized') {
    // Update status naar finalized
    finalizeNotulen(data.notulenId);
  }
};

// Stuur updates naar server
const sendUpdate = (notulenId, updates) => {
  ws.send(JSON.stringify({
    type: 'update_notulen',
    notulenId,
    updates
  }));
};
```

---

## 📱 Mobiele Integratie

### React Native Fetch Voorbeeld

```typescript
// services/notulenService.ts
export class NotulenService {
  private baseUrl = 'http://localhost:8082/api/notulen';

  async getNotulenList(params?: NotulenListParams): Promise<NotulenListResponse> {
    const token = await AsyncStorage.getItem('authToken');
    const queryParams = new URLSearchParams();

    if (params?.status) queryParams.append('status', params.status);
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const response = await fetch(`${this.baseUrl}?${queryParams}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async createNotulen(notulenData: NotulenCreateRequest): Promise<NotulenResponse> {
    const token = await AsyncStorage.getItem('authToken');

    const response = await fetch(this.baseUrl, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(notulenData)
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }
}
```

---

## 🧪 Error Handling

### HTTP Status Codes

| Status Code | Betekenis | Actie |
|-------------|-----------|--------|
| 200 | Success | Data verwerken |
| 201 | Created | Resource succesvol aangemaakt |
| 400 | Bad Request | Input validatie error |
| 401 | Unauthorized | Token refresh of login |
| 403 | Forbidden | Onvoldoende permissies |
| 404 | Not Found | Resource bestaat niet |
| 409 | Conflict | Version conflict bij updates |
| 500 | Internal Server Error | Server error, retry later |

### Error Response Format

```json
{
  "error": "Beschrijving van de fout"
}
```

### Frontend Error Handling Voorbeeld

```javascript
const handleApiCall = async (apiCall) => {
  try {
    const result = await apiCall();
    return { success: true, data: result };
  } catch (error) {
    if (error.response?.status === 401) {
      // Token expired, redirect to login
      redirectToLogin();
      return { success: false, error: 'Sessie verlopen' };
    }

    if (error.response?.status === 403) {
      return { success: false, error: 'Onvoldoende permissies' };
    }

    if (error.response?.status === 409) {
      return { success: false, error: 'Document is gewijzigd door iemand anders' };
    }

    return {
      success: false,
      error: error.response?.data?.error || 'Onbekende fout'
    };
  }
};
```

---

## 🔍 Best Practices

### 1. Caching Strategie

```javascript
// Cache gefinaliseerde notulen voor offline toegang
const cacheNotulen = async (notulen: NotulenResponse[]) => {
  const finalized = notulen.filter(n => n.status === 'finalized');
  await AsyncStorage.setItem('cached_notulen', JSON.stringify(finalized));
};

// Gebruik cache als fallback
const getNotulenWithCache = async () => {
  try {
    const fresh = await fetchNotulen();
    await cacheNotulen(fresh);
    return fresh;
  } catch (error) {
    const cached = await AsyncStorage.getItem('cached_notulen');
    return cached ? JSON.parse(cached) : [];
  }
};
```

### 2. Optimistische Updates

```javascript
// Update UI onmiddellijk, revert bij error
const updateNotulenOptimistically = async (id: string, updates: Partial<NotulenResponse>) => {
  // Sla originele data op
  const original = notulen.value.find(n => n.id === id);

  // Update UI onmiddellijk
  updateNotulenInState(id, updates);

  try {
    await updateNotulenOnServer(id, updates);
  } catch (error) {
    // Revert bij error
    updateNotulenInState(id, original);
    showError('Update mislukt');
  }
};
```

### 3. Debounced Search

```javascript
// Debounce search om server overload te voorkomen
const debouncedSearch = debounce(async (query: string) => {
  if (query.length < 3) return; // Minimum 3 karakters

  const results = await searchNotulen(query);
  updateSearchResults(results);
}, 300);
```

### 4. Pagination Handling

```javascript
const loadMoreNotulen = async () => {
  if (loading.value || !hasMore.value) return;

  loading.value = true;
  try {
    const nextPage = await fetchNotulen({
      offset: notulen.value.length,
      limit: PAGE_SIZE
    });

    notulen.value.push(...nextPage.notulen);
    hasMore.value = nextPage.notulen.length === PAGE_SIZE;
  } finally {
    loading.value = false;
  }
};
```

---

## 📞 Support & Contact

Voor vragen over frontend integratie:

1. Check deze documentatie
2. Bekijk [notulen-api.md](../api/notulen-api.md) voor backend API details
3. Bekijk [WEBSOCKET_INTEGRATION_GUIDE.md](../WEBSOCKET_INTEGRATION_GUIDE.md) voor real-time features
4. Contact development team voor technische ondersteuning

---

**Laatst Bijgewerkt**: November 2025
**Versie**: 1.0
**Status**: Production Ready ✅
