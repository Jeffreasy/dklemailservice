# Quick Start: Steps API voor Frontend

**Last Updated**: 2025-11-09  
**Backend URL**: `http://localhost:8080` (local) | `https://dklemailservice.onrender.com` (prod)

---

## 🚀 Snel aan de slag

### 1. Haal Total Steps op (REST API)
```typescript
const response = await fetch('http://localhost:8080/api/total-steps');
const data = await response.json();
console.log(data.total_steps); // 98450
```

### 2. WebSocket Real-time Connectie
```typescript
const ws = new WebSocket('ws://localhost:8080/ws/steps');

ws.onopen = () => {
    // Subscribe naar updates
    ws.send(JSON.stringify({
        type: 'subscribe',
        channels: ['step_updates', 'total_updates', 'leaderboard_updates']
    }));
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Update:', data);
};
```

---

## 📡 Beschikbare Endpoints

### Public Endpoints (Geen Auth)

#### GET /api/total-steps
Haal totaal aantal stappen op
```bash
curl http://localhost:8080/api/total-steps
```
Response:
```json
{
  "total_steps": 98450,
  "year": 0
}
```

#### GET /api/funds-distribution
Haal fondsverdeling op
```bash
curl http://localhost:8080/api/funds-distribution
```

#### GET /api/health
Service health check
```bash
curl http://localhost:8080/api/health
```

### WebSocket Endpoint

#### WS /ws/steps
Real-time steps updates
```
ws://localhost:8080/ws/steps
```

**Supported Message Types**:
- `welcome` - Ontvangen bij connectie
- `subscribe` - Subscribe naar channels
- `step_update` - Individual step change
- `total_update` - Total steps update
- `leaderboard_update` - Top N leaderboard
- `badge_earned` - Achievement notification
- `ping/pong` - Keep-alive

---

## 💡 Frontend Voorbeelden

### React Hook
```typescript
import { useEffect, useState } from 'react';

function useStepsData() {
    const [totalSteps, setTotalSteps] = useState(0);
    const [ws, setWs] = useState<WebSocket | null>(null);

    useEffect(() => {
        // Haal initiële data op
        fetch('http://localhost:8080/api/total-steps')
            .then(res => res.json())
            .then(data => setTotalSteps(data.total_steps));

        // Connect WebSocket voor updates
        const socket = new WebSocket('ws://localhost:8080/ws/steps');
        
        socket.onopen = () => {
            socket.send(JSON.stringify({
                type: 'subscribe',
                channels: ['total_updates']
            }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === 'total_update') {
                setTotalSteps(data.total_steps);
            }
        };

        setWs(socket);

        return () => socket.close();
    }, []);

    return { totalSteps, ws };
}

// Gebruik in component
function StepsDisplay() {
    const { totalSteps } = useStepsData();
    return <div>Total Steps: {totalSteps.toLocaleString()}</div>;
}
```

### Vue Composable
```typescript
import { ref, onMounted, onUnmounted } from 'vue';

export function useSteps() {
    const totalSteps = ref(0);
    let ws: WebSocket | null = null;

    onMounted(async () => {
        // REST API
        const res = await fetch('http://localhost:8080/api/total-steps');
        const data = await res.json();
        totalSteps.value = data.total_steps;

        // WebSocket
        ws = new WebSocket('ws://localhost:8080/ws/steps');
        ws.onopen = () => {
            ws?.send(JSON.stringify({
                type: 'subscribe',
                channels: ['total_updates']
            }));
        };
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === 'total_update') {
                totalSteps.value = data.total_steps;
            }
        };
    });

    onUnmounted(() => {
        ws?.close();
    });

    return { totalSteps };
}
```

---

## 🧪 Test Data Beschikbaar

Lokale Docker heeft test data:
- **5 test participants**
- **98,450 total steps**
- **Leaderboard data**

Gebruik dit voor UI development en testing!

---

## 🔍 Debugging

### Probleem: "total_steps": 0

**Oplossing**:
```bash
# Check of test data er is
curl http://localhost:8080/api/total-steps

# Als 0, voeg test data toe:
docker exec -i dkl-postgres psql -U postgres -d dklemailservice < scripts/add_test_steps_data.sql

# Verify
curl http://localhost:8080/api/total-steps
# Expected: {"total_steps":98450}
```

### Probleem: WebSocket verbindt niet

**Check**:
1. Is backend running? `curl http://localhost:8080/api/health`
2. Browser console errors?
3. Correcte URL? `ws://localhost:8080/ws/steps` (niet `wss://`)

### Probleem: CORS errors

**Oplossing**: Check `.env`:
```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

---

## 📚 Meer Documentatie

- [Complete WebSocket Integration Guide](./WEBSOCKET_STEPS_INTEGRATION.md)
- [WebSocket Fix Details](../WEBSOCKET_500_FIX.md)
- [Test Results](../WEBSOCKET_TEST_RESULTS.md)
- [Frontend Data Solution](../FRONTEND_STEPS_DATA_SOLUTION.md)

---

**Status**: ✅ KLAAR VOOR FRONTEND DEVELOPMENT  
**Data Available**: YES (98,450 steps)  
**WebSocket**: OPERATIONAL  
**REST API**: OPERATIONAL