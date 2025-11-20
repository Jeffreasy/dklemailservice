# Authentication API

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
    owner_id UUID NOT NULL, -- Can be gebruiker or participant ID
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
    afstand: '6 KM',
    ondersteuning: 'Nee',
    heeft_vervoer: false,
    terms: true,
    want_account: true,
    wachtwoord: 'SecurePassword123!'
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

## 🔒 **Advanced Security Features**

### Rate Limiting Details

The API implements comprehensive rate limiting to prevent abuse:

| Endpoint/Flow | Limit | Window | Scope |
|---------------|-------|--------|-------|
| **Login attempts** | 5 attempts | 15 minutes | Per email address + IP |
| **Contact forms** | 5 submissions | 1 hour | Per IP address |
| **Public registration** | 3 registrations | 1 hour | Per IP address |
| **Password reset** | 3 requests | 1 hour | Per email address |
| **Forgot password** | 3 requests | 1 hour | Per email address |

**Implementation:**
```go
// Rate limiting keys are constructed as:
// "login:{email}" - for login attempts
// "contact_ip:{ip}" - for contact forms
// "aanmelding_ip:{ip}" - for registrations
```

### JWT Token Revocation Strategy

**Access Tokens:** Not stored server-side, invalidated by expiration (20 minutes)
**Refresh Tokens:** Stored in database with revocation capability

```sql
-- Refresh token structure
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL, -- Can be gebruiker.id or participant.id
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    is_revoked BOOLEAN DEFAULT FALSE
);
```

**Revocation on Logout:**
```go
// Logout revokes refresh token
func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
    // Revoke refresh token in database
    // Clear httpOnly cookie
    // Audit log the logout
}
```

### Audit Logging

Comprehensive audit logging captures all authentication and authorization events:

**Authentication Events:**
- `LOGIN_SUCCESS` / `LOGIN_FAILED`
- `LOGOUT`
- `TOKEN_REFRESHED`
- `PASSWORD_CHANGED`

**Authorization Events:**
- `ACCESS_GRANTED` / `ACCESS_DENIED`
- `ROLE_ASSIGNED` / `ROLE_REVOKED`
- `USER_CREATED` / `USER_UPDATED` / `USER_DELETED`

**Audit Event Structure:**
```json
{
  "event_type": "LOGIN_SUCCESS",
  "timestamp": "2025-11-18T12:51:38Z",
  "actor_id": "05591b87-42eb-4005-9abe-0be4a4388326",
  "actor_email": "admin@dekoninklijkeloop.nl",
  "ip_address": "172.19.0.1",
  "user_agent": "curl/8.13.0",
  "result": "SUCCESS",
  "metadata": {
    "permissions_count": 127,
    "roles_count": 2,
    "user_type": "gebruiker"
  }
}
```

### Security Headers Configuration

**CORS Configuration:**
```go
// Configured in main.go
allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
app.Use(cors.New(cors.Config{
    AllowOrigins:     allowedOrigins,
    AllowCredentials: true,
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
}))
```

**Implemented Security Headers:**
```go
// SecurityHeadersMiddleware in handlers/middleware.go
func SecurityHeadersMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Basis beveiligingsheaders
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")

        // HSTS (HTTP Strict Transport Security) - alleen in productie
        if os.Getenv("ENV") == "production" {
            c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }

        // Content Security Policy - configureerbaar via environment variable
        if csp := os.Getenv("CONTENT_SECURITY_POLICY"); csp != "" {
            c.Set("Content-Security-Policy", csp)
        } else {
            // Default CSP voor development
            c.Set("Content-Security-Policy", "default-src 'self'")
        }

        // Referrer Policy
        if rp := os.Getenv("REFERRER_POLICY"); rp != "" {
            c.Set("Referrer-Policy", rp)
        } else {
            c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        }

        // Permissions Policy (voorheen Feature Policy)
        if pp := os.Getenv("PERMISSIONS_POLICY"); pp != "" {
            c.Set("Permissions-Policy", pp)
        }

        return c.Next()
    }
}
```

---

## 👥 **RBAC Administration Endpoints**

### User Management

**List Users:**
```
GET /api/users?limit=50&offset=0
Authorization: Bearer <admin_token>
Permissions: user:read
```

**Create User:**
```
POST /api/users
Authorization: Bearer <admin_token>
Permissions: user:write
Content-Type: application/json

{
  "email": "newuser@example.com",
  "naam": "New User",
  "password": "SecurePass123!",
  "is_actief": true,
  "newsletter_subscribed": false
}
```

**Update User:**
```
PUT /api/users/{id}
Authorization: Bearer <admin_token>
Permissions: user:write
```

**Delete User:**
```
DELETE /api/users/{id}
Authorization: Bearer <admin_token>
Permissions: user:delete
```

### Role Management

**Get User Roles:**
```
GET /api/users/{id}/roles
Authorization: Bearer <admin_token>
Permissions: user:read
```

**Assign Role to User:**
```
POST /api/users/{id}/roles
Authorization: Bearer <admin_token>
Permissions: user:manage_roles
Content-Type: application/json

{
  "role_id": "d209a318-50b3-4ebe-8635-242069d40b3a",
  "expires_at": "2026-01-01T00:00:00Z" // optional
}
```

**Assign Multiple Roles:**
```
PUT /api/users/{id}/roles
Authorization: Bearer <admin_token>
Permissions: admin (full access)
Content-Type: application/json

{
  "role_ids": [
    "d209a318-50b3-4ebe-8635-242069d40b3a",
    "4ee27340-91e5-4508-b28d-c3e3db4f972f"
  ]
}
```

**Remove Role from User:**
```
DELETE /api/users/{id}/roles/{roleId}
Authorization: Bearer <admin_token>
Permissions: user:manage_roles
```

**Get User Permissions:**
```
GET /api/users/{id}/permissions
Authorization: Bearer <admin_token>
Permissions: user:read
```

---

## 🚧 **Future Enhancements (Production Recommendations)**

### Change Email
```
POST /api/auth/change-email
Authorization: Bearer <token>
Content-Type: application/json

{
  "new_email": "newemail@example.com",
  "password": "current_password"
}
```

### Delete Account (GDPR)

Verwijdert een gebruikersaccount volledig volgens GDPR richtlijnen. Alle persoonlijke gegevens worden permanent verwijderd.

**Endpoint:** `DELETE /api/auth/account`

**Authentication:** Required (JWT)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**

```json
{
  "password": "current_password",
  "reason": "optional_deletion_reason"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Je account is succesvol verwijderd. Alle persoonlijke gegevens zijn permanent verwijderd."
}
```

**Error Responses:**
- `400 Bad Request` - Missing password or invalid input
- `401 Unauthorized` - Invalid password or not authenticated
- `404 Not Found` - Account not found
- `429 Too Many Requests` - Rate limit exceeded

**Security Features:**
- Password verification required
- Rate limiting (1 request per user per hour)
- Comprehensive audit logging
- Complete data deletion (tokens, roles, personal data)
- Automatic logout after deletion

**GDPR Compliance:**
- All personal data is permanently deleted
- No data retention after account deletion
- Audit trail maintained for compliance
- User consent verification via password

### Advanced Session Management

**List Active Sessions:**
```
GET /api/auth/sessions
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "sessions": [
    {
      "id": "session-uuid",
      "device_info": {
        "browser": "Chrome",
        "browser_version": "119.0",
        "os": "Windows",
        "os_version": "",
        "device_type": "desktop",
        "platform": "Windows"
      },
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "login_time": "2025-11-18T10:30:00Z",
      "last_activity": "2025-11-18T10:45:00Z",
      "is_current": false,
      "display_name": "Chrome on Windows",
      "location_info": "Netherlands"
    }
  ]
}
```

**Revoke Specific Session:**
```
DELETE /api/auth/sessions/:sessionId
Authorization: Bearer <token>
```

**Revoke All Other Sessions:**
```
POST /api/auth/sessions/revoke-others
Authorization: Bearer <token>
```

**Enhanced Logout with Session Control:**
```
POST /api/auth/logout
Authorization: Bearer <token>
Content-Type: application/json

{
  "revoke_all_sessions": true,  // Optional: revoke all sessions
  "session_id": "session-uuid"   // Optional: revoke specific session
}
```

**Token Rotation:**
- Refresh tokens are rotated on each refresh (✅ implemented)
- Access tokens could be rotated more frequently (not implemented)

### OAuth/Social Login (Future Enhancement)

**Google OAuth:**
```
GET /api/auth/google
GET /api/auth/google/callback
```

**Microsoft OAuth:**
```
GET /api/auth/microsoft
GET /api/auth/microsoft/callback
```

### Device Fingerprinting (Future Enhancement)

```javascript
// Client-side device fingerprinting
const fingerprint = {
  userAgent: navigator.userAgent,
  language: navigator.language,
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  screenResolution: `${screen.width}x${screen.height}`,
  colorDepth: screen.colorDepth,
  platform: navigator.platform
};
```

---

## 📊 **Monitoring & Metrics**

### Authentication Metrics

The system tracks authentication metrics for monitoring:

- Login success/failure rates
- Token refresh frequency
- Rate limit hits
- Geographic login patterns
- Device type distribution

### Health Check Endpoint

```
GET /api/health
```

Returns authentication system status including:
- Database connectivity
- Redis cache status
- Rate limiter status
- SMTP configuration
- Template validation

---

## 🔧 **Configuration Reference**

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | - | **Required:** JWT signing secret (min 32 chars) |
| `JWT_TOKEN_EXPIRY` | `20m` | Access token expiration |
| `ALLOWED_ORIGINS` | `http://localhost:3000,http://localhost:5173` | CORS allowed origins |
| `EMAIL_RATE_LIMIT` | `10` | Email sending rate limit |
| `CONTACT_RATE_LIMIT` | `5` | Contact form rate limit |
| `REGISTRATION_RATE_LIMIT` | `3` | Registration rate limit |
| `RATE_LIMIT_WINDOW` | `3600` | Rate limit window in seconds |

### Database Tables

**Core Authentication Tables:**
- `gebruikers` - Admin/staff users
- `participants` - App users
- `refresh_tokens` - Token storage
- `roles` - RBAC roles
- `permissions` - RBAC permissions
- `user_roles` - User-role assignments
- `role_permissions` - Role-permission assignments

---

## 📋 **Implementation Checklist for Production**

### ✅ **Implemented Features**
- [x] Dual user authentication (admin/staff + participants)
- [x] JWT-based authentication with refresh tokens
- [x] Full RBAC system with roles and permissions
- [x] Rate limiting on authentication endpoints
- [x] Comprehensive audit logging
- [x] CORS configuration
- [x] Token revocation on logout
- [x] Platform-specific token storage guidelines
- [x] Automatic token refresh strategy
- [x] User management endpoints
- [x] Role assignment/removal endpoints
- [x] Forgot password flow
- [x] Email verification system
- [x] Multiple session management with device tracking
- [x] Session validation middleware
- [x] Selective session revocation

### 🚧 **Recommended for Production**
- [x] Account deletion (GDPR compliance)
- [x] Security headers (HSTS, CSP, etc.)
- [ ] OAuth/Social login integration
- [x] Multiple session management
- [ ] Device fingerprinting
- [ ] Advanced monitoring/metrics
- [ ] Certificate pinning
- [x] Token rotation for access tokens

This comprehensive authentication system provides enterprise-grade security while maintaining simplicity for mobile and web applications.

For more information on RBAC implementation, see the permission middleware and service implementations in the codebase.