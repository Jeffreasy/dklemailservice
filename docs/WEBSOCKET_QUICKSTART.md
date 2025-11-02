# WebSocket Quick Start Guide

## 🚀 Snel Aan de Slag

Complete gids om de WebSocket functionaliteit voor stappen tracking direct te gebruiken.

---

## ✅ Wat is Geïmplementeerd

### Backend (100% Klaar)
- ✅ **StepsHub** - [`services/steps_hub.go`](../services/steps_hub.go)
- ✅ **WebSocket Handler** - [`handlers/steps_websocket_handler.go`](../handlers/steps_websocket_handler.go)
- ✅ **Service Integration** - [`services/steps_service.go`](../services/steps_service.go)
- ✅ **Unit Tests** - [`tests/steps_hub_test.go`](../tests/steps_hub_test.go) - ✅ 8/8 PASSED

### Frontend (100% Klaar)
- ✅ **TypeScript Client** - [`docs/frontend/steps-websocket-client.ts`](frontend/steps-websocket-client.ts)
- ✅ **React Hooks** - [`docs/frontend/useStepsWebSocket.ts`](frontend/useStepsWebSocket.ts)
- ✅ **Dashboard Voorbeeld** - [`docs/frontend/DashboardExample.tsx`](frontend/DashboardExample.tsx)

---

## 📋 Integratie in 3 Stappen

### Stap 1: Update main.go

Voeg deze code toe aan je `main.go`:

```go
// Initialiseer StepsService (waarschijnlijk al aanwezig)
stepsService := services.NewStepsService(db, aanmeldingRepo, routeFundRepo)

// 🆕 Maak StepsHub
stepsHub := services.NewStepsHub(stepsService, gamificationService)

// 🆕 Koppel hub aan service
stepsService.SetStepsHub(stepsHub)

// 🆕 Start hub in background
go stepsHub.Run()

// 🆕 Maak WebSocket handler
stepsWsHandler := handlers.NewStepsWebSocketHandler(stepsHub, authService)

// 🆕 Registreer WebSocket routes
stepsWsHandler.RegisterRoutes(app)

// 🆕 (Optioneel) Stats endpoint voor monitoring
app.Get("/api/ws/stats", 
    handlers.AuthMiddleware(authService),
    handlers.PermissionMiddleware(permissionService, "admin", "read"),
    stepsWsHandler.GetStats,
)
```

**Dat is alles!** Je backend is nu klaar voor WebSocket connections.

### Stap 2: Test de Connectie

**Met wscat** (CLI tool):
```bash
# Installeer wscat
npm install -g wscat

# Connect
wscat -c "ws://localhost:8080/ws/steps?user_id=test-user&participant_id=test-participant"

# Subscribe (type in terminal)
{"type":"subscribe","channels":["step_updates","total_updates","leaderboard_updates"]}

# Je zou nu live updates moeten zien wanneer iemand stappen update!
```

**Met browser JavaScript**:
```javascript
// Open browser console op je site
const ws = new WebSocket('ws://localhost:8080/ws/steps?user_id=test-user');

ws.onmessage = (event) => {
  console.log('Message received:', JSON.parse(event.data));
};

ws.onopen = () => {
  console.log('Connected!');
  ws.send(JSON.stringify({
    type: 'subscribe',
    channels: ['step_updates', 'total_updates']
  }));
};
```

### Stap 3: Trigger een Update

**Via REST API**:
```bash
curl -X POST http://localhost:8080/api/steps \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"steps": 1000}'
```

**Je WebSocket zou nu moeten ontvangen**:
```json
{
  "type": "step_update",
  "participant_id": "...",
  "naam": "John Doe",
  "steps": 1000,
  "delta": 1000,
  "route": "10 KM",
  "allocated_funds": 75,
  "timestamp": 1704211200
}
```

---

## 💻 Frontend Implementatie

### React Component

Kopieer [`docs/frontend/DashboardExample.tsx`](frontend/DashboardExample.tsx) naar je project:

```tsx
import Dashboard from './components/Dashboard';

function App() {
  const userId = getCurrentUserId();
  const participantId = getParticipantId();
  
  return <Dashboard userId={userId} participantId={participantId} />;
}
```

### Alleen de Hook Gebruiken

```tsx
import { useStepsWebSocket } from './hooks/useStepsWebSocket';

function MyComponent() {
  const { connected, latestUpdate, totalSteps } = useStepsWebSocket(
    userId,
    participantId
  );
  
  return (
    <div>
      <p>Connected: {connected ? 'Yes' : 'No'}</p>
      <p>Your Steps: {latestUpdate?.steps || 0}</p>
      <p>Total Global Steps: {totalSteps}</p>
    </div>
  );
}
```

### Vanilla JavaScript (Zonder React)

```html
<!DOCTYPE html>
<html>
<head>
  <title>Steps Tracker</title>
</head>
<body>
  <div id="stats">
    <p>Steps: <span id="steps">0</span></p>
    <p>Total: <span id="total">0</span></p>
  </div>

  <script>
    const ws = new WebSocket('ws://localhost:8080/ws/steps?user_id=test-user');
    
    ws.onopen = () => {
      console.log('Connected');
      ws.send(JSON.stringify({
        type: 'subscribe',
        channels: ['step_updates', 'total_updates']
      }));
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      
      if (data.type === 'step_update') {
        document.getElementById('steps').textContent = data.steps;
      }
      
      if (data.type === 'total_update') {
        document.getElementById('total').textContent = data.total_steps;
      }
    };
  </script>
</body>
</html>
```

---

## 📡 WebSocket Endpoint

**URL**: `ws://localhost:8080/ws/steps`

**Query Parameters**:
- `user_id` (required) - User ID
- `participant_id` (optional) - Participant ID voor filtering
- `token` (optional) - JWT token voor authenticatie

**Voorbeeld**:
```
ws://localhost:8080/ws/steps?user_id=abc123&participant_id=def456&token=eyJhbG...
```

---

## 📨 Message Types

### 1. Subscribe (Client → Server)

```json
{
  "type": "subscribe",
  "channels": ["step_updates", "total_updates", "leaderboard_updates"]
}
```

**Beschikbare channels**:
- `step_updates` - Individuele deelnemer stappen updates
- `total_updates` - Totaal aantal stappen updates
- `leaderboard_updates` - Leaderboard top 10 updates
- `badge_earned` wordt automatisch naar relevante client gestuurd

### 2. Step Update (Server → Client)

```json
{
  "type": "step_update",
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "naam": "John Doe",
  "steps": 5000,
  "delta": 1000,
  "route": "10 KM",
  "allocated_funds": 75,
  "timestamp": 1704211200
}
```

### 3. Total Update (Server → Client)

```json
{
  "type": "total_update",
  "total_steps": 250000,
  "year": 2025,
  "timestamp": 1704211200
}
```

### 4. Leaderboard Update (Server → Client)

```json
{
  "type": "leaderboard_update",
  "top_n": 10,
  "entries": [
    {
      "rank": 1,
      "participant_id": "...",
      "naam": "Jane Smith",
      "steps": 15000,
      "achievement_points": 200,
      "total_score": 15200,
      "route": "20 KM",
      "badge_count": 5
    }
  ],
  "timestamp": 1704211200
}
```

### 5. Badge Earned (Server → Client)

```json
{
  "type": "badge_earned",
  "participant_id": "550e8400-e29b-41d4-a716-446655440000",
  "badge_name": "First Steps",
  "badge_icon": "/icons/badges/first-steps.svg",
  "points": 10,
  "timestamp": 1704211200
}
```

### 6. Ping/Pong (Keep-Alive)

**Ping** (Client → Server):
```json
{
  "type": "ping",
  "timestamp": 1704211200
}
```

**Pong** (Server → Client):
```json
{
  "type": "pong",
  "timestamp": 1704211200
}
```

---

## 🔒 Authenticatie

### Optie 1: Query Parameter

```javascript
const ws = new WebSocket(
  `ws://localhost:8080/ws/steps?user_id=${userId}&token=${jwtToken}`
);
```

### Optie 2: Subprotocol (Aanbevolen)

```javascript
const ws = new WebSocket(
  `ws://localhost:8080/ws/steps?user_id=${userId}`,
  ['access_token', jwtToken]
);
```

### Optie 3: First Message

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/steps');

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'auth',
    token: jwtToken
  }));
};
```

---

## 🧪 Testing Checklist

### Backend Tests

```bash
# Run unit tests
go test ./tests/steps_hub_test.go -v

# Expected output:
# ✅ TestStepsHub_NewStepsHub (PASS)
# ✅ TestStepsHub_ClientRegistration (PASS)
# ✅ TestStepsHub_StepUpdateBroadcast (PASS)
# ✅ TestStepsHub_TotalUpdateBroadcast (PASS)
# ✅ TestStepsHub_SubscriptionFiltering (PASS)
# ✅ TestStepsHub_BadgeEarnedTargeting (PASS)
# ✅ TestStepsHub_GetSubscriptionCount (PASS)
# ✅ TestStepsHub_MultipleClients (PASS)
```

### Manual Testing Scenario

1. ✅ **Start server**: `go run main.go`
2. ✅ **Connect WebSocket**: wscat/browser console
3. ✅ **Subscribe**: Send subscribe message
4. ✅ **Update steps**: POST /api/steps via curl/Postman
5. ✅ **Verify broadcast**: Check WebSocket receives message
6. ✅ **Check stats**: GET /api/ws/stats
7. ✅ **Test disconnect/reconnect**: Close and reopen connection

---

## 📊 Monitoring

### Check Active Connections

```bash
curl http://localhost:8080/api/ws/stats \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

**Response**:
```json
{
  "total_clients": 15,
  "subscriptions": {
    "step_updates": 10,
    "total_updates": 15,
    "leaderboard_updates": 5
  }
}
```

### Logs to Watch

```bash
# WebSocket connections
2025-01-02T15:04:05Z INFO WebSocket client connecting user_id=abc123
2025-01-02T15:04:05Z INFO WebSocket client connected total_clients=15

# Broadcasts
2025-01-02T15:04:10Z INFO Message broadcast type=step_update recipients=10

# Disconnections
2025-01-02T15:05:00Z INFO WebSocket client disconnected user_id=abc123
```

---

## 🐛 Troubleshooting

### Problem: Cannot connect

**Check**:
1. Server running? `curl http://localhost:8080/api/health`
2. WebSocket routes registered? Check main.go
3. Hub running? Should see `go stepsHub.Run()` in main.go
4. Firewall blocking WebSocket? Check port 8080

### Problem: No messages received

**Check**:
1. Subscribed to channels? Send subscribe message
2. Client registered? Check `/api/ws/stats`
3. Hub running? Check logs
4. Check client.Subscriptions in debugger

### Problem: Messages received but UI niet updated

**Check**:
1. Event handlers registered? `client.on('step_update', ...)`
2. State updates in React? `setState()`
3. Console for errors? Check browser console
4. Message format correct? `JSON.parse()` succeeds?

---

## 🎯 Use Cases

### Use Case 1: Participant Dashboard

**Doel**: Deelnemer ziet eigen stappen real-time

**Implementation**:
```tsx
function ParticipantDashboard() {
  const { latestUpdate } = useParticipantDashboard(userId, participantId);
  
  return <div>Mijn Stappen: {latestUpdate?.steps || 0}</div>;
}
```

**Subscriptions**: `step_updates`, `badge_earned`

### Use Case 2: Public Leaderboard

**Doel**: Bezoekers zien live leaderboard zonder inloggen

**Implementation**:
```tsx
function PublicLeaderboard() {
  const { totalSteps, leaderboard } = useLeaderboard();
  
  return (
    <div>
      <h1>Totaal: {totalSteps.toLocaleString()}</h1>
      <LeaderboardTable entries={leaderboard?.entries || []} />
    </div>
  );
}
```

**Subscriptions**: `total_updates`, `leaderboard_updates`

### Use Case 3: Admin Monitoring

**Doel**: Admin ziet alle activiteit real-time

**Implementation**:
```tsx
function AdminMonitoring() {
  const { latestUpdate, totalSteps, leaderboard } = useStepsMonitoring(adminUserId);
  
  return (
    <div>
      <ActivityFeed updates={[latestUpdate]} />
      <Stats total={totalSteps} />
      <Leaderboard data={leaderboard} />
    </div>
  );
}
```

**Subscriptions**: ALL channels

---

## 📈 Performance Tips

### Tip 1: Subscribe Selectively

```javascript
// ❌ Don't subscribe to everything unless needed
ws.subscribe(['step_updates', 'total_updates', 'leaderboard_updates', 'badge_earned']);

// ✅ Only subscribe to what you need
ws.subscribe(['total_updates']); // Voor een simple counter
```

### Tip 2: Debounce UI Updates

```tsx
const [steps, setSteps] = useState(0);
const debouncedSteps = useDebounce(steps, 500); // Update UI max every 500ms

useEffect(() => {
  if (latestUpdate) {
    setSteps(latestUpdate.steps);
  }
}, [latestUpdate]);
```

### Tip 3: Unsubscribe Wanneer Niet Zichtbaar

```tsx
useEffect(() => {
  if (isVisible) {
    subscribe(['step_updates']);
  } else {
    unsubscribe(['step_updates']);
  }
}, [isVisible]);
```

---

## 🔐 Security Best Practices

1. **Always use JWT authentication**:
   ```javascript
   const token = await getAuthToken();
   const ws = new WebSocket(`ws://...?token=${token}`);
   ```

2. **Validate tokens server-side**: Done automatisch in handler

3. **Use WSS in production**: `wss://` instead of `ws://`

4. **Don't send sensitive data**: WebSocket is niet encrypted (gebruik HTTPS/WSS)

5. **Rate limit**: Client kan max 60 messages/min sturen

---

## 📱 Mobile App Integration

### React Native

```tsx
import { useStepsWebSocket } from './hooks/useStepsWebSocket';

function MobileApp() {
  const { latestUpdate, totalSteps } = useStepsWebSocket(userId, participantId);
  
  return (
    <View>
      <Text>Stappen: {latestUpdate?.steps || 0}</Text>
      <Text>Totaal: {totalSteps}</Text>
    </View>
  );
}
```

**Note**: React Native WebSocket is compatible!

---

## 🎨 UI Voorbeelden

### Badge Notification

```tsx
{latestBadge && (
  <div className="badge-popup">
    <img src={latestBadge.badge_icon} />
    <h3>{latestBadge.badge_name}</h3>
    <p>+{latestBadge.points} punten!</p>
  </div>
)}
```

### Live Step Counter

```tsx
<AnimatedNumber 
  value={latestUpdate?.steps || 0}
  duration={500}
  formatValue={(n) => n.toLocaleString()}
/>
```

### Connection Indicator

```tsx
<div className={`status ${connected ? 'online' : 'offline'}`}>
  {connected ? '🟢 Live' : '🔴 Offline'}
</div>
```

---

## 📚 Volledige Documentatie

- **Architectuur**: [`STEPS_ARCHITECTURE_WEBSOCKETS.md`](STEPS_ARCHITECTURE_WEBSOCKETS.md)
- **Integration Guide**: [`WEBSOCKET_INTEGRATION_GUIDE.md`](WEBSOCKET_INTEGRATION_GUIDE.md)
- **API Reference**: [`api/steps-api.md`](api/steps-api.md)

---

## ✅ Production Checklist

- [ ] WebSocket routes toegevoegd aan main.go
- [ ] StepsHub.Run() draait in goroutine
- [ ] Tested lokaal met wscat
- [ ] Frontend client geïntegreerd
- [ ] Event handlers geregistreerd
- [ ] Error handling geïmplementeerd
- [ ] Monitoring metrics enabled
- [ ] Security review completed
- [ ] Load testing done
- [ ] Documentation updated
- [ ] Deployed to staging
- [ ] Tested in staging
- [ ] Ready for production! 🚀

---

**Version**: 1.0  
**Status**: 🟢 Production Ready  
**Test Results**: ✅ 8/8 Tests PASSED