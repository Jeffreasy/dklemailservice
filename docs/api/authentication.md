# Authentication API

JWT-based authentication with Role-Based Access Control (RBAC).

## Endpoints

### Register User

Register a new user account.

**Endpoint:** `POST /api/auth/register`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "naam": "John Doe",
  "telefoon": "+31612345678"
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "naam": "John Doe",
      "created_at": "2025-01-08T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**Validation Rules:**
- Email must be valid and unique
- Password minimum 8 characters
- Required fields: email, password, naam

---

### Login

Authenticate and receive JWT tokens.

**Endpoint:** `POST /api/auth/login`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "naam": "John Doe",
      "roles": ["user"]
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Missing required fields

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
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or expired refresh token

---

### Get Current User

Retrieve authenticated user information.

**Endpoint:** `GET /api/auth/me`

**Authentication:** Required (JWT)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "naam": "John Doe",
    "telefoon": "+31612345678",
    "roles": ["user", "moderator"],
    "permissions": ["read:users", "write:posts"],
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
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
  "message": "Successfully logged out"
}
```

---

### Update Password

Change user password.

**Endpoint:** `PUT /api/auth/password`

**Authentication:** Required (JWT)

**Request Body:**

```json
{
  "current_password": "OldPassword123!",
  "new_password": "NewSecurePassword456!"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Password updated successfully"
}
```

**Error Responses:**
- `401 Unauthorized` - Current password is incorrect
- `400 Bad Request` - New password doesn't meet requirements

---

## JWT Token Structure

### Access Token

**Expiration:** 24 hours

**Claims:**
```json
{
  "sub": "1",
  "email": "user@example.com",
  "roles": ["user"],
  "permissions": ["read:users"],
  "exp": 1640000000,
  "iat": 1639913600
}
```

### Refresh Token

**Expiration:** 7 days

**Claims:**
```json
{
  "sub": "1",
  "type": "refresh",
  "exp": 1640518400,
  "iat": 1639913600
}
```

---

## Role-Based Access Control (RBAC)

### Default Roles

1. **admin** - Full system access
2. **moderator** - Content moderation
3. **user** - Standard user access
4. **guest** - Limited read-only access

### Permission Format

Permissions follow the pattern: `action:resource`

**Examples:**
- `read:users` - Can read user data
- `write:posts` - Can create/edit posts
- `delete:comments` - Can delete comments
- `manage:roles` - Can assign roles

### Checking Permissions

Include the JWT token in requests to protected endpoints:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/protected-endpoint
```

The server validates:
1. Token signature and expiration
2. User roles
3. Required permissions for the endpoint

---

## Security Best Practices

### Token Storage

**Frontend:**
- Store access token in memory (Redux/Context)
- Store refresh token in httpOnly cookie (preferred) or secure localStorage
- Never expose tokens in URLs

**Example (JavaScript):**
```javascript
// Store access token in memory
const [accessToken, setAccessToken] = useState(null);

// Refresh token in httpOnly cookie (handled by server)
// Or secure localStorage as fallback
localStorage.setItem('refresh_token', refreshToken);
```

### Token Refresh Flow

1. Access token expires after 24 hours
2. Use refresh token to get new access token
3. If refresh token expired, require re-login
4. Implement automatic token refresh before expiration

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
  setAccessToken(data.data.access_token);
  localStorage.setItem('refresh_token', data.data.refresh_token);
}
```

### Password Requirements

- Minimum 8 characters
- Must include:
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one number
  - At least one special character

### Rate Limiting

Authentication endpoints are rate-limited:
- **Login:** 5 attempts per 15 minutes per IP
- **Register:** 3 accounts per hour per IP
- **Password Reset:** 3 attempts per hour per email

---

## Error Codes

| Code | Description |
|------|-------------|
| `AUTH_001` | Invalid credentials |
| `AUTH_002` | Token expired |
| `AUTH_003` | Invalid token |
| `AUTH_004` | Missing authorization header |
| `AUTH_005` | Insufficient permissions |
| `AUTH_006` | Account locked |
| `AUTH_007` | Email already exists |
| `AUTH_008` | Invalid refresh token |

---

## Examples

### Complete Authentication Flow

```javascript
// 1. Login
const loginResponse = await fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});

const { data } = await loginResponse.json();
const { access_token, refresh_token } = data;

// 2. Make authenticated request
const userResponse = await fetch('/api/auth/me', {
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});

// 3. Refresh token when needed
const refreshResponse = await fetch('/api/auth/refresh', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ refresh_token })
});

// 4. Logout
await fetch('/api/auth/logout', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${access_token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ refresh_token })
});
```

---

For more information on RBAC implementation, see [Permissions API](./PERMISSIONS.md).