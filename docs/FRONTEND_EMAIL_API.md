# Frontend Email API Documentatie

Deze documentatie beschrijft alle email-gerelateerde endpoints voor gebruik in de frontend.

## Base URLs

- **Development (Docker):** `http://localhost:8082`
- **Production:** `https://dklemailservice.onrender.com`

---

# 📧 PART 1: PUBLIC EMAIL SUBMISSION (Geen Authenticatie)

Deze endpoints zijn beschikbaar voor publiek gebruik en vereisen GEEN authenticatie.

## 1. Contact Formulier Versturen

**Endpoint:** `POST /api/contact-email`

Verstuurt een contact formulier en verzend
t bevestiging emails naar zowel de gebruiker als admin.

**Request:**
```typescript
interface ContactFormulier {
  naam: string;              // Verplicht
  email: string;             // Verplicht
  telefoon?: string;         // Optioneel
  bericht: string;           // Verplicht
  privacy_akkoord: boolean;  // Verplicht: moet true zijn
  test_mode?: boolean;       // Optioneel: Voor testen zonder echte emails
}
```

**Voorbeeld Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/contact-email', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    naam: 'Jan Jansen',
    email: 'jan@example.com',
    telefoon: '06-12345678',
    bericht: 'Ik heb een vraag over het evenement',
    privacy_akkoord: true
  })
});

const result = await response.json();
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Je bericht is verzonden! Je ontvangt ook een bevestiging per email."
}
```

**Response (Test Mode):**
```json
{
  "success": true,
  "message": "[TEST MODE] Je bericht is verwerkt (geen echte email verzonden).",
  "test_mode": true
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "Naam, email en bericht zijn verplicht"
}
```

**Validatie Regels:**
- ✅ `naam` mag niet leeg zijn
- ✅ `email` mag niet leeg zijn en moet geldig format hebben
- ✅ `bericht` mag niet leeg zijn
- ✅ `privacy_akkoord` moet `true` zijn

---

## 2. Aanmelding/Registratie Formulier Versturen

**Endpoint:** `POST /api/aanmelding-email`

Verstuurt een registratie formulier voor deelname aan het evenement.

**Request:**
```typescript
interface AanmeldingFormulier {
  naam: string;              // Verplicht
  email: string;             // Verplicht
  telefoon?: string;         // Optioneel
  rol: string;               // Verplicht: "deelnemer" of "vrijwilliger"
  afstand: string;           // Verplicht: "5km", "10km", "15km"
  ondersteuning?: string;    // Optioneel: Extra ondersteuning behoeften
  bijzonderheden?: string;   // Optioneel: Aanvullende opmerkingen
  terms: boolean;            // Verplicht: moet true zijn
  test_mode?: boolean;       // Optioneel: Voor testen
}
```

**Voorbeeld Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/aanmelding-email', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    naam: 'Maria Smith',
    email: 'maria@example.com',
    telefoon: '06-98765432',
    rol: 'deelnemer',
    afstand: '10km',
    ondersteuning: 'Rolstoel toegankelijkheid',
    bijzonderheden: 'Eerste keer deelnemer',
    terms: true
  })
});

const result = await response.json();
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Je aanmelding is verzonden! Je ontvangt ook een bevestiging per email."
}
```

**Validatie Regels:**
- ✅ `naam` mag niet leeg zijn
- ✅ `email` mag niet leeg zijn en moet geldig format hebben (@, .)
- ✅ `rol` mag niet leeg zijn
- ✅ `afstand` mag niet leeg zijn
- ✅ `terms` moet `true` zijn

---

## 3. Whisky for Charity Order Email

**Endpoint:** `POST /api/wfc/order-email`

Verstuurt order bevestiging emails voor Whisky for Charity bestellingen.

**Authenticatie:** Vereist `X-API-Key` header (niet JWT).

**Request:**
```typescript
interface WFCOrderRequest {
  order_id: string;              // Verplicht: Unieke order ID
  customer_name: string;         // Verplicht
  customer_email: string;        // Verplicht
  customer_address?: string;     // Optioneel
  customer_city?: string;        // Optioneel
  customer_postal?: string;      // Optioneel
  customer_country?: string;     // Optioneel
  total_amount: number;          // Totaal bedrag
  items: WFCOrderItem[];         // Order items
}

interface WFCOrderItem {
  name: string;        // Product naam
  quantity: number;    // Aantal
  price: number;       // Prijs per stuk
  total: number;       // Totaal (quantity * price)
}
```

**Voorbeeld Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/wfc/order-email', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': 'your-api-key-here'
  },
  body: JSON.stringify({
    order_id: 'WFC-2024-001',
    customer_name: 'John Doe',
    customer_email: 'john@example.com',
    customer_address: 'Hoofdstraat 1',
    customer_city: 'Amsterdam',
    customer_postal: '1000AA',
    customer_country: 'Nederland',
    total_amount: 49.99,
    items: [
      {
        name: 'Premium Whisky',
        quantity: 1,
        price: 49.99,
        total: 49.99
      }
    ]
  })
});

const result = await response.json();
```

**Response:**
```json
{
  "success": true,
  "customer_email_sent": true,
  "admin_email_sent": true,
  "order_id": "WFC-2024-001"
}
```

---

# 🔐 PART 2: ADMIN EMAIL MANAGEMENT (Authenticatie Vereist)

Deze endpoints vereisen JWT authenticatie via `Authorization: Bearer <token>` header.

## Contact Formulieren Beheren

### 1. Lijst Ophalen

**Endpoint:** `GET /api/contact`

**Query Parameters:**
- `limit` (number, optioneel): Aantal resultaten (1-100, default 10)
- `offset` (number, optioneel): Offset voor paginering (default 0)

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/contact?limit=20&offset=0',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const contacts = await response.json();
```

**Response:**
```typescript
interface ContactFormulier {
  id: string;
  naam: string;
  email: string;
  telefoon: string | null;
  bericht: string;
  status: string;                    // "nieuw", "in_behandeling", "beantwoord", "gesloten"
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
  antwoorden?: ContactAntwoord[];    // Optioneel bij detail view
}

ContactFormulier[]
```

---

### 2. Specifiek Contact Ophalen

**Endpoint:** `GET /api/contact/:id`

**Request:**
```typescript
const contactId = '123e4567-e89b-12d3-a456-426614174000';
const response = await fetch(
  `https://dklemailservice.onrender.com/api/contact/${contactId}`,
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const contact = await response.json();
```

**Response:** Bevat het contact met alle antwoorden:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "naam": "Jan Jansen",
  "email": "jan@example.com",
  "bericht": "Vraag over evenement",
  "status": "beantwoord",
  "created_at": "2024-11-01T10:00:00Z",
  "antwoorden": [
    {
      "id": "456e7890-e89b-12d3-a456-426614174001",
      "contact_id": "123e4567-e89b-12d3-a456-426614174000",
      "tekst": "Bedankt voor je vraag...",
      "verzond_door": "admin@dekoninklijkeloop.nl",
      "email_verzonden": true,
      "created_at": "2024-11-01T11:00:00Z"
    }
  ]
}
```

---

### 3. Contact Status Updaten

**Endpoint:** `PUT /api/contact/:id`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/contact/${contactId}`,
  {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      status: 'in_behandeling',
      notities: 'Behandeling gestart'
    })
  }
);
```

**Response:** Geüpdatet contact object

---

### 4. Antwoord Toevoegen

**Endpoint:** `POST /api/contact/:id/antwoord`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/contact/${contactId}/antwoord`,
  {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      tekst: 'Bedankt voor je vraag. Hier is het antwoord...'
    })
  }
);
```

**Wat gebeurt er:**
1. Antwoord wordt opgeslagen in database
2. Email wordt automatisch naar gebruiker verstuurd
3. Status wordt geüpdatet naar "beantwoord"

---

### 5. Filter Op Status

**Endpoint:** `GET /api/contact/status/:status`

**Valid Statuses:** `nieuw`, `in_behandeling`, `beantwoord`, `gesloten`

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/contact/status/nieuw',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const newContacts = await response.json();
```

---

### 6. Contact Verwijderen

**Endpoint:** `DELETE /api/contact/:id`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/contact/${contactId}`,
  {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
```

**Response:**
```json
{
  "success": true,
  "message": "Contactformulier succesvol verwijderd"
}
```

---

## Aanmeldingen/Registraties Beheren

### 1. Lijst Ophalen

**Endpoint:** `GET /api/aanmelding`

**Query Parameters:**
- `limit` (number, optioneel): Aantal resultaten (1-100, default 10)
- `offset` (number, optioneel): Offset voor paginering (default 0)

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/aanmelding?limit=20&offset=0',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const aanmeldingen = await response.json();
```

**Response:**
```typescript
interface Aanmelding {
  id: string;
  naam: string;
  email: string;
  telefoon: string | null;
  rol: string;                       // "deelnemer" of "vrijwilliger"
  afstand: string;                   // "5km", "10km", "15km"
  ondersteuning: string | null;
  bijzonderheden: string | null;
  status: string;                    // "nieuw", "beantwoord", etc.
  terms: boolean;
  test_mode: boolean;
  email_verzonden: boolean;
  email_verzonden_op: string | null;
  behandeld_door: string | null;
  behandeld_op: string | null;
  notities: string | null;
  created_at: string;
  updated_at: string;
  gebruiker_id: string | null;      // Link naar gebruiker account
  antwoorden?: AanmeldingAntwoord[];
}

Aanmelding[]
```

---

### 2. Specifieke Aanmelding Ophalen

**Endpoint:** `GET /api/aanmelding/:id`

Haalt één aanmelding op inclusief alle antwoorden.

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/aanmelding/${aanmeldingId}`,
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const aanmelding = await response.json();
```

---

### 3. Filter Op Rol

**Endpoint:** `GET /api/aanmelding/rol/:rol`

**Valid Roles:** `deelnemer`, `vrijwilliger`

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/aanmelding/rol/deelnemer',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const deelnemers = await response.json();
```

---

### 4. Aanmelding Updaten

**Endpoint:** `PUT /api/aanmelding/:id`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/aanmelding/${aanmeldingId}`,
  {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      status: 'beantwoord',
      notities: 'Aanmelding verwerkt'
    })
  }
);
```

---

### 5. Antwoord Toevoegen

**Endpoint:** `POST /api/aanmelding/:id/antwoord`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/aanmelding/${aanmeldingId}/antwoord`,
  {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      tekst: 'Bedankt voor je aanmelding! Je bent ingeschreven voor 10km...'
    })
  }
);
```

---

### 6. Aanmelding Verwijderen

**Endpoint:** `DELETE /api/aanmelding/:id`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/aanmelding/${aanmeldingId}`,
  {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
```

---

# 📨 PART 3: INCOMING EMAIL MANAGEMENT

Beheer inkomende emails die automatisch worden opgehaald van mail accounts.

## 1. Lijst Emails Ophalen (Paginated)

**Endpoint:** `GET /api/mail`

**Query Parameters:**
- `limit` (number, 1-100): Aantal emails (default 10)
- `offset` (number): Offset voor paginering (default 0)

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/mail?limit=20&offset=0',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const result = await response.json();
```

**Response:**
```typescript
interface PaginatedMailResponse {
  emails: MailResponse[];
  totalCount: number;
}

interface MailResponse {
  id: string;
  message_id: string;
  sender: string;                    // From address
  to: string;                        // To address
  subject: string;
  html: string;                      // Email body (HTML)
  content_type: string;
  received_at: string;               // ISO timestamp
  uid: string;
  account_type: string;              // "info" of "inschrijving"
  read: boolean;                     // Is processed
  processed_at: string | null;
  created_at: string;
  updated_at: string;
}
```

---

## 2. Specifieke Email Ophalen

**Endpoint:** `GET /api/mail/:id`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/mail/${emailId}`,
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const email = await response.json();
```

**Response:** Single `MailResponse` object met volledige HTML body.

---

## 3. Onverwerkte Emails Ophalen

**Endpoint:** `GET /api/mail/unprocessed`

Haalt alle nog niet verwerkte emails op.

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/mail/unprocessed',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const unprocessed = await response.json();
```

---

## 4. Filter Op Account Type

**Endpoint:** `GET /api/mail/account/:type`

**Valid Types:** `info`, `inschrijving`

**Query Parameters:**
- `limit` (number): Aantal emails (default 10)
- `offset` (number): Offset

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/mail/account/info?limit=50',
  {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const infoEmails = await response.json();
```

---

## 5. Markeer Als Verwerkt

**Endpoint:** `PUT /api/mail/:id/processed`

Markeert een email als verwerkt/gelezen.

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/mail/${emailId}/processed`,
  {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
```

**Response:**
```json
{
  "message": "Email 123e4567-e89b-12d3-a456-426614174000 gemarkeerd als verwerkt"
}
```

---

## 6. Email Verwijderen

**Endpoint:** `DELETE /api/mail/:id`

**Request:**
```typescript
const response = await fetch(
  `https://dklemailservice.onrender.com/api/mail/${emailId}`,
  {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
```

---

## 7. Nieuwe Emails Ophalen (Manual Fetch)

**Endpoint:** `POST /api/mail/fetch`

Manueel ophalen van nieuwe emails van de mailserver.

**Request:**
```typescript
const response = await fetch(
  'https://dklemailservice.onrender.com/api/mail/fetch',
  {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  }
);
const result = await response.json();
```

**Response:**
```json
{
  "message": "5 emails opgehaald, 5 succesvol opgeslagen",
  "fetchTime": "2024-11-01T22:00:00Z"
}
```

---

# 🎨 React Voorbeelden

## Contact Formulier Component

```typescript
import React, { useState } from 'react';

interface ContactFormData {
  naam: string;
  email: string;
  telefoon: string;
  bericht: string;
  privacy_akkoord: boolean;
}

const ContactForm: React.FC = () => {
  const [formData, setFormData] = useState<ContactFormData>({
    naam: '',
    email: '',
    telefoon: '',
    bericht: '',
    privacy_akkoord: false
  });
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        'https://dklemailservice.onrender.com/api/contact-email',
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(formData)
        }
      );

      const result = await response.json();

      if (response.ok && result.success) {
        setSuccess(true);
        // Reset form
        setFormData({
          naam: '',
          email: '',
          telefoon: '',
          bericht: '',
          privacy_akkoord: false
        });
      } else {
        setError(result.error || 'Er is een fout opgetreden');
      }
    } catch (err) {
      setError('Er is een fout opgetreden bij het versturen');
      console.error('Contact form error:', err);
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="success-message">
        <h2>Bedankt!</h2>
        <p>Je bericht is verzonden. Je ontvangt een bevestiging per email.</p>
        <button onClick={() => setSuccess(false)}>Nieuw bericht</button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="naam">Naam *</label>
        <input
          id="naam"
          type="text"
          value={formData.naam}
          onChange={(e) => setFormData({...formData, naam: e.target.value})}
          required
        />
      </div>

      <div>
        <label htmlFor="email">Email *</label>
        <input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({...formData, email: e.target.value})}
          required
        />
      </div>

      <div>
        <label htmlFor="telefoon">Telefoon</label>
        <input
          id="telefoon"
          type="tel"
          value={formData.telefoon}
          onChange={(e) => setFormData({...formData, telefoon: e.target.value})}
        />
      </div>

      <div>
        <label htmlFor="bericht">Bericht *</label>
        <textarea
          id="bericht"
          value={formData.bericht}
          onChange={(e) => setFormData({...formData, bericht: e.target.value})}
          rows={5}
          required
        />
      </div>

      <div>
        <label>
          <input
            type="checkbox"
            checked={formData.privacy_akkoord}
            onChange={(e) => setFormData({...formData, privacy_akkoord: e.target.checked})}
            required
          />
          Ik ga akkoord met het privacybeleid *
        </label>
      </div>

      {error && <div className="error-message">{error}</div>}

      <button type="submit" disabled={loading}>
        {loading ? 'Versturen...' : 'Verstuur'}
      </button>
    </form>
  );
};

export default ContactForm;
```

---

## Aanmelding Formulier Component

```typescript
import React, { useState } from 'react';

interface AanmeldingFormData {
  naam: string;
  email: string;
  telefoon: string;
  rol: string;
  afstand: string;
  ondersteuning: string;
  bijzonderheden: string;
  terms: boolean;
}

const AanmeldingForm: React.FC = () => {
  const [formData, setFormData] = useState<AanmeldingFormData>({
    naam: '',
    email: '',
    telefoon: '',
    rol: 'deelnemer',
    afstand: '10km',
    ondersteuning: '',
    bijzonderheden: '',
    terms: false
  });
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        'https://dklemailservice.onrender.com/api/aanmelding-email',
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(formData)
        }
      );

      const result = await response.json();

      if (response.ok && result.success) {
        setSuccess(true);
      } else {
        setError(result.error || 'Er is een fout opgetreden');
      }
    } catch (err) {
      setError('Er is een fout opgetreden bij het versturen');
      console.error('Aanmelding error:', err);
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="success-message">
        <h2>Aanmelding Geslaagd!</h2>
        <p>Je aanmelding is verzonden. Je ontvangt een bevestiging per email.</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Naam *</label>
        <input
          type="text"
          value={formData.naam}
          onChange={(e) => setFormData({...formData, naam: e.target.value})}
          required
        />
      </div>

      <div>
        <label>Email *</label>
        <input
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({...formData, email: e.target.value})}
          required
        />
      </div>

      <div>
        <label>Telefoon</label>
        <input
          type="tel"
          value={formData.telefoon}
          onChange={(e) => setFormData({...formData, telefoon: e.target.value})}
        />
      </div>

      <div>
        <label>Ik meld me aan als *</label>
        <select
          value={formData.rol}
          onChange={(e) => setFormData({...formData, rol: e.target.value})}
          required
        >
          <option value="deelnemer">Deelnemer</option>
          <option value="vrijwilliger">Vrijwilliger</option>
        </select>
      </div>

      <div>
        <label>Afstand *</label>
        <select
          value={formData.afstand}
          onChange={(e) => setFormData({...formData, afstand: e.target.value})}
          required
        >
          <option value="5km">5 km</option>
          <option value="10km">10 km</option>
          <option value="15km">15 km</option>
        </select>
      </div>

      <div>
        <label>Heb je extra ondersteuning nodig?</label>
        <input
          type="text"
          value={formData.ondersteuning}
          onChange={(e) => setFormData({...formData, ondersteuning: e.target.value})}
          placeholder="Bijv. rolstoel toegankelijkheid"
        />
      </div>

      <div>
        <label>Bijzonderheden</label>
        <textarea
          value={formData.bijzonderheden}
          onChange={(e) => setFormData({...formData, bijzonderheden: e.target.value})}
          rows={3}
          placeholder="Eventuele opmerkingen"
        />
      </div>

      <div>
        <label>
          <input
            type="checkbox"
            checked={formData.terms}
            onChange={(e) => setFormData({...formData, terms: e.target.checked})}
            required
          />
          Ik ga akkoord met de voorwaarden *
        </label>
      </div>

      {error && <div className="error">{error}</div>}

      <button type="submit" disabled={loading}>
        {loading ? 'Aanmelden...' : 'Aanmelden'}
      </button>
    </form>
  );
};

export default AanmeldingForm;
```

---

## Admin Email Dashboard Component

```typescript
import React, { useEffect, useState } from 'react';

interface Email {
  id: string;
  sender: string;
  subject: string;
  received_at: string;
  read: boolean;
  account_type: string;
}

const EmailDashboard: React.FC = () => {
  const [emails, setEmails] = useState<Email[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(0);
  const limit = 20;

  const token = localStorage.getItem('auth_token');

  useEffect(() => {
    fetchEmails();
  }, [page]);

  const fetchEmails = async () => {
    setLoading(true);
    try {
      const response = await fetch(
        `https://dklemailservice.onrender.com/api/mail?limit=${limit}&offset=${page * limit}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      );

      const data = await response.json();
      setEmails(data.emails);
      setTotalCount(data.totalCount);
    } catch (error) {
      console.error('Error fetching emails:', error);
    } finally {
      setLoading(false);
    }
  };

  const markAsRead = async (emailId: string) => {
    try {
      await fetch(
        `https://dklemailservice.onrender.com/api/mail/${emailId}/processed`,
        {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      );
      // Refresh list
      fetchEmails();
    } catch (error) {
      console.error('Error marking email as read:', error);
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div className="email-dashboard">
      <h2>Inbox ({totalCount} emails)</h2>
      
      <div className="email-list">
        {emails.map((email) => (
          <div 
            key={email.id} 
            className={`email-item ${email.read ? 'read' : 'unread'}`}
          >
            <div className="email-header">
              <strong>{email.sender}</strong>
              <span className="account-type">{email.account_type}</span>
            </div>
            <div className="email-subject">{email.subject}</div>
            <div className="email-date">
              {new Date(email.received_at).toLocaleString('nl-NL')}
            </div>
            {!email.read && (
              <button onClick={() => markAsRead(email.id)}>
                Markeer als gelezen
              </button>
            )}
          </div>
        ))}
      </div>

      <div className="pagination">
        <button 
          onClick={() => setPage(p => Math.max(0, p - 1))}
          disabled={page === 0}
        >
          Vorige
        </button>
        <span>Pagina {page + 1}</span>
        <button 
          onClick={() => setPage(p => p + 1)}
          disabled={(page + 1) * limit >= totalCount}
        >
          Volgende
        </button>
      </div>
    </div>
  );
};

export default EmailDashboard;
```

---

## Admin Contact Management Component

```typescript
import React, { useEffect, useState } from 'react';

interface Contact {
  id: string;
  naam: string;
  email: string;
  bericht: string;
  status: string;
  created_at: string;
}

const ContactManagement: React.FC = () => {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [filter, setFilter] = useState<string>('nieuw');
  const token = localStorage.getItem('auth_token');

  useEffect(() => {
    fetchContacts();
  }, [filter]);

  const fetchContacts = async () => {
    try {
      const url = filter === 'all' 
        ? 'https://dklemailservice.onrender.com/api/contact'
        : `https://dklemailservice.onrender.com/api/contact/status/${filter}`;

      const response = await fetch(url, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      const data = await response.json();
      setContacts(data);
    } catch (error) {
      console.error('Error fetching contacts:', error);
    }
  };

  const updateStatus = async (contactId: string, newStatus: string) => {
    try {
      await fetch(
        `https://dklemailservice.onrender.com/api/contact/${contactId}`,
        {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ status: newStatus })
        }
      );
      fetchContacts();
    } catch (error) {
      console.error('Error updating contact:', error);
    }
  };

  return (
    <div className="contact-management">
      <h2>Contact Formulieren</h2>
      
      <div className="filters">
        <button onClick={() => setFilter('all')}>Alle</button>
        <button onClick={() => setFilter('nieuw')}>Nieuw</button>
        <button onClick={() => setFilter('in_behandeling')}>In Behandeling</button>
        <button onClick={() => setFilter('beantwoord')}>Beantwoord</button>
        <button onClick={() => setFilter('gesloten')}>Gesloten</button>
      </div>

      <table>
        <thead>
          <tr>
            <th>Naam</th>
            <th>Email</th>
            <th>Bericht</th>
            <th>Status</th>
            <th>Datum</th>
            <th>Acties</th>
          </tr>
        </thead>
        <tbody>
          {contacts.map((contact) => (
            <tr key={contact.id}>
              <td>{contact.naam}</td>
              <td>{contact.email}</td>
              <td>{contact.bericht.substring(0, 50)}...</td>
              <td>{contact.status}</td>
              <td>{new Date(contact.created_at).toLocaleDateString('nl-NL')}</td>
              <td>
                <select
                  value={contact.status}
                  onChange={(e) => updateStatus(contact.id, e.target.value)}
                >
                  <option value="nieuw">Nieuw</option>
                  <option value="in_behandeling">In Behandeling</option>
                  <option value="beantwoord">Beantwoord</option>
                  <option value="gesloten">Gesloten</option>
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ContactManagement;
```

---

# 📋 TypeScript Types

```typescript
// types/email.ts

// Public Submission Types
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

export interface WFCOrderRequest {
  order_id: string;
  customer_name: string;
  customer_email: string;
  customer_address?: string;
  customer_city?: string;
  customer_postal?: string;
  customer_country?: string;
  total_amount: number;
  items: WFCOrderItem[];
}

export interface WFCOrderItem {
  name: string;
  quantity: number;
  price: number;
  total: number;
}

// Admin Management Types
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

// API Response Types
export interface EmailSubmissionResponse {
  success: boolean;
  message: string;
  test_mode?: boolean;
  error?: string;
}

export interface WFCOrderResponse {
  success: boolean;
  customer_email_sent: boolean;
  admin_email_sent: boolean;
  order_id: string;
}
```

---

# 🔑 Belangrijke Opmerkingen

## Test Mode
Alle publieke submission endpoints ondersteunen test mode:

```typescript
// Via header
fetch(url, {
  headers: {
    'X-Test-Mode': 'true'
  }
});

// Of via body
fetch(url, {
  body: JSON.stringify({
    ...formData,
    test_mode: true
  })
});
```

**In test mode:**
- ✅ Validatie wordt uitgevoerd
- ✅ Response wordt gegeven
- ❌ Geen emails verzonden
- ❌ Geen data opgeslagen in database (contact/aanmelding)

## Email Templates
Het systeem gebruikt HTML email templates:
- [`contact_email.html`](../templates/contact_email.html) - Voor contact emails
- [`aanmelding_email.html`](../templates/aanmelding_email.html) - Voor aanmelding emails
- [`wfc_order_confirmation.html`](../templates/wfc_order_confirmation.html) - Voor WFC orders
- [`wfc_order_admin.html`](../templates/wfc_order_admin.html) - Voor WFC admin notificaties

## Rate Limiting
- Contact/Aanmelding endpoints hebben rate limiting
- Bij te veel requests krijg je een 429 Too Many Requests response
- Implementeer exponential backoff bij retries

## CORS
Toegestane origins zijn geconfigureerd in de backend:
- `https://www.dekoninklijkeloop.nl`
- `https://dekoninklijkeloop.nl`
- `https://admin.dekoninklijkeloop.nl`
- `http://localhost:3000`
- `http://localhost:5173`

## Authenticatie voor Admin Endpoints
Alle `/api/contact`, `/api/aanmelding`, en `/api/mail` endpoints vereisen:
1. JWT token via `Authorization: Bearer <token>` header
2. Token krijg je via `/api/auth/login` endpoint
3. Zie [`FRONTEND_QUICKSTART.md`](FRONTEND_QUICKSTART.md) voor auth details

## WFC API Key
Voor WFC orders gebruik je een speciale API key (niet JWT):
```typescript
headers: {
  'X-API-Key': process.env.WFC_API_KEY
}
```

---

# ⚡ Performance Tips

1. **Debounce form submissions** om dubbele verzendingen te voorkomen
2. **Show loading states** tijdens email verzending
3. **Implement retry logic** met exponential backoff
4. **Cache admin data** met SWR of React Query
5. **Validate client-side** voor betere UX (maar server validate is leidend)
6. **Use pagination** voor grote lijsten

---

# 🧪 Testing Voorbeelden

## Jest Test voor Contact Form

```typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import ContactForm from './ContactForm';

describe('ContactForm', () => {
  it('should submit contact form successfully', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          success: true,
          message: 'Email verzonden'
        })
      })
    );

    render(<ContactForm />);
    
    fireEvent.change(screen.getByLabelText(/naam/i), {
      target: { value: 'Test User' }
    });
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'test@example.com' }
    });
    fireEvent.change(screen.getByLabelText(/bericht/i), {
      target: { value: 'Test message' }
    });
    fireEvent.click(screen.getByLabelText(/privacy/i));
    
    fireEvent.click(screen.getByText(/verstuur/i));

    await waitFor(() => {
      expect(screen.getByText(/bedankt/i)).toBeInTheDocument();
    });
  });

  it('should show error for invalid email', async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: false,
        json: () => Promise.resolve({
          success: false,
          error: 'Ongeldig email adres'
        })
      })
    );

    render(<ContactForm />);
    
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'invalid-email' }
    });
    fireEvent.click(screen.getByText(/verstuur/i));

    await waitFor(() => {
      expect(screen.getByText(/ongeldig email/i)).toBeInTheDocument();
    });
  });
});
```

---

# 📊 Production Status & Data

## Email Accounts (Auto-Fetch)
- **info@dekoninklijkeloop.nl** - Algemene vragen (account_type: "info")
- **inschrijving@dekoninklijkeloop.nl** - Aanmeldingen (account_type: "inschrijving")

## Auto-Fetch Interval
Emails worden automatisch elke 15 minuten opgehaald (configureerbaar).

## Email Processing Flow
1. User vult formulier in op website
2. Frontend POST naar `/api/contact-email` of `/api/aanmelding-email`
3. Backend opslaat in database en stuurt 2 emails:
   - Bevestiging naar gebruiker
   - Notificatie naar admin email
4. Admin kan reageren via admin panel
5. Antwoord wordt automatisch per email verstuurd

---

# 🆘 Error Handling

```typescript
async function submitContactForm(data: ContactFormulier): Promise<EmailSubmissionResponse> {
  try {
    const response = await fetch(
      'https://dklemailservice.onrender.com/api/contact-email',
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
        signal: AbortSignal.timeout(10000) // 10 second timeout
      }
    );

    const result = await response.json();

    if (!response.ok) {
      throw new Error(result.error || 'Unknown error');
    }

    return result;
  } catch (error) {
    if (error instanceof TypeError && error.message === 'Failed to fetch') {
      throw new Error('Geen verbinding met server. Controleer je internetverbinding.');
    }
    
    if (error.name === 'AbortError') {
      throw new Error('Verzoek duurde te lang. Probeer opnieuw.');
    }

    throw error;
  }
}
```

---

# 🔄 Next.js Server Actions

```typescript
// app/actions/contact.ts
'use server';

import { z } from 'zod';

const contactSchema = z.object({
  naam: z.string().min(1, 'Naam is verplicht'),
  email: z.string().email('Ongeldig email adres'),
  telefoon: z.string().optional(),
  bericht: z.string().min(10, 'Bericht moet minimaal 10 karakters bevatten'),
  privacy_akkoord: z.boolean().refine(val => val === true, {
    message: 'Je moet akkoord gaan met het privacybeleid'
  })
});

export async function submitContactForm(formData: FormData) {
  const data = {
    naam: formData.get('naam'),
    email: formData.get('email'),
    telefoon: formData.get('telefoon') || '',
    bericht: formData.get('bericht'),
    privacy_akkoord: formData.get('privacy_akkoord') === 'on'
  };

  // Validate
  const validated = contactSchema.parse(data);

  // Submit to API
  const response = await fetch(
    'https://dklemailservice.onrender.com/api/contact-email',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(validated)
    }
  );

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to submit');
  }

  return await response.json();
}
```

---

# 🎯 Best Practices

## Client-Side Validation
```typescript
const validateEmail = (email: string): boolean => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
};

const validatePhone = (phone: string): boolean => {
  // Dutch phone number pattern
  const re = /^((\+|00)?31)?[0-9]{9,10}$/;
  return phone === '' || re.test(phone.replace(/[\s-]/g, ''));
};
```

## Form State Management
```typescript
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const contactSchema = z.object({
  naam: z.string().min(1, 'Naam is verplicht'),
  email: z.string().email('Ongeldig email adres'),
  bericht: z.string().min(10, 'Bericht te kort'),
  privacy_akkoord: z.boolean().refine(val => val, 'Verplicht veld')
});

type ContactFormData = z.infer<typeof contactSchema>;

const ContactFormWithValidation = () => {
  const { register, handleSubmit, formState: { errors } } = useForm<ContactFormData>({
    resolver: zodResolver(contactSchema)
  });

  const onSubmit = async (data: ContactFormData) => {
    // Submit logic
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* Form fields with error handling */}
    </form>
  );
};
```

## Success Feedback
```typescript
const [submitted, setSubmitted] = useState(false);

if (submitted) {
  return (
    <div className="success-card">
      <svg className="checkmark" viewBox="0 0 52 52">
        <circle className="checkmark-circle" cx="26" cy="26" r="25"/>
        <path className="checkmark-check" d="M14 27l7 7 16-16"/>
      </svg>
      <h2>Verzonden!</h2>
      <p>Je ontvangt een bevestiging per email</p>
    </div>
  );
}
```

---

# 🔒 Security Considerations

1. **Nooit gevoelige data in localStorage** bewaren (behalve tokens)
2. **Sanitize user input** client-side (XSS prevention)
3. **HTTPS alleen** in production
4. **CSRF protection** via SameSite cookies
5. **Rate limit awareness** - Implementeer retry logic
6. **API key veilig bewaren** (alleen server-side voor WFC)

---

# Vragen?

Voor vragen of problemen:
- Contact handlers: [`handlers/contact_handler.go`](../handlers/contact_handler.go)
- Aanmelding handlers: [`handlers/aanmelding_handler.go`](../handlers/aanmelding_handler.go)
- Email handler: [`handlers/email_handler.go`](../handlers/email_handler.go)
- Mail handler: [`handlers/mail_handler.go`](../handlers/mail_handler.go)
- WFC handler: [`handlers/wfc_order_handler.go`](../handlers/wfc_order_handler.go)