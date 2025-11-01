# 🎨 Frontend Integration Guide

> **Version:** 2.0 | **Backend:** V1.48.0+ | **Last Updated:** 2025-11-01

Complete guide voor het integreren van de DKL Email Service backend API met je frontend applicatie.

---

## 📋 Table of Contents

1. [Quick Start](#-quick-start)
2. [Environment Setup](#-environment-setup)
3. [Authentication](#-authentication)
4. [API Client](#-api-client)
5. [API Endpoints](#-api-endpoints)
6. [Service Layers](#-service-layers)
7. [Permission-Based UI](#-permission-based-ui)
8. [Development Workflow](#-development-workflow)
9. [Troubleshooting](#-troubleshooting)

---

## 🚀 Quick Start

### 3 Steps to Integration

**Step 1: Start Backend**
```bash
# In backend directory
docker-compose -f docker-compose.dev.yml up -d
```
✅ Backend now running on `http://localhost:8082`

**Step 2: Configure Frontend**

Create `.env.development` in your **frontend** project:

**For Vite/React:**
```env
VITE_API_BASE_URL=http://localhost:8082/api
```

**For Next.js:**
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8082/api
```

**Step 3: Setup API Client**

```typescript
// src/config/api.ts
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8082/api';

export const apiClient = {
  async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const token = localStorage.getItem('auth_token');
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers,
      },
    });
    
    if (!response.ok) {
      if (response.status === 401) {
        localStorage.clear();
        window.location.href = '/login';
      }
      throw new Error(`HTTP ${response.status}`);
    }
    
    return response.json();
  },
  
  get: <T>(endpoint: string) => apiClient.fetch<T>(endpoint, { method: 'GET' }),
  post: <T>(endpoint: string, data: any) => apiClient.fetch<T>(endpoint, {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  put: <T>(endpoint: string, data: any) => apiClient.fetch<T>(endpoint, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  delete: <T>(endpoint: string) => apiClient.fetch<T>(endpoint, { method: 'DELETE' }),
};
```

---

## ⚙️ Environment Setup

### Development vs Production

**Development (.env.development):**
```env
VITE_API_BASE_URL=http://localhost:8082/api
VITE_WS_URL=ws://localhost:8082/api/chat/ws
VITE_ENV=development
```

**Production (.env.production):**
```env
VITE_API_BASE_URL=https://dklemailservice.onrender.com/api
VITE_WS_URL=wss://dklemailservice.onrender.com/api/chat/ws
VITE_ENV=production
```

### CORS Configuration

Backend is **already configured** for:
- `http://localhost:3000` ✓ (React/Next.js default)
- `http://localhost:5173` ✓ (Vite default)
- `https://admin.dekoninklijkeloop.nl` ✓ (Production admin)
- `https://www.dekoninklijkeloop.nl` ✓ (Production site)

**Running on different port?** Add to backend `.env`:
```env
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:YOUR_PORT
```

---

## 🔐 Authentication

### Complete AuthProvider Implementation

```typescript
// src/features/auth/contexts/AuthProvider.tsx
import React, { createContext, useState, useEffect, useCallback } from 'react';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8082/api';

interface User {
  id: string;
  email: string;
  naam: string;
  rol: string;
  permissions: Array<{ resource: string; action: string }>;
  roles: Array<{ id: string; name: string; description: string }>;
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

  // Schedule token refresh (15 min = 5 min before 20 min expiry)
  const scheduleRefresh = useCallback(() => {
    if (refreshTimer) clearTimeout(refreshTimer);
    const timer = setTimeout(() => {
      console.log('🔄 Auto-refreshing token...');
      refreshToken();
    }, 15 * 60 * 1000);
    setRefreshTimer(timer);
  }, [refreshTimer]);

  const loadUser = async () => {
    try {
      const token = localStorage.getItem('auth_token');
      if (!token) {
        setLoading(false);
        return;
      }

      const response = await fetch(`${API_BASE_URL}/auth/profile`, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) {
        if (response.status === 401) {
          const newToken = await refreshToken();
          if (newToken) {
            await loadUser();
            return;
          }
        }
        throw new Error('Failed to load profile');
      }

      const userData = await response.json();
      setUser(userData);
      scheduleRefresh();
    } catch (error) {
      console.error('Error loading user:', error);
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
    } finally {
      setLoading(false);
    }
  };

  const login = async (email: string, password: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email,
          wachtwoord: password,  // Backend expects Dutch field!
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Login failed' }));
        return { 
          success: false, 
          error: errorData.error || `Login failed (${response.status})` 
        };
      }

      const data = await response.json();
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('refresh_token', data.refresh_token);
      setUser(data.user);
      scheduleRefresh();
      
      return { success: true };
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Network error' 
      };
    }
  };

  const refreshToken = async (): Promise<string | null> => {
    try {
      const refresh = localStorage.getItem('refresh_token');
      if (!refresh) throw new Error('No refresh token');

      const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refresh }),
      });

      if (!response.ok) throw new Error('Token refresh failed');

      const data = await response.json();
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('refresh_token', data.refresh_token);
      scheduleRefresh();
      
      return data.token;
    } catch (error) {
      await logout();
      return null;
    }
  };

  const logout = async () => {
    try {
      const token = localStorage.getItem('auth_token');
      if (token) {
        await fetch(`${API_BASE_URL}/auth/logout`, {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${token}` },
        }).catch(() => {});
      }
    } finally {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
      setUser(null);
      if (refreshTimer) clearTimeout(refreshTimer);
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, isAuthenticated: !!user, login, logout, refreshToken }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};
```

### usePermissions Hook

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

## 💻 API Client

### TypeScript Types

```typescript
// src/types/api.ts
export interface User {
  id: string;
  email: string;
  naam: string;
  rol: string;
  permissions: Permission[];
  roles: Role[];
  is_actief: boolean;
}

export interface Permission {
  resource: string;
  action: string;
}

export interface Role {
  id: string;
  name: string;
  description: string;
}

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

export interface Album {
  id: string;
  title: string;
  description?: string;
  cover_photo_id?: string;
  visible: boolean;
  order_number: number;
}

export interface Contact {
  id: string;
  naam: string;
  email: string;
  bericht: string;
  status: 'nieuw' | 'in_behandeling' | 'beantwoord' | 'gesloten';
  beantwoord: boolean;
  antwoorden_count: number;  // V1.48 denormalized
  created_at: string;
}
```

---

## 📡 API Endpoints

### Authentication

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@dekoninklijkeloop.nl",
  "wachtwoord": "password"  # Dutch field name!
}

Response:
{
  "success": true,
  "token": "eyJhbGc...",
  "refresh_token": "...",
  "user": {
    "id": "uuid",
    "email": "admin@dekoninklijkeloop.nl",
    "naam": "Admin",
    "permissions": [...],
    "roles": [...]
  }
}
```

```http
GET /api/auth/profile
Authorization: Bearer <token>

Response: <User object with permissions>
```

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}

Response:
{
  "success": true,
  "token": "new-token",
  "refresh_token": "new-refresh-token"
}
```

### Contact Management

**Permissions:** `contact:read`, `contact:write`, `contact:delete`

```http
GET /api/contact?limit=10&offset=0
Authorization: Bearer <token>

GET /api/contact/status/nieuw
Authorization: Bearer <token>

GET /api/contact/:id
Authorization: Bearer <token>

PUT /api/contact/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "in_behandeling",
  "notities": "Bezig met afhandelen"
}

POST /api/contact/:id/antwoord
Authorization: Bearer <token>
Content-Type: application/json

{
  "tekst": "Bedankt voor uw bericht..."
}

DELETE /api/contact/:id
Authorization: Bearer <token>
```

### Photos & Albums

**Permissions:** `photo:read/write/delete`, `album:read/write/delete`

```http
# Public endpoints (no auth)
GET /api/photos?year=2024&title=search
GET /api/albums
GET /api/albums/:id/photos

# Admin endpoints (auth required)
GET /api/photos/admin?limit=10&offset=0
POST /api/photos
{
  "cloudinary_public_id": "samples/photo123",
  "url": "https://...",
  "title": "Photo Title",
  "visible": true
}

PUT /api/photos/:id
DELETE /api/photos/:id

GET /api/albums/admin
POST /api/albums
{
  "title": "Album Title",
  "description": "...",
  "visible": true
}

PUT /api/albums/:id
DELETE /api/albums/:id

POST /api/albums/:id/photos
{
  "photo_id": "uuid",
  "order_number": 0
}

PUT /api/albums/:id/photos/reorder
{
  "photo_order": [
    {"photo_id": "uuid1", "order_number": 0},
    {"photo_id": "uuid2", "order_number": 1}
  ]
}
```

### Videos, Partners, Sponsors

Similar CRUD patterns for all content resources:

```http
GET /api/videos          # Public
GET /api/partners        # Public
GET /api/sponsors        # Public

GET /api/videos/admin    # Admin (requires permission)
POST /api/videos         # Create
PUT /api/videos/:id      # Update
DELETE /api/videos/:id   # Delete
```

**See:** [`AUTH_AND_RBAC.md`](AUTH_AND_RBAC.md) for complete permission requirements

---

## 🎯 Service Layers

### Photo Service Example

```typescript
// src/services/photoService.ts
import { apiClient } from '@/config/api';
import type { Photo } from '@/types/api';

export const photoService = {
  async getPhotos(filters?: { year?: number; title?: string }): Promise<Photo[]> {
    const params = new URLSearchParams(filters as any);
    return apiClient.get(`/photos?${params}`);
  },

  async getPhoto(id: string): Promise<Photo> {
    return apiClient.get(`/photos/${id}`);
  },

  async createPhoto(data: Partial<Photo>): Promise<Photo> {
    return apiClient.post('/photos', data);
  },

  async updatePhoto(id: string, data: Partial<Photo>): Promise<Photo> {
    return apiClient.put(`/photos/${id}`, data);
  },

  async deletePhoto(id: string): Promise<void> {
    return apiClient.delete(`/photos/${id}`);
  },
};
```

### Album Service Example

```typescript
// src/services/albumService.ts
export const albumService = {
  async getAlbums(includeCovers = false): Promise<Album[]> {
    return apiClient.get(`/albums?include_covers=${includeCovers}`);
  },

  async getAlbumPhotos(albumId: string): Promise<AlbumPhoto[]> {
    return apiClient.get(`/albums/${albumId}/photos`);
  },

  async createAlbum(data: Partial<Album>): Promise<Album> {
    return apiClient.post('/albums', data);
  },

  async addPhotoToAlbum(
    albumId: string, 
    photoId: string, 
    orderNumber = 0
  ): Promise<void> {
    return apiClient.post(`/albums/${albumId}/photos`, {
      photo_id: photoId,
      order_number: orderNumber
    });
  },

  async reorderPhotos(
    albumId: string, 
    photos: Array<{ photo_id: string; order_number: number }>
  ): Promise<void> {
    return apiClient.put(`/albums/${albumId}/photos/reorder`, {
      photo_order: photos
    });
  },

  async deleteAlbum(id: string): Promise<void> {
    return apiClient.delete(`/albums/${id}`);
  },
};
```

### Contact Service Example

```typescript
// src/services/contactService.ts
export const contactService = {
  async getMessages(limit = 10, offset = 0): Promise<Contact[]> {
    return apiClient.get(`/contact?limit=${limit}&offset=${offset}`);
  },

  async getMessagesByStatus(
    status: 'nieuw' | 'in_behandeling' | 'beantwoord' | 'gesloten'
  ): Promise<Contact[]> {
    return apiClient.get(`/contact/status/${status}`);
  },

  async updateStatus(
    id: string, 
    status: string, 
    notities?: string
  ): Promise<Contact> {
    return apiClient.put(`/contact/${id}`, { status, notities });
  },

  async addReply(id: string, tekst: string): Promise<void> {
    return apiClient.post(`/contact/${id}/antwoord`, { tekst });
  },

  async deleteMessage(id: string): Promise<void> {
    return apiClient.delete(`/contact/${id}`);
  },
};
```

---

## 🎨 Permission-Based UI

### Component with Permissions

```typescript
// src/pages/PhotosOverview.tsx
import { useEffect, useState } from 'react';
import { usePermissions } from '@/hooks/usePermissions';
import { photoService } from '@/services/photoService';

export const PhotosOverview = () => {
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const { hasPermission } = usePermissions();

  useEffect(() => {
    loadPhotos();
  }, []);

  const loadPhotos = async () => {
    try {
      setLoading(true);
      const data = await photoService.getPhotos();
      setPhotos(data);
    } catch (error) {
      console.error('Failed to load photos:', error);
    } finally {
      setLoading(false);
    }
  };

  const canWrite = hasPermission('photo', 'write');
  const canDelete = hasPermission('photo', 'delete');

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Photos</h1>
      {canWrite && <button onClick={handleAdd}>Add Photo</button>}
      
      <div className="photo-grid">
        {photos.map(photo => (
          <div key={photo.id}>
            <img src={photo.url} alt={photo.title} />
            {canWrite && <button onClick={() => handleEdit(photo)}>Edit</button>}
            {canDelete && <button onClick={() => handleDelete(photo.id)}>Delete</button>}
          </div>
        ))}
      </div>
    </div>
  );
};
```

### Permission Patterns

```typescript
// Check single permission
if (hasPermission('photo', 'write')) {
  // Show edit UI
}

// Check any of multiple permissions
if (hasAnyPermission('photo:write', 'photo:delete')) {
  // Show management panel
}

// Check all required permissions
if (hasAllPermissions('admin:access', 'user:manage_roles')) {
  // Show admin features
}

// Conditional rendering
{hasPermission('photo', 'delete') && <DeleteButton />}
```

---

## 🔄 Development Workflow

### Daily Workflow

**Morning:**
```bash
# 1. Start backend (in backend directory)
docker-compose -f docker-compose.dev.yml up -d

# 2. Verify backend running
curl http://localhost:8082/api/health

# 3. Start frontend (in frontend directory)
npm run dev
```

**Evening:**
```bash
# Backend (optional - can keep running)
docker-compose -f docker-compose.dev.yml down
```

### Test Credentials (Local Only)

**Admin Account:**
```
Email:    admin@dekoninklijkeloop.nl
Password: admin
```

**Test Login:**
```bash
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"admin"}'
```

### Database Access

```bash
# Via command line
docker exec -it dkl-postgres psql -U postgres -d dklemailservice

# Via GUI (pgAdmin/DBeaver)
Host: localhost
Port: 5433
Database: dklemailservice
User: postgres
Password: postgres
```

---

## 🐛 Troubleshooting

### CORS Errors

```
Access-Control-Allow-Origin error
```

**Solution:**
1. Check frontend runs on `localhost:3000` or `:5173`
2. If different port, add to backend `ALLOWED_ORIGINS`
3. Restart backend: `docker-compose -f docker-compose.dev.yml restart app`

### Connection Refused

```
ERR_CONNECTION_REFUSED
```

**Solution:**
```bash
# Check backend status
docker-compose -f docker-compose.dev.yml ps

# Start if not running
docker-compose -f docker-compose.dev.yml up -d
```

### 401 Unauthorized

**Causes:**
1. **Token expired** - Auto-refresh should handle this
2. **No token** - Check localStorage has `auth_token`
3. **Invalid token** - Clear localStorage and re-login
4. **Wrong endpoint** - Verify URL format

**Check:**
```javascript
// In browser console
console.log(localStorage.getItem('auth_token'));
console.log(import.meta.env.VITE_API_BASE_URL);
```

### Environment Variable Issues

**Problem:** URL shows literal "VITE_API_BASE_URL=..." string

**Solution:**
```typescript
// ❌ WRONG
const url = "VITE_API_BASE_URL=http://localhost:8082/api";

// ✅ CORRECT
const url = import.meta.env.VITE_API_BASE_URL;
```

**Also check:**
1. No quotes in `.env` file
2. No trailing slash in URL
3. Restart dev server after `.env` changes

### Backend Not Running

```bash
# Check status
docker-compose -f docker-compose.dev.yml ps

# View logs
docker logs dkl-email-service --tail 50

# Restart
docker-compose -f docker-compose.dev.yml restart app
```

---

## ✅ Integration Checklist

### Backend Setup
- [ ] Docker running (`docker-compose up -d`)
- [ ] Health endpoint responds (`/api/health`)
- [ ] Database migrations applied
- [ ] Test data loaded

### Frontend Setup
- [ ] Environment variables configured
- [ ] API client implemented
- [ ] AuthProvider integrated
- [ ] usePermissions hook available
- [ ] Route guards implemented

### Testing
- [ ] Login works
- [ ] Token refresh works
- [ ] Permission checks work
- [ ] Protected routes work
- [ ] API calls include Authorization header
- [ ] Error handling works

---

## 📊 API Endpoint Summary

**Public Endpoints** (no auth):
- Health, contact form submit, registration form, public content (albums, photos, videos, sponsors)

**Protected Endpoints** (auth + permission required):
- Contact management, user management, content management (admin versions), RBAC management

**Total Endpoints:** 50+  
**RBAC Protected:** 40+  
**Permission Resources:** 19  
**Permission Actions:** 58

See [`AUTH_AND_RBAC.md`](AUTH_AND_RBAC.md) for complete permission matrix.

---

## 🎓 Best Practices

### DO ✅

- Use environment variables for API URLs
- Include error handling in all API calls
- Implement loading states
- Check permissions before rendering UI
- Use service layers (don't call API directly in components)
- Clear tokens on logout
- Handle 401 errors gracefully

### DON'T ❌

- Hardcode API URLs
- Trust only frontend permission checks
- Forget to add Authorization header
- Store passwords (only hashed tokens)
- Use SELECT * in queries (backend optimized for specific fields)
- Manually set `updated_at` (backend triggers handle it)

---

## 📚 Related Documentation

- [`AUTH_AND_RBAC.md`](AUTH_AND_RBAC.md) - Complete RBAC system
- [`DATABASE_REFERENCE.md`](DATABASE_REFERENCE.md) - Database schema
- [`README.md`](README.md) - Documentation hub

---

## 🚀 Production URLs

**API Base URL:**
```
https://dklemailservice.onrender.com/api
```

**WebSocket (Chat):**
```
wss://dklemailservice.onrender.com/api/chat/ws
```

---

**Status:** ✅ Production Ready  
**Backend Version:** V1.48.0+  
**CORS:** Configured for common ports  
**Auth System:** JWT + RBAC fully operational  
**Last Verified:** 2025-11-01