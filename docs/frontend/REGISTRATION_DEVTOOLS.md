# Frontend DevTools - Registratie API

## 📋 Overzicht

Deze gids laat zien hoe je de registratie API test vanuit de frontend met DevTools. De DKL Email Service biedt een **duaal registratiesysteem** aan met twee types accounts:

1. **Full Account** - Met app toegang en gebruikersaccount
2. **Temporary Account** - Alleen voor het evenementjaar

## 🔗 API Endpoints

### Base URL
```
http://localhost:8080/api/public
```

### Beschikbare Endpoints

| Endpoint | Methode | Beschrijving |
|----------|---------|--------------|
| `/aanmelden` | POST | Registreer nieuwe participant |
| `/upgrade-to-full-account` | POST | Upgrade temporary naar full account |
| `/events/active` | GET | Haal actief evenement op |

---

## 🚀 1. Nieuwe Registratie - Full Account

### Request

```javascript
// Open DevTools Console (F12) en voer uit:
fetch('http://localhost:8080/api/public/aanmelden', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    // Persoonlijke gegevens
    naam: "Jeffrey Test",
    email: "laventejeffrey@gmail.com",
    telefoon: "0612345678",
    
    // Event keuzes
    rol: "Deelnemer",              // Opties: "Deelnemer", "Begeleider", "Vrijwilliger"
    afstand: "10 KM",              // Opties: "2.5 KM", "6 KM", "10 KM", "15 KM"
    ondersteuning: "Nee",          // Opties: "Ja", "Nee", "Anders"
    bijzonderheden: "",            // Verplicht als ondersteuning = "Ja" of "Anders"
    
    // Account type - BELANGRIJK!
    want_account: true,            // true = Full Account, false = Temporary
    wachtwoord: "TestWachtwoord123!", // Verplicht voor Full Account (min 8 chars)
    
    // Voorwaarden
    terms: true,
    
    // Test mode
    test_mode: true
  })
})
.then(res => res.json())
.then(data => console.log('✅ Registratie succesvol:', data))
.catch(err => console.error('❌ Fout:', err));
```

### Verwachte Response (Success)

```json
{
  "success": true,
  "message": "Je bent succesvol ingeschreven met een volledig account! Je hebt nu toegang tot de DKL Step App.",
  "participant_id": "uuid-hier",
  "registration_id": "uuid-hier",
  "account_type": "full",
  "has_app_access": true,
  "gebruiker_id": "uuid-hier",
  "event_name": "De Koninklijke Loop",
  "event_date": "15-06-2026"
}
```

---

## 🔄 2. Nieuwe Registratie - Temporary Account

### Request

```javascript
fetch('http://localhost:8080/api/public/aanmelden', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    naam: "Temporary User",
    email: "temp@example.com",
    telefoon: "",
    rol: "Deelnemer",
    afstand: "6 KM",
    ondersteuning: "Nee",
    bijzonderheden: "",
    
    // Temporary Account - GEEN wachtwoord nodig
    want_account: false,  // false = Temporary (alleen voor dit jaar)
    
    terms: true,
    test_mode: true
  })
})
.then(res => res.json())
.then(data => console.log('✅ Temp registratie succesvol:', data))
.catch(err => console.error('❌ Fout:', err));
```

### Verwachte Response (Success)

```json
{
  "success": true,
  "message": "Je bent succesvol ingeschreven voor 2026! Je registratie is geldig voor dit evenementjaar.",
  "participant_id": "uuid-hier",
  "registration_id": "uuid-hier",
  "account_type": "temporary",
  "has_app_access": false,
  "event_name": "De Koninklijke Loop",
  "event_date": "15-06-2026"
}
```

---

## ⬆️ 3. Upgrade naar Full Account

Upgrade een bestaand temporary account naar een full account.

### Request

```javascript
fetch('http://localhost:8080/api/public/upgrade-to-full-account', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: "temp@example.com",      // Email van bestaand temporary account
    wachtwoord: "NieuwWachtwoord123!" // Min 8 karakters
  })
})
.then(res => res.json())
.then(data => console.log('✅ Upgrade succesvol:', data))
.catch(err => console.error('❌ Fout:', err));
```

### Verwachte Response (Success)

```json
{
  "success": true,
  "message": "Je account is geüpgraded! Je hebt nu toegang tot de DKL Step App.",
  "gebruiker_id": "uuid-hier",
  "participant_id": "uuid-hier",
  "has_app_access": true
}
```

---

## 📊 4. Haal Actief Evenement Op

```javascript
fetch('http://localhost:8080/api/public/events/active')
  .then(res => res.json())
  .then(data => console.log('📅 Actief event:', data))
  .catch(err => console.error('❌ Fout:', err));
```

---

## ⚠️ Error Responses

### Validatie Fout

```json
{
  "error": "naam is verplicht",
  "code": "VALIDATION_ERROR"
}
```

### Email Bestaat Al (Full Account)

```json
{
  "error": "Er bestaat al een account met dit e-mailadres",
  "code": "EMAIL_EXISTS"
}
```

### Al Ingeschreven (Temporary - dit jaar)

```json
{
  "error": "Je bent al ingeschreven voor dit jaar met dit e-mailadres",
  "code": "ALREADY_REGISTERED"
}
```

### Database Fout

```json
{
  "error": "Registratie mislukt: kon event registratie niet aanmaken: ...",
  "code": "REGISTRATION_FAILED"
}
```

---

## 🧪 Test Scenario's

### Scenario 1: Full Account Registratie

```javascript
// Test volledige registratie flow
async function testFullAccountRegistration() {
  const testData = {
    naam: "Test User " + Date.now(),
    email: `test${Date.now()}@example.com`,
    telefoon: "0612345678",
    rol: "Deelnemer",
    afstand: "10 KM",
    ondersteuning: "Nee",
    bijzonderheden: "",
    want_account: true,
    wachtwoord: "TestPass123!",
    terms: true,
    test_mode: true
  };

  try {
    const response = await fetch('http://localhost:8080/api/public/aanmelden', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testData)
    });
    
    const result = await response.json();
    
    if (result.success) {
      console.log('✅ SUCCES - Full Account:', result);
      console.log('   Participant ID:', result.participant_id);
      console.log('   Gebruiker ID:', result.gebruiker_id);
      console.log('   App Access:', result.has_app_access);
    } else {
      console.error('❌ FOUT:', result);
    }
  } catch (error) {
    console.error('❌ Network Error:', error);
  }
}

testFullAccountRegistration();
```

### Scenario 2: Temporary Account Registratie

```javascript
async function testTemporaryAccountRegistration() {
  const testData = {
    naam: "Temp User " + Date.now(),
    email: `temp${Date.now()}@example.com`,
    telefoon: "",
    rol: "Begeleider",
    afstand: "15 KM",
    ondersteuning: "Anders",
    bijzonderheden: "Ik help met begeleiding",
    want_account: false,  // Temporary!
    terms: true,
    test_mode: true
  };

  try {
    const response = await fetch('http://localhost:8080/api/public/aanmelden', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(testData)
    });
    
    const result = await response.json();
    
    if (result.success) {
      console.log('✅ SUCCES - Temporary Account:', result);
      console.log('   Account Type:', result.account_type);
      console.log('   App Access:', result.has_app_access);
    } else {
      console.error('❌ FOUT:', result);
    }
  } catch (error) {
    console.error('❌ Network Error:', error);
  }
}

testTemporaryAccountRegistration();
```

### Scenario 3: Upgrade Flow

```javascript
async function testUpgradeFlow() {
  // Stap 1: Maak eerst een temporary account
  const tempEmail = `temp${Date.now()}@example.com`;
  
  const tempData = {
    naam: "To Upgrade",
    email: tempEmail,
    telefoon: "",
    rol: "Deelnemer",
    afstand: "6 KM",
    ondersteuning: "Nee",
    bijzonderheden: "",
    want_account: false,
    terms: true,
    test_mode: true
  };

  console.log('📝 Stap 1: Maak temporary account...');
  const tempResponse = await fetch('http://localhost:8080/api/public/aanmelden', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(tempData)
  });
  
  const tempResult = await tempResponse.json();
  console.log('✅ Temporary gemaakt:', tempResult);

  // Stap 2: Upgrade naar full account
  console.log('⬆️ Stap 2: Upgrade naar full account...');
  await new Promise(resolve => setTimeout(resolve, 1000)); // Wacht 1 sec
  
  const upgradeData = {
    email: tempEmail,
    wachtwoord: "UpgradePass123!"
  };

  const upgradeResponse = await fetch('http://localhost:8080/api/public/upgrade-to-full-account', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(upgradeData)
  });
  
  const upgradeResult = await upgradeResponse.json();
  console.log('✅ Upgrade compleet:', upgradeResult);
}

testUpgradeFlow();
```

---

## 📝 Validatie Regels

### Verplichte Velden

| Veld | Verplicht Voor | Voorwaarde |
|------|----------------|------------|
| `naam` | Alle | Altijd |
| `email` | Alle | Altijd, geldig email formaat |
| `telefoon` | Begeleider, Vrijwilliger | Alleen voor deze rollen |
| `rol` | Alle | Moet "Deelnemer", "Begeleider" of "Vrijwilliger" zijn |
| `afstand` | Alle | Moet "2.5 KM", "6 KM", "10 KM" of "15 KM" zijn |
| `ondersteuning` | Alle | Moet "Ja", "Nee" of "Anders" zijn |
| `bijzonderheden` | Bij ondersteuning | Verplicht als ondersteuning = "Ja" of "Anders" |
| `wachtwoord` | Full Account | Verplicht als `want_account` = true, min 8 chars |
| `terms` | Alle | Moet `true` zijn |

### Backend Validatie

De backend valideert:
- Email format en duplicaten
- Wachtwoord sterkte (min 8 karakters voor full accounts)
- Telefoon voor Begeleider/Vrijwilliger rollen
- Bijzonderheden bij ondersteuning = "Ja" of "Anders"
- Terms acceptance

---

## 🔍 Database Status Checken

```javascript
// Bekijk welke statussen bestaan
async function checkRegistrationStatuses() {
  // Dit zou een admin endpoint moeten zijn, maar voor debugging:
  console.log('Beschikbare statussen:');
  console.log('- registered (standaard na nieuwe registratie)');
  console.log('- nieuw (nieuwe aanmelding)');
  console.log('- confirmed (bevestigd door admin)');
  console.log('- cancelled (geannuleerd)');
  console.log('- waiting_list (op wachtlijst)');
  console.log('- attended (aanwezig geweest)');
  console.log('- no_show (niet verschenen)');
}

checkRegistrationStatuses();
```

---

## 🎨 Frontend Integratie Voorbeeld

### React/Vue Formulier

```javascript
const handleRegistration = async (formData) => {
  try {
    const response = await fetch('http://localhost:8080/api/public/aanmelden', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        naam: formData.naam,
        email: formData.email,
        telefoon: formData.telefoon || "",
        rol: formData.rol,
        afstand: formData.afstand,
        ondersteuning: formData.ondersteuning,
        bijzonderheden: formData.bijzonderheden || "",
        want_account: formData.createAccount, // Checkbox waarde
        wachtwoord: formData.createAccount ? formData.wachtwoord : undefined,
        terms: formData.terms,
        test_mode: import.meta.env.DEV // Automatisch in dev mode
      })
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registratie mislukt');
    }

    const result = await response.json();
    
    // Toon success message
    if (result.account_type === 'full') {
      alert(`Succesvol geregistreerd! Je kunt nu inloggen in de app met ${result.email}`);
    } else {
      alert(`Succesvol ingeschreven voor ${new Date().getFullYear()}!`);
    }
    
    return result;
    
  } catch (error) {
    console.error('Registratie fout:', error);
    alert(error.message);
    throw error;
  }
};
```

---

## 🐛 Troubleshooting

### CORS Errors

Als je CORS fouten krijgt:
```javascript
// Check of je origin in de allowed origins staat
console.log('Allowed Origins:', [
  'http://localhost:3000',
  'http://localhost:5173',
  'http://localhost:8082'
]);
```

### Network Errors

```javascript
// Test of de API bereikbaar is
fetch('http://localhost:8080/api/public/events/active')
  .then(res => res.ok ? console.log('✅ API is online') : console.error('❌ API error'))
  .catch(() => console.error('❌ Kan API niet bereiken'));
```

### Rollback Gebeurt Bij Fouten

Sinds de fix zijn rollbacks geïmplementeerd. Als je een fout krijgt:
- Bij Full Account: Gebruiker EN Participant worden verwijderd
- Bij Temporary Account: Participant wordt verwijderd
- Geen "halve" registraties meer in de database

---

## 📞 Support

Voor vragen over de API:
- Backend logs: `docker-compose logs app --tail 50 --follow`
- Database check: Zie [`docs/architecture/DATABASE.md`](../architecture/DATABASE.md)
- API documentatie: [`docs/api/`](../api/)