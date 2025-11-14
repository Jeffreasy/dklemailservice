# WebSocket Steps Integration Guide

## Overzicht

Dit document beschrijft hoe je de WebSocket verbinding voor real-time steps tracking integreert in de frontend.

## Probleem Opgelost

**Fix voor:** `WebSocket connection to 'wss://dklemailservice.onrender.com/ws/steps?user_id=public&token=' failed`

**Oorzaken:**
1. ❌ Lege token parameter (`token=`) werd niet correct afgehandeld
2. ❌ Publieke toegang was niet goed geconfigureerd
3. ❌ URL pad was inconsistent tussen frontend en backend

**Oplossing:**
1. ✅ Backend accepteert nu lege tokens als anonieme verbinding
2. ✅ Publieke toegang is volledig ondersteund
3. ✅ Beide `/ws/steps` en `/api/ws/steps` paden werken

## WebSocket Endpoints

### Publieke Toegang (Aanbevolen voor Public Website)

```javascript
// Zonder token - volledig publiek
const ws = new WebSocket('wss://dklemailservice.onrender.com/ws/steps');

// OF met user_id parameter
const ws = new WebSocket('wss://dklemailservice.onrender.com/ws/steps?user_id=public');

// Development
const ws = new WebSocket('ws://localhost:8080/ws/steps');
```

### Geauthenticeerde Toegang (Voor Admin/User Dashboard)

```javascript
const token = localStorage.getItem('access_token');
const ws = new WebSocket(`wss://dklemailservice.onrender.com/ws/steps?token=${token}`);
```

## Beschikbare Channels

### Voor Publieke Gebruikers
- `total_updates` - Totaal aantal stappen van alle deelnemers
- `leaderboard_updates` - Publieke leaderboard wijzigingen

### Voor Geauthenticeerde Gebruikers (extra channels)
- `step_updates` - Individuele stappen updates
- `badge_earned` - Persoonlijke badge notificaties

## Implementatie Voorbeelden

### 1. React/Next.js - Public Steps Counter

```typescript
// hooks/usePublicSteps.ts
import { useState, useEffect } from 'react';

interface StepsData {
  totalSteps: number;
  connected: boolean;
}

export function usePublicSteps(): StepsData {
  const [totalSteps, setTotalSteps] = useState(0);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    // Publieke WebSocket verbinding (geen authenticatie nodig)
    const ws = new WebSocket(
      process.env.NEXT_PUBLIC_WS_URL || 'wss://dklemailservice.onrender.com/ws/steps'
    );

    ws.onopen = () => {
      console.log('✅ Steps WebSocket connected (public)');
      setConnected(true);
      
      // Subscribe naar publieke channels
      ws.send(JSON.stringify({
        type: 'subscribe',
        channels: ['total_updates', 'leaderboard_updates']
      }));
    };

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        
        // Welcome message
        if (message.type === 'welcome') {
          console.log('📢', message.message);
        }
        
        // Total steps update
        if (message.type === 'total_update') {
          setTotalSteps(message.data.total_steps);
        }
      } catch (error) {
        console.error('❌ Error parsing WebSocket message:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
      setConnected(false);
    };

    ws.onclose = () => {
      console.log('🔌 WebSocket closed');
      setConnected(false);
    };

    return () => {
      ws.close();
    };
  }, []);

  return { totalSteps, connected };
}
```

**Usage in Component:**

```typescript
// components/StepsCounter.tsx
import { usePublicSteps } from '@/hooks/usePublicSteps';

export function StepsCounter() {
  const { totalSteps, connected } = usePublicSteps();

  return (
    <div className="steps-counter">
      <h2>Totaal Aantal Stappen</h2>
      <div className="count">
        {totalSteps.toLocaleString('nl-NL')}
      </div>
      {!connected && (
        <span className="status">Verbinding maken...</span>
      )}
    </div>
  );
}
```

### 2. Vue/Nuxt - Public Leaderboard

```typescript
// composables/usePublicLeaderboard.ts
import { ref, onMounted, onUnmounted } from 'vue';

export function usePublicLeaderboard() {
  const leaderboard = ref([]);
  const connected = ref(false);
  let ws: WebSocket | null = null;

  onMounted(() => {
    const config = useRuntimeConfig();
    ws = new WebSocket(
      config.public.wsUrl || 'wss://dklemailservice.onrender.com/ws/steps'
    );

    ws.onopen = () => {
      console.log('✅ Leaderboard WebSocket connected');
      connected.value = true;
      
      ws?.send(JSON.stringify({
        type: 'subscribe',
        channels: ['leaderboard_updates']
      }));
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      if (message.type === 'leaderboard_update') {
        leaderboard.value = message.data.entries;
      }
    };

    ws.onerror = () => {
      connected.value = false;
    };

    ws.onclose = () => {
      connected.value = false;
    };
  });

  onUnmounted(() => {
    ws?.close();
  });

  return {
    leaderboard,
    connected
  };
}
```

### 3. Vanilla JavaScript - Simple Integration

```javascript
// public-steps.js
class PublicStepsTracker {
  constructor(elementId) {
    this.element = document.getElementById(elementId);
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.connect();
  }

  connect() {
    const wsUrl = 'wss://dklemailservice.onrender.com/ws/steps';
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('✅ Steps tracker connected');
      this.reconnectAttempts = 0;
      
      // Subscribe to public channels
      this.ws.send(JSON.stringify({
        type: 'subscribe',
        channels: ['total_updates']
      }));
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      if (message.type === 'total_update') {
        this.updateDisplay(message.data.total_steps);
      }
    };

    this.ws.onclose = () => {
      console.log('🔌 Connection closed, attempting reconnect...');
      this.reconnect();
    };

    this.ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };
  }

  reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
      
      console.log(`⏳ Reconnecting in ${delay}ms...`);
      setTimeout(() => this.connect(), delay);
    } else {
      console.error('❌ Max reconnection attempts reached');
    }
  }

  updateDisplay(totalSteps) {
    if (this.element) {
      this.element.textContent = totalSteps.toLocaleString('nl-NL');
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

// Usage
const tracker = new PublicStepsTracker('steps-counter');
```

## Message Types

### Welcome Message (ontvangen bij verbinding)

```json
{
  "type": "welcome",
  "message": "Connected to StepsHub! Send {...} to receive updates",
  "available_channels": ["step_updates", "total_updates", "leaderboard_updates", "badge_earned"],
  "timestamp": 1704723600
}
```

### Total Steps Update

```json
{
  "type": "total_update",
  "data": {
    "total_steps": 1500000,
    "total_participants": 250,
    "year": 2025
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

### Leaderboard Update

```json
{
  "type": "leaderboard_update",
  "data": {
    "period": "daily",
    "entries": [
      {
        "rank": 1,
        "participant_id": "uuid",
        "naam": "Jane Smith",
        "steps": 25000
      }
    ]
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

## Best Practices

### 1. Connection Management

```typescript
// ✅ GOOD - Proper cleanup
useEffect(() => {
  const ws = new WebSocket(url);
  // ... setup handlers
  
  return () => {
    ws.close(); // Always cleanup!
  };
}, []);

// ❌ BAD - Memory leak
useEffect(() => {
  const ws = new WebSocket(url);
  // ... setup handlers
  // Missing cleanup!
}, []);
```

### 2. Error Handling

```typescript
// ✅ GOOD - Robust error handling
ws.onmessage = (event) => {
  try {
    const message = JSON.parse(event.data);
    handleMessage(message);
  } catch (error) {
    console.error('Failed to parse message:', error);
  }
};

// ❌ BAD - No error handling
ws.onmessage = (event) => {
  const message = JSON.parse(event.data); // Can throw!
  handleMessage(message);
};
```

### 3. Reconnection Logic

```typescript
// ✅ GOOD - Exponential backoff
const reconnect = () => {
  attempts++;
  const delay = Math.min(1000 * Math.pow(2, attempts), 30000);
  setTimeout(() => connect(), delay);
};

// ❌ BAD - Aggressive reconnection
const reconnect = () => {
  setTimeout(() => connect(), 100); // Too fast!
};
```

## Environment Variables

```bash
# .env.local (Next.js)
NEXT_PUBLIC_WS_URL=wss://dklemailservice.onrender.com/ws/steps

# .env (Nuxt)
NUXT_PUBLIC_WS_URL=wss://dklemailservice.onrender.com/ws/steps

# Development
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws/steps
```

## Troubleshooting

### Probleem: Connection Failed

**Symptomen:**
```
WebSocket connection to 'wss://...' failed
```

**Oplossingen:**
1. ✅ Controleer of de URL correct is (`/ws/steps` niet `/api/ws/steps` voor frontend)
2. ✅ Verwijder lege token parameters (`?token=`)
3. ✅ Gebruik `wss://` in productie, `ws://` in development
4. ✅ Check CORS instellingen in backend

### Probleem: No Messages Received

**Symptomen:**
```
Connected but no updates coming through
```

**Oplossingen:**
1. ✅ Verstuur subscribe bericht na verbinding
2. ✅ Check of juiste channels zijn gespecificeerd
3. ✅ Verify message parsing in `onmessage` handler

### Probleem: Too Many Reconnections

**Symptomen:**
```
Rapid reconnection attempts flooding logs
```

**Oplossingen:**
1. ✅ Implementeer exponential backoff
2. ✅ Stel maximum aantal pogingen in
3. ✅ Cleanup oude connectie voordat nieuwe maken

## Testing

### Test Public Connection

```bash
# Install wscat
npm install -g wscat

# Test public connection
wscat -c wss://dklemailservice.onrender.com/ws/steps

# Should receive welcome message
# Send subscribe:
{"type":"subscribe","channels":["total_updates"]}
```

### Monitor in Browser DevTools

```javascript
// Open console and run:
const ws = new WebSocket('wss://dklemailservice.onrender.com/ws/steps');
ws.onmessage = (e) => console.log('Received:', JSON.parse(e.data));
ws.onopen = () => {
  console.log('Connected!');
  ws.send(JSON.stringify({
    type: 'subscribe',
    channels: ['total_updates', 'leaderboard_updates']
  }));
};
```

## Support

Voor vragen of problemen:
- Zie [WebSocket API Documentation](../api/WEBSOCKET.md)
- Zie [Steps Gamification API](../api/STEPS_GAMIFICATION.md)
- Check backend logs voor connection errors

## Changelog

**2025-01-09:**
- ✅ Fixed lege token handling
- ✅ Added public access support  
- ✅ Added `/ws/steps` alias route
- ✅ Improved error handling and logging
- ✅ Updated documentation met correcte URLs