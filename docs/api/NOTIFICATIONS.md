# Notifications API

API endpoints voor het beheren van gebruikersnotificaties.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

## Authentication

Alle endpoints vereisen JWT authenticatie:

```
Authorization: Bearer <your_jwt_token>
```

---

## Endpoints

### List User Notifications

Haal notificaties op voor de ingelogde gebruiker.

**Endpoint:** `GET /api/notifications`

**Authentication:** Vereist

**Query Parameters:**
- `unread` (optional): Filter ongelezen notificaties (true/false)
- `type` (optional): Filter op notificatie type
- `limit` (optional): Aantal resultaten (default: 20, max: 100)
- `offset` (optional): Aantal over te slaan (default: 0)

**Response:** `200 OK`

```json
{
  "notifications": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "type": "achievement_unlocked",
      "title": "🏆 Achievement Unlocked!",
      "message": "Je hebt de '5K Hero' achievement behaald!",
      "data": {
        "achievement_id": "uuid",
        "achievement_name": "5K Hero",
        "points": 50
      },
      "read": false,
      "read_at": null,
      "created_at": "2025-01-08T14:30:00Z"
    },
    {
      "id": "uuid",
      "user_id": "uuid",
      "type": "contact_response",
      "title": "📧 Nieuw Antwoord",
      "message": "Je hebt een antwoord ontvangen op je contactformulier",
      "data": {
        "contact_id": "uuid",
        "response_preview": "Bedankt voor je bericht..."
      },
      "read": true,
      "read_at": "2025-01-08T15:00:00Z",
      "created_at": "2025-01-08T14:00:00Z"
    }
  ],
  "total": 25,
  "unread_count": 5
}
```

---

### Get Notification by ID

**Endpoint:** `GET /api/notifications/:id`

**Authentication:** Vereist

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "type": "badge_awarded",
  "title": "🎖️ Nieuwe Badge!",
  "message": "Je hebt de 'Marathon Master' badge verdiend!",
  "data": {
    "badge_id": "uuid",
    "badge_name": "Marathon Master",
    "badge_tier": "gold",
    "badge_image": "https://res.cloudinary.com/..."
  },
  "read": false,
  "created_at": "2025-01-08T14:30:00Z"
}
```

---

### Mark as Read

**Endpoint:** `PUT /api/notifications/:id/read`

**Authentication:** Vereist

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Notification marked as read",
  "notification_id": "uuid",
  "read_at": "2025-01-08T16:00:00Z"
}
```

---

### Mark All as Read

**Endpoint:** `PUT /api/notifications/read-all`

**Authentication:** Vereist

**Request Body (optional):**

```json
{
  "type": "achievement_unlocked"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "All notifications marked as read",
  "count": 12
}
```

---

### Delete Notification

**Endpoint:** `DELETE /api/notifications/:id`

**Authentication:** Vereist

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Notification deleted successfully"
}
```

---

### Delete All Read Notifications

**Endpoint:** `DELETE /api/notifications/read`

**Authentication:** Vereist

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Read notifications deleted",
  "count": 8
}
```

---

### Create Notification (Admin)

**Endpoint:** `POST /api/notifications`

**Authentication:** Vereist - `notification:write` permissie

**Request Body:**

```json
{
  "user_id": "uuid",
  "type": "system_announcement",
  "title": "Belangrijk Bericht",
  "message": "Het evenement is verplaatst naar volgende week",
  "data": {
    "announcement_id": "uuid",
    "priority": "high"
  },
  "send_email": true
}
```

**Notification Types:**
- `achievement_unlocked` - Achievement behaald
- `badge_awarded` - Badge toegekend
- `contact_response` - Antwoord op contactformulier
- `registration_confirmed` - Registratie bevestigd
- `event_reminder` - Evenement herinnering
- `chat_mention` - Genoemd in chat
- `system_announcement` - Systeem aankondiging
- `custom` - Custom notificatie

**Response:** `201 Created`

---

### Broadcast Notification (Admin)

Verstuur notificatie naar meerdere gebruikers.

**Endpoint:** `POST /api/notifications/broadcast`

**Authentication:** Vereist - `notification:write` permissie

**Request Body:**

```json
{
  "user_ids": ["uuid1", "uuid2", "uuid3"],
  "type": "system_announcement",
  "title": "Evenement Update",
  "message": "Het evenement start over 1 week!",
  "data": {
    "event_id": "uuid",
    "priority": "medium"
  },
  "send_email": false
}
```

**Alternative - All Users:**

```json
{
  "broadcast_all": true,
  "type": "system_announcement",
  "title": "Platform Onderhoud",
  "message": "Het platform zal vanavond om 22:00 offline zijn voor onderhoud"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Broadcast sent successfully",
  "recipients_count": 150,
  "failed_count": 0
}
```

---

## Notification Preferences

### Get User Preferences

**Endpoint:** `GET /api/notifications/preferences`

**Authentication:** Vereist

**Response:** `200 OK`

```json
{
  "user_id": "uuid",
  "preferences": {
    "email_notifications": true,
    "push_notifications": true,
    "achievement_notifications": true,
    "chat_notifications": true,
    "event_notifications": true,
    "newsletter_notifications": true
  }
}
```

---

### Update Preferences

**Endpoint:** `PUT /api/notifications/preferences`

**Authentication:** Vereist

**Request Body:**

```json
{
  "email_notifications": false,
  "push_notifications": true,
  "achievement_notifications": true
}
```

**Response:** `200 OK`

---

## Real-time Notifications

### Server-Sent Events (SSE)

Voor real-time notificatie updates zonder WebSocket:

**Endpoint:** `GET /api/notifications/stream`

**Authentication:** JWT via query parameter

**Connection:**

```javascript
const token = localStorage.getItem('access_token');
const eventSource = new EventSource(
  `http://localhost:8080/api/notifications/stream?token=${token}`
);

eventSource.onmessage = (event) => {
  const notification = JSON.parse(event.data);
  showNotification(notification);
};

eventSource.onerror = (error) => {
  console.error('SSE error:', error);
  eventSource.close();
};
```

**Event Format:**

```json
{
  "id": "uuid",
  "type": "achievement_unlocked",
  "title": "Achievement Unlocked!",
  "message": "Je hebt een nieuwe achievement behaald",
  "timestamp": "2025-01-08T14:30:00Z"
}
```

---

## Integration Examples

### React Notification Center

```typescript
import { useState, useEffect } from 'react';
import api from './api';

export function NotificationCenter() {
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    loadNotifications();
    
    // Setup SSE for real-time updates
    const token = localStorage.getItem('access_token');
    const eventSource = new EventSource(
      `/api/notifications/stream?token=${token}`
    );

    eventSource.onmessage = (event) => {
      const notification = JSON.parse(event.data);
      setNotifications(prev => [notification, ...prev]);
      setUnreadCount(prev => prev + 1);
    };

    return () => eventSource.close();
  }, []);

  const loadNotifications = async () => {
    const response = await api.get('/api/notifications');
    setNotifications(response.data.notifications);
    setUnreadCount(response.data.unread_count);
  };

  const markAsRead = async (id: string) => {
    await api.put(`/api/notifications/${id}/read`);
    setNotifications(prev =>
      prev.map(n => n.id === id ? { ...n, read: true } : n)
    );
    setUnreadCount(prev => Math.max(0, prev - 1));
  };

  const markAllAsRead = async () => {
    await api.put('/api/notifications/read-all');
    setNotifications(prev => prev.map(n => ({ ...n, read: true })));
    setUnreadCount(0);
  };

  return (
    <div className="notification-center">
      <div className="header">
        <h2>Notificaties</h2>
        {unreadCount > 0 && (
          <span className="badge">{unreadCount}</span>
        )}
        <button onClick={markAllAsRead}>Mark All Read</button>
      </div>
      
      <div className="notifications-list">
        {notifications.map(notification => (
          <NotificationItem
            key={notification.id}
            notification={notification}
            onRead={() => markAsRead(notification.id)}
          />
        ))}
      </div>
    </div>
  );
}
```

---

### Vue Notification Toast

```vue
<template>
  <Transition name="slide">
    <div v-if="visible" class="notification-toast" :class="type">
      <div class="icon">{{ getIcon(type) }}</div>
      <div class="content">
        <h4>{{ title }}</h4>
        <p>{{ message }}</p>
      </div>
      <button @click="close" class="close-btn">×</button>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const props = defineProps({
  notification: Object
});

const visible = ref(false);
const type = ref(props.notification.type);
const title = ref(props.notification.title);
const message = ref(props.notification.message);

const getIcon = (type) => {
  const icons = {
    achievement_unlocked: '🏆',
    badge_awarded: '🎖️',
    contact_response: '📧',
    event_reminder: '📅',
    chat_mention: '💬',
    system_announcement: '📢'
  };
  return icons[type] || '🔔';
};

onMounted(() => {
  visible.value = true;
  setTimeout(() => {
    visible.value = false;
  }, 5000);
});

const close = () => {
  visible.value = false;
};
</script>

<style scoped>
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>
```

---

## Notification Data Structure

### Common Data Fields

Elke notificatie type heeft specifieke data:

#### Achievement Unlocked

```json
{
  "type": "achievement_unlocked",
  "data": {
    "achievement_id": "uuid",
    "achievement_name": "5K Hero",
    "points": 50,
    "icon_url": "https://..."
  }
}
```

#### Badge Awarded

```json
{
  "type": "badge_awarded",
  "data": {
    "badge_id": "uuid",
    "badge_name": "Marathon Master",
    "badge_tier": "gold",
    "badge_image": "https://..."
  }
}
```

#### Contact Response

```json
{
  "type": "contact_response",
  "data": {
    "contact_id": "uuid",
    "response_preview": "Bedankt voor je bericht...",
    "responder_name": "Admin Name"
  }
}
```

#### Chat Mention

```json
{
  "type": "chat_mention",
  "data": {
    "channel_id": "uuid",
    "channel_name": "Team Chat",
    "message_id": "uuid",
    "mentioned_by": "uuid",
    "mentioned_by_name": "Jane Doe",
    "message_preview": "@john wat denk jij?"
  }
}
```

---

## Telegram Integration

Notificaties kunnen ook via Telegram worden verzonden als dit is geconfigureerd.

### Telegram Priority Levels

Notificaties worden gefilterd op basis van prioriteit:

| Priority | Telegram | Email | In-App |
|----------|----------|-------|--------|
| `low` | ❌ | ❌ | ✅ |
| `medium` | ✅ | ❌ | ✅ |
| `high` | ✅ | ✅ | ✅ |
| `critical` | ✅ | ✅ | ✅ |

**Configuration:**

```env
NOTIFICATION_MIN_PRIORITY=medium
NOTIFICATION_THROTTLE=15m
```

---

## Notification Throttling

Voorkom spam door vergelijkbare notificaties te throttlen:

**Rules:**
- Zelfde type binnen 15 minuten: Gecombineerd
- Achievement unlocks: Max 1 per minuut
- Chat mentions: Max 5 per 5 minuten
- System announcements: Geen throttling

**Example Combined Notification:**

```json
{
  "type": "achievements_unlocked",
  "title": "🏆 3 Achievements Unlocked!",
  "message": "Je hebt 3 nieuwe achievements behaald",
  "data": {
    "achievements": [
      {"name": "First Step", "points": 10},
      {"name": "1K Steps", "points": 25},
      {"name": "5K Steps", "points": 50}
    ],
    "total_points": 85
  }
}
```

---

## Push Notifications (Future)

Gereserveerd voor toekomstige implementatie van web push notifications.

**Planned Endpoints:**
- `POST /api/notifications/subscribe` - Subscribe to push
- `DELETE /api/notifications/unsubscribe` - Unsubscribe from push
- `POST /api/notifications/test-push` - Test push notification

---

## Best Practices

### Frontend Implementation

1. **Polling Interval**: Poll elke 30-60 seconden voor nieuwe notificaties
2. **Badge Count**: Toon unread count in UI
3. **Auto-dismiss**: Verberg toast notificaties na 5 seconden
4. **Sound**: Speel geluid af voor high-priority notificaties
5. **Persistence**: Bewaar gelezen status lokaal

### Backend Generation

```go
// Create notification
notification := &models.Notification{
    UserID:  userID,
    Type:    "achievement_unlocked",
    Title:   "🏆 Achievement Unlocked!",
    Message: "You've unlocked the 5K Hero achievement!",
    Data: map[string]interface{}{
        "achievement_id": achievementID,
        "points": 50,
    },
}

notificationService.Create(ctx, notification)
```

---

## Error Codes

| Code | HTTP Status | Beschrijving |
|------|-------------|--------------|
| `NOTIFICATION_NOT_FOUND` | 404 | Notificatie niet gevonden |
| `UNAUTHORIZED` | 401 | Geen toegang tot notificatie |
| `INVALID_TYPE` | 400 | Ongeldig notificatie type |
| `BROADCAST_FAILED` | 500 | Broadcast mislukt |

---

## Testing

```bash
# List notifications
curl http://localhost:8080/api/notifications \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get unread notifications
curl "http://localhost:8080/api/notifications?unread=true" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Mark as read
curl -X PUT http://localhost:8080/api/notifications/NOTIFICATION_ID/read \
  -H "Authorization: Bearer YOUR_TOKEN"

# Create notification (admin)
curl -X POST http://localhost:8080/api/notifications \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID",
    "type": "system_announcement",
    "title": "Test Notification",
    "message": "This is a test"
  }'
```

---

Voor meer informatie:
- [Authentication API](./AUTHENTICATION.md)
- [Chat API](./CHAT.md)
- [API Overview](./README.md)