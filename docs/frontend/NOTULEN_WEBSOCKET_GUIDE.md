# Notulen WebSocket Real-time Updates - Frontend Integratie Guide

Complete guide voor het implementeren van real-time notulen updates via WebSocket.

---

## 🎯 Wat Je Krijgt

Met WebSocket real-time updates zie je **instant** wanneer:
- ✅ Iemand anders de notulen bewerkt
- ✅ Notulen worden gefinaliseerd
- ✅ Notulen worden gearchiveerd  
- ✅ Notulen worden verwijderd

**Geen polling nodig!** Alle wijzigingen komen automatisch binnen.

---

## 🚀 Quick Start (5 minuten)

### Stap 1: Setup WebSocket Service

```typescript
// services/notulenWebSocket.ts

export class NotulenWebSocketService {
  private ws: WebSocket | null = null;

  constructor(
    private wsUrl: string,
    private getToken: () => string | null
  ) {}

  connect(notulenId?: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const token = this.getToken();
      let url = `${this.wsUrl}/api/ws/notulen`;
      
      // Add query parameters
      const params = new URLSearchParams();
      if (token) params.append('token', token);
      if (notulenId) params.append('notulen_id', notulenId);
      
      if (params.toString()) url += `?${params.toString()}`;

      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log('🟢 Notulen WebSocket connected');
        resolve();
      };

      this.ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        this.handleEvent(data);
      };

      this.ws.onerror = (error) => {
        console.error('🔴 WebSocket error:', error);
        reject(error);
      };

      this.ws.onclose = () => {
        console.log('🟡 WebSocket disconnected');
        this.reconnect();
      };
    });
  }

  private handleEvent(event: any): void {
    console.log('📨 WebSocket event:', event.type);
    // Emit to listeners (implement your event system)
    window.dispatchEvent(new CustomEvent('notulen-ws', { detail: event }));
  }

  private reconnect(): void {
    setTimeout(() => {
      console.log('🔄 Reconnecting...');
      this.connect().catch(console.error);
    }, 3000);
  }

  disconnect(): void {
    this.ws?.close();
  }
}
```

### Stap 2: React Component

```typescript
// components/NotulenLiveEditor.tsx

import { useEffect, useState } from 'react';

export function NotulenLiveEditor({ notulenId }: { notulenId: string }) {
  const [notulen, setNotulen] = useState<any>(null);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    // 1. Connect WebSocket
    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080';
    const token = localStorage.getItem('auth_token');
    
    const ws = new WebSocket(
      `${wsUrl}/api/ws/notulen?token=${token}&notulen_id=${notulenId}`
    );

    ws.onopen = () => {
      console.log('✅ WebSocket connected');
      setWsConnected(true);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      switch (data.type) {
        case 'welcome':
          console.log('👋 Welcome message:', data.message);
          break;
          
        case 'notulen_updated':
          console.log('📝 Notulen updated by:', data.userId);
          setNotulen(data.data); // Update state with new data
          break;
          
        case 'notulen_finalized':
          console.log('✅ Notulen finalized');
          setNotulen(prev => ({ ...prev, status: 'finalized' }));
          break;
          
        case 'notulen_deleted':
          console.log('🗑️ Notulen deleted');
          // Redirect user weg van deleted document
          window.location.href = '/notulen';
          break;
      }
    };

    ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
      setWsConnected(false);
    };

    ws.onclose = () => {
      console.log('🔌 WebSocket closed');
      setWsConnected(false);
    };

    // Cleanup
    return () => {
      ws.close();
    };
  }, [notulenId]);

  return (
    <div>
      {/* Status indicator */}
      <div className={`ws-indicator ${wsConnected ? 'connected' : 'disconnected'}`}>
        {wsConnected ? '🟢 Live' : '🔴 Offline'}
      </div>

      {/* Editor */}
      {notulen && (
        <div>
          <h1>{notulen.titel}</h1>
          <p>Status: {notulen.status}</p>
          <p>Versie: {notulen.versie}</p>
        </div>
      )}
    </div>
  );
}
```

---

## 📡 WebSocket Event Reference

### Inkomende Events (Server → Client)

#### 1. `welcome`
**Wanneer**: Direct na connectie  
**Data**:
```json
{
  "type": "welcome",
  "message": "Connected to NotulenHub! You will receive real-time updates...",
  "user_id": "your-uuid",
  "notulen_id": "notulen-uuid",
  "timestamp": 1699123456
}
```

**Gebruik**: Bevestig connectie, toon notification

#### 2. `notulen_updated`
**Wanneer**: Iemand bewerkt de notulen  
**Data**:
```json
{
  "type": "notulen_updated",
  "notulenId": "uuid",
  "userId": "uuid",
  "data": {
    "id": "uuid",
    "titel": "Updated Titel",
    "vergadering_datum": "2025-11-04",
    "status": "draft",
    "versie": 2,
    "created_by_name": "Jeffrey Admin",
    "updated_by_name": "Salih Bestuur",
    // ... volledige notulen data
  },
  "timestamp": "2025-11-04T23:00:00Z"
}
```

**Gebruik**: Update UI met nieuwe data

#### 3. `notulen_finalized`
**Wanneer**: Notulen worden gefinaliseerd  
**Data**:
```json
{
  "type": "notulen_finalized",
  "notulenId": "uuid",
  "userId": "uuid",
  "timestamp": "2025-11-04T23:00:00Z"
}
```

**Gebruik**: Update status naar 'finalized', disable editing

#### 4. `notulen_archived`
**Wanneer**: Notulen worden gearchiveerd  
**Data**:
```json
{
  "type": "notulen_archived",
  "notulenId": "uuid",
  "userId": "uuid",
  "timestamp": "2025-11-04T23:00:00Z"
}
```

**Gebruik**: Update status naar 'archived', verberg van active list

#### 5. `notulen_deleted`
**Wanneer**: Notulen worden verwijderd  
**Data**:
```json
{
  "type": "notulen_deleted",
  "notulenId": "uuid",
  "userId": "uuid",
  "timestamp": "2025-11-04T23:00:00Z"
}
```

**Gebruik**: Remove van UI, redirect als viewing deleted item

#### 6. `ping`
**Wanneer**: Server heartbeat (elke 30 seconden)  
**Data**:
```json
{
  "type": "ping",
  "timestamp": 1699123456
}
```

**Gebruik**: Respond met `pong`, check connection health

### Uitgaande Events (Client → Server)

#### 1. `ping`
```json
{
  "type": "ping"
}
```

**Response**: Server stuurt `pong` terug

#### 2. `subscribe` (Toekomstige feature)
```json
{
  "type": "subscribe",
  "notulen_id": "uuid"
}
```

#### 3. `unsubscribe` (Toekomstige feature)
```json
{
  "type": "unsubscribe",
  "notulen_id": "uuid"
}
```

---

## 🏗️ Implementatie Voorbeelden

### React Context Provider

```typescript
// contexts/NotulenWebSocketContext.tsx

import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { NotulenWebSocketService } from '@/services/notulenWebSocket';
import { NotulenWebSocketEvent } from '@/types/notulen';

interface NotulenWSContextType {
  connected: boolean;
  lastEvent: NotulenWebSocketEvent | null;
  connect: (notulenId?: string) => void;
  disconnect: () => void;
}

const NotulenWSContext = createContext<NotulenWSContextType | null>(null);

export function NotulenWebSocketProvider({ children }: { children: React.ReactNode }) {
  const [connected, setConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState<NotulenWebSocketEvent | null>(null);
  const [wsService] = useState(() => 
    new NotulenWebSocketService(
      process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080',
      () => localStorage.getItem('auth_token')
    )
  );

  const connect = useCallback((notulenId?: string) => {
    wsService.connect(notulenId)
      .then(() => setConnected(true))
      .catch(() => setConnected(false));
  }, [wsService]);

  const disconnect = useCallback(() => {
    wsService.disconnect();
    setConnected(false);
  }, [wsService]);

  useEffect(() => {
    // Listen to all events
    wsService.on('*', (event: NotulenWebSocketEvent) => {
      setLastEvent(event);
    });

    return () => {
      wsService.disconnect();
    };
  }, [wsService]);

  return (
    <NotulenWSContext.Provider value={{ connected, lastEvent, connect, disconnect }}>
      {children}
    </NotulenWSContext.Provider>
  );
}

export function useNotulenWebSocket() {
  const context = useContext(NotulenWSContext);
  if (!context) {
    throw new Error('useNotulenWebSocket must be used within NotulenWebSocketProvider');
  }
  return context;
}
```

### Usage in App

```typescript
// pages/_app.tsx or layout.tsx

import { NotulenWebSocketProvider } from '@/contexts/NotulenWebSocketContext';

export default function App({ Component, pageProps }) {
  return (
    <NotulenWebSocketProvider>
      <Component {...pageProps} />
    </NotulenWebSocketProvider>
  );
}
```

```typescript
// pages/notulen/[id].tsx

import { useNotulenWebSocket } from '@/contexts/NotulenWebSocketContext';

export default function NotulenDetail({ id }: { id: string }) {
  const { connected, lastEvent, connect } = useNotulenWebSocket();

  useEffect(() => {
    connect(id); // Connect for this specific notulen
  }, [id, connect]);

  useEffect(() => {
    if (lastEvent?.type === 'notulen_updated') {
      // Handle update
      console.log('Notulen updated:', lastEvent.data);
    }
  }, [lastEvent]);

  return (
    <div>
      <div className="status">
        {connected ? '🟢 Live updates' : '🔴 Disconnected'}
      </div>
      {/* Rest of component */}
    </div>
  );
}
```

---

## 🎨 UI Feedback Patterns

### 1. Update Toast Notifications

```typescript
// components/NotulenUpdateToast.tsx

import { useEffect } from 'react';
import { toast } from 'react-hot-toast';
import { useNotulenWebSocket } from '@/contexts/NotulenWebSocketContext';

export function NotulenUpdateToast() {
  const { lastEvent } = useNotulenWebSocket();

  useEffect(() => {
    if (!lastEvent) return;

    switch (lastEvent.type) {
      case 'notulen_updated':
        toast('📝 Notulen bijgewerkt door andere gebruiker', {
          icon: '🔄',
          duration: 3000,
        });
        break;

      case 'notulen_finalized':
        toast.success('✅ Notulen zijn gefinaliseerd!');
        break;

      case 'notulen_archived':
        toast('📦 Notulen zijn gearchiveerd', {
          icon: '📦',
        });
        break;

      case 'notulen_deleted':
        toast.error('🗑️ Notulen zijn verwijderd!');
        break;
    }
  }, [lastEvent]);

  return null;
}
```

### 2. Live Editing Indicator

```typescript
// components/LiveEditingBadge.tsx

import { useEffect, useState } from 'react';

export function LiveEditingBadge({ wsService, notulenId }) {
  const [editors, setEditors] = useState<Set<string>>(new Set());

  useEffect(() => {
    const handleUpdate = (event) => {
      if (event.notulenId === notulenId && event.userId) {
        setEditors(prev => new Set([...prev, event.userId]));
        
        // Remove after 10 seconds of inactivity
        setTimeout(() => {
          setEditors(prev => {
            const next = new Set(prev);
            next.delete(event.userId);
            return next;
          });
        }, 10000);
      }
    };

    wsService.on('notulen_updated', handleUpdate);
    return () => wsService.off('notulen_updated', handleUpdate);
  }, [wsService, notulenId]);

  if (editors.size === 0) return null;

  return (
    <div className="live-editing-badge">
      <span className="pulse"></span>
      {editors.size} {editors.size === 1 ? 'persoon' : 'personen'} aan het bewerken
    </div>
  );
}
```

### 3. Version Change Animation

```typescript
// components/VersionBadge.tsx

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';

export function VersionBadge({ version }) {
  const [highlight, setHighlight] = useState(false);

  useEffect(() => {
    setHighlight(true);
    const timer = setTimeout(() => setHighlight(false), 2000);
    return () => clearTimeout(timer);
  }, [version]);

  return (
    <motion.div
      className="version-badge"
      animate={highlight ? { scale: [1, 1.2, 1] } : {}}
      transition={{ duration: 0.3 }}
    >
      v{version}
      {highlight && <span className="new-badge">Nieuw!</span>}
    </motion.div>
  );
}
```

---

## 🔧 Advanced Patterns

### 1. Conflict Resolution

```typescript
// hooks/useNotulenConflictResolution.ts

export function useNotulenConflictResolution(notulenId: string) {
  const [localVersion, setLocalVersion] = useState<number>(1);
  const [serverVersion, setServerVersion] = useState<number>(1);
  const [hasConflict, setHasConflict] = useState(false);

  useEffect(() => {
    const handleUpdate = (event: any) => {
      if (event.type === 'notulen_updated' && event.data) {
        const newServerVersion = event.data.versie;
        
        if (newServerVersion > localVersion) {
          setServerVersion(newServerVersion);
          
          // Check for conflict
          if (localVersion !== serverVersion) {
            setHasConflict(true);
            
            // Show conflict resolution UI
            const shouldReload = confirm(
              'Iemand anders heeft deze notulen bewerkt. Wil je de nieuwste versie laden?'
            );
            
            if (shouldReload) {
              window.location.reload();
            }
          } else {
            // No conflict, just update
            setLocalVersion(newServerVersion);
          }
        }
      }
    };

    window.addEventListener('notulen-ws', handleUpdate as any);
    return () => window.removeEventListener('notulen-ws', handleUpdate as any);
  }, [localVersion, serverVersion, notulenId]);

  return { hasConflict, localVersion, serverVersion };
}
```

### 2. Optimistic Updates met Rollback

```typescript
// hooks/useOptimisticNotulenUpdate.ts

import { useState } from 'react';
import { notulenApi } from '@/api/notulenApi';

export function useOptimisticNotulenUpdate(initialData: any) {
  const [data, setData] = useState(initialData);
  const [isSaving, setIsSaving] = useState(false);

  const update = async (changes: any) => {
    // 1. Save current state for rollback
    const previousData = data;
    
    // 2. Optimistically update UI
    setData({ ...data, ...changes });
    setIsSaving(true);

    try {
      // 3. Send to server (WebSocket will broadcast to others)
      await notulenApi.update(data.id, changes);
      
      // 4. Success - server will send updated data via WebSocket
      setIsSaving(false);
      
    } catch (error) {
      // 5. Error - rollback to previous state
      console.error('Update failed, rolling back');
      setData(previousData);
      setIsSaving(false);
      throw error;
    }
  };

  return { data, setData, update, isSaving };
}
```

### 3. Multi-user Presence Tracking

```typescript
// components/ActiveEditors.tsx

import { useEffect, useState } from 'react';

interface Editor {
  userId: string;
  userName?: string;
  lastActive: number;
}

export function ActiveEditors({ wsService, notulenId }) {
  const [editors, setEditors] = useState<Map<string, Editor>>(new Map());

  useEffect(() => {
    const handleUpdate = (event: any) => {
      if (event.notulenId === notulenId && event.userId) {
        setEditors(prev => {
          const next = new Map(prev);
          next.set(event.userId, {
            userId: event.userId,
            userName: event.userName,
            lastActive: Date.now(),
          });
          return next;
        });
      }
    };

    wsService.on('notulen_updated', handleUpdate);

    // Cleanup stale editors elke 5 seconden
    const interval = setInterval(() => {
      const now = Date.now();
      setEditors(prev => {
        const next = new Map(prev);
        for (const [userId, editor] of next) {
          if (now - editor.lastActive > 30000) { // 30 seconds timeout
            next.delete(userId);
          }
        }
        return next;
      });
    }, 5000);

    return () => {
      wsService.off('notulen_updated', handleUpdate);
      clearInterval(interval);
    };
  }, [wsService, notulenId]);

  const editorsList = Array.from(editors.values());

  if (editorsList.length === 0) return null;

  return (
    <div className="active-editors">
      <h4>Actieve editors:</h4>
      <ul>
        {editorsList.map(editor => (
          <li key={editor.userId}>
            <span className="user-avatar"></span>
            {editor.userName || 'Anonymous'}
            <span className="active-dot"></span>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

---

## 🎭 Connection States

### State Machine

```typescript
// hooks/useWebSocketState.ts

type WSState = 'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'error';

export function useWebSocketState() {
  const [state, setState] = useState<WSState>('disconnected');
  const [error, setError] = useState<string | null>(null);

  const connect = async (wsService: any) => {
    setState('connecting');
    setError(null);

    try {
      await wsService.connect();
      setState('connected');
    } catch (err: any) {
      setState('error');
      setError(err.message);
    }
  };

  const reconnect = async (wsService: any) => {
    setState('reconnecting');
    await connect(wsService);
  };

  return { state, error, connect, reconnect };
}
```

### Visual Indicators

```typescript
// components/WebSocketStatus.tsx

export function WebSocketStatus({ state }: { state: WSState }) {
  const indicators = {
    disconnected: { icon: '⚫', text: 'Niet verbonden', color: 'gray' },
    connecting: { icon: '🟡', text: 'Verbinden...', color: 'yellow' },
    connected: { icon: '🟢', text: 'Live updates', color: 'green' },
    reconnecting: { icon: '🟠', text: 'Opnieuw verbinden...', color: 'orange' },
    error: { icon: '🔴', text: 'Verbindingsfout', color: 'red' },
  };

  const indicator = indicators[state];

  return (
    <div className={`ws-status ws-status-${indicator.color}`}>
      <span className="icon">{indicator.icon}</span>
      <span className="text">{indicator.text}</span>
    </div>
  );
}
```

---

## 🧪 Testing

### Mock WebSocket voor Tests

```typescript
// __mocks__/websocket.ts

export class MockWebSocket {
  onopen: (() => void) | null = null;
  onmessage: ((event: any) => void) | null = null;
  onerror: ((error: any) => void) | null = null;
  onclose: (() => void) | null = null;

  constructor(public url: string) {
    // Simulate connection after delay
    setTimeout(() => {
      if (this.onopen) this.onopen();
    }, 100);
  }

  send(data: string): void {
    console.log('Mock WS send:', data);
  }

  close(): void {
    if (this.onclose) this.onclose();
  }

  // Test helper: simulate server message
  simulateMessage(data: any): void {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
}

// Usage in tests
global.WebSocket = MockWebSocket as any;
```

### React Testing Library

```typescript
// __tests__/NotulenEditor.test.tsx

import { render, screen, waitFor } from '@testing-library/react';
import { NotulenEditor } from '@/components/NotulenEditor';
import { MockWebSocket } from '@/__mocks__/websocket';

describe('NotulenEditor', () => {
  it('should receive WebSocket updates', async () => {
    const mockWs = new MockWebSocket('ws://test');
    global.WebSocket = jest.fn(() => mockWs) as any;

    render(<NotulenEditor notulenId="test-id" />);

    // Simulate server update
    mockWs.simulateMessage({
      type: 'notulen_updated',
      data: {
        id: 'test-id',
        titel: 'Updated by server',
        versie: 2
      }
    });

    await waitFor(() => {
      expect(screen.getByText('Updated by server')).toBeInTheDocument();
    });
  });
});
```

---

## 🚨 Error Handling & Recovery

### Reconnection Strategy

```typescript
// services/reconnectionStrategy.ts

export class ExponentialBackoffReconnection {
  private attempts = 0;
  private maxAttempts = 10;
  private baseDelay = 1000;
  private maxDelay = 30000;

  async reconnect(connectFn: () => Promise<void>): Promise<void> {
    while (this.attempts < this.maxAttempts) {
      try {
        await connectFn();
        this.attempts = 0; // Reset on success
        return;
      } catch (error) {
        this.attempts++;
        const delay = Math.min(
          this.baseDelay * Math.pow(2, this.attempts),
          this.maxDelay
        );
        
        console.log(`Reconnect attempt ${this.attempts}, waiting ${delay}ms`);
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }

    throw new Error('Max reconnection attempts reached');
  }

  reset(): void {
    this.attempts = 0;
  }
}
```

### Network Status Integration

```typescript
// hooks/useNetworkAwareWebSocket.ts

import { useEffect, useState } from 'react';

export function useNetworkAwareWebSocket(wsService: any) {
  const [isOnline, setIsOnline] = useState(navigator.onLine);

  useEffect(() => {
    const handleOnline = () => {
      console.log('🟢 Network back online, reconnecting...');
      setIsOnline(true);
      wsService.connect().catch(console.error);
    };

    const handleOffline = () => {
      console.log('🔴 Network offline');
      setIsOnline(false);
      wsService.disconnect();
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, [wsService]);

  return { isOnline };
}
```

---

## 📦 Complete Setup Voorbeeld

### Next.js 14 App Router

```typescript
// app/notulen/[id]/page.tsx

'use client';

import { useEffect, useState } from 'react';
import { NotulenWebSocketService } from '@/services/notulenWebSocket';
import { notulenApi } from '@/api/notulenApi';
import { NotulenResponse } from '@/types/notulen';

export default function NotulenPage({ params }: { params: { id: string } }) {
  const [notulen, setNotulen] = useState<NotulenResponse | null>(null);
  const [wsConnected, setWsConnected] = useState(false);

  useEffect(() => {
    // 1. Fetch initial data
    notulenApi.getById(params.id).then(setNotulen);

    // 2. Setup WebSocket
    const wsService = new NotulenWebSocketService(
      process.env.NEXT_PUBLIC_WS_URL!,
      () => localStorage.getItem('auth_token')
    );

    wsService.connect(params.id)
      .then(() => setWsConnected(true))
      .catch(console.error);

    // 3. Listen for updates
    wsService.on('notulen_updated', (event) => {
      if (event.data) {
        setNotulen(event.data);
      }
    });

    wsService.on('notulen_finalized', (event) => {
      if (event.notulenId === params.id) {
        setNotulen(prev => prev ? { ...prev, status: 'finalized' } : null);
      }
    });

    // 4. Cleanup
    return () => {
      wsService.disconnect();
    };
  }, [params.id]);

  if (!notulen) return <div>Laden...</div>;

  return (
    <div className="notulen-page">
      {/* Status bar */}
      <div className="status-bar">
        <div className={`ws-indicator ${wsConnected ? 'connected' : 'disconnected'}`}>
          {wsConnected ? '🟢 Live' : '🔴 Offline'}
        </div>
        <div className="version">v{notulen.versie}</div>
        <div className="status">{notulen.status}</div>
      </div>

      {/* Content */}
      <h1>{notulen.titel}</h1>
      <p><strong>Datum:</strong> {new Date(notulen.vergadering_datum).toLocaleDateString('nl-NL')}</p>
      <p><strong>Locatie:</strong> {notulen.locatie}</p>
      
      {/* Participants */}
      <div className="participants">
        <h3>Aanwezigen</h3>
        <ul>
          {notulen.aanwezigen?.map((name, i) => (
            <li key={i}>{name}</li>
          ))}
        </ul>
      </div>

      {/* Agenda Items */}
      <div className="agenda">
        <h3>Agenda</h3>
        {notulen.agenda_items?.map((item, i) => (
          <div key={i} className="agenda-item">
            <h4>{item.title}</h4>
            <p>{item.details}</p>
          </div>
        ))}
      </div>

      {/* Actions */}
      {notulen.status === 'draft' && (
        <div className="actions">
          <button onClick={() => {/* edit logic */}}>
            ✏️ Bewerken
          </button>
          <button onClick={async () => {
            await notulenApi.finalize(params.id);
          }}>
            ✅ Finaliseren
          </button>
        </div>
      )}
    </div>
  );
}
```

---

## 🎬 Stap-voor-Stap Implementatie

### Week 1: Basic Setup

```typescript
// ✅ Day 1-2: Setup types en API client
import { NotulenResponse } from '@/types/notulen';
import { notulenApi } from '@/api/notulenApi';

// ✅ Day 3-4: Basic CRUD components
function NotulenList() {
  const [notulen, setNotulen] = useState([]);
  
  useEffect(() => {
    notulenApi.list().then(res => setNotulen(res.notulen));
  }, []);
  
  return <div>{/* List UI */}</div>;
}

// ✅ Day 5: Create/Edit forms
function NotulenForm() {
  const handleSubmit = async (data) => {
    await notulenApi.create(data);
  };
  
  return <form>{/* Form fields */}</form>;
}
```

### Week 2: WebSocket Integration

```typescript
// ✅ Day 1-2: WebSocket service
const wsService = new NotulenWebSocketService(WS_URL, getToken);
await wsService.connect();

// ✅ Day 3-4: Event handlers
wsService.on('notulen_updated', handleUpdate);
wsService.on('notulen_finalized', handleFinalized);

// ✅ Day 5: UI feedback
<WebSocketStatus state={wsConnected ? 'connected' : 'disconnected'} />
```

### Week 3: Polish & Testing

```typescript
// ✅ Day 1-2: Error handling
try {
  await notulenApi.update(id, data);
} catch (error) {
  handleNotulenError(error);
}

// ✅ Day 3-4: Optimistic updates
updateOptimistically(data);

// ✅ Day 5: Tests
test('should handle WebSocket updates', async () => {
  // Test implementation
});
```

---

## 💡 Tips & Tricks

### 1. Debugging WebSocket

```typescript
// Enable verbose logging
const wsService = new NotulenWebSocketService(WS_URL, getToken);

wsService.on('*', (event) => {
  console.log('📨 [WebSocket]', {
    type: event.type,
    timestamp: new Date(event.timestamp).toISOString(),
    data: event.data
  });
});
```

### 2. Performance Monitoring

```typescript
// Monitor WebSocket latency
wsService.on('notulen_updated', (event) => {
  const serverTime = new Date(event.timestamp).getTime();
  const clientTime = Date.now();
  const latency = clientTime - serverTime;
  
  console.log(`⏱️ Update latency: ${latency}ms`);
  
  if (latency > 1000) {
    console.warn('⚠️ High latency detected');
  }
});
```

### 3. Offline Support

```typescript
// Queue updates wanneer offline
class OfflineQueue {
  private queue: any[] = [];

  add(update: any): void {
    this.queue.push(update);
    localStorage.setItem('notulen-offline-queue', JSON.stringify(this.queue));
  }

  async flush(notulenApi: any): Promise<void> {
    while (this.queue.length > 0) {
      const update = this.queue.shift();
      try {
        await notulenApi.update(update.id, update.data);
      } catch (error) {
        // Re-add to queue on failure
        this.queue.unshift(update);
        throw error;
      }
    }
    localStorage.removeItem('notulen-offline-queue');
  }
}
```

---

## 🌟 Production Checklist

- [ ] ✅ WebSocket URL correct geconfigureerd
- [ ] ✅ Authentication token wordt meegestuurd
- [ ] ✅ R econnection logic geïmplementeerd
- [ ] ✅ Error handling voor alle scenarios
- [ ] ✅ Loading states voor alle operations
- [ ] ✅ Optimistic updates voor betere UX
- [ ] ✅ Network status monitoring
- [ ] ✅ Offline queue voor updates
- [ ] ✅ Tests voor kritieke flows
- [ ] ✅ Logging & monitoring
- [ ] ✅ Permission checks in UI
- [ ] ✅ Conflict resolution UI

---

## 📞 Support & Resources

**Documentatie**:
- [`NOTULEN_API_COMPLETE.md`](./NOTULEN_API_COMPLETE.md) - Complete API reference
- [`notulen-system-documentation.md`](../notulen-system-documentation.md) - Backend docs

**Backend Endpoints**:
- REST API: `http://localhost:8080/api/notulen`
- WebSocket: `ws://localhost:8080/api/ws/notulen`

**Test Credentials**:
```
Email: admin@dekoninklijkeloop.nl
Password: [zie .env file]
```

---

**Laatst Bijgewerkt**: November 2025  
**Versie**: 1.0  
**Status**: Complete Guide ✅