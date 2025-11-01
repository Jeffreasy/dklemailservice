# Frontend API Complete Overview

Overzicht van alle beschikbare APIs voor de DKL Email Service frontend.

## 🌐 Base URLs

- **Development (Docker):** `http://localhost:8082`
- **Production:** `https://dklemailservice.onrender.com`

---

## 📚 Documentatie Per Module

### 1. 📸 [Albums & Photos API](FRONTEND_ALBUMS_API.md)
Photo albums met Cloudinary integratie

**Key Endpoints:**
- `GET /api/albums` - Alle zichtbare albums
- `GET /api/albums?include_covers=true` - Met cover photos
- `GET /api/albums/:id/photos` - Photos per album

**Production Data:**
- 3 albums (DKL 2025, DKL-2024, Voorbereidingen)
- 45 totale photos
- Cloudinary hosted images

### 2. 🎥 [Videos API](FRONTEND_VIDEOS_API.md)
Video gallery met Streamable embeds

**Key Endpoints:**
- `GET /api/videos` - Alle zichtbare videos

**Production Data:**
- 5 videos
- Streamable platform
- Auto-generated thumbnails

### 3. 🚧 [Under Construction API](FRONTEND_UNDER_CONSTRUCTION_API.md)
Maintenance mode / "Website in onderhoud"

**Key Endpoints:**
- `GET /api/under-construction/active` - Check maintenance status

**Use Cases:**
- Geplande onderhoud
- Site updates
- Coming soon pages
- Progress tracking

### 4. 📧 [Email API](FRONTEND_EMAIL_API.md)
Complete email functionaliteit (submission + management)

**Public Endpoints (Geen Auth):**
- `POST /api/contact-email` - Contact formulier
- `POST /api/aanmelding-email` - Registratie formulier
- `POST /api/wfc/order-email` - WFC order emails (API key)

**Admin Endpoints (Auth Required):**
- Contact management (CRUD)
- Aanmelding management (CRUD)
- Incoming email management
- Reply functionaliteit

---

## 🔐 Authenticatie Types

### Public Endpoints
**Geen authenticatie nodig:**
- `/api/albums`
- `/api/videos`
- `/api/under-construction/active`
- `/api/contact-email` (POST)
- `/api/aanmelding-email` (POST)

### JWT Authenticatie
**Bearer Token Required:**
```typescript
headers: {
  'Authorization': `Bearer ${token}`
}
```

**Endpoints:**
- Alle `/api/contact/*` (beheer)
- Alle `/api/aanmelding/*` (beheer)
- Alle `/api/mail/*` (inbox)
- Admin CRUD operaties

**Token verkrijgen:**
```typescript
// Login
const response = await fetch('https://dklemailservice.onrender.com/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'admin@example.com',
    wachtwoord: 'password'
  })
});

const { access_token, refresh_token } = await response.json();
```

### API Key Authenticatie
**X-API-Key Header:**
```typescript
headers: {
  'X-API-Key': process.env.WFC_API_KEY
}
```

**Alleen voor:**
- `/api/wfc/order-email`

---

## 📊 Quick Reference

### Status van Modules

| Module | Endpoints | Auth | Production Data |
|--------|-----------|------|-----------------|
| **Albums** | 9 | Public + Admin | 3 albums, 45 photos |
| **Videos** | 5 | Public + Admin | 5 videos |
| **Under Construction** | 5 | Public + Admin | 1 config (inactive) |
| **Email (Public)** | 3 | None | Live submission |
| **Email (Admin)** | 12+ | JWT | Management tools |

---

## 🚀 Quick Start Templates

### Public Website Integration

```typescript
// App.tsx - Main frontend entry
import { useEffect, useState } from 'react';

function App() {
  const [maintenanceMode, setMaintenanceMode] = useState(false);

  useEffect(() => {
    // Check maintenance mode
    fetch('https://dklemailservice.onrender.com/api/under-construction/active')
      .then(res => {
        if (res.ok) {
          res.json().then(data => {
            setMaintenanceMode(true);
            // Show maintenance page
          });
        }
      });
  }, []);

  if (maintenanceMode) {
    return <MaintenancePage />;
  }

  return (
    <div>
      <AlbumsGallery />
      <VideoGallery />
      <ContactForm />
      <AanmeldingForm />
    </div>
  );
}
```

### Data Fetching Hook

```typescript
// hooks/usePublicData.ts
import { useEffect, useState } from 'react';

export function usePublicData<T>(endpoint: string) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const baseUrl = 'https://dklemailservice.onrender.com';
    
    fetch(`${baseUrl}${endpoint}`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then(data => {
        setData(data);
        setLoading(false);
      })
      .catch(err => {
        setError(err);
        setLoading(false);
      });
  }, [endpoint]);

  return { data, loading, error };
}

// Usage
const { data: albums, loading } = usePublicData<Album[]>('/api/albums?include_covers=true');
const { data: videos } = usePublicData<Video[]>('/api/videos');
```

---

## 🧪 Test Scripts

Alle modules hebben dedicated test scripts:

| Script | Beschrijving | Gebruik |
|--------|--------------|---------|
| [`test_albums.ps1`](../test_albums.ps1) | Test album endpoints | `.\test_albums.ps1 -Environment production` |
| [`test_videos.ps1`](../test_videos.ps1) | Test video endpoints | `.\test_videos.ps1 -Environment production` |
| [`test_under_construction.ps1`](../test_under_construction.ps1) | Test maintenance mode | `.\test_under_construction.ps1 -Environment production` |
| [`test_email_api.ps1`](../test_email_api.ps1) | Test email endpoints | `.\test_email_api.ps1 -Environment production -TestType public` |

**All-in-one test:**
```powershell
# Test everything in production
.\test_albums.ps1 -Environment production
.\test_videos.ps1 -Environment production
.\test_under_construction.ps1 -Environment production
.\test_email_api.ps1 -Environment production -TestType public
```

---

## 📦 Complete TypeScript Type Definitions

### Combined Types File

```typescript
// types/api.ts

// ============================================
// ALBUMS & PHOTOS
// ============================================
export interface Album {
  id: string;
  title: string;
  description: string | null;
  cover_photo_id: string | null;
  visible: boolean;
  order_number: number;
  created_at: string;
  updated_at: string;
}

export interface Photo {
  id: string;
  title: string;
  description: string | null;
  url: string;
  thumbnail_url: string | null;
  cloudinary_id: string;
  visible: boolean;
  width: number | null;
  height: number | null;
  format: string | null;
  created_at: string;
  updated_at: string;
}

export interface AlbumWithCover extends Album {
  cover_photo: Photo | null;
}

export interface PhotoWithAlbumInfo extends Photo {
  order_number: number | null;
}

// ============================================
// VIDEOS
// ============================================
export interface Video {
  id: string;
  video_id: string;
  url: string;
  title: string;
  description: string | null;
  thumbnail_url: string | null;
  visible: boolean;
  order_number: number;
  created_at: string;
  updated_at: string;
}

// ============================================
// UNDER CONSTRUCTION
// ============================================
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

// ============================================
// EMAIL - PUBLIC SUBMISSION
// ============================================
export interface ContactFormulier {
  naam: string;
  email: string;
  telefoon?: string;
  bericht: string;
  privacy_akkoord: boolean;
  test_mode?: boolean;
}

export interface AanmeldingFormulier {
  naam: string;
  email: string;
  telefoon?: string;
  rol: 'deelnemer' | 'vrijwilliger';
  afstand: '5km' | '10km' | '15km';
  ondersteuning?: string;
  bijzonderheden?: string;
  terms: boolean;
  test_mode?: boolean;
}

export interface EmailSubmissionResponse {
  success: boolean;
  message: string;
  test_mode?: boolean;
  error?: string;
}

// ============================================
// EMAIL - ADMIN MANAGEMENT
// ============================================
export interface Contact {
  id: string;
  naam: string;
  email: string;
  telefoon: string | null;
  bericht: string;
  status: 'nieuw' | 'in_behandeling' | 'beantwoord' | 'gesloten';
  privacy_akkoord: boolean;
  notities: string | null;
  beantwoord: boolean;
  antwoord_tekst: string | null;
  antwoord_datum: string | null;
  antwoord_door: string | null;
  behandeld_door: string | null;
  behandeld_op: string | null;
  created_at: string;
  updated_at: string;
  antwoorden?: ContactAntwoord[];
}

export interface ContactAntwoord {
  id: string;
  contact_id: string;
  tekst: string;
  verzond_door: string;
  email_verzonden: boolean;
  created_at: string;
}

export interface Aanmelding {
  id: string;
  naam: string;
  email: string;
  telefoon: string | null;
  rol: string;
  afstand: string;
  ondersteuning: string | null;
  bijzonderheden: string | null;
  status: string;
  terms: boolean;
  test_mode: boolean;
  email_verzonden: boolean;
  email_verzonden_op: string | null;
  behandeld_door: string | null;
  behandeld_op: string | null;
  notities: string | null;
  created_at: string;
  updated_at: string;
  gebruiker_id: string | null;
  antwoorden?: AanmeldingAntwoord[];
}

export interface AanmeldingAntwoord {
  id: string;
  aanmelding_id: string;
  tekst: string;
  verzond_door: string;
  email_verzonden: boolean;
  created_at: string;
}

export interface MailResponse {
  id: string;
  message_id: string;
  sender: string;
  to: string;
  subject: string;
  html: string;
  content_type: string;
  received_at: string;
  uid: string;
  account_type: 'info' | 'inschrijving';
  read: boolean;
  processed_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface PaginatedMailResponse {
  emails: MailResponse[];
  totalCount: number;
}
```

---

## 🔄 Complete API Client Class

```typescript
// lib/api-client.ts
const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://dklemailservice.onrender.com';

export class DKLApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string = BASE_URL) {
    this.baseUrl = baseUrl;
  }

  setToken(token: string) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // ============================================
  // PUBLIC ENDPOINTS (No Auth)
  // ============================================

  // Albums
  async getAlbums(includeCovers = false) {
    const query = includeCovers ? '?include_covers=true' : '';
    return this.request<Album[]>(`/api/albums${query}`);
  }

  async getAlbumPhotos(albumId: string) {
    return this.request<PhotoWithAlbumInfo[]>(`/api/albums/${albumId}/photos`);
  }

  // Videos
  async getVideos() {
    return this.request<Video[]>('/api/videos');
  }

  // Under Construction
  async getMaintenanceStatus() {
    try {
      return await this.request<UnderConstruction>('/api/under-construction/active');
    } catch (error) {
      // 404 means no maintenance mode
      return null;
    }
  }

  // Email Submission
  async submitContactForm(data: ContactFormulier) {
    return this.request<EmailSubmissionResponse>('/api/contact-email', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async submitAanmelding(data: AanmeldingFormulier) {
    return this.request<EmailSubmissionResponse>('/api/aanmelding-email', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // ============================================
  // ADMIN ENDPOINTS (Auth Required)
  // ============================================

  // Contact Management
  async listContacts(limit = 10, offset = 0) {
    return this.request<Contact[]>(`/api/contact?limit=${limit}&offset=${offset}`);
  }

  async getContact(id: string) {
    return this.request<Contact>(`/api/contact/${id}`);
  }

  async updateContact(id: string, data: { status?: string; notities?: string }) {
    return this.request<Contact>(`/api/contact/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async addContactReply(contactId: string, tekst: string) {
    return this.request<ContactAntwoord>(`/api/contact/${contactId}/antwoord`, {
      method: 'POST',
      body: JSON.stringify({ tekst }),
    });
  }

  async deleteContact(id: string) {
    return this.request<{ success: boolean; message: string }>(`/api/contact/${id}`, {
      method: 'DELETE',
    });
  }

  // Aanmelding Management
  async listAanmeldingen(limit = 10, offset = 0) {
    return this.request<Aanmelding[]>(`/api/aanmelding?limit=${limit}&offset=${offset}`);
  }

  async getAanmelding(id: string) {
    return this.request<Aanmelding>(`/api/aanmelding/${id}`);
  }

  async updateAanmelding(id: string, data: { status?: string; notities?: string }) {
    return this.request<Aanmelding>(`/api/aanmelding/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Incoming Emails
  async listEmails(limit = 10, offset = 0) {
    return this.request<PaginatedMailResponse>(`/api/mail?limit=${limit}&offset=${offset}`);
  }

  async getEmail(id: string) {
    return this.request<MailResponse>(`/api/mail/${id}`);
  }

  async markEmailAsProcessed(id: string) {
    return this.request<{ message: string }>(`/api/mail/${id}/processed`, {
      method: 'PUT',
    });
  }

  async listUnprocessedEmails() {
    return this.request<MailResponse[]>('/api/mail/unprocessed');
  }
}

// Export singleton instance
export const apiClient = new DKLApiClient();

// Usage:
// import { apiClient } from '@/lib/api-client';
// const albums = await apiClient.getAlbums(true);
```

---

## 🎯 Common Patterns

### Loading States

```typescript
const [data, setData] = useState(null);
const [loading, setLoading] = useState(true);
const [error, setError] = useState(null);

useEffect(() => {
  apiClient.getAlbums(true)
    .then(setData)
    .catch(setError)
    .finally(() => setLoading(false));
}, []);

if (loading) return <LoadingSpinner />;
if (error) return <ErrorMessage error={error} />;
if (!data) return <NoData />;

return <AlbumsGrid albums={data} />;
```

### Form Submission with Feedback

```typescript
const [submitting, setSubmitting] = useState(false);
const [submitted, setSubmitted] = useState(false);
const [error, setError] = useState<string | null>(null);

const handleSubmit = async (data: ContactFormulier) => {
  setSubmitting(true);
  setError(null);

  try {
    await apiClient.submitContactForm(data);
    setSubmitted(true);
  } catch (err) {
    setError(err.message);
  } finally {
    setSubmitting(false);
  }
};

if (submitted) {
  return <SuccessMessage />;
}

return (
  <form onSubmit={handleSubmit}>
    {/* form fields */}
    {error && <ErrorAlert message={error} />}
    <button disabled={submitting}>
      {submitting ? 'Sending...' : 'Send'}
    </button>
  </form>
);
```

### Pagination

```typescript
const [page, setPage] = useState(0);
const [emails, setEmails] = useState<MailResponse[]>([]);
const [total, setTotal] = useState(0);
const limit = 20;

const fetchPage = async (pageNum: number) => {
  const result = await apiClient.listEmails(limit, pageNum * limit);
  setEmails(result.emails);
  setTotal(result.totalCount);
};

useEffect(() => {
  fetchPage(page);
}, [page]);

const totalPages = Math.ceil(total / limit);
```

---

## ⚡ Performance Best Practices

### 1. Image Optimization (Albums)
```typescript
// Use thumbnails for grids
<img src={photo.thumbnail_url || photo.url} loading="lazy" />

// Use full size for lightbox/modal
<img src={photo.url} alt={photo.title} />
```

### 2. Video Lazy Loading
```typescript
// Show thumbnail first, load iframe on click
const [playing, setPlaying] = useState(false);

{!playing ? (
  <img 
    src={`https://cdn-cf-east.streamable.com/image/${video.video_id}.jpg`}
    onClick={() => setPlaying(true)}
  />
) : (
  <iframe src={video.url} allowFullScreen />
)}
```

### 3. Caching Strategy
```typescript
// SWR for data fetching
import useSWR from 'swr';

const { data: albums } = useSWR('/api/albums?include_covers=true', fetcher, {
  revalidateOnFocus: false,
  revalidateOnReconnect: false,
  refreshInterval: 0, // Albums don't change often
});

const { data: maintenanceStatus } = useSWR('/api/under-construction/active', fetcher, {
  refreshInterval: 60000, // Check every minute
  shouldRetryOnError: false,
});
```

### 4. Request Deduplication
```typescript
// React Query
import { useQuery } from '@tanstack/react-query';

const useAlbums = () => {
  return useQuery({
    queryKey: ['albums'],
    queryFn: () => apiClient.getAlbums(true),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};
```

---

## 🔒 Security Checklist

- ✅ **HTTPS only in production**
- ✅ **Sanitize all user input** (XSS prevention)
- ✅ **Validate on both client and server**
- ✅ **Store tokens securely** (httpOnly cookies if possible)
- ✅ **Implement CSRF protection**
- ✅ **Rate limit awareness** (429 handling)
- ✅ **Never expose API keys** client-side
- ✅ **Use environment variables** for sensitive config

---

## 📱 Responsive Design Guidelines

### Breakpoints
```css
/* Mobile First */
.container {
  padding: 1rem;
}

/* Tablet */
@media (min-width: 768px) {
  .container {
    padding: 2rem;
  }
}

/* Desktop */
@media (min-width: 1024px) {
  .container {
    padding: 3rem;
  }
}
```

### Grid Layouts
```css
/* Albums/Videos Grid */
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 2rem;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
}
```

---

## 🆘 Error Handling Strategy

```typescript
// Centralized error handler
class ApiError extends Error {
  constructor(
    public message: string,
    public status?: number,
    public code?: string
  ) {
    super(message);
  }
}

async function handleApiCall<T>(
  apiCall: () => Promise<T>,
  errorMessage = 'Er is een fout opgetreden'
): Promise<T> {
  try {
    return await apiCall();
  } catch (error) {
    if (error instanceof ApiError) {
      // Log to monitoring service
      console.error('API Error:', error);
      
      // Show user-friendly message
      if (error.status === 429) {
        throw new Error('Te veel verzoeken. Probeer het later opnieuw.');
      }
      if (error.status === 401) {
        throw new Error('Niet geautoriseerd. Log opnieuw in.');
      }
      if (error.status >= 500) {
        throw new Error('Server fout. Probeer het later opnieuw.');
      }
    }
    
    throw new Error(errorMessage);
  }
}
```

---

## 📖 Meer Documentatie

- [FRONTEND_QUICKSTART.md](FRONTEND_QUICKSTART.md) - Snelstart gids
- [AUTH_AND_RBAC.md](AUTH_AND_RBAC.md) - Authenticatie & permissies
- [DATABASE_REFERENCE.md](DATABASE_REFERENCE.md) - Database schema

---

## 🎓 Volgende Stappen

1. **Kies je framework** (React, Next.js, Vue, etc.)
2. **Lees relevante documentatie** voor jouw use case
3. **Test endpoints** met de test scripts
4. **Implementeer met voorbeeldcode** uit de docs
5. **Deploy en monitor** met error tracking

---

## ❓ Support

Voor vragen of problemen:
- Check de specifieke module documentatie
- Bekijk de handler code in [`handlers/`](../handlers/)
- Test met de PowerShell scripts
- Contact het backend team

---

## 🎉 Success!

Je hebt nu toegang tot:
- ✅ **4 complete API modules** (Albums, Videos, Maintenance, Email)
- ✅ **15+ endpoints** (public + admin)
- ✅ **Complete TypeScript types**
- ✅ **React component voorbeelden**
- ✅ **Next.js integratie**
- ✅ **Test scripts voor alle modules**
- ✅ **Production-tested endpoints**

Happy coding! 🚀