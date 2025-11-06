# Notulen API - Complete Frontend Documentatie

Complete integratie guide voor het notulen systeem met REST API en WebSocket real-time updates.

---

## 📋 Inhoudsopgave

1. [Quick Start](#quick-start)
2. [TypeScript Types](#typescript-types)
3. [REST API Endpoints](#rest-api-endpoints)
4. [WebSocket Real-time Updates](#websocket-real-time-updates)
5. [Complete Voorbeelden](#complete-voorbeelden)
6. [Error Handling](#error-handling)

---

## 🚀 Quick Start

### Installatie

```bash
npm install axios socket.io-client
# of
yarn add axios socket.io-client
```

### Environment Setup

```env
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

---

## 📦 TypeScript Types

```typescript
// types/notulen.ts

export interface AgendaItem {
  title: string;
  details?: string;
}

export interface Besluit {
  besluit: string;
  teams?: Record<string, any>;
}

export interface Actiepunt {
  actie: string;
  verantwoordelijke: string | string[];
}

export interface NotulenCreateRequest {
  titel: string;
  vergadering_datum: string; // YYYY-MM-DD format
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  
  // NIEUWE: UUID-based participant tracking
  aanwezigen_gebruikers?: string[]; // Array van user UUIDs
  afwezigen_gebruikers?: string[];  // Array van user UUIDs
  aanwezigen_gasten?: string[];     // Array van guest namen
  afwezigen_gasten?: string[];      // Array van guest namen
  
  // LEGACY: Nog steeds ondersteund voor backwards compatibility
  aanwezigen?: string[];
  afwezigen?: string[];
  
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
}

export interface NotulenUpdateRequest {
  titel?: string;
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  
  // Participant updates
  aanwezigen_gebruikers?: string[];
  afwezigen_gebruikers?: string[];
  aanwezigen_gasten?: string[];
  afwezigen_gasten?: string[];
  aanwezigen?: string[]; // Legacy
  afwezigen?: string[];  // Legacy
  
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
}

export interface NotulenResponse {
  id: string;
  titel: string;
  vergadering_datum: string;
  locatie?: string;
  voorzitter?: string;
  notulist?: string;
  
  // Participant data
  aanwezigen?: string[];         // Legacy combined names
  afwezigen?: string[];          // Legacy combined names
  aanwezigen_gebruikers?: string[]; // UUIDs van registered users
  afwezigen_gebruikers?: string[];  // UUIDs van registered users
  aanwezigen_gasten?: string[];     // Namen van guests
  afwezigen_gasten?: string[];      // Namen van guests
  
  agenda_items?: AgendaItem[];
  besluiten?: Besluit[];
  actiepunten?: Actiepunt[];
  notities?: string;
  
  status: 'draft' | 'finalized' | 'archived';
  versie: number;
  
  // User info with resolved names
  created_by: string;
  created_by_name?: string;
  created_at: string;
  updated_at: string;
  updated_by?: string;
  updated_by_name?: string;
  finalized_at?: string;
  finalized_by?: string;
  finalized_by_name?: string;
}

export interface NotulenListResponse {
  notulen: NotulenResponse[];
  total: number;
  limit: number;
  offset: number;
}

export interface NotulenVersie {
  id: string;
  notulen_id: string;
  versie: number;
  titel: string;
  vergadering_datum: string;
  // ... same fields as NotulenResponse
  gewijzigd_door: string;
  gewijzigd_op: string;
  wijziging_reden?: string;
}

// WebSocket event types
export interface NotulenWebSocketEvent {
  type: 'notulen_updated' | 'notulen_finalized' | 'notulen_archived' | 'notulen_deleted' | 'welcome' | 'ping' | 'pong';
  notulenId?: string;
  userId?: string;
  data?: NotulenResponse; // Only for 'notulen_updated'
  timestamp: string;
  message?: string; // For welcome messages
}
```

---

## 🌐 REST API Endpoints

### Base URL
```
http://localhost:8080/api/notulen
```

### 1. Create Notulen

**POST** `/api/notulen`

**Headers**:
```http
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "titel": "Vergadering Bestuur November 2025",
  "vergadering_datum": "2025-11-04",
  "locatie": "Teams",
  "voorzitter": "Salih",
  "notulist": "Jeffrey",
  
  "aanwezigen_gebruikers": [
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440002"
  ],
  "afwezigen_gebruikers": [],
  "aanwezigen_gasten": ["Externe adviseur"],
  "afwezigen_gasten": [],
  
  "agenda_items": [
    {
      "title": "Opening en mededelingen",
      "details": "Welkom en agendacheck"
    },
    {
      "title": "DKL 2026 Planning",
      "details": "Datum, locatie en doelstelling bespreken"
    }
  ],
  
  "besluiten": [
    {
      "besluit": "Datum vastgesteld op 16 mei 2026",
      "teams": {
        "Organisatie": ["Salih", "Angelique"]
      }
    }
  ],
  
  "actiepunten": [
    {
      "actie": "Partners benaderen voor samenwerking",
      "verantwoordelijke": "Salih"
    }
  ],
  
  "notities": "Constructieve vergadering met goede energie."
}
```

**Response** (201 Created):
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "titel": "Vergadering Bestuur November 2025",
  "vergadering_datum": "2025-11-04",
  "status": "draft",
  "versie": 1,
  "created_by": "550e8400-e29b-41d4-a716-446655440001",
  "created_by_name": "Jeffrey Admin",
  "created_at": "2025-11-04T23:00:00Z",
  "updated_at": "2025-11-04T23:00:00Z"
}
```

### 2. List Notulen

**GET** `/api/notulen?status=draft&limit=20&offset=0`

**Query Parameters**:
- `status`: Filter op status (`draft`, `finalized`, `archived`)
- `date_from`: Start datum (YYYY-MM-DD)
- `date_to`: Eind datum (YYYY-MM-DD)
- `limit`: Max resultaten (default: 50)
- `offset`: Pagination offset (default: 0)

**Response** (200 OK):
```json
{
  "notulen": [
    {
      "id": "uuid",
      "titel": "Vergadering...",
      "vergadering_datum": "2025-11-04",
      "status": "draft",
      "versie": 1,
      "created_by_name": "Jeffrey Admin",
      "created_at": "2025-11-04T23:00:00Z"
    }
  ],
  "total": 15,
  "limit": 20,
  "offset": 0
}
```

### 3. Get Specific Notulen

**GET** `/api/notulen/{id}`

**Query Parameters** (optional):
- `format=markdown` - Voor markdown export

**Response** (200 OK):
```json
{
  "id": "uuid",
  "titel": "Vergadering Bestuur November 2025",
  "vergadering_datum": "2025-11-04",
  "locatie": "Teams",
  "voorzitter": "Salih",
  "notulist": "Jeffrey",
  
  "aanwezigen_gebruikers": ["uuid1", "uuid2"],
  "afwezigen_gebruikers": [],
  "aanwezigen_gasten": ["Externe adviseur"],
  "afwezigen_gasten": [],
  
  "agenda_items": [...],
  "besluiten": [...],
  "actiepunten": [...],
  "notities": "...",
  
  "status": "draft",
  "versie": 1,
  "created_by": "uuid",
  "created_by_name": "Jeffrey Admin",
  "created_at": "2025-11-04T23:00:00Z",
  "updated_at": "2025-11-04T23:00:00Z"
}
```

### 4. Update Notulen

**PUT** `/api/notulen/{id}`

**Request Body** (partial updates):
```json
{
  "titel": "Updated titel",
  "aanwezigen_gebruikers": ["uuid1", "uuid2", "uuid3"],
  "agenda_items": [
    {
      "title": "Nieuw agendapunt",
      "details": "Details hier"
    }
  ]
}
```

**Response** (200 OK): Updated NotulenResponse

### 5. Finalize Notulen

**POST** `/api/notulen/{id}/finalize`

**Request Body** (optional):
```json
{
  "wijziging_reden": "Notulen goedgekeurd door bestuur"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Notulen finalized successfully"
}
```

### 6. Archive Notulen

**POST** `/api/notulen/{id}/archive`

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Notulen archived successfully"
}
```

### 7. Delete Notulen

**DELETE** `/api/notulen/{id}`

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Notulen deleted successfully"
}
```

### 8. Search Notulen

**GET** `/api/notulen/search?q=planning&status=finalized`

**Query Parameters**:
- `q`: Zoekterm (required)
- `status`, `date_from`, `date_to`, `limit`, `offset`: Zie List endpoint

**Response**: Same as List endpoint

### 9. Get Versions

**GET** `/api/notulen/{id}/versions`

**Response** (200 OK):
```json
{
  "versions": [
    {
      "versie": 3,
      "titel": "Vergadering...",
      "gewijzigd_door": "uuid",
      "gewijzigd_op": "2025-11-04T23:00:00Z",
      "wijziging_reden": "Automatic version snapshot"
    },
    {
      "versie": 2,
      "titel": "Vergadering...",
      "gewijzigd_door": "uuid",
      "gewijzigd_op": "2025-11-04T22:00:00Z"
    }
  ]
}
```

### 10. Get Specific Version

**GET** `/api/notulen/{id}/versions/{version}`

**Response** (200 OK): NotulenVersie object

### 11. List Public Notulen

**GET** `/api/notulen/public?limit=10`

**Description**: Publieke endpoint (geen auth) - Alleen finalized notulen

**Response**: Same as List endpoint, maar alleen finalized items

---

## 🔌 WebSocket Real-time Updates

### Connection Setup

```typescript
// services/notulenWebSocket.ts

import { NotulenWebSocketEvent, NotulenResponse } from '@/types/notulen';

export class NotulenWebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;
  private listeners: Map<string, Function[]> = new Map();

  constructor(
    private baseUrl: string = 'ws://localhost:8080',
    private getToken: () => string | null
  ) {}

  /**
   * Connect to WebSocket voor real-time notulen updates
   * @param notulenId - Optioneel: specifiek notulen document om te volgen
   */
  connect(notulenId?: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const token = this.getToken();
      
      // Build WebSocket URL with optional filters
      let wsUrl = `${this.baseUrl}/api/ws/notulen`;
      const params = new URLSearchParams();
      
      if (token) params.append('token', token);
      if (notulenId) params.append('notulen_id', notulenId);
      
      if (params.toString()) {
        wsUrl += `?${params.toString()}`;
      }

      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        console.log('Notulen WebSocket connected');
        this.reconnectAttempts = 0;
        resolve();
      };

      this.ws.onclose = () => {
        console.log('Notulen WebSocket disconnected');
        this.handleReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('Notulen WebSocket error:', error);
        reject(error);
      };

      this.ws.onmessage = (event) => {
        try {
          const data: NotulenWebSocketEvent = JSON.parse(event.data);
          this.handleMessage(data);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };
    });
  }

  /**
   * Disconnect from WebSocket
   */
  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  /**
   * Subscribe to specific event types
   */
  on(eventType: string, callback: Function): void {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, []);
    }
    this.listeners.get(eventType)!.push(callback);
  }

  /**
   * Unsubscribe from event
   */
  off(eventType: string, callback: Function): void {
    const callbacks = this.listeners.get(eventType);
    if (callbacks) {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }

  /**
   * Send ping to server (heartbeat)
   */
  ping(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type: 'ping' }));
    }
  }

  private handleMessage(event: NotulenWebSocketEvent): void {
    console.log('Notulen WebSocket event:', event.type);

    // Call registered listeners
    const callbacks = this.listeners.get(event.type) || [];
    callbacks.forEach(callback => callback(event));

    // Call global listeners
    const globalCallbacks = this.listeners.get('*') || [];
    globalCallbacks.forEach(callback => callback(event));
  }

  private handleReconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`);
      
      setTimeout(() => {
        this.connect().catch(console.error);
      }, this.reconnectDelay);
    } else {
      console.error('Max reconnect attempts reached');
    }
  }
}
```

### WebSocket Event Types

| Event Type | Wanneer | Data Field |
|------------|---------|------------|
| `welcome` | Bij connectie | `message`, `user_id`, `notulen_id` |
| `notulen_updated` | Notulen gewijzigd | Volledige `NotulenResponse` |
| `notulen_finalized` | Notulen gefinaliseerd | `notulenId`, `userId` |
| `notulen_archived` | Notulen gearchiveerd | `notulenId`, `userId` |
| `notulen_deleted` | Notulen verwijderd | `notulenId`, `userId` |
| `ping` | Server heartbeat | `timestamp` |
| `pong` | Client response | `timestamp` |

---

## 💻 Complete Voorbeelden

### React Hook: useNotulen

```typescript
// hooks/useNotulen.ts

import { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { NotulenWebSocketService } from '@/services/notulenWebSocket';
import { NotulenResponse, NotulenListResponse, NotulenCreateRequest } from '@/types/notulen';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080';

export function useNotulen(notulenId?: string) {
  const [notulen, setNotulen] = useState<NotulenResponse | null>(null);
  const [notulenList, setNotulenList] = useState<NotulenResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const [wsService] = useState(() => 
    new NotulenWebSocketService(WS_URL, () => localStorage.getItem('auth_token'))
  );

  // Fetch single notulen
  const fetchNotulen = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('auth_token');
      const response = await axios.get<NotulenResponse>(
        `${API_URL}/api/notulen/${id}`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );
      
      setNotulen(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch notulen');
    } finally {
      setLoading(false);
    }
  }, []);

  // Fetch notulen list
  const fetchNotulenList = useCallback(async (params?: {
    status?: string;
    date_from?: string;
    date_to?: string;
    limit?: number;
    offset?: number;
  }) => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('auth_token');
      const queryParams = new URLSearchParams(params as any).toString();
      
      const response = await axios.get<NotulenListResponse>(
        `${API_URL}/api/notulen${queryParams ? `?${queryParams}` : ''}`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );
      
      setNotulenList(response.data.notulen);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch notulen list');
    } finally {
      setLoading(false);
    }
  }, []);

  // Create notulen
  const createNotulen = useCallback(async (data: NotulenCreateRequest) => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('auth_token');
      const response = await axios.post<NotulenResponse>(
        `${API_URL}/api/notulen`,
        data,
        {
          headers: { 
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      
      return response.data;
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create notulen');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Update notulen
  const updateNotulen = useCallback(async (id: string, data: Partial<NotulenCreateRequest>) => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('auth_token');
      const response = await axios.put<NotulenResponse>(
        `${API_URL}/api/notulen/${id}`,
        data,
        {
          headers: { 
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      
      return response.data;
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update notulen');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Finalize notulen
  const finalizeNotulen = useCallback(async (id: string, reason?: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('auth_token');
      await axios.post(
        `${API_URL}/api/notulen/${id}/finalize`,
        { wijziging_reden: reason },
        {
          headers: { 
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to finalize notulen');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Delete notulen
  const deleteNotulen = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const token = localStorage.getItem('auth_token');
      await axios.delete(
        `${API_URL}/api/notulen/${id}`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete notulen');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Setup WebSocket listeners
  useEffect(() => {
    // Connect to WebSocket
    wsService.connect(notulenId).catch(console.error);

    // Handle notulen updates
    const handleUpdate = (event: NotulenWebSocketEvent) => {
      if (event.data) {
        setNotulen(event.data);
        
        // Also update in list if present
        setNotulenList(prev => 
          prev.map(n => n.id === event.data!.id ? event.data! : n)
        );
      }
    };

    // Handle finalization
    const handleFinalized = (event: NotulenWebSocketEvent) => {
      if (notulen && event.notulenId === notulen.id) {
        setNotulen(prev => prev ? { ...prev, status: 'finalized' } : null);
      }
      
      setNotulenList(prev =>
        prev.map(n => n.id === event.notulenId ? { ...n, status: 'finalized' } : n)
      );
    };

    // Handle deletion
    const handleDeleted = (event: NotulenWebSocketEvent) => {
      if (notulen && event.notulenId === notulen.id) {
        setNotulen(null);
      }
      
      setNotulenList(prev => prev.filter(n => n.id !== event.notulenId));
    };

    // Register listeners
    wsService.on('notulen_updated', handleUpdate);
    wsService.on('notulen_finalized', handleFinalized);
    wsService.on('notulen_deleted', handleDeleted);

    // Cleanup
    return () => {
      wsService.off('notulen_updated', handleUpdate);
      wsService.off('notulen_finalized', handleFinalized);
      wsService.off('notulen_deleted', handleDeleted);
      wsService.disconnect();
    };
  }, [notulenId, wsService, notulen]);

  return {
    notulen,
    notulenList,
    loading,
    error,
    fetchNotulen,
    fetchNotulenList,
    createNotulen,
    updateNotulen,
    finalizeNotulen,
    deleteNotulen,
  };
}
```

### React Component Example

```typescript
// components/NotulenEditor.tsx

import React, { useEffect, useState } from 'react';
import { useNotulen } from '@/hooks/useNotulen';
import { NotulenCreateRequest } from '@/types/notulen';

export function NotulenEditor({ notulenId }: { notulenId?: string }) {
  const {
    notulen,
    loading,
    error,
    fetchNotulen,
    updateNotulen,
    finalizeNotulen
  } = useNotulen(notulenId);

  const [formData, setFormData] = useState<Partial<NotulenCreateRequest>>({});

  useEffect(() => {
    if (notulenId) {
      fetchNotulen(notulenId);
    }
  }, [notulenId, fetchNotulen]);

  useEffect(() => {
    if (notulen) {
      setFormData({
        titel: notulen.titel,
        locatie: notulen.locatie,
        voorzitter: notulen.voorzitter,
        notulist: notulen.notulist,
        aanwezigen_gebruikers: notulen.aanwezigen_gebruikers,
        afwezigen_gebruikers: notulen.afwezigen_gebruikers,
        aanwezigen_gasten: notulen.aanwezigen_gasten,
        afwezigen_gasten: notulen.afwezigen_gasten,
        agenda_items: notulen.agenda_items,
        besluiten: notulen.besluiten,
        actiepunten: notulen.actiepunten,
        notities: notulen.notities,
      });
    }
  }, [notulen]);

  const handleSave = async () => {
    if (!notulenId) return;
    
    try {
      await updateNotulen(notulenId, formData);
      alert('Notulen opgeslagen! 📝');
    } catch (error) {
      console.error('Save failed:', error);
    }
  };

  const handleFinalize = async () => {
    if (!notulenId) return;
    
    if (confirm('Notulen finaliseren? Dit kan niet ongedaan worden gemaakt.')) {
      try {
        await finalizeNotulen(notulenId, 'Goedgekeurd door bestuur');
        alert('Notulen gefinaliseerd! ✅');
      } catch (error) {
        console.error('Finalize failed:', error);
      }
    }
  };

  if (loading) return <div>Laden...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!notulen) return <div>Geen notulen gevonden</div>;

  const isFinalized = notulen.status === 'finalized';

  return (
    <div className="notulen-editor">
      <h1>
        {formData.titel}
        <span className="badge">{notulen.status}</span>
        <span className="version">v{notulen.versie}</span>
      </h1>

      {/* Real-time indicator */}
      <div className="realtime-indicator">
        🟢 Live updates actief
      </div>

      <form>
        <input
          type="text"
          value={formData.titel || ''}
          onChange={(e) => setFormData({ ...formData, titel: e.target.value })}
          disabled={isFinalized}
          placeholder="Titel"
        />

        <input
          type="text"
          value={formData.locatie || ''}
          onChange={(e) => setFormData({ ...formData, locatie: e.target.value })}
          disabled={isFinalized}
          placeholder="Locatie"
        />

        {/* Participants Section */}
        <div className="participants">
          <h3>Aanwezigen</h3>
          <div>
            <label>Geregistreerde gebruikers (UUIDs):</label>
            {/* TODO: Multi-select component voor users */}
            <pre>{JSON.stringify(formData.aanwezigen_gebruikers, null, 2)}</pre>
          </div>
          
          <div>
            <label>Gasten:</label>
            <input
              type="text"
              placeholder="Naam1, Naam2, Naam3"
              value={formData.aanwezigen_gasten?.join(', ') || ''}
              onChange={(e) => 
                setFormData({ 
                  ...formData, 
                  aanwezigen_gasten: e.target.value.split(',').map(s => s.trim()).filter(Boolean)
                })
              }
              disabled={isFinalized}
            />
          </div>
        </div>

        {/* Agenda Items */}
        <div className="agenda-items">
          <h3>Agenda Items</h3>
          {formData.agenda_items?.map((item, index) => (
            <div key={index} className="agenda-item">
              <input
                type="text"
                value={item.title}
                onChange={(e) => {
                  const newItems = [...(formData.agenda_items || [])];
                  newItems[index] = { ...item, title: e.target.value };
                  setFormData({ ...formData, agenda_items: newItems });
                }}
                disabled={isFinalized}
              />
            </div>
          ))}
        </div>

        {/* Action Buttons */}
        {!isFinalized && (
          <div className="actions">
            <button type="button" onClick={handleSave}>
              💾 Opslaan
            </button>
            <button type="button" onClick={handleFinalize}>
              ✅ Finaliseren
            </button>
          </div>
        )}
      </form>

      {/* Version History */}
      <div className="version-info">
        <p>
          Gemaakt door: {notulen.created_by_name} op {new Date(notulen.created_at).toLocaleString('nl-NL')}
        </p>
        {notulen.updated_by_name && (
          <p>
            Laatst bewerkt door: {notulen.updated_by_name} op {new Date(notulen.updated_at).toLocaleString('nl-NL')}
          </p>
        )}
        {notulen.finalized_by_name && (
          <p>
            Gefinaliseerd door: {notulen.finalized_by_name} op {new Date(notulen.finalized_at!).toLocaleString('nl-NL')}
          </p>
        )}
      </div>
    </div>
  );
}
```

### Vue 3 Composable: useNotulen

```typescript
// composables/useNotulen.ts

import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { NotulenWebSocketService } from '@/services/notulenWebSocket';
import type { NotulenResponse, NotulenCreateRequest } from '@/types/notulen';

export function useNotulen(notulenId?: string) {
  const notulen = ref<NotulenResponse | null>(null);
  const notulenList = ref<NotulenResponse[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const wsService = new NotulenWebSocketService(
    import.meta.env.VITE_WS_URL,
    () => localStorage.getItem('auth_token')
  );

  const fetchNotulen = async (id: string) => {
    loading.value = true;
    error.value = null;

    try {
      const token = localStorage.getItem('auth_token');
      const response = await axios.get<NotulenResponse>(
        `${import.meta.env.VITE_API_URL}/api/notulen/${id}`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );

      notulen.value = response.data;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch notulen';
    } finally {
      loading.value = false;
    }
  };

  const updateNotulen = async (id: string, data: Partial<NotulenCreateRequest>) => {
    loading.value = true;
    error.value = null;

    try {
      const token = localStorage.getItem('auth_token');
      const response = await axios.put<NotulenResponse>(
        `${import.meta.env.VITE_API_URL}/api/notulen/${id}`,
        data,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update notulen';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  // Setup WebSocket
  onMounted(() => {
    wsService.connect(notulenId).catch(console.error);

    // Handle updates
    wsService.on('notulen_updated', (event: any) => {
      if (event.data) {
        notulen.value = event.data;
      }
    });

    // Handle finalization
    wsService.on('notulen_finalized', (event: any) => {
      if (notulen.value && event.notulenId === notulen.value.id) {
        notulen.value.status = 'finalized';
      }
    });
  });

  onUnmounted(() => {
    wsService.disconnect();
  });

  return {
    notulen,
    notulenList,
    loading,
    error,
    fetchNotulen,
    updateNotulen,
  };
}
```

### Axios API Client

```typescript
// api/notulenApi.ts

import axios, { AxiosInstance } from 'axios';
import { 
  NotulenResponse, 
  NotulenListResponse, 
  NotulenCreateRequest,
  NotulenUpdateRequest 
} from '@/types/notulen';

export class NotulenApiClient {
  private client: AxiosInstance;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.client = axios.create({
      baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add auth interceptor
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  // Create notulen
  async create(data: NotulenCreateRequest): Promise<NotulenResponse> {
    const response = await this.client.post<NotulenResponse>('/api/notulen', data);
    return response.data;
  }

  // Get notulen by ID
  async getById(id: string, asMarkdown = false): Promise<NotulenResponse | string> {
    const params = asMarkdown ? { format: 'markdown' } : {};
    const response = await this.client.get<NotulenResponse | string>(
      `/api/notulen/${id}`,
      { params }
    );
    return response.data;
  }

  // List notulen
  async list(params?: {
    status?: 'draft' | 'finalized' | 'archived';
    date_from?: string;
    date_to?: string;
    limit?: number;
    offset?: number;
  }): Promise<NotulenListResponse> {
    const response = await this.client.get<NotulenListResponse>('/api/notulen', { params });
    return response.data;
  }

  // Update notulen
  async update(id: string, data: NotulenUpdateRequest): Promise<NotulenResponse> {
    const response = await this.client.put<NotulenResponse>(`/api/notulen/${id}`, data);
    return response.data;
  }

  // Finalize notulen
  async finalize(id: string, reason?: string): Promise<void> {
    await this.client.post(`/api/notulen/${id}/finalize`, {
      wijziging_reden: reason
    });
  }

  // Archive notulen
  async archive(id: string): Promise<void> {
    await this.client.post(`/api/notulen/${id}/archive`);
  }

  // Delete notulen
  async delete(id: string): Promise<void> {
    await this.client.delete(`/api/notulen/${id}`);
  }

  // Search notulen
  async search(query: string, params?: {
    status?: string;
    date_from?: string;
    date_to?: string;
    limit?: number;
    offset?: number;
  }): Promise<NotulenListResponse> {
    const response = await this.client.get<NotulenListResponse>(
      '/api/notulen/search',
      { params: { q: query, ...params } }
    );
    return response.data;
  }

  // Get version history
  async getVersions(id: string): Promise<{ versions: any[] }> {
    const response = await this.client.get(`/api/notulen/${id}/versions`);
    return response.data;
  }

  // Get specific version
  async getVersion(id: string, version: number): Promise<any> {
    const response = await this.client.get(`/api/notulen/${id}/versions/${version}`);
    return response.data;
  }

  // List public notulen (no auth required)
  async listPublic(params?: {
    date_from?: string;
    date_to?: string;
    limit?: number;
  }): Promise<NotulenListResponse> {
    const response = await this.client.get<NotulenListResponse>(
      '/api/notulen/public',
      { params }
    );
    return response.data;
  }
}

// Export singleton instance
export const notulenApi = new NotulenApiClient(
  process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
);
```

---

## ⚡ Real-time Update Flow

### Auto-save met Debounce + WebSocket

```typescript
// hooks/useNotulenAutoSave.ts

import { useEffect, useRef, useCallback } from 'react';
import { useNotulen } from './useNotulen';
import { NotulenUpdateRequest } from '@/types/notulen';

export function useNotulenAutoSave(
  notulenId: string,
  formData: Partial<NotulenUpdateRequest>,
  debounceMs = 2000
) {
  const { updateNotulen } = useNotulen();
  const timeoutRef = useRef<NodeJS.Timeout>();
  const lastSavedRef = useRef<string>('');

  const save = useCallback(async () => {
    const currentData = JSON.stringify(formData);
    
    // Skip if no changes
    if (currentData === lastSavedRef.current) {
      return;
    }

    try {
      await updateNotulen(notulenId, formData);
      lastSavedRef.current = currentData;
      console.log('✅ Auto-saved at', new Date().toLocaleTimeString());
    } catch (error) {
      console.error('❌ Auto-save failed:', error);
    }
  }, [notulenId, formData, updateNotulen]);

  useEffect(() => {
    // Clear existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // Set new timeout
    timeoutRef.current = setTimeout(save, debounceMs);

    // Cleanup
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [formData, save, debounceMs]);

  return { save };
}
```

### Collaborative Editing Indicator

```typescript
// components/CollaborativeIndicator.tsx

import React, { useEffect, useState } from 'react';
import { NotulenWebSocketEvent } from '@/types/notulen';

export function CollaborativeIndicator({ 
  wsService,
  currentNotulenId 
}: {
  wsService: any;
  currentNotulenId: string;
}) {
  const [activeEditors, setActiveEditors] = useState<string[]>([]);

  useEffect(() => {
    const handleUpdate = (event: NotulenWebSocketEvent) => {
      if (event.notulenId === currentNotulenId && event.userId) {
        // Track who's editing
        setActiveEditors(prev => {
          if (!prev.includes(event.userId!)) {
            return [...prev, event.userId!];
          }
          return prev;
        });

        // Remove after 30 seconds of inactivity
        setTimeout(() => {
          setActiveEditors(prev => prev.filter(id => id !== event.userId));
        }, 30000);
      }
    };

    wsService.on('notulen_updated', handleUpdate);

    return () => {
      wsService.off('notulen_updated', handleUpdate);
    };
  }, [wsService, currentNotulenId]);

  if (activeEditors.length === 0) return null;

  return (
    <div className="collaborative-indicator">
      <span className="pulse-dot"></span>
      {activeEditors.length} {activeEditors.length === 1 ? 'gebruiker' : 'gebruikers'} aan het bewerken
    </div>
  );
}
```

---

## 🎨 UI Components

### Notulen List Component

```typescript
// components/NotulenList.tsx

import React, { useEffect } from 'react';
import { useNotulen } from '@/hooks/useNotulen';

export function NotulenList() {
  const { notulenList, loading, error, fetchNotulenList } = useNotulen();

  useEffect(() => {
    fetchNotulenList({ status: 'draft', limit: 20 });
  }, [fetchNotulenList]);

  if (loading) return <div>Laden...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="notulen-list">
      <h2>Notulen Overzicht</h2>
      
      <div className="filters">
        <button onClick={() => fetchNotulenList({ status: 'draft' })}>
          📝 Concepten
        </button>
        <button onClick={() => fetchNotulenList({ status: 'finalized' })}>
          ✅ Gefinaliseerd
        </button>
        <button onClick={() => fetchNotulenList({ status: 'archived' })}>
          📦 Gearchiveerd
        </button>
      </div>

      <ul>
        {notulenList.map(notulen => (
          <li key={notulen.id}>
            <a href={`/notulen/${notulen.id}`}>
              <h3>{notulen.titel}</h3>
              <p>
                {new Date(notulen.vergadering_datum).toLocaleDateString('nl-NL')}
                {' • '}
                <span className={`status ${notulen.status}`}>
                  {notulen.status}
                </span>
                {' • '}
                v{notulen.versie}
              </p>
              <p className="meta">
                Door {notulen.created_by_name} op {new Date(notulen.created_at).toLocaleString('nl-NL')}
              </p>
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

---

## 🔐 Authentication & Permissions

```typescript
// utils/auth.ts

export function hasNotulenPermission(
  permissions: string[],
  action: 'read' | 'write' | 'delete' | 'finalize' | 'archive'
): boolean {
  return permissions.some(p => 
    p === `notulen:${action}` || 
    p === 'notulen:*' || 
    p === '*:*'
  );
}

// Usage in components
function NotulenEditor() {
  const permissions = usePermissions(); // Get from auth context
  
  const canEdit = hasNotulenPermission(permissions, 'write');
  const canFinalize = hasNotulenPermission(permissions, 'finalize');
  const canDelete = hasNotulenPermission(permissions, 'delete');

  return (
    <div>
      {canEdit && <button>Bewerken</button>}
      {canFinalize && <button>Finaliseren</button>}
      {canDelete && <button>Verwijderen</button>}
    </div>
  );
}
```

---

## ⚠️ Error Handling

```typescript
// utils/notulenErrors.ts

export class NotulenError extends Error {
  constructor(
    message: string,
    public code: string,
    public statusCode?: number
  ) {
    super(message);
    this.name = 'NotulenError';
  }
}

export function handleNotulenError(error: any): NotulenError {
  if (error.response) {
    const { status, data } = error.response;
    
    switch (status) {
      case 401:
        return new NotulenError('Niet geauthenticeerd', 'UNAUTHORIZED', 401);
      case 403:
        return new NotulenError('Geen toestemming', 'FORBIDDEN', 403);
      case 404:
        return new NotulenError('Notulen niet gevonden', 'NOT_FOUND', 404);
      case 400:
        return new NotulenError(
          data.error || 'Ongeldige data',
          'VALIDATION_ERROR',
          400
        );
      default:
        return new NotulenError(
          data.error || 'Server fout',
          'SERVER_ERROR',
          status
        );
    }
  }

  return new NotulenError('Netwerk fout', 'NETWORK_ERROR');
}

// Usage
try {
  await notulenApi.create(data);
} catch (error) {
  const notulenError = handleNotulenError(error);
  
  if (notulenError.code === 'UNAUTHORIZED') {
    // Redirect to login
    router.push('/login');
  } else if (notulenError.code === 'VALIDATION_ERROR') {
    // Show validation errors
    setFormErrors(notulenError.message);
  } else {
    // Show generic error
    toast.error(notulenError.message);
  }
}
```

---

## 🧪 Testing Examples

### Jest Unit Tests

```typescript
// __tests__/notulenApi.test.ts

import { notulenApi } from '@/api/notulenApi';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('NotulenApiClient', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should create notulen', async () => {
    const mockResponse = {
      data: {
        id: 'uuid',
        titel: 'Test Notulen',
        status: 'draft',
        versie: 1
      }
    };

    mockedAxios.post.mockResolvedValue(mockResponse);

    const result = await notulenApi.create({
      titel: 'Test Notulen',
      vergadering_datum: '2025-11-04'
    });

    expect(result.titel).toBe('Test Notulen');
    expect(mockedAxios.post).toHaveBeenCalledWith(
      '/api/notulen',
      expect.objectContaining({ titel: 'Test Notulen' })
    );
  });

  it('should handle errors', async () => {
    mockedAxios.post.mockRejectedValue({
      response: {
        status: 400,
        data: { error: 'Invalid data' }
      }
    });

    await expect(
      notulenApi.create({ titel: '', vergadering_datum: '2025-11-04' })
    ).rejects.toThrow();
  });
});
```

---

## 📱 Mobile Integration (React Native)

```typescript
// services/NotulenService.native.ts

import { NotulenResponse } from '@/types/notulen';

export class NotulenServiceNative {
  private ws: WebSocket | null = null;
  
  async connect(token: string, notulenId?: string): Promise<void> {
    const url = `ws://your-api.com/api/ws/notulen?token=${token}${notulenId ? `&notulen_id=${notulenId}` : ''}`;
    
    this.ws = new WebSocket(url);
    
    this.ws.onopen = () => {
      console.log('Connected to notulen updates');
    };
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.handleNotulenUpdate(data);
    };
  }
  
  private handleNotulenUpdate(event: any): void {
    if (event.type === 'notulen_updated') {
      // Update local state/cache
      this.updateLocalNotulen(event.data);
    }
  }
  
  private updateLocalNotulen(notulen: NotulenResponse): void {
    // Implement your state management update logic
    // e.g., Redux, Zustand, React Query, etc.
  }
}
```

---

## 🎯 Best Practices

### ✅ DO's

- ✅ Gebruik TypeScript voor type safety
- ✅ Implementeer auto-save met debounce
- ✅ Handle WebSocket reconnections
- ✅ Cache notulen data lokaal
- ✅ Show loading states tijdens operations
- ✅ Validate data voor verzenden
- ✅ Optimistic UI updates bij mutaties

### ❌ DON'Ts

- ❌ Bewaar sensitive data in state zonder encryptie
- ❌ Negeer WebSocket disconnects
- ❌ Update finalized notulen (server rejected dit)
- ❌ Fetch zonder error handling
- ❌ Skip permission checks in UI
- ❌ Ignore version conflicts

---

## 🔄 Complete Workflow Voorbeeld

```typescript
// Example: Notulen management workflow

import { useState } from 'react';
import { notulenApi } from '@/api/notulenApi';
import { NotulenWebSocketService } from '@/services/notulenWebSocket';

export function NotulenWorkflow() {
  const [notulenId, setNotulenId] = useState<string | null>(null);
  const [wsService] = useState(() => 
    new NotulenWebSocketService(
      process.env.NEXT_PUBLIC_WS_URL!,
      () => localStorage.getItem('auth_token')
    )
  );

  // Step 1: Create draft notulen
  const createDraft = async () => {
    const draft = await notulenApi.create({
      titel: 'Bestuursvergadering November 2025',
      vergadering_datum: '2025-11-04',
      locatie: 'Teams',
      voorzitter: 'Salih',
      aanwezigen_gebruikers: ['uuid1', 'uuid2'],
      aanwezigen_gasten: ['External advisor'],
    });

    setNotulenId(draft.id);
    
    // Step 2: Connect WebSocket for real-time updates
    await wsService.connect(draft.id);
    
    // Step 3: Listen for updates from other users
    wsService.on('notulen_updated', (event) => {
      console.log('Someone else updated the notulen!', event.data);
      // Update UI automatically
    });

    return draft;
  };

  // Step 4: Auto-save during editing
  const handleEdit = async (changes: Partial<NotulenUpdateRequest>) => {
    if (!notulenId) return;
    
    // Optimistic update
    // ... update local state
    
    try {
      const updated = await notulenApi.update(notulenId, changes);
      // Confirm update succeeded
      console.log('✅ Updated successfully');
    } catch (error) {
      // Rollback optimistic update
      console.error('❌ Update failed, rolling back');
    }
  };

  // Step 5: Finalize when done
  const handleFinalize = async () => {
    if (!notulenId) return;

    try {
      await notulenApi.finalize(notulenId, 'Goedgekeurd door bestuur');
      console.log('✅ Notulen finalized');
      
      // WebSocket will broadcast this to all connected clients
      // Other users will see status change in real-time
    } catch (error) {
      console.error('❌ Finalize failed');
    }
  };

  return (
    <div>
      <button onClick={createDraft}>1️⃣ Create Draft</button>
      <button onClick={() => handleEdit({ titel: 'New title' })}>2️⃣ Edit</button>
      <button onClick={handleFinalize}>3️⃣ Finalize</button>
    </div>
  );
}
```

---

## 📊 State Management Integratie

### Redux Toolkit

```typescript
// store/notulenSlice.ts

import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { notulenApi } from '@/api/notulenApi';
import { NotulenResponse } from '@/types/notulen';

interface NotulenState {
  current: NotulenResponse | null;
  list: NotulenResponse[];
  loading: boolean;
  error: string | null;
}

const initialState: NotulenState = {
  current: null,
  list: [],
  loading: false,
  error: null,
};

export const fetchNotulenList = createAsyncThunk(
  'notulen/fetchList',
  async (params?: any) => {
    const response = await notulenApi.list(params);
    return response.notulen;
  }
);

export const updateNotulen = createAsyncThunk(
  'notulen/update',
  async ({ id, data }: { id: string; data: any }) => {
    return await notulenApi.update(id, data);
  }
);

const notulenSlice = createSlice({
  name: 'notulen',
  initialState,
  reducers: {
    // WebSocket real-time updates
    notulenUpdatedViaWebSocket: (state, action: PayloadAction<NotulenResponse>) => {
      state.current = action.payload;
      
      // Update in list if present
      const index = state.list.findIndex(n => n.id === action.payload.id);
      if (index !== -1) {
        state.list[index] = action.payload;
      }
    },
    
    notulenFinalizedViaWebSocket: (state, action: PayloadAction<string>) => {
      if (state.current?.id === action.payload) {
        state.current.status = 'finalized';
      }
      
      const index = state.list.findIndex(n => n.id === action.payload);
      if (index !== -1) {
        state.list[index].status = 'finalized';
      }
    },
    
    notulenDeletedViaWebSocket: (state, action: PayloadAction<string>) => {
      if (state.current?.id === action.payload) {
        state.current = null;
      }
      
      state.list = state.list.filter(n => n.id !== action.payload);
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchNotulenList.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchNotulenList.fulfilled, (state, action) => {
        state.list = action.payload;
        state.loading = false;
      })
      .addCase(fetchNotulenList.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to fetch';
        state.loading = false;
      });
  },
});

export const {
  notulenUpdatedViaWebSocket,
  notulenFinalizedViaWebSocket,
  notulenDeletedViaWebSocket,
} = notulenSlice.actions;

export default notulenSlice.reducer;
```

### React Query

```typescript
// hooks/useNotulenQuery.ts

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notulenApi } from '@/api/notulenApi';

export function useNotulenQuery(id: string) {
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ['notulen', id],
    queryFn: () => notulenApi.getById(id),
    staleTime: 30000, // 30 seconds
  });

  const updateMutation = useMutation({
    mutationFn: (data: any) => notulenApi.update(id, data),
    onSuccess: (updated) => {
      // Update cache
      queryClient.setQueryData(['notulen', id], updated);
      
      // Invalidate list
      queryClient.invalidateQueries({ queryKey: ['notulen', 'list'] });
    },
  });

  const finalizeMutation = useMutation({
    mutationFn: (reason?: string) => notulenApi.finalize(id, reason),
    onSuccess: () => {
      // Refetch to get updated status
      queryClient.invalidateQueries({ queryKey: ['notulen', id] });
    },
  });

  return {
    notulen: data,
    isLoading,
    error,
    update: updateMutation.mutate,
    finalize: finalizeMutation.mutate,
  };
}
```

---

## 🌍 Internationalization

```typescript
// i18n/nl/notulen.json

{
  "notulen": {
    "title": "Notulen",
    "create": "Nieuwe notulen",
    "edit": "Bewerken",
    "delete": "Verwijderen",
    "finalize": "Finaliseren",
    "archive": "Archiveren",
    
    "status": {
      "draft": "Concept",
      "finalized": "Gefinaliseerd",
      "archived": "Gearchiveerd"
    },
    
    "fields": {
      "titel": "Titel",
      "vergadering_datum": "Vergaderdatum",
      "locatie": "Locatie",
      "voorzitter": "Voorzitter",
      "notulist": "Notulist",
      "aanwezigen": "Aanwezigen",
      "afwezigen": "Afwezigen"
    },
    
    "messages": {
      "created": "Notulen succesvol aangemaakt",
      "updated": "Notulen bijgewerkt",
      "finalized": "Notulen gefinaliseerd",
      "deleted": "Notulen verwijderd",
      "error": "Er is een fout opgetreden"
    }
  }
}
```

---

## 📈 Performance Tips

### 1. Lazy Loading

```typescript
import dynamic from 'next/dynamic';

const NotulenEditor = dynamic(() => import('@/components/NotulenEditor'), {
  loading: () => <div>Laden...</div>,
  ssr: false, // Disable SSR voor WebSocket components
});
```

### 2. Pagination & Infinite Scroll

```typescript
import { useInfiniteQuery } from '@tanstack/react-query';

export function useInfiniteNotulen() {
  return useInfiniteQuery({
    queryKey: ['notulen', 'infinite'],
    queryFn: ({ pageParam = 0 }) =>
      notulenApi.list({ limit: 20, offset: pageParam }),
    getNextPageParam: (lastPage, allPages) => {
      const nextOffset = allPages.length * 20;
      return nextOffset < lastPage.total ? nextOffset : undefined;
    },
  });
}
```

### 3. Optimistic Updates

```typescript
const updateMutation = useMutation({
  mutationFn: (data: any) => notulenApi.update(id, data),
  onMutate: async (newData) => {
    // Cancel outgoing refetches
    await queryClient.cancelQueries({ queryKey: ['notulen', id] });

    // Snapshot previous value
    const previous = queryClient.getQueryData(['notulen', id]);

    // Optimistically update
    queryClient.setQueryData(['notulen', id], (old: any) => ({
      ...old,
      ...newData
    }));

    return { previous };
  },
  onError: (err, newData, context) => {
    // Rollback on error
    queryClient.setQueryData(['notulen', id], context?.previous);
  },
});
```

---

## 🔧 Troubleshooting

### WebSocket niet verbinden?

```typescript
// Debug WebSocket connection
const wsService = new NotulenWebSocketService(WS_URL, getToken);

wsService.connect(notulenId)
  .then(() => console.log('✅ Connected'))
  .catch((error) => {
    console.error('❌ Connection failed:', error);
    
    // Check:
    // 1. Is token valid?
    // 2. Is backend running?
    // 3. Is CORS configured correctly?
    // 4. Check browser console for errors
  });
```

### Updates niet ontvangen?

```typescript
// Verify WebSocket is listening
wsService.on('*', (event) => {
  console.log('Received event:', event.type, event);
});

// Send test ping
wsService.ping();
```

### Permission errors?

```typescript
// Check user permissions
const checkPermissions = async () => {
  const token = localStorage.getItem('auth_token');
  const response = await axios.get('/api/auth/profile', {
    headers: { Authorization: `Bearer ${token}` }
  });
  
  console.log('User permissions:', response.data.permissions);
  // Should include: notulen:read, notulen:write, etc.
};
```

---

## 📞 Support

Voor vragen:
1. Check deze documentatie
2. Review [`notulen-system-documentation.md`](../notulen-system-documentation.md)
3. Check backend logs voor errors
4. Test endpoints met Postman/Thunder Client

---

**Laatst Bijgewerkt**: November 2025  
**Versie**: 1.0  
**Status**: Production Ready ✅