# User Roles API - Frontend Integratie

## Endpoint: GET /api/users/:id/roles

Deze endpoint wordt gebruikt om alle rollen van een gebruiker op te halen.

### Route Registratie
De route wordt geregistreerd in [`handlers/user_handler.go:40`](../../handlers/user_handler.go:40):
```go
app.Get("/api/users/:id/roles", AuthMiddleware(h.authService), PermissionMiddleware(h.permissionService, "user", "read"), h.GetUserRoles)
```

### Vereiste Headers
- `Authorization: Bearer <jwt_token>` (verplicht)
- `Content-Type: application/json`

### Vereiste Permissies
De gebruiker moet de permissie `user:read` hebben.

### URL Parameters
- `:id` - De UUID van de gebruiker (bijvoorbeeld: `0197cfc3-7ca2-403b-ae4d-32627cd47222`)

### Response Format

**Succes (200 OK):**
```json
[
  {
    "id": "uuid",
    "user_id": "0197cfc3-7ca2-403b-ae4d-32627cd47222",
    "role_id": "uuid",
    "role": {
      "id": "uuid",
      "name": "admin",
      "display_name": "Administrator",
      "description": "Full system access"
    },
    "assigned_by": "uuid",
    "assigned_at": "2025-01-15T10:30:00Z",
    "expires_at": null,
    "is_active": true
  }
]
```

**Errors:**
- `401 Unauthorized` - Geen geldig JWT token
- `403 Forbidden` - Gebruiker heeft geen `user:read` permissie
- `500 Internal Server Error` - Database of server fout

### Frontend Implementatie

#### JavaScript/TypeScript
```typescript
async function getUserRoles(userId: string): Promise<UserRole[]> {
  const response = await fetch(`${API_BASE_URL}/api/users/${userId}/roles`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    }
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
}
```

#### React Hook Example
```typescript
import { useState, useEffect } from 'react';

interface UserRole {
  id: string;
  user_id: string;
  role_id: string;
  role: Role;
  assigned_by?: string;
  assigned_at: string;
  expires_at?: string;
  is_active: boolean;
}

interface Role {
  id: string;
  name: string;
  display_name: string;
  description?: string;
}

function useUserRoles(userId: string) {
  const [roles, setRoles] = useState<UserRole[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchUserRoles() {
      try {
        setLoading(true);
        const response = await fetch(`${process.env.REACT_APP_API_URL}/api/users/${userId}/roles`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const data = await response.json();
        setRoles(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
        setRoles([]);
      } finally {
        setLoading(false);
      }
    }

    if (userId) {
      fetchUserRoles();
    }
  }, [userId]);

  return { roles, loading, error };
}

// Usage in component
function UserRolesDisplay({ userId }: { userId: string }) {
  const { roles, loading, error } = useUserRoles(userId);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h3>User Roles</h3>
      <ul>
        {roles.map(userRole => (
          <li key={userRole.id}>
            {userRole.role.display_name} ({userRole.role.name})
            {userRole.expires_at && ` - Expires: ${new Date(userRole.expires_at).toLocaleDateString()}`}
          </li>
        ))}
      </ul>
    </div>
  );
}
```

### Gerelateerde Endpoints

#### Toewijzen van een enkele rol
**POST /api/users/:id/roles**
```typescript
async function assignRole(userId: string, roleId: string, expiresAt?: string) {
  const response = await fetch(`${API_BASE_URL}/api/users/${userId}/roles`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      role_id: roleId,
      expires_at: expiresAt // Optional: RFC3339 format
    })
  });
  
  return await response.json();
}
```
**Vereiste permissie:** `user:manage_roles`

#### Toewijzen van meerdere rollen (Admin only)
**PUT /api/users/:id/roles**
```typescript
async function assignMultipleRoles(userId: string, roleIds: string[]) {
  const response = await fetch(`${API_BASE_URL}/api/users/${userId}/roles`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      role_ids: roleIds
    })
  });
  
  return await response.json();
}
```
**Vereiste permissie:** Admin permissie (via AdminPermissionMiddleware)

#### Verwijderen van een rol
**DELETE /api/users/:id/roles/:roleId**
```typescript
async function removeRole(userId: string, roleId: string) {
  const response = await fetch(`${API_BASE_URL}/api/users/${userId}/roles/${roleId}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    }
  });
  
  return await response.json();
}
```
**Vereiste permissie:** `user:manage_roles`

#### Ophalen van gebruiker permissies
**GET /api/users/:id/permissions**
```typescript
async function getUserPermissions(userId: string) {
  const response = await fetch(`${API_BASE_URL}/api/users/${userId}/permissions`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    }
  });
  
  return await response.json();
}
```
**Vereiste permissie:** `user:read`

## Troubleshooting

### "Method Not Allowed" Error
Als je een "Method Not Allowed" fout krijgt:

1. **Controleer de HTTP methode:** Zorg ervoor dat je `GET` gebruikt voor ophalen van rollen
2. **Controleer de URL:** De URL moet exact `/api/users/{uuid}/roles` zijn
3. **Controleer Headers:** Zorg voor correcte `Authorization` en `Content-Type` headers
4. **Controleer CORS:** Als je vanaf een andere domeinnaam aanroept, controleer CORS instellingen

### Debug Checklist
```typescript
// Log de volledige request voor debugging
console.log('Request URL:', `${API_BASE_URL}/api/users/${userId}/roles`);
console.log('Method:', 'GET');
console.log('Headers:', {
  'Authorization': `Bearer ${token.substring(0, 20)}...`,
  'Content-Type': 'application/json'
});

try {
  const response = await fetch(url, options);
  console.log('Response status:', response.status);
  console.log('Response headers:', response.headers);
  const data = await response.json();
  console.log('Response data:', data);
} catch (error) {
  console.error('Fetch error:', error);
}
```

### Veelvoorkomende Fouten

1. **405 Method Not Allowed:**
   - Verkeerde HTTP methode gebruikt (bijv. POST in plaats van GET)
   - URL typo of extra/missende slashes

2. **401 Unauthorized:**
   - JWT token ontbreekt of is verlopen
   - Token niet in Authorization header

3. **403 Forbidden:**
   - Gebruiker heeft geen `user:read` permissie
   - Token is geldig maar gebruiker heeft onvoldoende rechten

4. **404 Not Found:**
   - Gebruiker ID bestaat niet in database
   - URL pad is incorrect

5. **500 Internal Server Error:**
   - Database connectie problemen
   - Server-side fout in handler logic

## Backend Implementatie Details

De handler [`GetUserRoles`](../../handlers/user_handler.go:154) in `handlers/user_handler.go`:

```go
func (h *UserHandler) GetUserRoles(c *fiber.Ctx) error {
	userID := c.Params("id")

	userRoles, err := h.userRoleRepo.GetByUserIDWithRoles(c.Context(), userID)
	if err != nil {
		logger.Error("Fout bij ophalen user roles", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kon user roles niet ophalen",
		})
	}

	return c.JSON(userRoles)
}
```

### Repository Method
De data wordt opgehaald via [`repository.UserRoleRepository.GetByUserIDWithRoles()`](../../repository/user_role_repository.go) die:
- Alle actieve user-role relations ophaalt voor de gegeven user ID
- De volledige role informatie mee-laadt via preloading
- Alleen actieve (`is_active = true`) relaties retourneert