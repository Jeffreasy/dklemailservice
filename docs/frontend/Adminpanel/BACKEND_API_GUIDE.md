# Admin Panel Backend API Gids

## 🎯 Overzicht

Deze gids beschrijft hoe het admin panel communiceert met de DKL Email Service backend API, inclusief authenticatie, autorisatie, en RBAC management.

---

## ✅ CORS Configuratie (GEFIXE!)

De backend is nu geconfigureerd om requests van het admin panel te accepteren.

### Toegestane Origins

```javascript
// In main.go (regel 369) en .env.local (regel 76)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:5174,http://localhost:8080
```

**Verificatie via Docker logs:**
```json
{
  "lvl":"INFO",
  "msg":"CORS geconfigureerd",
  "origins":["http://localhost:3000","http://localhost:5173","http://localhost:5174","http://localhost:8080"]
}
```

✅ **Admin Panel (http://localhost:5174)** kan nu veilig communiceren met de backend!

---

## 🔐 Authenticatie Flow

### 1. Login Endpoint

**Request:**
```http
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "admin@dekoninklijkeloop.nl",
  "wachtwoord": "YourPassword123"
}
```

**Response (Success):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user-uuid",
    "email": "admin@dekoninklijkeloop.nl",
    "naam": "Admin User",
    "permissions": [
      {"resource": "admin", "action": "access"},
      {"resource": "user", "action": "read"},
      {"resource": "user", "action": "write"}
    ],
    "roles": [
      {
        "id": "role-uuid",
        "name": "admin",
        "description": "Volledige systeem toegang"
      }
    ],
    "is_actief": true
  }
}
```

### 2. Token Storage (Frontend)

**Juiste implementatie:**
```typescript
// ✅ CORRECT: Access token in memory
const [accessToken, setAccessToken] = useState<string | null>(null);

// ✅ CORRECT: Refresh token in localStorage
localStorage.setItem('refresh_token', refreshToken);

// ❌ INCORRECT: Access token in localStorage (security risk!)
// localStorage.setItem('access_token', token); // DON'T DO THIS
```

### 3. API Requests

**Met Authorization header:**
```typescript
const response = await fetch('http://localhost:8080/api/auth/profile', {
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  }
});
```

### 4. Token Refresh

**Automatische refresh bij 401:**
```typescript
// In axios interceptor
apiClient.interceptors.response.use(
  response => response,
  async (error) => {
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true;
      
      // Refresh token
      const refreshToken = localStorage.getItem('refresh_token');
      const { data } = await axios.post('/api/auth/refresh', {
        refresh_token: refreshToken
      });
      
      // Update token en retry
      error.config.headers.Authorization = `Bearer ${data.token}`;
      return apiClient(error.config);
    }
    return Promise.reject(error);
  }
);
```

---

## 🛡️ RBAC Management API

### Base URL
```
http://localhost:8080/api/rbac
```

### Permissions Endpoints

#### Lijst Alle Permissions (Gegroepeerd)
```http
GET /api/rbac/permissions?group_by_resource=true&limit=1000
Authorization: Bearer <token>
```

**Response:**
```json
{
  "groups": [
    {
      "resource": "user",
      "permissions": [
        {
          "id": "uuid",
          "resource": "user",
          "action": "read",
          "description": "Gebruikers bekijken",
          "is_system_permission": true
        }
      ],
      "count": 4
    }
  ],
  "total": 127
}
```

#### Permission Aanmaken
```http
POST /api/rbac/permissions
Authorization: Bearer <token>
Content-Type: application/json

{
  "resource": "custom_resource",
  "action": "custom_action",
  "description": "Custom permission beschrijving"
}
```

#### Permission Bijwerken
```http
PUT /api/rbac/permissions/:id
Authorization: Bearer <token>

{
  "description": "Nieuwe beschrijving"
}
```

**Note:** `resource` en `action` kunnen NIET worden gewijzigd.

#### Permission Verwijderen
```http
DELETE /api/rbac/permissions/:id
Authorization: Bearer <token>
```

**Beperking:** Systeempermissions (`is_system_permission: true`) kunnen niet worden verwijderd.

---

### Roles Endpoints

#### Lijst Alle Rollen
```http
GET /api/rbac/roles?limit=50&offset=0
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "role-uuid",
    "name": "admin",
    "description": "Volledige systeem toegang",
    "is_system_role": true,
    "permissions": [
      {
        "id": "perm-uuid",
        "resource": "admin",
        "action": "access",
        "description": "Volledige admin toegang"
      }
    ],
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
]
```

#### Rol Aanmaken
```http
POST /api/rbac/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "content_editor",
  "description": "Content beheerder"
}
```

#### Rol Bijwerken
```http
PUT /api/rbac/roles/:id
Authorization: Bearer <token>

{
  "description": "Nieuwe beschrijving"
}
```

**Beperking:** Systeemrollen kunnen niet worden bewerkt.

#### Rol Verwijderen
```http
DELETE /api/rbac/roles/:id
Authorization: Bearer <token>
```

**Beperking:** Systeemrollen kunnen niet worden verwijderd.

---

### Role-Permission Management

#### Permissions Toewijzen aan Rol (Bulk)
```http
PUT /api/rbac/roles/:roleId/permissions
Authorization: Bearer <token>
Content-Type: application/json

{
  "permission_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Role permissions bijgewerkt",
  "added_count": 2,
  "removed_count": 1,
  "total_requested": 3
}
```

**Gedrag:** 
- Verwijdert eerst alle huidige permissions
- Voegt dan de nieuwe permissions toe
- Ideaal voor "Save All" functionaliteit in UI

#### Individuele Permission Toevoegen
```http
POST /api/rbac/roles/:roleId/permissions/:permissionId
Authorization: Bearer <token>
```

#### Individuele Permission Verwijderen
```http
DELETE /api/rbac/roles/:roleId/permissions/:permissionId
Authorization: Bearer <token>
```

---

### Menu Permission Management

#### Menu Permissions Ophalen
```http
GET /api/rbac/menu/permissions
Authorization: Bearer <token>
```

**Response:**
```json
{
  "permissions": [
    {
      "id": "uuid",
      "resource": "menu",
      "action": "dashboard",
      "description": "Dashboard menu item"
    },
    {
      "id": "uuid",
      "resource": "menu",
      "action": "emails",
      "description": "Emails menu item"
    }
  ],
  "total": 15
}
```

#### Menu Permission Matrix (Alle Rollen × Alle Menu Items)
```http
GET /api/rbac/menu/matrix
Authorization: Bearer <token>
```

**Response:**
```json
{
  "roles": [...],
  "permissions": [...],
  "matrix": {
    "role-uuid-1": {
      "perm-uuid-dashboard": true,
      "perm-uuid-emails": false
    }
  }
}
```

**Gebruik:** Perfect voor een checkbox grid UI.

#### Menu Permissions voor Specifieke Rol
```http
GET /api/rbac/roles/:roleId/menu-permissions
Authorization: Bearer <token>
```

#### Menu Permissions Bijwerken
```http
PUT /api/rbac/roles/:roleId/menu-permissions
Authorization: Bearer <token>
Content-Type: application/json

{
  "permission_actions": ["dashboard", "emails", "users"]
}
```

**Note:** Gebruikt action names in plaats van UUIDs voor frontend gemak.

---

### User-Role Management

#### Lijst Gebruikers met Rollen
```http
GET /api/rbac/users?limit=50&offset=0&search=admin
Authorization: Bearer <token>
```

**Response:**
```json
{
  "users": [
    {
      "id": "user-uuid",
      "email": "admin@example.com",
      "naam": "Admin User",
      "roles": [
        {
          "id": "role-uuid",
          "name": "admin",
          "description": "Full system access",
          "assigned_at": "2025-01-01T10:00:00Z",
          "expires_at": null,
          "is_active": true
        }
      ]
    }
  ],
  "total": 1
}
```

#### Rol Toewijzen aan Gebruiker
```http
POST /api/rbac/users/:userId/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "role_id": "role-uuid",
  "expires_at": "2026-01-01T00:00:00Z"  // Optional
}
```

#### Rol Intrekken van Gebruiker
```http
DELETE /api/rbac/users/:userId/roles/:roleId
Authorization: Bearer <token>
```

#### User-Role Status Bijwerken (Activate/Deactivate)
```http
PUT /api/rbac/users/:userId/roles/:roleId
Authorization: Bearer <token>

{
  "is_active": false,
  "expires_at": "2026-06-01T00:00:00Z"  // Optional
}
```

---

## 📊 Session Management API

### Lijst Actieve Sessies
```http
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
        "device_type": "desktop"
      },
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "login_time": "2025-11-18T10:30:00Z",
      "last_activity": "2025-11-18T10:45:00Z",
      "is_current": false,
      "display_name": "Chrome on Windows"
    }
  ]
}
```

### Specifieke Sessie Intrekken
```http
DELETE /api/auth/sessions/:sessionId
Authorization: Bearer <token>
```

### Alle Andere Sessies Intrekken
```http
POST /api/auth/sessions/revoke-others
Authorization: Bearer <token>
```

---

## 🔧 TypeScript Integration Voorbeelden

### RBAC Service
```typescript
// src/services/rbac-management.service.ts
import { apiClient } from '@/api/axios-client';

export const rbacManagementService = {
  // PERMISSIONS
  async listPermissions(groupByResource = true) {
    const { data } = await apiClient.get('/rbac/permissions', {
      params: { group_by_resource: groupByResource, limit: 1000 }
    });
    return data;
  },

  async createPermission(permission: {
    resource: string;
    action: string;
    description: string;
  }) {
    const { data } = await apiClient.post('/rbac/permissions', permission);
    return data;
  },

  // ROLES
  async listRoles() {
    const { data } = await apiClient.get('/rbac/roles', {
      params: { limit: 100, offset: 0 }
    });
    return data;
  },

  async createRole(role: { name: string; description: string }) {
    const { data } = await apiClient.post('/rbac/roles', role);
    return data;
  },

  // ROLE-PERMISSIONS
  async updateRolePermissions(roleId: string, permissionIds: string[]) {
    const { data } = await apiClient.put(`/rbac/roles/${roleId}/permissions`, {
      permission_ids: permissionIds
    });
    return data;
  },

  // MENU PERMISSIONS
  async getMenuPermissionMatrix() {
    const { data } = await apiClient.get('/rbac/menu/matrix');
    return data;
  },

  async updateRoleMenuPermissions(roleId: string, actions: string[]) {
    const { data } = await apiClient.put(`/rbac/roles/${roleId}/menu-permissions`, {
      permission_actions: actions
    });
    return data;
  },

  // USER-ROLES
  async listUsersWithRoles(limit = 50, offset = 0, search = '') {
    const { data } = await apiClient.get('/rbac/users', {
      params: { limit, offset, search }
    });
    return data;
  },

  async assignRoleToUser(userId: string, roleId: string, expiresAt?: string) {
    const { data } = await apiClient.post(`/rbac/users/${userId}/roles`, {
      role_id: roleId,
      expires_at: expiresAt
    });
    return data;
  },

  async revokeRoleFromUser(userId: string, roleId: string) {
    const { data } = await apiClient.delete(`/rbac/users/${userId}/roles/${roleId}`);
    return data;
  }
};
```

### React Hook: usePermissions
```typescript
// src/hooks/use-permissions.ts
import { useQuery } from '@tanstack/react-query';
import { rbacManagementService } from '@/services/rbac-management.service';

export function usePermissions(groupByResource = true) {
  return useQuery({
    queryKey: ['permissions', groupByResource],
    queryFn: () => rbacManagementService.listPermissions(groupByResource),
    staleTime: 5 * 60 * 1000, // 5 minuten (match backend cache)
  });
}

export function useRoles() {
  return useQuery({
    queryKey: ['roles'],
    queryFn: () => rbacManagementService.listRoles(),
    staleTime: 5 * 60 * 1000,
  });
}

export function useUsersWithRoles(params: {
  limit?: number;
  offset?: number;
  search?: string;
}) {
  return useQuery({
    queryKey: ['users-with-roles', params],
    queryFn: () => rbacManagementService.listUsersWithRoles(
      params.limit,
      params.offset,
      params.search
    ),
    staleTime: 2 * 60 * 1000, // 2 minuten
  });
}
```

### Permission Guard Component
```typescript
// src/components/auth/permission-guard.tsx
import { useAuth } from '@/context/auth-context';

interface PermissionGuardProps {
  resource: string;
  action: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function PermissionGuard({
  resource,
  action,
  children,
  fallback = null
}: PermissionGuardProps) {
  const { hasPermission } = useAuth();

  if (!hasPermission(resource, action)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}

// Usage:
// <PermissionGuard resource="user" action="write">
//   <CreateUserButton />
// </PermissionGuard>
```

---

## 🎨 Permission Structuur

### Naming Convention

Format: `<resource>:<action>`

### Admin Panel Permissions

#### Menu Access (Sidebar Visibility)
```
menu:dashboard      - Dashboard toegang
menu:emails         - Emails menu item
menu:contacts       - Contacten menu item  
menu:participants   - Deelnemers menu item
menu:users          - Gebruikers menu item
menu:roles          - Rollen & Permissions menu item
menu:events         - Evenementen menu item
menu:cms            - CMS menu item
menu:newsletter     - Newsletter menu item
menu:settings       - Instellingen menu item
```

#### Resource Permissions (CRUD Operations)
```
admin:access        - Volledige admin toegang (superuser)
user:read           - Gebruikers bekijken
user:write          - Gebruikers aanmaken/bewerken
user:delete         - Gebruikers verwijderen
user:manage_roles   - Gebruikersrollen beheren
contact:read        - Contactformulieren bekijken
contact:write       - Contactformulieren bewerken
contact:delete      - Contactformulieren verwijderen
email:read          - Emails bekijken
email:write         - Emails markeren als verwerkt
email:delete        - Emails verwijderen
email:fetch         - Nieuwe emails ophalen
newsletter:read     - Nieuwsbrieven bekijken
newsletter:write    - Nieuwsbrieven aanmaken/bewerken
newsletter:send     - Nieuwsbrieven verzenden
newsletter:delete   - Nieuwsbrieven verwijderen
```

### Standaard Rollen

#### Admin
- Alle permissions (127+ permissions)
- Volledig CRUD op alle resources
- User & role management
- Kan niet worden verwijderd (is_system_role: true)

#### Staff
- Alleen read permissions op resources
- `staff:access` permission
- Kan niet worden verwijderd (is_system_role: true)

#### User
- `chat:read`, `chat:write`
- Eigen profiel beheer
- Kan niet worden verwijderd (is_system_role: true)

---

## 🔍 Permission Checking

### Backend Middleware

**AuthMiddleware** ([`handlers/middleware.go`](../../handlers/middleware.go)):
1. Valideert JWT token
2. Haalt gebruiker op uit database
3. Zet `userID` in Fiber context

**PermissionMiddleware** ([`handlers/permission_middleware.go`](../../handlers/permission_middleware.go)):
1. Haalt `userID` uit context
2. Check Redis cache: `perm:{user_id}:{resource}:{action}`
3. Fallback naar database bij cache miss
4. Cache TTL: 5 minuten

**AdminOrStaffPermissionMiddleware:**
- Controleert `admin:access` OF `staff:access`
- Gebruikt voor alle `/api/rbac/*` routes

### Frontend Permission Checking

**In Context:**
```typescript
const { hasPermission } = useAuth();

// Instant check (geen API call)
if (hasPermission('user', 'write')) {
  // Toon "Create User" knop
}
```

**In Component:**
```typescript
<PermissionGuard resource="user" action="delete">
  <DeleteUserButton userId={userId} />
</PermissionGuard>
```

**In Route:**
```typescript
<Route
  element={<ProtectedRoute requiredPermission={{ resource: 'admin', action: 'access' }} />}
>
  <Route path="/admin" element={<AdminPanel />} />
</Route>
```

---

## ⚡ Performance Optimizations

### Caching Strategy

**Backend (Redis):**
- Permission checks: 5 minuten cache
- Cache key: `perm:{user_id}:{resource}:{action}`
- Automatic invalidation bij rol/permissie wijzigingen

**Frontend (React Query):**
```typescript
// Cache permissions voor 5 minuten
const { data: permissions } = useQuery({
  queryKey: ['permissions'],
  queryFn: fetchPermissions,
  staleTime: 5 * 60 * 1000, // Match backend cache
});

// Cache role details voor 2 minuten
const { data: roles } = useQuery({
  queryKey: ['roles'],
  queryFn: fetchRoles,
  staleTime: 2 * 60 * 1000,
});
```

### Batch Operations

**Vermijd:**
```typescript
// ❌ BAD: Multiple API calls in loop
for (const permId of permissionIds) {
  await addPermissionToRole(roleId, permId);
}
```

**Gebruik:**
```typescript
// ✅ GOOD: Single bulk update
await updateRolePermissions(roleId, permissionIds);
```

---

## 🚨 Error Handling

### Error Response Format

```json
{
  "error": "Beschrijvende foutmelding",
  "code": "ERROR_CODE"
}
```

### Common Error Codes

| HTTP | Code | Oplossing |
|------|------|-----------|
| 401 | `TOKEN_EXPIRED` | Refresh token via `/api/auth/refresh` |
| 401 | `INVALID_CREDENTIALS` | Check email/wachtwoord |
| 403 | `PERMISSION_DENIED` | User heeft niet de juiste permissie |
| 403 | `USER_INACTIVE` | Account is gedeactiveerd |
| 404 | `NOT_FOUND` | Resource bestaat niet |
| 409 | `DUPLICATE_EMAIL` | Email bestaat al |
| 409 | `ALREADY_ASSIGNED` | Rol is al toegewezen |
| 429 | `RATE_LIMIT_EXCEEDED` | Te veel requests, wacht even |

### Frontend Error Handling

```typescript
try {
  await rbacService.assignRoleToUser(userId, roleId);
  toast.success('Rol succesvol toegewezen');
} catch (error) {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.error || 'Er is een fout opgetreden';
    const code = error.response?.data?.code;
    
    // Specifieke foutafhandeling
    switch (code) {
      case 'PERMISSION_DENIED':
        toast.error('Je hebt geen rechten voor deze actie');
        break;
      case 'TOKEN_EXPIRED':
        // Auto-handled door axios interceptor
        break;
      default:
        toast.error(message);
    }
  }
}
```

---

## 🔒 Security Best Practices

### 1. Token Management

**✅ DO:**
- Access token in memory (React state/context)
- Refresh token in localStorage of httpOnly cookie
- Clear tokens on logout
- Automatic token refresh voor 401 responses

**❌ DON'T:**
- Access token in localStorage (XSS risk)
- Tokens in URL parameters (leaks via logs/history)
- Tokens in sessionStorage (lost on tab close)
- Hardcode tokens in code

### 2. Permission Checks

**✅ DO:**
- Check permissions before rendering UI elements
- Validate permissions server-side (backend middleware)
- Cache permissions in context na login
- Use granular permissions (`user:write` not `admin`)

**❌ DON'T:**
- Only check permissions client-side (security hole)
- Assume admin = all permissions (use explicit checks)
- Skip permission checks voor "small" features

### 3. Rate Limiting

**Backend Limits:**
- Login: 5 pogingen / 15 minuten (per email)
- Forgot password: 3 requests / uur (per email)
- Contact form: 5 submissions / uur (per IP)

**Frontend Handling:**
```typescript
if (error.response?.status === 429) {
  toast.error('Te veel verzoeken, probeer het over een paar minuten opnieuw');
  
  // Disable submit button for X seconds
  setIsSubmitting(true);
  setTimeout(() => setIsSubmitting(false), 60000);
}
```

---

## 🧪 Testing

### Manual API Testing

**1. Test Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dekoninklijkeloop.nl",
    "wachtwoord": "YourPassword"
  }'
```

**2. Test Protected Endpoint:**
```bash
TOKEN="eyJhbGci..."

curl http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

**3. Test CORS:**
```bash
curl -H "Origin: http://localhost:5174" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type,Authorization" \
     -X OPTIONS \
     http://localhost:8080/api/auth/login -v
```

**Expected Response:**
```
< Access-Control-Allow-Origin: http://localhost:5174
< Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
< Access-Control-Allow-Credentials: true
```

### Frontend Integration Testing

**Browser DevTools:**
1. Open http://localhost:5174
2. Open DevTools → Network tab
3. Login met admin credentials
4. Check request headers:
   - `Origin: http://localhost:5174`
   - `Authorization: Bearer ...`
5. Check response headers:
   - `Access-Control-Allow-Origin: http://localhost:5174`

**Expected Behavior:**
- ✅ No CORS errors
- ✅ 200 OK response
- ✅ User object in response
- ✅ Token set in memory

---

## 📝 Database Schema Referentie

### RBAC Tabellen

**roles:**
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**permissions:**
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    is_system_permission BOOLEAN DEFAULT FALSE,
    UNIQUE(resource, action)
);
```

**role_permissions:**
```sql
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);
```

**user_roles:**
```sql
CREATE TABLE user_roles (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES gebruikers(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(user_id, role_id)
);
```

### Belangrijke Views

**user_permissions view:**
```sql
CREATE OR REPLACE VIEW user_permissions AS
SELECT
    ur.user_id,
    u.email,
    r.name as role_name,
    p.resource,
    p.action
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
JOIN gebruikers u ON ur.user_id = u.id
WHERE ur.is_active = true;
```

---

## 🐛 Troubleshooting

### CORS Errors

**Symptoom:**
```
Access to XMLHttpRequest has been blocked by CORS policy
```

**Diagnose:**
```bash
# Check CORS configuratie in logs
docker logs dkl-email-service 2>&1 | findstr "CORS"

# Verwacht: "origins":["...","http://localhost:5174",...]
```

**Oplossing:**
1. Controleer `.env.local` bevat `http://localhost:5174`
2. Restart Docker: `docker-compose restart app`
3. Clear browser cache
4. Hard refresh (Ctrl+Shift+R)

### 401 Unauthorized Errors

**Symptoom:**
```json
{
  "error": "Token expired",
  "code": "TOKEN_EXPIRED"
}
```

**Diagnose:**
```typescript
// Check token in DevTools → Application → Local Storage
console.log('Refresh token:', localStorage.getItem('refresh_token'));
```

**Oplossing:**
```typescript
// Forceer token refresh
const refreshToken = localStorage.getItem('refresh_token');
const response = await fetch('/api/auth/refresh', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ refresh_token: refreshToken })
});
```

### Permission Denied Errors

**Symptoom:**
```json
{
  "error": "Permission denied",
  "code": "PERMISSION_DENIED"
}
```

**Diagnose:**
```typescript
// Check user permissions
const profile = await fetch('/api/auth/profile', {
  headers: { 'Authorization': `Bearer ${token}` }
});
console.log('User permissions:', profile.permissions);
```

**Oplossing:**
1. Check of gebruiker de juiste rol heeft
2. Check of rol de juiste permissions heeft
3. Clear browser cache + logout/login
4. Contact admin om rol toe te wijzen

---

## 🚀 Quick Start Guide

### Voor Developers

**1. Clone repository:**
```bash
git clone https://github.com/dekoninklijkeloop/dklemailservice.git
cd dklemailservice
```

**2. Copy environment file:**
```bash
cp .env.local .env
# OF voor Windows:
copy .env.local .env
```

**3. Start Docker services:**
```bash
docker-compose up -d
```

**4. Verify backend:**
```bash
curl http://localhost:8080/api/health
```

**Expected:**
```json
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected"
}
```

**5. Test login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"admin123"}'
```

**6. Start admin panel (separate project):**
```bash
cd ../dkl-admin-panel
npm install
npm run dev
# Opens on http://localhost:5174
```

**7. Login in browser:**
- Navigate to http://localhost:5174
- Login met admin credentials
- Geen CORS errors! ✅

---

## 📚 Gerelateerde Documentatie

- [AUTHENTICATION.md](../../api/archive/AUTHENTICATION.md) - Volledige auth documentatie
- [PERMISSIONS.md](../../api/archive/PERMISSIONS.md) - RBAC API details
- [Admin Panel README](./README.md) - Frontend admin panel docs
- [API_DOCUMENTATION.md](../../api/API_DOCUMENTATION.md) - Complete API overzicht

---

## 🆘 Support

**Bij problemen:**
1. Check Docker logs: `docker logs dkl-email-service --tail 100`
2. Check backend health: `curl http://localhost:8080/api/health`
3. Check CORS logs: `docker logs dkl-email-service 2>&1 | findstr "CORS"`
4. Clear browser cache en cookies
5. Logout/login opnieuw

**Contact:**
- GitHub Issues: [DKL Email Service](https://github.com/dekoninklijkeloop/dklemailservice/issues)
- Email: info@dekoninklijkeloop.nl
