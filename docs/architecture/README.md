# Architecture Documentation

System architecture and design documentation for the DKL Email Service.

## Overview

The DKL Email Service is a comprehensive backend system built with Go, providing:
- RESTful API endpoints
- WebSocket real-time communication
- JWT-based authentication with RBAC
- PostgreSQL database with advanced features
- Redis caching
- Cloudinary media integration
- Email template management

## Architecture Components

### [Database Schema](./DATABASE.md)
Complete PostgreSQL 17 database design including:
- **All table structures** (30+ tables volledig gedocumenteerd)
- **Complete relationships** (Foreign keys, indexes, constraints)
- **RBAC implementation** (Roles, permissions, user_roles)
- **Gamification tables** (Achievements, badges, leaderboards)
- **CMS tables** (Albums, photos, videos, partners, sponsors)
- **Event system** (Events, registrations, participants with geofencing)
- **Migration strategy** (Versioned SQL migrations)
- **Performance optimization** (Connection pooling, query optimization)
- **Backup & monitoring** (Automated backups, health checks)

### Authentication & Security
Security architecture featuring:
- **JWT token management** (Access & refresh tokens)
- **Role-Based Access Control** (Granulaire permissies)
- **Permission system** (Resource:action pattern)
- **Session management** (Redis-backed sessions)
- **Security best practices** (Password hashing, SQL injection prevention)
- See: [Authentication API](../api/AUTHENTICATION.md), [Permissions API](../api/PERMISSIONS.md)

### WebSocket Real-time Features
Real-time communication design:
- **Connection management** (Auth via JWT query param)
- **Message broadcasting** (Steps updates, Notulen updates, Chat)
- **Channel architecture** (Multiple WebSocket endpoints)
- **Scalability considerations** (Redis pub/sub, connection pooling)
- See: [WebSocket API](../api/WEBSOCKET.md), [Steps & Gamification](../api/STEPS_GAMIFICATION.md)

### Service Layer Architecture
Business logic organization:
- **Service patterns** (EmailService, AuthService, ChatService, etc.)
- **Repository pattern** (Clean data access abstraction)
- **Dependency injection** (ServiceFactory pattern)
- **Error handling** (Structured logging, audit trails)
- **Background workers** (Email batching, auto-fetcher, newsletter)

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Web Client  │  │ Mobile Client│  │  Admin Panel │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└────────────┬──────────────┬────────────────┬───────────────┘
             │              │                │
             │ REST API     │ WebSocket      │ REST API
             │              │                │
┌────────────┴──────────────┴────────────────┴───────────────┐
│                      API Gateway Layer                       │
│  ┌────────────────────────────────────────────────────┐    │
│  │           JWT Authentication Middleware             │    │
│  │           Rate Limiting & CORS                      │    │
│  └────────────────────────────────────────────────────┘    │
└────────────┬──────────────┬────────────────┬───────────────┘
             │              │                │
┌────────────┴──────────────┴────────────────┴───────────────┐
│                    Application Layer (Go)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Handlers   │  │   Services   │  │  WebSocket   │     │
│  │              │  │              │  │   Manager    │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                  │                  │              │
│  ┌──────┴──────────────────┴──────────────────┴───────┐    │
│  │              Repository Layer                       │    │
│  │        (Data Access & Business Logic)               │    │
│  └───────────────────────────┬─────────────────────────┘    │
└──────────────────────────────┼──────────────────────────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
┌───────┴────────┐  ┌─────────┴────────┐  ┌─────────┴────────┐
│   PostgreSQL   │  │      Redis       │  │   Cloudinary     │
│   Database     │  │   Cache/Queue    │  │  Media Storage   │
│                │  │                  │  │                  │
│ • Users        │  │ • Sessions       │  │ • Images         │
│ • Permissions  │  │ • Rate Limits    │  │ • Videos         │
│ • Content      │  │ • WebSocket      │  │ • Transformations│
│ • Events       │  │   State          │  │                  │
└────────────────┘  └──────────────────┘  └──────────────────┘
```

## Technology Stack

### Backend
- **Language**: Go 1.23.0
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL 17
- **Cache**: Redis 7+
- **WebSocket**: Fiber WebSocket (Gorilla WebSocket compatible)

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Deployment**: Render.com
- **CI/CD**: GitHub Actions

### Third-Party Services
- **Media Storage**: Cloudinary
- **Email**: SMTP (configurable)
- **Monitoring**: Custom logging + ELK Stack ready

## Design Patterns

### Repository Pattern
Abstracts data access logic from business logic:

```go
type Repository interface {
    Create(ctx context.Context, entity interface{}) error
    GetByID(ctx context.Context, id int) (interface{}, error)
    Update(ctx context.Context, entity interface{}) error
    Delete(ctx context.Context, id int) error
}
```

### Service Layer Pattern
Encapsulates business logic:

```go
type UserService struct {
    repo       UserRepository
    authRepo   AuthRepository
    mailService EmailService
}

func (s *UserService) RegisterUser(ctx context.Context, req RegisterRequest) error {
    // Business logic here
}
```

### Middleware Pattern
Request/response processing pipeline:

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
}
```

## Security Architecture

### Defense in Depth

1. **Network Layer**
   - HTTPS/WSS only in production
   - Rate limiting per IP
   - CORS configuration

2. **Application Layer**
   - JWT authentication
   - RBAC authorization
   - Input validation
   - SQL injection prevention (parameterized queries)
   - XSS protection

3. **Data Layer**
   - Encrypted passwords (bcrypt)
   - Database user permissions
   - Encrypted connections
   - Regular backups

### Authentication Flow

```
┌─────────┐         ┌─────────┐         ┌──────────┐
│ Client  │         │   API   │         │ Database │
└────┬────┘         └────┬────┘         └────┬─────┘
     │                   │                    │
     │ POST /auth/login  │                    │
     │──────────────────>│                    │
     │                   │ Verify credentials │
     │                   │───────────────────>│
     │                   │                    │
     │                   │<───────────────────│
     │                   │ Generate JWT       │
     │                   │                    │
     │  {access_token,   │                    │
     │   refresh_token}  │                    │
     │<──────────────────│                    │
     │                   │                    │
     │ GET /api/resource │                    │
     │ Bearer: token     │                    │
     │──────────────────>│                    │
     │                   │ Validate JWT       │
     │                   │ Check permissions  │
     │                   │───────────────────>│
     │                   │                    │
     │     Response      │<───────────────────│
     │<──────────────────│                    │
```

## Scalability Considerations

### Horizontal Scaling
- Stateless application design
- Session storage in Redis
- Load balancer ready
- Database connection pooling

### Caching Strategy
- Redis for frequently accessed data
- In-memory caching for static content
- Cache invalidation patterns

### Database Optimization
- Proper indexing strategy
- Query optimization
- Connection pooling
- Read replicas support (future)

## Monitoring & Logging

### Logging Levels
- **DEBUG**: Development troubleshooting
- **INFO**: General information
- **WARN**: Warning conditions
- **ERROR**: Error conditions
- **FATAL**: Critical failures

### Audit Logging
All security-relevant actions are logged:
- Authentication attempts
- Permission changes
- Data modifications
- Admin actions

### Metrics
Key performance indicators:
- Request latency
- Error rates
- Database query performance
- WebSocket connections
- Cache hit ratio

## Error Handling Strategy

### Error Types

1. **Client Errors (4xx)**
   - Input validation errors
   - Authentication failures
   - Permission denials
   - Resource not found

2. **Server Errors (5xx)**
   - Database connection errors
   - External service failures
   - Unexpected exceptions

### Error Response Format

```json
{
  "success": false,
  "error": "Human-readable error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "validation details"
  },
  "timestamp": "2025-01-08T10:00:00Z"
}
```

## Development Principles

### Code Organization
- Clean architecture
- Separation of concerns
- DRY (Don't Repeat Yourself)
- SOLID principles

### Testing Strategy
- Unit tests for business logic
- Integration tests for APIs
- E2E tests for critical flows
- Mock external dependencies

### Documentation Standards
- Code comments for complex logic
- API documentation (this docs folder)
- Architecture decision records
- Changelog maintenance

## Future Enhancements

### Planned Features
- GraphQL API support
- gRPC for internal services
- Event-driven architecture
- Microservices migration path
- Kubernetes deployment
- Service mesh integration

### Performance Improvements
- Database read replicas
- CDN integration
- Advanced caching strategies
- Query optimization
- Batch processing

---

For detailed information on specific components, refer to the individual architecture documents listed above.