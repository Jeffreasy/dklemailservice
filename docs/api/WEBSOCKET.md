# WebSocket API

Real-time communication for Notulen and Steps applications.

## Connection

### WebSocket URLs

```
Development: ws://localhost:8080/ws/{channel}
Production: wss://api.dklemailservice.com/ws/{channel}
```

**Available Channels:**
- `/ws/notulen` - Meeting notes real-time updates
- `/ws/steps` - Steps application real-time updates

### Authentication

WebSocket connections require JWT authentication via query parameter:

```javascript
// Steps WebSocket
const wsSteps = new WebSocket('ws://localhost:8080/ws/steps?token=YOUR_JWT_TOKEN');

// Notulen WebSocket
const wsNotulen = new WebSocket('ws://localhost:8080/api/ws/notulen?token=YOUR_JWT_TOKEN');
```

---

## Notulen WebSocket

Real-time updates for meeting notes (Notulen) system.

### Connection Example

```javascript
const token = 'your_jwt_token_here';
const ws = new WebSocket(`ws://localhost:8080/api/ws/notulen?token=${token}`);

ws.onopen = () => {
  console.log('Connected to Notulen WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleNotulenUpdate(message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket connection closed');
};
```

### Message Types

#### 1. Notulen Created

Sent when a new meeting note is created.

```json
{
  "type": "notulen_created",
  "data": {
    "id": 123,
    "titel": "Team Meeting",
    "datum": "2025-01-08T10:00:00Z",
    "aanwezig": ["John Doe", "Jane Smith"],
    "afwezig": ["Bob Johnson"],
    "verslag": "Meeting notes content...",
    "actiepunten": ["Action item 1", "Action item 2"],
    "created_by": 1,
    "created_at": "2025-01-08T10:00:00Z"
  },
  "timestamp": "2025-01-08T10:00:00Z"
}
```

#### 2. Notulen Updated

Sent when a meeting note is updated.

```json
{
  "type": "notulen_updated",
  "data": {
    "id": 123,
    "titel": "Team Meeting - Updated",
    "datum": "2025-01-08T10:00:00Z",
    "verslag": "Updated meeting notes...",
    "updated_at": "2025-01-08T11:00:00Z"
  },
  "timestamp": "2025-01-08T11:00:00Z"
}
```

#### 3. Notulen Deleted

Sent when a meeting note is deleted.

```json
{
  "type": "notulen_deleted",
  "data": {
    "id": 123
  },
  "timestamp": "2025-01-08T12:00:00Z"
}
```

---

## Steps WebSocket

Real-time updates for Steps application tracking.

### Connection Example

```javascript
const token = 'your_jwt_token_here';
const ws = new WebSocket(`ws://localhost:8080/ws/steps?token=${token}`);

ws.onopen = () => {
  console.log('Connected to Steps WebSocket');
  
  // Subscribe to specific user updates
  ws.send(JSON.stringify({
    type: 'subscribe',
    user_id: 123
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleStepsUpdate(message);
};
```

### Message Types

#### 1. Steps Progress Update

Sent when a user's step count is updated.

```json
{
  "type": "steps_progress",
  "data": {
    "user_id": 123,
    "date": "2025-01-08",
    "steps": 5000,
    "goal": 10000,
    "percentage": 50,
    "calories": 250,
    "distance_km": 3.5
  },
  "timestamp": "2025-01-08T14:30:00Z"
}
```

#### 2. Steps Goal Achieved

Sent when a user reaches their daily step goal.

```json
{
  "type": "steps_goal_achieved",
  "data": {
    "user_id": 123,
    "date": "2025-01-08",
    "steps": 10000,
    "goal": 10000,
    "achievement_unlocked": true,
    "badge_id": 5
  },
  "timestamp": "2025-01-08T18:45:00Z"
}
```

#### 3. Leaderboard Update

Sent when the leaderboard changes.

```json
{
  "type": "leaderboard_update",
  "data": {
    "period": "daily",
    "top_users": [
      {
        "user_id": 123,
        "naam": "John Doe",
        "steps": 15000,
        "rank": 1
      },
      {
        "user_id": 456,
        "naam": "Jane Smith",
        "steps": 12000,
        "rank": 2
      }
    ]
  },
  "timestamp": "2025-01-08T19:00:00Z"
}
```

---

## Client-Side Implementation

### React Example

```javascript
import { useEffect, useState } from 'react';

function useNotulenWebSocket(token) {
  const [notulen, setNotulen] = useState([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws/notulen?token=${token}`);

    ws.onopen = () => {
      console.log('WebSocket connected');
      setConnected(true);
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      switch (message.type) {
        case 'notulen_created':
          setNotulen(prev => [...prev, message.data]);
          break;
        case 'notulen_updated':
          setNotulen(prev => 
            prev.map(n => n.id === message.data.id ? message.data : n)
          );
          break;
        case 'notulen_deleted':
          setNotulen(prev => prev.filter(n => n.id !== message.data.id));
          break;
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setConnected(false);
    };

    ws.onclose = () => {
      console.log('WebSocket closed');
      setConnected(false);
    };

    return () => {
      ws.close();
    };
  }, [token]);

  return { notulen, connected };
}
```

### Vue Example

```javascript
import { ref, onMounted, onUnmounted } from 'vue';

export function useStepsWebSocket(token) {
  const steps = ref(0);
  const connected = ref(false);
  let ws = null;

  onMounted(() => {
    ws = new WebSocket(`ws://localhost:8080/ws/steps?token=${token}`);

    ws.onopen = () => {
      connected.value = true;
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      if (message.type === 'steps_progress') {
        steps.value = message.data.steps;
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      connected.value = false;
    };

    ws.onclose = () => {
      connected.value = false;
    };
  });

  onUnmounted(() => {
    if (ws) {
      ws.close();
    }
  });

  return { steps, connected };
}
```

---

## Connection Management

### Reconnection Strategy

Implement automatic reconnection with exponential backoff:

```javascript
class WebSocketClient {
  constructor(url, token) {
    this.url = url;
    this.token = token;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
  }

  connect() {
    this.ws = new WebSocket(`${this.url}?token=${this.token}`);

    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;
    };

    this.ws.onclose = () => {
      this.reconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    this.ws.onmessage = (event) => {
      this.handleMessage(JSON.parse(event.data));
    };
  }

  reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
      
      console.log(`Reconnecting in ${delay}ms...`);
      setTimeout(() => this.connect(), delay);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    }
  }

  close() {
    if (this.ws) {
      this.ws.close();
    }
  }

  handleMessage(message) {
    // Override this method in your implementation
    console.log('Received:', message);
  }
}
```

### Heartbeat / Keep-Alive

Implement ping/pong to keep connection alive:

```javascript
let heartbeatInterval;

ws.onopen = () => {
  // Send ping every 30 seconds
  heartbeatInterval = setInterval(() => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'ping' }));
    }
  }, 30000);
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === 'pong') {
    console.log('Received pong');
  } else {
    handleMessage(message);
  }
};

ws.onclose = () => {
  clearInterval(heartbeatInterval);
};
```

---

## Error Handling

### Connection Errors

```json
{
  "type": "error",
  "error": "Authentication failed",
  "code": "AUTH_INVALID",
  "timestamp": "2025-01-08T10:00:00Z"
}
```

### Common Error Codes

| Code | Description | Action |
|------|-------------|--------|
| `AUTH_INVALID` | Invalid or expired token | Re-authenticate |
| `AUTH_REQUIRED` | No token provided | Provide valid token |
| `PERMISSION_DENIED` | Insufficient permissions | Check user roles |
| `RATE_LIMIT` | Too many connections | Implement backoff |
| `SERVICE_UNAVAILABLE` | Server maintenance | Retry later |

---

## Performance Considerations

### Message Throttling

Implement client-side throttling for high-frequency updates:

```javascript
function throttle(func, delay) {
  let lastCall = 0;
  return function(...args) {
    const now = Date.now();
    if (now - lastCall >= delay) {
      lastCall = now;
      func.apply(this, args);
    }
  };
}

const throttledUpdate = throttle((data) => {
  updateUI(data);
}, 100); // Max 10 updates per second

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  throttledUpdate(message.data);
};
```

### Connection Pooling

For multiple WebSocket connections, implement a connection pool:

```javascript
class WebSocketPool {
  constructor() {
    this.connections = new Map();
  }

  getConnection(channel, token) {
    const key = `${channel}:${token}`;
    
    if (!this.connections.has(key)) {
      const ws = new WebSocket(`ws://localhost:8080/ws/${channel}?token=${token}`);
      this.connections.set(key, ws);
    }
    
    return this.connections.get(key);
  }

  closeConnection(channel, token) {
    const key = `${channel}:${token}`;
    const ws = this.connections.get(key);
    
    if (ws) {
      ws.close();
      this.connections.delete(key);
    }
  }
}
```

---

## Testing

### Manual Testing with wscat

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c "ws://localhost:8080/ws/notulen?token=YOUR_JWT_TOKEN"

# Send test message
> {"type":"ping"}

# Receive response
< {"type":"pong","timestamp":"2025-01-08T10:00:00Z"}
```

### Automated Testing

```javascript
describe('WebSocket Connection', () => {
  it('should connect successfully with valid token', (done) => {
    const ws = new WebSocket('ws://localhost:8080/ws/notulen?token=VALID_TOKEN');
    
    ws.onopen = () => {
      expect(ws.readyState).toBe(WebSocket.OPEN);
      ws.close();
      done();
    };
  });

  it('should reject connection with invalid token', (done) => {
    const ws = new WebSocket('ws://localhost:8080/ws/notulen?token=INVALID');
    
    ws.onerror = () => {
      done();
    };
  });
});
```

---

For implementation examples, see the [WebSocket Integration Guide](../guides/WEBSOCKET.md).