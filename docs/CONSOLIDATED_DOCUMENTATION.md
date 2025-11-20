# DKL Email Service - Complete Backend Documentation

**Version:** 1.2.0
**Date:** 2025-11-19
**Status:** ✅ PRODUCTION READY

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture](#architecture)
3. [Database Schema](#database-schema)
4. [API Reference](#api-reference)
5. [Authentication & Authorization](#authentication--authorization)
6. [Real-time Features (WebSocket)](#real-time-features-websocket)
7. [Services & Business Logic](#services--business-logic)
8. [Data Models](#data-models)
9. [Configuration](#configuration)
10. [Deployment & Infrastructure](#deployment--infrastructure)
11. [Development Guide](#development-guide)
12. [Testing](#testing)
13. [Troubleshooting](#troubleshooting)
14. [Frontend Documentation](#frontend-documentation)
    - [Website](#website)
    - [Stepsapp](#stepsapp)
    - [Adminpanel](#adminpanel)

---

## System Overview

The DKL Email Service is a comprehensive participant management and communication platform for the De Koninklijke Loop event. Built with Go and modern web technologies, it provides a complete backend solution for event management, participant registration, real-time communication, and administrative operations.

### Key Features

- **Dual Registration System**: Full accounts (with app access) and temporary accounts
- **Role-Based Access Control (RBAC)**: Granular permissions system with automatic role assignment
- **Real-time Communication**: WebSocket-based steps tracking and leaderboard updates
- **Email Management**: Automated email sending, IMAP fetching, and template system
- **Content Management System (CMS)**: Media management, newsletters, event management
- **Gamification**: Achievements, badges, and leaderboard system
- **Multi-channel Notifications**: Email, Telegram bot integration
- **Rate Limiting & Security**: Comprehensive security measures and API protection

### Technology Stack

- **Backend**: Go 1.23.0 with Fiber web framework
- **Database**: PostgreSQL 17 with Redis caching
- **Real-time**: WebSocket connections with custom hubs
- **Email**: SMTP sending, IMAP fetching, template system
- **Media**: Cloudinary integration for image/video management
- **Deployment**: Docker containers on Render platform
- **Monitoring**: Prometheus metrics, structured logging

---

## Architecture

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Handlers      │    │   Services      │    │  Repository     │
│   (HTTP/WS)     │◄──►│  (Business      │◄──►│  (Data Access)  │
│                 │    │   Logic)        │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Fiber Router  │    │   Redis Cache   │    │  PostgreSQL DB  │
│   Middleware    │    │   Sessions      │    │  Migrations     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Core Components

#### 1. Handlers Layer
- **HTTP Handlers**: REST API endpoints for all system operations
- **WebSocket Handlers**: Real-time communication for steps tracking and notulen
- **Middleware**: Authentication, authorization, rate limiting, CORS
- **Validation**: Request validation and error handling

#### 2. Services Layer
- **AuthService**: JWT token management, user authentication
- **EmailService**: SMTP sending, template rendering
- **StepsService**: Participant steps tracking and leaderboard
- **GamificationService**: Achievements and badges management
- **PermissionService**: RBAC permission checking with caching
- **MailFetcher**: IMAP email fetching and processing

#### 3. Repository Layer
- **Data Access Objects**: PostgreSQL implementations for all entities
- **Query Optimization**: Efficient database queries with proper indexing
- **Transaction Management**: Database transaction handling
- **Connection Pooling**: Optimized database connections

#### 4. Configuration Layer
- **Environment Variables**: Flexible configuration management
- **Database Config**: PostgreSQL connection settings
- **Redis Config**: Caching and session management
- **Cloudinary Config**: Media storage configuration

### Data Flow

1. **Request Processing**:
   - HTTP/WebSocket request → Middleware validation → Handler
   - Handler → Service layer → Repository layer → Database
   - Response ← Service layer ← Repository layer ← Database

2. **Real-time Updates**:
   - Steps updates → StepsHub → WebSocket broadcast
   - Leaderboard changes → Cached calculations → Real-time push

3. **Email Processing**:
   - IMAP fetcher → Email decoder → Database storage
   - Template rendering → SMTP sending → Metrics tracking

---

## Database Schema

### Core Tables Overview

#### User Management & Authentication
- **`gebruikers`**: Admin/staff accounts with RBAC roles
- **`participants`**: Participant data (person information)
- **`event_registrations`**: Event-specific data (steps, role, distance)
- **`refresh_tokens`**: JWT refresh token storage
- **`access_tokens`**: API access token management with session linking
- **`sessions`**: Multi-device session tracking with device fingerprinting
- **`password_reset_tokens`**: Password reset token management
- **`email_verification_tokens`**: Email verification token storage

#### RBAC System
- **`roles`**: Available system roles
- **`permissions`**: Granular permission definitions
- **`user_roles`**: User-role assignments
- **`role_permissions`**: Role-permission mappings

#### Content Management
- **`events`**: Event definitions with geofences and configuration
- **`videos`**, **`photos`**, **`albums`**: Media management
- **`partners`**, **`sponsors`**: Partner and sponsor management
- **`newsletters`**: Newsletter content and scheduling
- **`contact_formulieren`**: Contact form submissions
- **`notulen`**: Meeting notes and minutes

#### Gamification
- **`achievements`**, **`badges`**: Achievement definitions
- **`user_achievements`**, **`user_badges`**: User progress tracking
- **`leaderboards`**: Dynamic leaderboard calculations

#### Communication
- **`incoming_email`**: IMAP-fetched emails
- **`verzonden_email`**: Sent email tracking
- **`notifications`**: System notifications
- **`auto_responses`**: Automated email responses

### Key Schema Evolution

#### V28: Participant Refactor (Foundation)
**Impact**: Major schema restructuring
- Moved event-specific data from `participants` to `event_registrations`
- Separated person data from participation data
- Enabled multi-event participation per person

#### V30: Dual Registration System
**Impact**: New account types and RBAC integration
- Added `account_type` ('full' or 'temporary') to participants
- Added `has_app_access` boolean flag
- Integrated RBAC with automatic role assignment
- Added participant upgrade mechanism

#### V34: Legacy Column Removal
**Impact**: Permanent data removal (non-reversible)
- Removed 14 legacy columns from `participants` table
- Cleaned up database schema after V28-V31 migrations
- Requires full database backup before deployment

#### V35-V41: Permission and User Management
**Impact**: Enhanced RBAC and user management
- Added staff user permissions and roles
- Created password reset token system
- Implemented email verification tokens
- Added menu permissions for UI access control

#### V42-V44: Authentication Enhancements
**Impact**: Advanced authentication features
- Created access tokens table for server-side token management
- Added email verification token system
- Implemented password reset token functionality
- Enhanced security with token rotation capabilities

#### V45-V47: Session Management
**Impact**: Multi-device session tracking
- Created sessions table for device fingerprinting
- Added session_id to access tokens for session linking
- Implemented session expiration and activity tracking
- Enhanced security with IP address and user agent logging

### Database Relationships

```
participants (1) ──── (N) event_registrations (N) ──── (1) events
     │                        │
     │                        │
     └─── (1) gebruikers       └─── (N) achievements
          (optional)                │
                                   │
                                   └─── (N) badges
```

### Indexes and Performance

#### Critical Indexes
```sql
-- Participant lookups
CREATE INDEX idx_participants_email ON participants(email);
CREATE INDEX idx_participants_account_type ON participants(account_type);

-- Event registration queries
CREATE INDEX idx_event_registrations_participant_id ON event_registrations(participant_id);
CREATE INDEX idx_event_registrations_event_id ON event_registrations(event_id);

-- RBAC performance
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);

-- Steps and leaderboard
CREATE INDEX idx_event_registrations_steps ON event_registrations(steps);
CREATE INDEX idx_event_registrations_distance_route ON event_registrations(distance_route);
```

---

## API Reference

### Authentication Endpoints

#### POST `/api/auth/login`
Authenticate user and return JWT tokens with session creation
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

#### POST `/api/auth/logout`
Revoke current session and access token

#### POST `/api/auth/refresh`
Refresh access token using refresh token
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### GET `/api/auth/profile`
Get current user profile (authenticated)

#### POST `/api/auth/reset-password`
Reset password (authenticated user)

#### POST `/api/auth/forgot-password`
Request password reset email

#### POST `/api/auth/reset-password-with-token`
Reset password using reset token

#### POST `/api/auth/send-verification`
Send email verification

#### POST `/api/auth/verify-email`
Verify email address

#### POST `/api/auth/resend-verification`
Resend email verification

#### GET `/api/auth/sessions`
List all active sessions for current user

#### DELETE `/api/auth/sessions/{sessionId}`
Revoke specific session

#### POST `/api/auth/sessions/revoke-others`
Revoke all other sessions

### Participant Management

#### GET `/api/participants`
List participants with pagination and filtering
**Query Parameters:**
- `limit`: Number of results (default: 10)
- `offset`: Pagination offset (default: 0)
- `event_id`: Filter by event
- `search`: Search by name or email

#### GET `/api/participants/{id}`
Get participant details

#### POST `/api/participants`
Create new participant (admin only)

#### PUT `/api/participants/{id}`
Update participant information (admin only)

#### DELETE `/api/participants/{id}`
Delete participant (admin only)

#### GET `/api/participants/{id}/registrations`
Get participant event registrations

#### POST `/api/public/aanmelden`
Public participant registration
```json
{
  "naam": "John Doe",
  "email": "john@example.com",
  "telefoon": "0612345678",
  "afstand": "10 KM",
  "rol": "Deelnemer",
  "ondersteuning": "Nee",
  "want_account": true,
  "wachtwoord": "securepassword"
}
```

### Event Management

#### GET `/api/events/active`
Get currently active event

#### POST `/api/events`
Create new event (admin only)
```json
{
  "name": "DKL 2025",
  "description": "De Koninklijke Loop 2025",
  "start_time": "2025-06-01T10:00:00Z",
  "end_time": "2025-06-01T16:00:00Z",
  "geofences": [
    {
      "name": "Start",
      "type": "start",
      "lat": 52.3676,
      "long": 4.9041,
      "radius": 100
    }
  ]
}
```

### Steps Tracking

#### POST `/api/steps/update`
Update participant steps
```json
{
  "delta_steps": 500
}
```

#### GET `/api/steps/dashboard`
Get participant steps dashboard

#### GET `/api/leaderboard`
Get current leaderboard
**Query Parameters:**
- `limit`: Number of results (default: 50)
- `route`: Filter by distance route

### WebSocket Endpoints

#### `/ws/steps`
Real-time steps tracking and leaderboard updates
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/steps?token=jwt_token');

// Subscribe to updates
ws.send(JSON.stringify({
  type: 'subscribe',
  channels: ['step_updates', 'leaderboard_updates']
}));
```

#### `/api/ws/notulen`
Real-time meeting notes collaboration

### CMS Endpoints

#### Albums Management
- `GET /api/albums` - List visible albums
- `POST /api/albums` - Create album (admin)
- `PUT /api/albums/{id}` - Update album (admin)
- `DELETE /api/albums/{id}` - Delete album (admin)

#### Photos Management
- `GET /api/photos` - List visible photos
- `POST /api/photos` - Upload photo (admin)
- `PUT /api/photos/{id}` - Update photo (admin)

#### Content Management
- `GET /api/videos` - List videos
- `POST /api/videos` - Upload video (admin)
- `PUT /api/videos/{id}` - Update video (admin)
- `DELETE /api/videos/{id}` - Delete video (admin)

- `GET /api/partners` - List partners
- `POST /api/partners` - Create partner (admin)
- `PUT /api/partners/{id}` - Update partner (admin)
- `DELETE /api/partners/{id}` - Delete partner (admin)

- `GET /api/sponsors` - List sponsors
- `POST /api/sponsors` - Create sponsor (admin)
- `PUT /api/sponsors/{id}` - Update sponsor (admin)
- `DELETE /api/sponsors/{id}` - Delete sponsor (admin)

- `GET /api/newsletters` - List newsletters
- `POST /api/newsletters` - Create newsletter (admin)
- `PUT /api/newsletters/{id}` - Update newsletter (admin)
- `DELETE /api/newsletters/{id}` - Delete newsletter (admin)

#### Chat System
- `GET /api/chat/channels` - List chat channels
- `POST /api/chat/channels` - Create channel (admin)
- `GET /api/chat/channels/{id}/messages` - Get channel messages
- `POST /api/chat/messages` - Send message
- `PUT /api/chat/messages/{id}` - Edit message
- `DELETE /api/chat/messages/{id}` - Delete message
- `POST /api/chat/messages/{id}/reactions` - Add reaction

#### Notulen (Meeting Notes)
- `GET /api/notulen` - List meeting notes
- `POST /api/notulen` - Create meeting notes
- `GET /api/notulen/{id}` - Get specific notes
- `PUT /api/notulen/{id}` - Update notes
- `DELETE /api/notulen/{id}` - Delete notes
- `WebSocket /ws/notulen` - Real-time collaboration

#### Auto Responses
- `GET /api/mail/autoresponse` - List auto responses
- `POST /api/mail/autoresponse` - Create auto response
- `PUT /api/mail/autoresponse/{id}` - Update auto response
- `DELETE /api/mail/autoresponse/{id}` - Delete auto response

#### Image Management
- `POST /api/images/upload` - Upload image
- `GET /api/images` - List images
- `DELETE /api/images/{id}` - Delete image

### Email Management

#### POST `/api/contact-email`
Send contact form email
```json
{
  "naam": "John Doe",
  "email": "john@example.com",
  "bericht": "Contact message",
  "privacy_akkoord": true
}
```

#### GET `/api/mail/unprocessed`
Get unprocessed emails (admin)

#### POST `/api/mail/send`
Send custom email (admin)

### RBAC Management

#### GET `/api/roles`
List all roles (admin)

#### GET `/api/permissions`
List all permissions (admin)

#### POST `/api/rbac/roles`
Create new role (admin)

#### PUT `/api/rbac/roles/{id}`
Update role (admin)

#### DELETE `/api/rbac/roles/{id}`
Delete role (admin)

#### POST `/api/rbac/roles/{id}/permissions`
Assign permissions to role (admin)

#### DELETE `/api/rbac/roles/{id}/permissions/{permissionId}`
Remove permission from role (admin)

#### GET `/api/users`
List users with roles (admin)

#### POST `/api/users/{id}/roles`
Assign roles to user (admin)

#### DELETE `/api/users/{id}/roles/{roleId}`
Remove role from user (admin)

#### GET `/api/rbac/roles/{id}/users`
List users with specific role

### Health & Monitoring

#### GET `/api/health`
System health check
```json
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "version": "1.1.0",
  "uptime": "2h45m"
}
```

#### GET `/metrics/email`
Email service metrics

#### GET `/metrics/rate-limits`
Rate limiting statistics

---

## Authentication & Authorization

### JWT Token Structure

```json
{
  "email": "user@example.com",
  "roles": ["participant_user", "participant_guide"],
  "rbac_active": true,
  "exp": 1638360000,
  "iat": 1638273600
}
```

### RBAC Permission System

#### Core Roles
- **`participant_user`**: Base role for full account participants
- **`participant_guide`**: Guides with additional permissions
- **`participant_volunteer`**: Volunteers with event support permissions
- **`admin`**: Full system access
- **`content_manager`**: CMS access

#### Permission Categories
- **`app:access`**: Application access permissions
- **`steps:track`**: Steps tracking and viewing
- **`leaderboard:view`**: Leaderboard access
- **`achievements:*`**: Achievement management
- **`admin:*`**: Administrative operations

### Authentication Flow

1. **Login Request** → Validate credentials
2. **Participant Check** → Verify account type and app access
3. **RBAC Loading** → Load user roles and permissions
4. **Token Generation** → Create JWT with roles array
5. **Permission Caching** → Cache permissions in Redis (5min TTL)

### Session Management

- **Access Tokens**: Short-lived (15 minutes) with server-side storage
- **Refresh Tokens**: Long-lived (30 days)
- **Session Tracking**: Multi-device session management with device fingerprinting
- **Device Information**: Comprehensive login tracking (IP, user agent, location)
- **Session Revocation**: Individual and bulk session revocation
- **Security Headers**: Enhanced security with proper HTTP headers

---

## Real-time Features (WebSocket)

### WebSocket Hubs

#### StepsHub
- Manages real-time steps updates with subscription filtering
- Broadcasts leaderboard changes to subscribed clients
- Handles participant connections with connection pooling
- Maintains active connection pools for efficient broadcasting
- Supports channel-based subscriptions (step_updates, leaderboard_updates)

#### NotulenHub
- Real-time meeting notes collaboration with multi-user editing
- Version conflict resolution and concurrent editing support
- Connection management with session tracking
- Real-time synchronization across multiple clients

### Message Types

#### Subscription Management
```json
{
  "type": "subscribe",
  "channels": ["step_updates", "leaderboard_updates"],
  "filters": {
    "event_id": "uuid",
    "distance_route": "10 KM"
  }
}
```

#### Steps Tracking
```json
{
  "type": "step_update",
  "participant_id": "uuid",
  "naam": "John Doe",
  "steps": 12500,
  "delta": 500,
  "timestamp": "2025-11-19T10:30:00Z"
}
```

#### Leaderboard Updates
```json
{
  "type": "leaderboard_update",
  "entries": [
    {
      "id": "uuid",
      "naam": "Jane Runner",
      "route": "10 KM",
      "steps": 15000,
      "rank": 1
    }
  ]
}
```

#### Achievement Notifications
```json
{
  "type": "badge_earned",
  "participant_id": "uuid",
  "badge_name": "First Steps",
  "badge_description": "Completed first 1000 steps"
}
```

### Connection Management

- **Authentication**: JWT token validation on connection
- **Keep-alive**: Ping/pong every 30 seconds
- **Reconnection**: Automatic reconnection handling
- **Broadcasting**: Efficient message broadcasting to multiple clients

---

## Services & Business Logic

### Core Services

#### AuthService
- JWT token generation and validation with access token rotation
- Password hashing and verification with bcrypt
- User registration and account management
- Session management with multi-device support
- Email verification token handling
- Password reset token management
- Device fingerprinting and security tracking

#### EmailService
- SMTP email sending with templates
- Email template rendering
- Bulk email operations
- Email metrics and tracking

#### StepsService
- Steps tracking and validation
- Leaderboard calculations
- Route-based filtering
- Real-time broadcast integration

#### GamificationService
- Achievement and badge management
- Progress tracking
- Automatic badge awarding
- Leaderboard integration

#### PermissionService
- RBAC permission checking with Redis caching
- Role hierarchy management and inheritance
- Participant app access validation
- Dynamic permission evaluation
- Menu permission checking for UI access control

#### ChatService
- Real-time chat channel management
- Message threading and reactions
- User presence tracking
- Channel participant management
- Message history and search

#### ImageService
- Cloudinary integration for media management
- Image upload, transformation, and optimization
- Secure URL generation with access controls
- Batch image processing
- CDN delivery optimization

### Background Services

#### MailFetcher
- IMAP email fetching
- Email parsing and decoding
- Automatic processing
- Error handling and retry logic

#### EmailAutoFetcher
- Scheduled email fetching
- Multiple account support
- Processing queue management
- Metrics collection

#### NewsletterService
- RSS feed processing
- Content filtering and formatting
- Scheduled sending
- Subscriber management

### Service Factory Pattern

The ServiceFactory implements a comprehensive dependency injection pattern that manages all service lifecycles and their interdependencies.

```go
type ServiceFactory struct {
    // Core Services
    AuthService         AuthService
    EmailService        *EmailService
    PermissionService   PermissionService
    ChatService         ChatService
    ImageService        *ImageService

    // Background Services
    EmailAutoFetcher    EmailAutoFetcherInterface
    NewsletterService   *NewsletterService
    TelegramBotService  *TelegramBotService
    NotulenService      *NotulenService

    // Infrastructure
    RateLimiter         RateLimiterInterface
    EmailMetrics        *EmailMetrics
    EmailBatcher        *EmailBatcher
    RedisClient         *redis.Client
    Hub                 *Hub

    // Business Services
    GamificationService *GamificationService
    NewsletterSender    *NewsletterSender
}
```

---

## Data Models

### Core Models

#### Participant Model
```go
type Participant struct {
    ID                uuid.UUID    `json:"id"`
    Naam              string       `json:"naam"`
    Email             string       `json:"email"`
    Telefoon          *string      `json:"telefoon"`
    AccountType       AccountType  `json:"account_type"`
    HasAppAccess      bool         `json:"has_app_access"`
    RegistrationYear  *int         `json:"registration_year"`
    GebruikerID       *uuid.UUID   `json:"gebruiker_id"`
    Terms             bool         `json:"terms"`
    TestMode          bool         `json:"test_mode"`
    CreatedAt         time.Time    `json:"created_at"`
    UpdatedAt         time.Time    `json:"updated_at"`
}
```

#### EventRegistration Model
```go
type EventRegistration struct {
    ID                uuid.UUID      `json:"id"`
    EventID           uuid.UUID      `json:"event_id"`
    ParticipantID     uuid.UUID      `json:"participant_id"`
    ParticipantRoleName string       `json:"participant_role_name"`
    DistanceRoute     *string        `json:"distance_route"`
    Steps             int            `json:"steps"`
    Ondersteuning     *string        `json:"ondersteuning"`
    Status            *string        `json:"status"`
    RegisteredAt      time.Time      `json:"registered_at"`
    // ... additional fields
}
```

#### Authentication Models
```go
type AccessToken struct {
    ID        string    `json:"id"`
    OwnerID   string    `json:"owner_id"`
    SessionID *string   `json:"session_id"`
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    RevokedAt *time.Time `json:"revoked_at"`
    IsRevoked bool      `json:"is_revoked"`
}

type Session struct {
    ID           string       `json:"id"`
    OwnerID      string       `json:"owner_id"`
    AccessToken  string       `json:"access_token"`
    DeviceInfo   DeviceInfo   `json:"device_info"`
    IPAddress    string       `json:"ip_address"`
    UserAgent    string       `json:"user_agent"`
    IsActive     bool         `json:"is_active"`
    IsCurrent    bool         `json:"is_current"`
    ExpiresAt    time.Time    `json:"expires_at"`
    LastActivity time.Time    `json:"last_activity"`
}

type DeviceInfo struct {
    Browser        string `json:"browser"`
    BrowserVersion string `json:"browser_version"`
    OS             string `json:"os"`
    OSVersion      string `json:"os_version"`
    DeviceType     string `json:"device_type"`
    Platform       string `json:"platform"`
}
```

#### RBAC Models
```go
type RBACRole struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    IsActive    bool      `json:"is_active"`
}

type Permission struct {
    ID          uuid.UUID `json:"id"`
    Resource    string    `json:"resource"`
    Action      string    `json:"action"`
    Description string    `json:"description"`
}
```

### Relationship Models

#### AlbumPhoto Model
```go
type AlbumPhoto struct {
    AlbumID   uuid.UUID `json:"album_id"`
    PhotoID   uuid.UUID `json:"photo_id"`
    OrderNumber int      `json:"order_number"`
}
```

#### UserAchievement Model
```go
type UserAchievement struct {
    UserID       uuid.UUID `json:"user_id"`
    AchievementID uuid.UUID `json:"achievement_id"`
    EarnedAt     time.Time `json:"earned_at"`
    Progress     int       `json:"progress"`
}
```

---

## Configuration

### Environment Variables

#### Required Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure_password
DB_NAME=dklemailservice
DB_SSL_MODE=require

# JWT
JWT_SECRET=your-256-bit-secret-key

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=app-password
SMTP_FROM=noreply@yourdomain.com

# Admin
ADMIN_EMAIL=admin@yourdomain.com
```

#### Optional Variables
```bash
# Redis (optional)
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379

# Cloudinary (optional)
CLOUDINARY_CLOUD_NAME=your-cloud
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret

# Email Fetching (optional)
INFO_EMAIL=info@yourdomain.com
INFO_EMAIL_PASSWORD=password
IMAP_SERVER=imap.gmail.com

# Rate Limiting
EMAIL_RATE_LIMIT=10
CONTACT_RATE_LIMIT=5
REGISTRATION_RATE_LIMIT=3

# Logging
LOG_LEVEL=INFO
```

### Configuration Loading

#### Database Configuration
```go
type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func LoadDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvAsInt("DB_PORT", 5432),
        User:     getEnv("DB_USER", "postgres"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   getEnv("DB_NAME", "dklemailservice"),
        SSLMode:  getEnv("DB_SSL_MODE", "disable"),
    }
}
```

#### Logger Configuration
```go
type LoggerConfig struct {
    Level      string
    Format     string
    Output     string
    ELKEnabled bool
    ELKURL     string
}

func LoadLoggerConfig() *LoggerConfig {
    return &LoggerConfig{
        Level:      getEnv("LOG_LEVEL", "INFO"),
        Format:     getEnv("LOG_FORMAT", "json"),
        Output:     getEnv("LOG_OUTPUT", "stdout"),
        ELKEnabled: getEnvAsBool("ELK_ENABLED", false),
        ELKURL:     os.Getenv("ELK_URL"),
    }
}
```

---

## Deployment & Infrastructure

### Docker Deployment

#### Dockerfile (Multi-stage)
```dockerfile
# Build stage for production
FROM golang:1.24-alpine AS builder-prod
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main-prod -ldflags="-w -s" .

# Build stage for development
FROM golang:1.24-alpine AS builder-dev
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main-dev .

# Runtime stage
FROM alpine:3.19 AS runtime
RUN apk --no-cache add ca-certificates sqlite openssh-server
COPY --from=builder-prod /app/main-prod ./main
COPY --from=builder-dev /app/main-dev ./main-dev
COPY --from=builder-prod /app/templates ./templates
EXPOSE 8080
CMD ["./main"]
```

#### Docker Compose (Development)
```yaml
services:
  postgres:
    image: postgres:17
    environment:
      POSTGRES_DB: dklemailservice
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5433:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6380:6379"
    command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      REDIS_HOST: redis
      # ... other env vars
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
```

### Render Deployment

#### render.yaml
```yaml
services:
  - type: web
    name: dklautomatie-backend
    env: go
    buildCommand: go build -o server
    startCommand: ./server
    envVars:
      - key: ADMIN_EMAIL
        sync: false
      - key: SMTP_HOST
        sync: false
      # ... other secrets
    healthCheckPath: /api/health

  - type: pserv
    name: dklautomatie-db
    engine: postgres
    version: 14
    ipAllowList: []
    plan: free
    autoDeploy: true
    envVars:
      - key: POSTGRES_USER
        generateValue: true
      - key: POSTGRES_PASSWORD
        generateValue: true
      - key: POSTGRES_DB
        value: dklemailservice

  - type: redis
    name: dklautomatie-redis
    ipAllowList: []
    plan: free
    autoDeploy: true
```

### Infrastructure Components

#### PostgreSQL Database
- **Version**: 17 (production), 14 (Render)
- **Extensions**: UUID, PostGIS (optional)
- **Backup**: Automated daily backups with point-in-time recovery
- **Connection Pooling**: Built-in connection pooling with optimized settings
- **Performance**: Advanced indexing, query optimization, and partitioning support

#### Redis Cache
- **Version**: 7.x with advanced features
- **Usage**: Session storage, permission caching, rate limiting, real-time data
- **Persistence**: RDB snapshots with configurable intervals
- **Memory**: LRU eviction, memory optimization for high-throughput scenarios
- **Clustering**: Support for Redis Cluster for horizontal scaling

#### Load Balancing
- **Render**: Automatic load balancing
- **Custom**: Nginx reverse proxy recommended
- **SSL**: Automatic SSL certificate management

### Scaling Strategy

#### Vertical Scaling
- Increase Render service plan
- Upgrade database instance size
- Add more Redis memory

#### Horizontal Scaling
- Multiple application instances
- Read replicas for database
- Redis cluster for caching
- Load balancer configuration

---

## Development Guide

### Project Structure

```
dklemailservice/
├── main.go                 # Application entry point
├── handlers/               # HTTP request handlers
├── services/               # Business logic
├── models/                 # Data models
├── repository/             # Data access layer
├── config/                 # Configuration management
├── database/               # Migrations and scripts
├── templates/              # Email templates
├── tests/                  # Test files
├── docs/                   # Documentation
└── scripts/                # Utility scripts
```

### Development Setup

#### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- PostgreSQL client (optional)
- Git

#### Local Development
```bash
# Clone repository
git clone <repository-url>
cd dklemailservice

# Copy environment file
cp .env.example .env

# Start dependencies
docker-compose up -d postgres redis

# Install dependencies
go mod download

# Run application
go run main.go

# Run tests
go test ./...
```

### Code Standards

#### Go Code Guidelines
- Use `gofmt` for formatting
- Follow standard Go naming conventions
- Add comprehensive error handling
- Include unit tests for all functions
- Use meaningful variable and function names
- Add comments for complex logic

#### Database Guidelines
- Use migrations for schema changes
- Include rollback scripts
- Add appropriate indexes
- Document complex queries
- Use transactions for data consistency

#### API Design Guidelines
- RESTful endpoint naming
- Consistent JSON response format
- Proper HTTP status codes
- Comprehensive error messages
- Version API endpoints when breaking changes

### Database Migrations

#### Creating New Migrations
```bash
# Create migration file
touch database/migrations/V47__new_feature.sql

# Migration structure
-- +migrate Up
-- SQL for applying migration

-- +migrate Down
-- SQL for rolling back migration
```

#### Migration Best Practices
- Test migrations on development database first
- Include rollback scripts for all migrations
- Document any breaking changes
- Backup database before applying destructive migrations
- Test rollback procedures

---

## Testing

### Test Structure

#### Unit Tests
```go
func TestAuthService_Login(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    service := NewAuthService(mockRepo)

    // Test
    token, err := service.Login("user@example.com", "password")

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

#### Integration Tests
```go
func TestParticipantRegistration_Integration(t *testing.T) {
    // Setup database
    db := setupTestDB(t)
    defer db.Close()

    // Create repositories and services
    repo := repository.NewParticipantRepository(db)
    service := services.NewParticipantService(repo)

    // Test full registration flow
    participant, err := service.RegisterParticipant(registrationData)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, participant)
}
```

#### API Tests
```go
func TestCreateParticipant_API(t *testing.T) {
    app := fiber.New()
    handler := setupParticipantHandler(t)

    // Create test request
    req := httptest.NewRequest("POST", "/api/participants", strings.NewReader(participantJSON))
    req.Header.Set("Content-Type", "application/json")

    // Execute request
    resp, err := app.Test(req)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
}
```

### Test Categories

#### Handler Tests (35+ test files)
- HTTP request/response validation with comprehensive coverage
- Authentication middleware testing with JWT validation
- Authorization middleware testing with RBAC permissions
- Error response formatting and status code validation
- Rate limiting verification and threshold testing
- WebSocket handler testing with connection management
- CORS middleware testing with origin validation

#### Service Tests (25+ test files)
- Business logic validation with edge case coverage
- External service integration (SMTP, Redis, Cloudinary)
- Data transformation and serialization testing
- Error handling scenarios with proper error propagation
- Service factory pattern testing with dependency injection
- Background service testing (email fetcher, newsletter sender)

#### Repository Tests (20+ test files)
- Database query correctness with complex JOIN operations
- Transaction management and rollback testing
- Connection pooling and database connection handling
- SQL injection prevention with parameterized queries
- Migration testing with up/down script validation
- Repository factory pattern testing

#### Integration Tests (15+ test files)
- End-to-end API flows with real database operations
- Database migration testing with schema validation
- External service mocking (email, Redis, Cloudinary)
- Performance validation with load testing
- WebSocket integration testing with real-time features
- Authentication flow integration testing

#### Security Tests (10+ test files)
- Authentication security testing (token validation, session management)
- Authorization testing with permission matrix validation
- Rate limiting security testing with abuse prevention
- Input validation and sanitization testing
- SQL injection and XSS prevention testing

### Test Execution

#### Running Tests
```bash
# Run all tests
go test ./...

# Run specific package
go test ./handlers -v

# Run with coverage
go test -cover ./...

# Run integration tests only
go test -tags=integration ./...
```

#### Test Configuration
```bash
# Test database setup
export DB_HOST=localhost
export DB_PORT=5433
export DB_NAME=dklemailservice_test

# Enable test mode
export APP_ENV=test
export LOG_LEVEL=debug
```

### Test Coverage Goals

- **Handlers**: >85% coverage (35+ handler test files)
- **Services**: >90% coverage (25+ service test files)
- **Models**: >75% coverage (comprehensive model validation)
- **Repository**: >95% coverage (20+ repository test files)
- **Integration**: >80% coverage (15+ integration test files)
- **Security**: >90% coverage (10+ security test files)
- **Overall**: >85% coverage (80+ test files total)

---

## Troubleshooting

### Common Issues

#### Database Connection Issues
**Symptoms:**
- Application fails to start
- "connection refused" errors
- Migration failures

**Solutions:**
```bash
# Check database status
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection
psql -h localhost -p 5433 -U postgres -d dklemailservice

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

#### Authentication Problems
**Symptoms:**
- Login failures
- JWT token errors
- Permission denied errors

**Solutions:**
```bash
# Check JWT secret
echo $JWT_SECRET

# Verify user exists
psql -d dklemailservice -c "SELECT * FROM gebruikers WHERE email = 'user@example.com';"

# Check user roles
psql -d dklemailservice -c "SELECT * FROM user_roles WHERE user_id = (SELECT id FROM gebruikers WHERE email = 'user@example.com');"
```

#### Email Sending Issues
**Symptoms:**
- Emails not being sent
- SMTP connection errors
- Template rendering failures

**Solutions:**
```bash
# Test SMTP connection
telnet smtp.gmail.com 587

# Check email configuration
echo "SMTP_HOST: $SMTP_HOST"
echo "SMTP_USER: $SMTP_USER"

# Verify email templates exist
ls -la templates/
```

#### WebSocket Connection Issues
**Symptoms:**
- Real-time updates not working
- WebSocket connection failures
- Steps not updating in real-time

**Solutions:**
```bash
# Check WebSocket endpoint
curl -I http://localhost:8080/ws/steps

# Verify JWT token
# Check application logs for WebSocket errors

# Test WebSocket connection
wscat -c "ws://localhost:8080/ws/steps?token=YOUR_JWT_TOKEN"
```

### Performance Issues

#### High Memory Usage
**Symptoms:**
- Application consuming excessive memory
- Out of memory errors
- Slow response times

**Solutions:**
```bash
# Check memory usage
docker stats

# Profile application
go tool pprof http://localhost:8080/debug/pprof/heap

# Optimize goroutines
# Check for memory leaks in WebSocket connections
```

#### Slow Database Queries
**Symptoms:**
- API responses taking too long
- Database CPU usage high
- Timeout errors

**Solutions:**
```bash
# Check slow queries
psql -d dklemailservice -c "SELECT * FROM pg_stat_activity WHERE state = 'active';"

# Analyze query plans
EXPLAIN ANALYZE SELECT * FROM participants WHERE email = 'user@example.com';

# Check indexes
psql -d dklemailservice -c "\di"
```

#### High CPU Usage
**Symptoms:**
- Application consuming excessive CPU
- Slow response times
- System becoming unresponsive

**Solutions:**
```bash
# Profile CPU usage
go tool pprof http://localhost:8080/debug/pprof/profile

# Check goroutine count
# Optimize concurrent operations
# Review rate limiting configuration
```

### Monitoring and Logs

#### Application Logs
```bash
# View application logs
docker-compose logs -f app

# Filter by level
docker-compose logs app | grep ERROR

# Search for specific errors
docker-compose logs app | grep "WebSocket"
```

#### Database Logs
```bash
# PostgreSQL logs
docker-compose logs postgres

# Query execution statistics
psql -d dklemailservice -c "SELECT * FROM pg_stat_user_tables ORDER BY n_tup_ins DESC;"
```

#### System Monitoring
```bash
# Health check
curl http://localhost:8080/api/health

# Metrics endpoint
curl http://localhost:8080/metrics

# Database connections
psql -d dklemailservice -c "SELECT count(*) FROM pg_stat_activity;"
```

### Emergency Procedures

#### Application Restart
```bash
# Restart application
docker-compose restart app

# Full rebuild
docker-compose down
docker-compose up --build
```

#### Database Recovery
```bash
# Create backup
pg_dump -h localhost -p 5433 -U postgres dklemailservice > backup.sql

# Restore from backup
psql -h localhost -p 5433 -U postgres -d dklemailservice < backup.sql
```

#### Rollback Deployment
```bash
# Rollback to previous version (Render)
# Use Render dashboard to redeploy previous version

# Manual rollback
git checkout previous-commit-hash
docker-compose up --build
```

---

## Frontend Documentation

The DKL Email Service backend supports three main frontend applications, each serving different user groups and use cases.

### Website
**Location:** `docs/frontend/Website/`

The public-facing website provides information about De Koninklijke Loop event, participant registration, and general event information.

- **Target Users**: General public, potential participants
- **Key Features**: Event information, registration forms, news, partner displays
- **Backend Integration**: Public API endpoints, CMS content delivery
- **Planning Documents**: See `docs/frontend/Website/` for detailed specifications

### Stepsapp
**Location:** `docs/frontend/Stepsapp/`

The mobile/web application for participants to track their walking/running steps during the event.

- **Target Users**: Registered participants with app access
- **Key Features**: Step tracking, real-time leaderboard, achievements, gamification
- **Backend Integration**: WebSocket real-time updates, participant authentication, steps API
- **Planning Documents**: See `docs/frontend/Stepsapp/` for detailed specifications

### Adminpanel
**Location:** `docs/frontend/Adminpanel/`

The administrative dashboard for event organizers and staff to manage all aspects of the event.

- **Target Users**: Event staff, administrators, content managers
- **Key Features**: Participant management, CMS content editing, analytics, system monitoring
- **Backend Integration**: Full admin API access, RBAC permissions, real-time monitoring
- **Planning Documents**: See `docs/frontend/Adminpanel/` for detailed specifications

### Frontend-Backend Integration

#### Authentication Flow
- JWT-based authentication with role-based access control
- Session management with multi-device support
- Secure API communication with proper CORS configuration

#### Real-time Features
- WebSocket connections for live updates (steps, leaderboard, chat)
- Subscription-based message filtering for performance optimization
- Connection pooling for efficient resource management

#### API Design Principles
- RESTful endpoints with consistent naming conventions
- Comprehensive error handling with detailed error messages
- Rate limiting and security headers for protection
- Pagination and filtering for large datasets

---

## Conclusion

This consolidated documentation provides a comprehensive overview of the DKL Email Service backend system. The platform has evolved from a simple email service to a full-featured participant management system with advanced features including dual registration, RBAC authorization, real-time WebSocket communication, and comprehensive gamification.

### Key Achievements

- ✅ **Advanced Authentication**: Access tokens, sessions, email verification, and security headers
- ✅ **Dual Registration System**: Full accounts with app access and temporary accounts
- ✅ **RBAC Authorization**: Complete role-based access control with granular permissions
- ✅ **Real-time Features**: WebSocket-based steps tracking with subscription filtering and connection pooling
- ✅ **Comprehensive API**: 40+ handlers with 100+ endpoints covering all system functionality
- ✅ **Enhanced Service Layer**: Factory pattern with 15+ services for business logic
- ✅ **Complete Database Schema**: Through V47 migrations with advanced features
- ✅ **Production Ready**: Docker deployment on Render with monitoring and logging
- ✅ **Comprehensive Testing**: 80+ test files with 85%+ coverage across all components
- ✅ **Complete Documentation**: Detailed technical documentation and troubleshooting guides
- ✅ **Frontend Integration**: Support for Website, Stepsapp, and Adminpanel applications

### System Status

🟢 **FULLY OPERATIONAL AND PRODUCTION READY**

The system successfully handles participant registration, real-time tracking, email communication, content management, and administrative operations for the De Koninklijke Loop event.

### Future Enhancements

- Mobile application integration
- Advanced analytics and reporting
- Multi-event support
- Enhanced gamification features
- API rate limiting improvements
- Performance optimizations

For additional details, refer to the individual documentation files in the `docs/` directory or contact the development team.