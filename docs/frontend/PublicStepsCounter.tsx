/**
 * PublicStepsCounter - KANT-EN-KLAAR component voor live stappen op website
 * 
 * Gebruik:
 * ```tsx
 * <PublicStepsCounter />
 * ```
 * 
 * Features:
 * - Live stappen counter via WebSocket
 * - Auto-reconnect
 * - Fallback naar polling
 * - Geen login vereist
 */

import React, { useState, useEffect, useRef } from 'react';

interface PublicStepsCounterProps {
  apiUrl?: string;  // Default: auto-detect
}

export const PublicStepsCounter: React.FC<PublicStepsCounterProps> = ({ 
  apiUrl = window.location.origin 
}) => {
  const [totalSteps, setTotalSteps] = useState<number>(0);
  const [connected, setConnected] = useState<boolean>(false);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  // Fetch via REST API (fallback)
  const fetchStepsREST = async () => {
    try {
      const response = await fetch(`${apiUrl}/api/total-steps`);
      const data = await response.json();
      setTotalSteps(data.total_steps);
      setLastUpdate(new Date());
      console.log('[REST] Total steps:', data.total_steps);
    } catch (error) {
      console.error('[REST] Error:', error);
    }
  };

  // WebSocket connection
  useEffect(() => {
    const wsUrl = apiUrl.replace('https:', 'wss:').replace('http:', 'ws:') + '/api/ws/steps?user_id=public';
    
    const connect = () => {
      console.log('[WebSocket] Connecting to:', wsUrl);
      
      try {
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
          console.log('[WebSocket] ✅ Connected!');
          setConnected(true);

          // ✨ SUBSCRIBE - DIT IS CRUCIAAL!
          const subscribeMsg = {
            type: 'subscribe',
            channels: ['total_updates', 'step_updates', 'leaderboard_updates']
          };
          
          ws.send(JSON.stringify(subscribeMsg));
          console.log('[WebSocket] ✅ Subscribed to channels');

          // Stop polling wanneer WebSocket werkt
          if (pollingRef.current) {
            clearInterval(pollingRef.current);
            pollingRef.current = null;
          }
        };

        ws.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            console.log('[WebSocket] 📨 Received:', message.type);

            switch (message.type) {
              case 'welcome':
                console.log('[WebSocket] 👋 Welcome message:', message.message);
                break;

              case 'total_update':
                setTotalSteps(message.total_steps);
                setLastUpdate(new Date());
                console.log('[WebSocket] 🔄 Total updated:', message.total_steps);
                break;

              case 'step_update':
                console.log('[WebSocket] 👟 Step update:', message.naam, '+', message.delta);
                // Triggers automatic total_update from backend
                break;

              case 'pong':
                console.log('[WebSocket] 💓 Pong');
                break;
            }
          } catch (error) {
            console.error('[WebSocket] Parse error:', error);
          }
        };

        ws.onerror = (error) => {
          console.error('[WebSocket] ❌ Error:', error);
        };

        ws.onclose = () => {
          console.log('[WebSocket] 👋 Disconnected');
          setConnected(false);
          wsRef.current = null;

          // Start polling fallback
          if (!pollingRef.current) {
            console.log('[WebSocket] Starting polling fallback');
            pollingRef.current = setInterval(fetchStepsREST, 5000);
          }

          // Reconnect after delay
          setTimeout(connect, 5000);
        };

      } catch (error) {
        console.error('[WebSocket] Connection error:', error);
        
        // Start polling fallback
        if (!pollingRef.current) {
          pollingRef.current = setInterval(fetchStepsREST, 5000);
        }
      }
    };

    // Initial fetch
    fetchStepsREST();

    // Try WebSocket connection
    connect();

    // Cleanup
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, [apiUrl]);

  return (
    <div className="public-steps-counter">
      <div className="counter-status">
        <span className={`status-dot ${connected ? 'connected' : 'disconnected'}`}>
          {connected ? '🟢' : '🔴'}
        </span>
        <span className="status-text">
          {connected ? 'Live Updates' : 'Updates Every 5s'}
        </span>
      </div>

      <div className="counter-display">
        <div className="counter-label">Totaal Gelopen Stappen</div>
        <div className="counter-value">
          {totalSteps.toLocaleString('nl-NL')}
        </div>
        {lastUpdate && (
          <div className="counter-timestamp">
            Laatst bijgewerkt: {lastUpdate.toLocaleTimeString('nl-NL')}
          </div>
        )}
      </div>

      <style>{`
        .public-steps-counter {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          border-radius: 16px;
          padding: 32px;
          color: white;
          text-align: center;
          box-shadow: 0 10px 40px rgba(102, 126, 234, 0.4);
        }

        .counter-status {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 8px;
          margin-bottom: 24px;
          font-size: 14px;
          opacity: 0.9;
        }

        .status-dot {
          font-size: 12px;
        }

        .counter-label {
          font-size: 18px;
          opacity: 0.9;
          margin-bottom: 12px;
          font-weight: 500;
        }

        .counter-value {
          font-size: 64px;
          font-weight: bold;
          line-height: 1;
          margin-bottom: 12px;
          text-shadow: 0 2px 10px rgba(0,0,0,0.2);
        }

        .counter-timestamp {
          font-size: 12px;
          opacity: 0.7;
        }

        @media (max-width: 768px) {
          .counter-value {
            font-size: 48px;
          }
        }
      `}</style>
    </div>
  );
};

export default PublicStepsCounter;
```

---

## 🚀 DIRECT GEBRUIKEN

**In je App.tsx of HomePage.tsx**:

```tsx
import PublicStepsCounter from './components/PublicStepsCounter';

function HomePage() {
  return (
    <div>
      <h1>De Koninklijke Loop</h1>
      
      {/* ✨ LIVE STAPPEN COUNTER */}
      <PublicStepsCounter />
      
      {/* Rest van je page */}
    </div>
  );
}
```

**Klaar! Bezoekers zien nu**number updates LIVE** zonder refresh!

---

## 🔧 INSTANT FIX

**Als je bestaande `usePublicStepsCounter.ts` hebt**:

Voeg ALLEEN deze regel toe in `ws.onopen`:

```typescript
ws.onopen = () => {
  setConnected(true);
  
  // ✨ VOEG DIT TOE:
  ws.send(JSON.stringify({
    type: 'subscribe',
    channels: ['total_updates']
  }));
  
  console.log('✅ Subscribed to total_updates');
};
```

**Dat is alles!** Berichten zouden nu moeten binnen komen.

---

## ✅ VERIFICATIE CHECKLIST

Na fix zou je in console moeten zien:

```
[usePublic StepsCounter] Connecting to: ws://...
[usePublicStepsCounter] ✅ WebSocket connected
✅ Subscribed to total_updates                    ← NIEUW
[WebSocket] 📨 Received: welcome                 ← NIEUW (na deploy)
... wait for someone to update steps ...
[WebSocket] 📨 Received: total_update            ← NIEUW
[WebSocket] 🔄 Total updated: 1311              ← NIEUW
```

**En in UI**: Counter springt van 1261 → 1311 **INSTANT**!

---

**Fix**: Add subscribe message in onopen handler  
**Impact**: Immediate - berichten komen binnen  
**Status**: Copy-paste ready code provided
</result>
</attempt_completion>