# Frontend Under Construction (Maintenance Mode) API Documentatie

Deze documentatie beschrijft hoe je de maintenance mode status kunt ophalen en gebruiken in de frontend.

## Base URLs

- **Development (Docker):** `http://localhost:8082`
- **Production:** `https://dklemailservice.onrender.com`

## Beschikbare Endpoints

### 1. Get Active Maintenance Mode Status

**Endpoint:** `GET /api/under-construction/active`

Haalt de actieve maintenance mode configuratie op, als deze actief is.

**Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/under-construction/active');

if (response.ok) {
  const maintenanceMode = await response.json();
  // Toon maintenance page
} else if (response.status === 404) {
  // Geen maintenance mode - website is normaal beschikbaar
}
```

**Response (wanneer actief):**
```typescript
interface UnderConstruction {
  id: number;                          // ID van de maintenance mode record
  is_active: boolean;                  // Of maintenance mode actief is
  title: string;                       // Hoofdtitel (bijv. "Website in onderhoud")
  message: string;                     // Bericht aan gebruikers
  footer_text: string | null;          // Optionele footer tekst
  logo_url: string | null;             // URL naar logo
  expected_date: string | null;        // Verwachte datum van terugkeer (ISO 8601)
  social_links: SocialLink[] | null;   // Array van social media links
  progress_percentage: number | null;  // Voortgang percentage (0-100)
  contact_email: string | null;        // Contact email
  newsletter_enabled: boolean;         // Of newsletter inschrijving actief is
  created_at: string;                  // ISO timestamp
  updated_at: string;                  // ISO timestamp
}

interface SocialLink {
  platform: string;  // Bijv. "Twitter", "Instagram", "YouTube"
  url: string;       // URL naar social media profiel
}
```

**Voorbeeld Response (Actieve Maintenance Mode):**
```json
{
  "id": 1,
  "is_active": true,
  "title": "Website in onderhoud",
  "message": "We stomen ons klaar voor De Koninklijke Loop 2026, op dit moment is de website helaas niet bereikbaar",
  "footer_text": "Bedankt voor uw geduld!",
  "logo_url": "https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png",
  "expected_date": "2026-01-31T18:00:00Z",
  "social_links": [
    {
      "platform": "Twitter",
      "url": "https://twitter.com/koninklijkeloop"
    },
    {
      "platform": "Instagram",
      "url": "https://instagram.com/koninklijkeloop"
    },
    {
      "platform": "YouTube",
      "url": "https://www.youtube.com/@DeKoninklijkeLoop"
    }
  ],
  "progress_percentage": 85,
  "contact_email": "info@koninklijkeloop.nl",
  "newsletter_enabled": false,
  "created_at": "2025-09-26T17:37:22.197854Z",
  "updated_at": "2025-10-09T21:00:29.391392Z"
}
```

**Response (Geen Maintenance Mode):**
- **Status Code:** `404 Not Found`
- **Body:** `{ "error": "No active under construction found" }`

---

## React Voorbeelden

### Basic Maintenance Mode Check

```typescript
import React, { useEffect, useState } from 'react';

interface MaintenanceMode {
  title: string;
  message: string;
  logo_url: string | null;
  expected_date: string | null;
  progress_percentage: number | null;
}

const App: React.FC = () => {
  const [isMaintenanceMode, setIsMaintenanceMode] = useState(false);
  const [maintenanceData, setMaintenanceData] = useState<MaintenanceMode | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkMaintenanceMode = async () => {
      try {
        const response = await fetch(
          'https://dklemailservice.onrender.com/api/under-construction/active'
        );
        
        if (response.ok) {
          const data = await response.json();
          setIsMaintenanceMode(true);
          setMaintenanceData(data);
        } else if (response.status === 404) {
          // No maintenance mode active
          setIsMaintenanceMode(false);
        }
      } catch (error) {
        console.error('Error checking maintenance mode:', error);
        // On error, assume site is available
        setIsMaintenanceMode(false);
      } finally {
        setLoading(false);
      }
    };

    checkMaintenanceMode();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (isMaintenanceMode && maintenanceData) {
    return <MaintenancePage data={maintenanceData} />;
  }

  return <NormalApp />;
};

export default App;
```

### Full Maintenance Page Component

```typescript
import React from 'react';

interface MaintenancePageProps {
  data: {
    title: string;
    message: string;
    footer_text: string | null;
    logo_url: string | null;
    expected_date: string | null;
    social_links: Array<{ platform: string; url: string }> | null;
    progress_percentage: number | null;
    contact_email: string | null;
    newsletter_enabled: boolean;
  };
}

const MaintenancePage: React.FC<MaintenancePageProps> = ({ data }) => {
  const formatDate = (dateString: string | null): string => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('nl-NL', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getSocialIcon = (platform: string): string => {
    const icons: Record<string, string> = {
      'Twitter': '🐦',
      'Instagram': '📷',
      'YouTube': '▶️',
      'Facebook': '📘',
      'LinkedIn': '💼'
    };
    return icons[platform] || '🔗';
  };

  return (
    <div className="maintenance-container">
      {data.logo_url && (
        <div className="logo">
          <img src={data.logo_url} alt="Logo" />
        </div>
      )}
      
      <h1>{data.title}</h1>
      <p className="message">{data.message}</p>
      
      {data.expected_date && (
        <div className="expected-date">
          <h3>Verwachte terugkeer</h3>
          <p>{formatDate(data.expected_date)}</p>
        </div>
      )}
      
      {data.progress_percentage !== null && (
        <div className="progress">
          <div className="progress-label">
            Voortgang: {data.progress_percentage}%
          </div>
          <div className="progress-bar">
            <div 
              className="progress-fill"
              style={{ width: `${data.progress_percentage}%` }}
            />
          </div>
        </div>
      )}
      
      {data.social_links && data.social_links.length > 0 && (
        <div className="social-links">
          <h3>Volg ons</h3>
          <div className="social-icons">
            {data.social_links.map((link, index) => (
              <a 
                key={index}
                href={link.url}
                target="_blank"
                rel="noopener noreferrer"
                className="social-link"
              >
                <span className="icon">{getSocialIcon(link.platform)}</span>
                <span className="platform">{link.platform}</span>
              </a>
            ))}
          </div>
        </div>
      )}
      
      {data.contact_email && (
        <div className="contact">
          <p>
            Vragen? Neem contact op via{' '}
            <a href={`mailto:${data.contact_email}`}>
              {data.contact_email}
            </a>
          </p>
        </div>
      )}
      
      {data.newsletter_enabled && (
        <div className="newsletter">
          <h3>Blijf op de hoogte</h3>
          <p>Schrijf je in voor onze nieuwsbrief</p>
          <form>
            <input 
              type="email" 
              placeholder="Je email adres"
              required
            />
            <button type="submit">Inschrijven</button>
          </form>
        </div>
      )}
      
      {data.footer_text && (
        <footer>
          <p>{data.footer_text}</p>
        </footer>
      )}
    </div>
  );
};

export default MaintenancePage;
```

### Periodic Check with Polling

```typescript
import React, { useEffect, useState, useCallback } from 'react';

const POLL_INTERVAL = 60000; // Check every minute

const useMaintenanceMode = () => {
  const [isMaintenanceMode, setIsMaintenanceMode] = useState(false);
  const [maintenanceData, setMaintenanceData] = useState(null);

  const checkMaintenanceMode = useCallback(async () => {
    try {
      const response = await fetch(
        'https://dklemailservice.onrender.com/api/under-construction/active'
      );
      
      if (response.ok) {
        const data = await response.json();
        setIsMaintenanceMode(true);
        setMaintenanceData(data);
      } else if (response.status === 404) {
        setIsMaintenanceMode(false);
        setMaintenanceData(null);
      }
    } catch (error) {
      console.error('Error checking maintenance mode:', error);
    }
  }, []);

  useEffect(() => {
    // Initial check
    checkMaintenanceMode();

    // Set up polling
    const interval = setInterval(checkMaintenanceMode, POLL_INTERVAL);

    return () => clearInterval(interval);
  }, [checkMaintenanceMode]);

  return { isMaintenanceMode, maintenanceData };
};

// Usage
const App: React.FC = () => {
  const { isMaintenanceMode, maintenanceData } = useMaintenanceMode();

  if (isMaintenanceMode && maintenanceData) {
    return <MaintenancePage data={maintenanceData} />;
  }

  return <NormalApp />;
};
```

---

## Next.js App Router Voorbeeld

### Layout with Maintenance Check

```typescript
// app/layout.tsx
import { MaintenanceChecker } from './components/MaintenanceChecker';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="nl">
      <body>
        <MaintenanceChecker>
          {children}
        </MaintenanceChecker>
      </body>
    </html>
  );
}
```

```typescript
// app/components/MaintenanceChecker.tsx
'use client';

import { useEffect, useState } from 'react';
import MaintenancePage from './MaintenancePage';

export function MaintenanceChecker({ children }: { children: React.ReactNode }) {
  const [status, setStatus] = useState<'loading' | 'maintenance' | 'active'>('loading');
  const [data, setData] = useState(null);

  useEffect(() => {
    fetch('https://dklemailservice.onrender.com/api/under-construction/active')
      .then(async (res) => {
        if (res.ok) {
          const data = await res.json();
          setStatus('maintenance');
          setData(data);
        } else {
          setStatus('active');
        }
      })
      .catch(() => setStatus('active'));
  }, []);

  if (status === 'loading') {
    return <div>Loading...</div>;
  }

  if (status === 'maintenance' && data) {
    return <MaintenancePage data={data} />;
  }

  return <>{children}</>;
}
```

### Server-Side Check (Middleware)

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export async function middleware(request: NextRequest) {
  // Skip check for static files and API routes
  if (
    request.nextUrl.pathname.startsWith('/_next') ||
    request.nextUrl.pathname.startsWith('/api') ||
    request.nextUrl.pathname.includes('.')
  ) {
    return NextResponse.next();
  }

  try {
    const response = await fetch(
      'https://dklemailservice.onrender.com/api/under-construction/active',
      { next: { revalidate: 60 } } // Cache for 1 minute
    );

    if (response.ok) {
      // Maintenance mode is active, redirect to maintenance page
      if (request.nextUrl.pathname !== '/maintenance') {
        return NextResponse.redirect(new URL('/maintenance', request.url));
      }
    } else {
      // No maintenance mode
      if (request.nextUrl.pathname === '/maintenance') {
        return NextResponse.redirect(new URL('/', request.url));
      }
    }
  } catch (error) {
    console.error('Error checking maintenance mode:', error);
  }

  return NextResponse.next();
}

export const config = {
  matcher: '/((?!_next/static|_next/image|favicon.ico).*)',
};
```

---

## CSS Voorbeelden

### Modern Maintenance Page Styling

```css
.maintenance-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  text-align: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.logo img {
  max-width: 200px;
  margin-bottom: 2rem;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.2));
}

.maintenance-container h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  text-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.message {
  font-size: 1.25rem;
  max-width: 600px;
  margin: 0 auto 2rem;
  line-height: 1.6;
}

.expected-date {
  background: rgba(255,255,255,0.1);
  padding: 1.5rem;
  border-radius: 12px;
  margin: 2rem 0;
  backdrop-filter: blur(10px);
}

.progress {
  width: 100%;
  max-width: 400px;
  margin: 2rem auto;
}

.progress-bar {
  width: 100%;
  height: 30px;
  background: rgba(255,255,255,0.2);
  border-radius: 15px;
  overflow: hidden;
  margin-top: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #10b981, #34d399);
  transition: width 0.3s ease;
  border-radius: 15px;
}

.social-links {
  margin: 2rem 0;
}

.social-icons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.social-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(255,255,255,0.1);
  border-radius: 25px;
  text-decoration: none;
  color: white;
  transition: all 0.3s;
  backdrop-filter: blur(10px);
}

.social-link:hover {
  background: rgba(255,255,255,0.2);
  transform: translateY(-2px);
}

.contact {
  margin: 2rem 0;
}

.contact a {
  color: white;
  text-decoration: underline;
}

.newsletter {
  background: rgba(255,255,255,0.1);
  padding: 2rem;
  border-radius: 12px;
  max-width: 500px;
  margin: 2rem auto;
  backdrop-filter: blur(10px);
}

.newsletter form {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

.newsletter input {
  flex: 1;
  padding: 0.75rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
}

.newsletter button {
  padding: 0.75rem 1.5rem;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 8px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s;
}

.newsletter button:hover {
  transform: scale(1.05);
}

footer {
  margin-top: auto;
  padding-top: 2rem;
  opacity: 0.8;
}
```

---

## TypeScript Types

```typescript
// types/maintenance.ts
export interface UnderConstruction {
  id: number;
  is_active: boolean;
  title: string;
  message: string;
  footer_text: string | null;
  logo_url: string | null;
  expected_date: string | null;
  social_links: SocialLink[] | null;
  progress_percentage: number | null;
  contact_email: string | null;
  newsletter_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface SocialLink {
  platform: string;
  url: string;
}

export interface MaintenanceStatus {
  isActive: boolean;
  data: UnderConstruction | null;
}
```

---

## Belangrijke Opmerkingen

### 🚦 Status Codes
- **200 OK:** Maintenance mode IS actief - toon maintenance page
- **404 Not Found:** Geen maintenance mode - website is normaal beschikbaar
- **500 Internal Server Error:** Probleem met server - toon fallback

### 🔄 Polling Strategy
1. **Initial Check:** Bij laden van de applicatie
2. **Periodic Polling:** Check elke 1-5 minuten
3. **On Focus:** Check wanneer gebruiker terug naar tab komt
4. **Retry Logic:** Bij netwerk errors, probeer opnieuw na 30 seconden

### ⚡ Performance Tips
1. **Cache responses** voor 1 minuut om server te ontlasten
2. **Use middleware** voor Next.js om server-side te checken
3. **Graceful fallback** bij API errors (toon normale site)
4. **LocalStorage** voor laatste bekende status

### 🎨 UX Best Practices
1. **Toon progress** als percentage beschikbaar is
2. **Geef verwachte tijd** van terugkeer
3. **Bied alternatieven** zoals social media links
4. **Contact optie** voor urgente vragen
5. **Branding behouden** met logo en kleuren

### 📱 Responsive Design
- Mobile-first approach
- Touch-friendly buttons (min 44x44px)
- Leesbare font sizes (min 16px)
- Werkende social links op mobiel

### 🔐 Authenticatie
- **Geen authenticatie nodig** - public endpoint
- Check is beschikbaar voor alle gebruikers

---

## Production Status (November 2024)

**Current Status:** ✅ **NO MAINTENANCE MODE ACTIVE**

De website is normaal beschikbaar. De maintenance mode configuration is wel aanwezig in de database maar staat op `is_active: false`.

### Configuratie in Database:
```json
{
  "id": 1,
  "is_active": false,
  "title": "Website in onderhoud",
  "message": "We stomen ons klaar voor De Koninklijke Loop 2026...",
  "expected_date": "2026-01-31T18:00:00Z",
  "progress_percentage": 85,
  "contact_email": "info@koninklijkeloop.nl"
}
```

---

## Error Handling

```typescript
async function checkMaintenanceMode(): Promise<MaintenanceStatus> {
  try {
    const response = await fetch(
      'https://dklemailservice.onrender.com/api/under-construction/active',
      {
        // Add timeout
        signal: AbortSignal.timeout(5000)
      }
    );
    
    if (response.ok) {
      const data = await response.json();
      return { isActive: true, data };
    } else if (response.status === 404) {
      return { isActive: false, data: null };
    } else {
      throw new Error(`HTTP ${response.status}`);
    }
  } catch (error) {
    console.error('Maintenance check failed:', error);
    // On error, assume site is available (fail open)
    return { isActive: false, data: null };
  }
}
```

---

## Testing

```typescript
// Example test with Jest
describe('Maintenance Mode', () => {
  it('should show maintenance page when active', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          title: 'Under Maintenance',
          message: 'We will be back soon'
        })
      })
    );

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Under Maintenance')).toBeInTheDocument();
    });
  });

  it('should show normal app when no maintenance', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: false,
        status: 404
      })
    );

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Welcome')).toBeInTheDocument();
    });
  });
});
```

---

## Vragen?

Voor vragen of problemen met de API, neem contact op met het backend team of check de API documentatie in [`handlers/under_construction_handler.go`](../handlers/under_construction_handler.go).