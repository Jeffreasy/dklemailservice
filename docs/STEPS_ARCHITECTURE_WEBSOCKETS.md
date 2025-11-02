# Stappenteller Backend Architectuur & WebSocket Implementatie

## 📋 Overzicht

Dit document beschrijft de volledige backend architectuur van de stappenteller functionaliteit en het implementatieplan voor real-time WebSocket updates.

**Laatste Update**: 2025-01-02  
**Status**: 🟢 Production Ready (HTTP) | 🟡 WebSocket Planning

---

## 🏗️ Huidige Architectuur

### 1. Database Schema

#### Aanmeldingen Tabel
```sql
CREATE TABLE aanmeldingen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    naam VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefoon VARCHAR(50),
    rol VARCHAR(100),
    afstand VARCHAR(50),          -- Route: "6 KM", "10 KM", "15 KM", "20 KM"
    ondersteuning TEXT,
    bijzonderheden TEXT,
    terms BOOLEAN NOT NULL,
    email_verzonden BOOLEAN DEFAULT false,
    email_verzonden_op TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'nieuw',
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP WITH TIME ZONE,
    notities TEXT,
    test_mode BOOLEAN DEFAULT false,
    steps INTEGER DEFAULT 0,      -- 🎯 Stappen counter
    gebruiker_id UUID,            -- Link naar gebruikersaccount
    
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_gebruiker_id (gebruiker_id)
);
```

#### Route Funds Tabel
```sql
CREATE TABLE route_funds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route VARCHAR(50) NOT NULL UNIQUE,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Standaard data
INSERT INTO route_funds (route, amount) VALUES
    ('6 KM', 50),
    ('10 KM', 75),
    ('15 KM', 100),
    ('20 KM', 125);
```

#### Gamification Tabellen (V1_51)
```sql
-- Badges
CREATE TABLE badges (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    icon_url VARCHAR(500),
    criteria JSONB NOT NULL DEFAULT '{}',
    points INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0
);

-- Verdiende Achievements
CREATE TABLE participant_achievements (
    id UUID PRIMARY KEY,
    participant_id UUID REFERENCES aanmeldingen(id) ON DELETE CASCADE,
    badge_id UUID REFERENCES badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(participant_id, badge_id)
);

-- Leaderboard View
CREATE VIEW leaderboard_view AS
SELECT 
    a.id,
    a.naam,
    a.afstand as route,
    a.steps,
    COALESCE(SUM(b.points), 0) as achievement_points,
    a.steps + COALESCE(SUM(b.points), 0) as total_score,
    RANK() OVER (ORDER BY (a.steps + COALESCE(SUM(b.points), 0)) DESC) as rank,
    COUNT(pa.id) as badge_count
FROM aanmeldingen a
LEFT JOIN participant_achievements pa ON a.id = pa.participant_id
LEFT JOIN badges b ON pa.badge_id = b.id AND b.is_active = true
GROUP BY a.id, a.naam, a.afstand, a.steps
ORDER BY total_score DESC;
```

### 2. Models Layer

**Bestand**: [`models/aanmelding.go`](../models/aanmelding.go)

```go
type Aanmelding struct {
    ID            string     `json:"id" gorm:"primaryKey;type:uuid"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
    Naam          string     `json:"naam"`
    Email         string     `json:"email"`
    Telefoon      string     `json:"telefoon"`
    Rol           string     `json:"rol"`
    Afstand       string     `json:"afstand"`        // Route
    Steps         int        `json:"steps"`          // 🎯 Stappen counter
    GebruikerID   *string    `json:"gebruiker_id"`   // User account link
    // ... andere velden
}
```

### 3. Repository Layer

**Bestand**: [`repository/aanmelding_repository.go`](../repository/aanmelding_repository.go)

```go
type AanmeldingRepository interface {
    Create(ctx context.Context, aanmelding *Aanmelding) error
    GetByID(ctx context.Context, id string) (*Aanmelding, error)
    Update(ctx context.Context, aanmelding *Aanmelding) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*Aanmelding, error)
    FindByEmail(ctx context.Context, email string) ([]*Aanmelding, error)
    FindByStatus(ctx context.Context, status string) ([]*Aanmelding, error)
}
```

**Key Features**:
- ✅ Context-aware met timeouts
- ✅ Error handling via `handleError()`
- ✅ PostgreSQL optimized queries
- ✅ Index support voor snelle lookups

### 4. Service Layer

**Bestand**: [`services/steps_service.go`](../services/steps_service.go)

```go
type StepsService struct {
    db             *gorm.DB
    aanmeldingRepo repository.AanmeldingRepository
    routeFundRepo  repository.RouteFundRepository
}
```

#### Core Functionaliteit

**Stappen Updates**:
```go
// UpdateSteps werkt stappen bij voor een deelnemer (delta)
// - Haalt deelnemer op via ID
// - Voegt delta toe aan huidige stappen
// - Voorkomt negatieve stappen (minimum = 0)
// - Slaat wijzigingen op
func (s *StepsService) UpdateSteps(participantID string, deltaSteps int) (*models.Aanmelding, error)

// UpdateStepsByUserID werkt stappen bij via gebruiker ID
// - Voor ingelogde deelnemers die hun eigen stappen bijwerken
// - Zoekt eerst aanmelding via gebruiker_id
func (s *StepsService) UpdateStepsByUserID(userID string, deltaSteps int) (*models.Aanmelding, error)
```

**Dashboard Data**:
```go
// GetParticipantDashboard haalt dashboard data op
// Retourneert:
// - Deelnemer gegevens (naam, email, stappen)
// - Berekende allocated funds gebaseerd op route
func (s *StepsService) GetParticipantDashboard(participantID string) (*models.Aanmelding, int, error)

// GetParticipantDashboardByUserID - via gebruiker ID
func (s *StepsService) GetParticipantDashboardByUserID(userID string) (*models.Aanmelding, int, error)
```

**Statistieken**:
```go
// GetTotalSteps - totaal aantal stappen
// - year = 0: ALLE stappen van ALLE deelnemers
// - year > 0: filter op jaar van aanmelding
func (s *StepsService) GetTotalSteps(year int) (int, error)

// GetFundsDistributionProportional - fondsverdeling per route
// - Berekent proportioneel gebaseerd op aantal deelnemers
func (s *StepsService) GetFundsDistributionProportional() (map[string]int, int, error)
```

**Admin Functies**:
```go
// Route Fund beheer
func (s *StepsService) GetRouteFunds() ([]*models.RouteFund, error)
func (s *StepsService) CreateRouteFund(route string, amount int) (*models.RouteFund, error)
func (s *StepsService) UpdateRouteFund(route string, amount int) (*models.RouteFund, error)
func (s *StepsService) DeleteRouteFund(route string) error
```

### 5. Handler Layer

**Bestand**: [`handlers/steps_handler.go`](../handlers/steps_handler.go)

#### REST API Endpoints

| Method | Endpoint | Auth | Permission | Beschrijving |
|--------|----------|------|------------|--------------|
| POST | `/api/steps` | ✅ | `steps:write` | Update stappen voor ingelogde deelnemer |
| POST | `/api/steps/:id` | ✅ | `steps:write` | Update stappen voor specifieke deelnemer (admin) |
| GET | `/api/participant/dashboard` | ✅ | `steps:read` | Dashboard voor ingelogde deelnemer |
| GET | `/api/participant/:id/dashboard` | ✅ | `steps:read` | Dashboard voor specifieke deelnemer (admin) |
| GET | `/api/total-steps` | 🔓 | - | Totaal aantal stappen (PUBLIC) |
| GET | `/api/funds-distribution` | 🔓 | - | Fondsverdeling (PUBLIC) |
| GET | `/api/steps/admin/route-funds` | ✅ | `steps:write` | Alle route funds ophalen |
| POST | `/api/steps/admin/route-funds` | ✅ | `steps:write` | Route fund aanmaken |
| PUT | `/api/steps/admin/route-funds/:route` | ✅ | `steps:write` | Route fund bijwerken |
| DELETE | `/api/steps/admin/route-funds/:route` | ✅ | `steps:write` | Route fund verwijderen |

#### Request/Response Voorbeelden

**POST /api/steps** - Update stappen:
```json
// Request
{
    "steps": 1500  // Delta waarde (kan negatief zijn)
}

// Response
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "naam": "John Doe",
    "email": "john@example.com",
    "afstand": "10 KM",
    "steps": 3500,  // Nieuwe totaal
    "created_at": "2024-03-20T15:04:05Z",
    "updated_at": "2024-03-20T16:00:00Z"
}
```

**GET /api/participant/dashboard** - Dashboard data:
```json
{
    "steps": 3500,
    "route": "10 KM",
    "allocatedFunds": 75,
    "naam": "John Doe",
    "email": "john@example.com"
}
```

**GET /api/total-steps?year=2025** - Totaal stappen:
```json
{
    "total_steps": 250000,
    "year": 2025
}
```

### 6. Security & RBAC

**Permissions** (V1_45 migratie):
```sql
INSERT INTO permissions (resource, action, description) VALUES
('steps', 'read', 'Stappen data bekijken'),
('steps', 'write', 'Stappen bijwerken');
```

**Role Assignments**:
- **Admin**: `steps:read`, `steps:write` (alle endpoints)
- **Staff**: `steps:read`, `steps:write` (alle endpoints)
- **Deelnemer**: `steps:read`, `steps:write` (alleen eigen data)
- **Public**: `/api/total-steps`, `/api/funds-distribution` (geen auth)

**Authorization Flow**:
```
Request → JWT Middleware → Permission Middleware → Handler
           (userID)         (checks permission)       (logic)
```

---

## 🔌 WebSocket Implementatie Plan

### Huidige WebSocket Infrastructure

**Bestand**: [`services/websocket_hub.go`](../services/websocket_hub.go)

De applicatie heeft al een WebSocket Hub voor chat functionaliteit:

```go
type Hub struct {
    Clients      map[*Client]bool
    Broadcast    chan []byte
    Register     chan *Client
    Unregister   chan *Client
    ChatService  ChatService
    GetChannelHub func(channelID string) *Hub
}

type Client struct {
    Hub    *Hub
    Conn   *websocket.Conn
    Send   chan []byte
    UserID string
}
```

**Features**:
- ✅ Client registration/unregistration
- ✅ Broadcast messaging
- ✅ User presence tracking
- ✅ Typing indicators
- ✅ Channel-based routing

### Nieuwe Steps WebSocket Architectuur

#### 1. StepsHub - Dedicated Hub voor Stappen Updates

**Nieuw bestand**: `services/steps_hub.go`

```go
package services

import (
    "encoding/json"
    "github.com/gofiber/websocket/v2"
)

// StepsHub beheert WebSocket connecties voor stappen updates
type StepsHub struct {
    // Registered clients
    Clients map[*StepsClient]bool
    
    // Broadcast channels
    StepUpdate      chan *StepUpdateMessage
    TotalUpdate     chan *TotalUpdateMessage
    LeaderboardUpdate chan *LeaderboardUpdateMessage
    
    // Registration channels
    Register   chan *StepsClient
    Unregister chan *StepsClient
    
    // Services
    StepsService *StepsService
    GamificationService *GamificationService
}

// StepsClient represents a WebSocket client
type StepsClient struct {
    Hub          *StepsHub
    Conn         *websocket.Conn
    Send         chan []byte
    UserID       string
    ParticipantID string
    Subscriptions map[string]bool  // Subscription types
}

// Message Types
type MessageType string

const (
    MessageTypeStepUpdate       MessageType = "step_update"
    MessageTypeTotalUpdate      MessageType = "total_update"
    MessageTypeLeaderboardUpdate MessageType = "leaderboard_update"
    MessageTypeBadgeEarned      MessageType = "badge_earned"
    MessageTypeSubscribe        MessageType = "subscribe"
    MessageTypeUnsubscribe      MessageType = "unsubscribe"
)

// StepUpdateMessage - individuele deelnemer update
type StepUpdateMessage struct {
    Type          MessageType `json:"type"`
    ParticipantID string      `json:"participant_id"`
    Naam          string      `json:"naam"`
    Steps         int         `json:"steps"`
    Delta         int         `json:"delta"`
    Route         string      `json:"route"`
    AllocatedFunds int        `json:"allocated_funds"`
    Timestamp     int64       `json:"timestamp"`
}

// TotalUpdateMessage - totaal stappen update
type TotalUpdateMessage struct {
    Type       MessageType `json:"type"`
    TotalSteps int         `json:"total_steps"`
    Year       int         `json:"year"`
    Timestamp  int64       `json:"timestamp"`
}

// LeaderboardUpdateMessage - leaderboard update
type LeaderboardUpdateMessage struct {
    Type      MessageType `json:"type"`
    TopN      int         `json:"top_n"`
    Entries   []LeaderboardEntry `json:"entries"`
    Timestamp int64       `json:"timestamp"`
}

type LeaderboardEntry struct {
    Rank            int    `json:"rank"`
    ParticipantID   string `json:"participant_id"`
    Naam            string `json:"naam"`
    Steps           int    `json:"steps"`
    AchievementPoints int  `json:"achievement_points"`
    TotalScore      int    `json:"total_score"`
    Route           string `json:"route"`
}

// BadgeEarnedMessage - badge verdiend notificatie
type BadgeEarnedMessage struct {
    Type          MessageType `json:"type"`
    ParticipantID string      `json:"participant_id"`
    BadgeName     string      `json:"badge_name"`
    BadgeIcon     string      `json:"badge_icon"`
    Points        int         `json:"points"`
    Timestamp     int64       `json:"timestamp"`
}
```

#### 2. Hub Run Loop

```go
func NewStepsHub(stepsService *StepsService, gamificationService *GamificationService) *StepsHub {
    return &StepsHub{
        Clients:               make(map[*StepsClient]bool),
        StepUpdate:            make(chan *StepUpdateMessage, 256),
        TotalUpdate:           make(chan *TotalUpdateMessage, 256),
        LeaderboardUpdate:     make(chan *LeaderboardUpdateMessage, 256),
        Register:              make(chan *StepsClient),
        Unregister:            make(chan *StepsClient),
        StepsService:          stepsService,
        GamificationService:   gamificationService,
    }
}

func (h *StepsHub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.Clients[client] = true
            
        case client := <-h.Unregister:
            if _, ok := h.Clients[client]; ok {
                delete(h.Clients, client)
                close(client.Send)
            }
            
        case update := <-h.StepUpdate:
            h.broadcastStepUpdate(update)
            
        case update := <-h.TotalUpdate:
            h.broadcastTotalUpdate(update)
            
        case update := <-h.LeaderboardUpdate:
            h.broadcastLeaderboardUpdate(update)
        }
    }
}

func (h *StepsHub) broadcastStepUpdate(update *StepUpdateMessage) {
    message, _ := json.Marshal(update)
    
    for client := range h.Clients {
        // Check if client is subscribed to step_updates
        if client.Subscriptions["step_updates"] {
            select {
            case client.Send <- message:
            default:
                close(client.Send)
                delete(h.Clients, client)
            }
        }
    }
}

func (h *StepsHub) broadcastTotalUpdate(update *TotalUpdateMessage) {
    message, _ := json.Marshal(update)
    
    for client := range h.Clients {
        if client.Subscriptions["total_updates"] {
            select {
            case client.Send <- message:
            default:
                close(client.Send)
                delete(h.Clients, client)
            }
        }
    }
}

func (h *StepsHub) broadcastLeaderboardUpdate(update *LeaderboardUpdateMessage) {
    message, _ := json.Marshal(update)
    
    for client := range h.Clients {
        if client.Subscriptions["leaderboard_updates"] {
            select {
            case client.Send <- message:
            default:
                close(client.Send)
                delete(h.Clients, client)
            }
        }
    }
}
```

#### 3. Client Read/Write Pumps

```go
func (c *StepsClient) ReadPump() {
    defer func() {
        c.Hub.Unregister <- c
        c.Conn.Close()
    }()
    
    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }
        
        var msg map[string]interface{}
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }
        
        msgType, ok := msg["type"].(string)
        if !ok {
            continue
        }
        
        switch MessageType(msgType) {
        case MessageTypeSubscribe:
            c.handleSubscribe(msg)
        case MessageTypeUnsubscribe:
            c.handleUnsubscribe(msg)
        }
    }
}

func (c *StepsClient) WritePump() {
    defer func() {
        c.Conn.Close()
    }()
    
    for message := range c.Send {
        err := c.Conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            return
        }
    }
}

func (c *StepsClient) handleSubscribe(msg map[string]interface{}) {
    if channels, ok := msg["channels"].([]interface{}); ok {
        for _, ch := range channels {
            if channel, ok := ch.(string); ok {
                c.Subscriptions[channel] = true
            }
        }
    }
}

func (c *StepsClient) handleUnsubscribe(msg map[string]interface{}) {
    if channels, ok := msg["channels"].([]interface{}); ok {
        for _, ch := range channels {
            if channel, ok := ch.(string); ok {
                delete(c.Subscriptions, channel)
            }
        }
    }
}
```

#### 4. WebSocket Handler

**Nieuw bestand**: `handlers/steps_websocket_handler.go`

```go
package handlers

import (
    "dklautomationgo/services"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
)

type StepsWebSocketHandler struct {
    stepsHub    *services.StepsHub
    authService services.AuthService
}

func NewStepsWebSocketHandler(
    stepsHub *services.StepsHub,
    authService services.AuthService,
) *StepsWebSocketHandler {
    return &StepsWebSocketHandler{
        stepsHub:    stepsHub,
        authService: authService,
    }
}

// RegisterRoutes registreert WebSocket routes
func (h *StepsWebSocketHandler) RegisterRoutes(app *fiber.App) {
    // WebSocket upgrade middleware
    app.Use("/ws/steps", func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })
    
    // WebSocket endpoint
    app.Get("/ws/steps", websocket.New(h.HandleWebSocket))
}

func (h *StepsWebSocketHandler) HandleWebSocket(c *websocket.Conn) {
    // Extract user info from query params or headers
    // (Token validation should happen before WebSocket upgrade)
    userID := c.Query("user_id")
    participantID := c.Query("participant_id")
    
    client := &services.StepsClient{
        Hub:           h.stepsHub,
        Conn:          c,
        Send:          make(chan []byte, 256),
        UserID:        userID,
        ParticipantID: participantID,
        Subscriptions: make(map[string]bool),
    }
    
    h.stepsHub.Register <- client
    
    go client.WritePump()
    client.ReadPump()
}
```

#### 5. Integratie met StepsService

**Aanpassing**: `services/steps_service.go`

```go
type StepsService struct {
    db             *gorm.DB
    aanmeldingRepo repository.AanmeldingRepository
    routeFundRepo  repository.RouteFundRepository
    stepsHub       *StepsHub  // ✨ Nieuwe field
}

func (s *StepsService) UpdateSteps(participantID string, deltaSteps int) (*models.Aanmelding, error) {
    // Bestaande logica...
    participant, err := s.aanmeldingRepo.GetByID(context.TODO(), participantID)
    if err != nil {
        return nil, fmt.Errorf("deelnemer niet gevonden: %w", err)
    }
    
    newSteps := participant.Steps + deltaSteps
    if newSteps < 0 {
        newSteps = 0
    }
    participant.Steps = newSteps
    
    if err := s.aanmeldingRepo.Update(context.TODO(), participant); err != nil {
        return nil, fmt.Errorf("kon stappen niet bijwerken: %w", err)
    }
    
    // ✨ NIEUWE: Broadcast WebSocket update
    if s.stepsHub != nil {
        allocatedFunds := s.CalculateAllocatedFunds(participant.Afstand)
        
        s.stepsHub.StepUpdate <- &StepUpdateMessage{
            Type:           MessageTypeStepUpdate,
            ParticipantID:  participant.ID,
            Naam:           participant.Naam,
            Steps:          participant.Steps,
            Delta:          deltaSteps,
            Route:          participant.Afstand,
            AllocatedFunds: allocatedFunds,
            Timestamp:      time.Now().Unix(),
        }
        
        // Update totaal stappen
        totalSteps, _ := s.GetTotalSteps(0)
        s.stepsHub.TotalUpdate <- &TotalUpdateMessage{
            Type:       MessageTypeTotalUpdate,
            TotalSteps: totalSteps,
            Year:       0,
            Timestamp:  time.Now().Unix(),
        }
        
        // Update leaderboard (top 10)
        go s.updateLeaderboard()
    }
    
    return participant, nil
}

func (s *StepsService) updateLeaderboard() {
    // Query leaderboard_view
    var entries []LeaderboardEntry
    s.db.Table("leaderboard_view").
        Limit(10).
        Find(&entries)
    
    if s.stepsHub != nil {
        s.stepsHub.LeaderboardUpdate <- &LeaderboardUpdateMessage{
            Type:      MessageTypeLeaderboardUpdate,
            TopN:      10,
            Entries:   entries,
            Timestamp: time.Now().Unix(),
        }
    }
}
```

#### 6. Frontend WebSocket Client

**JavaScript/TypeScript Voorbeeld**:

```typescript
class StepsWebSocketClient {
    private ws: WebSocket | null = null;
    private reconnectTimeout: number = 1000;
    private maxReconnectTimeout: number = 30000;
    private callbacks: Map<string, Function[]> = new Map();
    
    constructor(
        private url: string,
        private token: string,
        private userId: string,
        private participantId?: string
    ) {}
    
    connect(): void {
        const queryParams = new URLSearchParams({
            user_id: this.userId,
            ...(this.participantId && { participant_id: this.participantId })
        });
        
        this.ws = new WebSocket(`${this.url}?${queryParams}`);
        
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.reconnectTimeout = 1000;
            
            // Subscribe to channels
            this.subscribe(['step_updates', 'total_updates', 'leaderboard_updates']);
        };
        
        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleMessage(message);
            } catch (error) {
                console.error('Failed to parse message:', error);
            }
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket disconnected, reconnecting...');
            this.reconnect();
        };
    }
    
    private reconnect(): void {
        setTimeout(() => {
            this.reconnectTimeout = Math.min(
                this.reconnectTimeout * 2,
                this.maxReconnectTimeout
            );
            this.connect();
        }, this.reconnectTimeout);
    }
    
    subscribe(channels: string[]): void {
(!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            console.warn('Cannot subscribe, WebSocket not connected');
            return;
        }
        
        this.ws.send(JSON.stringify({
            type: 'subscribe',
            channels: channels
        }));
    }
    
    unsubscribe(channels: string[]): void {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            return;
        }
        
        this.ws.send(JSON.stringify({
            type: 'unsubscribe',
            channels: channels
        }));
    }
    
    on(messageType: string, callback: Function): void {
        if (!this.callbacks.has(messageType)) {
            this.callbacks.set(messageType, []);
        }
        this.callbacks.get(messageType)!.push(callback);
    }
    
    private handleMessage(message: any): void {
        const callbacks = this.callbacks.get(message.type);
        if (callbacks) {
            callbacks.forEach(callback => callback(message));
        }
    }
    
    disconnect(): void {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }
}

// Gebruik voorbeeld
const client = new StepsWebSocketClient(
    'ws://localhost:8080/ws/steps',
    'jwt-token-here',
    'user-id-123',
    'participant-id-456'
);

// Event listeners
client.on('step_update', (data) => {
    console.log('Step update:', data);
    updateUI(data);
});

client.on('total_update', (data) => {
    console.log('Total steps:', data.total_steps);
    updateTotalDisplay(data.total_steps);
});

client.on('leaderboard_update', (data) => {
    console.log('Leaderboard:', data.entries);
    updateLeaderboard(data.entries);
});

client.on('badge_earned', (data) => {
    console.log('Badge earned!', data.badge_name);
    showBadgeNotification(data);
});

client.connect();
```

**React Hook Voorbeeld**:

```typescript
import { useEffect, useState, useRef } from 'react';

interface StepUpdate {
    type: string;
    participant_id: string;
    naam: string;
    steps: number;
    delta: number;
    route: string;
    allocated_funds: number;
    timestamp: number;
}

interface TotalUpdate {
    type: string;
    total_steps: number;
    year: number;
    timestamp: number;
}

export function useStepsWebSocket(userId: string, participantId?: string) {
    const [connected, setConnected] = useState(false);
    const [latestUpdate, setLatestUpdate] = useState<StepUpdate | null>(null);
    const [totalSteps, setTotalSteps] = useState<number>(0);
    const wsRef = useRef<StepsWebSocketClient | null>(null);
    
    useEffect(() => {
        const token = localStorage.getItem('token');
        if (!token || !userId) return;
        
        const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws/steps`;
        const client = new StepsWebSocketClient(wsUrl, token, userId, participantId);
        
        client.on('step_update', (data: StepUpdate) => {
            setLatestUpdate(data);
        });
        
        client.on('total_update', (data: TotalUpdate) => {
            setTotalSteps(data.total_steps);
        });
        
        wsRef.current = client;
        client.connect();
        setConnected(true);
        
        return () => {
            client.disconnect();
            setConnected(false);
        };
    }, [userId, participantId]);
    
    return {
        connected,
        latestUpdate,
        totalSteps,
        subscribe: (channels: string[]) => wsRef.current?.subscribe(channels),
        unsubscribe: (channels: string[]) => wsRef.current?.unsubscribe(channels)
    };
}
```

---

## 📊 Data Flow Diagram

```
┌────────────────────────────────────────────────────────────────┐
│                      Frontend Clients                           │
│ (Web App, Mobile App, Admin Dashboard, Public Leaderboard)     │
└────────────┬────────────────────────────────┬──────────────────┘
             │                                │
             │ HTTP REST                      │ WebSocket
             │ (POST /api/steps)              │ (ws://*/ws/steps)
             ▼                                ▼
┌────────────────────────────────────────────────────────────────┐
│                      Handler Layer                              │
│  ┌──────────────────┐      ┌──────────────────────────┐       │
│  │ StepsHandler     │      │ StepsWebSocketHandler     │       │
│  │ - UpdateSteps()  │      │ - HandleWebSocket()       │       │
│  └────────┬─────────┘      └────────┬─────────────────┘       │
└───────────┼──────────────────────────┼─────────────────────────┘
            │                          │ Register Client
            ▼                          ▼
┌────────────────────────────────────────────────────────────────┐
│                      Service Layer                              │
│  ┌──────────────────┐      ┌──────────────────────────┐       │
│  │ StepsService     │◄─────┤ StepsHub                 │       │
│  │ - UpdateSteps()  │      │ - Clients map            │       │
│  │ - GetTotal...()  │      │ - Broadcast channels     │       │
│  └────────┬─────────┘      └────────┬─────────────────┘       │
└───────────┼──────────────────────────┼─────────────────────────┘
            │                          │ Broadcast updates
            ▼                          ▼
┌────────────────────────────────────────────────────────────────┐
│                    Repository Layer                             │
│  ┌────────────────────────────────────────────────────────┐   │
│  │ AanmeldingRepository, RouteFundRepository              │   │
│  └────────────┬───────────────────────────────────────────┘   │
└───────────────┼─────────────────────────────────────────────────┘
                │ SQL queries
                ▼
┌────────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐│
│  │ aanmeldingen │  │ route_funds  │  │ leaderboard_view     ││
│  │ - steps      │  │ - route      │  │ (computed view)      ││
│  └──────────────┘  └──────────────┘  └──────────────────────┘│
└────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Implementatie Stappenplan

### Fase 1: Foundation (1-2 dagen)

- [ ] **Stap 1**: Creëer StepsHub structuur (`services/steps_hub.go`)
- [ ] **Stap 2**: Implementeer Client management (ReadPump/WritePump)
- [ ] **Stap 3**: Unit tests voor hub en client

### Fase 2: Integration (2-3 dagen)

- [ ] **Stap 4**: Handler implementatie (`handlers/steps_websocket_handler.go`)
- [ ] **Stap 5**: Service integratie (StepsHub in StepsService)
- [ ] **Stap 6**: Dependency injection en configuratie

### Fase 3: Features (2-3 dagen)

- [ ] **Stap 7**: Advanced features (badges, leaderboard, stats)
- [ ] **Stap 8**: Performance optimizations
- [ ] **Stap 9**: Security (JWT, permissions, rate limiting)

### Fase 4: Frontend (3-4 dagen)

- [ ] **Stap 10**: JavaScript/TypeScript client library
- [ ] **Stap 11**: React/Vue integration hooks
- [ ] **Stap 12**: Mobile app integration

### Fase 5: Testing & Deployment (2-3 dagen)

- [ ] **Stap 13**: Integration en load tests
- [ ] **Stap 14**: Monitoring & logging setup
- [ ] **Stap 15**: Documentation
- [ ] **Stap 16**: Production deployment

**Totale geschatte tijd**: 10-15 dagen

---

## 🔧 Configuration

```bash
# WebSocket configuratie
WEBSOCKET_ENABLED=true
WEBSOCKET_MAX_CONNECTIONS=10000
WEBSOCKET_READ_BUFFER_SIZE=1024
WEBSOCKET_WRITE_BUFFER_SIZE=1024
WEBSOCKET_PING_INTERVAL=30s
WEBSOCKET_PONG_TIMEOUT=60s

# Rate limiting
WEBSOCKET_RATE_LIMIT_ENABLED=true
WEBSOCKET_RATE_LIMIT_MESSAGES_PER_MINUTE=60

# Broadcasting
WEBSOCKET_BROADCAST_BUFFER_SIZE=256
WEBSOCKET_BATCH_MESSAGES=true
WEBSOCKET_BATCH_INTERVAL=100ms
```

---

## 📈 Performance Targets

| Metric | Target | Notes |
|--------|--------|-------|
| Max Concurrent Connections | 10,000 | Per server instance |
| Message Latency | < 50ms | 95th percentile |
| Message Throughput | 10,000 msg/s | Per server |
| Memory per Connection | < 10KB | Base overhead |
| CPU Usage | < 70% | At peak load |

---

## 🔐 Security Checklist

- [ ] JWT validation voor WebSocket handshake
- [ ] Permission checking per message type
- [ ] Rate limiting per client (60 msg/min)
- [ ] DDoS protection (connection limits)
- [ ] Input validation voor alle messages
- [ ] Secure WebSocket (wss:// in productie)
- [ ] CORS configuratie
- [ ] Origin checking

---

## 📊 Monitoring & Metrics

**Prometheus Metrics**:
- `websocket_connections_total` - Active connections
- `websocket_messages_sent_total` - Messages sent by type
- `websocket_messages_received_total` - Messages received
- `websocket_connection_duration_seconds` - Connection duration histogram
- `websocket_errors_total` - Errors by type

---

## 🧪 Testing Strategy

### Unit Tests
```bash
go test ./services/steps_hub_test.go -v
go test ./services/steps_service_test.go -v
go test ./handlers/steps_websocket_handler_test.go -v
```

### Integration Tests
```bash
go test ./tests/websocket_integration_test.go -v
```

### Load Tests
```bash
# 1000 concurrent connections, 10000 messages
go test ./tests/websocket_load_test.go -bench=. -benchtime=10s
```

---

## 📚 Samenvatting

### ✅ Huidige Status
- Complete HTTP REST API voor stappen tracking
- Database schema met stappen, route funds, gamification
- RBAC security met permissions
- Service layer met business logic
- Repository pattern met PostgreSQL

### 🔜 WebSocket Implementatie
- **Dedicated StepsHub** voor real-time updates
- **4 message types**: step_update, total_update, leaderboard_update, badge_earned  
- **Subscription system** voor selective updates
- **Auto-reconnect** in frontend clients
- **Broadcasting** naar alle geabonneerde clients
- **Performance optimized** voor 10,000+ connections

### 🎯 Voordelen
1. **Real-time updates** - No polling needed
2. **Efficient** - Bidirectional communication
3. **Scalable** - Handles thousands of connections
4. **Selective** - Subscribe only to needed channels
5. **Reliable** - Auto-reconnect & error handling

---

**Document Versie**: 1.0  
**Laatste Update**: 2025-01-02  
**Status**: 📝 Ready for Implementation
        if