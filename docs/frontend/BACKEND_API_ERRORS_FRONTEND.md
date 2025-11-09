# Backend API Errors - Frontend Gids

**Document Versie:** 1.0  
**Datum:** 2025-01-08  
**Voor:** Frontend Developers  

---

## 🎯 Quick Summary

**2 Backend API Issues geanalyseerd:**

| Endpoint | Status | Actie Vereist | Prioriteit |
|----------|--------|---------------|-----------|
| `GET /api/participant/emails` | 🔴 **Bug** | Wacht op backend fix | HIGH |
| `GET /api/under-construction/active` | 🟢 **Geen bug** | Frontend aanpassen | LOW |

---

## 1️⃣ Participant Emails - Backend Bug

### Probleem

```
GET http://localhost:8080/api/participant/emails
Status: 500 Internal Server Error
```

**Waar gebruikt:**
1. [`EmailInbox.tsx:164`](../../src/features/email/components/EmailInbox.tsx:164) - Email suggesties laden
2. [`adminEmailService.ts:625`](../../src/features/email/adminEmailService.ts:625) - Fallback in `fetchAanmeldingenEmails()`
3. Admin email compose - TO field autocomplete

### Root Cause

Route ordering bug in backend. De route wordt geregistreerd NÁ de generieke `/:id` route, waardoor het nooit wordt bereikt.

### ✋ Wat moet JIJ doen?

**NIETS!** Dit is een backend issue. Wacht op backend fix.

**Verwachte fix:** Backend verplaatst route naar juiste volgorde.

### Wanneer gefixed is

Na backend deployment test je:

```typescript
// Test in adminEmailService.ts
const emails = await getParticipantEmails();
console.log('Emails loaded:', emails.participant_emails.length);
```

**Verwachte response:**
```json
{
  "participant_emails": ["email1@example.com", "email2@example.com"],
  "system_emails": ["info@dekoninklijkeloop.nl", "inschrijving@dekoninklijkeloop.nl"],
  "all_emails": ["email1@...", "email2@...", "info@...", "inschrijving@..."],
  "counts": {
    "participants": 2,
    "system": 2,
    "total": 4
  }
}
```

### Tijdelijke Workaround (Optioneel)

Als je niet kunt wachten, voeg tijdelijke fallback toe:

```typescript
// In adminEmailService.ts
export async function getParticipantEmails() {
  try {
    const response = await api.get('/participant/emails');
    return response.data;
  } catch (error) {
    console.error('Failed to fetch participant emails, using fallback:', error);
    
    // Tijdelijke fallback: haal alle participants op en extract emails
    try {
      const participants = await api.get('/participant?limit=1000');
      const participantEmails = participants.data.map(p => p.email);
      
      return {
        participant_emails: participantEmails,
        system_emails: ['info@dekoninklijkeloop.nl', 'inschrijving@dekoninklijkeloop.nl'],
        all_emails: [...participantEmails, 'info@dekoninklijkeloop.nl', 'inschrijving@dekoninklijkeloop.nl'],
        counts: {
          participants: participantEmails.length,
          system: 2,
          total: participantEmails.length + 2
        }
      };
    } catch (fallbackError) {
      console.error('Fallback also failed:', fallbackError);
      throw error;
    }
  }
}
```

⚠️ **LET OP:** Verwijder deze workaround zodra backend gefixed is!

---

## 2️⃣ Under Construction - Geen Bug

### Probleem

```
GET http://localhost:8080/api/under-construction/active
Status: 404 Not Found
```

**Waar gebruikt:**
- [`UnderConstructionForm.tsx:39`](../../src/features/settings/UnderConstructionForm.tsx:39) (of soortgelijk)
- Settings pagina maintenance mode feature

### Root Cause

**DIT IS GEEN BUG!** 

Backend returnt 404 wanneer:
- Er is **geen** actieve "under construction" mode
- Website is **niet** in maintenance

Dit is **correct gedrag by design**.

### ✅ Wat moet JIJ doen?

**FIX DE FRONTEND ERROR HANDLING!**

#### Huidige Code (Problematisch)

```typescript
// ❌ SLECHT - Logt elke 404 als error
const fetchUnderConstruction = async () => {
  try {
    const response = await api.get('/under-construction/active');
    setUnderConstruction(response.data);
  } catch (error) {
    console.error('Failed to fetch under construction:', error); // ← SLECHT!
    setError('Kon maintenance mode niet ophalen');
  }
};
```

**Probleem:**
- 404 wordt gelogd als ERROR (misleidend!)
- Gebruiker ziet foutmelding terwijl alles werkt
- Logs vervuild met "fake" errors

#### Gecorrigeerde Code (Correct)

```typescript
// ✅ GOED - Handle 404 als normale state
const fetchUnderConstruction = async () => {
  try {
    const response = await api.get('/under-construction/active');
    
    // 200 OK: Maintenance mode is actief
    setUnderConstruction(response.data);
    setIsActive(true);
    
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 404) {
      // 404: Maintenance mode staat UIT - dit is NORMAAL
      setUnderConstruction(null);
      setIsActive(false);
      // Geen error logging!
      return;
    }
    
    // Alleen log andere errors (500, network errors, etc.)
    console.error('Server error fetching under construction:', error);
    setError('Kon maintenance mode niet ophalen');
  }
};
```

#### In TypeScript met Types

```typescript
interface UnderConstructionResponse {
  id: number;
  is_active: boolean;
  title: string;
  message: string;
  footer_text?: string;
  logo_url?: string;
  expected_date?: string;
  social_links?: any[];
  progress_percentage?: number;
  contact_email?: string;
  newsletter_enabled?: boolean;
  created_at: string;
  updated_at: string;
}

interface UnderConstructionState {
  isActive: boolean;
  data: UnderConstructionResponse | null;
  loading: boolean;
  error: string | null;
}

const fetchUnderConstruction = async (): Promise<void> => {
  setLoading(true);
  setError(null);
  
  try {
    const response = await api.get<UnderConstructionResponse>(
      '/under-construction/active'
    );
    
    setState({
      isActive: true,
      data: response.data,
      loading: false,
      error: null
    });
    
  } catch (error) {
    if (axios.isAxiosError(error)) {
      // 404 = Maintenance mode OFF (normal)
      if (error.response?.status === 404) {
        setState({
          isActive: false,
          data: null,
          loading: false,
          error: null
        });
        return;
      }
      
      // Other errors
      setState({
        isActive: false,
        data: null,
        loading: false,
        error: `Server error: ${error.response?.status || 'unknown'}`
      });
    } else {
      setState({
        isActive: false,
        data: null,
        loading: false,
        error: 'Network error'
      });
    }
  }
};
```

### UI Verbeteringen (Optioneel)

In plaats van error te tonen bij 404, toon duidelijk bericht:

```typescript
// In je component
{loading ? (
  <Spinner />
) : error ? (
  <Alert variant="error">{error}</Alert>
) : !isActive ? (
  <Alert variant="info">
    ℹ️ Maintenance mode staat momenteel UIT. 
    De website is normaal bereikbaar voor bezoekers.
  </Alert>
) : (
  <UnderConstructionForm data={data} onSave={handleSave} />
)}
```

---

## 📋 Action Items Checklist

### Voor Jou (Frontend Developer)

#### Immediate (Doen NU)
- [ ] **Fix Under Construction Error Handling**
  - [ ] Open `underConstructionClient.ts` (of vergelijkbaar bestand)
  - [ ] Update error handling om 404 als normale state te behandelen
  - [ ] Verwijder error logging voor 404 responses
  - [ ] Test in dev environment
  - [ ] Commit met message: `fix: handle 404 as normal state for under-construction endpoint`

#### Waiting on Backend Fix
- [ ] **Test Participant Emails na Backend Deploy**
  - [ ] Wacht op notificatie van backend team over fix deployment
  - [ ] Test `/api/participant/emails` endpoint
  - [ ] Verify email suggesties werken in admin panel
  - [ ] Verify fallback mechanisme werkt in email inbox
  - [ ] Verwijder tijdelijke workaround (als gebruikt)

#### Optional Improvements
- [ ] Voeg loading state toe voor under construction fetch
- [ ] Voeg info message toe wanneer maintenance mode UIT staat
- [ ] Voeg retry logic toe voor participant emails (indien backend traag is)

---

## 🧪 Testing Scripts

### Test 1: Under Construction (na fix)

```typescript
// In browser console tijdens dev
const testUnderConstruction = async () => {
  console.log('Testing under-construction endpoint...');
  
  try {
    const response = await fetch('http://localhost:8080/api/under-construction/active');
    
    if (response.status === 404) {
      console.log('✅ PASS: 404 returned (maintenance OFF)');
      const data = await response.json();
      console.log('Response:', data);
    } else if (response.status === 200) {
      console.log('✅ PASS: 200 returned (maintenance ON)');
      const data = await response.json();
      console.log('Active maintenance:', data);
    } else {
      console.error('❌ FAIL: Unexpected status:', response.status);
    }
  } catch (error) {
    console.error('❌ FAIL: Network error:', error);
  }
};

testUnderConstruction();
```

### Test 2: Participant Emails (na backend fix)

```typescript
// In browser console met auth token
const testParticipantEmails = async () => {
  console.log('Testing participant emails endpoint...');
  
  const token = localStorage.getItem('token'); // of waar je token opslaat
  
  try {
    const response = await fetch('http://localhost:8080/api/participant/emails', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.status === 200) {
      const data = await response.json();
      console.log('✅ PASS: Emails loaded successfully');
      console.log('Participant emails:', data.participant_emails?.length);
      console.log('System emails:', data.system_emails?.length);
      console.log('Total emails:', data.all_emails?.length);
      console.log('Full response:', data);
    } else {
      console.error('❌ FAIL: Status:', response.status);
      const error = await response.json();
      console.error('Error:', error);
    }
  } catch (error) {
    console.error('❌ FAIL: Network error:', error);
  }
};

testParticipantEmails();
```

### Test 3: Email Inbox Flow

```typescript
// Manual test in admin panel
// 1. Login als admin
// 2. Navigeer naar Email Inbox
// 3. Open console (F12)
// 4. Klik "Aanmeldingen" tab
// 5. Observeer console logs

// Verwacht:
// - Geen errors in console
// - Emails laden succesvol
// - Participant suggesties werken in compose form
```

---

## 🔧 API Client Updates

Als je een centrale API client gebruikt, update deze ook:

```typescript
// api/client/underConstructionClient.ts

export const underConstructionClient = {
  /**
   * Get active under construction record
   * Returns null if no maintenance mode is active (404)
   */
  async getActive(): Promise<UnderConstructionResponse | null> {
    try {
      const response = await api.get<UnderConstructionResponse>(
        '/under-construction/active'
      );
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 404) {
        // No active maintenance - this is normal
        return null;
      }
      // Re-throw other errors
      throw error;
    }
  },
  
  // Other methods...
};

// api/client/participantClient.ts

export const participantClient = {
  /**
   * Get all participant emails + system emails
   * Used for email autocomplete in admin panel
   */
  async getEmails(): Promise<ParticipantEmailsResponse> {
    const response = await api.get<ParticipantEmailsResponse>(
      '/participant/emails'
    );
    return response.data;
  },
  
  // Other methods...
};

interface ParticipantEmailsResponse {
  participant_emails: string[];
  system_emails: string[];
  all_emails: string[];
  counts: {
    participants: number;
    system: number;
    total: number;
  };
}
```

---

## 📞 Wie te Contacteren

### Voor Backend Issues
- [Jeffrey] - Backend developer
- Slack: `#backend-team`
- Issue: "Participant emails route ordering bug"

### Voor Frontend Issues  
- [Jouw Team Lead]
- Slack: `#frontend-team`
- Issue: "Under construction error handling"

---

## 📚 Related Documentation

### Backend Docs
- [Backend API Errors Full Analysis](../api/BACKEND_API_ERRORS.md)
- [API Quick Reference](../api/QUICK_REFERENCE.md)
- [Authentication Guide](../api/AUTHENTICATION.md)

### Frontend Docs (Maak aan indien nodig)
- API Client Documentation
- Error Handling Best Practices
- Admin Panel Guide

---

## 📝 Change Log

| Datum | Versie | Auteur | Wijziging |
|-------|--------|--------|-----------|
| 2025-01-08 | 1.0 | Kilo Code | Initiële frontend gids |

---

## 💡 Tips & Best Practices

### Error Handling Strategy

```typescript
// ✅ GOED: Specifiek per HTTP status
catch (error) {
  if (axios.isAxiosError(error)) {
    switch (error.response?.status) {
      case 404:
        // Handle as normal "not found" state
        return handleNotFound();
      case 401:
        // Redirect to login
        return redirectToLogin();
      case 403:
        // Show permission denied
        return showPermissionDenied();
      case 500:
        // Show server error
        return showServerError();
      default:
        return showGenericError();
    }
  }
  // Network error
  return showNetworkError();
}
```

### Loading States

```typescript
// ✅ GOED: Altijd loading state tonen
const [state, setState] = useState({
  data: null,
  loading: true,  // Start with loading
  error: null
});

// In je fetch functie
setState(prev => ({ ...prev, loading: true }));
try {
  const data = await fetchData();
  setState({ data, loading: false, error: null });
} catch (error) {
  setState({ data: null, loading: false, error: error.message });
}
```

### TypeScript Types

```typescript
// ✅ GOED: Type alle API responses
import type { 
  UnderConstructionResponse, 
  ParticipantEmailsResponse,
  ApiError 
} from '@/types/api';

// Gebruik generics voor type safety
async function apiGet<T>(url: string): Promise<T> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new ApiError(response.status, await response.text());
  }
  return response.json();
}
```
