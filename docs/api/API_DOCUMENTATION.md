
# DKL Email Service API Documentation

Complete API reference for the DKL Email Service - a comprehensive event management and participant tracking system.

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [Events & Registrations](#events--registrations)
- [Steps & Gamification](#steps--gamification)
- [Content Management System (CMS)](#content-management-system-cms)
- [Notifications](#notifications)
- [Permissions & RBAC](#permissions--rbac)
- [Auto Responses](#auto-responses)
- [WebSocket APIs](#websocket-apis)
- [Error Handling](#error-handling)
- [Common Patterns](#common-patterns)
- [Testing](#testing)

## Overview

The DKL Email Service provides a comprehensive API for managing events, participant registrations, step tracking, gamification features, and content management. The system supports dual user types (admin/staff and participants) with role-based access control.

### Base URLs

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

### Key Features

- **Dual Authentication System**: Separate login flows for admin/staff and participants
- **Event Management**: Complete event lifecycle from creation to completion
- **Real-time Tracking**: WebSocket-based step tracking and leaderboard updates
- **Gamification**: Achievements, badges, and leaderboards
- **Content Management**: Videos, photos, partners, sponsors, and more
- **Notifications**: Real-time user notifications and email alerts
- **RBAC**: Comprehensive role-based access control system

### Architecture

The API follows RESTful principles with JSON responses and JWT authentication. Real-time features are implemented via WebSocket connections.

## Getting Started

### Prerequisites

- Valid JWT token for authenticated endpoints
- API key for metrics and special endpoints
- WebSocket support for real-time features

### Authentication Flow

1. **Login**: Obtain JWT token via `/api/auth/login`
2. **Use Token**: Include in `Authorization: Bearer <token>` header
3. **Refresh**: Use refresh token to get new access token before expiration

### Quick Example

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Use API
curl http://localhost:8080/api/events \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

# Authentication

JWT-based authentication with Role-Based Access Control (RBAC) supporting dual user types: traditional admin/staff accounts and participant app users.

## User Types

The system supports two distinct user types with different authentication and authorization models:

### 1. Gebruiker (Admin/Staff) - Traditional Admin Accounts
- **Table:** `gebruikers`
- **Authentication:** Email + password hash
- **Authorization:** Full RBAC system with roles and permissions
- **Access:** Admin panel, CMS, user management
- **JWT Claims:** `rbac_active: true`, roles array from database

### 2. Participant - App Users
- **Table:** `participants`
- **Authentication:** Email + password hash (only for full accounts)
- **Authorization:** Fixed permissions via `participantPermissionsMap`
- **Access:** Mobile app features (steps tracking, achievements, etc.)
- **JWT Claims:** `rbac_active: false`, fixed `["participant_user"]` role
- **Account Types:**
  - **Full Account:** `account_type = 'full'`, `has_app_access = true`, linked to `gebruikers` table
  - **Temporary Account:** `account_type = 'temporary'`, no app access, event-only registration

## Endpoints

### Public Registration

Register a new participant for an event with optional full account creation.

**Endpoint:** `POST /api/public/aanmelden`

**Authentication:** Not required

**Request Body:**

```json
{
  "naam": "John Doe",
  "email": "user@example.com",
  "telefoon": "+31612345678",
  "rol": "Deelnemer",
  "afstand": "6 KM",
  "ondersteuning": "Nee",
  "bijzonderheden": "",
  "heeft_vervoer": false,
  "terms": true,
  "want_account": true,
  "wachtwoord": "SecurePassword123!",
  "event_id": "f1a75cc7-303e-4207-b501-8eea557bff33"
}
```

**Response:** `201 Created` (Full Account)

```json
{
  "success": true,
  "message": "Je bent succesvol ingeschreven met een volledig account! Je hebt nu toegang tot de DKL Step App.",
  "participantID": "uuid",
  "registrationID": "uuid",
  "accountType": "full",
  "hasAppAccess": true,
  "gebruikerID": "uuid",
  "eventName": "De Koninklijke Loop 2026",
  "eventDate": "01-06-2026"
}
```

**Response:** `201 Created` (Temporary Account)

```json
{
  "success": true,
  "message": "Je bent succesvol ingeschreven voor het evenement.",
  "participantID": "uuid",
  "registrationID": "uuid",
  "accountType": "temporary",
  "hasAppAccess": false,
  "eventName": "De Koninklijke Loop 2026",
  "eventDate": "01-06-2026"
}
```

**Validation Rules:**
- Email must be valid and unique (for full accounts)
- Password minimum 8 characters (required for full accounts)
- Required fields: naam, email, rol, afstand, ondersteuning, heeft_vervoer, terms
- Phone required for Begeleider and Vrijwilliger roles
- afstand must be one of: "2.5 KM", "6 KM", "10 KM", "15 KM"
- rol must be one of: "Deelnemer", "Begeleider", "Vrijwilliger"
- ondersteuning must be one of: "Ja", "Nee", "Anders"
- wantAccount: true = full account, false = temporary

---

### Upgrade to Full Account

Upgrade a temporary registration to a full account with app access.

**Endpoint:** `POST /api/public/upgrade-to-full-account`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "wachtwoord": "NewSecurePassword123!"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Je account is geüpgraded! Je hebt nu toegang tot de DKL Step App.",
  "participantID": "uuid",
  "gebruikerID": "uuid",
  "hasAppAccess": true
}
```

---

### Login

Authenticate and receive JWT tokens. Supports both admin/staff accounts and participant accounts.

**Endpoint:** `POST /api/auth/login`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "wachtwoord": "SecurePassword123!"
}
```

**Response:** `200 OK` (Admin/Staff Account)

```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "naam": "John Doe",
    "permissions": [
      {"resource": "admin", "action": "access"},
      {"resource": "user", "action": "read"}
    ],
    "roles": [
      {
        "id": "role-uuid",
        "name": "admin",
        "description": "Full system access"
      }
    ],
    "is_actief": true
  }
}
```

**Response:** `200 OK` (Participant Account)

```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "naam": "",
    "permissions": [
      {"resource": "app", "action": "access"},
      {"resource": "steps", "action": "view_own"},
      {"resource": "participant", "action": "view_own"}
    ],
    "roles": [
      {
        "id": "participant-role",
        "name": "participant_user",
        "description": "Participant with app access"
      }
    ],
    "is_actief": true
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Missing required fields
- `403 Forbidden` - Account inactive or no app access

---

### Refresh Token

Obtain a new access token using a refresh token.

**Endpoint:** `POST /api/auth/refresh`

**Authentication:** Not required (uses refresh token)

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or expired refresh token

---

### Get Current User Profile

Retrieve authenticated user information and permissions.

**Endpoint:** `GET /api/auth/profile`

**Authentication:** Required (JWT)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK` (Admin/Staff User)

```json
{
  "id": "uuid",
  "naam": "John Doe",
  "email": "user@example.com",
  "permissions": [
    {"resource": "admin", "action": "access"},
    {"resource": "user", "action": "read"}
  ],
  "roles": [
    {
      "id": "role-uuid",
      "name": "admin",
      "description": "Full system access",
      "assigned_at": "2025-01-01T10:00:00Z",
      "is_active": true
    }
  ],
  "is_actief": true,
  "laatste_login": "2025-01-08T10:00:00Z",
  "created_at": "2025-01-01T10:00:00Z"
}
```

---

### Logout

Invalidate the refresh token and logout.

**Endpoint:** `POST /api/auth/logout`

**Authentication:** Required (JWT)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "message": "Logout succesvol"
}
```

---

### Reset Password

Change user password for authenticated users.

**Endpoint:** `POST /api/auth/reset-password`

**Authentication:** Required (JWT)

**Request Body:**

```json
{
  "huidig_wachtwoord": "OldPassword123!",
  "nieuw_wachtwoord": "NewSecurePassword456!"
}
```

**Response:** `200 OK`

```json
{
  "message": "Wachtwoord succesvol gewijzigd"
}
```

**Error Responses:**
- `401 Unauthorized` - Current password is incorrect
- `400 Bad Request` - New password doesn't meet requirements

---

### Forgot Password

Request a password reset token via email. Supports both admin/staff and participant accounts.

**Endpoint:** `POST /api/auth/forgot-password`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Response:** `200 OK` (Always returns success for security - no email enumeration)

```json
{
  "message": "Als er een account bestaat met dit email adres, ontvang je binnenkort een email met instructies om je wachtwoord te resetten."
}
```

**Error Responses:**
- `429 Too Many Requests` - Rate limit exceeded

---

### Reset Password with Token

Reset password using a token received via email.

**Endpoint:** `POST /api/auth/reset-password-with-token`

**Authentication:** Not required

**Request Body:**

```json
{
  "token": "reset_token_from_email",
  "new_password": "NewSecurePassword123!"
}
```

**Response:** `200 OK`

```json
{
  "message": "Je wachtwoord is succesvol gewijzigd. Je kunt nu inloggen met je nieuwe wachtwoord."
}
```

**Error Responses:**
- `400 Bad Request` - Invalid or expired token, or password doesn't meet requirements

---

## JWT Token Structure

### Access Token

**Expiration:** 20 minutes (configurable via `JWT_TOKEN_EXPIRY` environment variable)

**Claims Structure:**

```go
type JWTClaims struct {
    Email      string   `json:"email"`
    Roles      []string `json:"roles"`       // RBAC roles from user_roles table
    RBACActive bool     `json:"rbac_active"` // Indicates if RBAC system is active
    jwt.RegisteredClaims
}
```

**Admin/Staff Token Claims:**
```json
{
  "email": "admin@example.com",
  "roles": ["admin", "staff"],
  "rbac_active": true,
  "exp": 1640000000,
  "iat": 1639913600,
  "nbf": 1639913600,
  "iss": "dklemailservice",
  "sub": "user-uuid"
}
```

**Participant Token Claims:**
```json
{
  "email": "participant@example.com",
  "roles": ["participant_user"],
  "rbac_active": false,
  "exp": 1640000000,
  "iat": 1639913600,
  "nbf": 1639913600,
  "iss": "dklemailservice",
  "sub": "participant-uuid"
}
```

### Refresh Token

**Expiration:** 7 days

**Database Storage:** Refresh tokens are stored in the `refresh_tokens` table with:
- `owner_id`: User or participant UUID
- `token`: Base64-encoded random bytes
- `expires_at`: Expiration timestamp
- `is_revoked`: Revocation flag

---

## Role-Based Access Control (RBAC)

### User Types and Permissions

#### Admin/Staff Users (Gebruikers)
- Stored in `gebruikers` table
- Full RBAC with roles and permissions from database
- Permissions checked via `PermissionService.HasPermission()`
- Roles assigned via `user_roles` table

#### Participant Users
- Stored in `participants` table
- Fixed role: `participant_user`
- Predefined permissions from `participantPermissionsMap`
- Limited to participant-specific actions

### Permission Format

Permissions follow the pattern: `resource:action`

#### Admin/Staff Permission Examples:
- `admin:access` - Full admin access
- `user:read` - Can read user data
- `user:write` - Can create/edit users
- `user:delete` - Can delete users
- `event:read` - Can read event data
- `event:write` - Can create/edit events

#### Participant Permission Examples:
- `app:access` - Can access the mobile app
- `steps:view_own` - Can view own step data
- `steps:create` - Can create step entries
- `participant:view_own` - Can view own participant data
- `participant:update_own` - Can update own participant data
- `leaderboard:view` - Can view leaderboards
- `events:view` - Can view events
- `events:register` - Can register for events

### RBAC Roles

#### System Roles (Auto-assigned)

| Role | Auto-Assigned | Description | Permissions |
|------|---------------|-------------|-------------|
| `participant_user` | ✅ Always (full accounts) | Basic participant access | 8 core permissions |
| `participant_guide` | ✅ If event role = "Begeleider" | Guide with moderation rights | + community:moderate |
| `participant_volunteer` | ✅ If event role = "Vrijwilliger" | Volunteer with support rights | + event:support |

#### Custom Admin Roles

| Role | Description | Typical Permissions |
|------|-------------|-------------------|
| `admin` | Full system access | All permissions |
| `staff` | Staff access | Most permissions except user management |
| `moderator` | Content moderation | Limited permissions |

### Permission Checking

Include the JWT token in requests to protected endpoints:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/protected-endpoint
```

The server validates:
1. Token signature and expiration
2. User authentication (gebruiker or participant)
3. Required permissions for the endpoint via middleware
4. For participants: Uses predefined permission list
5. For admin/staff: Checks database RBAC tables

**Permission Check Flow:**
```
Request → AuthMiddleware → PermissionMiddleware → Handler
              ↓                    ↓
           JWT Check         Permission Check
           User Load         Redis Cache Lookup
                          Database Fallback
```

---

## Security Best Practices

### Token Storage Security

Secure token storage is critical for maintaining authentication security across different platforms. Below are platform-specific guidelines for the website, admin panel, and React Native applications.

#### 🌐 **Web Browsers (Website & Admin Panel)**

**Access Token:**
- Store in memory only (Redux, Context, Zustand, etc.)
- Never persist to localStorage or sessionStorage
- Clear on page refresh/window close

**Refresh Token:**
- Store in httpOnly, secure, SameSite cookies (preferred)
- Fallback: secure localStorage with encryption
- Never store in sessionStorage (cleared on tab close)

**Implementation Example:**
```javascript
// React Context + Cookies (Recommended)
import Cookies from 'js-cookie';

const AuthContext = createContext();

function AuthProvider({ children }) {
  const [accessToken, setAccessToken] = useState(null);

  useEffect(() => {
    // Load refresh token from httpOnly cookie (server-side)
    // Access token loaded from memory only
  }, []);

  const login = async (credentials) => {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials),
      credentials: 'include' // Include httpOnly cookies
    });

    const data = await response.json();
    setAccessToken(data.token); // Store in memory only

    // Refresh token automatically stored in httpOnly cookie by server
  };

  const logout = async () => {
    await fetch('/api/auth/logout', {
      method: 'POST',
      credentials: 'include'
    });
    setAccessToken(null);
  };

  return (
    <AuthContext.Provider value={{ accessToken, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
```

**Security Considerations:**
- Use `credentials: 'include'` for all authenticated requests
- Implement automatic token refresh before expiration
- Clear tokens on logout and handle 401 responses
- Use HTTPS only (enforce with HSTS headers)

#### 📱 **React Native iOS**

**Access Token:**
- Keychain Services with `kSecAttrAccessibleWhenUnlocked`
- Never store in AsyncStorage or plain text

**Refresh Token:**
- Keychain Services with `kSecAttrAccessibleAfterFirstUnlock`
- Encrypt with device-specific keys

**Implementation Example:**
```javascript
// iOS Keychain Storage
import * as Keychain from 'react-native-keychain';

class TokenManager {
  static async storeTokens(accessToken, refreshToken) {
    try {
      // Store access token (unlocked when device is unlocked)
      await Keychain.setGenericPassword(
        'accessToken',
        accessToken,
        {
          service: 'com.dkl.app',
          accessControl: Keychain.ACCESS_CONTROL.BIOMETRY_ANY,
          accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED
        }
      );

      // Store refresh token (accessible after first unlock)
      await Keychain.setGenericPassword(
        'refreshToken',
        refreshToken,
        {
          service: 'com.dkl.app.refresh',
          accessible: Keychain.ACCESSIBLE.AFTER_FIRST_UNLOCK
        }
      );
    } catch (error) {
      console.error('Failed to store tokens:', error);
    }
  }

  static async getAccessToken() {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: 'com.dkl.app'
      });
      return credentials ? credentials.password : null;
    } catch (error) {
      console.error('Failed to get access token:', error);
      return null;
    }
  }

  static async getRefreshToken() {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: 'com.dkl.app.refresh'
      });
      return credentials ? credentials.password : null;
    } catch (error) {
      console.error('Failed to get refresh token:', error);
      return null;
    }
  }

  static async clearTokens() {
    try {
      await Keychain.resetGenericPassword({ service: 'com.dkl.app' });
      await Keychain.resetGenericPassword({ service: 'com.dkl.app.refresh' });
    } catch (error) {
      console.error('Failed to clear tokens:', error);
    }
  }
}

// Usage in Auth Context
const refreshAccessToken = async () => {
  const refreshToken = await TokenManager.getRefreshToken();
  if (!refreshToken) {
    logout();
    return;
  }

  try {
    const response = await fetch('/api/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });

    if (response.ok) {
      const data = await response.json();
      await TokenManager.storeTokens(data.token, data.refresh_token);
      setAccessToken(data.token);
    } else {
      logout();
    }
  } catch (error) {
    logout();
  }
};
```

**iOS Security Features:**
- Face ID/Touch ID integration for token access
- Device encryption protection
- Automatic token cleanup on app uninstall
- Background app refresh protection

#### 🤖 **React Native Android**

**Access Token:**
- EncryptedSharedPreferences or Android Keystore
- Use AndroidX Security library

**Refresh Token:**
- Android Keystore with KeyStore protection
- AES encryption with device-specific keys

**Implementation Example:**
```javascript
// Android Secure Storage
import EncryptedStorage from 'react-native-encrypted-storage';

class TokenManager {
  static async storeTokens(accessToken, refreshToken) {
    try {
      // Store access token with encryption
      await EncryptedStorage.setItem(
        'access_token',
        JSON.stringify({
          token: accessToken,
          timestamp: Date.now()
        })
      );

      // Store refresh token with additional encryption
      await EncryptedStorage.setItem(
        'refresh_token',
        JSON.stringify({
          token: refreshToken,
          timestamp: Date.now()
        })
      );
    } catch (error) {
      console.error('Failed to store tokens:', error);
    }
  }

  static async getAccessToken() {
    try {
      const data = await EncryptedStorage.getItem('access_token');
      if (data) {
        const parsed = JSON.parse(data);
        // Check if token is not too old (optional)
        return parsed.token;
      }
      return null;
    } catch (error) {
      console.error('Failed to get access token:', error);
      return null;
    }
  }

  static async getRefreshToken() {
    try {
      const data = await EncryptedStorage.getItem('refresh_token');
      if (data) {
        const parsed = JSON.parse(data);
        return parsed.token;
      }
      return null;
    } catch (error) {
      console.error('Failed to get refresh token:', error);
      return null;
    }
  }

  static async clearTokens() {
    try {
      await EncryptedStorage.removeItem('access_token');
      await EncryptedStorage.removeItem('refresh_token');
    } catch (error) {
      console.error('Failed to clear tokens:', error);
    }
  }
}

// Alternative: Android Keystore (more secure)
import RNSecureStorage from 'rn-secure-storage';

class SecureTokenManager {
  static async storeTokens(accessToken, refreshToken) {
    try {
      await RNSecureStorage.setItem('access_token', accessToken, {
        keychainService: 'dkl_app_tokens'
      });
      await RNSecureStorage.setItem('refresh_token', refreshToken, {
        keychainService: 'dkl_app_refresh'
      });
    } catch (error) {
      console.error('Failed to store tokens securely:', error);
    }
  }

  static async getAccessToken() {
    try {
      return await RNSecureStorage.getItem('access_token', {
        keychainService: 'dkl_app_tokens'
      });
    } catch (error) {
      return null;
    }
  }

  static async clearTokens() {
    try {
      await RNSecureStorage.removeItem('access_token', {
        keychainService: 'dkl_app_tokens'
      });
      await RNSecureStorage.removeItem('refresh_token', {
        keychainService: 'dkl_app_refresh'
      });
    } catch (error) {
      console.error('Failed to clear tokens:', error);
    }
  }
}
```

**Android Security Features:**
- Android Keystore protection
- Device encryption (File-Based Encryption)
- Biometric authentication integration
- Automatic token cleanup on app data clear

#### 🔄 **Token Refresh Strategy**

**Automatic Refresh (All Platforms):**
```javascript
// Universal token refresh logic
class AuthManager {
  constructor() {
    this.refreshPromise = null;
    this.tokenExpiryBuffer = 5 * 60 * 1000; // 5 minutes before expiry
  }

  async getValidAccessToken() {
    const token = await TokenManager.getAccessToken();

    if (!token) {
      throw new Error('No access token available');
    }

    // Check if token is close to expiry
    const decoded = jwt_decode(token);
    const now = Date.now() / 1000;
    const timeUntilExpiry = decoded.exp - now;

    if (timeUntilExpiry < this.tokenExpiryBuffer) {
      return this.refreshAccessToken();
    }

    return token;
  }

  async refreshAccessToken() {
    // Prevent multiple simultaneous refresh requests
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = this._performRefresh();

    try {
      const newToken = await this.refreshPromise;
      return newToken;
    } finally {
      this.refreshPromise = null;
    }
  }

  async _performRefresh() {
    const refreshToken = await TokenManager.getRefreshToken();

    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await fetch('/api/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken })
    });

    if (!response.ok) {
      // Refresh failed - logout user
      await this.logout();
      throw new Error('Token refresh failed');
    }

    const data = await response.json();
    await TokenManager.storeTokens(data.token, data.refresh_token);

    return data.token;
  }

  async logout() {
    try {
      await TokenManager.clearTokens();
      // Navigate to login screen
    } catch (error) {
      console.error('Logout error:', error);
    }
  }
}
```

#### 🚨 **Critical Security Warnings**

**Never:**
- Store tokens in plain text files
- Log tokens in console or error messages
- Send tokens over non-HTTPS connections
- Store tokens in URL parameters
- Cache tokens in browser memory without automatic cleanup

**Always:**
- Use platform-specific secure storage
- Implement automatic token refresh
- Handle token expiration gracefully
- Clear tokens on logout
- Validate token signatures on the server
- Use short-lived access tokens (20 minutes)
- Rotate refresh tokens on each refresh

**Platform-Specific Considerations:**
- **Web**: Use httpOnly cookies when possible, implement CSRF protection
- **iOS**: Enable Face ID/Touch ID for enhanced security
- **Android**: Use Android Keystore for maximum protection
- **All Platforms**: Implement certificate pinning for API calls

### Token Refresh Flow

1. Access token expires after 20 minutes (configurable)
2. Use refresh token to get new access token
3. If refresh token expired (7 days), require re-login
4. Implement automatic token refresh before expiration

### Access Token Rotation

**Server-side Token Storage:** Access tokens are now stored in the `access_tokens` database table for server-side validation and revocation.

**Token Rotation Features:**
- Access tokens are stored server-side with expiration timestamps
- Tokens can be revoked before natural expiration
- All user access tokens are revoked on logout
- Middleware validates tokens against database in addition to JWT signature validation

**Database Table:**
```sql
CREATE TABLE access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL, -- Can be gebruiker.id or participant.id
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP NULL,
    is_revoked BOOLEAN DEFAULT FALSE
);
```

**Security Benefits:**
- Immediate token revocation on logout
- Protection against replay attacks
- Enhanced session control
- Compliance with security best practices

**Example:**
```javascript
async function refreshAccessToken() {
  const refreshToken = localStorage.getItem('refresh_token');

  const response = await fetch('/api/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken })
  });

  const data = await response.json();
  setAccessToken(data.token);
  localStorage.setItem('refresh_token', data.refresh_token);
}
```

### Password Requirements

- Minimum 8 characters
- For full accounts: Required during registration and upgrade
- For admin/staff accounts: Set during user creation by administrators

### Rate Limiting

Authentication endpoints are rate-limited:
- **Login:** Rate limited by email address (configurable)
- **Public Registration:** No specific rate limiting
- **Password Reset:** No specific rate limiting beyond authentication requirements
- **Forgot Password:** Rate limited by email address (3 requests per hour)

---

## Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| `400` | `INVALID_INPUT` | Invalid request data |
| `400` | `VALIDATION_ERROR` | Validation failed |
| `401` | `INVALID_CREDENTIALS` | Invalid email or password |
| `401` | `NO_AUTH_HEADER` | Missing Authorization header |
| `401` | `INVALID_AUTH_HEADER` | Invalid Authorization header format |
| `401` | `TOKEN_EXPIRED` | JWT token has expired |
| `401` | `TOKEN_MALFORMED` | JWT token is malformed |
| `401` | `TOKEN_SIGNATURE_INVALID` | JWT token signature invalid |
| `401` | `REFRESH_TOKEN_INVALID` | Invalid or expired refresh token |
| `403` | `USER_INACTIVE` | User account is inactive |
| `403` | `UNAUTHORIZED` | User not authorized for this action |
| `403` | `FORBIDDEN` | Insufficient permissions |
| `404` | `EVENT_NOT_FOUND` | Event not found |
| `409` | `EMAIL_EXISTS` | Email already exists |
| `409` | `ALREADY_REGISTERED` | Already registered for this event |
| `429` | `RATE_LIMIT_EXCEEDED` | Too many requests |

---

## Examples

### Complete Authentication Flow

```javascript
// 1. Register (for participants)
const registerResponse = await fetch('/api/public/aanmelden', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    naam: 'John Doe',
    email: 'user@example.com',
    telefoon: '+31612345678',
    rol: 'Deelnemer',
    afstand: '10km',
    ondersteuning: 'Nee',
    heeft_vervoer: false,
    terms: true,
    want_account: true,
    wachtwoord: 'SecurePassword123!',
    event_id: 'event-uuid'
  })
});

// 2. Login
const loginResponse = await fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    wachtwoord: 'SecurePassword123!'
  })
});

const loginData = await loginResponse.json();
const { token: access_token, refresh_token, user } = loginData;

// 3. Make authenticated request
const profileResponse = await fetch('/api/auth/profile', {
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});

// 4. Refresh token when needed
const refreshResponse = await fetch('/api/auth/refresh', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ refresh_token })
});

const refreshData = await refreshResponse.json();
const new_access_token = refreshData.token;
const new_refresh_token = refreshData.refresh_token;

// 5. Logout
await fetch('/api/auth/logout', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});
```

---

## Implementation Details

### Database Tables

**Core Authentication Tables:**
- `gebruikers` - Admin/staff accounts
- `participants` - Participant accounts
- `refresh_tokens` - Token storage
- `password_reset_tokens` - Password reset token storage
- `email_verification_tokens` - Email verification token storage
- `roles` - RBAC roles
- `permissions` - RBAC permissions
- `user_roles` - User-role assignments
- `role_permissions` - Role-permission assignments

### Key Services

**AuthService:**
- `Login()` - Dual authentication logic
- `ValidateToken()` - JWT validation
- `RefreshAccessToken()` - Token refresh
- `GetUserFromToken()` - User resolution

**PermissionService:**
- `HasPermission()` - Permission checking with caching
- `GetUserPermissions()` - Permission retrieval
- `GetUserRoles()` - Role retrieval

### Middleware

**AuthMiddleware:** JWT token validation
**PermissionMiddleware:** Permission checking for endpoints

### Caching Strategy

- Redis cache for permissions (5-minute TTL)
- Cache key: `perm:{user_id}:{resource}:{action}`
- Automatic cache invalidation on role/permission changes

---

# Events & Registrations

API endpoints voor het beheren van evenementen en registraties.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

## Authentication

Admin endpoints vereisen JWT authenticatie en permissies:

```
Authorization: Bearer <your_jwt_token>
```

---

## Events Management

### List Events (Public)

Haal alle events op (publiek toegankelijk).

**Endpoint:** `GET /api/events`

**Authentication:** Niet vereist

**Query Parameters:**
- `status` (optional): Filter op status (`upcoming`, `active`, `completed`, `cancelled`)
- `is_active` (optional): Filter op actieve events (true/false)

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "De Koninklijke Loop 2025",
    "description": "Jaarlijks hardloop evenement voor KWF",
    "start_time": "2025-06-15T09:00:00Z",
    "end_time": "2025-06-15T15:00:00Z",
    "status": "upcoming",
    "status_description": "Binnenkort",
    "geofences": [
      {
        "type": "start",
        "lat": 52.370216,
        "long": 4.895168,
        "radius": 50,
        "name": "Start Lijn"
      },
      {
        "type": "checkpoint",
        "lat": 52.375216,
        "long": 4.900168,
        "radius": 30,
        "name": "Checkpoint 1"
      },
      {
        "type": "finish",
        "lat": 52.380216,
        "long": 4.905168,
        "radius": 50,
        "name": "Finish"
      }
    ],
    "event_config": {
      "max_participants": 500,
      "registration_open": true,
      "allow_team_registration": true,
      "require_medical_certificate": false
    },
    "is_active": true
  }
]
```

---

### Get Active Event (Public)

Haal het huidige actieve event op.

**Endpoint:** `GET /api/events/active`

**Authentication:** Niet vereist

**Response:** `200 OK`

Retourneert hetzelfde format als List Events, maar alleen het actieve event.

**Response wanneer geen actief event:**

```json
{
  "message": "No active event found"
}
```

---

### Get Event Details (Public)

**Endpoint:** `GET /api/events/:id`

**Authentication:** Niet vereist

**Response:** `200 OK`

Retourneert event details inclusief registratie statistieken.

---

### Create Event (Admin)

**Endpoint:** `POST /api/events`

**Authentication:** Vereist - `event:write` permissie

**Request Body:**

```json
{
  "name": "De Koninklijke Loop 2025",
  "description": "Jaarlijks hardloop evenement voor KWF",
  "start_time": "2025-06-15T09:00:00Z",
  "end_time": "2025-06-15T15:00:00Z",
  "status": "upcoming",
  "geofences": [
    {
      "type": "start",
      "lat": 52.370216,
      "long": 4.895168,
      "radius": 50,
      "name": "Start Lijn"
    }
  ],
  "event_config": {
    "max_participants": 500,
    "registration_open": true
  },
  "is_active": true
}
```

**Validation:**
- `name`: Verplicht, max 255 karakters
- `start_time`: Verplicht, moet in de toekomst zijn
- `geofences`: Verplicht, minimaal 1 start en 1 finish geofence
- `status`: Optioneel, default `upcoming`

**Response:** `201 Created`

---

### Update Event (Admin)

**Endpoint:** `PUT /api/events/:id`

**Authentication:** Vereist - `event:write` permissie

**Request Body:**

```json
{
  "name": "Updated Event Name",
  "status": "active",
  "is_active": true
}
```

**Response:** `200 OK`

---

### Delete Event (Admin)

**Endpoint:** `DELETE /api/events/:id`

**Authentication:** Vereist - `event:delete` permissie

**Response:** `200 OK`

**Note:** Cascade delete verwijdert ook alle gerelateerde registraties.

---

## Event Registrations

### Get Event Registrations

Haal alle registraties voor een specifiek event op.

**Endpoint:** `GET /api/events/:id/registrations`

**Authentication:** Vereist - `event:read` permissie

**Query Parameters:**
- `status` (optional): Filter op registratie status
- `tracking_status` (optional): Filter op tracking status
- `role` (optional): Filter op deelnemer rol
- `route` (optional): Filter op afstand route

**Response:** `200 OK`

```json
{
  "event_id": "uuid",
  "event_name": "De Koninklijke Loop 2025",
  "total_registrations": 150,
  "registrations": [
    {
      "id": "uuid",
      "event_id": "uuid",
      "participant_id": "uuid",
      "participant_name": "John Doe",
      "status": "confirmed",
      "status_description": "Bevestigd",
      "tracking_status": "registered",
      "distance_route": "10km",
      "participant_role": "deelnemer",
      "participant_role_description": "Evenement deelnemer",
      "registered_at": "2025-01-01T10:00:00Z",
      "steps": 0,
      "total_distance": 0
    }
  ]
}
```

---

### Get Registration Details

**Endpoint:** `GET /api/registration/:id`

**Authentication:** Vereist - `registration:read` permissie

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "event_id": "uuid",
  "event_name": "De Koninklijke Loop 2025",
  "participant_id": "uuid",
  "participant_name": "John Doe",
  "status": "confirmed",
  "status_description": "Bevestigd",
  "tracking_status": "in_progress",
  "distance_route": "10km",
  "participant_role": "deelnemer",
  "registered_at": "2025-01-01T10:00:00Z",
  "check_in_time": "2025-06-15T08:45:00Z",
  "start_time": "2025-06-15T09:05:00Z",
  "finish_time": null,
  "steps": 5000,
  "total_distance": 3.5,
  "last_location_update": "2025-06-15T09:30:00Z"
}
```

---

### Update Registration Status

**Endpoint:** `PUT /api/registration/:id`

**Authentication:** Vereist - `registration:write` permissie

**Request Body:**

```json
{
  "status": "confirmed",
  "tracking_status": "checked_in",
  "notes": "Deelnemer is aanwezig"
}
```

**Available Statuses:**
- Registration Status: `pending`, `confirmed`, `cancelled`, `completed`
- Tracking Status: `registered`, `checked_in`, `started`, `in_progress`, `finished`, `dnf`

**Response:** `200 OK`

---

### Filter Registrations by Role

**Endpoint:** `GET /api/registration/rol/:rol`

**Authentication:** Vereist - `registration:read` permissie

**Parameters:**
- `rol`: Participant role (`deelnemer`, `vrijwilliger`, `begeleider`)

**Response:** `200 OK`

Retourneert array van registraties gefilterd op rol.

---

## Geofencing

### Geofence Types

Events gebruiken geofences voor locatie tracking tijdens het evenement.

**Geofence Structure:**

```json
{
  "type": "start|checkpoint|finish",
  "lat": 52.370216,
  "long": 4.895168,
  "radius": 50,
  "name": "Checkpoint Name"
}
```

**Types:**
- `start`: Start locatie (verplicht)
- `checkpoint`: Tussenstation (optioneel, meerdere mogelijk)
- `finish`: Finish locatie (verplicht)

**Radius:**
- Opgegeven in meters
- Aanbevolen: 30-100 meter voor checkpoints, 50-200 meter voor start/finish

---

### Location Updates During Event

**Endpoint:** `POST /api/registration/:id/location`

**Authentication:** Vereist - `registration:write` permissie

**Request Body:**

```json
{
  "lat": 52.370216,
  "long": 4.895168,
  "accuracy": 10,
  "timestamp": "2025-06-15T09:30:00Z"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "geofence_triggered": true,
  "geofence_type": "checkpoint",
  "geofence_name": "Checkpoint 1",
  "tracking_status": "in_progress",
  "total_distance": 3.5
}
```

**Automatic Status Updates:**
- Bij start geofence: status → `started`
- Bij checkpoint: distance increment
- Bij finish geofence: status → `finished`

---

## Event Status Lifecycle

### Status Flow

```
upcoming → active → completed
    ↓
cancelled (mogelijk op elk moment)
```

**Status Descriptions:**
- `upcoming`: Event is nog niet begonnen
- `active`: Event is momenteel gaande
- `completed`: Event is afgelopen
- `cancelled`: Event is geannuleerd

---

## Event Configuration

### EventConfig Object

Het `event_config` JSONB veld ondersteunt flexibele configuratie:

```json
{
  "max_participants": 500,
  "registration_open": true,
  "allow_team_registration": true,
  "require_medical_certificate": false,
  "early_bird_deadline": "2025-05-01T00:00:00Z",
  "registration_fee": 25.00,
  "age_restrictions": {
    "min_age": 16,
    "max_age": null
  },
  "routes": [
    {
      "name": "5km",
      "distance": 5.0,
      "max_participants": 200
    },
    {
      "name": "10km",
      "distance": 10.0,
      "max_participants": 200
    },
    {
      "name": "15km",
      "distance": 15.0,
      "max_participants": 100
    }
  ]
}
```

---

## Statistics & Reporting

### Event Statistics

**Endpoint:** `GET /api/events/:id/stats`

**Authentication:** Vereist - `event:read` permissie

**Response:** `200 OK`

```json
{
  "event_id": "uuid",
  "event_name": "De Koninklijke Loop 2025",
  "total_registrations": 150,
  "confirmed_count": 140,
  "cancelled_count": 5,
  "pending_count": 5,
  "checked_in_count": 120,
  "started_count": 115,
  "finished_count": 105,
  "dnf_count": 10,
  "total_steps": 1500000,
  "total_distance": 1050,
  "average_steps_per_participant": 10000,
  "by_route": {
    "5km": {
      "count": 50,
      "steps": 250000
    },
    "10km": {
      "count": 70,
      "steps": 700000
    },
    "15km": {
      "count