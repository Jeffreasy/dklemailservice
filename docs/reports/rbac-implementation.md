# RBAC (Role-Based Access Control) Frontend Documentatie

## Overzicht

De DKL Email Service gebruikt een geavanceerd Role-Based Access Control (RBAC) systeem geïmplementeerd met Redis caching voor optimale prestaties. Dit document beschrijft hoe frontend ontwikkelaars moeten omgaan met authenticatie, autorisatie, permissies en de Redis-gebaseerde features.

## 🔐 Authenticatie & Autorisatie

### JWT Token Gebruik

Alle beveiligde endpoints vereisen JWT token authenticatie:

```javascript
// Headers voor API calls
const headers = {
  'Authorization': `Bearer ${jwtToken}`,
  'Content-Type': 'application/json'
};
```

### Token Vernieuwing

```javascript
// Controleer token expiratie
const isTokenExpired = (token) => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.exp * 1000 < Date.now();
  } catch {
    return true;
  }
};

// Automatische token refresh
const refreshToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');
  const response = await fetch('/api/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken })
  });

  if (response.ok) {
    const data = await response.json();
    localStorage.setItem('jwtToken', data.token);
    return data.token;
  }

  // Redirect naar login als refresh faalt
  window.location.href = '/login';
  throw new Error('Token refresh failed');
};
```

## 👥 Rollen & Permissies

### Systeemrollen

| Rol | Beschrijving | Permissies |
|-----|-------------|------------|
| `admin` | Volledige beheerder | Alle permissies |
| `staff` | Ondersteunend personeel | Beperkte beheerrechten |
| `user` | Standaard gebruiker | Basis chat permissies |
| `owner` | Chat kanaal eigenaar | Volledige chat beheer |
| `chat_admin` | Chat beheerder | Chat moderatie |
| `member` | Chat lid | Basis chat toegang |
| `deelnemer` | Evenement deelnemer | Geen speciale permissies |
| `begeleider` | Evenement begeleider | Geen speciale permissies |
| `vrijwilliger` | Evenement vrijwilliger | Geen speciale permissies |

### Permissie Structuur

Permissies volgen het patroon: `resource:action`

#### Beschikbare Resources & Actions

**Contact Management:**
- `contact:read` - Contactformulieren bekijken
- `contact:write` - Contactformulieren bewerken
- `contact:delete` - Contactformulieren verwijderen

**Aanmeldingen:**
- `aanmelding:read` - Aanmeldingen bekijken
- `aanmelding:write` - Aanmeldingen bewerken
- `aanmelding:delete` - Aanmeldingen verwijderen

**Nieuwsbrieven:**
- `newsletter:read` - Nieuwsbrieven bekijken
- `newsletter:write` - Nieuwsbrieven aanmaken/bewerken
- `newsletter:send` - Nieuwsbrieven verzenden
- `newsletter:delete` - Nieuwsbrieven verwijderen

**Email Management:**
- `email:read` - Inkomende emails bekijken
- `email:write` - Emails bewerken
- `email:delete` - Emails verwijderen
- `email:fetch` - Nieuwe emails ophalen

**Gebruikersbeheer:**
- `user:read` - Gebruikers bekijken
- `user:write` - Gebruikers aanmaken/bewerken
- `user:delete` - Gebruikers verwijderen
- `user:manage_roles` - Gebruikersrollen beheren

**Chat:**
- `chat:read` - Chat kanalen bekijken
- `chat:write` - Berichten verzenden
- `chat:manage_channel` - Kanalen beheren
- `chat:moderate` - Berichten modereren

**Admin Email:**
- `admin_email:send` - Emails verzenden namens admin

## 🔍 Permissie Controle

### Frontend Permissie Checking

```javascript
// Permissie checker utility
class PermissionChecker {
  constructor(userPermissions = []) {
    this.permissions = new Set(userPermissions);
  }

  hasPermission(resource, action) {
    return this.permissions.has(`${resource}:${action}`);
  }

  hasAnyPermission(...permissions) {
    return permissions.some(perm => {
      const [resource, action] = perm.split(':');
      return this.hasPermission(resource, action);
    });
  }

  hasAllPermissions(...permissions) {
    return permissions.every(perm => {
      const [resource, action] = perm.split(':');
      return this.hasPermission(resource, action);
    });
  }
}

// Gebruik in componenten
const permissionChecker = new PermissionChecker(user.permissions);

if (permissionChecker.hasPermission('contact', 'write')) {
  // Toon bewerk knoppen
}

if (permissionChecker.hasAnyPermission('admin_email:send', 'newsletter:send')) {
  // Toon verzend opties
}
```

### React Hook voor Permissies

```javascript
// usePermissions.js
import { useContext, useMemo } from 'react';
import { AuthContext } from './AuthContext';

export const usePermissions = () => {
  const { user } = useContext(AuthContext);

  const permissions = useMemo(() => {
    if (!user?.permissions) return new Set();

    return new Set(
      user.permissions.map(p =>
        `${p.resource}:${p.action}`
      )
    );
  }, [user?.permissions]);

  const hasPermission = (resource, action) => {
    return permissions.has(`${resource}:${action}`);
  };

  const hasAnyPermission = (...perms) => {
    return perms.some(perm => {
      const [resource, action] = perm.split(':');
      return hasPermission(resource, action);
    });
  };

  const hasAllPermissions = (...perms) => {
    return perms.every(perm => {
      const [resource, action] = perm.split(':');
      return hasPermission(resource, action);
    });
  };

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    permissions: Array.from(permissions)
  };
};

// Gebruik in componenten
const MyComponent = () => {
  const { hasPermission, hasAnyPermission } = usePermissions();

  return (
    <div>
      {hasPermission('contact', 'read') && (
        <ContactList />
      )}

      {hasAnyPermission('contact:write', 'aanmelding:write') && (
        <EditButton />
      )}
    </div>
  );
};
```

## 🚀 API Endpoints

### Authenticatie Endpoints

#### POST /api/auth/login
```javascript
const login = async (email, password) => {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });

  const data = await response.json();

  if (response.ok) {
    // Sla token op
    localStorage.setItem('jwtToken', data.token);
    // Sla refresh token op indien beschikbaar
    if (data.refresh_token) {
      localStorage.setItem('refreshToken', data.refresh_token);
    }
    // Haal gebruikersinfo op inclusief permissies
    await loadUserProfile();
  }

  return data;
};
```

#### GET /api/auth/profile
```javascript
const loadUserProfile = async () => {
  const token = localStorage.getItem('jwtToken');
  const response = await fetch('/api/auth/profile', {
    headers: { 'Authorization': `Bearer ${token}` }
  });

  if (response.ok) {
    const user = await response.json();
    // user object bevat:
    // {
    //   id: "user-uuid",
    //   email: "user@example.com",
    //   naam: "Gebruikersnaam",
    //   permissions: [
    //     { resource: "contact", action: "read" },
    //     { resource: "contact", action: "write" },
    //     { resource: "user", action: "read" },
    //     // ... meer permissies
    //   ],
    //   roles: [
    //     { id: "role-uuid", name: "admin", description: "..." }
    //   ]
    // }
    setUser(user);
    return user;
  }

  throw new Error('Failed to load profile');
};
```

### Login Systeem Integratie met RBAC

#### Complete Login Flow

```javascript
// AuthContext.js
import React, { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Controleer bij app start of gebruiker ingelogd is
  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      // Controleer token geldigheid en laad user data
      loadUserProfile().catch(() => {
        // Token ongeldig, logout
        logout();
      });
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email, password) => {
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      const data = await response.json();

      if (response.ok) {
        // Sla tokens op
        localStorage.setItem('jwtToken', data.token);
        if (data.refresh_token) {
          localStorage.setItem('refreshToken', data.refresh_token);
        }

        // Laad user profile met permissies
        await loadUserProfile();
        return { success: true };
      } else {
        return { success: false, error: data.error };
      }
    } catch (error) {
      return { success: false, error: 'Netwerk fout' };
    }
  };

  const loadUserProfile = async () => {
    const token = localStorage.getItem('jwtToken');
    const response = await fetch('/api/auth/profile', {
      headers: { 'Authorization': `Bearer ${token}` }
    });

    if (response.ok) {
      const userData = await response.json();
      setUser(userData);
      setLoading(false);
      return userData;
    } else if (response.status === 401) {
      // Token verlopen, probeer refresh
      await refreshToken();
    } else {
      throw new Error('Failed to load profile');
    }
  };

  const refreshToken = async () => {
    const refreshToken = localStorage.getItem('refreshToken');
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
        localStorage.setItem('jwtToken', data.token);
        await loadUserProfile();
      } else {
        logout();
      }
    } catch {
      logout();
    }
  };

  const logout = () => {
    localStorage.removeItem('jwtToken');
    localStorage.removeItem('refreshToken');
    setUser(null);
    setLoading(false);
  };

  const isAuthenticated = () => {
    return user !== null && localStorage.getItem('jwtToken');
  };

  return (
    <AuthContext.Provider value={{
      user,
      login,
      logout,
      isAuthenticated,
      loading
    }}>
      {children}
    </AuthContext.Provider>
  );
};
```

#### Login Component met RBAC

```javascript
// Login.js
import React, { useState } from 'react';
import { useAuth } from './AuthContext';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const result = await login(email, password);

    if (result.success) {
      // Redirect naar dashboard of home
      navigate('/dashboard');
    } else {
      setError(result.error);
    }

    setLoading(false);
  };

  return (
    <div className="login-form">
      <h2>Inloggen</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Wachtwoord:</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <div className="error">{error}</div>}
        <button type="submit" disabled={loading}>
          {loading ? 'Bezig met inloggen...' : 'Inloggen'}
        </button>
      </form>
    </div>
  );
};

export default Login;
```

#### Permission-Based Navigation

```javascript
// App.js
import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './AuthContext';
import { usePermissions } from './usePermissions';
import Login from './Login';
import Dashboard from './Dashboard';
import AdminPanel from './AdminPanel';
import ContactManager from './ContactManager';
import UserManager from './UserManager';

const ProtectedRoute = ({ children, requiredPermission }) => {
  const { isAuthenticated, loading } = useAuth();
  const { hasPermission } = usePermissions();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated()) {
    return <Navigate to="/login" />;
  }

  if (requiredPermission && !hasPermission(...requiredPermission.split(':'))) {
    return <Navigate to="/access-denied" />;
  }

  return children;
};

const AppRoutes = () => {
  const { isAuthenticated } = useAuth();

  return (
    <Routes>
      <Route path="/login" element={
        isAuthenticated() ? <Navigate to="/dashboard" /> : <Login />
      } />

      <Route path="/dashboard" element={
        <ProtectedRoute>
          <Dashboard />
        </ProtectedRoute>
      } />

      <Route path="/contacts" element={
        <ProtectedRoute requiredPermission="contact:read">
          <ContactManager />
        </ProtectedRoute>
      } />

      <Route path="/users" element={
        <ProtectedRoute requiredPermission="user:read">
          <UserManager />
        </ProtectedRoute>
      } />

      <Route path="/admin" element={
        <ProtectedRoute requiredPermission="admin:access">
          <AdminPanel />
        </ProtectedRoute>
      } />

      <Route path="/" element={<Navigate to="/dashboard" />} />
    </Routes>
  );
};

const App = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;
```

#### Dashboard met Permission-Based UI

```javascript
// Dashboard.js
import React from 'react';
import { useAuth } from './AuthContext';
import { usePermissions } from './usePermissions';

const Dashboard = () => {
  const { user, logout } = useAuth();
  const { hasPermission, hasAnyPermission } = usePermissions();

  return (
    <div className="dashboard">
      <header>
        <h1>Welkom, {user.naam}!</h1>
        <button onClick={logout}>Uitloggen</button>
      </header>

      <div className="dashboard-grid">
        {/* Contact sectie - alleen zichtbaar als gebruiker contact permissies heeft */}
        {hasPermission('contact', 'read') && (
          <div className="card">
            <h3>Contacten</h3>
            <p>Beheer contactformulieren</p>
            <a href="/contacts">Bekijken</a>
            {hasPermission('contact', 'write') && (
              <span className="badge">Bewerk rechten</span>
            )}
          </div>
        )}

        {/* Gebruikers sectie - alleen voor staff/admin */}
        {hasPermission('user', 'read') && (
          <div className="card">
            <h3>Gebruikers</h3>
            <p>Beheer gebruikersaccounts</p>
            <a href="/users">Beheren</a>
          </div>
        )}

        {/* Admin sectie - alleen voor admin */}
        {hasPermission('admin', 'access') && (
          <div className="card admin-card">
            <h3>Admin Panel</h3>
            <p>Systeembeheer</p>
            <a href="/admin">Toegang</a>
          </div>
        )}

        {/* Nieuwsbrief sectie */}
        {hasAnyPermission('newsletter:read', 'newsletter:write') && (
          <div className="card">
            <h3>Nieuwsbrieven</h3>
            <p>Beheer nieuwsbrief verzending</p>
            <a href="/newsletters">Beheren</a>
          </div>
        )}
      </div>

      {/* Debug info voor development */}
      {process.env.NODE_ENV === 'development' && (
        <details className="debug-info">
          <summary>Debug Info</summary>
          <pre>{JSON.stringify(user, null, 2)}</pre>
        </details>
      )}
    </div>
  );
};

export default Dashboard;
```

### Beveiligde Endpoints met Permissie Controle

#### GET /api/contact (vereist: contact:read)
```javascript
const loadContacts = async () => {
  const token = localStorage.getItem('jwtToken');
  const response = await fetch('/api/contact', {
    headers: { 'Authorization': `Bearer ${token}` }
  });

  if (response.status === 403) {
    throw new Error('Geen toegang tot contacten');
  }

  return await response.json();
};
```

#### POST /api/contact/:id/antwoord (vereist: contact:write)
```javascript
const replyToContact = async (contactId, replyText) => {
  const token = localStorage.getItem('jwtToken');
  const response = await fetch(`/api/contact/${contactId}/antwoord`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ tekst: replyText })
  });

  return await response.json();
};
```

#### GET /api/users (vereist: user:read)
```javascript
const loadUsers = async () => {
  const token = localStorage.getItem('jwtToken');
  const response = await fetch('/api/users', {
    headers: { 'Authorization': `Bearer ${token}` }
  });

  return await response.json();
};
```

## 🔄 Redis Caching & Performance

### Wat is Redis?

Redis is een in-memory data structure store die wordt gebruikt voor:
- **Snelle permissie caching**: Gebruikerspermissies worden gecached voor betere performance
- **Distributed rate limiting**: API rate limits worden gedeeld over meerdere servers
- **Session storage**: JWT tokens en sessie data

### Frontend Impact

```javascript
// Permissies worden automatisch gecached door backend
// Frontend hoeft alleen maar de API te gebruiken

// Rate limiting wordt automatisch afgehandeld
// Frontend krijgt 429 status bij overschrijding

const handleApiCall = async (url, options = {}) => {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('jwtToken')}`,
      ...options.headers
    }
  });

  if (response.status === 429) {
    // Rate limit overschreden
    showRateLimitError();
    return;
  }

  if (response.status === 403) {
    // Permissie geweigerd
    showPermissionError();
    return;
  }

  return await response.json();
};
```

### Cache Invalidatie

```javascript
// Wanneer een gebruiker zijn rol verandert,
// wordt de cache automatisch ongeldig gemaakt door de backend
// Frontend hoeft alleen maar de permissies opnieuw op te halen

const refreshPermissions = async () => {
  const user = await loadUserProfile();
  // Update permission checker
  permissionChecker.updatePermissions(user.permissions);
};
```

## 🛡️ Error Handling

### HTTP Status Codes

| Status | Betekenis | Frontend Actie |
|--------|-----------|----------------|
| 401 | Niet geautoriseerd | Redirect naar login |
| 403 | Geen toegang | Toon permissie error |
| 429 | Rate limit overschreden | Toon wacht bericht |
| 500 | Server error | Toon algemene error |

### Error Response Format

```javascript
// Standaard error response
{
  "error": "Geen toegang",
  "code": "PERMISSION_DENIED",
  "details": {
    "required_permission": "contact:write",
    "user_permissions": ["contact:read"]
  }
}
```

### Frontend Error Handling

```javascript
const handleApiError = (error, response) => {
  switch (response.status) {
    case 401:
      // Token verlopen of ongeldig
      logout();
      break;

    case 403:
      // Permissie geweigerd
      showPermissionDenied(error.details);
      break;

    case 429:
      // Rate limit
      showRateLimitMessage();
      break;

    default:
      // Andere errors
      showGenericError(error.message);
  }
};
```

## 🔧 Development Best Practices

### 1. Permissie Checks in Componenten

```javascript
// Voorkom rendering van niet-toegankelijke UI
const AdminPanel = () => {
  const { hasPermission } = usePermissions();

  if (!hasPermission('admin', 'access')) {
    return <AccessDenied />;
  }

  return <AdminDashboard />;
};
```

### 2. API Call Wrappers

```javascript
// Centraliseer API calls met automatische error handling
class ApiClient {
  constructor(baseUrl = '/api') {
    this.baseUrl = baseUrl;
  }

  async request(endpoint, options = {}) {
    const token = localStorage.getItem('jwtToken');

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        ...options.headers
      }
    });

    if (!response.ok) {
      const error = await response.json();
      throw new ApiError(response.status, error);
    }

    return await response.json();
  }

  // Convenience methods
  get(endpoint) { return this.request(endpoint); }
  post(endpoint, data) { return this.request(endpoint, { method: 'POST', body: JSON.stringify(data) }); }
  put(endpoint, data) { return this.request(endpoint, { method: 'PUT', body: JSON.stringify(data) }); }
  delete(endpoint) { return this.request(endpoint, { method: 'DELETE' }); }
}

// Gebruik
const api = new ApiClient();

try {
  const contacts = await api.get('/contact');
} catch (error) {
  handleApiError(error);
}
```

### 3. Permission-Based Routing

```javascript
// React Router guard
const ProtectedRoute = ({ permission, children }) => {
  const { hasPermission } = usePermissions();
  const navigate = useNavigate();

  useEffect(() => {
    if (!hasPermission(...permission.split(':'))) {
      navigate('/access-denied');
    }
  }, [hasPermission, permission, navigate]);

  return children;
};

// Gebruik
<Route path="/admin" element={
  <ProtectedRoute permission="admin:access">
    <AdminPanel />
  </ProtectedRoute>
} />
```

### 4. Loading States & Permissions

```javascript
const ContactManager = () => {
  const { hasPermission } = usePermissions();
  const [contacts, setContacts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (hasPermission('contact', 'read')) {
      loadContacts();
    } else {
      setLoading(false);
    }
  }, [hasPermission]);

  const loadContacts = async () => {
    try {
      const data = await api.get('/contact');
      setContacts(data);
    } catch (error) {
      // Handle error
    } finally {
      setLoading(false);
    }
  };

  if (!hasPermission('contact', 'read')) {
    return <AccessDenied />;
  }

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div>
      {contacts.map(contact => (
        <ContactItem
          key={contact.id}
          contact={contact}
          canEdit={hasPermission('contact', 'write')}
          canDelete={hasPermission('contact', 'delete')}
        />
      ))}
    </div>
  );
};
```

## 📊 Monitoring & Debugging

### Health Checks

```javascript
// Controleer Redis en systeem health
const checkSystemHealth = async () => {
  const response = await fetch('/api/health');
  const health = await response.json();

  return {
    system: health.status === 'ok',
    redis: health.redis_connected || false,
    database: health.database_connected || false
  };
};
```

### Permission Debugging

```javascript
// Debug permissies in development
const debugPermissions = () => {
  if (process.env.NODE_ENV === 'development') {
    console.log('Current user permissions:', user.permissions);
    console.log('Available routes:', getAvailableRoutes(user.permissions));
  }
};
```

## 🔄 Migration Guide

### Van Legacy naar RBAC

Als je migreert van het oude systeem naar RBAC:

1. **Update authentication flow** om permissies op te halen
2. **Vervang role-based checks** met permission-based checks
3. **Update error handling** voor 403 responses
4. **Test alle protected routes** met verschillende rollen

### Breaking Changes

- Sommige endpoints vereisen nu specifiekere permissies
- Error responses bevatten meer details over vereiste permissies
- Rate limiting werkt nu distributed via Redis

## 📚 Resources

- [RBAC Concepten](https://en.wikipedia.org/wiki/Role-based_access_control)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Redis Documentation](https://redis.io/documentation)

---

**Let op:** Deze documentatie is specifiek voor de frontend integratie met het RBAC systeem. Backend ontwikkelaars dienen de backend RBAC documentatie te raadplegen voor implementatie details.