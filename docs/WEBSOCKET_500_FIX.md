# WebSocket 500 Error Fix

## Problem Identified
The WebSocket endpoint at `wss://dklemailservice.onrender.com/ws/steps` was returning HTTP 500 errors.

## Root Cause
In [`handlers/steps_websocket_handler.go`](../handlers/steps_websocket_handler.go), line 95 was attempting to call `c.Locals("userID")` on a `*websocket.Conn` object:

```go
// WRONG - This caused the 500 error
func (h *StepsWebSocketHandler) HandleWebSocket(c *websocket.Conn) {
    // ...
    if uid := c.Locals("userID"); uid != nil {  // ❌ websocket.Conn doesn't have Locals()
        // ...
    }
}
```

The `Locals()` method only exists on `*fiber.Ctx`, not on `*websocket.Conn`. After the WebSocket upgrade, the handler receives a websocket connection object, not a Fiber context.

## Solution Implemented
Refactored the WebSocket handler to properly handle authentication **before** the upgrade:

1. **Combined Handler Approach**: Instead of using separate middleware and handler, created a unified handler that:
   - Validates authentication **before** WebSocket upgrade (when we still have `*fiber.Ctx`)
   - Passes the authenticated user ID directly to the WebSocket connection handler
   - Upgrades to WebSocket with the authentication context

2. **Code Changes**:
```go
// CORRECT - Fixed implementation
func (h *StepsWebSocketHandler) RegisterRoutes(app *fiber.App) {
    wsHandler := func(c *fiber.Ctx) error {
        // Check WebSocket upgrade
        if !websocket.IsWebSocketUpgrade(c) {
            return fiber.ErrUpgradeRequired
        }

        // Validate token BEFORE upgrade (we still have fiber.Ctx here)
        var authenticatedUserID string
        token := c.Query("token")
        if token == "" {
            authHeader := c.Get("Authorization")
            if strings.HasPrefix(authHeader, "Bearer ") {
                token = strings.TrimPrefix(authHeader, "Bearer ")
            }
        }
        
        if token != "" {
            userID, err := h.authService.ValidateToken(token)
            if err == nil {
                authenticatedUserID = userID
            }
        }

        // Upgrade with context
        return websocket.New(func(conn *websocket.Conn) {
            h.handleWebSocketConnection(conn, authenticatedUserID)
        })(c)
    }

    app.Get("/api/ws/steps", wsHandler)
    app.Get("/ws/steps", wsHandler)
}

// Handler now receives authenticated user ID as parameter
func (h *StepsWebSocketHandler) handleWebSocketConnection(c *websocket.Conn, authenticatedUserID string) {
    userID := c.Query("user_id")
    
    // Use authenticated user ID if not in query params
    if userID == "" {
        userID = authenticatedUserID
    }
    
    // ... rest of handler logic
}
```

## Benefits
1. ✅ **Fixes 500 Error**: No more invalid method calls on websocket.Conn
2. ✅ **Proper Authentication**: Token validation happens at the correct time
3. ✅ **Anonymous Access**: Still allows unauthenticated connections for public features
4. ✅ **Clean Architecture**: Authentication logic is cleanly separated from WebSocket handling

## Testing
After deploying this fix:
1. WebSocket connections should successfully upgrade without 500 errors
2. Authenticated connections work with `?token=xxx` query parameter or `Authorization: Bearer xxx` header
3. Anonymous connections work for public leaderboard/stats viewing
4. Welcome message should be sent upon connection

## Related Files
- [`handlers/steps_websocket_handler.go`](../handlers/steps_websocket_handler.go) - Main fix location
- [`services/steps_hub.go`](../services/steps_hub.go) - WebSocket hub (unchanged)
- [`services/steps_service.go`](../services/steps_service.go) - Steps service (unchanged)
- [`main.go`](../main.go) - WebSocket initialization (unchanged)

## Frontend Integration
The frontend can now connect using:
```typescript
const ws = new WebSocket(`wss://dklemailservice.onrender.com/ws/steps?token=${authToken}`);
// or
const ws = new WebSocket('wss://dklemailservice.onrender.com/ws/steps'); // anonymous
```

See [`docs/frontend/WEBSOCKET_STEPS_INTEGRATION.md`](./frontend/WEBSOCKET_STEPS_INTEGRATION.md) for complete integration guide.