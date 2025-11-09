# Code Examples

Practical code examples for integrating with the DKL Email Service.

## Table of Contents

1. [API Client Setup](#api-client-setup)
2. [Authentication Examples](#authentication-examples)
3. [WebSocket Integration](#websocket-integration)
4. [File Upload Examples](#file-upload-examples)
5. [Steps Tracking](#steps-tracking-examples)
6. [Event Registration](#event-registration-examples)
7. [Chat Integration](#chat-integration-examples)
8. [Notification System](#notification-examples)

---

## API Client Setup

### TypeScript API Client (Complete)

```typescript
// src/services/api.ts
import axios, { AxiosInstance, AxiosError } from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

class ApiClient {
  private client: AxiosInstance;
  private refreshing: Promise<string> | null = null;

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
      timeout: 10000,
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as any;

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const newToken = await this.refreshToken();
            originalRequest.headers.Authorization = `Bearer ${newToken}`;
            return this.client(originalRequest);
          } catch (refreshError) {
            this.clearAuth();
            window.location.href = '/login';
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  private async refreshToken(): Promise<string> {
    if (this.refreshing) {
      return this.refreshing;
    }

    this.refreshing = (async () => {
      const refreshToken = localStorage.getItem('refresh_token');
      if (!refreshToken) {
        throw new Error('No refresh token');
      }

      const response = await axios.post(`${API_URL}/api/auth/refresh`, {
        refresh_token: refreshToken,
      });

      const { access_token, refresh_token: newRefreshToken } = response.data.data;
      
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', newRefreshToken);
      
      this.refreshing = null;
      return access_token;
    })();

    return this.refreshing;
  }

  private clearAuth() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  }

  // API methods
  get<T = any>(url: string, config?: any) {
    return this.client.get<ApiResponse<T>>(url, config);
  }

  post<T = any>(url: string, data?: any, config?: any) {
    return this.client.post<ApiResponse<T>>(url, data, config);
  }

  put<T = any>(url: string, data?: any, config?: any) {
    return this.client.put<ApiResponse<T>>(url, data, config);
  }

  delete<T = any>(url: string, config?: any) {
    return this.client.delete<ApiResponse<T>>(url, config);
  }
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  error?: string;
}

export const api = new ApiClient();
```

---

## Authentication Examples

### Complete Auth Service

```typescript
// src/services/auth.service.ts
import { api } from './api';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface User {
  id: string;
  email: string;
  naam: string;
  telefoon?: string;
  roles: string[];
  permissions: string[];
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

class AuthService {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/api/auth/login', credentials);
    const { access_token, refresh_token } = response.data.data;
    
    localStorage.setItem('access_token', access_token);
    localStorage.setItem('refresh_token', refresh_token);
    
    return response.data.data;
  }

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('refresh_token');
    
    try {
      await api.post('/api/auth/logout', { refresh_token: refreshToken });
    } finally {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  }

  async getCurrentUser(): Promise<User> {
    const response = await api.get<User>('/api/auth/profile');
    return response.data.data;
  }

  async updatePassword(currentPassword: string, newPassword: string): Promise<void> {
    await api.post('/api/auth/reset-password', {
      huidig_wachtwoord: currentPassword,
      nieuw_wachtwoord: newPassword,
    });
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem('access_token');
  }

  hasPermission(permission: string): boolean {
    const user = this.getCurrentUserSync();
    return user?.permissions?.includes(permission) || false;
  }

  hasRole(role: string): boolean {
    const user = this.getCurrentUserSync();
    return user?.roles?.includes(role) || false;
  }

  private getCurrentUserSync(): User | null {
    // Get from state management (Redux/Zustand/Context)
    return null; // Implement based on your state solution
  }
}

export const authService = new AuthService();
```

### React Auth Context

```typescript
// src/context/AuthContext.tsx
import React, { createContext, useContext, useState, useEffect } from 'react';
import { authService, User } from '../services/auth.service';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  hasPermission: (permission: string) => boolean;
  hasRole: (role: string) => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadUser();
  }, []);

  const loadUser = async () => {
    if (!authService.isAuthenticated()) {
      setLoading(false);
      return;
    }

    try {
      const userData = await authService.getCurrentUser();
      setUser(userData);
    } catch (err) {
      console.error('Failed to load user:', err);
      authService.logout();
    } finally {
      setLoading(false);
    }
  };

  const login = async (email: string, password: string) => {
    try {
      setError(null);
      const response = await authService.login({ email, password });
      setUser(response.user);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed');
      throw err;
    }
  };

  const logout = async () => {
    await authService.logout();
    setUser(null);
  };

  const hasPermission = (permission: string) => {
    return user?.permissions?.includes(permission) || false;
  };

  const hasRole = (role: string) => {
    return user?.roles?.includes(role) || false;
  };

  return (
    <AuthContext.Provider
      value={{ user, loading, error, login, logout, hasPermission, hasRole }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
```

### Vue Auth Composable

```typescript
// src/composables/useAuth.ts
import { ref, computed } from 'vue';
import { authService, type User } from '@/services/auth.service';

const user = ref<User | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

export function useAuth() {
  const loadUser = async () => {
    if (!authService.isAuthenticated()) {
      loading.value = false;
      return;
    }

    try {
      user.value = await authService.getCurrentUser();
    } catch (err) {
      console.error('Failed to load user:', err);
      await authService.logout();
    } finally {
      loading.value = false;
    }
  };

  const login = async (email: string, password: string) => {
    try {
      error.value = null;
      const response = await authService.login({ email, password });
      user.value = response.user;
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Login failed';
      throw err;
    }
  };

  const logout = async () => {
    await authService.logout();
    user.value = null;
  };

  const isAuthenticated = computed(() => !!user.value);
  
  const hasPermission = (permission: string) => {
    return user.value?.permissions?.includes(permission) || false;
  };

  const hasRole = (role: string) => {
    return user.value?.roles?.includes(role) || false;
  };

  return {
    user: computed(() => user.value),
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    isAuthenticated,
    login,
    logout,
    loadUser,
    hasPermission,
    hasRole,
  };
}
```

### Authentication Flow

```typescript
// Login
const loginResponse = await api.post('/api/auth/login', {
  email: 'user@example.com',
  password: 'password123',
});

const { access_token, refresh_token } = loginResponse.data.data;
localStorage.setItem('access_token', access_token);
localStorage.setItem('refresh_token', refresh_token);

// Get current user
const user = await api.get('/api/auth/me');
console.log(user.data.data);

// Logout
await api.post('/api/auth/logout', {
  refresh_token: localStorage.getItem('refresh_token'),
});
localStorage.removeItem('access_token');
localStorage.removeItem('refresh_token');
```

### WebSocket Connection

```typescript
// Connect to WebSocket
const token = localStorage.getItem('access_token');
const ws = new WebSocket(`ws://localhost:8080/ws/notulen?token=${token}`);

ws.onopen = () => {
  console.log('Connected');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
  
  switch (message.type) {
    case 'notulen_created':
      console.log('New notulen:', message.data);
      break;
    case 'notulen_updated':
      console.log('Updated notulen:', message.data);
      break;
    case 'notulen_deleted':
      console.log('Deleted notulen ID:', message.data.id);
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected');
};
```

### File Upload

```typescript
// Upload single image
const uploadImage = async (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await api.post('/api/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',

---

## Steps Tracking Examples

### React Steps Dashboard

```typescript
// src/components/StepsDashboard.tsx
import { useState, useEffect } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';
import { api } from '../services/api';

export function StepsDashboard({ registrationId }: { registrationId: string }) {
  const [steps, setSteps] = useState(0);
  const [distance, setDistance] = useState(0);
  const [rank, setRank] = useState(0);
  const [updating, setUpdating] = useState(false);

  const { connected } = useWebSocket('steps', {
    onMessage: (message) => {
      if (message.type === 'steps_update') {
        setSteps(message.data.steps);
        setDistance(message.data.total_distance);
      }
    }
  });

  const addSteps = async (delta: number) => {
    setUpdating(true);
    try {
      await api.post(`/api/registration/${registrationId}/steps`, {
        delta_steps: delta
      });
    } finally {
      setUpdating(false);
    }
  };

  return (
    <div className="steps-dashboard">
      <h2>{steps.toLocaleString()} stappen</h2>
      <p>{distance.toFixed(2)} km - Rank #{rank}</p>
      <button onClick={() => addSteps(1000)} disabled={updating}>
        +1000 stappen
      </button>
    </div>
  );
}
```

---

## Event Registration Example

```typescript
// src/services/registration.service.ts
import { api } from './api';

interface RegistrationData {
  naam: string;
  email: string;
  telefoon: string;
  rol: string;
  afstand: string;
  privacy_akkoord: boolean;
  terms: boolean;
}

export const registrationService = {
  async register(eventId: string, data: RegistrationData) {
    return api.post('/api/register', {
      ...data,
      event_id: eventId
    });
  },

  async getRegistrationDetails(registrationId: string) {
    return api.get(`/api/registration/${registrationId}`);
  },

  async updateSteps(registrationId: string, deltaSteps: number) {
    return api.post(`/api/registration/${registrationId}/steps`, {
      delta_steps: deltaSteps
    });
  }
};
```

---

## Complete Integration Example

### Full-Stack Event Tracker

```typescript
// src/pages/EventTracker.tsx
import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useWebSocket } from '../hooks/useWebSocket';
import { api } from '../services/api';

export function EventTracker() {
  const { user, hasPermission } = useAuth();
  const [event, setEvent] = useState(null);
  const [registration, setRegistration] = useState(null);
  const [leaderboard, setLeaderboard] = useState([]);

  // WebSocket for real-time updates
  const { connected } = useWebSocket('steps', {
    onMessage: (message) => {
      if (message.type === 'leaderboard_update') {
        setLeaderboard(message.data.top_users);
      }
    }
  });

  useEffect(() => {
    loadActiveEvent();
  }, []);

  const loadActiveEvent = async () => {
    const eventRes = await api.get('/api/events/active');
    setEvent(eventRes.data);
    
    if (user && eventRes.data.id) {
      const regRes = await api.get(`/api/registration/user/${user.id}`);
      setRegistration(regRes.data);
    }

    const leaderRes = await api.get('/api/leaderboard?period=daily&limit=10');
    setLeaderboard(leaderRes.data.entries);
  };

  return (
    <div className="event-tracker">
      <h1>{event?.name}</h1>
      
      {registration && (
        <StepsDashboard registrationId={registration.id} />
      )}

      <div className="leaderboard">
        <h2>Top 10 Today</h2>
        {leaderboard.map((entry, index) => (
          <div key={entry.participant_id} className="leaderboard-entry">
            <span className="rank">#{index + 1}</span>
            <span className="name">{entry.naam}</span>
            <span className="steps">{entry.steps.toLocaleString()}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

Voor meer voorbeelden en integratie patronen, zie:
- [Frontend Integration Guide](../guides/FRONTEND_INTEGRATION.md)
- [API Quick Reference](../api/QUICK_REFERENCE.md)
- [WebSocket API](../api/WEBSOCKET.md)

    },
  });
  
  return response.data.data.url;
};

// Usage
const file = document.querySelector('input[type="file"]').files[0];
const imageUrl = await uploadImage(file);
console.log('Uploaded to:', imageUrl);
```

### React Hooks

#### useApi Hook
```typescript
import { useState, useEffect } from 'react';
import api from './api';

export function useApi<T>(endpoint: string) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const response = await api.get(endpoint);
        setData(response.data.data);
        setError(null);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [endpoint]);

  return { data, loading, error };
}

// Usage
function AlbumsList() {
  const { data: albums, loading, error } = useApi<Album[]>('/api/albums');

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  
  return (
    <div>
      {albums?.map(album => (
        <div key={album.id}>{album.titel}</div>
      ))}
    </div>
  );
}
```

#### useWebSocket Hook
```typescript
import { useEffect, useState } from 'react';

export function useWebSocket(channel: string) {
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState<any[]>([]);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    const ws = new WebSocket(`ws://localhost:8080/ws/${channel}?token=${token}`);

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages(prev => [...prev, message]);
    };

    return () => ws.close();
  }, [channel]);

  return { connected, messages };
}

// Usage
function NotulenLive() {
  const { connected, messages } = useWebSocket('notulen');

  return (
    <div>
      <div>Status: {connected ? 'Connected' : 'Disconnected'}</div>
      {messages.map((msg, i) => (
        <div key={i}>{msg.type}: {JSON.stringify(msg.data)}</div>
      ))}
    </div>
  );
}
```

### Error Handling

```typescript
// Centralized error handler
export function handleApiError(error: any): string {
  if (error.response) {
    switch (error.response.status) {
      case 400:
        return 'Invalid request. Please check your input.';
      case 401:
        return 'Please log in to continue.';
      case 403:
        return 'You do not have permission to perform this action.';
      case 404:
        return 'Resource not found.';
      case 429:
        return 'Too many requests. Please try again later.';
      case 500:
        return 'Server error. Please try again later.';
      default:
        return error.response.data?.error || 'An error occurred.';
    }
  } else if (error.request) {
    return 'Network error. Please check your connection.';
  }
  return error.message || 'An unexpected error occurred.';
}

// Usage
try {
  await api.post('/api/albums', albumData);
} catch (error) {
  const errorMessage = handleApiError(error);
  alert(errorMessage);
}
```

### Form Validation

```typescript
// Validation helper
export const validators = {
  email: (value: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(value) ? null : 'Invalid email address';
  },
  
  password: (value: string) => {
    if (value.length < 8) return 'Password must be at least 8 characters';
    if (!/[A-Z]/.test(value)) return 'Password must contain uppercase letter';
    if (!/[a-z]/.test(value)) return 'Password must contain lowercase letter';
    if (!/[0-9]/.test(value)) return 'Password must contain number';
    return null;
  },
  
  required: (value: any) => {
    return value ? null : 'This field is required';
  },
};

// Usage in form
const [errors, setErrors] = useState({});

const handleSubmit = (e) => {
  e.preventDefault();
  const newErrors = {
    email: validators.email(email) || validators.required(email),
    password: validators.password(password) || validators.required(password),
  };
  
  if (Object.values(newErrors).some(error => error)) {
    setErrors(newErrors);
    return;
  }
  
  // Submit form
};
```

### Pagination

```typescript
// Paginated list component
function PaginatedList() {
  const [page, setPage] = useState(1);
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const limit = 20;

  useEffect(() => {
    const fetchData = async () => {
      const response = await api.get('/api/albums', {
        params: { page, limit },
      });
      setData(response.data.data);
      setTotal(response.data.pagination.total);
    };
    fetchData();
  }, [page]);

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      {/* Render data */}
      {data.map(item => <div key={item.id}>{item.titel}</div>)}
      
      {/* Pagination controls */}
      <div>
        <button 
          onClick={() => setPage(p => Math.max(1, p - 1))}
          disabled={page === 1}
        >
          Previous
        </button>
        <span>Page {page} of {totalPages}</span>
        <button 
          onClick={() => setPage(p => Math.min(totalPages, p + 1))}
          disabled={page === totalPages}
        >
          Next
        </button>
      </div>
    </div>
  );
}
```

### Search and Filter

```typescript
// Search with debounce
import { useState, useEffect } from 'react';

function SearchableList() {
  const [search, setSearch] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const timer = setTimeout(async () => {
      if (search) {
        setLoading(true);
        try {
          const response = await api.get('/api/albums', {
            params: { search },
          });
          setResults(response.data.data);
        } finally {
          setLoading(false);
        }
      } else {
        setResults([]);
      }
    }, 500); // 500ms debounce

    return () => clearTimeout(timer);
  }, [search]);

  return (
    <div>
      <input
        type="text"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search albums..."
      />
      {loading && <div>Searching...</div>}
      {results.map(album => (
        <div key={album.id}>{album.titel}</div>
      ))}
    </div>
  );
}
```

## Environment Configuration Examples

### React (.env)
```env
REACT_APP_API_URL=http://localhost:8080
REACT_APP_WS_URL=ws://localhost:8080
REACT_APP_CLOUDINARY_CLOUD_NAME=your_cloud_name
```

### Vue (.env)
```env
VUE_APP_API_URL=http://localhost:8080
VUE_APP_WS_URL=ws://localhost:8080
```

### Next.js (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## Testing Examples

### API Service Test
```typescript
import { describe, it, expect, jest } from '@jest/globals';
import api from './api';

jest.mock('./api');

describe('Album Service', () => {
  it('should fetch albums', async () => {
    const mockAlbums = [
      { id: 1, titel: 'Test Album' },
    ];
    
    (api.get as jest.Mock).mockResolvedValue({
      data: { data: mockAlbums },
    });

    const response = await api.get('/api/albums');
    expect(response.data.data).toEqual(mockAlbums);
  });
});
```

## More Resources

For complete implementation guides, see:
- [Frontend Integration Guide](../guides/FRONTEND_INTEGRATION.md)
- [API Documentation](../api/README.md)
- [WebSocket Guide](../api/WEBSOCKET.md)
- [Authentication Guide](../api/AUTHENTICATION.md)