# 🔧 Frontend Authentication Fix - Complete Guide

## ❌ Current Problem

The frontend is trying to use **TWO different authentication systems**:

1. **Supabase OIDC** (❌ FAILING - `Custom OIDC provider "custom" not allowed`)
2. **Backend JWT** (❌ Gets 401 errors)

### Error Analysis

```javascript
// ❌ WRONG - Frontend trying Supabase auth
installHook.js:1 Supabase sign in failed: AuthApiError: Custom OIDC provider "custom" not allowed
    at async N$.signInWithIdToken (GoTrueClient.js:457:25)
    at async p (AuthProvider.tsx:102:44)

// ❌ WRONG - Backend returns 401
dklemailservice.onrender.com/api/auth/login:1 Failed to load resource: the server responded with a status of 401

// ✅ GOOD - Backend permissions are loaded correctly
AuthProvider.tsx:175 ✅ Backend permissions loaded: 75 permissions

// ❌ WRONG - Photos fail because auth is broken
PhotosOverview.tsx Error loading data: Error: Failed to fetch photos
```

## ✅ Solution: Pure Backend JWT Auth

The backend uses **JWT tokens with RBAC** (see [`docs/AUTH_SYSTEM.md`](./AUTH_SYSTEM.md)). The frontend must **REMOVE all Supabase code** and use pure backend JWT.

---

## 📝 Step 1: Fix AuthProvider.tsx

Replace the entire AuthProvider with this correct implementation:

```typescript
// src/features/auth/contexts/AuthProvider.tsx
import React, { createContext, useState, useEffect, useCallback } from 'react';

// Get API base URL from environment
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://dklemailservice.onrender.com/api';

console.log('🔍 DEBUG: VITE_API_BASE_URL =', import.meta.env.VITE_API_BASE_URL);
console.log('✅ Using API Base URL:', API_BASE_URL);

interface User {
  id: string;
  email: string;
  naam: string;
  rol: string;
  permissions: Array<{ resource: string; action: string }>;
  roles: Array<{
    id: string;
    name: string;
    description: string;
  }>;
  is_actief: boolean;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<string | null>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshTimer, setRefreshTimer] = useState<NodeJS.Timeout | null>(null);

  // Load user on mount
  useEffect(() => {
    loadUser();
  }, []);

  // Schedule token refresh (15 minutes = 5 min before 20 min expiry)
  const scheduleRefresh = useCallback(() => {
    if (refreshTimer) {
      clearTimeout(refreshTimer);
    }
    const timer = setTimeout(() => {
      console.log('🔄 Auto-refreshing token...');
      refreshToken();
    }, 15 * 60 * 1000); // 15 minutes
    setRefreshTimer(timer);
  }, [refreshTimer]);

  // Load user profile from backend
  const loadUser = async () => {
    try {
      const token = localStorage.getItem('auth_token');
      if (!token) {
        setLoading(false);
        return;
      }

      const response = await fetch(`${API_BASE_URL}/auth/profile`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          // Try to refresh token
          const newToken = await refreshToken();
          if (newToken) {
            // Retry with new token
            await loadUser();
            return;
          }
        }
        throw new Error('Failed to load profile');
      }

      const userData = await response.json();
      setUser(userData);
      console.log('✅ Backend permissions loaded:', userData.permissions?.length || 0, 'permissions');
      
      // Schedule token refresh
      scheduleRefresh();
    } catch (error) {
      console.error('Error loading user:', error);
      // Clear invalid tokens
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
    } finally {
      setLoading(false);
    }
  };

  // Login with backend JWT
  const login = async (email: string, password: string): Promise<{ success: boolean; error?: string }> => {
    try {
      console.log('🔐 Login attempt for:', email);
      
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          wachtwoord: password, // Backend expects Dutch field name
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Login failed' }));
        console.error('❌ Login failed:', response.status, errorData);
        return { 
          success: false, 
          error: errorData.error || `Login failed (${response.status})` 
        };
      }

      const data = await response.json();
      console.log('✅ Login successful');

      // Store tokens
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('refresh_token', data.refresh_token);

      // Set user data
      setUser(data.user);

      // Schedule token refresh
      scheduleRefresh();

      return { success: true };
    } catch (error) {
      console.error('❌ Login error:', error);
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Network error' 
      };
    }
  };

  // Refresh access token
  const refreshToken = async (): Promise<string | null> => {
    try {
      const refresh = localStorage.getItem('refresh_token');
      if (!refresh) {
        throw new Error('No refresh token');
      }

      const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          refresh_token: refresh,
        }),
      });

      if (!response.ok) {
        throw new Error('Token refresh failed');
      }

      const data = await response.json();
      
      // Store new tokens
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('refresh_token', data.refresh_token);

      // Schedule next refresh
      scheduleRefresh();

      console.log('✅ Token refreshed successfully');
      return data.token;
    } catch (error) {
      console.error('❌ Token refresh failed:', error);
      // Logout on refresh failure
      await logout();
      return null;
    }
  };

  // Logout
  const logout = async () => {
    try {
      // Call backend logout endpoint
      const token = localStorage.getItem('auth_token');
      if (token) {
        await fetch(`${API_BASE_URL}/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }).catch(() => {
          // Ignore errors, logout anyway
        });
      }
    } finally {
      // Clear local state
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
      setUser(null);
      if (refreshTimer) {
        clearTimeout(refreshTimer);
      }
      console.log('✅ Logged out');
    }
  };

  const value: AuthContextType = {
    user,
    loading,
    isAuthenticated: !!user,
    login,
    logout,
    refreshToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
```

---

## 📝 Step 2: Create usePermissions Hook

```typescript
// src/hooks/usePermissions.ts
import { useAuth } from '@/features/auth/contexts/AuthProvider';

export const usePermissions = () => {
  const { user } = useAuth();

  const hasPermission = (resource: string, action: string): boolean => {
    if (!user?.permissions) return false;
    return user.permissions.some(
      p => p.resource === resource && p.action === action
    );
  };

  const hasAnyPermission = (...permissions: string[]): boolean => {
    return permissions.some(perm => {
      const [resource, action] = perm.split(':');
      return hasPermission(resource, action);
    });
  };

  const hasAllPermissions = (...permissions: string[]): boolean => {
    return permissions.every(perm => {
      const [resource, action] = perm.split(':');
      return hasPermission(resource, action);
    });
  };

  const permissions = user?.permissions?.map(p => `${p.resource}:${p.action}`) || [];

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    permissions,
  };
};
```

---

## 📝 Step 3: Create API Client with Auth

```typescript
// src/lib/api-client.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://dklemailservice.onrender.com/api';

class APIClient {
  private getAuthHeader(): Record<string, string> {
    const token = localStorage.getItem('auth_token');
    if (token) {
      return { 'Authorization': `Bearer ${token}` };
    }
    return {};
  }

  async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...this.getAuthHeader(),
        ...options.headers,
      },
    });

    if (!response.ok) {
      if (response.status === 401) {
        // Token expired, redirect to login
        localStorage.removeItem('auth_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
      }
      const error = await response.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.fetch<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.fetch<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint: string, data: any): Promise<T> {
    return this.fetch<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.fetch<T>(endpoint, { method: 'DELETE' });
  }
}

export const apiClient = new APIClient();
```

---

## 📝 Step 4: Fix Photo Service

```typescript
// src/services/photoService.ts
import { apiClient } from '@/lib/api-client';

export interface Photo {
  id: string;
  cloudinary_public_id: string;
  url: string;
  thumbnail_url?: string;
  title?: string;
  description?: string;
  visible: boolean;
  order_number: number;
}

export const photoService = {
  async getPhotos(): Promise<Photo[]> {
    return apiClient.get<Photo[]>('/photos');
  },

  async getPhoto(id: string): Promise<Photo> {
    return apiClient.get<Photo>(`/photos/${id}`);
  },

  async createPhoto(data: Partial<Photo>): Promise<Photo> {
    return apiClient.post<Photo>('/photos', data);
  },

  async updatePhoto(id: string, data: Partial<Photo>): Promise<Photo> {
    return apiClient.put<Photo>(`/photos/${id}`, data);
  },

  async deletePhoto(id: string): Promise<void> {
    return apiClient.delete<void>(`/photos/${id}`);
  },
};
```

---

## 📝 Step 5: Fix PhotosOverview Component

```typescript
// src/pages/PhotosOverview.tsx
import React, { useEffect, useState } from 'react';
import { usePermissions } from '@/hooks/usePermissions';
import { photoService, Photo } from '@/services/photoService';

export const PhotosOverview: React.FC = () => {
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { hasPermission } = usePermissions();

  useEffect(() => {
    loadPhotos();
  }, []);

  const loadPhotos = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await photoService.getPhotos();
      setPhotos(data);
      console.log('✅ Photos loaded:', data.length);
    } catch (err) {
      console.error('❌ Error loading photos:', err);
      setError(err instanceof Error ? err.message : 'Failed to load photos');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Loading photos...</div>;
  if (error) return <div>Error: {error}</div>;

  const canWrite = hasPermission('photo', 'write');
  const canDelete = hasPermission('photo', 'delete');

  return (
    <div>
      <h1>Photos</h1>
      {canWrite && <button onClick={() => {/* Add photo */}}>Add Photo</button>}
      <div className="photo-grid">
        {photos.map(photo => (
          <div key={photo.id}>
            <img src={photo.url} alt={photo.title} />
            {canWrite && <button>Edit</button>}
            {canDelete && <button>Delete</button>}
          </div>
        ))}
      </div>
    </div>
  );
};
```

---

## 📝 Step 6: Environment Configuration

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api

# .env.production
VITE_API_BASE_URL=https://dklemailservice.onrender.com/api
```

**IMPORTANT:** No quotes, no trailing slashes!

---

## 📝 Step 7: Remove Supabase

1. **Remove Supabase packages:**
   ```bash
   npm uninstall @supabase/supabase-js
   ```

2. **Remove Supabase config files:**
   - Delete `src/lib/supabase.ts` (if exists)
   - Remove any Supabase initialization code

3. **Clean up imports:**
   - Remove all `import { supabase }` statements
   - Remove all `.signInWithIdToken()` calls

---

## 🧪 Testing Checklist

```bash
# 1. Test backend is running
curl https://dklemailservice.onrender.com/api/health

# 2. Test login
curl -X POST https://dklemailservice.onrender.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"your_password"}'

# Expected response:
# {
#   "success": true,
#   "token": "eyJhbGc...",
#   "refresh_token": "...",
#   "user": {
#     "id": "uuid",
#     "email": "admin@dekoninklijkeloop.nl",
#     "naam": "Admin",
#     "rol": "admin",
#     "permissions": [...]
#   }
# }

# 3. Test profile endpoint
curl https://dklemailservice.onrender.com/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🎯 RBAC Permission Examples

```typescript
// Check single permission
if (hasPermission('photo', 'write')) {
  // Show edit button
}

// Check any of multiple permissions
if (hasAnyPermission('photo:write', 'photo:delete')) {
  // Show management panel
}

// Check all required permissions
if (hasAllPermissions('admin:access', 'user:manage_roles')) {
  // Show admin panel
}

// Conditional rendering
{hasPermission('photo', 'delete') && <DeleteButton />}
```

---

## 📚 Complete Permissions List

According to [`V1_21__seed_rbac_data.sql`](../database/migrations/V1_21__seed_rbac_data.sql):

| Resource | Actions | Description |
|----------|---------|-------------|
| `admin` | `access` | Full admin access |
| `staff` | `access` | Staff level access |
| `contact` | `read`, `write`, `delete` | Contact form management |
| `aanmelding` | `read`, `write`, `delete` | Registration management |
| `user` | `read`, `write`, `delete`, `manage_roles` | User management |
| `photo` | `read`, `write`, `delete` | Photo management |
| `album` | `read`, `write`, `delete` | Album management |
| `video` | `read`, `write`, `delete` | Video management |
| `newsletter` | `read`, `write`, `send`, `delete` | Newsletter management |
| `email` | `read`, `write`, `delete`, `fetch` | Email management |
| `chat` | `read`, `write`, `manage_channel`, `moderate` | Chat features |

**Total:** 19 resources, 58+ permissions

---

## ✅ Success Indicators

After implementing these fixes, you should see:

```javascript
// In browser console:
✅ Using API Base URL: https://dklemailservice.onrender.com/api
✅ Login successful
✅ Backend permissions loaded: 75 permissions
✅ Photos loaded: 25
✅ Token refreshed successfully
```

**No more:**
- ❌ Supabase errors
- ❌ 401 errors on `/api/auth/login`
- ❌ Failed to fetch photos

---

## 🚀 Deployment

1. Set environment variables in your hosting platform (Vercel/Netlify):
   ```
   VITE_API_BASE_URL=https://dklemailservice.onrender.com/api
   ```

2. Rebuild and deploy frontend

3. Test login at production URL

---

## 📞 Support

See [`docs/AUTH_SYSTEM.md`](./AUTH_SYSTEM.md) for complete RBAC documentation including:
- JWT token structure
- Permission middleware
- Role management
- Troubleshooting guide

---

**Status:** ✅ Ready for Implementation  
**Backend Compatibility:** V1.48.0+  
**Last Updated:** 2025-11-01