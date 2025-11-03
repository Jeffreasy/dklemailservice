# DKL Email Service - Architecture Documentation

## 📐 Project Structure

This project follows Clean Architecture principles with clear separation of concerns.

```
dklemailservice/
├── api/                    # API definitions & JWT middleware
├── config/                 # Configuration management
├── database/               # Database migrations & scripts
├── docs/                   # Project documentation
├── handlers/               # HTTP request handlers (Presentation Layer)
├── logger/                 # Structured logging utilities
├── models/                 # Data models & entities
├── public/                 # Static assets (favicon, etc.)
├── repository/             # Data Access Layer (DAL)
├── scripts/                # Deployment & utility scripts
├── services/               # Business Logic Layer (BLL)
├── templates/              # Email HTML templates
├── tests/                  # Unit & integration tests
├── testscripts/            # Test automation scripts
├── utils/                  # Helper utilities
└── main.go                 # Application entry point
```

## 🏗️ Architecture Layers

### 1. Presentation Layer (`handlers/`)
HTTP request handlers that process incoming requests and return responses.

**Responsibilities:**
- Request validation
- Response formatting
- HTTP status codes
- Authentication & authorization checks

**Key Files:**
- `auth_handler.go` - Authentication & authorization
- `email_handler.go` - Email sending endpoints
- `*_handler.go` - Resource-specific handlers

### 2. Business Logic Layer (`services/`)
Core business logic and orchestration between repositories.

**Responsibilities:**
- Business rules enforcement
- Data transformation
- External service integration
- Email processing & batching

**Key Files:**
- `auth_service.go` - User authentication
- `email_service.go` - Email operations
- `permission_service.go` - RBAC permissions
- `*_service.go` - Domain-specific services

### 3. Data Access Layer (`repository/`)
Database operations and data persistence.

**Responsibilities:**
- CRUD operations
- Query building
- Transaction management
- Database connection handling

**Pattern:** Repository Pattern
- Each entity has its own repository
- Interface-based for testability
- Factory pattern for initialization

### 4. Models (`models/`)
Data structures and domain entities.

**Types:**
- Database entities (matching DB schema)
- Request/Response DTOs
- Business domain models

## 🔄 Request Flow

```
HTTP Request
    ↓
[handlers/] → HTTP Handler
    ↓ (calls)
[services/] → Business Service
    ↓ (calls)
[repository/] → Repository
    ↓ (queries)
[PostgreSQL Database]
    ↓ (returns)
[models/] → Data Models
    ↓ (transforms)
HTTP Response
```

## 🔐 Security Architecture

### Authentication
- JWT-based authentication
- Refresh token mechanism
- Secure password hashing (bcrypt)

### Authorization (RBAC)
- Role-Based Access Control
- Permission-based middleware
- Fine-grained resource access

### Endpoints:
- `handlers/auth_handler.go` - Login, logout, token refresh
- `handlers/permission_handler.go` - RBAC management
- `services/auth_service.go` - Auth business logic
- `services/permission_service.go` - Permission checks

## 📧 Email Processing

### Components:
1. **SMTP Clients** (`services/smtp_client.go`)
   - Multiple SMTP configurations
   - Template-based emails
   - Retry logic

2. **Email Batching** (`services/email_batcher.go`)
   - Batch processing
   - Rate limiting
   - Queue management

3. **Email Fetching** (`services/mail_fetcher.go`)
   - IMAP integration
   - Auto-fetcher daemon
   - Email parsing

4. **Templates** (`templates/`)
   - HTML email templates
   - Dynamic content injection

## 🗄️ Database

**Type:** PostgreSQL
**ORM:** None (using standard library `database/sql`)

### Migrations:
- Location: `database/migrations/`
- Format: `V{version}__{description}.sql`
- Management: Custom migration manager

### Key Tables:
- `gebruikers` - Users
- `aanmeldingen` - Registrations
- `contact` - Contact forms
- `rbac_*` - RBAC system
- `chat_*` - Chat system
- `events` - Event tracking & geofencing
- `event_participants` - Event participant tracking
- Content tables (albums, photos, videos, etc.)

## 🔌 External Integrations

### Services:
- **Cloudinary** - Image & media management
- **Telegram Bot** - Notifications
- **Redis** - Rate limiting & caching
- **ELK Stack** - Centralized logging
- **Prometheus** - Metrics collection

### Email Providers:
- Multiple SMTP configurations
- IMAP for incoming emails
- Template-based rendering

## 📊 Monitoring & Logging

### Logging (`logger/`)
- Structured JSON logging
- Multiple log levels
- ELK integration
- Request tracing

### Metrics (`handlers/metrics_handler.go`)
- Prometheus metrics
- Email statistics
- Rate limit tracking
- Health checks

## 🧪 Testing

**Location:** `tests/`

**Types:**
- Unit tests (service layer)
- Integration tests (handlers)
- Repository tests (with test DB)
- Email processing tests

**Coverage:**
- Core business logic
- HTTP handlers
- Repository operations
- Authentication & authorization

## 🚀 Deployment

### Docker:
- `Dockerfile` - Multi-stage build
- `docker-compose.yml` - Production setup
- `docker-compose.dev.yml` - Development setup

### Platforms:
- **Render** - Primary deployment (`render.yaml`)
- **Local** - Development environment
- **Docker** - Containerized deployment

### Configuration:
- Environment variables (`.env`)
- Configuration struct (`config/`)
- Secrets management

## 📝 API Endpoints

### Public Endpoints:
- `POST /api/contact-email` - Contact form
- `POST /api/aanmelding-email` - Registration form
- `GET /api/health` - Health check
- `GET /api/program` - Program schedule
- `GET /api/total-steps` - Steps counter
- `GET /api/events` - Events overview
- `GET /api/events/active` - Active event details

### Protected Endpoints (require auth):
- `/api/auth/*` - Authentication
- `/api/admin/*` - Admin operations
- `/api/rbac/*` - RBAC management
- `/api/events/*` - Event management (admin)
- `/api/events/:id/participants` - Event tracking
- Content management endpoints

## 🔧 Configuration

### Environment Variables:
- Database connection
- SMTP credentials (multiple)
- JWT secret
- Redis connection
- External API keys
- Feature flags

### Files:
- `.env` - Local environment (gitignored)
- `.env.example` - Template
- `config/` - Configuration management

## 🎯 Design Patterns

1. **Repository Pattern** - Data access abstraction
2. **Factory Pattern** - Service & repository initialization
3. **Dependency Injection** - Loose coupling
4. **Middleware Pattern** - Request processing pipeline
5. **Observer Pattern** - WebSocket hub for real-time updates
6. **Strategy Pattern** - Multiple SMTP providers

## 📚 Key Dependencies

```go
// Web Framework
github.com/gofiber/fiber/v2

// Database
github.com/lib/pq

// Authentication
github.com/golang-jwt/jwt/v5

// External Services
github.com/cloudinary/cloudinary-go/v2
github.com/redis/go-redis/v9

// Monitoring
github.com/prometheus/client_golang
```

## 🔄 Development Workflow

1. **Local Development:**
   ```bash
   docker-compose -f docker-compose.dev.yml up
   ```

2. **Testing:**
   ```bash
   go test ./...
   ```

3. **Building:**
   ```bash
   go build -o dklemailservice
   ```

4. **Deployment:**
   - Git push triggers Render deployment
   - Docker images built & deployed

## 📖 Further Documentation

- See `docs/` folder for detailed API documentation
- `docs/FRONTEND_INTEGRATION.md` - Frontend integration guide
- `docs/DATABASE_REFERENCE.md` - Database schema reference
- `docs/AUTH_AND_RBAC.md` - Authentication & authorization guide

---

**Maintained by:** DKL Development Team
**Last Updated:** 2025-02