# Frontend Integration Guide

Complete guide for integrating your frontend application with the DKL Email Service API.

## Quick Start

### 1. Environment Configuration

Create a `.env` file in your frontend project:

```env
# API Configuration
REACT_APP_API_URL=http://localhost:8080
REACT_APP_WS_URL=ws://localhost:8080

# Production
# REACT_APP_API_URL=https://api.yourdomain.com
# REACT_APP_WS_URL=wss://api.yourdomain.com
```

### 2. API Client Setup

```typescript
// src/services/api.ts
import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Create axios instance
export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});

// Request interceptor - Add JWT token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor - Handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If 401 and not already retried
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const response = await axios.post(`${API_URL}/api/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token } = response.data.data;
        
        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);
        
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Refresh failed, redirect to login
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);
```

## Authentication

### Login

```typescript
// src/services/auth.service.ts
import { apiClient } from './api';

interface LoginRequest {
  email: string;
  password: string;
}

interface LoginResponse {
  success: boolean;
  data: {
    user: {
      id: number;
      email: string;
      naam: string;
      roles: string[];
    };
    access_token: string;
    refresh_token: string;
    expires_in: number;
  };
}

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>(
      '/api/auth/login',
      credentials
    );
    
    // Store tokens
    const { access_token, refresh_token } = response.data.data;
    localStorage.setItem('access_token', access_token);
    localStorage.setItem('refresh_token', refresh_token);
    
    return response.data;
  },

  async register(data: RegisterRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>(
      '/api/auth/register',
      data
    );
    
    const { access_token, refresh_token } = response.data.data;
    localStorage.setItem('access_token', access_token);
    localStorage.setItem('refresh_token', refresh_token);
    
    return response.data;
  },

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('refresh_token');
    
    try {
      await apiClient.post('/api/auth/logout', {
        refresh_token: refreshToken,
      });
    } finally {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  },

  async getCurrentUser() {
    const response = await apiClient.get('/api/auth/me');
    return response.data.data;
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('access_token');
  },
};
```

### React Authentication Hook

```typescript
// src/hooks/useAuth.ts
import { createContext, useContext, useState, useEffect } from 'react';
import { authService } from '../services/auth.service';

interface User {
  id: number;
  email: string;
  naam: string;
  roles: string[];
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Load user on mount
    if (authService.isAuthenticated()) {
      authService.getCurrentUser()
        .then(setUser)
        .catch(() => {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const response = await authService.login({ email, password });
    setUser(response.data.user);
  };

  const logout = async () => {
    await authService.logout();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAuthenticated: !!user,
      }}
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

### Protected Routes

```typescript
// src/components/ProtectedRoute.tsx
import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRoles?: string[];
}

export function ProtectedRoute({ children, requiredRoles }: ProtectedRouteProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <LoadingSpinner />;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRoles && !requiredRoles.some(role => user.roles.includes(role))) {
    return <Navigate to="/unauthorized" replace />;
  }

  return <>{children}</>;
}
```

## WebSocket Integration

### WebSocket Hook

```typescript
// src/hooks/useWebSocket.ts
import { useEffect, useState, useCallback, useRef } from 'react';

interface WebSocketMessage {
  type: string;
  data: any;
  timestamp: string;
}

interface UseWebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

export function useWebSocket(channel: string, options: UseWebSocketOptions = {}) {
  const [connected, setConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();

  const {
    onMessage,
    onConnect,
    onDisconnect,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5,
  } = options;

  const connect = useCallback(() => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      console.error('No access token available');
      return;
    }

    const wsUrl = process.env.REACT_APP_WS_URL || 'ws://localhost:8080';
    const ws = new WebSocket(`${wsUrl}/ws/${channel}?token=${token}`);

    ws.onopen = () => {
      console.log(`WebSocket connected to ${channel}`);
      setConnected(true);
      reconnectAttemptsRef.current = 0;
      onConnect?.();
    };

    ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        setLastMessage(message);
        onMessage?.(message);
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log(`WebSocket disconnected from ${channel}`);
      setConnected(false);
      onDisconnect?.();

      // Attempt reconnection
      if (reconnectAttemptsRef.current < maxReconnectAttempts) {
        reconnectAttemptsRef.current++;
        const delay = reconnectInterval * Math.pow(2, reconnectAttemptsRef.current - 1);
        
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log(`Reconnecting... (attempt ${reconnectAttemptsRef.current})`);
          connect();
        }, delay);
      }
    };

    wsRef.current = ws;
  }, [channel, onMessage, onConnect, onDisconnect, reconnectInterval, maxReconnectAttempts]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  }, []);

  const send = useCallback((message: any) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected');
    }
  }, []);

  useEffect(() => {
    connect();
    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    connected,
    lastMessage,
    send,
    disconnect,
    reconnect: connect,
  };
}
```

### Using WebSocket in Components

```typescript
// src/components/NotulenList.tsx
import { useWebSocket } from '../hooks/useWebSocket';
import { useState, useEffect } from 'react';

export function NotulenList() {
  const [notulen, setNotulen] = useState([]);

  const { connected, lastMessage } = useWebSocket('notulen', {
    onMessage: (message) => {
      switch (message.type) {
        case 'notulen_created':
          setNotulen(prev => [...prev, message.data]);
          break;
        case 'notulen_updated':
          setNotulen(prev =>
            prev.map(n => n.id === message.data.id ? message.data : n)
          );
          break;
        case 'notulen_deleted':
          setNotulen(prev => prev.filter(n => n.id !== message.data.id));
          break;
      }
    },
    onConnect: () => console.log('Connected to Notulen updates'),
    onDisconnect: () => console.log('Disconnected from Notulen updates'),
  });

  return (
    <div>
      <div>Status: {connected ? 'Connected' : 'Disconnected'}</div>
      {/* Render notulen list */}
    </div>
  );
}
```

## API Services

### Generic API Service

```typescript
// src/services/api.service.ts
import { apiClient } from './api';

export class ApiService<T> {
  constructor(private endpoint: string) {}

  async getAll(params?: any): Promise<T[]> {
    const response = await apiClient.get(this.endpoint, { params });
    return response.data.data;
  }

  async getById(id: number): Promise<T> {
    const response = await apiClient.get(`${this.endpoint}/${id}`);
    return response.data.data;
  }

  async create(data: Partial<T>): Promise<T> {
    const response = await apiClient.post(this.endpoint, data);
    return response.data.data;
  }

  async update(id: number, data: Partial<T>): Promise<T> {
    const response = await apiClient.put(`${this.endpoint}/${id}`, data);
    return response.data.data;
  }

  async delete(id: number): Promise<void> {
    await apiClient.delete(`${this.endpoint}/${id}`);
  }
}
```

### Specific Services

```typescript
// src/services/notulen.service.ts
import { ApiService } from './api.service';

interface Notulen {
  id: number;
  titel: string;
  datum: string;
  aanwezig: string[];
  afwezig: string[];
  verslag: string;
  actiepunten: string[];
}

export const notulenService = new ApiService<Notulen>('/api/notulen');

// src/services/album.service.ts
interface Album {
  id: number;
  titel: string;
  beschrijving: string;
  cover_foto_url: string;
}

export const albumService = new ApiService<Album>('/api/albums');
```

## File Upload

### Image Upload to Cloudinary

```typescript
// src/services/upload.service.ts
import { apiClient } from './api';

export const uploadService = {
  async uploadImage(file: File, folder?: string): Promise<string> {
    const formData = new FormData();
    formData.append('file', file);
    if (folder) {
      formData.append('folder', folder);
    }

    const response = await apiClient.post('/api/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / (progressEvent.total || 1)
        );
        console.log('Upload progress:', percentCompleted);
      },
    });

    return response.data.data.url;
  },

  async uploadMultiple(files: File[]): Promise<string[]> {
    const uploads = files.map(file => this.uploadImage(file));
    return Promise.all(uploads);
  },
};
```

### Upload Component

```typescript
// src/components/ImageUpload.tsx
import { useState } from 'react';
import { uploadService } from '../services/upload.service';

export function ImageUpload({ onUploadComplete }) {
  const [uploading, setUploading] = useState(false);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    setUploading(true);
    try {
      const urls = await uploadService.uploadMultiple(Array.from(files));
      onUploadComplete(urls);
    } catch (error) {
      console.error('Upload failed:', error);
      alert('Upload failed');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <input
        type="file"
        multiple
        accept="image/*"
        onChange={handleFileChange}
        disabled={uploading}
      />
      {uploading && <div>Uploading...</div>}
    </div>
  );
}
```

## Error Handling

### Global Error Handler

```typescript
// src/utils/errorHandler.ts
export function handleApiError(error: any): string {
  if (error.response) {
    // Server responded with error
    const { status, data } = error.response;
    
    switch (status) {
      case 400:
        return data.error || 'Invalid request';
      case 401:
        return 'Authentication required';
      case 403:
        return 'Permission denied';
      case 404:
        return 'Resource not found';
      case 429:
        return 'Too many requests. Please try again later.';
      case 500:
        return 'Server error. Please try again later.';
      default:
        return data.error || 'An error occurred';
    }
  } else if (error.request) {
    // Request made but no response
    return 'Network error. Please check your connection.';
  } else {
    // Error in request setup
    return error.message || 'An error occurred';
  }
}
```

## Best Practices

### 1. Token Management
- Store access token in memory (state)
- Store refresh token in httpOnly cookie (server-set) or secure localStorage
- Implement automatic token refresh
- Clear tokens on logout

### 2. Error Handling
- Use try-catch blocks
- Display user-friendly error messages
- Log errors for debugging
- Implement retry logic for network errors

### 3. Loading States
```typescript
const [loading, setLoading] = useState(false);
const [error, setError] = useState<string | null>(null);

try {
  setLoading(true);
  const data = await apiService.getData();
  // Handle success
} catch (err) {
  setError(handleApiError(err));
} finally {
  setLoading(false);
}
```

### 4. Caching
Use React Query or SWR for data caching:

```typescript
import { useQuery } from 'react-query';

function useNotulen() {
  return useQuery('notulen', () => notulenService.getAll(), {
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 10 * 60 * 1000, // 10 minutes
  });
}
```

### 5. TypeScript Types
Generate types from API schema or define manually:

```typescript
// src/types/api.types.ts
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}
```

## Testing

### Mock API Calls

```typescript
// src/__mocks__/api.ts
export const apiClient = {
  get: jest.fn(),
  post: jest.fn(),
  put: jest.fn(),
  delete: jest.fn(),
};

// src/services/__tests__/auth.service.test.ts
import { authService } from '../auth.service';
import { apiClient } from '../api';

jest.mock('../api');

describe('AuthService', () => {
  it('should login successfully', async () => {
    const mockResponse = {
      data: {
        success: true,
        data: {
          user: { id: 1, email: 'test@example.com' },
          access_token: 'token',
          refresh_token: 'refresh',
        },
      },
    };

    (apiClient.post as jest.Mock).mockResolvedValue(mockResponse);

    const result = await authService.login({
      email: 'test@example.com',
      password: 'password',
    });

    expect(result.success).toBe(true);
    expect(result.data.user.email).toBe('test@example.com');
  });
});
```

## Examples

See the [examples folder](../examples/) for complete implementation examples:
- [API Client Setup](../examples/api-client.ts)
- [Authentication Flow](../examples/auth-example.tsx)
- [WebSocket Integration](../examples/websocket-example.tsx)
- [File Upload](../examples/upload-example.tsx)

## Troubleshooting

### CORS Issues
Ensure the backend has your frontend URL in CORS configuration:
```
FRONTEND_URL=http://localhost:3000,https://yourdomain.com
```

### Token Expiration
Implement automatic token refresh in axios interceptor (shown above)

### WebSocket Connection Failures
- Check token is valid
- Verify WebSocket URL protocol (ws:// or wss://)
- Check firewall/proxy settings

---

For more details, see:
- [API Documentation](../api/README.md)
- [WebSocket Documentation](../api/WEBSOCKET.md)
- [Authentication Documentation](../api/AUTHENTICATION.md)