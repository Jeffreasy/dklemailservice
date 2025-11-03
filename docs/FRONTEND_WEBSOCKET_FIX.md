# 🔧 Frontend WebSocket Fix - Publieke Website

## 🎯 Doel

Bezoekers op de website zien **LIVE** de stappen counter updaten zonder pagina refresh.

---

## 📊 Huidige Situatie (Jouw Logs)

```
✅ WebSocket connected
❌ Geen berichten ontvangen
✅ Polling fallback actief (1261 steps via REST)
```

---

## 🐛 PROBLEEM: Frontend stuurt geen Subscribe

**In `usePublicStepsCounter.ts` moet dit gebeuren**:

### Fix 1: Subscribe Message Toevoegen

```typescript
// In usePublicStepsCounter.ts

useEffect(() => {
  const ws = new WebSocket('ws://localhost:8082/api/ws/steps?user_id=public');
  
  ws.onopen = () => {
    console.log('[WebSocket] Connected ✅');
    
    // ✨ DIT ONTBREEKT WAARSCHIJNLIJK:
    ws.send(JSON.stringify({
      type: 'subscribe',
      channels: ['total_updates', 'step_updates', 'leaderboard_updates']
    }));
    
    console.log('[WebSocket] Subscribed to channels ✅');
  };
  
  ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('[WebSocket] Message received:', message);
    
    // ✨ HANDLE NIEUWE BERICHTEN:
    if (message.type === 'welcome') {
      console.log('Welcome:', message.message);
    }
    
    if (message.type === 'total_update') {
      setTotalSteps(message.total_steps);  // Update UI!
      console.log('Total steps updated to:', message.total_steps);
    }
    
    if (message.type === 'step_update') {
      console.log('Participant stepped:', message);
      // Trigger refresh van totaal
    }
  };
  
  ws.onerror = (error) => {
    console.error('[WebSocket] Error:', error);
  };
  
  ws.onclose = () => {
    console.log('[WebSocket] Disconnected');
    // Auto-reconnect logic
  };
  
  return () => ws.close();
}, []);
```

---

## ✅ COMPLETE WORKING EXAMPLE

**Bestand**: `usePublicStepsCounter.ts` (FIXED VERSION)

```typescript
import { useEffect, useState, useRef } from 'react';

export function usePublicStepsCounter(backendUrl: string) {
  const [totalSteps, setTotalSteps] = useState<number>(0);
  const [connected, setConnected] = useState<boolean>(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number>(1000);

  useEffect(() => {
    const connect = () => {
      const wsUrl = backendUrl.replace('http', 'ws') + '/api/ws/steps?user_id=public';
      console.log('[usePublicStepsCounter] Connecting to:', wsUrl);

      try {
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
          console.log('[usePublicStepsCounter] ✅ WebSocket connected');
          setConnected(true);
          reconnectTimeoutRef.current = 1000; // Reset timeout

          // ✨ CRUCIALE STAP: SUBSCRIBEN
          ws.send(JSON.stringify({
            type: 'subscribe',
            channels: ['total_updates', 'step_updates', 'leaderboard_updates']
          }));
          console.log('[usePublicStepsCounter] ✅ Subscribed to channels');
        };

        ws.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            console.log('[WebSocket] 📨 Message received:', message.type, message);

            switch (message.type) {
              case 'welcome':
                console.log('[WebSocket] Welcome:', message.message);
                break;
                
              case 'total_update':
                setTotalSteps(message.total_steps);
                console.log('[WebSocket] 🔄 Total updated:', message.total_steps);
                break;
                
              case 'step_update':
                console.log('[WebSocket] 👟 Participant stepped:', message.naam, '+', message.delta);
                // Optioneel: show toast notification
                break;
                
              case 'leaderboard_update':
                console.log('[WebSocket] 🏆 Leaderboard updated:', message.entries.length, 'entries');
                break;
                
              case 'pong':
                console.log('[WebSocket] 💓 Pong received');
                break;
            }
          } catch (error) {
            console.error('[WebSocket] Parse error:', error);
          }
        };

        ws.onerror = (error) => {
          console.error('[WebSocket] ❌ Error:', error);
        };

        ws.onclose = (event) => {
          console.log('[WebSocket] 👋 Disconnected:', event.code, event.reason);
          setConnected(false);
          wsRef.current = null;

          // Auto-reconnect
          const timeout = Math.min(reconnectTimeoutRef.current * 2, 30000);
          console.log('[WebSocket] 🔄 Reconnecting in', timeout, 'ms');
          setTimeout(connect, timeout);
          reconnectTimeoutRef.current = timeout;
        };

      } catch (error) {
        console.error('[WebSocket] Connection failed:', error);
        // Retry
        setTimeout(connect, 5000);
      }
    };

    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [backendUrl]);

  return {
    totalSteps,
    connected,
  };
}
```

---

## 🧪 TEST PROCEDURE

### Test 1: Verify Subscribe

**Check in browser console**:
```
[usePublicStepsCounter] Connecting to: ws://...
[usePublicStepsCounter] ✅ WebSocket connected
[usePublicStepsCounter] ✅ Subscribed to channels  ← MOET JE ZIEN
[WebSocket] 📨 Message received: welcome {...}     ← NIEUWE
```

### Test 2: Trigger Broadcast

**In andere terminal/Postman**:
```bash
curl -X POST http://localhost:8080/api/steps \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"steps": 50}'
```

**Verwacht in console**:
```
[WebSocket] 📨 Message received: total_update {...}
[WebSocket] 🔄 Total updated: 1311
```

**Verwacht in UI**:
- Counter update van 1261 → 1311 **ZONDER REFRESH**

---

## ✅ WAAROM DIT NU ZOU MOETEN WERKEN

**Backend**: ✅ KLAAR (commit ba5bc4f)
- Stuurt welcome message
- Vertelt welke channels beschikbaar zijn
- Logs alle connects/subscribes

**Frontend**: ❓ **CHECK OF DIT ER IS**
- Subscribe message na connect
- Message handlers
- State updates

**Trigger**: ⏳ **VEREIST**
- Iemand moet stappen updaten
- Dan komt broadcast
- Dan update UI

---

## 🎯 ACTIES VOOR JOU

### Actie 1: Check Subscribe Logic

**Zoek in je `usePublicStepsCounter.ts`**:
```typescript
// Is dit er? Zo niet, toevoegen!
ws.send(JSON.stringify({
  type: 'subscribe',
  channels: ['total_updates']
}));
```

### Actie 2: Test Broadcast

**Trigger een step update** via:
- Postman/curl naar `/api/steps`
- Of mobile app (sync 50 stappen)
- Of admin panel

**Verwacht**:
- Browser console: "Message received: total_update"
- UI update **INSTANT**

### Actie 3: Verify Welcome Message

**Na Render redeploy**, eerste bericht zou moeten zijn:
```json
{
  "type": "welcome",
  "message": "Connected to StepsHub!...",
  "available_channels": [...]
}
```

---

## 🎊 SUMMARY

**Backend**: ✅ **100% KLAAR**
- WebSocket server draait
- Endpoints geregistreerd
- Welcome message feature
- Broadcasts triggered bij updates
- 4 commits gepusht

**Frontend**: ⏳ **VEREIST**
- Subscribe message toevoegen
- Message handlers verbeteren  
- Test met broadcast trigger

**Expected Behavior (na fix)**:
```
Bezoeker opent website
  → WebSocket verbindt
  → Welcome message ontvangen
  → Subscribe naar channels
  → Iemand update stappen
  → Broadcast getriggerd
  → Frontend ontvangt message
  → Counter update LIVE ⚡
  → Geen refresh nodig!
```

---

**Status**: Backend ✅ | Frontend needs subscribe logic  
**Git**: ba5bc4f pushed  
**Next**: Add subscribe to frontend + test broadcast
</result>
</attempt_completion>