# 🔧 Frontend AuthProvider Fix

## ❌ Probleem Gedetecteerd

De frontend gebruikt de environment variabele **verkeerd**:

```
URL: http://localhost:3000/VITE_API_BASE_URL=http://localhost:8082/api/api/auth/login
                           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
                           Dit is letterlijk de STRING, niet de WAARDE!
```

**Twee problemen:**
1. Environment variabele wordt als literal string gebruikt
2. Dubbel `/api` in de URL

## ✅ Oplossing

### Stap 1: Check je .env.development

Moet er **EXACT** zo uitzien:

```env
VITE_API_BASE_URL=http://localhost:8082/api
```

**NIET:**
```env
VITE_API_BASE_URL="http://localhost:8082/api"  # Geen quotes!
```

### Stap 2: Fix AuthProvider.tsx

**Vervang dit:**
```typescript
// ❌ FOUT - Literal string
const baseURL = "VITE_API_BASE_URL=http://localhost:8082/api";

// ❌ FOUT - Dubbel /api
fetch(`${baseURL}/api/auth/login`)
```

**Met dit:**
```typescript
// ✅ CORRECT - Haal waarde uit environment
const baseURL = import.meta.env.VITE_API_BASE_URL;

// ✅ CORRECT - Alleen /auth/login (baseURL bevat al /api)
fetch(`${baseURL}/auth/login`)
```

### Stap 3: Complete AuthProvider Voorbeeld

```typescript
// src/providers/AuthProvider.tsx
import React, { createContext, useContext, useState } from 'react';

// API BASE URL uit environment
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8082/api';

interface AuthContextType {
  user: any;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState<string | null>(
    localStorage.getItem('auth_token')
  );

  const login = async (email: string, password: string) => {
    try {
      // BELANGRIJK: API verwacht 'wachtwoord' (Nederlands), niet 'password'!
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          wachtwoord: password,  // <-- LET OP: 'wachtwoord'
        }),
      });

      if (!response.ok) {
        const error = await response.text();
        throw new Error(error || 'Login failed');
      }

      const data = await response.json();
      
      // Sla token op
      setToken(data.token);
      setUser(data.user);
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('refresh_token', data.refresh_token);
      
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('auth_token');
    localStorage.removeItem('refresh_token');
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
```

### Stap 4: Gebruik in LoginPage

```typescript
// src/pages/LoginPage.tsx
import { useState } from 'react';
import { useAuth } from '../providers/AuthProvider';

export const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    try {
      await login(email, password);
      // Redirect naar dashboard na succesvolle login
      window.location.href = '/dashboard';
    } catch (err: any) {
      setError(err.message || 'Login mislukt');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input 
        type="email" 
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="admin@dekoninklijkeloop.nl"
      />
      <input 
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="admin"
      />
      {error && <div className="error">{error}</div>}
      <button type="submit">Login</button>
    </form>
  );
};
```

## 🔍 Debug Checklist

### 1. Check Environment Variabele wordt Goed Geladen

```typescript
// In browser console
console.log('API URL:', import.meta.env.VITE_API_BASE_URL);
// Verwacht: "http://localhost:8082/api"
// NIET: undefined
// NIET: "VITE_API_BASE_URL=http://localhost:8082/api"
```

### 2. Check Vite Config (indien issue blijft)

**vite.config.ts:**
```typescript
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  
  console.log('Mode:', mode);
  console.log('API URL:', env.VITE_API_BASE_URL);
  
  return {
    // ... rest van config
  };
});
```

### 3. Herstart Development Server

```bash
# Stop frontend
Ctrl+C

# Herstart
npm run dev

# Check console output voor loaded env vars
```

### 4. Test API Direct

```bash
# Test of backend bereikbaar is
curl http://localhost:8082/api/health

# Test login direct (zou moeten werken)
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dekoninklijkeloop.nl","wachtwoord":"admin"}'
```

## 🎯 Common Frontend Mistakes

### Mistake 1: Quotes in .env
```env
# ❌ FOUT
VITE_API_BASE_URL="http://localhost:8082/api"

# ✅ CORRECT
VITE_API_BASE_URL=http://localhost:8082/api
```

### Mistake 2: Dubbel /api
```typescript
// ❌ FOUT
const url = `${baseURL}/api/auth/login`;
// Wordt: http://localhost:8082/api/api/auth/login

// ✅ CORRECT
const url = `${baseURL}/auth/login`;
// Wordt: http://localhost:8082/api/auth/login
```

### Mistake 3: Password vs Wachtwoord
```typescript
// ❌ FOUT
body: JSON.stringify({ email, password })

// ✅ CORRECT
body: JSON.stringify({ email, wachtwoord: password })
```

### Mistake 4: Environment niet herladen
```bash
# Na wijzigen .env file, HERSTART dev server
npm run dev
```

## 📝 Volledige Werkende Setup

**1. .env.development** (in frontend project root):
```env
VITE_API_BASE_URL=http://localhost:8082/api
```

**2. src/config/api.ts**:
```typescript
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8082/api';

console.log('✓ API configured:', API_BASE_URL);
```

**3. src/services/auth.ts**:
```typescript
import { API_BASE_URL } from '../config/api';

export const login = async (email: string, password: string) => {
  const url = `${API_BASE_URL}/auth/login`;
  console.log('Login URL:', url); // Debug
  
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ 
      email, 
      wachtwoord: password  // Nederlands!
    }),
  });

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
  }

  return await response.json();
};
```

**4. Gebruik**:
```typescript
const handleLogin = async (email: string, password: string) => {
  try {
    const data = await login(email, password);
    console.log('Success:', data);
    localStorage.setItem('auth_token', data.token);
  } catch (error) {
    console.error('Error:', error);
  }
};
```

## ✅ Verificatie

Na de fix, deze URL moet kloppen:

```
http://localhost:8082/api/auth/login
                    ^^^^            Geen dubbel /api
```

Test in browser console:
```javascript
fetch('http://localhost:8082/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@dekoninklijkeloop.nl',
    wachtwoord: 'admin'
  })
})
.then(r => r.json())
.then(console.log)
// Moet token teruggeven!
```

---

**Conclusie:** Dit is 100% een frontend configuratie issue. De backend werkt perfect!