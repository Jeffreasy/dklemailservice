# Permissions & RBAC API

API endpoints voor Role-Based Access Control (RBAC) beheer.

## Base URL

```
Development: http://localhost:8080/api
Production: https://api.dklemailservice.com/api
```

## Authentication

Alle endpoints vereisen JWT authenticatie en admin permissies:

```
Authorization: Bearer <your_jwt_token>
```

---

## Permissions

### List All Permissions

**Endpoint:** `GET /api/permissions`

**Authentication:** Vereist - `permission:read` permissie

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "contact:read",
    "description": "Kan contactformulieren bekijken",
    "resource": "contact",
    "action": "read",
    "created_at": "2025-01-08T10:00:00Z"
  },
  {
    "id": "uuid",
    "name": "user:write",
    "description": "Kan gebruikers aanmaken en bewerken",
    "resource": "user",
    "action": "write",
    "created_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### Get Permission by ID

**Endpoint:** `GET /api/permissions/:id`

**Authentication:** Vereist - `permission:read` permissie

---

### Create Permission

**Endpoint:** `POST /api/permissions`

**Authentication:** Vereist - `permission:write` permissie

**Request Body:**

```json
{
  "name": "custom_resource:custom_action",
  "description": "Custom permission description",
  "resource": "custom_resource",
  "action": "custom_action"
}
```

**Naming Convention:**
- Format: `resource:action`
- Resources: `contact`, `user`, `event`, `chat`, etc.
- Actions: `read`, `write`, `delete`, `manage`

**Response:** `201 Created`

---

### Update Permission

**Endpoint:** `PUT /api/permissions/:id`

**Authentication:** Vereist - `permission:write` permissie

**Request Body:**

```json
{
  "description": "Updated description"
}
```

**Note:** `name`, `resource`, en `action` kunnen niet worden gewijzigd om consistentie te behouden.

---

### Delete Permission

**Endpoint:** `DELETE /api/permissions/:id`

**Authentication:** Vereist - `permission:delete` permissie

**Note:** Cascade delete verwijdert ook alle role-permission koppelingen.

---

## Roles

### List All Roles

**Endpoint:** `GET /api/roles`

**Authentication:** Vereist - `role:read` permissie

**Response:** `200 OK`

```json
[
  {
    "id": "uuid",
    "name": "admin",
    "description": "Volledige systeem toegang",
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  },
  {
    "id": "uuid",
    "name": "moderator",
    "description": "Content moderatie toegang",
    "created_at": "2025-01-08T10:00:00Z",
    "updated_at": "2025-01-08T10:00:00Z"
  }
]
```

---

### Get Role by ID

**Endpoint:** `GET /api/roles/:id`

**Authentication:** Vereist - `role:read` permissie

**Response:** `200 OK`

```json
{
  "id": "uuid",
  "name": "admin",
  "description": "Volledige systeem toegang",
  "permissions": [
    {
      "id": "uuid",
      "name": "contact:read",
      "description": "Kan contactformulieren bekijken"
    },
    {
      "id": "uuid",
      "name": "contact:write",
      "description": "Kan contactformulieren bewerken"
    }
  ],
  "created_at": "2025-01-08T10:00:00Z",
  "updated_at": "2025-01-08T10:00:00Z"
}
```

---

### Create Role

**Endpoint:** `POST /api/roles`

**Authentication:** Vereist - `role:write` permissie

**Request Body:**

```json
{
  "name": "content_editor",
  "description": "Kan content beheren maar geen gebruikersrechten wijzigen"
}
```

**Validation:**
- `name`: Verplicht, uniek, lowercase, geen spaties
- `description`: Verplicht

**Response:** `201 Created`

---

### Update Role

**Endpoint:** `PUT /api/roles/:id`

**Authentication:** Vereist - `role:write` permissie

**Request Body:**

```json
{
  "description": "Updated role description"
}
```

**Note:** `name` kan niet worden gewijzigd om referentie integriteit te behouden.

---

### Delete Role

**Endpoint:** `DELETE /api/roles/:id`

**Authentication:** Vereist - `role:delete` permissie

**Note:** Kan niet default roles verwijderen (`admin`, `user`, `moderator`, `guest`).

---

## Role Permissions

### Get Role Permissions

**Endpoint:** `GET /api/roles/:id/permissions`

**Authentication:** Vereist - `role:read` permissie

**Response:** `200 OK`

```json
{
  "role_id": "uuid",
  "role_name": "moderator",
  "permissions": [
    {
      "id": "uuid",
      "name": "chat:moderate",
      "description": "Kan chat berichten modereren",
      "resource": "chat",
      "action": "moderate"
    }
  ]
}
```

---

### Assign Permissions to Role

**Endpoint:** `POST /api/roles/:id/permissions`

**Authentication:** Vereist - `role:write` permissie

**Request Body:**

```json
{
  "permission_ids": [
    "uuid1",
    "uuid2",
    "uuid3"
  ]
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Permissions assigned successfully",
  "role_id": "uuid",
  "permissions_added": 3
}
```

---

### Update Role Permissions (Batch)

Vervang alle permissies van een rol in één keer.

**Endpoint:** `PUT /api/roles/:id/permissions`

**Authentication:** Vereist - `role:write` permissie

**Request Body:**

```json
{
  "permission_ids": [
    "uuid1",
    "uuid2",
    "uuid3"
  ]
}
```

**Note:** Verwijdert eerst alle bestaande permissies en voegt dan de nieuwe toe.

---

### Remove Permission from Role

**Endpoint:** `DELETE /api/roles/:id/permissions/:permission_id`

**Authentication:** Vereist - `role:write` permissie

---

## User Roles

### Get User Roles

**Endpoint:** `GET /api/users/:id/roles`

**Authentication:** Vereist - `user:read` permissie

**Response:** `200 OK`

```json
{
  "user_id": "uuid",
  "user_email": "user@example.com",
  "roles": [
    {
      "id": "uuid",
      "name": "user",
      "description": "Standaard gebruiker toegang",
      "assigned_at": "2025-01-01T10:00:00Z",
      "assigned_by": "uuid"
    },
    {
      "id": "uuid",
      "name": "moderator",
      "description": "Content moderatie toegang",
      "assigned_at": "2025-01-05T14:30:00Z",
      "assigned_by": "uuid"
    }
  ]
}
```

---

### Assign Role to User

**Endpoint:** `POST /api/users/:id/roles`

**Authentication:** Vereist - `user:manage_roles` permissie

**Request Body:**

```json
{
  "role_id": "uuid"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Role assigned successfully",
  "user_id": "uuid",
  "role_id": "uuid",
  "role_name": "moderator"
}
```

---

### Assign Multiple Roles to User

**Endpoint:** `POST /api/users/:id/roles/batch`

**Authentication:** Vereist - `user:manage_roles` permissie

**Request Body:**

```json
{
  "role_ids": [
    "uuid1",
    "uuid2"
  ]
}
```

---

### Remove Role from User

**Endpoint:** `DELETE /api/users/:id/roles/:role_id`

**Authentication:** Vereist - `user:manage_roles` permissie

**Note:** Laatste admin rol kan niet worden verwijderd om lock-out te voorkomen.

---

## Permission Checking

### Check User Permission

Controleer of een gebruiker een specifieke permissie heeft.

**Endpoint:** `GET /api/permissions/check`

**Authentication:** Vereist

**Query Parameters:**
- `user_id`: User ID (optional, default: current user)
- `resource`: Resource naam (required)
- `action`: Action naam (required)

**Response:** `200 OK`

```json
{
  "has_permission": true,
  "user_id": "uuid",
  "permission": "contact:write",
  "granted_via": ["admin", "moderator"]
}
```

---

## Default Roles & Permissions

### Admin Role

Full system access:

```
- contact:read, contact:write, contact:delete
- user:read, user:write, user:delete, user:manage_roles
- event:read, event:write, event:delete
- chat:read, chat:write, chat:moderate, chat:manage_channel
- newsletter:read, newsletter:write, newsletter:send, newsletter:delete
- email:read, email:write, email:delete, email:fetch
- permission:read, permission:write, permission:delete
- role:read, role:write, role:delete
- All CMS permissions (album, photo, video, partner, sponsor, etc.)
- All gamification permissions (achievement, badge, leaderboard)
```

---

### Moderator Role

Content moderation access:

```
- contact:read, contact:write
- chat:read, chat:write, chat:moderate
- email:read
- album:read, album:write
- photo:read, photo:write
- video:read, video:write
```

---

### User Role

Standard user access:

```
- chat:read, chat:write
- Own profile management
```

---

### Guest Role

Limited read-only access:

```
- Public content viewing only
- No authenticated endpoints
```

---

## Permission Caching

### Redis Cache Strategy

Permissions worden gecached voor optimale performance:

- **Cache Key**: `permissions:user:{user_id}:resource:{resource}:action:{action}`
- **TTL**: 5 minuten
- **Invalidation**: Bij rol/permissie wijzigingen

**Benefits:**
- Sub-milliseconde permission checks
- Verminderde database load
- Horizontaal schaalbaar

**Fallback:**
- Bij Redis failure: Direct database lookup
- Geen functionaliteit verlies
- Performance degradatie zonder cache

---

## Audit Logging

Alle RBAC wijzigingen worden gelogd:

```json
{
  "timestamp": "2025-01-08T10:00:00Z",
  "action": "role_assigned",
  "user_id": "uuid",
  "target_user_id": "uuid",
  "role_id": "uuid",
  "role_name": "moderator",
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0..."
}
```

**Logged Actions:**
- Role assignments/removals
- Permission changes
- Role creation/deletion
- Permission creation/deletion

---

## Best Practices

### Permission Design

1. **Granularity**: Gebruik specifieke permissies (`contact:read` i.p.v. `admin`)
2. **Naming**: Volg `resource:action` conventie consistent
3. **Documentation**: Documenteer elke permissie duidelijk
4. **Least Privilege**: Ken minimale vereiste permissies toe

### Role Management

1. **Role Hierarchy**: Creëer logische rol hiërarchie
2. **Default Roles**: Wijzig niet de standaard rollen
3. **Custom Roles**: Maak specifieke rollen voor use cases
4. **Regular Review**: Controleer regelmatig rol toewijzingen

### Security Considerations

1. **Admin Protection**: Bescherm tegen onbedoelde admin verwijdering
2. **Audit Trail**: Log alle wijzigingen voor compliance
3. **Cache Invalidation**: Invalideer cache bij wijzigingen
4. **Token Refresh**: Forceer token refresh na rol wijzigingen

---

## Frontend Integration

### Permission Check Hook (React)

```typescript
import { useState, useEffect } from 'react';
import { useAuth } from './useAuth';

export function usePermission(resource: string, action: string) {
  const { user } = useAuth();
  const [hasPermission, setHasPermission] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) {
      setHasPermission(false);
      setLoading(false);
      return;
    }

    // Check permission from user object
    const permissionKey = `${resource}:${action}`;
    const has = user.permissions?.includes(permissionKey) || false;
    
    setHasPermission(has);
    setLoading(false);
  }, [user, resource, action]);

  return { hasPermission, loading };
}

// Usage
function ContactList() {
  const { hasPermission } = usePermission('contact', 'write');

  return (
    <div>
      {hasPermission && (
        <button onClick={handleCreate}>Create Contact</button>
      )}
    </div>
  );
}
```

---

### Role-Based Component (Vue)

```vue
<template>
  <div v-if="hasAnyRole(requiredRoles)">
    <slot></slot>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useAuth } from '@/composables/useAuth';

const props = defineProps({
  requiredRoles: {
    type: Array,
    required: true
  }
});

const { user } = useAuth();

const hasAnyRole = computed(() => {
  if (!user.value) return false;
  return props.requiredRoles.some(role => 
    user.value.roles?.includes(role)
  );
});
</script>

<!-- Usage -->
<RoleGuard :required-roles="['admin', 'moderator']">
  <AdminPanel />
</RoleGuard>
```

---

## Permission Middleware

De applicatie gebruikt middleware voor automatische permissie controle:

```go
// Handler voorbeeld
admin := app.Group("/api/admin", 
  AuthMiddleware(authService),
  PermissionMiddleware(permissionService, "admin", "read")
)
```

**Middleware Flow:**

```
Request → AuthMiddleware → PermissionMiddleware → Handler
            ↓                    ↓
         JWT Check         Permission Check
         User Load         Redis Cache Lookup
                          Database Fallback
```

---

## Resource Types

### Available Resources

| Resource | Description |
|----------|-------------|
| `contact` | Contactformulieren |
| `participant` | Deelnemers |
| `registration` | Event registraties |
| `user` | Gebruikers |
| `role` | Rollen |
| `permission` | Permissies |
| `event` | Evenementen |
| `chat` | Chat systeem |
| `newsletter` | Nieuwsbrieven |
| `email` | Email beheer |
| `album` | Foto albums |
| `photo` | Foto's |
| `video` | Video's |
| `partner` | Partners |
| `sponsor` | Sponsors |
| `achievement` | Achievements |
| `badge` | Badges |
| `leaderboard` | Leaderboards |
| `notulen` | Notulen |
| `steps` | Stappen tracking |

---

## Action Types

### Standard Actions

| Action | Description |
|--------|-------------|
| `read` | Mag resource bekijken/ophalen |
| `write` | Mag resource aanmaken/updaten |
| `delete` | Mag resource verwijderen |
| `manage` | Volledige controle over resource |

### Special Actions

| Action | Resource | Description |
|--------|----------|-------------|
| `send` | `newsletter` | Mag nieuwsbrieven verzenden |
| `fetch` | `email` | Mag emails ophalen |
| `moderate` | `chat` | Mag chat modereren |
| `manage_roles` | `user` | Mag gebruikersrollen beheren |
| `manage_channel` | `chat` | Mag chat kanalen beheren |

---

## Error Codes

| Code | HTTP Status | Beschrijving |
|------|-------------|--------------|
| `PERMISSION_DENIED` | 403 | Onvoldoende permissies |
| `ROLE_NOT_FOUND` | 404 | Rol niet gevonden |
| `PERMISSION_NOT_FOUND` | 404 | Permissie niet gevonden |
| `DUPLICATE_ROLE` | 409 | Rol bestaat al |
| `DUPLICATE_PERMISSION` | 409 | Permissie bestaat al |
| `INVALID_PERMISSION_FORMAT` | 400 | Ongeldige permissie naam |
| `CANNOT_DELETE_DEFAULT_ROLE` | 409 | Kan standaard rol niet verwijderen |
| `LAST_ADMIN_ROLE` | 409 | Kan laatste admin rol niet verwijderen |

---

## Testing

### Manual Permission Testing

```bash
# List all permissions
curl http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Create custom permission
curl -X POST http://localhost:8080/api/permissions \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "custom:action",
    "description": "Custom permission",
    "resource": "custom",
    "action": "action"
  }'

# Assign role to user
curl -X POST http://localhost:8080/api/users/USER_ID/roles \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": "ROLE_ID"
  }'

# Check permission
curl "http://localhost:8080/api/permissions/check?resource=contact&action=write" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Migration Guide

### Adding New Permissions

1. **Define Permission:**
   ```sql
   INSERT INTO permissions (name, description, resource, action)
   VALUES ('new_resource:read', 'Can view new resource', 'new_resource', 'read');
   ```

2. **Assign to Roles:**
   ```sql
   INSERT INTO role_permissions (role_id, permission_id)
   SELECT r.id, p.id
   FROM roles r, permissions p
   WHERE r.name = 'admin' AND p.name = 'new_resource:read';
   ```

3. **Update Middleware:**
   ```go
   writeGroup := admin.Group("", 
     PermissionMiddleware(permissionService, "new_resource", "write")
   )
   ```

4. **Clear Cache:**
   ```bash
   redis-cli FLUSHDB  # Development only
   # Or selective: redis-cli DEL permissions:user:*
   ```

---

## Performance Optimization

### Reducing Permission Checks

```typescript
// Bad: Check permission on every render
function Component() {
  const { hasPermission } = usePermission('contact', 'write');
  // Re-checks on every render
}

// Good: Check permission once on mount
function Component() {
  const [canWrite, setCanWrite] = useState(false);
  
  useEffect(() => {
    checkPermission('contact', 'write').then(setCanWrite);
  }, []);
}

// Better: Use context/state to share permission data
const { user } = useAuth();
const canWrite = user.permissions?.includes('contact:write');
```

---

### Batch Permission Loading

```typescript
// Load all user permissions at login
async function loadUserWithPermissions() {
  const response = await api.get('/api/auth/me');
  const user = response.data.data;
  
  // User object includes all permissions
  // Store in global state (Redux/Zustand/Context)
  setUser(user);
  
  // Now permission checks are instant lookups
}
```

---

Voor meer informatie:
- [Authentication API](./AUTHENTICATION.md)
- [User Management](./USERS.md)
- [API Overview](./README.md)