# WebSocket Integration Guide

## 📋 Overzicht

Deze guide beschrijft hoe je de StepsHub WebSocket functionaliteit integreert in de bestaande applicatie.

---

## ✅ Wat is al geïmplementeerd

### 1. Core Components

- ✅ **`services/steps_hub.go`** - StepsHub met message types en broadcasting
- ✅ **`handlers/steps_websocket_handler.go`** - WebSocket handler met auth
- ✅ **`services/steps_service.go`** - Updated met broadcast functies

### 2. Features

- ✅ Client registration/unregistration
- ✅ Subscription-based message filtering
- ✅ Ping/pong voor keep-alive
- ✅ 4 message types: step_update, total_update, leaderboard_update, badge_earned
- ✅ JWT authentication support
- ✅ Automatic broadcasts bij step updates

---

## 🔧 Integratie Stappen

### Stap 1: Update main.go

Voeg de volgende code toe aan `main.go`:

```go
package main

import (
    // ... existing imports
    "dklautomationgo/handlers"
    "dklautomationgo/services"
)

func main() {
    // ... existing initialization ...

    // Initialize repositories
    aanmeldingRepo := repository.NewPostgresAanmeldingRepository(baseRepo)
    routeFundRepo := repository.NewPostgresRouteFundRepository(baseRepo)
    
    // Initialize services
    authService := services.NewAuthService(/* ... */)
    gamificationService := services.NewGamificationService(/* ... */)
    
    // 🆕 Create StepsService
    stepsService := services.NewStepsService(db, aanmeldingRepo, routeFundRepo)
    
    // 🆕 Create StepsHub
    stepsHub := services.NewStepsHub(stepsService, gamificationService)
    
    // 🆕 Link StepsHub to StepsService
    stepsService.SetStepsHub(stepsHub)
    
    // 🆕 Start StepsHub (in goroutine)
    go stepsHub.Run()
    
    // Initialize handlers
    stepsHandler := handlers.NewStepsHandler(stepsService, authService, permissionService)
    
    // 🆕 Create WebSocket handler
    stepsWsHandler := handlers.NewStepsWebSocketHandler(stepsHub, authService)
    
    // Register routes
    stepsHandler.RegisterRoutes(app)
    
    // 🆕 Register WebSocket routes
    stepsWsHandler.RegisterRoutes(app)
    
    // 🆕 Optional: Stats endpoint (admin only)
    app.Get("/api/ws/stats", 
        handlers.AuthMiddleware(authService),
        handlers.PermissionMiddleware(permissionService, "admin", "read"),
        stepsWsHandler.GetStats,
    )
    
    // ... rest of application setup ...
    
    // Start server
    logger.Info("Starting server on :8080")
    if err := app.Listen(":8080"); err != nil {
        logger.Fatal("Failed to start server", "error", err)
    }
}
```

### Stap 2: Environment Variables (Optioneel)

Voeg aan `.env` toe voor configuratie:

```bash
# WebSocket Configuration
WEBSOCKET_ENABLED=true
WEBSOCKET_MAX_CONNECTIONS=10000
WEBSOCKET_PING_INTERVAL=30s
WEBSOCKET_READ_BUFFER_SIZE=1024
WEBSOCKET_WRITE_BUFFER_SIZE=1024
```

---

## 🧪 Testen

### Test 1: WebSocket Connectie

```bash
# Test WebSocket connectie
wscat -c "ws://localhost:8080/ws/steps?user_id=test-user&participant_id=test-participant"

# Subscribe to channels
{"type":"subscribe","channels":["step_updates","total_updates","leaderboard_updates"]}
```

### Test 2: Step Update via REST API

```bash
# Trigger een step update (in andere terminal)
curl -X POST http://localhost:8080/api/steps \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"steps": 1000}'

# Check WebSocket terminal - je zou een message moeten zien:
# {"type":"step_update","participant_id":"...","steps":1000, ...}
```

### Test 3: Check Stats

```bash
curl http://localhost:8080/api/ws/stats \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"

# Response:
# {
#   "total_clients": 5,
#   "subscriptions": {
#     "step_updates": 3,
#     "total_updates": 5,
#     "leaderboard_updates": 2
#   }
# }
```

---

## 📡 Frontend Integration

### JavaScript/TypeScript Client

Gebruik de client class uit het architectuur document:

```typescript
import { StepsWebSocketClient } from './lib/websocket-client';

const client = new StepsWebSocketClient(
    'ws://localhost:8080/ws/steps',
    localStorage.getItem('token'),
    userId,
    participantId
);

// Subscribe to events
client.on('step_update', (data) => {
    console.log('Steps updated:', data.steps);
    updateDashboard(data);
});

client.on('total_update', (data) => {
    console.log('Total steps:', data.total_steps);
    updateTotalCounter(data.total_steps);
});

client.on('leaderboard_update', (data) => {
    console.log('Leaderboard:', data.entries);
    renderLeaderboard(data.entries);
});

// Connect
client.connect();
```

### React Hook

```typescript
import { useStepsWebSocket } from './hooks/useStepsWebSocket';

function Dashboard({ userId, participantId }) {
    const { connected, latestUpdate, totalSteps } = useStepsWebSocket(
        userId, 
        participantId
    );
    
    return (
        <div>
            <div>Status: {connected ? '🟢 Connected' : '🔴 Disconnected'}</div>
            <div>Your Steps: {latestUpdate?.steps || 0}</div>
            <div>Total Steps: {totalSteps}</div>
        </div>
    );
}
```

---

## 🔒 Security Checklist

- [x] JWT validation in WebSocket handshake
- [x] User isolation (participants only see own updates)
- [x] Subscription filtering
- [ ] Rate limiting per client (TODO in production)
- [ ] Connection limit per user (TODO in production)
- [ ] Origin checking (TODO for CORS)

---

## 📊 Monitoring

### Metrics om te tracken

1. **Active Connections**: `stepsHub.GetClientCount()`
2. **Subscriptions**: `stepsHub.GetSubscriptionCount()`
3. **Message Rate**: Track broadcasts per second
4. **Connection Duration**: Track time.Since(clientConnectTime)
5. **Error Rate**: Track websocket errors

### Logging

Key events die gelogd worden:

```go
logger.Info("WebSocket client connecting", "user_id", userID)
logger.Info("WebSocket client connected", "total_clients", count)
logger.Info("WebSocket client disconnected", "user_id", userID)
logger.Error("WebSocket auth failed", "error", err)
```

---

## 🚨 Troubleshooting

### Issue: WebSocket upgrade fails

**Symptom**: `426 Upgrade Required`

**Solution**:
```go
// Zorg dat Upgrade headers aanwezig zijn
app.Use("/ws/steps", func(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
        return c.Next()
    }
    return fiber.ErrUpgradeRequired
})
```

### Issue: No messages received

**Debug Steps**:
1. Check of client geregistreerd is: `stepsHub.GetClientCount()`
2. Check subscriptions: Client moet subscriben via `{"type":"subscribe","channels":[...]}`
3. Check of StepsHub.Run() draait (moet in goroutine)
4. Check logs voor errors

### Issue: Memory leak

**Symptom**: Memory usage groeit continu

**Check**:
1. Clients worden correct unregistered bij disconnect
2. Send channels worden gesloten
3. Goroutines eindigen na disconnect

---

## 🎯 Next Steps

### Short Term (Week 1-2)

- [ ] Implementeer rate limiting per client
- [ ] Add connection limits per user
- [ ] Implement message queue voor burst traffic
- [ ] Add Prometheus metrics
- [ ] Write integration tests

### Medium Term (Week 3-4)

- [ ] Add Redis pub/sub voor multi-instance support
- [ ] Implement message persistence
- [ ] Add replay functionality
- [ ] Build admin dashboard voor monitoring
- [ ] Load testing

### Long Term (Month 2-3)

- [ ] Auto-scaling based on connection count
- [ ] Geographic distribution
- [ ] Message compression
- [ ] Custom protocol optimization

---

## 📚 Resources

- **Architecture Doc**: [`docs/STEPS_ARCHITECTURE_WEBSOCKETS.md`](STEPS_ARCHITECTURE_WEBSOCKETS.md)
- **API Doc**: [`docs/api/steps-api.md`](api/steps-api.md)
- **Fiber WebSocket**: https://docs.gofiber.io/api/middleware/websocket/
- **WebSocket RFC**: https://datatracker.ietf.org/doc/html/rfc6455

---

## ✅ Deployment Checklist

### Pre-Deployment

- [ ] Code review completed
- [ ] Unit tests passing
- [ ] Integration tests passing
- [ ] Load tests completed
- [ ] Security audit done
- [ ] Documentation updated

### Deployment

- [ ] Deploy to staging
- [ ] Smoke tests on staging
- [ ] Monitor metrics for 24h
- [ ] Deploy to production (canary)
- [ ] Monitor metrics for 48h
- [ ] Full production rollout

### Post-Deployment

- [ ] Monitor error rates
- [ ] Monitor connection count
- [ ] Monitor message latency
- [ ] User feedback collected
- [ ] Performance tuning if needed

---

**Document Version**: 1.0  
**Last Updated**: 2025-01-02  
**Status**: 🟢 Ready for Integration